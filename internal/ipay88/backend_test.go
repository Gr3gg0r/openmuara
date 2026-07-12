package ipay88

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func newTestDispatcher(t *testing.T) *webhook.Dispatcher {
	t.Helper()
	d := webhook.NewDispatcherFromBuilder("http://example.com/webhook", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("payload"), nil }, nil)
	d.ProviderName = ProviderName
	return d
}

func TestBackendHandlerValidSuccess(t *testing.T) {
	p := newTestProvider(t)
	d := newTestDispatcher(t)
	p.SetDispatcher(d)
	seedPaymentRequest(t, p, "REF-BACK-1")

	paymentID := "2"
	status := "1"
	sig := SignResponse(testMerchantKey, testMerchantCode, paymentID, "REF-BACK-1", "12.50", "MYR", status)

	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("PaymentId", paymentID)
	form.Set("RefNo", "REF-BACK-1")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", status)
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", sig)

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if strings.TrimSpace(rec.Body.String()) != "RECEIVEOK" {
		t.Errorf("body: want RECEIVEOK, got %q", rec.Body.String())
	}

	attempt, err := d.Store.Get("REF-BACK-1")
	if err != nil {
		t.Fatalf("store get: %v", err)
	}
	if attempt == nil {
		t.Fatal("expected webhook attempt to be recorded")
	}
}

func TestBackendHandlerInvalidSignature(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-BACK-2")

	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("PaymentId", "2")
	form.Set("RefNo", "REF-BACK-2")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", "1")
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", "invalid")

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestBackendHandlerNotFound(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("RefNo", "missing")
	form.Set("Status", "1")

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestBackendHandlerWrongMethod(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/ipay88/backend", nil)
	rec := httptest.NewRecorder()

	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}
