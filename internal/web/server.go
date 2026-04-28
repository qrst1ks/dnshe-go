package web

import (
	"context"
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

type Server struct {
	http.Server
	store  *config.Store
	syncer *ddns.Syncer
	logs   *logbuf.Buffer
	tmpl   *template.Template
}

type publicConfig struct {
	IntervalSeconds     int             `json:"interval_seconds"`
	TTL                 int             `json:"ttl"`
	DNSHE               publicDNSHE     `json:"dnshe"`
	IPv6                config.IPConfig `json:"ipv6"`
	APIKeyConfigured    bool            `json:"api_key_configured"`
	APISecretConfigured bool            `json:"api_secret_configured"`
}

type publicDNSHE struct {
	APIKeyMasked    string `json:"api_key_masked"`
	APISecretMasked string `json:"api_secret_masked"`
	APIBaseURL      string `json:"api_base_url"`
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
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/status", s.handleStatus)
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
			APIKeyMasked:    config.MaskSecret(cfg.DNSHE.APIKey),
			APISecretMasked: config.MaskSecret(cfg.DNSHE.APISecret),
			APIBaseURL:      cfg.DNSHE.APIBaseURL,
		},
		IPv6:                cfg.IPv6,
		APIKeyConfigured:    strings.TrimSpace(cfg.DNSHE.APIKey) != "",
		APISecretConfigured: strings.TrimSpace(cfg.DNSHE.APISecret) != "",
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func validateConfig(cfg config.Config) []string {
	var problems []string
	if strings.TrimSpace(cfg.DNSHE.APIKey) == "" {
		problems = append(problems, "DNSHE API Key 不能为空")
	}
	if strings.TrimSpace(cfg.DNSHE.APISecret) == "" {
		problems = append(problems, "DNSHE API Secret 不能为空")
	}
	if !cfg.IPv6.Enable {
		problems = append(problems, "IPv6 同步必须启用")
	}
	validateIPConfig := func(name string, ipCfg config.IPConfig) {
		if !ipCfg.Enable {
			return
		}
		if len(ipCfg.DomainsClean()) == 0 {
			problems = append(problems, name+" 域名不能为空")
		}
		switch ipCfg.Source {
		case "url":
			if len(ipCfg.URLs) == 0 {
				problems = append(problems, name+" 使用 URL 来源时 URL 不能为空")
			}
		case "cmd":
			if strings.TrimSpace(ipCfg.Command) == "" {
				problems = append(problems, name+" 使用命令来源时命令不能为空")
			}
		case "interface":
			// 网卡名允许为空，表示自动选择活动网卡。
		}
	}
	validateIPConfig("IPv6", cfg.IPv6)
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
