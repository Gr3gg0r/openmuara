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

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/server"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func validPaymentIntentRequest() PaymentIntentRequest {
	return PaymentIntentRequest{
		Amount:             5000,
		Currency:           "myr",
		PaymentMethodTypes: []string{"fpx"},
		Metadata: map[string]string{
			"order_id": "order-456",
		},
	}
}

func TestCreatePaymentIntentReturnsPaymentIntent(t *testing.T) {
	store := NewMemoryPaymentIntentStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreatePaymentIntentHandler(store, ledger, nil)

	body, _ := json.Marshal(validPaymentIntentRequest())
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var pi PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &pi); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if !strings.HasPrefix(pi.ID, "pi_test_") {
		t.Errorf("id: want pi_test_ prefix, got %q", pi.ID)
	}
	if pi.Status != "requires_confirmation" {
		t.Errorf("status: want requires_confirmation, got %q", pi.Status)
	}
	if pi.Amount != 5000 {
		t.Errorf("amount: want 5000, got %d", pi.Amount)
	}
	if pi.Currency != "myr" {
		t.Errorf("currency: want myr, got %q", pi.Currency)
	}
	if pi.ClientSecret == "" {
		t.Error("client_secret is empty")
	}
}

func TestCreatePaymentIntentRecordsTransaction(t *testing.T) {
	store := NewMemoryPaymentIntentStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreatePaymentIntentHandler(store, ledger, nil)

	body, _ := json.Marshal(validPaymentIntentRequest())
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var pi PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &pi); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	tx, ok, err := ledger.GetByReference(pi.ID)
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != "stripe" {
		t.Errorf("provider: want stripe, got %q", tx.Provider)
	}
	if tx.Amount != 50.00 {
		t.Errorf("amount: want 50.00, got %f", tx.Amount)
	}
	if tx.Currency != "MYR" {
		t.Errorf("currency: want MYR, got %q", tx.Currency)
	}
	if tx.Status != engine.TransactionStatusNew {
		t.Errorf("status: want new, got %q", tx.Status)
	}
}

func TestCreatePaymentIntentValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  PaymentIntentRequest
	}{
		{"missing amount", PaymentIntentRequest{Currency: "usd"}},
		{"missing currency", PaymentIntentRequest{Amount: 1000}},
		{"invalid amount", PaymentIntentRequest{Amount: 0, Currency: "usd"}},
		{"unsupported type", PaymentIntentRequest{Amount: 1000, Currency: "usd", PaymentMethodTypes: []string{"ideal"}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewCreatePaymentIntentHandler(NewMemoryPaymentIntentStore(), engine.NewMemoryStore(), nil)
			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

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
		})
	}
}

func TestGetPaymentIntentReturnsPaymentIntent(t *testing.T) {
	store := NewMemoryPaymentIntentStore()
	ledger := engine.NewMemoryStore()
	create := NewCreatePaymentIntentHandler(store, ledger, nil)

	body, _ := json.Marshal(validPaymentIntentRequest())
	createRec := httptest.NewRecorder()
	create.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/v1/payment_intents", bytes.NewReader(body)))

	var created PaymentIntent
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created pi: %v", err)
	}

	get := NewGetPaymentIntentHandler(store)
	rec := httptest.NewRecorder()
	get.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/payment_intents/"+created.ID, nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var fetched PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode fetched pi: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: created %q, fetched %q", created.ID, fetched.ID)
	}
}

func TestGetPaymentIntentNotFound(t *testing.T) {
	handler := NewGetPaymentIntentHandler(NewMemoryPaymentIntentStore())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/payment_intents/pi_test_missing", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error.Code != "resource_missing" {
		t.Errorf("error code: want resource_missing, got %q", resp.Error.Code)
	}
}

func TestConfirmPaymentIntentCardSucceeds(t *testing.T) {
	p := providerWithPaymentIntent(t)

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dispatcher := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	p.SetDispatcher(dispatcher)

	pi := createPaymentIntent(t, p)
	confirm := confirmPaymentIntent(t, p, pi.ID, "pm_card_visa")

	if confirm.Status != "succeeded" {
		t.Errorf("status: want succeeded, got %q", confirm.Status)
	}
	if confirm.PaymentMethod != "pm_card_visa" {
		t.Errorf("payment_method: want pm_card_visa, got %q", confirm.PaymentMethod)
	}

	tx, ok, _ := p.ledger.GetByReference(pi.ID)
	if !ok {
		t.Fatal("transaction not found in ledger")
	}
	if tx.Status != engine.TransactionStatusPaid {
		t.Errorf("ledger status: want paid, got %q", tx.Status)
	}

	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("webhook was not delivered")
	}
}

