package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/webhook"
)

// PaymentIntentStore persists PaymentIntents by ID.
type PaymentIntentStore interface {
	Save(id string, pi *PaymentIntent)
	Load(id string) (*PaymentIntent, bool)
}

// MemoryPaymentIntentStore is an in-memory PaymentIntentStore.
type MemoryPaymentIntentStore struct {
	mu   sync.RWMutex
	data map[string]*PaymentIntent
}

// NewMemoryPaymentIntentStore creates a new in-memory PaymentIntent store.
func NewMemoryPaymentIntentStore() *MemoryPaymentIntentStore {
	return &MemoryPaymentIntentStore{data: make(map[string]*PaymentIntent)}
}

// Save stores a PaymentIntent by ID.
func (s *MemoryPaymentIntentStore) Save(id string, pi *PaymentIntent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = pi
}

// Load retrieves a PaymentIntent by ID.
func (s *MemoryPaymentIntentStore) Load(id string) (*PaymentIntent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pi, ok := s.data[id]
	return pi, ok
}

// NewCreatePaymentIntentHandler returns POST /v1/payment_intents.
func NewCreatePaymentIntentHandler(store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		var req PaymentIntentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_json", "", "invalid JSON body")
			return
		}

		if err := validatePaymentIntentRequest(req); err != nil {
			writeStripeValidationError(w, err)
			return
		}

		pi := buildPaymentIntent(req)
		store.Save(pi.ID, pi)

		if _, _, err := ledger.CreateOrGet(paymentIntentToLedgerTransaction(r.Context(), pi)); err != nil {
			writeStripeInvalidRequestError(w, http.StatusInternalServerError, "", "", "failed to record transaction")
			return
		}

		dispatchPaymentIntentWebhook(r.Context(), dispatcher, pi.ID, webhook.PaymentStatusNew)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(pi) // #nosec G117 -- emulates Stripe API field name client_secret
	}
}

// NewGetPaymentIntentHandler returns GET /v1/payment_intents/{id}.
func NewGetPaymentIntentHandler(store PaymentIntentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := paymentIntentIDFromPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "id", "payment intent id is required")
			return
		}

		pi, ok := store.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "id", "payment intent not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(pi) // #nosec G117 -- emulates Stripe API field name client_secret
	}
}

// NewConfirmPaymentIntentHandler returns POST /v1/payment_intents/{id}/confirm.
func NewConfirmPaymentIntentHandler(store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := paymentIntentIDFromConfirmPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "id", "payment intent id is required")
			return
		}

		pi, ok := store.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "id", "payment intent not found")
			return
		}

		var req PaymentIntentConfirmRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_json", "", "invalid JSON body")
			return
		}

		switch {
		case isCardPaymentMethod(req.PaymentMethod):
			confirmCardPaymentIntent(w, r, pi, store, ledger, dispatcher, req.PaymentMethod)
		case isFPXPaymentMethod(req.PaymentMethod):
			confirmFPXPaymentIntent(w, r, pi, store, baseURL, req.PaymentMethod)
		default:
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "payment_method", fmt.Sprintf("no such payment_method: %s", req.PaymentMethod))
		}
	}
}

// NewCancelPaymentIntentHandler returns POST /v1/payment_intents/{id}/cancel.
func NewCancelPaymentIntentHandler(store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := paymentIntentIDFromCancelPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "id", "payment intent id is required")
			return
		}

		pi, ok := store.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "id", "payment intent not found")
			return
		}

		if pi.Status != "requires_confirmation" && pi.Status != "requires_action" {
			writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", fmt.Sprintf("payment intent status is %s", pi.Status))
			return
		}

		pi.Status = "canceled"
		store.Save(pi.ID, pi)

		if err := updateLedgerStatus(ledger, pi.ID, engine.TransactionStatusUnpaid); err != nil {
			writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", errcode.Message(err))
			return
		}

		dispatchPaymentIntentWebhook(r.Context(), dispatcher, pi.ID, webhook.PaymentStatusUnpaid)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(pi) // #nosec G117 -- emulates Stripe API field name client_secret
	}
}

func confirmCardPaymentIntent(w http.ResponseWriter, r *http.Request, pi *PaymentIntent, store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher, paymentMethod string) {
	if pi.Status != "requires_confirmation" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", fmt.Sprintf("payment intent status is %s", pi.Status))
		return
	}

	pi.Status = "succeeded"
	pi.PaymentMethod = paymentMethod
	store.Save(pi.ID, pi)

	if err := updateLedgerStatus(ledger, pi.ID, engine.TransactionStatusPaid); err != nil {
		writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", errcode.Message(err))
		return
	}

	dispatchPaymentIntentWebhook(r.Context(), dispatcher, pi.ID, webhook.PaymentStatusPaid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pi) // #nosec G117 -- emulates Stripe API field name client_secret
}

