package ipay88

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestBackendHandlerInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestBackendHandlerInvalidSignatureType(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-BACK-SIGTYPE")

	form := url.Values{}
	form.Set("RefNo", "REF-BACK-SIGTYPE")
	form.Set("Status", "1")
	form.Set("SignatureType", "MD5")

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestBackendHandlerPaymentIDFallback(t *testing.T) {
	p := newTestProvider(t)
	d := newTestDispatcher(t)
	p.SetDispatcher(d)
	seedPaymentRequest(t, p, "REF-BACK-FALL")

	sig := SignResponse(testMerchantKey, testMerchantCode, "2", "REF-BACK-FALL", "12.50", "MYR", "1")
	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("RefNo", "REF-BACK-FALL")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", "1")
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", sig)

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestBackendHandlerTransitionConflict(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-BACK-CONFLICT")
	if err := p.transitionTransaction("REF-BACK-CONFLICT", engine.TransactionStatusPaid); err != nil {
		t.Fatalf("seed paid: %v", err)
	}

	sig := SignResponse(testMerchantKey, testMerchantCode, "2", "REF-BACK-CONFLICT", "12.50", "MYR", "0")
	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("PaymentId", "2")
	form.Set("RefNo", "REF-BACK-CONFLICT")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", "0")
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", sig)

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("want 409, got %d", rec.Code)
	}
}

func TestBackendHandlerNilDispatcher(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-BACK-NODISP")

	sig := SignResponse(testMerchantKey, testMerchantCode, "2", "REF-BACK-NODISP", "12.50", "MYR", "1")
	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("PaymentId", "2")
	form.Set("RefNo", "REF-BACK-NODISP")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", "1")
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", sig)

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
}

func TestBackendHandlerDispatchError(t *testing.T) {
	p := newTestProvider(t)
	d := webhook.NewDispatcherFromBuilder("", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return nil, errors.New("boom") }, nil)
	d.ProviderName = ProviderName
	p.SetDispatcher(d)
	seedPaymentRequest(t, p, "REF-BACK-DISPERR")

	sig := SignResponse(testMerchantKey, testMerchantCode, "2", "REF-BACK-DISPERR", "12.50", "MYR", "1")
	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("PaymentId", "2")
	form.Set("RefNo", "REF-BACK-DISPERR")
	form.Set("Amount", "12.50")
	form.Set("Currency", "MYR")
	form.Set("Status", "1")
	form.Set("SignatureType", "SHA256")
	form.Set("Signature", sig)

	req := httptest.NewRequest(http.MethodPost, "/ipay88/backend", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.backendHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", rec.Code)
	}
}