func TestConfirmPaymentIntentFPXRequiresAction(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	confirm := confirmPaymentIntent(t, p, pi.ID, "pm_fpx_maybank")

	if confirm.Status != "requires_action" {
		t.Fatalf("status: want requires_action, got %q", confirm.Status)
	}
	if confirm.NextAction == nil || confirm.NextAction.Type != "redirect_to_url" {
		t.Fatal("expected redirect_to_url next_action")
	}
	if !strings.Contains(confirm.NextAction.RedirectToURL.URL, "/_admin/stripe/payment_intent/"+pi.ID) {
		t.Errorf("redirect url missing admin path: %q", confirm.NextAction.RedirectToURL.URL)
	}
}

func TestConfirmPaymentIntentUnknownToken(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	rec := httptest.NewRecorder()
	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: "pm_unknown"})
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/"+pi.ID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	NewConfirmPaymentIntentHandler(p.paymentIntents, p.ledger, nil, "").ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error.Type != "invalid_request_error" {
		t.Errorf("error type: want invalid_request_error, got %q", resp.Error.Type)
	}
	if resp.Error.Code != "resource_missing" {
		t.Errorf("error code: want resource_missing, got %q", resp.Error.Code)
	}
	if resp.Error.Param != "payment_method" {
		t.Errorf("error param: want payment_method, got %q", resp.Error.Param)
	}
}

func TestCancelPaymentIntent(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dispatcher := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	p.SetDispatcher(dispatcher)

	cancel := cancelPaymentIntent(t, p, pi.ID)
	if cancel.Status != "canceled" {
		t.Errorf("status: want canceled, got %q", cancel.Status)
	}

	tx, ok, _ := p.ledger.GetByReference(pi.ID)
	if !ok {
		t.Fatal("transaction not found in ledger")
	}
	if tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("ledger status: want unpaid, got %q", tx.Status)
	}

	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("webhook was not delivered")
	}
}

func TestCancelPaymentIntentInvalidTransition(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	// Confirm the card payment so it is no longer cancellable.
	confirmPaymentIntent(t, p, pi.ID, "pm_card_visa")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/"+pi.ID+"/cancel", nil)
	NewCancelPaymentIntentHandler(p.paymentIntents, p.ledger, nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Error.Code != "payment_intent_unexpected_state" {
		t.Errorf("error code: want payment_intent_unexpected_state, got %q", resp.Error.Code)
	}
}

func TestPaymentIntentAdminPageRendersFPX(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	page := NewPaymentIntentAdminPageHandler(p.paymentIntents)
	req := httptest.NewRequest(http.MethodGet, "/_admin/stripe/payment_intent/"+pi.ID, nil)
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "test-csrf"))
	rec := httptest.NewRecorder()
	page.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Maybank2U") {
		t.Error("admin page missing FPX bank selector")
	}
	if !strings.Contains(body, `name="csrf_token" value="test-csrf"`) {
		t.Error("admin page missing CSRF token")
	}
}

func TestPaymentIntentAdminPageRendersCard(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntentWithTypes(t, p, []string{"card"})

	page := NewPaymentIntentAdminPageHandler(p.paymentIntents)
	req := httptest.NewRequest(http.MethodGet, "/_admin/stripe/payment_intent/"+pi.ID, nil)
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "test-csrf"))
	rec := httptest.NewRecorder()
	page.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Confirm the card payment") {
		t.Error("admin page missing card confirmation prompt")
	}
	if strings.Contains(body, "Bank") {
		t.Error("card admin page should not show bank selector")
	}
}

func TestPaymentIntentAdminActionConfirm(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dispatcher := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	p.SetDispatcher(dispatcher)

	// Move to requires_action first to simulate FPX redirect.
	confirmPaymentIntent(t, p, pi.ID, "pm_fpx_maybank")

	action := NewPaymentIntentAdminActionHandler(p.paymentIntents, p.ledger, p.dispatcher)
	form := "action=confirm"
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/payment_intent/"+pi.ID, strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	req.Header.Set("X-CSRF-Token", "test-csrf")
	rec := httptest.NewRecorder()
	action.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var updated PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if updated.Status != "succeeded" {
		t.Errorf("status: want succeeded, got %q", updated.Status)
	}
	if updated.NextAction != nil {
		t.Error("next_action should be cleared")
	}

	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("webhook was not delivered")
	}
}

