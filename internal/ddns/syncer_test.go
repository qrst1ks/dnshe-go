package ddns

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/logbuf"
)

func TestRunOnceResolvesIPv6WithoutProviders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cfg := config.Default()
	cfg.IPv6.Source = "cmd"
	cfg.IPv6.Command = "echo 2001:db8::5"
	if err := config.Save(path, cfg); err != nil {
		t.Fatal(err)
	}
	store, err := config.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}

	syncer := NewSyncer(store, logbuf.New(20))
	status := syncer.RunOnce(context.Background(), true)
	if status.LastError != "" {
		t.Fatalf("last error = %q", status.LastError)
	}
	if status.CurrentIPv6 != "2001:db8::5" {
		t.Fatalf("ipv6 = %q", status.CurrentIPv6)
	}
	if len(status.Results) != 0 {
		t.Fatalf("results = %#v", status.Results)
	}
}
