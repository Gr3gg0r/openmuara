package stripe

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/server"
	"github.com/openmuara/openmuara/internal/webhook"
)

func validCreateRequest() CreateCheckoutSessionRequest {
	return CreateCheckoutSessionRequest{
		SuccessURL: "http://localhost/success",
		CancelURL:  "http://localhost/cancel",
		Mode:       "payment",
		LineItems: []LineItem{
			{
				PriceData: &PriceData{
					Currency:   "usd",
					UnitAmount: 1000,
					ProductData: struct {
						Name string `json:"name"`
					}{Name: "Test Product"},
				},
				Quantity: 2,
			},
		},
		CustomerEmail:     "customer@example.com",
		ClientReferenceID: "order-123",
	}
}

func TestCreateCheckoutSessionReturnsSession(t *testing.T) {
	// Given a valid create request
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When the handler is called
	handler.ServeHTTP(rec, httpReq)

	// Then it returns 201 with a session
	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if !strings.HasPrefix(session.ID, "cs_test_") {
		t.Errorf("id: want cs_test_ prefix, got %q", session.ID)
	}
	if session.Status != "open" {
		t.Errorf("status: want open, got %q", session.Status)
	}
	if session.PaymentStatus != "unpaid" {
		t.Errorf("payment_status: want unpaid, got %q", session.PaymentStatus)
	}
	if session.AmountTotal != 2000 {
		t.Errorf("amount_total: want 2000, got %d", session.AmountTotal)
	}
	if session.Currency != "usd" {
		t.Errorf("currency: want usd, got %q", session.Currency)
	}
	if session.URL == "" {
		t.Error("url is empty")
	}
}

