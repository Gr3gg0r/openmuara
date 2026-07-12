package server

import (
	"context"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/factory"
)

// fakeStripeProvider is a minimal test-only provider used to satisfy config
// validation when server tests enable Stripe. Importing the real internal/stripe
// package would create an import cycle because stripe imports server, so we
// register a lightweight factory instead.
type fakeStripeProvider struct{}

func (fakeStripeProvider) Name() string                 { return "stripe" }
func (fakeStripeProvider) Init(map[string]any) error    { return nil }
func (fakeStripeProvider) Routes() []provider.Route     { return nil }
func (fakeStripeProvider) ChargeHandler() http.Handler  { return nil }
func (fakeStripeProvider) WebhookHandler() http.Handler { return nil }
func (fakeStripeProvider) EscapeHandler() http.Handler  { return nil }
func (fakeStripeProvider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }
}

func init() {
	factory.Register("stripe", func(_ map[string]any) (provider.Provider, error) {
		return fakeStripeProvider{}, nil
	})
}
