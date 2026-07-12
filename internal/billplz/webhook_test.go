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

func TestPayloadBuilderFormUrlEncoded(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: b.ID})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	values, err := url.ParseQuery(string(payload))
	if err != nil {
		t.Fatalf("parse payload: %v", err)
	}
	if values.Get("id") != b.ID {
		t.Errorf("id: want %q, got %q", b.ID, values.Get("id"))
	}
	if values.Get("x_signature") == "" {
		t.Error("x_signature is empty")
	}
	if values.Get("amount") != "1000" {
		t.Errorf("amount: want 1000, got %q", values.Get("amount"))
	}

	// Verify the signature.
	formMap := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			formMap[k] = v[0]
		}
	}
	if !billplz.VerifyCallback(formMap, "muara-billplz-xsig-key") {
		t.Error("callback signature verification failed")
	}
}

func TestWebhookHandlerAcceptsValidSignature(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	builder := p.PayloadBuilder()
	payload, _ := builder(context.Background(), provider.Transaction{Reference: b.ID})

	handler := p.WebhookHandler()
	req := httptest.NewRequest(http.MethodPost, "/billplz/webhook", strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerRejectsInvalidSignature(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.WebhookHandler()

	form := url.Values{"id": {"bill-123"}, "x_signature": {"invalid"}}
	req := httptest.NewRequest(http.MethodPost, "/billplz/webhook", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestPayloadHeaders(t *testing.T) {
	p := newInitializedProvider(t)
	headers, err := p.PayloadHeaders(context.Background(), provider.Transaction{})
	if err != nil {
		t.Fatalf("payload headers: %v", err)
	}
	if headers["Content-Type"] != "application/x-www-form-urlencoded" {
		t.Errorf("content-type: want application/x-www-form-urlencoded, got %q", headers["Content-Type"])
	}
}
