package ddns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/dnshe"
	"github.com/qrst1ks/dnshe-go/internal/logbuf"
)

type RecordResult struct {
	Time       string `json:"time"`
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"`
	IP         string `json:"ip"`
	OldIP      string `json:"old_ip,omitempty"`
	RecordID   int    `json:"record_id,omitempty"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type Status struct {
	Running          bool           `json:"running"`
	LastRunStartedAt string         `json:"last_run_started_at,omitempty"`
	LastRunEndedAt   string         `json:"last_run_ended_at,omitempty"`
	LastDurationMs   int64          `json:"last_duration_ms,omitempty"`
	LastError        string         `json:"last_error,omitempty"`
	CurrentIPv6      string         `json:"current_ipv6,omitempty"`
	Results          []RecordResult `json:"results"`
}

type Syncer struct {
	store *config.Store
	logs  *logbuf.Buffer

	mu     sync.Mutex
	status Status
	cache  map[string]string
}

func NewSyncer(store *config.Store, logs *logbuf.Buffer) *Syncer {
	return &Syncer{
		store: store,
		logs:  logs,
		cache: map[string]string{},
	}
}

func (s *Syncer) Snapshot() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return cloneStatus(s.status)
}

func (s *Syncer) RunLoop(ctx context.Context, interval time.Duration) {
	s.RunOnce(ctx, true)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logs.Addf("INFO", "dnshe-go stopped")
			return
		case <-ticker.C:
			s.RunOnce(ctx, false)
		}
	}
}

func (s *Syncer) RunOnce(ctx context.Context, force bool) Status {
	if !s.tryStart() {
		s.logs.Addf("WARN", "sync already running")
		return s.Snapshot()
	}
	started := time.Now()
	status := Status{
		Running:          true,
		LastRunStartedAt: started.Format(time.RFC3339),
	}
	s.setStatus(status)
	defer func() {
		ended := time.Now()
		status.Running = false
		status.LastRunEndedAt = ended.Format(time.RFC3339)
		status.LastDurationMs = ended.Sub(started).Milliseconds()
		s.setStatus(status)
		s.finish()
	}()

	cfg := s.store.Get()
	if cfg.DNSHE.APIKey == "" || cfg.DNSHE.APISecret == "" {
		status.LastError = "DNSHE API key/secret not configured"
		s.logs.Addf("ERROR", "%s", status.LastError)
		return status
	}

	client := dnshe.NewClient(cfg.DNSHE)
	s.logs.Addf("INFO", "sync started")
	s.runIPv6(ctx, client, cfg, force, &status)
	if status.LastError == "" {
		s.logs.Addf("INFO", "sync finished in %dms", time.Since(started).Milliseconds())
	}
	return status
}

func (s *Syncer) runIPv6(ctx context.Context, client *dnshe.Client, cfg config.Config, force bool, status *Status) {
	if !cfg.IPv6.Enable {
		return
	}
	const recordType = "AAAA"
	domains := cfg.IPv6.DomainsClean()
	if len(domains) == 0 {
		return
	}

	ip, err := ResolveIPv6(ctx, cfg.IPv6)
	if err != nil {
		msg := fmt.Sprintf("resolve IPv6 failed: %v", err)
		status.LastError = msg
		s.logs.Addf("ERROR", "%s", msg)
		return
	}
	status.CurrentIPv6 = ip

	if !force && s.cache[recordType] == ip {
		s.logs.Addf("INFO", "IPv6 unchanged: %s", maskIP(ip))
		return
	}

	familyOK := true
	for _, domain := range domains {
		result, err := client.EnsureRecord(ctx, domain, recordType, ip, cfg.TTL)
		record := RecordResult{
			Time:       time.Now().Format(time.RFC3339),
			Domain:     domain,
			RecordType: recordType,
			IP:         ip,
			OldIP:      result.OldIP,
			RecordID:   result.RecordID,
			Status:     result.Status,
		}
		if err != nil {
			familyOK = false
			record.Status = "failed"
			record.Message = err.Error()
			status.LastError = fmt.Sprintf("%s %s failed: %v", domain, recordType, err)
			s.logs.Addf("ERROR", "update %s %s failed: %v", domain, recordType, err)
		} else if result.Status == "updated" {
			s.logs.Addf("INFO", "updated %s %s: %s -> %s", domain, recordType, maskIP(result.OldIP), maskIP(ip))
		} else {
			s.logs.Addf("INFO", "unchanged %s %s: %s", domain, recordType, maskIP(ip))
		}
		status.Results = append(status.Results, record)
	}
	if familyOK {
		s.cache[recordType] = ip
	}
}

func (s *Syncer) tryStart() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status.Running {
		return false
	}
	s.status.Running = true
	return true
}

func (s *Syncer) finish() {
	s.mu.Lock()
	s.status.Running = false
	s.mu.Unlock()
}

func (s *Syncer) setStatus(status Status) {
	s.mu.Lock()
	s.status = cloneStatus(status)
	s.mu.Unlock()
}

func cloneStatus(status Status) Status {
	status.Results = append([]RecordResult(nil), status.Results...)
	return status
}
