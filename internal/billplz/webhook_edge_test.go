package billplz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestWebhookMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodPost, "/billplz/webhook")

	req := httptest.NewRequest(http.MethodGet, "/billplz/webhook", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestWebhookInvalidForm(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.WebhookHandler()

	req := httptest.NewRequest(http.MethodPost, "/billplz/webhook", strings.NewReader("%ZZ"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestWebhookMissingSignature(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.WebhookHandler()

	form := url.Values{"id": {"bill-123"}}
	req := httptest.NewRequest(http.MethodPost, "/billplz/webhook", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestVerifyCallbackMissingSignature(t *testing.T) {
	values := map[string]string{"id": "bill-123"}
	if billplz.VerifyCallback(values, "secret") {
		t.Fatal("expected VerifyCallback to fail when x_signature is missing")
	}
}

func TestPayloadBuilderBillNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	builder := p.PayloadBuilder()

	_, err := builder(context.Background(), provider.Transaction{Reference: "missing"})
	if err == nil {
		t.Fatal("expected error for missing bill")
	}
}
