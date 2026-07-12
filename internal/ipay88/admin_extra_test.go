package ipay88

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

type errorReadCloser struct{}

func (e *errorReadCloser) Read([]byte) (int, error) { return 0, errors.New("read error") }
func (e *errorReadCloser) Close() error             { return nil }

func TestAdminPayPageMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-1", nil)
	rec := httptest.NewRecorder()
	p.adminPayPageHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}

func TestAdminPayPageInvalidStoredAmount(t *testing.T) {
	p := newTestProvider(t)
	p.saveRequest(PaymentRequest{RefNo: "REF-BAD-AMT", Amount: "abc", Currency: "MYR"})
	req := httptest.NewRequest(http.MethodGet, "/_admin/ipay88/pay/REF-BAD-AMT", nil)
	req.SetPathValue("refNo", "REF-BAD-AMT")
	rec := httptest.NewRecorder()
	p.adminPayPageHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", rec.Code)
	}
}

func TestAdminPayActionMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodGet, "/_admin/ipay88/pay/REF-1", nil)
	req.SetPathValue("refNo", "REF-1")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("want 405, got %d", rec.Code)
	}
}

func TestAdminPayActionInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-1", strings.NewReader("%"))
	req.SetPathValue("refNo", "REF-1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestAdminPayActionNotFound(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("outcome", "pay")
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/missing", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "missing")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", rec.Code)
	}
}

func TestAdminPayActionDefaultPaymentMethod(t *testing.T) {
	p := newTestProvider(t)
	client, rt := newFakeClient("RECEIVEOK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-DEFAULT-PAY")

	form := url.Values{}
	form.Set("outcome", "pay")
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-DEFAULT-PAY", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-DEFAULT-PAY")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if len(rt.requests) != 1 {
		t.Fatalf("want 1 backend post, got %d", len(rt.requests))
	}
	body := readBody(rt.requests[0])
	if !strings.Contains(body, "PaymentId=2") {
		t.Errorf("want default PaymentId=2, got %q", body)
	}
}

func TestAdminPayActionTransitionConflict(t *testing.T) {
	p := newTestProvider(t)
	client, _ := newFakeClient("RECEIVEOK")
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-CONFLICT")

	payForm := url.Values{}
	payForm.Set("outcome", "pay")
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-CONFLICT", strings.NewReader(payForm.Encode()))
	req.SetPathValue("refNo", "REF-CONFLICT")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.adminPayActionHandler().ServeHTTP(httptest.NewRecorder(), req)

	cancelForm := url.Values{}
	cancelForm.Set("outcome", "cancel")
	req = httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-CONFLICT", strings.NewReader(cancelForm.Encode()))
	req.SetPathValue("refNo", "REF-CONFLICT")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("want 409, got %d", rec.Code)
	}
}

func TestAdminPayActionBackendPostError(t *testing.T) {
	p := newTestProvider(t)
	client := &http.Client{Transport: &fakeRoundTripper{err: errors.New("boom")}}
	p.SetHTTPClient(client)
	seedPaymentRequest(t, p, "REF-BACKERR")

	form := url.Values{}
	form.Set("outcome", "pay")
	req := httptest.NewRequest(http.MethodPost, "/_admin/ipay88/pay/REF-BACKERR", strings.NewReader(form.Encode()))
	req.SetPathValue("refNo", "REF-BACKERR")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	p.adminPayActionHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", rec.Code)
	}
}

func TestTransitionTransactionNotFound(t *testing.T) {
	p := newTestProvider(t)
	err := p.transitionTransaction("missing", engine.TransactionStatusPaid)
	if !errors.Is(err, engine.ErrInvalidTransition) {
		t.Errorf("want ErrInvalidTransition, got %v", err)
	}
}

func TestTransitionTransactionInvalidTransition(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-TRANS-ERR")
	if err := p.transitionTransaction("REF-TRANS-ERR", engine.TransactionStatusPaid); err != nil {
		t.Fatalf("first transition: %v", err)
	}
	err := p.transitionTransaction("REF-TRANS-ERR", engine.TransactionStatusUnpaid)
	if !errors.Is(err, engine.ErrInvalidTransition) {
		t.Errorf("want ErrInvalidTransition, got %v", err)
	}
}

func TestTransitionTransactionStoreError(t *testing.T) {
	p := newTestProvider(t)
	p.SetStore(&fakeErrorStore{getErr: errors.New("db down")})
	err := p.transitionTransaction("REF-1", engine.TransactionStatusPaid)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTransitionTransactionCreateOrGetError(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-CREATE-ERR")
	p.SetStore(&fakeErrorStore{createErr: errors.New("db down")})
	err := p.transitionTransaction("REF-CREATE-ERR", engine.TransactionStatusPaid)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPostBackendCallbackInvalidURL(t *testing.T) {
	p := newTestProvider(t)
	req := PaymentRequest{BackendURL: "http://[::1]:namedport/"}
	err := p.postBackendCallback(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestPostBackendCallbackHTTPClientError(t *testing.T) {
	p := newTestProvider(t)
	client := &http.Client{Transport: &fakeRoundTripper{err: errors.New("network error")}}
	p.SetHTTPClient(client)
	req := PaymentRequest{BackendURL: "http://example.com/backend", Amount: "12.50", Currency: "MYR"}
	err := p.postBackendCallback(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPostBackendCallbackReadError(t *testing.T) {
	p := newTestProvider(t)
	client := &http.Client{Transport: &fakeRoundTripper{
		response: &http.Response{StatusCode: http.StatusOK, Body: &errorReadCloser{}, Header: make(http.Header)},
	}}
	p.SetHTTPClient(client)
	req := PaymentRequest{BackendURL: "http://example.com/backend", Amount: "12.50", Currency: "MYR"}
	err := p.postBackendCallback(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResponseValuesPaymentIDFallback(t *testing.T) {
	req := PaymentRequest{RefNo: "REF-1", PaymentID: "7", Amount: "12.50", Currency: "MYR"}
	v := responseValues(req, "M00001", "", "12.50", "MYR", "1", "key")
	if v.Get("PaymentId") != "7" {
		t.Errorf("want payment id 7, got %q", v.Get("PaymentId"))
	}

	req2 := PaymentRequest{RefNo: "REF-2", SelectedPaymentID: "5", PaymentID: "7", Amount: "12.50", Currency: "MYR"}
	v2 := responseValues(req2, "M00001", "", "12.50", "MYR", "1", "key")
	if v2.Get("PaymentId") != "5" {
		t.Errorf("want payment id 5, got %q", v2.Get("PaymentId"))
	}
}
