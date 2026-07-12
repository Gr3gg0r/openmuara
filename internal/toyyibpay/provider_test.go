package toyyibpay

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestProviderName(t *testing.T) {
	p := NewProvider()
	if got := p.Name(); got != ProviderName {
		t.Fatalf("Name() = %q, want %q", got, ProviderName)
	}
}

func TestProviderInitRequiresSecret(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{}); err == nil {
		t.Fatal("Init() expected error for missing secret")
	}
	if err := p.Init(map[string]any{"user_secret_key": ""}); err == nil {
		t.Fatal("Init() expected error for empty secret")
	}
}

func TestProviderInitAcceptsOptionalCategory(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{
		"user_secret_key": "secret",
		"category_code":   "CAT123",
	}); err != nil {
		t.Fatalf("Init() error: %v", err)
	}
	if p.secret != "secret" {
		t.Fatalf("secret not stored")
	}
	if p.defaultCategory != "CAT123" {
		t.Fatalf("default category not stored")
	}
}

func TestProviderRoutes(t *testing.T) {
	p := NewProvider()
	routes := p.Routes()
	if len(routes) == 0 {
		t.Fatal("Routes() returned empty table")
	}
	paths := make(map[string]bool)
	for _, r := range routes {
		paths[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{
		"POST /index.php/api/createCategory",
		"POST /index.php/api/getCategoryDetails",
		"POST /index.php/api/createBill",
		"POST /index.php/api/getBillTransactions",
		"POST /index.php/api/inactiveBill",
		"GET /_admin/toyyibpay/pay/{billCode}",
		"POST /_admin/toyyibpay/pay/{billCode}",
		"GET /toyyibpay/return",
		"POST /toyyibpay/webhook",
	} {
		if !paths[want] {
			t.Fatalf("missing route %q", want)
		}
	}
}

func TestProviderSetters(t *testing.T) {
	p := NewProvider()
	store := engine.NewMemoryStore()
	p.SetStore(store)
	if p.store != store {
		t.Fatal("SetStore did not store ledger")
	}
	p.SetBaseURL("http://localhost:9000")
	if p.baseURL != "http://localhost:9000" {
		t.Fatal("SetBaseURL did not store base URL")
	}
	p.SetDispatcher(&webhook.Dispatcher{})
	if p.dispatcher == nil {
		t.Fatal("SetDispatcher did not store dispatcher")
	}
}

func TestPayloadBuilderReturnsFormEncoded(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	body, err := p.PayloadBuilder()(context.Background(), provider.Transaction{
		Reference: bill.OrderID,
		Status:    string(engine.TransactionStatusPaid),
	})
	if err != nil {
		t.Fatalf("PayloadBuilder error: %v", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		t.Fatalf("parse payload: %v", err)
	}
	if values.Get("status") != "1" {
		t.Fatalf("status = %q, want 1", values.Get("status"))
	}
	if values.Get("billcode") != bill.BillCode {
		t.Fatalf("billcode = %q, want %q", values.Get("billcode"), bill.BillCode)
	}
	if values.Get("order_id") != bill.OrderID {
		t.Fatalf("order_id mismatch")
	}
	if values.Get("amount") != "1000" {
		t.Fatalf("amount = %q, want 1000", values.Get("amount"))
	}
	if values.Get("hash") == "" {
		t.Fatal("hash missing")
	}
	if !VerifyCallback(p.secret, values) {
		t.Fatal("callback hash did not verify")
	}
}

func TestPayloadHeaders(t *testing.T) {
	p := newTestProvider(t)
	headers, err := p.PayloadHeaders(context.Background(), provider.Transaction{})
	if err != nil {
		t.Fatalf("PayloadHeaders error: %v", err)
	}
	if headers["Content-Type"] != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q", headers["Content-Type"])
	}
}

func TestPayloadBuilderMissingBill(t *testing.T) {
	p := newTestProvider(t)
	_, err := p.PayloadBuilder()(context.Background(), provider.Transaction{Reference: "missing"})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not-found error, got %v", err)
	}
}
