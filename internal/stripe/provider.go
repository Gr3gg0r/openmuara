package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/ui"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
	"github.com/spf13/viper"
)

// ProviderName is the registered provider identifier.
const ProviderName = "stripe"

// RegisterWith registers the Stripe provider with the given registry. Tests can
// use this to build isolated registries instead of relying on the global default.
func RegisterWith(r *provider.Registry) {
	r.Register(NewProvider())
}

// Provider implements the generic provider interface for Stripe Checkout.
type Provider struct {
	publishableKey string
	secretKey      string
	webhookSecret  string
	sessions       SessionStore
	paymentIntents PaymentIntentStore
	ledger         engine.TransactionStore
	dispatcher     *webhook.Dispatcher
	baseURL        string
	configPath     string
	enabledEvents  []string
}

// NewProvider returns an unconfigured Stripe provider.
func NewProvider() *Provider {
	return &Provider{
		sessions:       NewMemorySessionStore(),
		paymentIntents: NewMemoryPaymentIntentStore(),
		ledger:         engine.NewMemoryStore(),
		enabledEvents:  defaultEnabledEvents(),
	}
}

func defaultEnabledEvents() []string {
	return []string{
		"checkout.session.completed",
		"checkout.session.expired",
		"payment_intent.created",
		"payment_intent.succeeded",
		"payment_intent.canceled",
	}
}

// Name returns the provider identifier.
func (p *Provider) Name() string { return ProviderName }

// Init validates and stores provider-specific configuration.
func (p *Provider) Init(cfg map[string]any) error {
	pk, ok := cfg["publishable_key"].(string)
	if !ok || pk == "" {
		return errcode.New(errcode.EConfigMissing, "stripe: publishable_key is required")
	}
	sk, ok := cfg["secret_key"].(string)
	if !ok || sk == "" {
		return errcode.New(errcode.EConfigMissing, "stripe: secret_key is required")
	}
	ws, ok := cfg["webhook_secret"].(string)
	if !ok || ws == "" {
		return errcode.New(errcode.EConfigMissing, "stripe: webhook_secret is required")
	}
	p.publishableKey = pk
	p.secretKey = sk
	p.webhookSecret = ws

	if events, ok := cfg["enabled_events"].([]any); ok {
		p.enabledEvents = toStringSlice(events)
	}
	return nil
}

// SetBaseURL configures the base URL used for session payment links.
func (p *Provider) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

// SetDispatcher configures the outbound webhook dispatcher.
func (p *Provider) SetDispatcher(dispatcher *webhook.Dispatcher) {
	p.dispatcher = dispatcher
	if dispatcher == nil {
		return
	}
	if len(p.enabledEvents) > 0 {
		dispatcher.EnabledEvents = p.enabledEvents
	}
}

// SetStore replaces the provider's transaction ledger.
func (p *Provider) SetStore(store engine.TransactionStore) {
	p.ledger = store
}

// SetConfigPath configures the path used to persist provider settings.
func (p *Provider) SetConfigPath(path string) {
	p.configPath = path
}

// Routes returns the Stripe route table.
func (p *Provider) Routes() []provider.Route {
	checkoutPayPage := NewCheckoutSessionPayPageHandler(p.sessions)
	checkoutPayAction := NewCheckoutSessionPayActionHandler(
		p.sessions, p.ledger, p.dispatcherIfSet(),
	)
	retrieveSession := p.sessionRetrieveHandler()
	createPI := NewCreatePaymentIntentHandler(
		p.paymentIntents, p.ledger, p.dispatcherIfSet(),
	)
	getPI := NewGetPaymentIntentHandler(p.paymentIntents)
	confirmPI := NewConfirmPaymentIntentHandler(
		p.paymentIntents, p.ledger, p.dispatcherIfSet(), p.baseURL,
	)
	cancelPI := NewCancelPaymentIntentHandler(
		p.paymentIntents, p.ledger, p.dispatcherIfSet(),
	)
	piAdminPage := NewPaymentIntentAdminPageHandler(p.paymentIntents)
	piAdminAction := NewPaymentIntentAdminActionHandler(
		p.paymentIntents, p.ledger, p.dispatcherIfSet(),
	)
	failSim := NewFailureSimulationHandler(p.sessions, p.ledger)
	cancelSim := NewCancelSimulationHandler(p.sessions, p.ledger)

	return []provider.Route{
		{Method: http.MethodPost, Path: "/v1/checkout/sessions", Handler: p.ChargeHandler()},
		{Method: http.MethodGet, Path: "/v1/checkout/sessions/{id}/pay", Handler: checkoutPayPage},
		{Method: http.MethodPost, Path: "/v1/checkout/sessions/{id}/pay", Handler: checkoutPayAction},
		{Method: http.MethodGet, Path: "/v1/checkout/sessions/", Handler: retrieveSession},
		{Method: http.MethodPost, Path: "/v1/payment_intents", Handler: createPI},
		{Method: http.MethodGet, Path: "/v1/payment_intents/", Handler: getPI},
		{Method: http.MethodPost, Path: "/v1/payment_intents/{id}/confirm", Handler: confirmPI},
		{Method: http.MethodPost, Path: "/v1/payment_intents/{id}/cancel", Handler: cancelPI},
		{Method: http.MethodGet, Path: "/_admin/stripe/payment_intent/{id}", Handler: piAdminPage},
		{Method: http.MethodPost, Path: "/_admin/stripe/payment_intent/{id}", Handler: piAdminAction},
		{Method: http.MethodPost, Path: "/_admin/stripe/success", Handler: p.successSimulationHandler()},
		{Method: http.MethodPost, Path: "/_admin/stripe/fail", Handler: failSim},
		{Method: http.MethodPost, Path: "/_admin/stripe/cancel", Handler: cancelSim},
	}
}

