// Package stripe emulates the Stripe Checkout payment gateway.
package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/audit"
	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/server"
	"github.com/Gr3gg0r/openmuara/internal/ui"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
	"github.com/google/uuid"
)

// SessionStore persists checkout sessions by ID.
type SessionStore interface {
	Save(id string, session *CheckoutSession)
	Load(id string) (*CheckoutSession, bool)
}

// MemorySessionStore is an in-memory SessionStore.
type MemorySessionStore struct {
	data map[string]*CheckoutSession
}

// NewMemorySessionStore creates a new in-memory session store.
func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{data: make(map[string]*CheckoutSession)}
}

// Save stores a session by ID.
func (s *MemorySessionStore) Save(id string, session *CheckoutSession) {
	s.data[id] = session
}

// Load retrieves a session by ID.
func (s *MemorySessionStore) Load(id string) (*CheckoutSession, bool) {
	session, ok := s.data[id]
	return session, ok
}

// NewCreateCheckoutSessionHandler returns the POST /v1/checkout/sessions handler.
func NewCreateCheckoutSessionHandler(sessions SessionStore, ledger engine.TransactionStore, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		var req CreateCheckoutSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_json", "", "invalid JSON body")
			return
		}

		if err := validateCreateRequest(req); err != nil {
			writeStripeValidationError(w, err)
			return
		}

		session := buildSession(req, resolveBaseURL(baseURL, r))
		sessions.Save(session.ID, session)

		if _, _, err := ledger.CreateOrGet(toLedgerTransaction(r.Context(), session, req)); err != nil {
			writeStripeInvalidRequestError(w, http.StatusInternalServerError, "", "", "failed to record transaction")
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "charge.created", "checkout_session", session.ID, audit.JSON(req), "ok")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(session)
	}
}

// NewGetCheckoutSessionHandler returns the GET /v1/checkout/sessions/:id handler.
func NewGetCheckoutSessionHandler(sessions SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := sessionIDFromRetrievePath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "session_id", "session id is required")
			return
		}

		session, ok := sessions.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "session_id", "session not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(session)
	}
}

// NewCheckoutSessionPayPageHandler returns the GET /v1/checkout/sessions/{id}/pay handler.
func NewCheckoutSessionPayPageHandler(sessions SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := sessionIDFromPayPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "session_id", "session id is required")
			return
		}

		session, ok := sessions.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "session_id", "session not found")
			return
		}

		data := checkoutPageData(session)
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}

		_ = ui.ServeStripeCheckoutPage(w, data)
	}
}

// NewCheckoutSessionPayActionHandler returns the POST /v1/checkout/sessions/{id}/pay handler.
func NewCheckoutSessionPayActionHandler(sessions SessionStore, ledger engine.TransactionStore, dispatcher Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeStripeInvalidRequestError(w, http.StatusMethodNotAllowed, "", "", "method not allowed")
			return
		}

		id := sessionIDFromPayPath(r.URL.Path)
		if id == "" {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_missing", "session_id", "session id is required")
			return
		}

		if !validateCheckoutCSRF(r) {
			writeStripeInvalidRequestError(w, http.StatusForbidden, "csrf_token_invalid", "", "csrf token missing or invalid")
			return
		}

		if err := r.ParseForm(); err != nil {
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "invalid_request", "", "invalid form")
			return
		}

		session, ok := sessions.Load(id)
		if !ok {
			writeStripeInvalidRequestError(w, http.StatusNotFound, "resource_missing", "session_id", "session not found")
			return
		}

		action := r.FormValue("action")
		switch action {
		case "confirm":
			completeCheckoutSession(w, r, session, sessions, ledger, dispatcher, webhook.PaymentStatusPaid)
		case "cancel":
			expireCheckoutSession(w, r, session, sessions, ledger, dispatcher, webhook.PaymentStatusUnpaid)
		default:
			writeStripeInvalidRequestError(w, http.StatusBadRequest, "parameter_invalid", "action", "action must be confirm or cancel")
		}
	}
}

func validateCheckoutCSRF(r *http.Request) bool {
	expected, ok := server.CSRFTokenFromCookie(r)
	if !ok {
		return false
	}
	actual := server.CSRFRequestToken(r)
	return actual != "" && strings.EqualFold(expected, actual)
}

