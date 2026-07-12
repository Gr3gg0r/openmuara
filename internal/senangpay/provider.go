package senangpay

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

// ProviderName is the registered provider identifier.
const ProviderName = "senangpay"

// RegisterWith registers the SenangPay provider with the given registry. Tests can
// use this to build isolated registries instead of relying on the global default.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
}

// Provider implements the generic provider interface for SenangPay.
type Provider struct {
	secret  string
	store   engine.TransactionStore
	baseURL string
}

// NewProvider returns an unconfigured SenangPay provider.
func NewProvider() *Provider {
	return &Provider{store: engine.NewMemoryStore()}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return ProviderName }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	secret, ok := cfg["secret_key"].(string)
	if !ok || secret == "" {
		return errcode.New(errcode.EConfigMissing, "senangpay: secret_key is required")
	}
	p.secret = secret
	return nil
}

// SetStore replaces the provider's transaction store.
func (p *Provider) SetStore(store engine.TransactionStore) {
	p.store = store
}

// SetBaseURL configures the base URL used for payment links.
func (p *Provider) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

// Routes returns the SenangPay route table.
func (p *Provider) Routes() []provider.Route {
	return []provider.Route{
		{Method: http.MethodPost, Path: "/senangpay/charge", Handler: NewChargeHandler(p.secret, p.store, p.baseURL)},
		{Method: http.MethodGet, Path: "/senangpay/callback", Handler: NewCallbackHandler(p.store)},
		{Method: http.MethodGet, Path: "/senangpay/query", Handler: NewStatusHandler(p.secret, p.store)},
		{Method: http.MethodPost, Path: "/senangpay/webhook", Handler: NewWebhookHandler(p.store)},
	}
}

// ChargeHandler returns the charge endpoint handler.
func (p *Provider) ChargeHandler() http.Handler {
	return NewChargeHandler(p.secret, p.store, p.baseURL)
}

// WebhookHandler returns the webhook receiver handler.
func (p *Provider) WebhookHandler() http.Handler {
	return NewWebhookHandler(p.store)
}

// PayloadBuilder returns a simple JSON payload builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		return json.Marshal(map[string]any{
			"provider":  ProviderName,
			"reference": tx.Reference,
			"status":    tx.Status,
		})
	}
}

// EscapeHandler returns nil; SenangPay uses the callback flow.
func (p *Provider) EscapeHandler() http.Handler { return nil }
