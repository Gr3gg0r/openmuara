package toyyibpay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestReturnHandlerValidation(t *testing.T) {
	p := newTestProvider(t)

	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("billName", "Bad")
	form.Set("billAmount", "100")
	form.Set("billReturnUrl", "://bad")
	form.Set("billCallbackUrl", "http://example.com/callback")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var resp BillCreateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	badBill := resp.Bill

	cases := []struct {
		name string
		path string
		want int
	}{
		{"missing_billcode", "/toyyibpay/return", http.StatusNotFound},
		{"invalid_return_url", "/toyyibpay/return?billcode=" + badBill.BillCode, http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			p.returnHandler().ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestBuildReturnURL(t *testing.T) {
	cases := []struct {
		name string
		base string
		q    url.Values
		want string
	}{
		{
			name: "merges_params",
			base: "http://merchant.example.com/return?existing=keep",
			q: func() url.Values {
				v := url.Values{}
				v.Set("status_id", "1")
				v.Set("billcode", "B")
				v.Set("order_id", "O")
				v.Set("transaction_id", "T")
				v.Set("msg", "ok")
				return v
			}(),
			want: "http://merchant.example.com/return?billcode=B&existing=keep&msg=ok&order_id=O&status_id=1&transaction_id=T",
		},
		{
			name: "skips_empty",
			base: "http://merchant.example.com/return",
			q: func() url.Values {
				v := url.Values{}
				v.Set("status_id", "1")
				return v
			}(),
			want: "http://merchant.example.com/return?status_id=1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildReturnURL(tc.base, tc.q)
			if err != nil {
				t.Fatalf("buildReturnURL: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuildReturnURLError(t *testing.T) {
	_, err := buildReturnURL("://bad", url.Values{})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
