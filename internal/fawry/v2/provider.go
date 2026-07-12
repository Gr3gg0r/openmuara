// Package v2 implements the Fawry V2 server notification format.
package v2

import (
	"context"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// Provider implements the Fawry v2 webhook receiver and payload builder.
type Provider struct {
	webhookSecret string
	store         engine.TransactionStore
	dispatcher    *webhook.Dispatcher
}

// NewProvider returns a new Fawry v2 provider.
func NewProvider(webhookSecret string) *Provider {
	return &Provider{webhookSecret: webhookSecret}
}

// SetStore sets the transaction store.
func (p *Provider) SetStore(s engine.TransactionStore) { p.store = s }

// SetDispatcher sets the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(d *webhook.Dispatcher) { p.dispatcher = d }

// WebhookHandler returns the v2 webhook receiver.
func (p *Provider) WebhookHandler() http.Handler { return NewWebhookHandler(p.webhookSecret) }

// PayloadBuilder returns the v2 webhook payload builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return NewPayloadBuilder(p.webhookSecret, p.store)
}
