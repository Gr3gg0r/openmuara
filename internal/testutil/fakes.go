// Package testutil provides reusable helpers and fakes for OpenMuara tests.
package testutil

import (
	"context"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// FakeDispatcher records every dispatch call without sending HTTP requests.
type FakeDispatcher struct {
	Calls []FakeDispatchCall
}

// FakeDispatchCall captures one call to FakeDispatcher.Dispatch.
type FakeDispatchCall struct {
	Ref    string
	Status webhook.PaymentStatus
}

// Dispatch implements webhook.Sender.
func (f *FakeDispatcher) Dispatch(_ context.Context, ref string, status webhook.PaymentStatus) (*webhook.Attempt, error) {
	f.Calls = append(f.Calls, FakeDispatchCall{Ref: ref, Status: status})
	return &webhook.Attempt{Ref: ref, Status: webhook.AttemptStatusDelivered}, nil
}

// FakeProvider is a minimal provider.Provider implementation for router tests.
type FakeProvider struct {
	ProviderName string
	InitErr      error
	RoutesFunc   func() []provider.Route
}

// Name returns the provider identifier.
func (f *FakeProvider) Name() string { return f.ProviderName }

// Init returns the configured error, if any.
func (f *FakeProvider) Init(_ map[string]any) error { return f.InitErr }

// Routes returns the configured routes or nil.
func (f *FakeProvider) Routes() []provider.Route {
	if f.RoutesFunc != nil {
		return f.RoutesFunc()
	}
	return nil
}

// ChargeHandler returns a no-op handler.
func (f *FakeProvider) ChargeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// WebhookHandler returns a no-op handler.
func (f *FakeProvider) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// PayloadBuilder returns an empty payload.
func (f *FakeProvider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(context.Context, provider.Transaction) ([]byte, error) {
		return []byte("{}"), nil
	}
}

// EscapeHandler returns nil.
func (f *FakeProvider) EscapeHandler() http.Handler { return nil }

// SetStore satisfies the optional setter used by the CLI.
func (f *FakeProvider) SetStore(_ engine.TransactionStore) {}

// SetBaseURL satisfies the optional setter used by the CLI.
func (f *FakeProvider) SetBaseURL(_ string) {}

// SetDispatcher satisfies the optional setter used by the CLI.
func (f *FakeProvider) SetDispatcher(_ *webhook.Dispatcher) {}
