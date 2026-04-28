package config

import (
	"reflect"
	"testing"
)

func TestCleanStringList(t *testing.T) {
	got := CleanStringList([]string{"a.example.com\n\n# comment\nb.example.com", "a.example.com"})
	want := []string{"a.example.com", "b.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestNormalizeDefaults(t *testing.T) {
	cfg := Default()
	cfg.Normalize()
	if cfg.IntervalSeconds != 300 {
		t.Fatalf("interval = %d", cfg.IntervalSeconds)
	}
	if cfg.TTL != 600 {
		t.Fatalf("ttl = %d", cfg.TTL)
	}
	if cfg.DNSHE.APIBaseURL != DefaultDNSHEBaseURL {
		t.Fatalf("base url = %q", cfg.DNSHE.APIBaseURL)
	}
	if cfg.IPv6.Source != "url" {
		t.Fatalf("ipv6 source = %q", cfg.IPv6.Source)
	}
	if !cfg.IPv6.Enable {
		t.Fatal("ipv6 should be enabled by default")
	}
}
