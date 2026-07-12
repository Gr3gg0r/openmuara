package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodNotAllowed(t *testing.T) {
	p := newTestProvider(t)

	cases := []struct {
		name   string
		h      http.Handler
		method string
		path   string
	}{
		{"createBill", p.billCreateHandler(), http.MethodGet, "/index.php/api/createBill"},
		{"getBillTransactions", p.billTransactionsHandler(), http.MethodGet, "/index.php/api/getBillTransactions"},
		{"inactiveBill", p.billInactiveHandler(), http.MethodGet, "/index.php/api/inactiveBill"},
		{"createCategory", p.categoryCreateHandler(), http.MethodGet, "/index.php/api/createCategory"},
		{"getCategoryDetails", p.categoryDetailsHandler(), http.MethodGet, "/index.php/api/getCategoryDetails"},
		{"payPage", p.payPageHandler(), http.MethodPost, "/_admin/toyyibpay/pay/xyz"},
		{"payPageAction", p.payPageActionHandler(), http.MethodGet, "/_admin/toyyibpay/pay/xyz"},
		{"return", p.returnHandler(), http.MethodPost, "/toyyibpay/return"},
		{"webhook", p.WebhookHandler(), http.MethodGet, "/toyyibpay/webhook"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			tc.h.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405", rec.Code)
			}
		})
	}
}

func TestChargeAndEscapeHandlers(t *testing.T) {
	p := newTestProvider(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/charge", nil)
	p.ChargeHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("ChargeHandler status = %d, want 200", rec.Code)
	}

	if p.EscapeHandler() != nil {
		t.Fatal("EscapeHandler should return nil")
	}
}
