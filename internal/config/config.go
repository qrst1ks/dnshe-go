package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const DefaultDNSHEBaseURL = "https://api005.dnshe.com/index.php?m=domain_hub"

type Config struct {
	IntervalSeconds int         `json:"interval_seconds"`
	TTL             int         `json:"ttl"`
	DNSHE           DNSHEConfig `json:"dnshe"`
	IPv6            IPConfig    `json:"ipv6"`
}

type DNSHEConfig struct {
	APIKey     string `json:"api_key"`
	APISecret  string `json:"api_secret"`
	APIBaseURL string `json:"api_base_url"`
}

type IPConfig struct {
	Enable    bool     `json:"enable"`
	Source    string   `json:"source"`
	URLs      []string `json:"urls"`
	Interface string   `json:"interface"`
	Command   string   `json:"command"`
	Domains   []string `json:"domains"`
}

type Store struct {
	path string
	mu   sync.RWMutex
	cfg  Config
}

func Default() Config {
	return Config{
		IntervalSeconds: 300,
		TTL:             600,
		DNSHE: DNSHEConfig{
			APIBaseURL: DefaultDNSHEBaseURL,
		},
		IPv6: IPConfig{
			Enable:  true,
			Source:  "url",
			Command: DefaultIPv6Command(),
			URLs: []string{
				"https://api-ipv6.ip.sb/ip",
				"https://ipv6.icanhazip.com",
			},
		},
	}
}

func NewStore(path string) (*Store, error) {
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}
	return &Store{path: path, cfg: cfg}, nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Get() Config {
	s.mu.RLock()
	cfg := s.cfg
	s.mu.RUnlock()
	cfg.ApplyEnv()
	cfg.Normalize()
	return cfg
}

func (s *Store) Save(cfg Config) error {
	cfg.Normalize()
	if err := Save(s.path, cfg); err != nil {
		return err
	}
	s.mu.Lock()
	s.cfg = cfg
	s.mu.Unlock()
	return nil
}

func Load(path string) (Config, error) {
	cfg := Default()
	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg.Normalize()
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(body, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	cfg.Normalize()
	return cfg, nil
}

func Save(path string, cfg Config) error {
	cfg.Normalize()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, body, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (c *Config) Normalize() {
	if c.IntervalSeconds <= 0 {
		c.IntervalSeconds = 300
	}
	if c.TTL <= 0 {
		c.TTL = 600
	}
	if strings.TrimSpace(c.DNSHE.APIBaseURL) == "" {
		c.DNSHE.APIBaseURL = DefaultDNSHEBaseURL
	}
	c.DNSHE.APIKey = strings.TrimSpace(c.DNSHE.APIKey)
	c.DNSHE.APISecret = strings.TrimSpace(c.DNSHE.APISecret)
	c.DNSHE.APIBaseURL = strings.TrimSpace(c.DNSHE.APIBaseURL)
	c.IPv6.Normalize()
}

func (c *Config) ApplyEnv() {
	if v := strings.TrimSpace(os.Getenv("DNSHE_API_KEY")); v != "" {
		c.DNSHE.APIKey = v
	}
	if v := strings.TrimSpace(os.Getenv("DNSHE_API_SECRET")); v != "" {
		c.DNSHE.APISecret = v
	}
	if v := strings.TrimSpace(os.Getenv("DNSHE_API_BASE_URL")); v != "" {
		c.DNSHE.APIBaseURL = v
	}
}

func (c IPConfig) DomainsClean() []string {
	return CleanStringList(c.Domains)
}

func (c *IPConfig) Normalize() {
	c.Source = strings.TrimSpace(strings.ToLower(c.Source))
	switch c.Source {
	case "url", "interface", "cmd":
	default:
		c.Source = "url"
	}
	c.Interface = strings.TrimSpace(c.Interface)
	c.Command = strings.TrimSpace(c.Command)
	if c.Command == "" {
		c.Command = DefaultIPv6Command()
	}
	c.URLs = CleanStringList(c.URLs)
	c.Domains = CleanStringList(c.Domains)
}

func DefaultIPv6Command() string {
	switch runtime.GOOS {
	case "windows":
		return `(Get-NetIPAddress -AddressFamily IPv6 | Where-Object {$_.AddressState -eq "Preferred" -and $_.IPAddress -notlike "fe80*" -and $_.IPAddress -ne "::1"} | Select-Object -First 1 -ExpandProperty IPAddress)`
	case "linux":
		return `ip -6 addr show scope global | awk '/inet6/{print $2; exit}' | cut -d/ -f1`
	default:
		return `ifconfig | awk '/inet6 / && $2 !~ /^fe80/ && $2 != "::1" {print $2; exit}'`
	}
}

func CleanStringList(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		for _, line := range strings.Split(value, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			out = append(out, line)
		}
	}
	return out
}

func IsMaskedSecret(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || strings.Contains(value, "***")
}

func MaskSecret(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	if len(value) <= 8 {
		return value[:1] + "***" + value[len(value)-1:]
	}
	return value[:3] + "***" + value[len(value)-3:]
}
