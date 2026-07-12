package fawry

import (
	"context"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	v1 "github.com/Gr3gg0r/openmuara/internal/fawry/v1"
	v2 "github.com/Gr3gg0r/openmuara/internal/fawry/v2"
	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// NewPlugin creates the Fawry plugin implementation.
func NewPlugin(merchantSecurityKey, webhookSecret string, store engine.TransactionStore, dispatcher *webhook.Dispatcher) plugin.Plugin {
	return plugin.NewBuiltinPlugin("fawry", "1.0.0", func(_ context.Context, reg *plugin.Registry) error {
		if err := reg.RegisterHandler("fawry_charge", func(_ plugin.Dependencies) (http.Handler, error) {
			return NewChargeHandler(merchantSecurityKey, store), nil
		}); err != nil {
			return err
		}
		v1p := v1.NewProvider(webhookSecret)
		v1p.SetStore(store)
		v1p.SetDispatcher(dispatcher)
		if err := reg.RegisterHandler("fawry_v1_webhook", func(_ plugin.Dependencies) (http.Handler, error) {
			return v1p.WebhookHandler(), nil
		}); err != nil {
			return err
		}
		v2p := v2.NewProvider(webhookSecret)
		v2p.SetStore(store)
		v2p.SetDispatcher(dispatcher)
		if err := reg.RegisterHandler("fawry_v2_webhook", func(_ plugin.Dependencies) (http.Handler, error) {
			return v2p.WebhookHandler(), nil
		}); err != nil {
			return err
		}
		if err := reg.RegisterHandler("fawry_escape_page", func(_ plugin.Dependencies) (http.Handler, error) {
			return NewEscapeHandler(), nil
		}); err != nil {
			return err
		}
		if err := reg.RegisterHandler("fawry_escape_action", func(_ plugin.Dependencies) (http.Handler, error) {
			return NewEscapeActionHandler(dispatcher, store), nil
		}); err != nil {
			return err
		}
		return nil
	})
}
