package ddns

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
)

var ipv6CandidateReg = regexp.MustCompile(`([0-9A-Fa-f:.]{2,})`)

func ResolveIPv6(ctx context.Context, cfg config.IPConfig) (string, error) {
	cfg.Normalize()
	switch cfg.Source {
	case "url":
		return resolveFromURLs(ctx, cfg.URLs)
	case "interface":
		return resolveFromInterface(cfg.Interface)
	case "cmd":
		return resolveFromCommand(ctx, cfg.Command)
	default:
		return "", fmt.Errorf("unknown IP source: %s", cfg.Source)
	}
}

func resolveFromURLs(ctx context.Context, urls []string) (string, error) {
	if len(urls) == 0 {
		return "", errors.New("no IP URL configured")
	}
	client := ipv6HTTPClient()
	var lastErr error
	for _, endpoint := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("User-Agent", "dnshe-go")
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("%s returned HTTP %d", endpoint, resp.StatusCode)
			continue
		}
		if readErr != nil {
			lastErr = readErr
			continue
		}
		ip := ExtractIPv6(string(body))
		if ip != "" {
			return ip, nil
		}
		lastErr = fmt.Errorf("no IPv6 address in response from %s", endpoint)
	}
	if lastErr == nil {
		lastErr = errors.New("failed to resolve IP")
	}
	return "", lastErr
}

func resolveFromInterface(name string) (string, error) {
	var interfaces []net.Interface
	if strings.TrimSpace(name) != "" {
		iface, err := net.InterfaceByName(name)
		if err != nil {
			return "", err
		}
		interfaces = []net.Interface{*iface}
	} else {
		var err error
		interfaces, err = net.Interfaces()
		if err != nil {
			return "", err
		}
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := ipFromAddr(addr)
			if ip == nil || !ip.IsGlobalUnicast() || ip.IsLinkLocalUnicast() {
				continue
			}
			if ip.To4() == nil {
				return ip.String(), nil
			}
		}
	}
	if name != "" {
		return "", fmt.Errorf("no IPv6 address found on interface %s", name)
	}
	return "", errors.New("no IPv6 address found on active interfaces")
}

func resolveFromCommand(ctx context.Context, command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", errors.New("empty command")
	}
	runCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(runCtx, "powershell", "-Command", command)
	} else {
		cmd = exec.CommandContext(runCtx, "sh", "-c", command)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	ip := ExtractIPv6(string(out))
	if ip == "" {
		return "", errors.New("command output does not contain IPv6 address")
	}
	return ip, nil
}

func ExtractIPv6(text string) string {
	for _, candidate := range ipv6CandidateReg.FindAllString(text, -1) {
		ip := net.ParseIP(candidate)
		if ip != nil && ip.To4() == nil {
			return ip.String()
		}
	}
	return ""
}

func ipFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}

func ipv6HTTPClient() *http.Client {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	return &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, "tcp6", addr)
			},
		},
	}
}

func maskIP(ip string) string {
	if ip == "" {
		return "-"
	}
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) > 4 {
			return strings.Join(parts[:2], ":") + ":****:" + parts[len(parts)-1]
		}
		return ip
	}
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".***.***"
	}
	return ip
}
