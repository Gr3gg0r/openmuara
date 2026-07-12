package billplz_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/openmuara/openmuara/internal/billplz"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func validProviderConfig() map[string]any {
	return map[string]any{
		"api_key":         "muara-billplz-api-key",
		"x_signature_key": "muara-billplz-xsig-key",
		"collection_id":   "muara-collection-id",
	}
}

func TestProviderNameReturnsBillplz(t *testing.T) {
	p := billplz.NewProvider()
	if got := p.Name(); got != "billplz" {
		t.Errorf("name: want billplz, got %q", got)
	}
}

func TestProviderInitWithValidConfigSucceeds(t *testing.T) {
	p := billplz.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestProviderInitMissingAPIKey(t *testing.T) {
	p := billplz.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "api_key")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestProviderInitMissingXSignatureKey(t *testing.T) {
	p := billplz.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "x_signature_key")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing x_signature_key")
	}
}

func TestProviderInitCollectionIDOptional(t *testing.T) {
	p := billplz.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "collection_id")
	if err := p.Init(cfg); err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestProviderRoutes(t *testing.T) {
	p := billplz.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}

	routes := p.Routes()
	want := []struct {
		path   string
		method string
	}{
		{"/api/v3/collections", http.MethodPost},
		{"/api/v3/collections/{id}", http.MethodGet},
		{"/api/v3/bills", http.MethodPost},
		{"/api/v3/bills/{id}", http.MethodGet},
		{"/api/v3/bills/{id}", http.MethodDelete},
		{"/api/v3/collections/{id}/payment_methods", http.MethodGet},
		{"/_admin/billplz/pay/{id}", http.MethodGet},
		{"/_admin/billplz/pay/{id}", http.MethodPost},
		{"/billplz/redirect", http.MethodGet},
		{"/billplz/webhook", http.MethodPost},
	}

	got := make(map[string][]string)
	for _, r := range routes {
		got[r.Path] = append(got[r.Path], r.Method)
	}
	for _, w := range want {
		methods, ok := got[w.path]
		if !ok {
			t.Errorf("missing route %s %s", w.method, w.path)
			continue
		}
		found := false
		for _, m := range methods {
			if m == w.method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("route %s missing method %s", w.path, w.method)
		}
	}
}

func TestProviderSetters(t *testing.T) {
	p := billplz.NewProvider()
	store := engine.NewMemoryStore()
	p.SetStore(store)

	d := webhook.NewDispatcherFromBuilder("http://127.0.0.1:1", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	p.SetDispatcher(d)

	if p.ChargeHandler() == nil {
		t.Fatal("expected charge handler after setters")
	}
}
