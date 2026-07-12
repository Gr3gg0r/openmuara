package ipay88

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequerySuccess(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-REQ-1")

	form := url.Values{}
	form.Set("payment_method", "2")
	form.Set("outcome", "pay")
	form.Set("csrf_token", "token")
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-REQ-1", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-REQ-1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client, _ := newFakeClient("RECEIVEOK")
	p.SetHTTPClient(client)
	p.adminPayActionHandler().ServeHTTP(httptest.NewRecorder(), req)

	form = url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("RefNo", "REF-REQ-1")
	form.Set("Amount", "12.50")
	req = httptest.NewRequest(http.MethodPost, "/ePayment/enquiry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.requeryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != string(RequeryStatusSuccess) {
		t.Errorf("body: want 00, got %q", rec.Body.String())
	}
}

func TestRequeryUnknown(t *testing.T) {
	p := newTestProvider(t)

	form := url.Values{}
	form.Set("MerchantCode", testMerchantCode)
	form.Set("RefNo", "REF-REQ-MISSING")
	form.Set("Amount", "12.50")
	req := httptest.NewRequest(http.MethodPost, "/ePayment/enquiry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.requeryHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != string(RequeryStatusFailure) {
		t.Errorf("body: want 01, got %q", rec.Body.String())
	}
}
