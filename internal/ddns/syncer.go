package ddns

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/cloudflare"
	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/dnshe"
	"github.com/qrst1ks/dnshe-go/internal/logbuf"
	"github.com/qrst1ks/dnshe-go/internal/provider"
)

type RecordResult struct {
	Time       string `json:"time"`
	Provider   string `json:"provider"`
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"`
	IP         string `json:"ip"`
	OldIP      string `json:"old_ip,omitempty"`
	RecordID   string `json:"record_id,omitempty"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type providerJob struct {
	name    string
	domains []string
	updater provider.Updater
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
	jobs, problems := buildProviderJobs(cfg)
	for _, problem := range problems {
		s.logs.Addf("ERROR", "%s", problem)
	}

	ip, err := ResolveIPv6(ctx, cfg.IPv6)
	if err != nil {
		msg := fmt.Sprintf("resolve IPv6 failed: %v", err)
		status.LastError = strings.Join(append(problems, msg), "; ")
		s.logs.Addf("ERROR", "%s", msg)
		return status
	}
	status.CurrentIPv6 = ip

	if len(jobs) == 0 {
		if len(problems) == 0 {
			s.logs.Addf("INFO", "no DNS provider enabled, IPv6 resolved without updates")
			return status
		}
		status.LastError = strings.Join(problems, "; ")
		return status
	}

	if !force && s.cache["AAAA"] == ip {
		s.logs.Addf("INFO", "IPv6 unchanged: %s", maskIP(ip))
		status.LastError = strings.Join(problems, "; ")
		return status
	}

	s.logs.Addf("INFO", "sync started")
	jobsOK := true
	for _, job := range jobs {
		ok, jobProblems := s.runIPv6(ctx, job, cfg.TTL, ip, &status)
		if !ok {
			jobsOK = false
			problems = append(problems, jobProblems...)
		}
	}
	if jobsOK {
		s.cache["AAAA"] = ip
	}
	if len(problems) == 0 {
		s.logs.Addf("INFO", "sync finished in %dms", time.Since(started).Milliseconds())
	} else {
		status.LastError = strings.Join(problems, "; ")
	}
	return status
}

func buildProviderJobs(cfg config.Config) ([]providerJob, []string) {
	var jobs []providerJob
	var problems []string

	if cfg.DNSHE.Enable {
		domains := cfg.DNSHE.DomainsClean()
		switch {
		case cfg.DNSHE.APIKey == "" || cfg.DNSHE.APISecret == "":
			problems = append(problems, "DNSHE API key/secret not configured")
		case len(domains) == 0:
			problems = append(problems, "DNSHE has no domains configured")
		default:
			jobs = append(jobs, providerJob{
				name:    "DNSHE",
				domains: domains,
				updater: dnshe.NewClient(cfg.DNSHE),
			})
		}
	}

	if cfg.Cloudflare.Enable {
		domains := cfg.Cloudflare.DomainsClean()
		switch {
		case cfg.Cloudflare.APIToken == "":
			problems = append(problems, "Cloudflare API token not configured")
		case len(domains) == 0:
			problems = append(problems, "Cloudflare has no domains configured")
		default:
			jobs = append(jobs, providerJob{
				name:    "Cloudflare",
				domains: domains,
				updater: cloudflare.NewClient(cfg.Cloudflare),
			})
		}
	}

	return jobs, problems
}

func (s *Syncer) runIPv6(ctx context.Context, job providerJob, ttl int, ip string, status *Status) (bool, []string) {
	const recordType = "AAAA"
	ok := true
	var problems []string

	for _, domain := range job.domains {
		result, err := job.updater.EnsureRecord(ctx, domain, recordType, ip, ttl)
		record := RecordResult{
			Time:       time.Now().Format(time.RFC3339),
			Provider:   job.name,
			Domain:     domain,
			RecordType: recordType,
			IP:         ip,
			OldIP:      result.OldIP,
			RecordID:   result.RecordID,
			Status:     result.Status,
		}
		if err != nil {
			ok = false
			record.Status = "failed"
			record.Message = err.Error()
			problems = append(problems, fmt.Sprintf("%s %s (%s) failed: %v", domain, recordType, job.name, err))
			s.logs.Addf("ERROR", "update %s %s via %s failed: %v", domain, recordType, job.name, err)
		} else {
			switch result.Status {
			case provider.StatusUpdated:
				s.logs.Addf("INFO", "updated %s %s via %s: %s -> %s", domain, recordType, job.name, maskIP(result.OldIP), maskIP(ip))
			case provider.StatusCreated:
				s.logs.Addf("INFO", "created %s %s via %s: %s", domain, recordType, job.name, maskIP(ip))
			default:
				s.logs.Addf("INFO", "unchanged %s %s via %s: %s", domain, recordType, job.name, maskIP(ip))
			}
		}
		status.Results = append(status.Results, record)
	}
	return ok, problems
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
