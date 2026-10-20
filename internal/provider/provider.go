package provider

import "context"

const (
	StatusUnchanged = "unchanged"
	StatusUpdated   = "updated"
	StatusCreated   = "created"
)

type Result struct {
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"`
	RecordID   string `json:"record_id,omitempty"`
	OldIP      string `json:"old_ip,omitempty"`
	NewIP      string `json:"new_ip"`
	Status     string `json:"status"`
}

type Updater interface {
	EnsureRecord(ctx context.Context, domain string, recordType string, ip string, ttl int) (Result, error)
}
