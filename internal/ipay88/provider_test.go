package ipay88

import (
	"context"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func validConfig() map[string]any {
	return map[string]any{
		"merchant_code": "M00001",
		"merchant_key":  "secret-key",
	}
}

func TestProviderName(t *testing.T) {
	p := NewProvider()
	if got := p.Name(); got != ProviderName {
		t.Errorf("name: want %q, got %q", ProviderName, got)
	}
}

func TestProviderInitValid(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
}

func TestProviderInitMissingMerchantCode(t *testing.T) {
	p := NewProvider()
	cfg := validConfig()
	delete(cfg, "merchant_code")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing merchant_code")
	}
}

func TestProviderInitMissingMerchantKey(t *testing.T) {
	p := NewProvider()
	cfg := validConfig()
	delete(cfg, "merchant_key")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing merchant_key")
	}
}

func TestProviderRoutes(t *testing.T) {
	p := NewProvider()
	_ = p.Init(validConfig())
	routes := p.Routes()
	if len(routes) != 6 {
		t.Fatalf("routes: want 6, got %d", len(routes))
	}
}

func TestProviderInterfaceMethods(t *testing.T) {
	p := NewProvider()
	_ = p.Init(validConfig())

	if p.ChargeHandler() == nil {
		t.Error("expected ChargeHandler to be non-nil")
	}
	if p.WebhookHandler() == nil {
		t.Error("expected WebhookHandler to be non-nil")
	}
	if p.EscapeHandler() != nil {
		t.Error("expected EscapeHandler to be nil")
	}

	p.SetStore(engine.NewMemoryStore())
	p.SetBaseURL("http://localhost")
	p.SetDispatcher(webhook.NewDispatcherFromBuilder("http://example.com", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("ok"), nil }, nil))

	builder := p.PayloadBuilder()
	if builder == nil {
		t.Fatal("expected PayloadBuilder to be non-nil")
	}

	headers, err := p.PayloadHeaders(context.Background(), provider.Transaction{})
	if err != nil {
		t.Fatalf("payload headers: %v", err)
	}
	if !strings.Contains(headers["Content-Type"], "x-www-form-urlencoded") {
		t.Errorf("unexpected content type: %q", headers["Content-Type"])
	}
}
