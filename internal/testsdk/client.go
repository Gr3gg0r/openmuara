// Package testsdk provides a provider-agnostic HTTP client for writing
// integration tests against a running OpenMuara server.
package testsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openmuara/openmuara/internal/api"
	"github.com/openmuara/openmuara/internal/webhook"
)

// Client talks to an OpenMuara HTTP server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a client for the given base URL.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// CreatePayment calls POST /v1/pay.
func (c *Client) CreatePayment(ctx context.Context, req api.PaymentRequest) (*api.PaymentResponse, error) {
	var resp api.PaymentResponse
	status, err := c.post(ctx, "/v1/pay", req, &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated && status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return &resp, nil
}

// GetPayment calls GET /v1/pay/{ref}.
func (c *Client) GetPayment(ctx context.Context, ref string) (*api.PaymentResponse, error) {
	var resp api.PaymentResponse
	status, err := c.get(ctx, fmt.Sprintf("/v1/pay/%s", ref), &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return &resp, nil
}

// Refund calls POST /v1/refund/{ref}.
func (c *Client) Refund(ctx context.Context, ref string) (*api.PaymentResponse, error) {
	var resp api.PaymentResponse
	status, err := c.post(ctx, fmt.Sprintf("/v1/refund/%s", ref), nil, &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return &resp, nil
}

// ListWebhooks calls GET /_admin/webhooks.
func (c *Client) ListWebhooks(ctx context.Context) ([]*webhook.Attempt, error) {
	var resp struct {
		Results []*webhook.Attempt `json:"results"`
	}
	status, err := c.get(ctx, "/_admin/webhooks", &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return resp.Results, nil
}

// ReplayWebhook calls POST /_admin/webhooks/{ref}/replay.
func (c *Client) ReplayWebhook(ctx context.Context, ref string) (*webhook.Attempt, error) {
	var resp webhook.Attempt
	status, err := c.post(ctx, fmt.Sprintf("/_admin/webhooks/%s/replay", ref), nil, &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return &resp, nil
}

// Scenario calls POST /_admin/scenario/{outcome}?ref=...
func (c *Client) Scenario(ctx context.Context, outcome, ref string) (*ScenarioResponse, error) {
	var resp ScenarioResponse
	status, err := c.post(ctx, fmt.Sprintf("/_admin/scenario/%s?ref=%s", outcome, ref), nil, &resp)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", status)
	}
	return &resp, nil
}

func (c *Client) get(ctx context.Context, path string, out any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return 0, err
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, body any, out any) (int, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bodyReader)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) (int, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if out != nil && resp.StatusCode < http.StatusBadRequest {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return resp.StatusCode, fmt.Errorf("decode response: %w", err)
		}
	}
	return resp.StatusCode, nil
}
