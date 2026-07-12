// Package defaultplugin provides a minimal example/fallback provider implementation.
package defaultplugin

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/provider"
)

func init() {
	RegisterWith(provider.Default())
}

// RegisterWith registers the default provider with the given registry. Tests can
// use this to build isolated registries instead of relying on the global default.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
	slog.Debug("registered provider", "name", "default")
}

// Provider is a minimal example/fallback provider.
type Provider struct{}

// NewProvider returns a new default provider.
func NewProvider() *Provider { return &Provider{} }

// Name returns the provider identifier.
func (p *Provider) Name() string { return "default" }

// Init succeeds with any config (or no config).
func (p *Provider) Init(_ map[string]any) error { return nil }

// Routes returns example routes.
func (p *Provider) Routes() []provider.Route {
	return []provider.Route{
		{Method: http.MethodGet, Path: "/default", Handler: p.statusHandler()},
		{Method: http.MethodPost, Path: "/default/charge", Handler: p.ChargeHandler()},
		{Method: http.MethodPost, Path: "/default/webhook", Handler: p.WebhookHandler()},
	}
}

// ChargeHandler returns a deterministic example charge response.
func (p *Provider) ChargeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := uuid.NewString()
		resp := map[string]any{
			"provider":       "default",
			"transaction_id": tid,
			"status":         "success",
		}
		audit.FromContext(r.Context()).Log(r.Context(), "charge.created", "transaction", tid, "", "ok")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}

// WebhookHandler acknowledges receipt.
func (p *Provider) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "acknowledged"})
	})
}

// PayloadBuilder returns a simple JSON payload.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		return json.Marshal(map[string]any{
			"provider":  "default",
			"reference": tx.Reference,
			"status":    tx.Status,
		})
	}
}

// EscapeHandler returns nil because the default plugin has no escape flow.
func (p *Provider) EscapeHandler() http.Handler { return nil }

func (p *Provider) statusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"provider": "default", "status": "ok"})
	})
}
