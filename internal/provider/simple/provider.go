// Package simple implements a YAML-driven provider runtime for OpenMuara.
// It turns a declarative gateway.yml into a provider.Provider that can validate
// requests, verify signatures, record transactions, render templated responses,
// and drive an escape/simulation page.
package simple

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// Provider implements provider.Provider from a GatewayConfig using the
// simple runtime.
type Provider struct {
	name       string
	cfg        plugin.GatewayConfig
	runtime    plugin.SimpleRuntime
	secret     string
	store      engine.TransactionStore
	dispatcher *webhook.Dispatcher
	baseURL    string
}

// NewProvider returns an unconfigured simple provider for the given config.
func NewProvider(cfg plugin.GatewayConfig) *Provider {
	return &Provider{
		name:    cfg.Metadata.Name,
		cfg:     cfg,
		store:   engine.NewMemoryStore(),
		runtime: runtimeConfig(cfg),
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return p.name }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	if p.cfg.Runtime == nil || p.cfg.Runtime.Type != "simple" {
		return errcode.New(errcode.EConfigMissing, fmt.Sprintf("%s: runtime.type must be simple", p.name))
	}

	secret, err := resolveSecret(p.name, p.cfg.Signature, cfg)
	if err != nil {
		return err
	}
	p.secret = secret

	if p.runtime.Currency == "" {
		p.runtime.Currency = "MYR"
	}
	if p.runtime.ReferenceField == "" {
		p.runtime.ReferenceField = "reference"
	}
	if p.runtime.AmountField == "" {
		p.runtime.AmountField = "amount"
	}
	return nil
}

// Routes returns the route table declared in gateway.yml with simple-runtime handlers.
func (p *Provider) Routes() []provider.Route {
	routes := make([]provider.Route, 0, len(p.cfg.Routes))
	for _, r := range p.cfg.Routes {
		routes = append(routes, provider.Route{
			Method:  r.Method,
			Path:    r.Path,
			Handler: p.HandlerFor(r),
		})
	}
	return routes
}

// HandlerFor returns the simple-runtime handler for a declarative route.
func (p *Provider) HandlerFor(r plugin.Route) http.Handler {
	return p.handlerFor(r)
}

// ChargeHandler returns the handler for the configured charge route, or a no-op
// handler if no charge route is configured.
func (p *Provider) ChargeHandler() http.Handler {
	for _, r := range p.cfg.Routes {
		if r.Action == p.runtime.ChargeRoute && p.runtime.ChargeRoute != "" {
			return p.handlerFor(r)
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// WebhookHandler returns the first webhook route handler, or a no-op handler.
func (p *Provider) WebhookHandler() http.Handler {
	for _, r := range p.cfg.Routes {
		if isWebhookAction(r.Action) {
			return p.handlerFor(r)
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// PayloadBuilder returns a JSON payload builder using the first configured webhook
// template. If no webhook template is configured, it returns a minimal builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return p.buildPayload
}

// EscapeHandler returns the escape page handler when configured.
func (p *Provider) EscapeHandler() http.Handler {
	if p.runtime.EscapePage == nil || !p.runtime.EscapePage.Enabled {
		return nil
	}
	for _, r := range p.cfg.Routes {
		if isEscapePageAction(r.Action) {
			return p.handlerFor(r)
		}
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

func runtimeConfig(cfg plugin.GatewayConfig) plugin.SimpleRuntime {
	if cfg.Runtime != nil && cfg.Runtime.Simple != nil {
		return *cfg.Runtime.Simple
	}
	return plugin.SimpleRuntime{}
}

func (p *Provider) amount(values map[string]any) float64 {
	if p.runtime.AmountField == "charge_items" {
		var total float64
		for _, it := range itemList(values) {
			total += it.Price * float64(it.Quantity)
		}
		return total
	}
	amount, _ := floatValue(values, p.runtime.AmountField)
	return amount
}
