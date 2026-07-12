package ipay88

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestResponseHandlerForwardsToResponseURL(t *testing.T) {
	p := newTestProvider(t)
	client, rt := newFakeClient("OK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-RESP-1")

	form := url.Values{}
	form.Set("RefNo", "REF-RESP-1")
	form.Set("Status", "1")
	form.Set("PaymentId", "2")

	req := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.responseHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	if len(rt.requests) != 1 {
		t.Fatalf("forwarded requests: want 1, got %d", len(rt.requests))
	}

	forwarded := rt.requests[0]
	if forwarded.URL.String() != "http://example.com/response" {
		t.Errorf("forward url: want example.com/response, got %q", forwarded.URL.String())
	}

	body := readBody(forwarded)
	if !strings.Contains(body, "Signature=") {
		t.Error("forwarded body missing signature")
	}
	if !strings.Contains(body, "Status=1") {
		t.Error("forwarded body missing status")
	}
	if !strings.Contains(body, "PaymentId=2") {
		t.Error("forwarded body missing payment id")
	}
}

func TestResponseHandlerNotFound(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("RefNo", "missing")
	form.Set("Status", "1")

	req := httptest.NewRequest(http.MethodPost, "/ipay88/response", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.responseHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}
