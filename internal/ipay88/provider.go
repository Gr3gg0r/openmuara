// Package ipay88 emulates the iPay88 payment gateway.
package ipay88

import (
	"context"
	"net/http"
	"sync"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// ProviderName is the registered provider identifier.
const ProviderName = "ipay88"

// RegisterWith registers the iPay88 provider with the given registry.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
}

// Provider implements provider.Provider for the iPay88 Malaysia ePayment gateway.
type Provider struct {
	mu           sync.RWMutex
	merchantCode string
	merchantKey  string
	store        engine.TransactionStore
	dispatcher   *webhook.Dispatcher
	baseURL      string
	requests     map[string]PaymentRequest
	httpClient   *http.Client
}

// NewProvider returns an unconfigured iPay88 provider.
func NewProvider() *Provider {
	return &Provider{
		store:      engine.NewMemoryStore(),
		requests:   make(map[string]PaymentRequest),
		httpClient: http.DefaultClient,
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return ProviderName }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	mc, ok := cfg["merchant_code"].(string)
	if !ok || mc == "" {
		return errcode.New(errcode.EConfigMissing, "ipay88: merchant_code is required")
	}
	mk, ok := cfg["merchant_key"].(string)
	if !ok || mk == "" {
		return errcode.New(errcode.EConfigMissing, "ipay88: merchant_key is required")
	}
	p.merchantCode = mc
	p.merchantKey = mk
	return nil
}

// SetStore replaces the provider's transaction ledger.
func (p *Provider) SetStore(store engine.TransactionStore) {
	p.store = store
}

// SetBaseURL configures the base URL used for payment links.
func (p *Provider) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

// SetDispatcher configures the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(dispatcher *webhook.Dispatcher) {
	p.dispatcher = dispatcher
}

// SetHTTPClient replaces the default HTTP client used for backend/response posts.
func (p *Provider) SetHTTPClient(client *http.Client) {
	p.httpClient = client
}

// Routes returns the iPay88 route table.
func (p *Provider) Routes() []provider.Route {
	return []provider.Route{
		{Method: http.MethodPost, Path: "/ePayment/entry.asp", Handler: p.entryHandler()},
		{Method: http.MethodPost, Path: "/ePayment/enquiry.asp", Handler: p.requeryHandler()},
		{Method: http.MethodGet, Path: "/_admin/ipay88/pay/{refNo}", Handler: p.adminPayPageHandler()},
		{Method: http.MethodPost, Path: "/_admin/ipay88/pay/{refNo}", Handler: p.adminPayActionHandler()},
		{Method: http.MethodPost, Path: "/ipay88/response", Handler: p.responseHandler()},
		{Method: http.MethodPost, Path: "/ipay88/backend", Handler: p.backendHandler()},
	}
}

// ChargeHandler returns the entry handler for the provider charge endpoint.
func (p *Provider) ChargeHandler() http.Handler {
	return p.entryHandler()
}

// WebhookHandler returns the backend callback handler.
func (p *Provider) WebhookHandler() http.Handler {
	return p.backendHandler()
}

// PayloadBuilder returns an iPay88 form-encoded webhook payload builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return p.buildPayload
}

// PayloadHeaders returns the headers for iPay88 form-encoded webhook payloads.
func (p *Provider) PayloadHeaders(_ context.Context, _ provider.Transaction) (map[string]string, error) {
	return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, nil
}

// EscapeHandler returns nil; iPay88 uses the local admin payment page.
func (p *Provider) EscapeHandler() http.Handler { return nil }

func (p *Provider) getRequest(ref string) (PaymentRequest, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	req, ok := p.requests[ref]
	return req, ok
}

func (p *Provider) saveRequest(req PaymentRequest) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.requests[req.RefNo] = req
}
