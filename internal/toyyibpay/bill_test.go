package toyyibpay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateBill(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	if bill.BillCode == "" {
		t.Fatal("bill code not generated")
	}
	if bill.BillStatus != "1" {
		t.Fatalf("bill status = %q, want 1", bill.BillStatus)
	}
	if !strings.HasPrefix(bill.BillPaymentLink, p.baseURL+"/_admin/toyyibpay/pay/") {
		t.Fatalf("payment link = %q", bill.BillPaymentLink)
	}
	if bill.OrderID != "ORDER-1" {
		t.Fatalf("order id = %q, want ORDER-1", bill.OrderID)
	}
}

func TestCreateBillInvalidSecret(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("userSecretKey", "wrong")
	form.Set("billName", "Test")
	form.Set("billAmount", "1000")
	form.Set("billReturnUrl", "http://example.com/return")
	form.Set("billCallbackUrl", "http://example.com/callback")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want unauthorized", rec.Code)
	}
}

func TestCreateBillMissingRequiredFields(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("billName", "Test")
	form.Set("billAmount", "1000")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want bad request", rec.Code)
	}
}

func TestGetBillTransactions(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("billCode", bill.BillCode)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/getBillTransactions", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billTransactionsHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var txs []BillTransaction
	if err := json.Unmarshal(rec.Body.Bytes(), &txs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("transactions = %d, want 1", len(txs))
	}
}

func TestInactiveBill(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("billCode", bill.BillCode)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/inactiveBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billInactiveHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	updated, _ := p.bills.GetByCode(bill.BillCode)
	if updated.BillStatus != "2" {
		t.Fatalf("bill status = %q, want 2", updated.BillStatus)
	}
}