func TestCreateCheckoutSessionRecordsTransaction(t *testing.T) {
	// Given a valid create request
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When the handler is called
	handler.ServeHTTP(rec, httpReq)

	// Then a transaction is recorded in the ledger
	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	tx, ok, err := ledger.GetByReference(session.ID)

	if err != nil {

		t.Fatalf("lookup transaction: %v", err)

	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != "stripe" {
		t.Errorf("provider: want stripe, got %q", tx.Provider)
	}
	if tx.Amount != 20.00 {
		t.Errorf("amount: want 20.00, got %f", tx.Amount)
	}
	if tx.Currency != "USD" {
		t.Errorf("currency: want USD, got %q", tx.Currency)
	}
	if tx.Status != engine.TransactionStatusNew {
		t.Errorf("status: want %q, got %q", engine.TransactionStatusNew, tx.Status)
	}
	if tx.Reference != session.ID {
		t.Errorf("reference mismatch: session %q, tx %q", session.ID, tx.Reference)
	}
	if len(tx.Items) != 1 {
		t.Fatalf("items: want 1, got %d", len(tx.Items))
	}
	if tx.Items[0].ItemCode != "Test Product" {
		t.Errorf("item code: want Test Product, got %q", tx.Items[0].ItemCode)
	}
}

func TestCreateCheckoutSessionMissingSuccessURL(t *testing.T) {
	// Given a request without success_url
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	req.SuccessURL = ""
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When the handler is called
	handler.ServeHTTP(rec, httpReq)

	// Then it returns 400
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateCheckoutSessionEmptyLineItems(t *testing.T) {
	// Given a request with no line items
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	req.LineItems = nil
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When the handler is called
	handler.ServeHTTP(rec, httpReq)

	// Then it returns 400
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateCheckoutSessionInvalidJSON(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader([]byte("not json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetCheckoutSessionReturnsSession(t *testing.T) {
	// Given an existing session
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	createHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	rec := httptest.NewRecorder()
	createHandler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var created CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created session: %v", err)
	}

	getHandler := NewGetCheckoutSessionHandler(sessions)
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/"+created.ID, nil)
	getRec := httptest.NewRecorder()

	// When GET is called
	getHandler.ServeHTTP(getRec, httpReq)

	// Then it returns the same session
	if getRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", getRec.Code, getRec.Body.String())
	}
	var fetched CheckoutSession
	if err := json.Unmarshal(getRec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode fetched session: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: created %q, fetched %q", created.ID, fetched.ID)
	}
}

func TestGetCheckoutSessionNotFound(t *testing.T) {
	sessions := NewMemorySessionStore()
	handler := NewGetCheckoutSessionHandler(sessions)

	httpReq := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/cs_test_missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateCheckoutSessionMethodNotAllowed(t *testing.T) {
	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	req := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestGetCheckoutSessionMethodNotAllowed(t *testing.T) {
	handler := NewGetCheckoutSessionHandler(NewMemorySessionStore())
	req := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/cs_test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestGetCheckoutSessionMissingID(t *testing.T) {
	handler := NewGetCheckoutSessionHandler(NewMemorySessionStore())
	req := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestValidateCreateRequestMissingPriceData(t *testing.T) {
	req := validCreateRequest()
	req.LineItems[0].PriceData = nil
	body, _ := json.Marshal(req)

	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestValidateCreateRequestMissingCurrency(t *testing.T) {
	req := validCreateRequest()
	req.LineItems[0].PriceData.Currency = ""
	body, _ := json.Marshal(req)

	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestValidateCreateRequestInvalidUnitAmount(t *testing.T) {
	req := validCreateRequest()
	req.LineItems[0].PriceData.UnitAmount = 0
	body, _ := json.Marshal(req)

	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestValidateCreateRequestMissingProductName(t *testing.T) {
	req := validCreateRequest()
	req.LineItems[0].PriceData.ProductData.Name = ""
	body, _ := json.Marshal(req)

	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestValidateCreateRequestInvalidQuantity(t *testing.T) {
	req := validCreateRequest()
	req.LineItems[0].Quantity = 0
	body, _ := json.Marshal(req)

	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestBuildSessionDefaultMode(t *testing.T) {
	req := validCreateRequest()
	req.Mode = ""
	session := buildSession(req, "http://localhost")

	if session.Mode != "payment" {
		t.Errorf("mode: want payment, got %q", session.Mode)
	}
}

func TestCreateCheckoutSessionPaymentMethodTypes(t *testing.T) {
	tests := []struct {
		name  string
		types []string
		want  string
	}{
		{"fpx", []string{"fpx"}, "fpx"},
		{"card", []string{"card"}, "card"},
		{"card and fpx", []string{"card", "fpx"}, "card"},
		{"fpx and card reversed", []string{"fpx", "card"}, "fpx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sessions := NewMemorySessionStore()
			ledger := engine.NewMemoryStore()
			handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

			req := validCreateRequest()
			req.PaymentMethodTypes = tc.types
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, httpReq)

			if rec.Code != http.StatusCreated {
				t.Fatalf("status: want 201, got %d, body: %s", rec.Code, rec.Body.String())
			}

			var session CheckoutSession
			if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if len(session.PaymentMethodTypes) == 0 || session.PaymentMethodTypes[0] != tc.want {
				t.Errorf("payment_method_types: want %v, got %v", tc.types, session.PaymentMethodTypes)
			}
		})
	}
}

func TestCreateCheckoutSessionDefaultsPaymentMethodTypes(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(session.PaymentMethodTypes) != 1 || session.PaymentMethodTypes[0] != "card" {
		t.Errorf("payment_method_types: want [card], got %v", session.PaymentMethodTypes)
	}
}

func TestCreateCheckoutSessionRejectsInvalidPaymentMethodTypes(t *testing.T) {
	tests := []struct {
		name  string
		types []string
	}{
		{"unsupported type", []string{"ideal"}},
		{"too many types", []string{"card", "fpx", "ideal"}},
		{"empty string", []string{""}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")

			req := validCreateRequest()
			req.PaymentMethodTypes = tc.types
			body, _ := json.Marshal(req)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
			}

			var resp ErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if resp.Error.Type != "invalid_request_error" {
				t.Errorf("error type: want invalid_request_error, got %q", resp.Error.Type)
			}
			if resp.Error.Param != "payment_method_types" {
				t.Errorf("error param: want payment_method_types, got %q", resp.Error.Param)
			}
		})
	}
}

func TestCheckoutSessionValidationReturnsStripeError(t *testing.T) {
	handler := NewCreateCheckoutSessionHandler(NewMemorySessionStore(), engine.NewMemoryStore(), "http://localhost")

	req := validCreateRequest()
	req.SuccessURL = ""
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error.Type != "invalid_request_error" {
		t.Errorf("error type: want invalid_request_error, got %q", resp.Error.Type)
	}
	if resp.Error.Param != "success_url" {
		t.Errorf("error param: want success_url, got %q", resp.Error.Param)
	}
}

func TestCheckoutPayPageRendersForFPX(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	req.PaymentMethodTypes = []string{"fpx"}
	body, _ := json.Marshal(req)
	createRec := httptest.NewRecorder()
	createHandler.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(createRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	pageHandler := NewCheckoutSessionPayPageHandler(sessions)
	pageReq := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/"+session.ID+"/pay", nil)
	pageReq = pageReq.WithContext(httputil.WithCSRFToken(pageReq.Context(), "test-csrf"))
	pageRec := httptest.NewRecorder()
	pageHandler.ServeHTTP(pageRec, pageReq)

	if pageRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", pageRec.Code, pageRec.Body.String())
	}
	bodyStr := pageRec.Body.String()
	if !strings.Contains(bodyStr, "Maybank2U") {
		t.Error("fpx checkout page missing bank selector")
	}
	if !strings.Contains(bodyStr, `name="csrf_token" value="test-csrf"`) {
		t.Error("checkout page missing CSRF token")
	}
}

func TestCheckoutPayPageRendersForCard(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	req.PaymentMethodTypes = []string{"card"}
	body, _ := json.Marshal(req)
	createRec := httptest.NewRecorder()
	createHandler.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(createRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	pageHandler := NewCheckoutSessionPayPageHandler(sessions)
	pageReq := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/"+session.ID+"/pay", nil)
	pageReq = pageReq.WithContext(httputil.WithCSRFToken(pageReq.Context(), "test-csrf"))
	pageRec := httptest.NewRecorder()
	pageHandler.ServeHTTP(pageRec, pageReq)

	if pageRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", pageRec.Code, pageRec.Body.String())
	}
	bodyStr := pageRec.Body.String()
	if !strings.Contains(bodyStr, "Card number") {
		t.Error("card checkout page missing card number field")
	}
}

func TestCheckoutPayActionConfirmUpdatesSessionAndDispatchesWebhook(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dispatcher := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	p.SetDispatcher(dispatcher)

	form := "action=confirm&bank=maybank2u"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionReq.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	actionReq.Header.Set("X-CSRF-Token", "test-csrf")
	actionRec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(p.sessions, p.ledger, p.dispatcher).ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}
	if loc := actionRec.Header().Get("Location"); loc != session.SuccessURL {
		t.Errorf("redirect: want %q, got %q", session.SuccessURL, loc)
	}

	updated, ok, _ := p.ledger.GetByReference(session.ID)
	if !ok {
		t.Fatal("transaction not found in ledger")
	}
	if updated.Status != engine.TransactionStatusPaid {
		t.Errorf("ledger status: want paid, got %q", updated.Status)
	}

	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("webhook was not delivered")
	}
}

func TestCheckoutPayActionCancelUpdatesSession(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	createRec := httptest.NewRecorder()
	createHandler.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(createRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	form := "action=cancel"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionReq.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	actionReq.Header.Set("X-CSRF-Token", "test-csrf")
	actionRec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(sessions, ledger, nil).ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}
	if loc := actionRec.Header().Get("Location"); loc != session.CancelURL {
		t.Errorf("redirect: want %q, got %q", session.CancelURL, loc)
	}

	updated, ok := sessions.Load(session.ID)
	if !ok {
		t.Fatal("session not found")
	}
	if updated.Status != "expired" {
		t.Errorf("session status: want expired, got %q", updated.Status)
	}
	if updated.PaymentStatus != "unpaid" {
		t.Errorf("payment status: want unpaid, got %q", updated.PaymentStatus)
	}
}

func TestCheckoutPayActionInvalidTransition(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	createRec := httptest.NewRecorder()
	createHandler.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(createRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	session.Status = "complete"
	sessions.Save(session.ID, &session)

	form := "action=confirm"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionReq.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	actionReq.Header.Set("X-CSRF-Token", "test-csrf")
	actionRec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(sessions, ledger, nil).ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}
}

func TestCheckoutPayActionCSRFMissing(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	createHandler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	createRec := httptest.NewRecorder()
	createHandler.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(createRec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	form := "action=confirm"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionRec := httptest.NewRecorder()

	wrapped := server.CSRFGuardMiddleware(server.CSRFGuardConfig{Enabled: true})(NewCheckoutSessionPayActionHandler(sessions, ledger, nil))
	wrapped.ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}
}

func TestCheckoutPayActionConfirmEventType(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	form := "action=confirm&bank=maybank2u"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionReq.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	actionReq.Header.Set("X-CSRF-Token", "test-csrf")
	actionRec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(p.sessions, p.ledger, nil).ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: session.ID, Status: "complete"})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if event.Type != "checkout.session.completed" {
		t.Errorf("event type: want checkout.session.completed, got %q", event.Type)
	}
	if event.Data.Object == nil || event.Data.Object.PaymentMethodTypes[0] != "card" {
		t.Errorf("payment_method_types missing in event payload")
	}
}

func TestCheckoutPayActionCancelEventType(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	form := "action=cancel"
	actionReq := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/"+session.ID+"/pay", strings.NewReader(form))
	actionReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	actionReq.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	actionReq.Header.Set("X-CSRF-Token", "test-csrf")
	actionRec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(p.sessions, p.ledger, nil).ServeHTTP(actionRec, actionReq)

	if actionRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", actionRec.Code, actionRec.Body.String())
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: session.ID, Status: "UNPAID"})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if event.Type != "checkout.session.expired" {
		t.Errorf("event type: want checkout.session.expired, got %q", event.Type)
	}
}
