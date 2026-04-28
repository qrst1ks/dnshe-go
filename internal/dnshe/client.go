package dnshe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/qrst1ks/dnshe-go/internal/config"
)

type Client struct {
	BaseURL    string
	APIKey     string
	APISecret  string
	HTTPClient *http.Client
}

type responseBodyError struct {
	Context     string
	StatusCode  int
	ContentType string
	Body        []byte
	Err         error
}

func (e responseBodyError) Error() string {
	preview := strings.TrimSpace(string(e.Body))
	if len(preview) > 300 {
		preview = preview[:300] + "..."
	}
	preview = strings.Join(strings.Fields(preview), " ")
	return fmt.Sprintf("%s: %v (http=%d content-type=%q body=%q)", e.Context, e.Err, e.StatusCode, e.ContentType, preview)
}

type Record struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type UpdateResult struct {
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"`
	RecordID   int    `json:"record_id"`
	OldIP      string `json:"old_ip,omitempty"`
	NewIP      string `json:"new_ip"`
	Status     string `json:"status"`
}

func NewClient(cfg config.DNSHEConfig) *Client {
	baseURL := strings.TrimSpace(cfg.APIBaseURL)
	if baseURL == "" {
		baseURL = config.DefaultDNSHEBaseURL
	}
	return &Client{
		BaseURL:   baseURL,
		APIKey:    strings.TrimSpace(cfg.APIKey),
		APISecret: strings.TrimSpace(cfg.APISecret),
		HTTPClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) EnsureRecord(ctx context.Context, domain string, recordType string, ip string, ttl int) (UpdateResult, error) {
	result := UpdateResult{
		Domain:     domain,
		RecordType: strings.ToUpper(recordType),
		NewIP:      ip,
	}

	subdomainID, err := c.FindSubdomain(ctx, domain)
	if err != nil {
		return result, err
	}
	if subdomainID == 0 {
		return result, fmt.Errorf("subdomain not found: %s", domain)
	}

	record, err := c.FindRecord(ctx, subdomainID, result.RecordType)
	if err != nil {
		return result, err
	}
	if record.ID == 0 {
		return result, fmt.Errorf("record not found: %s %s", domain, result.RecordType)
	}
	result.RecordID = record.ID
	result.OldIP = record.Content

	if record.Content == ip {
		result.Status = "unchanged"
		return result, nil
	}

	if _, err := c.UpdateRecord(ctx, record.ID, ip, ttl); err != nil {
		return result, err
	}
	result.Status = "updated"
	return result, nil
}

func (c *Client) FindSubdomain(ctx context.Context, domain string) (int, error) {
	endpoint, err := c.endpoint("subdomains", "list", nil)
	if err != nil {
		return 0, err
	}
	body, meta, err := c.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	var resp struct {
		Subdomains []struct {
			ID         int    `json:"id"`
			FullDomain string `json:"full_domain"`
		} `json:"subdomains"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, responseBodyError{
			Context:     "parse subdomains response",
			StatusCode:  meta.StatusCode,
			ContentType: meta.ContentType,
			Body:        body,
			Err:         err,
		}
	}

	for _, item := range resp.Subdomains {
		if strings.EqualFold(strings.TrimSpace(item.FullDomain), strings.TrimSpace(domain)) {
			return item.ID, nil
		}
	}
	return 0, nil
}

func (c *Client) FindRecord(ctx context.Context, subdomainID int, recordType string) (Record, error) {
	params := url.Values{}
	params.Set("subdomain_id", strconv.Itoa(subdomainID))
	endpoint, err := c.endpoint("dns_records", "list", params)
	if err != nil {
		return Record{}, err
	}
	body, meta, err := c.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Record{}, err
	}

	var resp struct {
		Records []Record `json:"records"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return Record{}, responseBodyError{
			Context:     "parse records response",
			StatusCode:  meta.StatusCode,
			ContentType: meta.ContentType,
			Body:        body,
			Err:         err,
		}
	}

	recordType = strings.ToUpper(strings.TrimSpace(recordType))
	for _, record := range resp.Records {
		if strings.EqualFold(record.Type, recordType) {
			return record, nil
		}
	}
	return Record{}, nil
}

func (c *Client) UpdateRecord(ctx context.Context, recordID int, ip string, ttl int) (map[string]any, error) {
	endpoint, err := c.endpoint("dns_records", "update", nil)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"record_id": recordID,
		"content":   ip,
		"ttl":       ttl,
	}
	body, meta, err := c.do(ctx, http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, responseBodyError{
			Context:     "parse update response",
			StatusCode:  meta.StatusCode,
			ContentType: meta.ContentType,
			Body:        body,
			Err:         err,
		}
	}
	return result, nil
}

func (c *Client) endpoint(endpoint, action string, params url.Values) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("endpoint", endpoint)
	q.Set("action", action)
	for key, values := range params {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type responseMeta struct {
	StatusCode  int
	ContentType string
}

func (c *Client) do(ctx context.Context, method string, endpoint string, payload any) ([]byte, responseMeta, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, responseMeta{}, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, responseMeta{}, err
	}
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-API-Secret", c.APISecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, responseMeta{}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	meta := responseMeta{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
	}
	if resp.StatusCode >= 400 {
		return nil, meta, fmt.Errorf("DNSHE API error %d: %s", resp.StatusCode, string(respBody))
	}
	if err := validateDNSHEResult(respBody); err != nil {
		return nil, meta, err
	}
	return respBody, meta, nil
}

func validateDNSHEResult(body []byte) error {
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	success := strings.TrimSpace(strings.ToLower(stringFromAny(m["success"])))
	if success == "false" || success == "0" || success == "no" {
		return fmt.Errorf("DNSHE API business error: %s", firstNonEmpty(stringFromAny(m["msg"]), stringFromAny(m["message"]), "unknown error"))
	}

	status := strings.TrimSpace(strings.ToLower(stringFromAny(m["status"])))
	if status == "error" || status == "failed" || status == "fail" {
		return fmt.Errorf("DNSHE API business error: %s", firstNonEmpty(stringFromAny(m["msg"]), stringFromAny(m["message"]), "unknown error"))
	}

	codeStr := strings.TrimSpace(stringFromAny(m["code"]))
	if codeStr != "" {
		if code, err := strconv.Atoi(codeStr); err == nil && code != 0 {
			return fmt.Errorf("DNSHE API code=%d: %s", code, firstNonEmpty(stringFromAny(m["msg"]), stringFromAny(m["message"]), "unknown error"))
		}
	}

	if errMsg := firstNonEmpty(stringFromAny(m["error"]), stringFromAny(m["err"])); errMsg != "" {
		return fmt.Errorf("DNSHE API business error: %s", errMsg)
	}
	return nil
}

func stringFromAny(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return fmt.Sprintf("%.0f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
