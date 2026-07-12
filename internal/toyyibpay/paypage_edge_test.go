package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestPayPageActionErrors(t *testing.T) {
	cases := []struct {
		name   string
		prep   func(*Provider)
		status string
		code   string
		want   int
	}{
		{
			name:   "invalid_status",
			prep:   nil,
			status: "2",
			want:   http.StatusBadRequest,
		},
		{
			name: "store_nil",
			prep: func(p *Provider) {
				p.store = nil
			},
			status: "1",
			want:   http.StatusInternalServerError,
		},
		{
			name: "tx_not_found",
			prep: func(p *Provider) {
				p.store = engine.NewMemoryStore()
			},
			status: "1",
			want:   http.StatusInternalServerError,
		},
		{
			name:   "nil_dispatcher",
			prep:   nil,
			status: "1",
			want:   http.StatusInternalServerError,
		},
		{
			name:   "missing_bill",
			prep:   nil,
			status: "1",
			code:   "missing",
			want:   http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := newTestProvider(t)
			bill := createTestBill(t, p)
			if tc.prep != nil {
				tc.prep(p)
			}

			code := bill.BillCode
			if tc.code != "" {
				code = tc.code
			}

			form := url.Values{}
			form.Set("status", tc.status)

			rec := httptest.NewRecorder()
			req := withCSRF(httptest.NewRequest(http.MethodPost, "/_admin/toyyibpay/pay/"+code, strings.NewReader(form.Encode())))
			req.SetPathValue("billCode", code)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			p.payPageActionHandler().ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestPayPageActionInvalidForm(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodPost, "/_admin/toyyibpay/pay/"+bill.BillCode, strings.NewReader("%zz=bad")))
	req.SetPathValue("billCode", bill.BillCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.payPageActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want bad request", rec.Code)
	}
}

func TestPayPageActionTransitionError(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	tx, _, _ := p.store.GetByReference(bill.OrderID)
	tx.Status = engine.TransactionStatusUnpaid
	if _, _, err := p.store.CreateOrGet(tx); err != nil {
		t.Fatalf("update tx: %v", err)
	}

	form := url.Values{}
	form.Set("status", "1")

	rec := httptest.NewRecorder()
	req := withCSRF(httptest.NewRequest(http.MethodPost, "/_admin/toyyibpay/pay/"+bill.BillCode, strings.NewReader(form.Encode())))
	req.SetPathValue("billCode", bill.BillCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.payPageActionHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want internal server error", rec.Code)
	}
}
