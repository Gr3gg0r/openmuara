// Package provider defines the common plugin contract for payment gateway
// emulations in OpenMuara.
package provider

import (
	"context"
	"net/http"
)

// Provider is the common contract implemented by every payment gateway plugin.
type Provider interface {
	// Name returns the stable provider identifier used in config and logs.
	Name() string

	// Init receives provider-specific configuration and returns an error if the
	// config is invalid. It is called once at server startup.
	Init(cfg map[string]any) error

	// Routes returns the route table this provider wants to register.
	Routes() []Route

	// ChargeHandler returns the HTTP handler for the provider's charge/payment endpoint.
	ChargeHandler() http.Handler

	// WebhookHandler returns the HTTP handler that receives provider webhooks.
	WebhookHandler() http.Handler

	// PayloadBuilder returns a function that builds a webhook payload for a transaction.
	PayloadBuilder() func(ctx context.Context, tx Transaction) ([]byte, error)

	// EscapeHandler returns an HTTP handler for the provider's escape/3-D Secure page,
	// or nil if the provider does not use escape pages.
	EscapeHandler() http.Handler
}

// VersionedProvider is an optional contract for providers that emulate more than
// one API version. The dashboard uses this to surface version selection to users.
type VersionedProvider interface {
	Provider
	// Versions returns the supported API versions, e.g. ["v1", "v2"].
	Versions() []string
	// CurrentVersion returns the version selected by the current configuration.
	CurrentVersion() string
}

// Route describes one HTTP route a provider wants to register.
type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}
