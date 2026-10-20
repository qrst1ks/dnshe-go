package web

import (
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/ddns"
	"github.com/qrst1ks/dnshe-go/internal/logbuf"
)

//go:embed assets/*
var webAssets embed.FS

type Server struct {
	http.Server
	store  *config.Store
	syncer *ddns.Syncer
	logs   *logbuf.Buffer
	tmpl   *template.Template
}

type publicConfig struct {
	IntervalSeconds int              `json:"interval_seconds"`
	TTL             int              `json:"ttl"`
	DNSHE           publicDNSHE      `json:"dnshe"`
	Cloudflare      publicCloudflare `json:"cloudflare"`
	IPv6            publicIPConfig   `json:"ipv6"`
}

type publicDNSHE struct {
	Enable          bool     `json:"enable"`
	APIKeyMasked    string   `json:"api_key_masked"`
	APISecretMasked string   `json:"api_secret_masked"`
	APIBaseURL      string   `json:"api_base_url"`
	Domains         []string `json:"domains"`
}

type publicCloudflare struct {
	Enable         bool     `json:"enable"`
	APITokenMasked string   `json:"api_token_masked"`
	APIBaseURL     string   `json:"api_base_url"`
	Domains        []string `json:"domains"`
}

type publicIPConfig struct {
	Source    string   `json:"source"`
	URLs      []string `json:"urls"`
	Interface string   `json:"interface"`
	Command   string   `json:"command"`
}

type netInterfaceInfo struct {
	Name      string   `json:"name"`
	Label     string   `json:"label"`
	Addresses []string `json:"addresses"`
}

func NewServer(listen string, store *config.Store, syncer *ddns.Syncer, logs *logbuf.Buffer) *Server {
	s := &Server{
		store:  store,
		syncer: syncer,
		logs:   logs,
		tmpl:   template.Must(template.New("index").Parse(indexHTML)),
	}
	mux := http.NewServeMux()
	assets := http.FileServer(http.FS(webAssets))
	mux.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		assets.ServeHTTP(w, r)
	})
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/secret", s.handleSecret)
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/interfaces", s.handleInterfaces)
	mux.HandleFunc("/api/run", s.handleRun)
	mux.HandleFunc("/api/logs/clear", s.handleClearLogs)
	s.Server = http.Server{
		Addr:              listen,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return s
}

func (s *Server) handleInterfaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "msg": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"interfaces": listIPv6Interfaces(),
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.tmpl.Execute(w, map[string]any{"Title": "dnshe-go"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":          true,
		"config_path": s.store.Path(),
		"config":      makePublicConfig(s.store.Get()),
		"sync":        s.syncer.Snapshot(),
		"logs":        s.logs.List(),
	})
}

func (s *Server) handleSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "msg": "method not allowed"})
		return
	}
	cfg := s.store.Get()
	var value string
	switch r.URL.Query().Get("id") {
	case "dnshe-api-key":
		value = cfg.DNSHE.APIKey
	case "dnshe-api-secret":
		value = cfg.DNSHE.APISecret
	case "cloudflare-api-token":
		value = cfg.Cloudflare.APIToken
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"ok": false, "msg": "unknown secret"})
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "value": value})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "msg": "method not allowed"})
		return
	}
	var next config.Config
	if err := json.NewDecoder(r.Body).Decode(&next); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "msg": "invalid json body"})
		return
	}
	current := s.store.Get()
	if config.IsMaskedSecret(next.DNSHE.APIKey) {
		next.DNSHE.APIKey = current.DNSHE.APIKey
	}
	if config.IsMaskedSecret(next.DNSHE.APISecret) {
		next.DNSHE.APISecret = current.DNSHE.APISecret
	}
	if config.IsMaskedSecret(next.Cloudflare.APIToken) {
		next.Cloudflare.APIToken = current.Cloudflare.APIToken
	}
	next.Normalize()
	if problems := validateConfig(next); len(problems) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "msg": strings.Join(problems, "；"), "problems": problems})
		return
	}
	if err := s.store.Save(next); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "msg": err.Error()})
		return
	}
	s.logs.Addf("INFO", "config saved")
	go s.syncer.RunOnce(context.Background(), true)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "config": makePublicConfig(s.store.Get())})
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "msg": "method not allowed"})
		return
	}
	go s.syncer.RunOnce(context.Background(), true)
	writeJSON(w, http.StatusAccepted, map[string]any{"ok": true})
}

