package stripe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/server"
	"github.com/Gr3gg0r/openmuara/internal/ui"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// NewPaymentIntentAdminPageHandler returns GET /_admin/stripe/payment_intent/{id}.
func NewPaymentIntentAdminPageHandler(store PaymentIntentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := paymentIntentIDFromAdminPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "id", "payment intent id is required")
			return
		}

		pi, ok := store.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "id", "payment intent not found")
			return
		}

		data := paymentIntentPageData(pi)
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}

		_ = ui.ServeStripePaymentIntentPage(w, data)
	}
}

// NewPaymentIntentAdminActionHandler returns POST /_admin/stripe/payment_intent/{id}.
func NewPaymentIntentAdminActionHandler(store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		if !validatePaymentIntentCSRF(r) {
			writeStripeInvalidRequestError(w, http.StatusForbidden, "csrf_token_invalid", "", "csrf token missing or invalid")
			return
		}

		if err := r.ParseForm(); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_request", "", "invalid form")
			return
		}

		id := paymentIntentIDFromAdminPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "id", "payment intent id is required")
			return
		}

		pi, ok := store.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "id", "payment intent not found")
			return
		}

		action := r.FormValue("action")
		switch action {
		case "confirm":
			adminConfirmPaymentIntent(w, r, pi, store, ledger, dispatcher)
		case "cancel":
			adminCancelPaymentIntent(w, r, pi, store, ledger, dispatcher)
		default:
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_invalid", "action", "action must be confirm or cancel")
		}
	}
}

func adminConfirmPaymentIntent(w http.ResponseWriter, r *http.Request, pi *PaymentIntent, store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher) {
	if pi.Status != "requires_confirmation" && pi.Status != "requires_action" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", fmt.Sprintf("payment intent status is %s", pi.Status))
		return
	}

	pi.Status = "succeeded"
	pi.NextAction = nil
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

func adminCancelPaymentIntent(w http.ResponseWriter, r *http.Request, pi *PaymentIntent, store PaymentIntentStore, ledger engine.TransactionStore, dispatcher Dispatcher) {
	if pi.Status != "requires_confirmation" && pi.Status != "requires_action" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "payment_intent_unexpected_state", "", fmt.Sprintf("payment intent status is %s", pi.Status))
		return
	}

	pi.Status = "canceled"
	pi.NextAction = nil
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

func validatePaymentIntentCSRF(r *http.Request) bool {
	expected, ok := server.CSRFTokenFromCookie(r)
	if !ok {
		return false
	}
	actual := server.CSRFRequestToken(r)
	return actual != "" && strings.EqualFold(expected, actual)
}

func paymentIntentIDFromAdminPath(path string) string {
	prefix := "/_admin/stripe/payment_intent/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimPrefix(path, prefix)
}

func paymentIntentPageData(pi *PaymentIntent) ui.StripePaymentIntentPageData {
	showCard := false
	showFPX := false
	for _, t := range pi.PaymentMethodTypes {
		switch t {
		case PaymentMethodTypeCard:
			showCard = true
		case PaymentMethodTypeFPX:
			showFPX = true
		}
	}

	return ui.StripePaymentIntentPageData{
		ID:                 pi.ID,
		Amount:             pi.Amount,
		Currency:           pi.Currency,
		AmountDisplay:      formatAmount(pi.Amount, pi.Currency),
		Status:             pi.Status,
		PaymentMethodTypes: pi.PaymentMethodTypes,
		ShowCard:           showCard,
		ShowFPX:            showFPX,
		FPXBanks:           uiFPXBanks(),
	}
}