// ChargeHandler returns the POST /v1/checkout/sessions handler.
func (p *Provider) ChargeHandler() http.Handler {
	return NewCreateCheckoutSessionHandler(p.sessions, p.ledger, p.baseURL)
}

func (p *Provider) sessionRetrieveHandler() http.Handler {
	return NewGetCheckoutSessionHandler(p.sessions)
}

func (p *Provider) dispatcherIfSet() Dispatcher {
	if p.dispatcher == nil {
		return nil
	}
	return p.dispatcher
}

func (p *Provider) successSimulationHandler() http.Handler {
	return NewSuccessSimulationHandler(p.sessions, p.ledger, p.dispatcherIfSet())
}

// WebhookConfigPageHandler returns the GET /_admin/stripe/webhooks handler.
func (p *Provider) WebhookConfigPageHandler() http.Handler {
	return p.webhookConfigPageHandler()
}

// WebhookConfigSaveHandler returns the POST /_admin/stripe/webhooks handler.
func (p *Provider) WebhookConfigSaveHandler() http.Handler {
	return p.webhookConfigSaveHandler()
}

func (p *Provider) webhookConfigPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		data := ui.StripeWebhooksPageData{
			URL:           p.webhookURL(),
			WebhookSecret: p.webhookSecret,
			Events:        p.webhookEventOptions(),
		}
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}

		_ = ui.ServeStripeWebhooksPage(w, data)
	}
}

func (p *Provider) webhookConfigSaveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_request", "", "invalid form")
			return
		}

		url := r.FormValue("webhook_url")
		events := r.Form["enabled_events"]
		if len(events) == 0 {
			events = []string{}
		}

		p.enabledEvents = events
		if p.dispatcher != nil {
			p.dispatcher.URL = url
			p.dispatcher.EnabledEvents = events
		}

		if err := p.saveWebhookConfig(url, events); err != nil {
			writeStripeInvalidRequestError(w, http.StatusInternalServerError, "", "", "failed to save config")
			return
		}

		http.Redirect(w, r, "/_admin/stripe/webhooks", http.StatusSeeOther)
	}
}

func (p *Provider) webhookURL() string {
	if p.dispatcher != nil {
		return p.dispatcher.URL
	}
	return ""
}

func (p *Provider) webhookEventOptions() []ui.StripeWebhookEventOption {
	all := defaultEnabledEvents()
	enabled := make(map[string]bool, len(p.enabledEvents))
	for _, e := range p.enabledEvents {
		enabled[e] = true
	}

	options := make([]ui.StripeWebhookEventOption, 0, len(all))
	for _, e := range all {
		options = append(options, ui.StripeWebhookEventOption{
			Name:    e,
			Checked: enabled[e],
		})
	}
	return options
}

func (p *Provider) saveWebhookConfig(url string, events []string) error {
	if p.configPath == "" {
		return nil
	}
	v := viper.New()
	v.SetConfigFile(p.configPath)
	if err := v.ReadInConfig(); err != nil {
		return errcode.Wrap(errcode.EConfigInvalid, "read config", err)
	}
	v.Set("providers.stripe.config.webhook_url", url)
	v.Set("providers.stripe.config.enabled_events", events)
	if err := v.WriteConfig(); err != nil {
		return errcode.Wrap(errcode.EConfigInvalid, "write config", err)
	}
	return nil
}

