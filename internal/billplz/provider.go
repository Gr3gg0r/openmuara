// Package billplz emulates the Billplz payment gateway.
package billplz

import (
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// ProviderName is the registered provider identifier.
const ProviderName = "billplz"

// RegisterWith registers the Billplz provider with the given registry.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
}

// Provider implements the generic provider interface for Billplz.
type Provider struct {
	apiKey              string
	xSignatureKey       string
	defaultCollectionID string
	baseURL             string
	collections         *collectionStore
	bills               *billStore
	store               engine.TransactionStore
	dispatcher          *webhook.Dispatcher
}

// NewProvider returns an unconfigured Billplz provider.
func NewProvider() *Provider {
	return &Provider{
		collections: newCollectionStore(),
		bills:       newBillStore(),
		store:       engine.NewMemoryStore(),
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return ProviderName }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	apiKey, ok := cfg["api_key"].(string)
	if !ok || apiKey == "" {
		return errcode.New(errcode.EConfigMissing, "billplz: api_key is required")
	}
	xSig, ok := cfg["x_signature_key"].(string)
	if !ok || xSig == "" {
		return errcode.New(errcode.EConfigMissing, "billplz: x_signature_key is required")
	}
	p.apiKey = apiKey
	p.xSignatureKey = xSig
	if cid, ok := cfg["collection_id"].(string); ok {
		p.defaultCollectionID = cid
	}
	return nil
}

// SetBaseURL configures the base URL used for payment links.
func (p *Provider) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

// SetDispatcher configures the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(dispatcher *webhook.Dispatcher) {
	p.dispatcher = dispatcher
}

// SetStore replaces the provider's transaction ledger.
func (p *Provider) SetStore(store engine.TransactionStore) {
	p.store = store
}

// Routes returns the Billplz route table.
func (p *Provider) Routes() []provider.Route {
	return []provider.Route{
		{Method: http.MethodPost, Path: "/api/v3/collections", Handler: p.collectionsHandler()},
		{Method: http.MethodGet, Path: "/api/v3/collections/{id}", Handler: p.collectionHandler()},
		{Method: http.MethodPost, Path: "/api/v3/bills", Handler: p.billsHandler()},
		{Method: http.MethodGet, Path: "/api/v3/bills/{id}", Handler: p.billHandler()},
		{Method: http.MethodDelete, Path: "/api/v3/bills/{id}", Handler: p.deleteBillHandler()},
		{Method: http.MethodGet, Path: "/api/v3/collections/{id}/payment_methods", Handler: p.paymentMethodsHandler()},
		{Method: http.MethodGet, Path: "/_admin/billplz/pay/{id}", Handler: NewPayPageHandler(p.bills)},
		{Method: http.MethodPost, Path: "/_admin/billplz/pay/{id}", Handler: p.payActionHandler()},
		{Method: http.MethodGet, Path: "/billplz/redirect", Handler: NewRedirectHandler(p.bills, p.xSignatureKey)},
		{Method: http.MethodPost, Path: "/billplz/webhook", Handler: NewWebhookHandler(p.xSignatureKey)},
	}
}

func (p *Provider) collectionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewCollectionsHandler(p.apiKey, p.collections).ServeHTTP(w, r)
	})
}

func (p *Provider) collectionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewCollectionHandler(p.apiKey, p.collections).ServeHTTP(w, r)
	})
}

func (p *Provider) billsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewBillsHandler(p.apiKey, p.defaultCollectionID, p.bills, p.collections, p.store, p.baseURL).ServeHTTP(w, r)
	})
}

func (p *Provider) billHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewBillHandler(p.apiKey, p.bills).ServeHTTP(w, r)
	})
}

func (p *Provider) deleteBillHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewDeleteBillHandler(p.apiKey, p.bills).ServeHTTP(w, r)
	})
}

func (p *Provider) paymentMethodsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewPaymentMethodsHandler(p.apiKey, p.collections).ServeHTTP(w, r)
	})
}

func (p *Provider) payActionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewPayActionHandler(p.bills, p.store, p.dispatcher, p.xSignatureKey).ServeHTTP(w, r)
	})
}

// ChargeHandler returns a no-op handler (Billplz uses /api/v3/bills).
func (p *Provider) ChargeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// WebhookHandler returns the Billplz webhook receiver handler.
func (p *Provider) WebhookHandler() http.Handler {
	return NewWebhookHandler(p.xSignatureKey)
}

// EscapeHandler returns nil; Billplz uses the local pay page route.
func (p *Provider) EscapeHandler() http.Handler { return nil }
