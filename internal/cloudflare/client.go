package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
	"github.com/qrst1ks/dnshe-go/internal/provider"
)

type Client struct {
	APIToken   string
	BaseURL    string
	HTTPClient *http.Client
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type record struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func NewClient(cfg config.CloudflareConfig) *Client {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.APIBaseURL), "/")
	if baseURL == "" {
		baseURL = config.DefaultCloudflareAPIBase
	}
	return &Client{
		APIToken: strings.TrimSpace(cfg.APIToken),
		BaseURL:  baseURL,
		HTTPClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) EnsureRecord(ctx context.Context, domain string, recordType string, ip string, ttl int) (provider.Result, error) {
	result := provider.Result{
		Domain:     domain,
		RecordType: strings.ToUpper(strings.TrimSpace(recordType)),
		NewIP:      ip,
	}
	fqdn := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")

	z, err := c.findZone(ctx, fqdn)
	if err != nil {
		return result, err
	}

	existing, found, err := c.findRecord(ctx, z.ID, fqdn, result.RecordType)
	if err != nil {
		return result, err
	}
	if found {
		result.RecordID = existing.ID
		result.OldIP = existing.Content
		if sameIP(existing.Content, ip) {
			result.Status = provider.StatusUnchanged
			return result, nil
		}
		if err := c.updateRecord(ctx, z.ID, existing, ip, ttl); err != nil {
			return result, err
		}
		result.Status = provider.StatusUpdated
		return result, nil
	}

	created, err := c.createRecord(ctx, z.ID, fqdn, result.RecordType, ip, ttl)
	if err != nil {
		return result, err
	}
	result.RecordID = created.ID
	result.Status = provider.StatusCreated
	return result, nil
}

func (c *Client) findZone(ctx context.Context, domain string) (zone, error) {
	candidates := zoneCandidates(domain)
	if len(candidates) == 0 {
		return zone{}, fmt.Errorf("invalid domain: %s", domain)
	}
	for _, candidate := range candidates {
		params := url.Values{}
		params.Set("name", candidate)
		params.Set("status", "active")
		var zones []zone
		if err := c.do(ctx, http.MethodGet, "/zones?"+params.Encode(), nil, &zones); err != nil {
			return zone{}, err
		}
		for _, z := range zones {
			if strings.EqualFold(strings.TrimSuffix(z.Name, "."), candidate) {
				return z, nil
			}
		}
	}
	return zone{}, fmt.Errorf("zone not found for %s", domain)
}

func (c *Client) findRecord(ctx context.Context, zoneID string, domain string, recordType string) (record, bool, error) {
	params := url.Values{}
	params.Set("type", recordType)
	params.Set("name", domain)
	var records []record
	path := "/zones/" + url.PathEscape(zoneID) + "/dns_records?" + params.Encode()
	if err := c.do(ctx, http.MethodGet, path, nil, &records); err != nil {
		return record{}, false, err
	}
	for _, item := range records {
		if strings.EqualFold(item.Type, recordType) &&
			strings.EqualFold(strings.TrimSuffix(item.Name, "."), strings.TrimSuffix(domain, ".")) {
			return item, true, nil
		}
	}
	return record{}, false, nil
}

func (c *Client) updateRecord(ctx context.Context, zoneID string, existing record, ip string, ttl int) error {
	payload := map[string]any{
		"type":    existing.Type,
		"name":    existing.Name,
		"content": ip,
		"ttl":     cloudflareTTL(ttl),
		"proxied": existing.Proxied,
	}
	path := "/zones/" + url.PathEscape(zoneID) + "/dns_records/" + url.PathEscape(existing.ID)
	var updated record
	return c.do(ctx, http.MethodPut, path, payload, &updated)
}

func (c *Client) createRecord(ctx context.Context, zoneID string, domain string, recordType string, ip string, ttl int) (record, error) {
	payload := map[string]any{
		"type":    recordType,
		"name":    domain,
		"content": ip,
		"ttl":     cloudflareTTL(ttl),
		"proxied": false,
	}
	path := "/zones/" + url.PathEscape(zoneID) + "/dns_records"
	var created record
	if err := c.do(ctx, http.MethodPost, path, payload, &created); err != nil {
		return record{}, err
	}
	return created, nil
}

func (c *Client) do(ctx context.Context, method string, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "dnshe-go")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	var envelope struct {
		Success bool            `json:"success"`
		Errors  []apiError      `json:"errors"`
		Result  json.RawMessage `json:"result"`
	}
	parseErr := json.Unmarshal(respBody, &envelope)

	if resp.StatusCode >= 400 {
		if parseErr == nil && len(envelope.Errors) > 0 {
			return fmt.Errorf("Cloudflare API HTTP %d: %s", resp.StatusCode, formatAPIErrors(envelope.Errors))
		}
		return fmt.Errorf("Cloudflare API HTTP %d: %s", resp.StatusCode, bodyPreview(respBody))
	}
	if parseErr != nil {
		return fmt.Errorf("parse Cloudflare response: %w (body=%q)", parseErr, bodyPreview(respBody))
	}
	if !envelope.Success {
		return fmt.Errorf("Cloudflare API error: %s", formatAPIErrors(envelope.Errors))
	}
	if out != nil && len(envelope.Result) > 0 {
		if err := json.Unmarshal(envelope.Result, out); err != nil {
			return fmt.Errorf("parse Cloudflare result: %w", err)
		}
	}
	return nil
}

func zoneCandidates(domain string) []string {
	domain = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return nil
	}
	candidates := make([]string, 0, len(parts)-1)
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "" {
			return nil
		}
		candidates = append(candidates, strings.Join(parts[i:], "."))
	}
	return candidates
}

// Cloudflare only accepts TTL 1 (auto) or >= 60.
func cloudflareTTL(ttl int) int {
	if ttl == 1 || ttl >= 60 {
		return ttl
	}
	return 60
}

func sameIP(a string, b string) bool {
	ipA := net.ParseIP(strings.TrimSpace(a))
	ipB := net.ParseIP(strings.TrimSpace(b))
	if ipA != nil && ipB != nil {
		return ipA.Equal(ipB)
	}
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func formatAPIErrors(errs []apiError) string {
	if len(errs) == 0 {
		return "unknown error"
	}
	parts := make([]string, 0, len(errs))
	for _, item := range errs {
		message := strings.TrimSpace(item.Message)
		if message == "" {
			message = "unknown error"
		}
		if item.Code != 0 {
			message = fmt.Sprintf("%d: %s", item.Code, message)
		}
		parts = append(parts, message)
	}
	return strings.Join(parts, "; ")
}

func bodyPreview(body []byte) string {
	preview := strings.Join(strings.Fields(string(body)), " ")
	if len(preview) > 300 {
		preview = preview[:300] + "..."
	}
	return preview
}