func TestPaymentIntentAdminActionCancel(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	action := NewPaymentIntentAdminActionHandler(p.paymentIntents, p.ledger, nil)
	form := "action=cancel"
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/payment_intent/"+pi.ID, strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:     "openmuara_csrf",
		Value:    "test-csrf",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	req.Header.Set("X-CSRF-Token", "test-csrf")
	rec := httptest.NewRecorder()
	action.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var updated PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if updated.Status != "canceled" {
		t.Errorf("status: want canceled, got %q", updated.Status)
	}
}

func TestPaymentIntentAdminActionCSRFMissing(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	action := NewPaymentIntentAdminActionHandler(p.paymentIntents, p.ledger, nil)
	form := "action=confirm"
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/payment_intent/"+pi.ID, strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	wrapped := server.CSRFGuardMiddleware(server.CSRFGuardConfig{Enabled: true})(action)
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestPaymentIntentPayloadEventType(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	cases := []struct {
		ref, status, want string
	}{
		{"pi_test_123", "NEW", "payment_intent.created"},
		{"pi_test_123", "PAID", "payment_intent.succeeded"},
		{"pi_test_123", "UNPAID", "payment_intent.canceled"},
		{"cs_test_123", "PAID", "checkout.session.completed"},
		{"cs_test_123", "UNPAID", "checkout.session.expired"},
	}

	for _, tc := range cases {
		got := p.PayloadEventType(tc.ref, tc.status)
		if got != tc.want {
			t.Errorf("PayloadEventType(%q, %q): want %q, got %q", tc.ref, tc.status, tc.want, got)
		}
	}
}

func TestPaymentIntentPayloadBuilder(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	pi := &PaymentIntent{
		ID:                 "pi_test_payload",
		Object:             "payment_intent",
		Amount:             1000,
		Currency:           "usd",
		Status:             "succeeded",
		ClientSecret:       "pi_test_payload_secret_test",
		PaymentMethodTypes: []string{"card"},
	}
	p.paymentIntents.Save(pi.ID, pi)

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: pi.ID, Status: "PAID"})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var event PaymentIntentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if event.Type != "payment_intent.succeeded" {
		t.Errorf("event type: want payment_intent.succeeded, got %q", event.Type)
	}
	if event.Data.Object == nil || event.Data.Object.ID != pi.ID {
		t.Error("event data object missing or mismatch")
	}
}

func providerWithPaymentIntent(t *testing.T) *Provider {
	t.Helper()
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	return p
}

func createPaymentIntent(t *testing.T, p *Provider) *PaymentIntent {
	t.Helper()
	return createPaymentIntentWithTypes(t, p, []string{"fpx"})
}

func createPaymentIntentWithTypes(t *testing.T, p *Provider, types []string) *PaymentIntent {
	t.Helper()
	req := validPaymentIntentRequest()
	req.PaymentMethodTypes = types
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	NewCreatePaymentIntentHandler(p.paymentIntents, p.ledger, p.dispatcher).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/payment_intents", bytes.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create pi: want 201, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var pi PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &pi); err != nil {
		t.Fatalf("decode pi: %v", err)
	}
	return &pi
}

func confirmPaymentIntent(t *testing.T, p *Provider, id, paymentMethod string) *PaymentIntent {
	t.Helper()
	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: paymentMethod})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/"+id+"/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	NewConfirmPaymentIntentHandler(p.paymentIntents, p.ledger, p.dispatcher, p.baseURL).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("confirm pi: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var pi PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &pi); err != nil {
		t.Fatalf("decode confirmed pi: %v", err)
	}
	return &pi
}

func cancelPaymentIntent(t *testing.T, p *Provider, id string) *PaymentIntent {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/"+id+"/cancel", nil)
	NewCancelPaymentIntentHandler(p.paymentIntents, p.ledger, p.dispatcher).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("cancel pi: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var pi PaymentIntent
	if err := json.Unmarshal(rec.Body.Bytes(), &pi); err != nil {
		t.Fatalf("decode canceled pi: %v", err)
	}
	return &pi
}
