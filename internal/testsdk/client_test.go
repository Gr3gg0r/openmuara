package testsdk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/api"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	_ "github.com/openmuara/openmuara/internal/provider/defaultplugin"
	"github.com/openmuara/openmuara/internal/server"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestClientCreateAndGetPayment(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL)
	ctx := context.Background()

	created, err := client.CreatePayment(ctx, api.PaymentRequest{
		Provider:  "default",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "sdk-ref-1",
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}
	if created.Reference != "sdk-ref-1" {
		t.Errorf("reference: want sdk-ref-1, got %s", created.Reference)
	}

	got, err := client.GetPayment(ctx, "sdk-ref-1")
	if err != nil {
		t.Fatalf("get payment: %v", err)
	}
	if got.Reference != "sdk-ref-1" {
		t.Errorf("get reference: want sdk-ref-1, got %s", got.Reference)
	}
}

func TestClientRefund(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL)
	ctx := context.Background()

	if _, err := client.CreatePayment(ctx, api.PaymentRequest{
		Provider:  "default",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "sdk-ref-2",
	}); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	if _, err := client.Scenario(ctx, "success", "sdk-ref-2"); err != nil {
		t.Fatalf("scenario success: %v", err)
	}

	refunded, err := client.Refund(ctx, "sdk-ref-2")
	if err != nil {
		t.Fatalf("refund: %v", err)
	}
	if refunded.Status != "refunded" {
		t.Errorf("status: want refunded, got %s", refunded.Status)
	}
}

func TestClientScenario(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL)
	ctx := context.Background()

	if _, err := client.CreatePayment(ctx, api.PaymentRequest{
		Provider:  "default",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "sdk-ref-3",
	}); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	resp, err := client.Scenario(ctx, "success", "sdk-ref-3")
	if err != nil {
		t.Fatalf("scenario: %v", err)
	}
	if resp.Status != "paid" {
		t.Errorf("status: want paid, got %s", resp.Status)
	}
}

func TestClientListWebhooks(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := NewClient(ts.URL)
	ctx := context.Background()

	attempts, err := client.ListWebhooks(ctx)
	if err != nil {
		t.Fatalf("list webhooks: %v", err)
	}
	if attempts == nil {
		t.Error("expected empty slice, got nil")
	}
}

func TestClientReplayWebhook(t *testing.T) {
	d := webhook.NewDispatcherFromBuilder("http://127.0.0.1:1", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Store = webhook.NewMemoryStore()
	_ = d.Store.Save(&webhook.Attempt{Ref: "ref-1", Status: webhook.AttemptStatusDelivered})

	router := server.NewRouter(server.RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		Dispatcher:       d,
		TransactionStore: engine.NewMemoryStore(),
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := NewClient(ts.URL)
	ctx := context.Background()

	attempt, err := client.ReplayWebhook(ctx, "ref-1")
	if err != nil {
		t.Fatalf("replay webhook: %v", err)
	}
	if attempt.Ref != "ref-1" {
		t.Errorf("ref: want ref-1, got %q", attempt.Ref)
	}
}

func TestClientErrors(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		call    func(*Client, context.Context) error
		wantErr bool
	}{
		{
			name:   "create payment unexpected status",
			status: http.StatusBadRequest,
			body:   "bad",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.CreatePayment(ctx, api.PaymentRequest{})
				return err
			},
			wantErr: true,
		},
		{
			name:   "get payment not found",
			status: http.StatusNotFound,
			body:   "not found",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.GetPayment(ctx, "missing")
				return err
			},
			wantErr: true,
		},
		{
			name:   "refund conflict",
			status: http.StatusConflict,
			body:   "conflict",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.Refund(ctx, "ref")
				return err
			},
			wantErr: true,
		},
		{
			name:   "list webhooks unexpected status",
			status: http.StatusInternalServerError,
			body:   "boom",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.ListWebhooks(ctx)
				return err
			},
			wantErr: true,
		},
		{
			name:   "replay webhook unexpected status",
			status: http.StatusBadRequest,
			body:   "bad",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.ReplayWebhook(ctx, "ref")
				return err
			},
			wantErr: true,
		},
		{
			name:   "scenario unexpected status",
			status: http.StatusBadRequest,
			body:   "bad",
			call: func(c *Client, ctx context.Context) error {
				_, err := c.Scenario(ctx, "success", "ref")
				return err
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			client := NewClient(srv.URL)
			err := tc.call(client, context.Background())
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestClientDecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.GetPayment(context.Background(), "ref")
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func TestClientNetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	srv.Close()

	client := NewClient(srv.URL)
	ctx := context.Background()

	if _, err := client.CreatePayment(ctx, api.PaymentRequest{}); err == nil {
		t.Error("expected network error for CreatePayment")
	}
	if _, err := client.GetPayment(ctx, "ref"); err == nil {
		t.Error("expected network error for GetPayment")
	}
	if _, err := client.Refund(ctx, "ref"); err == nil {
		t.Error("expected network error for Refund")
	}
	if _, err := client.ListWebhooks(ctx); err == nil {
		t.Error("expected network error for ListWebhooks")
	}
	if _, err := client.ReplayWebhook(ctx, "ref"); err == nil {
		t.Error("expected network error for ReplayWebhook")
	}
	if _, err := client.Scenario(ctx, "success", "ref"); err == nil {
		t.Error("expected network error for Scenario")
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	d := webhook.NewDispatcherFromBuilder("http://127.0.0.1:1", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Store = webhook.NewMemoryStore()

	router := server.NewRouter(server.RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		Dispatcher:       d,
		TransactionStore: engine.NewMemoryStore(),
	})
	return httptest.NewServer(router)
}
