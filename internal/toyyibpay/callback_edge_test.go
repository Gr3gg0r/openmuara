package toyyibpay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestVerifyCallback(t *testing.T) {
	secret := "secret"
	order := "ORDER-1"
	refno := "REF-1"
	hash := ComputeHash(secret, "1", order, refno)

	cases := []struct {
		name string
		form url.Values
		want bool
	}{
		{
			name: "empty_hash",
			form: func() url.Values {
				v := url.Values{}
				v.Set("status", "1")
				v.Set("order_id", order)
				v.Set("refno", refno)
				return v
			}(),
			want: false,
		},
		{
			name: "mismatch",
			form: func() url.Values {
				v := url.Values{}
				v.Set("hash", "bad")
				v.Set("status", "1")
				v.Set("order_id", order)
				v.Set("refno", refno)
				return v
			}(),
			want: false,
		},
		{
			name: "match",
			form: func() url.Values {
				v := url.Values{}
				v.Set("hash", hash)
				v.Set("status", "1")
				v.Set("order_id", order)
				v.Set("refno", refno)
				return v
			}(),
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := VerifyCallback(secret, tc.form); got != tc.want {
				t.Fatalf("VerifyCallback = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestWebhookHandlerValidation(t *testing.T) {
	t.Run("invalid_form", func(t *testing.T) {
		p := newTestProvider(t)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader("%zz=bad"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want bad request", rec.Code)
		}
	})

	t.Run("missing_hash", func(t *testing.T) {
		p := newTestProvider(t)
		bill := createTestBill(t, p)

		form := url.Values{}
		form.Set("status", "1")
		form.Set("order_id", bill.OrderID)
		form.Set("refno", "REF-1")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want bad request", rec.Code)
		}
	})

	t.Run("missing_order_id", func(t *testing.T) {
		p := newTestProvider(t)
		_ = createTestBill(t, p)

		form := url.Values{}
		form.Set("status", "1")
		form.Set("refno", "REF-1")
		form.Set("hash", ComputeHash(p.secret, "1", "", "REF-1"))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want not found", rec.Code)
		}
	})

	t.Run("store_nil", func(t *testing.T) {
		p := newTestProvider(t)
		bill := createTestBill(t, p)
		p.store = nil

		form := url.Values{}
		form.Set("status", "1")
		form.Set("order_id", bill.OrderID)
		form.Set("refno", "REF-1")
		form.Set("hash", ComputeHash(p.secret, "1", bill.OrderID, "REF-1"))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want not found", rec.Code)
		}
	})

	t.Run("transaction_not_found", func(t *testing.T) {
		p := newTestProvider(t)
		bill := createTestBill(t, p)
		p.store = engine.NewMemoryStore()

		form := url.Values{}
		form.Set("status", "1")
		form.Set("order_id", bill.OrderID)
		form.Set("refno", "REF-1")
		form.Set("hash", ComputeHash(p.secret, "1", bill.OrderID, "REF-1"))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want not found", rec.Code)
		}
	})

	t.Run("status_unpaid", func(t *testing.T) {
		p := newTestProvider(t)
		bill := createTestBill(t, p)

		form := url.Values{}
		form.Set("status", "3")
		form.Set("order_id", bill.OrderID)
		form.Set("refno", "REF-1")
		form.Set("hash", ComputeHash(p.secret, "3", bill.OrderID, "REF-1"))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/toyyibpay/webhook", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.WebhookHandler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
		assertTransactionStatus(t, p, bill.OrderID, engine.TransactionStatusUnpaid)
	})
}

func TestDispatchCallbackNilDispatcher(t *testing.T) {
	p := newTestProvider(t)
	bill := createTestBill(t, p)

	err := p.dispatchCallback(context.Background(), bill, "1")
	if err == nil || !strings.Contains(err.Error(), "dispatcher not configured") {
		t.Fatalf("expected dispatcher error, got %v", err)
	}
}
