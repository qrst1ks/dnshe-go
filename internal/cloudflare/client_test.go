package cloudflare

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/provider"
)

func TestEnsureRecordUpdatesChangedRecord(t *testing.T) {
	var updateCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("authorization = %q", got)
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/zones":
			if r.URL.Query().Get("status") != "active" {
				t.Fatalf("status = %q", r.URL.Query().Get("status"))
			}
			if r.URL.Query().Get("name") == "home.example.com" {
				writeResult(w, []map[string]any{})
				return
			}
			if r.URL.Query().Get("name") != "example.com" {
				t.Fatalf("unexpected zone query: %s", r.URL.Query().Get("name"))
			}
			writeResult(w, []map[string]any{{"id": "zone-1", "name": "example.com"}})
		case r.Method == http.MethodGet && r.URL.Path == "/zones/zone-1/dns_records":
			if r.URL.Query().Get("type") != "AAAA" || r.URL.Query().Get("name") != "home.example.com" {
				t.Fatalf("unexpected record query: %s", r.URL.RawQuery)
			}
			writeResult(w, []map[string]any{{
				"id": "record-1", "type": "AAAA", "name": "home.example.com",
				"content": "2001:db8::1", "ttl": 600, "proxied": false,
			}})
		case r.Method == http.MethodPut && r.URL.Path == "/zones/zone-1/dns_records/record-1":
			updateCalled = true
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["content"] != "2001:db8::2" {
				t.Fatalf("content = %#v", payload["content"])
			}
			if payload["proxied"] != false {
				t.Fatalf("proxied = %#v", payload["proxied"])
			}
			writeResult(w, map[string]any{"id": "record-1"})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	client := NewClient(config.CloudflareConfig{APIToken: "token", APIBaseURL: server.URL})
	result, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != provider.StatusUpdated {
		t.Fatalf("status = %s", result.Status)
	}
	if result.RecordID != "record-1" {
		t.Fatalf("record id = %q", result.RecordID)
	}
	if result.OldIP != "2001:db8::1" {
		t.Fatalf("old ip = %q", result.OldIP)
	}
	if !updateCalled {
		t.Fatal("expected update call")
	}
}

func TestEnsureRecordCreatesMissingRecord(t *testing.T) {
	var createCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/zones":
			if r.URL.Query().Get("name") == "home.example.com" {
				writeResult(w, []map[string]any{})
				return
			}
			writeResult(w, []map[string]any{{"id": "zone-1", "name": "example.com"}})
		case r.Method == http.MethodGet && r.URL.Path == "/zones/zone-1/dns_records":
			writeResult(w, []map[string]any{})
		case r.Method == http.MethodPost && r.URL.Path == "/zones/zone-1/dns_records":
			createCalled = true
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["type"] != "AAAA" || payload["name"] != "home.example.com" || payload["content"] != "2001:db8::2" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
			if payload["proxied"] != false {
				t.Fatalf("proxied = %#v", payload["proxied"])
			}
			writeResult(w, map[string]any{"id": "record-2"})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	client := NewClient(config.CloudflareConfig{APIToken: "token", APIBaseURL: server.URL})
	result, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != provider.StatusCreated {
		t.Fatalf("status = %s", result.Status)
	}
	if result.RecordID != "record-2" {
		t.Fatalf("record id = %q", result.RecordID)
	}
	if !createCalled {
		t.Fatal("expected create call")
	}
}

func TestEnsureRecordSkipsEquivalentIP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/zones":
			if r.URL.Query().Get("name") == "home.example.com" {
				writeResult(w, []map[string]any{})
				return
			}
			writeResult(w, []map[string]any{{"id": "zone-1", "name": "example.com"}})
		case r.Method == http.MethodGet && r.URL.Path == "/zones/zone-1/dns_records":
			writeResult(w, []map[string]any{{
				"id": "record-1", "type": "AAAA", "name": "home.example.com",
				"content": "2001:0db8:0000:0000:0000:0000:0000:0002", "ttl": 600,
			}})
		case r.Method == http.MethodPut, r.Method == http.MethodPost:
			t.Fatal("record should not be updated")
		}
	}))
	defer server.Close()

	client := NewClient(config.CloudflareConfig{APIToken: "token", APIBaseURL: server.URL})
	result, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != provider.StatusUnchanged {
		t.Fatalf("status = %s", result.Status)
	}
}

func TestEnsureRecordReportsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"errors":  []map[string]any{{"code": 9109, "message": "Invalid access token"}},
			"result":  nil,
		})
	}))
	defer server.Close()

	client := NewClient(config.CloudflareConfig{APIToken: "bad", APIBaseURL: server.URL})
	if _, err := client.EnsureRecord(context.Background(), "home.example.com", "AAAA", "2001:db8::2", 600); err == nil {
		t.Fatal("expected error")
	}
}

func TestCloudflareTTL(t *testing.T) {
	cases := map[int]int{0: 60, 1: 1, 30: 60, 60: 60, 600: 600}
	for input, want := range cases {
		if got := cloudflareTTL(input); got != want {
			t.Fatalf("cloudflareTTL(%d) = %d, want %d", input, got, want)
		}
	}
}

func TestZoneCandidates(t *testing.T) {
	got := zoneCandidates("nas.home.example.com")
	want := []string{"nas.home.example.com", "home.example.com", "example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
	if got := zoneCandidates("example"); got != nil {
		t.Fatalf("got %#v, want nil", got)
	}
}

func writeResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"errors":  []any{},
		"result":  result,
	})
}
