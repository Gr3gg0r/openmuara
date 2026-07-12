package ipay88

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
)

type fakeRoundTripper struct {
	requests []*http.Request
	response *http.Response
	err      error
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)
	if f.err != nil {
		return nil, f.err
	}
	return f.response, nil
}

func newFakeClient(response string) (*http.Client, *fakeRoundTripper) {
	rt := &fakeRoundTripper{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
			Header:     make(http.Header),
		},
	}
	return &http.Client{Transport: rt}, rt
}

func seedPaymentRequest(t *testing.T, p *Provider, ref string) {
	t.Helper()
	form := buildEntryForm(ref, "12.50")
	req := httptest.NewRequest(http.MethodPost, "/ePayment/entry.asp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.entryHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("seed entry failed: status %d, body %s", rec.Code, rec.Body.String())
	}
}

func TestAdminPayPageRenders(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-ADMIN-1")

	req := httptest.NewRequest(http.MethodGet, "/_admin/ipay88/pay/REF-ADMIN-1", nil)
	req.SetPathValue("refNo", "REF-ADMIN-1")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "token-123"))
	rec := httptest.NewRecorder()

	p.adminPayPageHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "REF-ADMIN-1") {
		t.Error("page missing reference")
	}
	if !strings.Contains(body, "token-123") {
		t.Error("page missing csrf token")
	}
}

func TestAdminPayPageNotFound(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/_admin/ipay88/pay/missing", nil)
	rec := httptest.NewRecorder()

	p.adminPayPageHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestAdminPayActionPay(t *testing.T) {
	p := newTestProvider(t)
	client, rt := newFakeClient("RECEIVEOK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-ADMIN-2")

	form := url.Values{}
	form.Set("payment_method", "2")
	form.Set("outcome", "pay")
	form.Set("csrf_token", "token")

	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-ADMIN-2", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-ADMIN-2")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.adminPayActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := p.store.GetByReference("REF-ADMIN-2")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusPaid {
		t.Errorf("status: want paid, got %q", tx.Status)
	}

	if len(rt.requests) != 1 {
		t.Fatalf("backend posts: want 1, got %d", len(rt.requests))
	}
	backendBody := rt.requests[0].URL.String()
	if backendBody != "http://example.com/backend" {
		t.Errorf("backend url: want example.com, got %q", backendBody)
	}
}

func TestAdminPayActionCancel(t *testing.T) {
	p := newTestProvider(t)
	client, rt := newFakeClient("RECEIVEOK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-ADMIN-3")

	form := url.Values{}
	form.Set("payment_method", "2")
	form.Set("outcome", "cancel")
	form.Set("csrf_token", "token")

	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-ADMIN-3", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-ADMIN-3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.adminPayActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := p.store.GetByReference("REF-ADMIN-3")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("status: want unpaid, got %q", tx.Status)
	}

	if len(rt.requests) != 1 {
		t.Fatalf("backend posts: want 1, got %d", len(rt.requests))
	}
}

func TestAdminPayActionBackendAckFailure(t *testing.T) {
	p := newTestProvider(t)
	client, _ := newFakeClient("FAIL")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-ADMIN-4")

	form := url.Values{}
	form.Set("payment_method", "2")
	form.Set("outcome", "pay")
	form.Set("csrf_token", "token")

	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-ADMIN-4", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-ADMIN-4")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	p.adminPayActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", rec.Code)
	}
}
