package duckdns

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const duckDNSUpdateURL = "https://www.duckdns.org/update"

type DuckDNSConnector struct {
	httpClient *http.Client
	token      string
}

func NewDuckDNSConnector(token string) *DuckDNSConnector {
	return &DuckDNSConnector{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		token:      token,
	}
}

func (c *DuckDNSConnector) SetTXTRecord(ctx context.Context, domain, txt string) (*http.Response, error) {
	return c.updateTXTRecord(ctx, domain, txt)
}

func (c *DuckDNSConnector) CleanTXTRecord(ctx context.Context, domain string) (*http.Response, error) {
	return c.updateTXTRecord(ctx, domain, "")
}

func (c *DuckDNSConnector) updateTXTRecord(ctx context.Context, domain, txt string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, duckDNSUpdateURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("domains", domain)
	q.Add("token", c.token)
	q.Add("dns_txt", txt)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update TXT record for domain %s: %w", domain, err)
	}

	return resp, nil
}
