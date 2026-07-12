package toyyibpay

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func assertErrcode(t *testing.T, err error, want errcode.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	var ec *errcode.Error
	if !errors.As(err, &ec) {
		t.Fatalf("expected *errcode.Error, got %T", err)
	}
	if ec.Code != want {
		t.Errorf("errcode: want %q, got %q", want, ec.Code)
	}
}

func newCreateBillHTTPRequest() *http.Request {
	form := url.Values{}
	form.Set("userSecretKey", "secret")
	form.Set("categoryCode", "cat")
	form.Set("billName", "Test")
	form.Set("billAmount", "1000")
	form.Set("billReturnUrl", "http://example.com/return")
	form.Set("billCallbackUrl", "http://example.com/callback")
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", nil)
	req.Form = form
	return req
}

func TestProviderInitMissingSecretKeyHasErrcode(t *testing.T) {
	p := NewProvider()
	err := p.Init(map[string]any{})
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestParseCreateBillInvalidAmountHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	req := newCreateBillHTTPRequest()
	req.Form.Set("billAmount", "0")
	_, err := p.parseCreateBill(req)
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestApplyIncomingCallbackMissingOrderIDHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := p.applyIncomingCallback(url.Values{"status": []string{"1"}})
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestApplyIncomingCallbackTransactionNotFoundHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := p.applyIncomingCallback(url.Values{"order_id": []string{"missing"}, "status": []string{"1"}})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestBuildCallbackPayloadMissingBillHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	_, err := p.buildCallbackPayload(context.Background(), provider.Transaction{Reference: "missing"})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestRecordPaymentOutcomeMissingTransactionHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := p.recordPaymentOutcome(Bill{OrderID: "missing"}, "1")
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestAuthenticateInvalidSecretHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	assertErrcode(t, errInvalidSecret, errcode.ESignatureMismatch)
}
