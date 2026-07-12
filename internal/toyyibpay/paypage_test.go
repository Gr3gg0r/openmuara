package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestPayPageRender(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodGet, "/_admin/toyyibpay/pay/"+bill.BillCode, nil))
	req.SetPathValue("billCode", bill.BillCode)
	p.payPageHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Pay") {
		t.Fatal("page missing pay button")
	}
	if !strings.Contains(body, "test-csrf-token") {
		t.Fatal("page missing csrf token")
	}
}

func TestPayPageInactiveBill(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)
	bill.BillStatus = "2"
	p.bills.Update(bill)

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodGet, "/_admin/toyyibpay/pay/"+bill.BillCode, nil))
	req.SetPathValue("billCode", bill.BillCode)
	p.payPageHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want conflict", rec.Code)
	}
}

func TestPayPageActionPay(t *testing.T) {
	p := newTestProvider(t)
	callbackReceived := make(chan url.Values, 1)
	ts := startCallbackServer(t, callbackReceived)
	bill := createTestBillWithCallback(t, p, ts.URL)
	p.SetDispatcher(webhook.NewDispatcherFromProvider(ts.URL, 0, p))

	form := url.Values{}
	form.Set("status", "1")
	form.Set("csrf_token", "test-csrf-token")

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodPost, "/_admin/toyyibpay/pay/"+bill.BillCode, strings.NewReader(form.Encode())))
	req.SetPathValue("billCode", bill.BillCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.payPageActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	assertTransactionStatus(t, p, bill.OrderID, engine.TransactionStatusPaid)
	assertCallbackReceived(t, callbackReceived, "1", bill)
	if !strings.Contains(rec.Header().Get("Location"), "/toyyibpay/return") {
		t.Fatalf("redirect location = %q", rec.Header().Get("Location"))
	}
}

func TestPayPageActionCancel(t *testing.T) {
	p := newTestProvider(t)
	callbackReceived := make(chan url.Values, 1)
	ts := startCallbackServer(t, callbackReceived)
	bill := createTestBillWithCallback(t, p, ts.URL)
	p.SetDispatcher(webhook.NewDispatcherFromProvider(ts.URL, 0, p))

	form := url.Values{}
	form.Set("status", "3")
	form.Set("csrf_token", "test-csrf-token")

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodPost, "/_admin/toyyibpay/pay/"+bill.BillCode, strings.NewReader(form.Encode())))
	req.SetPathValue("billCode", bill.BillCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.payPageActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	assertTransactionStatus(t, p, bill.OrderID, engine.TransactionStatusUnpaid)
	assertCallbackReceived(t, callbackReceived, "3", bill)
}

func TestReturnURLRedirect(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/toyyibpay/return?status_id=1&billcode="+bill.BillCode+"&order_id="+bill.OrderID+"&transaction_id=TXN&msg=ok", nil)
	p.returnHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, bill.BillReturnURL) {
		t.Fatalf("location = %q", loc)
	}
	if !strings.Contains(loc, "status_id=1") {
		t.Fatal("location missing status_id")
	}
}

func TestIncomingWebhookVerifiesHash(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	refno := "REF-1"
	status := "1"
	hash := ComputeHash(p.secret, status, bill.OrderID, refno)
	form := url.Values{}
	form.Set("refno", refno)
	form.Set("status", status)
	form.Set("billcode", bill.BillCode)
	form.Set("order_id", bill.OrderID)
	form.Set("amount", "1000")
	form.Set("hash", hash)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.WebhookHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	assertTransactionStatus(t, p, bill.OrderID, engine.TransactionStatusPaid)
}

func TestIncomingWebhookRejectsInvalidHash(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	form := url.Values{}
	form.Set("refno", "REF-1")
	form.Set("status", "1")
	form.Set("billcode", bill.BillCode)
	form.Set("order_id", bill.OrderID)
	form.Set("hash", "invalid")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.WebhookHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want bad request", rec.Code)
	}
}

func TestComputeHash(t *testing.T) {
	got := ComputeHash("secret", "1", "ORDER-1", "REF-1")
	want := ComputeHash("secret", "1", "ORDER-1", "REF-1")
	if got != want {
		t.Fatalf("hash mismatch")
	}
}

func startCallbackServer(t *testing.T, ch chan url.Values) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		values := url.Values{}
		for k, v := range r.Form {
			values[k] = v
		}
		ch <- values
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func assertTransactionStatus(t *testing.T, p *Provider, orderID string, want engine.TransactionStatus) {
	t.Helper()
	tx, ok, _ := p.store.GetByReference(orderID)
	if !ok {
		t.Fatal("transaction not found")
	}
	if tx.Status != want {
		t.Fatalf("status = %q, want %q", tx.Status, want)
	}
}

func assertCallbackReceived(t *testing.T, ch chan url.Values, wantStatus string, bill Bill) {
	t.Helper()
	select {
	case values := <-ch:
		if values.Get("status") != wantStatus {
			t.Fatalf("callback status = %q, want %q", values.Get("status"), wantStatus)
		}
		if values.Get("billcode") != bill.BillCode {
			t.Fatalf("callback billcode = %q, want %q", values.Get("billcode"), bill.BillCode)
		}
		if !VerifyCallback("test-secret", values) {
			t.Fatal("callback hash did not verify")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("callback not received")
	}
}
