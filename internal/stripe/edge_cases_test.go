package stripe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
)

func TestConfirmPaymentIntentDuplicateCardConfirm(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntentWithTypes(t, p, []string{"card"})
	confirmPaymentIntent(t, p, pi.ID, "pm_card_visa")

	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: "pm_card_visa"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/"+pi.ID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	NewConfirmPaymentIntentHandler(p.paymentIntents, p.ledger, nil, "").ServeHTTP(rec, req)

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

func TestCancelPaymentIntentAfterCancel(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)
	cancelPaymentIntent(t, p, pi.ID)

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

func TestConfirmPaymentIntentMissingPaymentMethod(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntent(t, p)

	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: ""})
	rec := httptest.NewRecorder()
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
	if resp.Error.Param != "payment_method" {
		t.Errorf("error param: want payment_method, got %q", resp.Error.Param)
	}
}

func TestConfirmPaymentIntentInvalidCardToken(t *testing.T) {
	p := providerWithPaymentIntent(t)
	pi := createPaymentIntentWithTypes(t, p, []string{"card"})

	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: "pm_card_invalid"})
	rec := httptest.NewRecorder()
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
	if resp.Error.Code != "resource_missing" {
		t.Errorf("error code: want resource_missing, got %q", resp.Error.Code)
	}
}

func TestConfirmPaymentIntentMissingSession(t *testing.T) {
	p := providerWithPaymentIntent(t)

	body, _ := json.Marshal(PaymentIntentConfirmRequest{PaymentMethod: "pm_card_visa"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/pi_test_missing/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	NewConfirmPaymentIntentHandler(p.paymentIntents, p.ledger, nil, "").ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCancelPaymentIntentMissingSession(t *testing.T) {
	p := providerWithPaymentIntent(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/payment_intents/pi_test_missing/cancel", nil)
	NewCancelPaymentIntentHandler(p.paymentIntents, p.ledger, nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutSessionDuplicateConfirm(t *testing.T) {
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

	action := NewCheckoutSessionPayActionHandler(sessions, ledger, nil)
	makeAction := func() *httptest.ResponseRecorder {
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
		rec := httptest.NewRecorder()
		action.ServeHTTP(rec, actionReq)
		return rec
	}

	if rec := makeAction(); rec.Code != http.StatusSeeOther {
		t.Fatalf("first confirm status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if rec := makeAction(); rec.Code != http.StatusConflict {
		t.Fatalf("duplicate confirm status: want 409, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutSessionDuplicateCancel(t *testing.T) {
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

	action := NewCheckoutSessionPayActionHandler(sessions, ledger, nil)
	makeAction := func() *httptest.ResponseRecorder {
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
		rec := httptest.NewRecorder()
		action.ServeHTTP(rec, actionReq)
		return rec
	}

	if rec := makeAction(); rec.Code != http.StatusSeeOther {
		t.Fatalf("first cancel status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if rec := makeAction(); rec.Code != http.StatusConflict {
		t.Fatalf("duplicate cancel status: want 409, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutPayPageMissingSession(t *testing.T) {
	sessions := NewMemorySessionStore()
	pageHandler := NewCheckoutSessionPayPageHandler(sessions)

	req := httptest.NewRequest(http.MethodGet, "/v1/checkout/sessions/cs_test_missing/pay", nil)
	rec := httptest.NewRecorder()
	pageHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutPayActionMissingSession(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	action := NewCheckoutSessionPayActionHandler(sessions, ledger, nil)

	form := "action=confirm&bank=maybank2u"
	req := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions/cs_test_missing/pay", strings.NewReader(form))
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

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCheckoutPayActionInvalidAction(t *testing.T) {
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

	form := "action=refund"
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
	rec := httptest.NewRecorder()
	NewCheckoutSessionPayActionHandler(sessions, ledger, nil).ServeHTTP(rec, actionReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestPaymentIntentAdminActionMissingSession(t *testing.T) {
	p := providerWithPaymentIntent(t)
	action := NewPaymentIntentAdminActionHandler(p.paymentIntents, p.ledger, nil)

	form := "action=confirm"
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/payment_intent/pi_test_missing", strings.NewReader(form))
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

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
