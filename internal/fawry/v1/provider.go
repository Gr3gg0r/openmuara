// Package v1 implements the legacy Fawry webhook format.
package v1

import (
	"context"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// Provider implements the Fawry v1 webhook receiver and payload builder.
type Provider struct {
	webhookSecret string
	store         engine.TransactionStore
	dispatcher    *webhook.Dispatcher
}

// NewProvider returns a new Fawry v1 provider.
func NewProvider(webhookSecret string) *Provider {
	return &Provider{webhookSecret: webhookSecret}
}

// SetStore sets the transaction store.
func (p *Provider) SetStore(s engine.TransactionStore) { p.store = s }

// SetDispatcher sets the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(d *webhook.Dispatcher) { p.dispatcher = d }

// WebhookHandler returns the v1 webhook receiver.
func (p *Provider) WebhookHandler() http.Handler { return NewWebhookHandler(p.webhookSecret) }

// PayloadBuilder returns the v1 webhook payload builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return NewPayloadBuilder(p.webhookSecret)
}
