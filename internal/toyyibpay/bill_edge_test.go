package toyyibpay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateBillValidation(t *testing.T) {
	p := newTestProvider(t)
	base := func() url.Values {
		v := url.Values{}
		v.Set("userSecretKey", p.secret)
		v.Set("billName", "Test")
		v.Set("billReturnUrl", "http://example.com/return")
		v.Set("billCallbackUrl", "http://example.com/callback")
		v.Set("billAmount", "1000")
		return v
	}

	cases := []struct {
		name string
		edit func(url.Values)
	}{
		{"missing_amount", func(v url.Values) { v.Del("billAmount") }},
		{"zero_amount", func(v url.Values) { v.Set("billAmount", "0") }},
		{"negative_amount", func(v url.Values) { v.Set("billAmount", "-1") }},
		{"invalid_amount", func(v url.Values) { v.Set("billAmount", "abc") }},
		{"missing_name", func(v url.Values) { v.Del("billName") }},
		{"missing_return", func(v url.Values) { v.Del("billReturnUrl") }},
		{"missing_callback", func(v url.Values) { v.Del("billCallbackUrl") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := base()
			tc.edit(v)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(v.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			p.billCreateHandler().ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestCreateBillDefaults(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "secret", "category_code": "CAT-DEF"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetBaseURL("http://localhost:9000")

	form := url.Values{}
	form.Set("userSecretKey", "secret")
	form.Set("billName", "Test")
	form.Set("billAmount", "500")
	form.Set("billReturnUrl", "http://example.com/return")
	form.Set("billCallbackUrl", "http://example.com/callback")
	form.Set("billPriceSetting", "1")
	form.Set("billPayorInfo", "1")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var resp BillCreateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Bill.CategoryCode != "CAT-DEF" {
		t.Fatalf("category = %q, want CAT-DEF", resp.Bill.CategoryCode)
	}
	if resp.Bill.BillPaymentChannel != "2" {
		t.Fatalf("channel = %q, want 2", resp.Bill.BillPaymentChannel)
	}
	if resp.Bill.BillPriceSetting != "1" || resp.Bill.BillPayorInfo != "1" {
		t.Fatalf("defaults not preserved: price=%s payor=%s", resp.Bill.BillPriceSetting, resp.Bill.BillPayorInfo)
	}
}

func TestRecordBillTransactionNoStore(t *testing.T) {
	p := newTestProvider(t)
	p.store = nil

	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("billName", "Test")
	form.Set("billAmount", "100")
	form.Set("billReturnUrl", "http://example.com/return")
	form.Set("billCallbackUrl", "http://example.com/callback")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
