// Package toyyibpay emulates the ToyyibPay payment gateway.
package toyyibpay

import (
	"context"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// ProviderName is the registered provider identifier.
const ProviderName = "toyyibpay"

// RegisterWith registers the ToyyibPay provider with the given registry.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
}

// Provider implements provider.Provider for the ToyyibPay gateway.
type Provider struct {
	secret          string
	defaultCategory string
	categories      *CategoryStore
	bills           *BillStore
	store           engine.TransactionStore
	baseURL         string
	dispatcher      *webhook.Dispatcher
}

// NewProvider returns an unconfigured ToyyibPay provider.
func NewProvider() *Provider {
	return &Provider{
		categories: NewCategoryStore(),
		bills:      NewBillStore(),
		store:      engine.NewMemoryStore(),
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return ProviderName }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	secret, ok := cfg["user_secret_key"].(string)
	if !ok || secret == "" {
		return errcode.New(errcode.EConfigMissing, "toyyibpay: user_secret_key is required")
	}
	p.secret = secret
	if cat, ok := cfg["category_code"].(string); ok {
		p.defaultCategory = cat
	}
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

// Routes returns the ToyyibPay route table.
func (p *Provider) Routes() []provider.Route {
	return []provider.Route{
		{Method: http.MethodPost, Path: "/index.php/api/createCategory", Handler: p.categoryCreateHandler()},
		{Method: http.MethodPost, Path: "/index.php/api/getCategoryDetails", Handler: p.categoryDetailsHandler()},
		{Method: http.MethodPost, Path: "/index.php/api/createBill", Handler: p.billCreateHandler()},
		{Method: http.MethodPost, Path: "/index.php/api/getBillTransactions", Handler: p.billTransactionsHandler()},
		{Method: http.MethodPost, Path: "/index.php/api/inactiveBill", Handler: p.billInactiveHandler()},
		{Method: http.MethodGet, Path: "/_admin/toyyibpay/pay/{billCode}", Handler: p.payPageHandler()},
		{Method: http.MethodPost, Path: "/_admin/toyyibpay/pay/{billCode}", Handler: p.payPageActionHandler()},
		{Method: http.MethodGet, Path: "/toyyibpay/return", Handler: p.returnHandler()},
		{Method: http.MethodPost, Path: "/toyyibpay/webhook", Handler: p.WebhookHandler()},
	}
}

// ChargeHandler returns a no-op handler; ToyyibPay uses the bills API.
func (p *Provider) ChargeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// PayloadBuilder returns a form-encoded callback body builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return p.buildCallbackPayload
}

// PayloadHeaders returns the callback content type header.
func (p *Provider) PayloadHeaders(_ context.Context, _ provider.Transaction) (map[string]string, error) {
	return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, nil
}

// EscapeHandler returns nil; ToyyibPay uses the local payment page.
func (p *Provider) EscapeHandler() http.Handler { return nil }