func confirmFPXPaymentIntent(w http.ResponseWriter, r *http.Request, pi *PaymentIntent, store PaymentIntentStore, baseURL string, paymentMethod string) {
	if pi.Status != "requires_confirmation" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", fmt.Sprintf("payment intent status is %s", pi.Status))
		return
	}

	pi.Status = "requires_action"
	pi.PaymentMethod = paymentMethod
	pi.NextAction = &PaymentIntentNextAction{
		Type: "redirect_to_url",
		RedirectToURL: &PaymentIntentRedirectToURL{
			URL: resolveBaseURL(baseURL, r) + "/_admin/stripe/payment_intent/" + pi.ID,
		},
	}
	store.Save(pi.ID, pi)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(pi) // #nosec G117 -- emulates Stripe API field name client_secret
}

func buildPaymentIntent(req PaymentIntentRequest) *PaymentIntent {
	id := "pi_test_" + strings.ReplaceAll(uuid.Must(uuid.NewRandom()).String(), "-", "")
	return &PaymentIntent{
		ID:                 id,
		Object:             "payment_intent",
		Amount:             req.Amount,
		Currency:           strings.ToLower(req.Currency),
		Status:             "requires_confirmation",
		ClientSecret:       id + "_secret_test",
		PaymentMethodTypes: normalizePaymentIntentMethodTypes(req.PaymentMethodTypes),
		Metadata:           req.Metadata,
		Livemode:           false,
		ConfirmationMethod: "automatic",
	}
}

func paymentIntentToLedgerTransaction(ctx context.Context, pi *PaymentIntent) engine.Transaction {
	return engine.NewTransaction(engine.Transaction{
		Provider:    ProviderName,
		Type:        "payment_intent",
		Amount:      float64(pi.Amount) / 100.0,
		Currency:    strings.ToUpper(pi.Currency),
		Status:      engine.TransactionStatusNew,
		CustomerRef: pi.PaymentMethod,
		Reference:   pi.ID,
		TraceID:     httputil.TraceIDFromContext(ctx),
	})
}

func dispatchPaymentIntentWebhook(ctx context.Context, dispatcher Dispatcher, ref string, status webhook.PaymentStatus) {
	if isNilDispatcher(dispatcher) {
		return
	}
	if _, err := dispatcher.Dispatch(ctx, ref, status); err != nil {
		audit.FromContext(ctx).Log(ctx, "payment_intent.webhook_failed", "payment_intent", ref, "", err.Error())
	}
}

func isNilDispatcher(d Dispatcher) bool {
	if d == nil {
		return true
	}
	if v, ok := d.(*webhook.Dispatcher); ok {
		return v == nil
	}
	return false
}

func paymentIntentIDFromPath(path string) string {
	return strings.TrimPrefix(path, "/v1/payment_intents/")
}

func paymentIntentIDFromConfirmPath(path string) string {
	prefix := "/v1/payment_intents/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(path, prefix), "/confirm")
}

func paymentIntentIDFromCancelPath(path string) string {
	prefix := "/v1/payment_intents/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(path, prefix), "/cancel")
}

func isPaymentIntentRef(ref string) bool {
	return strings.HasPrefix(ref, "pi_test_")
}

// PaymentIntentEvent mirrors a subset of the Stripe event object for PaymentIntents.
type PaymentIntentEvent struct {
	ID     string                 `json:"id"`
	Object string                 `json:"object"`
	Type   string                 `json:"type"`
	Data   PaymentIntentEventData `json:"data"`
}

// PaymentIntentEventData wraps the PaymentIntent inside a Stripe event.
type PaymentIntentEventData struct {
	Object *PaymentIntent `json:"object"`
}

func paymentIntentEventType(status string) string {
	switch status {
	case "PAID", "complete":
		return "payment_intent.succeeded"
	case "UNPAID", "canceled", "expired":
		return "payment_intent.canceled"
	case "NEW":
		return "payment_intent.created"
	default:
		return "payment_intent.created"
	}
}

func buildPaymentIntentEventPayload(pi *PaymentIntent, eventType string) ([]byte, error) {
	event := PaymentIntentEvent{
		ID:     "evt_test_" + pi.ID,
		Object: "event",
		Type:   eventType,
		Data:   PaymentIntentEventData{Object: pi.Clone()},
	}
	return json.Marshal(event)
}
