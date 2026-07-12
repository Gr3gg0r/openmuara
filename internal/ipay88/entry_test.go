package ipay88

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

const (
	testMerchantCode = "M00001"
	testMerchantKey  = "secret-key"
)

func newTestProvider(t *testing.T) *Provider {
	t.Helper()
	p := NewProvider()
	if err := p.Init(map[string]any{
		"merchant_code": testMerchantCode,
		"merchant_key":  testMerchantKey,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetStore(engine.NewMemoryStore())
	return p
}

func buildEntryForm(ref, amount string) url.Values {
	v := url.Values{}
	v.Set("MerchantCode", testMerchantCode)
	v.Set("PaymentId", "2")
	v.Set("RefNo", ref)
	v.Set("Amount", amount)
	v.Set("Currency", "MYR")
	v.Set("ProdDesc", "Test Product")
	v.Set("UserName", "Test User")
	v.Set("UserEmail", "test@example.com")
	v.Set("UserContact", "0123456789")
	v.Set("Remark", "")
	v.Set("Lang", "UTF8")
	v.Set("SignatureType", "SHA256")
	v.Set("ResponseURL", "http://example.com/response")
	v.Set("BackendURL", "http://example.com/backend")
	v.Set("Signature", SignRequest(testMerchantKey, testMerchantCode, ref, amount, "MYR"))
	return v
}

func TestEntryHandlerSuccess(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-1", "12.50")

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "/_admin/ipay88/pay/REF-ENTRY-1") {
		t.Errorf("unexpected redirect: %q", loc)
	}

	tx, ok, err := p.store.GetByReference("REF-ENTRY-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction")
	}
	if tx.Status != engine.TransactionStatusNew {
		t.Errorf("status: want new, got %q", tx.Status)
	}
}

func TestEntryHandlerInvalidSignature(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-2", "12.50")
	form.Set("Signature", "invalid")

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerWrongSignatureType(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-3", "12.50")
	form.Set("SignatureType", "MD5")

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerMissingField(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-4", "12.50")
	form.Del("RefNo")

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerRejectsPrivateResponseURL(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-5", "12.50")
	form.Set("ResponseURL", "http://127.0.0.1:9999/callback")
	form.Set("Signature", SignRequest(testMerchantKey, testMerchantCode, "REF-ENTRY-5", "12.50", "MYR"))

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerRejectsPrivateBackendURL(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-ENTRY-6", "12.50")
	form.Set("BackendURL", "http://192.168.1.1/backend")
	form.Set("Signature", SignRequest(testMerchantKey, testMerchantCode, "REF-ENTRY-6", "12.50", "MYR"))

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.entryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}
