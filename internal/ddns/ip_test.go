package ddns

import "testing"

func TestExtractIPv6(t *testing.T) {
	if got := ExtractIPv6("addr 2001:db8::42 ok"); got != "2001:db8::42" {
		t.Fatalf("ipv6 = %q", got)
	}
	if got := ExtractIPv6("203.0.113.7"); got != "" {
		t.Fatalf("non-IPv6 address should be ignored, got %q", got)
	}
	if got := ExtractIPv6("::ffff:192.168.1.2"); got != "" {
		t.Fatalf("mapped ipv6 should be ignored, got %q", got)
	}
}
