// Package fawry emulates the Fawry Express Checkout payment gateway.
package fawry

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	v1 "github.com/openmuara/openmuara/internal/fawry/v1"
	v2 "github.com/openmuara/openmuara/internal/fawry/v2"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// RegisterWith registers the Fawry provider with the given registry. Tests can
// use this to build isolated registries instead of relying on the global default.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
	slog.Debug("registered provider", "name", "fawry")
}

// versioned is the contract implemented by each Fawry API version.
type versioned interface {
	WebhookHandler() http.Handler
	PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error)
	SetStore(engine.TransactionStore)
	SetDispatcher(*webhook.Dispatcher)
}

// Provider implements provider.Provider for the Fawry Express Checkout gateway.
// It dispatches webhook and payload behavior to the configured API version.
type Provider struct {
	merchantCode        string
	merchantSecurityKey string
	webhookSecret       string
	store               engine.TransactionStore
	dispatcher          *webhook.Dispatcher
	version             string
	v1Impl              versioned
	v2Impl              versioned
}

// NewProvider returns an unconfigured Fawry provider with its own in-memory store.
func NewProvider() *Provider {
	return &Provider{store: engine.NewMemoryStore(), version: "v1"}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return "fawry" }

// Init validates and stores provider-specific config.
func (p *Provider) Init(cfg map[string]any) error {
	mc, ok := cfg["merchant_code"].(string)
	if !ok || mc == "" {
		return errcode.New(errcode.EConfigMissing, "fawry: merchant_code is required")
	}
	msk, ok := cfg["merchant_security_key"].(string)
	if !ok || msk == "" {
		return errcode.New(errcode.EConfigMissing, "fawry: merchant_security_key is required")
	}
	ws, ok := cfg["webhook_secret"].(string)
	if !ok || ws == "" {
		return errcode.New(errcode.EConfigMissing, "fawry: webhook_secret is required")
	}
	p.merchantCode = mc
	p.merchantSecurityKey = msk
	p.webhookSecret = ws

	version := "v1"
	if v, ok := cfg["version"].(string); ok && v != "" {
		version = v
	}
	if version != "v1" && version != "v2" {
		return errcode.New(errcode.EProviderVersionUnsupported, fmt.Sprintf("fawry: unsupported version %q", version))
	}
	p.version = version

	p.v1Impl = v1.NewProvider(ws)
	p.v2Impl = v2.NewProvider(ws)
	p.v1Impl.SetStore(p.store)
	p.v2Impl.SetStore(p.store)
	if p.dispatcher != nil {
		p.v1Impl.SetDispatcher(p.dispatcher)
		p.v2Impl.SetDispatcher(p.dispatcher)
	}
	return nil
}

// Routes returns the Fawry HTTP routes.
func (p *Provider) Routes() []provider.Route {
	charge := p.ChargeHandler()
	return []provider.Route{
		{Method: http.MethodPost, Path: "/fawry/charge", Handler: charge},
		{Method: http.MethodGet, Path: "/fawry/payment-status", Handler: NewStatusHandler(p.merchantCode, p.merchantSecurityKey, p.store)},
		{Method: http.MethodPost, Path: "/fawry/webhook", Handler: p.WebhookHandler()},
		{Method: http.MethodPost, Path: "/fawry/v1/charge", Handler: charge},
		{Method: http.MethodPost, Path: "/fawry/v2/charge", Handler: charge},
		{Method: http.MethodPost, Path: "/fawry/v1/webhook", Handler: p.v1Impl.WebhookHandler()},
		{Method: http.MethodPost, Path: "/fawry/v2/webhook", Handler: p.v2Impl.WebhookHandler()},
		{Method: http.MethodGet, Path: "/_admin/fawry-escape", Handler: p.EscapeHandler()},
		{Method: http.MethodPost, Path: "/_admin/fawry-escape", Handler: NewEscapeActionHandler(p.dispatcher, p.store)},
	}
}

// ChargeHandler returns the charge handler backed by the provider's store.
func (p *Provider) ChargeHandler() http.Handler {
	return NewChargeHandler(p.merchantSecurityKey, p.store)
}

// WebhookHandler returns the webhook handler for the configured version.
func (p *Provider) WebhookHandler() http.Handler {
	if p.version == "v2" {
		return p.v2Impl.WebhookHandler()
	}
	return p.v1Impl.WebhookHandler()
}

// PayloadBuilder returns the payload builder for the configured version.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	if p.version == "v2" {
		return p.v2Impl.PayloadBuilder()
	}
	return p.v1Impl.PayloadBuilder()
}

// EscapeHandler returns the Fawry 3-D Secure escape page handler.
func (p *Provider) EscapeHandler() http.Handler {
	return NewEscapeHandler()
}

// MerchantCode returns the configured merchant code.
func (p *Provider) MerchantCode() string { return p.merchantCode }

// MerchantSecurityKey returns the configured merchant security key.
func (p *Provider) MerchantSecurityKey() string { return p.merchantSecurityKey }

// Store returns the provider's transaction store.
func (p *Provider) Store() engine.TransactionStore { return p.store }

// Dispatcher returns the provider's webhook dispatcher.
func (p *Provider) Dispatcher() *webhook.Dispatcher { return p.dispatcher }

// V1Impl returns the Fawry v1 implementation.
func (p *Provider) V1Impl() versioned { return p.v1Impl }

// V2Impl returns the Fawry v2 implementation.
func (p *Provider) V2Impl() versioned { return p.v2Impl }

// Versions returns the supported Fawry API versions.
func (p *Provider) Versions() []string { return []string{"v1", "v2"} }

// CurrentVersion returns the configured Fawry API version.
func (p *Provider) CurrentVersion() string { return p.version }

// VerifyWebhookSignature checks the messageSignature embedded in a Fawry webhook payload.
func (p *Provider) VerifyWebhookSignature(payload []byte, _ map[string]string) (bool, error) {
	if p.version == "v2" {
		var body webhook.FawryV2Payload
		if err := json.Unmarshal(payload, &body); err != nil {
			return false, errcode.Wrap(errcode.EInvalidRequest, "failed to decode fawry v2 webhook payload", err)
		}
		if body.MessageSignature == "" {
			return false, errcode.New(errcode.ESignatureMissing, "fawry v2 webhook message signature missing")
		}
		return webhook.NewHMACSigner(p.webhookSecret).Verify(body, body.MessageSignature)
	}

	var body v1.WebhookBody
	if err := json.Unmarshal(payload, &body); err != nil {
		return false, errcode.Wrap(errcode.EInvalidRequest, "failed to decode fawry v1 webhook payload", err)
	}
	if body.MessageSignature == "" {
		return false, errcode.New(errcode.ESignatureMissing, "fawry v1 webhook message signature missing")
	}
	return v1.VerifySignature(body, p.webhookSecret), nil
}

// SetDispatcher configures the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(dispatcher *webhook.Dispatcher) {
	p.dispatcher = dispatcher
	if p.v1Impl != nil {
		p.v1Impl.SetDispatcher(dispatcher)
	}
	if p.v2Impl != nil {
		p.v2Impl.SetDispatcher(dispatcher)
	}
}

// SetStore replaces the provider's transaction store.
func (p *Provider) SetStore(store engine.TransactionStore) {
	p.store = store
	if p.v1Impl != nil {
		p.v1Impl.SetStore(store)
	}
	if p.v2Impl != nil {
		p.v2Impl.SetStore(store)
	}
}
