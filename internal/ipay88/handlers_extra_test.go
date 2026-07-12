package ipay88

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEntryHandlerMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/ePayment/entry.asp", nil)
	rec := httptest.NewRecorder()
	p.entryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}

func TestEntryHandlerInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.entryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerInvalidAmount(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-AMT-ERR", "abc")
	form.Set("Signature", SignRequest(testMerchantKey, testMerchantCode, "REF-AMT-ERR", "abc", "MYR"))

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.entryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestEntryHandlerStoreError(t *testing.T) {
	p := newTestProvider(t)
	form := buildEntryForm("REF-STORE-ERR", "12.50")
	p.SetStore(&fakeErrorStore{createErr: errors.New("db down")})

	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.entryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", rec.Code)
	}
}

func TestResponseHandlerMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/ipay88/response", nil)
	rec := httptest.NewRecorder()
	p.responseHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}

func TestResponseHandlerInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.responseHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestResponseHandlerStatusEmptyFallback(t *testing.T) {
	p := newTestProvider(t)
	client, rt := newFakeClient("OK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-RESP-FALL")

	stored, _ := p.getRequest("REF-RESP-FALL")
	stored.Status = "1"
	p.saveRequest(stored)

	form := url.Values{}
	form.Set("RefNo", "REF-RESP-FALL")
	req := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.responseHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if len(rt.requests) != 1 {
		t.Fatalf("want 1 forwarded, got %d", len(rt.requests))
	}
	body := readBody(rt.requests[0])
	if !strings.Contains(body, "Status=1") {
		t.Errorf("want Status=1 in body, got %q", body)
	}
}

func TestResponseHandlerInvalidResponseURL(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-RESP-URL")

	stored, _ := p.getRequest("REF-RESP-URL")
	stored.ResponseURL = "http://127.0.0.1/callback"
	p.saveRequest(stored)

	form := url.Values{}
	form.Set("RefNo", "REF-RESP-URL")
	r := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.responseHandler().ServeHTTP(rec, r)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestResponseHandlerPostError(t *testing.T) {
	p := newTestProvider(t)
	client := &http.Client{Transport: &fakeRoundTripper{err: errors.New("boom")}}
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-RESP-ERR")

	form := url.Values{}
	form.Set("RefNo", "REF-RESP-ERR")
	req := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.responseHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", rec.Code)
	}
}

func TestRequeryHandlerMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/ePayment/enquiry.asp", nil)
	rec := httptest.NewRecorder()
	p.requeryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}

func TestRequeryHandlerInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/ePayment/enquiry.asp", strings.NewReader("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.requeryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}
