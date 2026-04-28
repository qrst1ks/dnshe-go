package dnshe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qrst1ks/dnshe-go/internal/config"
)

func TestEnsureRecordUpdatesChangedRecord(t *testing.T) {
	var updateCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "key" || r.Header.Get("X-API-Secret") != "secret" {
			t.Fatalf("missing DNSHE headers")
		}
		switch r.URL.Query().Get("endpoint") + ":" + r.URL.Query().Get("action") {
		case "subdomains:list":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"subdomains": []map[string]any{{"id": 7, "full_domain": "home.example.com"}},
			})
		case "dns_records:list":
			if r.URL.Query().Get("subdomain_id") != "7" {
				t.Fatalf("subdomain_id = %s", r.URL.Query().Get("subdomain_id"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"records": []map[string]any{{"id": 11, "type": "AAAA", "content": "2001:db8::1"}},
			})
		case "dns_records:update":
			updateCalled = true
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["content"] != "2001:db8::2" {
				t.Fatalf("content = %#v", payload["content"])
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
		default:
			t.Fatalf("unexpected request: %s", r.URL.String())
		}
	}))
	defer server.Close()

	client := NewClient(config.DNSHEConfig{
		APIKey:     "key",
		APISecret:  "secret",
		APIBaseURL: server.URL + "/index.php?m=domain_hub",
	})
	result, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "updated" {
		t.Fatalf("status = %s", result.Status)
	}
	if !updateCalled {
		t.Fatal("expected update call")
	}
}

func TestEnsureRecordSkipsUnchangedRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("endpoint") + ":" + r.URL.Query().Get("action") {
		case "subdomains:list":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"subdomains": []map[string]any{{"id": 7, "full_domain": "home.example.com"}},
			})
		case "dns_records:list":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"records": []map[string]any{{"id": 11, "type": "AAAA", "content": "2001:db8::2"}},
			})
		case "dns_records:update":
			t.Fatal("update should not be called")
		}
	}))
	defer server.Close()

	client := NewClient(config.DNSHEConfig{APIBaseURL: server.URL + "/index.php?m=domain_hub"})
	result, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "unchanged" {
		t.Fatalf("status = %s", result.Status)
	}
}