// WebhookHandler returns a no-op handler; Stripe outbound webhooks are dispatched
// by the generic webhook dispatcher using PayloadBuilder().
func (p *Provider) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// PayloadBuilder returns a Stripe checkout.session.completed event payload builder.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		return p.buildPayload(tx.Reference, tx.Status)
	}
}

// PayloadEventType returns the Stripe event type for a transaction reference and status.
func (p *Provider) PayloadEventType(ref, status string) string {
	if isPaymentIntentRef(ref) {
		return paymentIntentEventType(status)
	}
	return eventTypeForStatus(status)
}

// PayloadHeaders returns the Stripe-Signature header for a webhook payload.
func (p *Provider) PayloadHeaders(_ context.Context, tx provider.Transaction) (map[string]string, error) {
	payload, err := p.buildPayload(tx.Reference, tx.Status)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Stripe-Signature": SignPayload(payload, p.webhookSecret, time.Now()),
	}, nil
}

// VerifyWebhookSignature checks the Stripe-Signature header on an outgoing webhook.
func (p *Provider) VerifyWebhookSignature(payload []byte, headers map[string]string) (bool, error) {
	h := headers["Stripe-Signature"]
	if h == "" {
		return false, nil
	}
	return true, VerifySignature(payload, h, p.webhookSecret)
}

// EscapeHandler returns nil; Stripe Checkout does not use a local escape page.
// The checkout page is served at the session URL instead.
func (p *Provider) EscapeHandler() http.Handler { return nil }

func (p *Provider) buildPayload(ref, status string) ([]byte, error) {
	if isPaymentIntentRef(ref) {
		pi, ok := p.paymentIntents.Load(ref)
		if !ok {
			return nil, errcode.New(errcode.ETransactionNotFound, fmt.Sprintf("stripe: payment intent not found for reference %q", ref))
		}
		return buildPaymentIntentEventPayload(pi, paymentIntentEventType(status))
	}

	session, ok := p.sessions.Load(ref)
	if ok {
		return p.buildEventPayload(session, status)
	}

	// Fallback for transactions that are not stored as Checkout Sessions.
	tx, ok, err := p.ledger.GetByReference(ref)
	if err != nil {
		return nil, errcode.Wrap(errcode.EInternal, fmt.Sprintf("stripe: lookup transaction for ref %q", ref), err)
	}
	if !ok {
		return nil, errcode.New(errcode.ETransactionNotFound, "stripe: session or transaction not found for reference "+ref)
	}
	return p.buildEventPayload(sessionFromTransaction(tx), status)
}

func (p *Provider) buildEventPayload(session *CheckoutSession, status string) ([]byte, error) {
	payload := session.Clone()
	payload.Status = mapSessionStatus(status)
	payload.PaymentStatus = mapPaymentStatus(status)

	// Use a deterministic event ID derived from the session so repeated
	// PayloadBuilder/PayloadHeaders calls produce the same bytes for the
	// same transaction (required for signature verification).
	event := Event{
		ID:     "evt_test_" + session.ID,
		Object: "event",
		Type:   eventTypeForStatus(status),
		Data:   EventData{Object: payload},
	}

	return json.Marshal(event)
}

func eventTypeForStatus(status string) string {
	if isCanceledStatus(status) {
		return "checkout.session.expired"
	}
	return "checkout.session.completed"
}

func sessionFromTransaction(tx engine.Transaction) *CheckoutSession {
	status := "open"
	paymentStatus := "unpaid"
	switch tx.Status {
	case engine.TransactionStatusPaid:
		status = "complete"
		paymentStatus = "paid"
	case engine.TransactionStatusUnpaid:
		status = "expired"
	}

	return &CheckoutSession{
		ID:                tx.Reference,
		Object:            "checkout.session",
		AmountTotal:       int64(tx.Amount * 100),
		Currency:          strings.ToLower(tx.Currency),
		CustomerEmail:     tx.CustomerRef,
		Mode:              "payment",
		PaymentStatus:     paymentStatus,
		Status:            status,
		ClientReferenceID: tx.Reference,
	}
}

func mapSessionStatus(status string) string {
	switch status {
	case "PAID", "complete":
		return "complete"
	case "UNPAID", "canceled", "expired":
		return "expired"
	default:
		return "open"
	}
}

func mapPaymentStatus(status string) string {
	switch status {
	case "PAID", "complete":
		return "paid"
	default:
		return "unpaid"
	}
}

func isCanceledStatus(status string) bool {
	switch status {
	case "UNPAID", "canceled", "expired":
		return true
	default:
		return false
	}
}

func toStringSlice(v []any) []string {
	out := make([]string, 0, len(v))
	for _, item := range v {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