func completeCheckoutSession(w http.ResponseWriter, r *http.Request, session *CheckoutSession, sessions SessionStore, ledger engine.TransactionStore, dispatcher Dispatcher, status webhook.PaymentStatus) {
	if session.Status != "open" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "session_invalid", "", fmt.Sprintf("session status is %s", session.Status))
		return
	}

	session.Status = "complete"
	session.PaymentStatus = "paid"
	sessions.Save(session.ID, session)

	if err := updateLedgerStatus(ledger, session.ID, engine.TransactionStatusPaid); err != nil {
		writeStripeInvalidRequestError(w, http.StatusConflict, "session_invalid", "", errcode.Message(err))
		return
	}

	dispatchCheckoutWebhook(r.Context(), dispatcher, session.ID, status)
	// #nosec G710 -- redirect target is caller-supplied success_url in emulation
	http.Redirect(w, r, session.SuccessURL, http.StatusSeeOther)
}

func expireCheckoutSession(w http.ResponseWriter, r *http.Request, session *CheckoutSession, sessions SessionStore, ledger engine.TransactionStore, dispatcher Dispatcher, status webhook.PaymentStatus) {
	if session.Status != "open" {
		writeStripeInvalidRequestError(w, http.StatusConflict, "session_invalid", "", fmt.Sprintf("session status is %s", session.Status))
		return
	}

	session.Status = "expired"
	session.PaymentStatus = "unpaid"
	sessions.Save(session.ID, session)

	if err := updateLedgerStatus(ledger, session.ID, engine.TransactionStatusUnpaid); err != nil {
		writeStripeInvalidRequestError(w, http.StatusConflict, "session_invalid", "", errcode.Message(err))
		return
	}

	dispatchCheckoutWebhook(r.Context(), dispatcher, session.ID, status)
	// #nosec G710 -- redirect target is caller-supplied cancel_url in emulation
	http.Redirect(w, r, session.CancelURL, http.StatusSeeOther)
}

func updateLedgerStatus(ledger engine.TransactionStore, ref string, target engine.TransactionStatus) error {
	tx, ok, err := ledger.GetByReference(ref)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "failed to lookup transaction", err)
	}
	if !ok {
		return errcode.New(errcode.ETransactionNotFound, "transaction not found")
	}
	if err := engine.Transition(&tx, target); err != nil {
		return err
	}
	// Clear the idempotency key so the update is applied rather than
	// short-circuited by an existing idempotency mapping.
	tx.IdempotencyKey = ""
	if _, _, err := ledger.CreateOrGet(tx); err != nil {
		return errcode.Wrap(errcode.EInternal, "failed to update transaction", err)
	}
	return nil
}

func dispatchCheckoutWebhook(ctx context.Context, dispatcher Dispatcher, ref string, status webhook.PaymentStatus) {
	if dispatcher == nil {
		return
	}
	if _, err := dispatcher.Dispatch(ctx, ref, status); err != nil {
		// Webhook delivery is best-effort; log and continue redirecting.
		audit.FromContext(ctx).Log(ctx, "checkout.webhook_failed", "checkout_session", ref, "", err.Error())
	}
}

func validateCreateRequest(req CreateCheckoutSessionRequest) error {
	if req.SuccessURL == "" {
		return errMissingParam("success_url")
	}
	if len(req.LineItems) == 0 {
		return errMissingParam("line_items")
	}
	for i, item := range req.LineItems {
		if item.PriceData == nil {
			return errMissingParam(fmt.Sprintf("line_items[%d].price_data", i))
		}
		if item.PriceData.Currency == "" {
			return errMissingParam(fmt.Sprintf("line_items[%d].price_data.currency", i))
		}
		if item.PriceData.UnitAmount <= 0 {
			return errInvalidParam(fmt.Sprintf("line_items[%d].price_data.unit_amount", i), "must be greater than 0")
		}
		if item.PriceData.ProductData.Name == "" {
			return errMissingParam(fmt.Sprintf("line_items[%d].price_data.product_data.name", i))
		}
		if item.Quantity <= 0 {
			return errInvalidParam(fmt.Sprintf("line_items[%d].quantity", i), "must be greater than 0")
		}
	}
	if err := validatePaymentMethodTypes(req.PaymentMethodTypes); err != nil {
		return errInvalidParam("payment_method_types", err.Error())
	}
	return nil
}