func (s *Server) handleClearLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "msg": "method not allowed"})
		return
	}
	s.logs.Clear()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func makePublicConfig(cfg config.Config) publicConfig {
	return publicConfig{
		IntervalSeconds: cfg.IntervalSeconds,
		TTL:             cfg.TTL,
		DNSHE: publicDNSHE{
			Enable:          cfg.DNSHE.Enable,
			APIKeyMasked:    config.MaskSecret(cfg.DNSHE.APIKey),
			APISecretMasked: config.MaskSecret(cfg.DNSHE.APISecret),
			APIBaseURL:      cfg.DNSHE.APIBaseURL,
			Domains:         cfg.DNSHE.DomainsClean(),
		},
		Cloudflare: publicCloudflare{
			Enable:         cfg.Cloudflare.Enable,
			APITokenMasked: config.MaskSecret(cfg.Cloudflare.APIToken),
			APIBaseURL:     cfg.Cloudflare.APIBaseURL,
			Domains:        cfg.Cloudflare.DomainsClean(),
		},
		IPv6: publicIPConfig{
			Source:    cfg.IPv6.Source,
			URLs:      cfg.IPv6.URLs,
			Interface: cfg.IPv6.Interface,
			Command:   cfg.IPv6.Command,
		},
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func validateConfig(cfg config.Config) []string {
	var problems []string
	if !cfg.DNSHE.Enable && !cfg.Cloudflare.Enable {
		problems = append(problems, "至少需要启用一个 DNS 服务商")
	}
	if cfg.DNSHE.Enable {
		if strings.TrimSpace(cfg.DNSHE.APIKey) == "" {
			problems = append(problems, "DNSHE API Key 不能为空")
		}
		if strings.TrimSpace(cfg.DNSHE.APISecret) == "" {
			problems = append(problems, "DNSHE API Secret 不能为空")
		}
		if len(cfg.DNSHE.DomainsClean()) == 0 {
			problems = append(problems, "DNSHE 域名不能为空")
		}
	}
	if cfg.Cloudflare.Enable {
		if strings.TrimSpace(cfg.Cloudflare.APIToken) == "" {
			problems = append(problems, "Cloudflare API Token 不能为空")
		}
		if len(cfg.Cloudflare.DomainsClean()) == 0 {
			problems = append(problems, "Cloudflare 域名不能为空")
		}
	}
	switch cfg.IPv6.Source {
	case "url":
		if len(cfg.IPv6.URLs) == 0 {
			problems = append(problems, "IPv6 使用 URL 来源时 URL 不能为空")
		}
	case "cmd":
		if strings.TrimSpace(cfg.IPv6.Command) == "" {
			problems = append(problems, "IPv6 使用命令来源时命令不能为空")
		}
	case "interface":
		// 网卡名允许为空，表示自动选择活动网卡。
	}
	return problems
}

func listIPv6Interfaces() []netInterfaceInfo {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	result := make([]netInterfaceInfo, 0, len(ifaces))
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		ipv6 := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			ip := ipFromAddr(addr)
			if ip == nil || ip.To4() != nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			ipv6 = append(ipv6, ip.String())
		}
		if len(ipv6) == 0 {
			continue
		}
		label := iface.Name + " - " + ipv6[0]
		result = append(result, netInterfaceInfo{
			Name:      iface.Name,
			Label:     label,
			Addresses: ipv6,
		})
	}
	return result
}

func ipFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}
