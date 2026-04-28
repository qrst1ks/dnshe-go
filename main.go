package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/ddns"
	"github.com/qrst1ks/dnshe-go/internal/logbuf"
	"github.com/qrst1ks/dnshe-go/internal/web"
)

var version = "dev"

func main() {
	versionFlag := flag.Bool("v", false, "print version")
	configPath := flag.String("c", defaultConfigPath(), "config file path")
	listen := flag.String("l", defaultListen(), "listen address")
	intervalOverride := flag.Int("f", 0, "override update frequency in seconds")
	noWeb := flag.Bool("noweb", false, "disable web server")
	once := flag.Bool("once", false, "run sync once and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	logger := logbuf.New(200)
	store, err := config.NewStore(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	syncer := ddns.NewSyncer(store, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *once {
		status := syncer.RunOnce(ctx, true)
		if status.LastError != "" {
			os.Exit(1)
		}
		return
	}

	if !*noWeb {
		server := web.NewServer(*listen, store, syncer, logger)
		go func() {
			logger.Addf("INFO", "web server listening on http://%s", displayListen(*listen))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Addf("ERROR", "web server stopped: %v", err)
				stop()
			}
		}()
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = server.Shutdown(shutdownCtx)
		}()
	}

	interval := time.Duration(store.Get().IntervalSeconds) * time.Second
	if *intervalOverride > 0 {
		interval = time.Duration(*intervalOverride) * time.Second
	}
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	logger.Addf("INFO", "dnshe-go started, config=%s, interval=%s", store.Path(), interval)
	syncer.RunLoop(ctx, interval)
}

func defaultConfigPath() string {
	if v := strings.TrimSpace(os.Getenv("CONFIG_PATH")); v != "" {
		return v
	}
	return "data/config.json"
}

func defaultListen() string {
	if v := strings.TrimSpace(os.Getenv("LISTEN")); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("PORT")); v != "" {
		if _, err := strconv.Atoi(v); err == nil {
			return "0.0.0.0:" + v
		}
	}
	return "127.0.0.1:9876"
}

func displayListen(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "127.0.0.1" + addr
	}
	if strings.HasPrefix(addr, "0.0.0.0:") {
		return "127.0.0.1:" + strings.TrimPrefix(addr, "0.0.0.0:")
	}
	return addr
}