func buildSession(req CreateCheckoutSessionRequest, baseURL string) *CheckoutSession {
	mode := req.Mode
	if mode == "" {
		mode = "payment"
	}

	var total int64
	for _, item := range req.LineItems {
		total += item.Quantity * item.PriceData.UnitAmount
	}

	currency := req.LineItems[0].PriceData.Currency
	id := "cs_test_" + strings.ReplaceAll(uuid.Must(uuid.NewRandom()).String(), "-", "")

	return &CheckoutSession{
		ID:                 id,
		Object:             "checkout.session",
		AmountTotal:        total,
		Currency:           strings.ToLower(currency),
		CustomerEmail:      req.CustomerEmail,
		Mode:               mode,
		PaymentMethodTypes: normalizePaymentMethodTypes(req.PaymentMethodTypes),
		PaymentStatus:      "unpaid",
		Status:             "open",
		SuccessURL:         req.SuccessURL,
		CancelURL:          req.CancelURL,
		URL:                baseURL + "/v1/checkout/sessions/" + id + "/pay",
		ClientReferenceID:  req.ClientReferenceID,
		Metadata:           req.Metadata,
	}
}

func toLedgerTransaction(ctx context.Context, session *CheckoutSession, req CreateCheckoutSessionRequest) engine.Transaction {
	items := make([]engine.TransactionItem, 0, len(req.LineItems))
	for _, item := range req.LineItems {
		items = append(items, engine.TransactionItem{
			ItemCode: item.PriceData.ProductData.Name,
			Price:    float64(item.PriceData.UnitAmount) / 100.0,
			Quantity: int(item.Quantity),
		})
	}

	return engine.NewTransaction(engine.Transaction{
		Provider:       ProviderName,
		Type:           "checkout_session",
		Amount:         float64(session.AmountTotal) / 100.0,
		Currency:       strings.ToUpper(session.Currency),
		Status:         engine.TransactionStatusNew,
		CustomerRef:    session.CustomerEmail,
		IdempotencyKey: session.ClientReferenceID,
		Reference:      session.ID,
		TraceID:        httputil.TraceIDFromContext(ctx),
		Items:          items,
	})
}

func resolveBaseURL(baseURL string, r *http.Request) string {
	if baseURL != "" {
		return baseURL
	}
	scheme := "http"
	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	} else if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost"
	}
	return scheme + "://" + host
}

func sessionIDFromRetrievePath(path string) string {
	return strings.TrimPrefix(path, "/v1/checkout/sessions/")
}

func sessionIDFromPayPath(path string) string {
	prefix := "/v1/checkout/sessions/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(path, prefix), "/pay")
}

func checkoutPageData(session *CheckoutSession) ui.StripeCheckoutPageData {
	showCard := false
	showFPX := false
	for _, t := range session.PaymentMethodTypes {
		switch t {
		case "card":
			showCard = true
		case "fpx":
			showFPX = true
		}
	}

	lineItems := []ui.StripeCheckoutLineItem{{
		Name:     "Total",
		Quantity: 1,
		Amount:   session.AmountTotal,
		Currency: session.Currency,
	}}

	return ui.StripeCheckoutPageData{
		ID:                      session.ID,
		AmountTotal:             session.AmountTotal,
		Currency:                session.Currency,
		AmountTotalDisplay:      formatAmount(session.AmountTotal, session.Currency),
		LineItems:               lineItems,
		PaymentMethodTypes:      session.PaymentMethodTypes,
		ShowCard:                showCard,
		ShowFPX:                 showFPX,
		ShowPaymentMethodToggle: showCard && showFPX,
		FPXBanks:                uiFPXBanks(),
	}
}

func formatAmount(amount int64, currency string) string {
	major := float64(amount) / 100.0
	return fmt.Sprintf("%.2f %s", major, strings.ToUpper(currency))
}

func uiFPXBanks() []ui.FPXBank {
	banks := make([]ui.FPXBank, len(FPXBanks))
	for i, b := range FPXBanks {
		banks[i] = ui.FPXBank{Code: b.Code, Name: b.Name}
	}
	return banks
}
