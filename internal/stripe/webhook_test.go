package stripe

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestPayloadBuilderReturnsEvent(t *testing.T) {
	// Given a provider with a session
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetBaseURL("http://localhost")

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	// When PayloadBuilder is called
	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{
		Reference: session.ID,
		Status:    "complete",
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	// Then it returns a checkout.session.completed event
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	if event.Object != "event" {
		t.Errorf("object: want event, got %q", event.Object)
	}
	if event.Type != "checkout.session.completed" {
		t.Errorf("type: want checkout.session.completed, got %q", event.Type)
	}
	if event.Data.Object == nil {
		t.Fatal("event data object is nil")
	}
	if event.Data.Object.ID != session.ID {
		t.Errorf("session id mismatch: want %q, got %q", session.ID, event.Data.Object.ID)
	}
	if event.Data.Object.Status != "complete" {
		t.Errorf("status: want complete, got %q", event.Data.Object.Status)
	}
}

func TestPayloadHeadersProducesStripeSignature(t *testing.T) {
	// Given a provider with a session
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

	// When PayloadHeaders is called
	headers, err := p.PayloadHeaders(context.Background(), provider.Transaction{
		Reference: session.ID,
		Status:    "complete",
	})
	if err != nil {
		t.Fatalf("build headers: %v", err)
	}

	// Then it returns a verifiable Stripe-Signature
	sig, ok := headers["Stripe-Signature"]
	if !ok {
		t.Fatal("Stripe-Signature header missing")
	}

	payload, err := p.PayloadBuilder()(context.Background(), provider.Transaction{
		Reference: session.ID,
		Status:    "complete",
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	if err := VerifySignature(payload, sig, "whsec_muara"); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestSuccessSimulationEndpoint(t *testing.T) {
	// Given a provider with a session and a dispatcher
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Store(true)
		if r.Header.Get("Stripe-Signature") == "" {
			t.Error("Stripe-Signature header missing")
		}
		payload, _ := io.ReadAll(r.Body)
		var event Event
		if err := json.Unmarshal(payload, &event); err != nil {
			t.Errorf("decode event: %v", err)
		}
		if event.Type != "checkout.session.completed" {
			t.Errorf("type: want checkout.session.completed, got %q", event.Type)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dispatcher := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	p.SetDispatcher(dispatcher)

	httpReq := httptest.NewRequest(http.MethodPost, "/_admin/stripe/success?session_id="+session.ID, nil)
	successRec := httptest.NewRecorder()
	p.successSimulationHandler().ServeHTTP(successRec, httpReq)

	// Then it returns success
	if successRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", successRec.Code, successRec.Body.String())
	}
	var resp map[string]string
	if err := json.Unmarshal(successRec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["status"] != "complete" {
		t.Errorf("status: want complete, got %q", resp["status"])
	}

	// And the webhook is delivered
	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("webhook was not delivered")
	}
}

func TestSuccessSimulationEndpointSessionNotFound(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	httpReq := httptest.NewRequest(http.MethodPost, "/_admin/stripe/success?session_id=cs_test_missing", nil)
	rec := httptest.NewRecorder()
	p.successSimulationHandler().ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestSuccessSimulationEndpointMethodNotAllowed(t *testing.T) {
	p := NewProvider()
	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/stripe/success?session_id=cs_1", nil)
	rec := httptest.NewRecorder()
	p.successSimulationHandler().ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestSuccessSimulationEndpointMissingSessionID(t *testing.T) {
	p := NewProvider()
	httpReq := httptest.NewRequest(http.MethodPost, "/_admin/stripe/success", nil)
	rec := httptest.NewRecorder()
	p.successSimulationHandler().ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestSuccessSimulationEndpointNoDispatcher(t *testing.T) {
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

	httpReq := httptest.NewRequest(http.MethodPost, "/_admin/stripe/success?session_id="+session.ID, nil)
	successRec := httptest.NewRecorder()
	p.successSimulationHandler().ServeHTTP(successRec, httpReq)

	if successRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", successRec.Code)
	}
}
