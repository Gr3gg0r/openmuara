package toyyibpay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestGetBillTransactionsValidation(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	cases := []struct {
		name string
		form url.Values
		want int
	}{
		{
			name: "invalid_secret",
			form: func() url.Values {
				v := url.Values{}
				v.Set("userSecretKey", "wrong")
				v.Set("billCode", bill.BillCode)
				return v
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "missing_bill",
			form: func() url.Values {
				v := url.Values{}
				v.Set("userSecretKey", p.secret)
				v.Set("billCode", "missing")
				return v
			}(),
			want: http.StatusNotFound,
		},
		{
			name: "no_bill_code",
			form: func() url.Values {
				v := url.Values{}
				v.Set("userSecretKey", p.secret)
				return v
			}(),
			want: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/index.php/api/getBillTransactions", strings.NewReader(tc.form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			p.billTransactionsHandler().ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestGetBillTransactionsUnpaidStatus(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	tx, _, _ := p.store.GetByReference(bill.OrderID)
	tx.Status = engine.TransactionStatusUnpaid
	if _, _, err := p.store.CreateOrGet(tx); err != nil {
		t.Fatalf("update tx: %v", err)
	}

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
	if len(txs) != 1 || txs[0].BillPaymentStatus != "3" {
		t.Fatalf("txs = %+v", txs)
	}
}

func TestGetBillTransactionsNoStore(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)
	p.store = nil

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
	if len(txs) != 0 {
		t.Fatalf("expected empty txs, got %+v", txs)
	}
}

func TestInactiveBillValidation(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	cases := []struct {
		name string
		form url.Values
		want int
	}{
		{
			name: "invalid_secret",
			form: func() url.Values {
				v := url.Values{}
				v.Set("userSecretKey", "wrong")
				v.Set("billCode", bill.BillCode)
				return v
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "missing_bill",
			form: func() url.Values {
				v := url.Values{}
				v.Set("userSecretKey", p.secret)
				v.Set("billCode", "missing")
				return v
			}(),
			want: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/index.php/api/inactiveBill", strings.NewReader(tc.form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			p.billInactiveHandler().ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}
