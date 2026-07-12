package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

func newTestProvider(t *testing.T) *Provider {
	t.Helper()
	p := NewProvider()
	if err := p.Init(map[string]any{"user_secret_key": "test-secret"}); err != nil {
		t.Fatalf("init provider: %v", err)
	}
	p.SetBaseURL("http://localhost:9000")
	return p
}

func createTestCategory(t *testing.T, p *Provider) Category {
	t.Helper()
	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("categoryName", "Test Category")
	form.Set("categoryDescription", "desc")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createCategory", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.categoryCreateHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("create category status = %d", rec.Code)
	}
	return p.categories.by[firstCategoryCode(t, p)]
}

func firstCategoryCode(t *testing.T, p *Provider) string {
	t.Helper()
	for code := range p.categories.by {
		return code
	}
	t.Fatal("no category created")
	return ""
}

func createTestBill(t *testing.T, p *Provider) Bill {
	t.Helper()
	return createTestBillWithCallback(t, p, "http://merchant.example.com/callback")
}

func createTestBillWithCallback(t *testing.T, p *Provider, callbackURL string) Bill {
	t.Helper()
	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("categoryCode", "CAT123")
	form.Set("billName", "Test Bill")
	form.Set("billDescription", "desc")
	form.Set("billAmount", "1000")
	form.Set("billReturnUrl", "http://merchant.example.com/return")
	form.Set("billCallbackUrl", callbackURL)
	form.Set("billPaymentChannel", "2")
	form.Set("billTo", "John")
	form.Set("billEmail", "john@example.com")
	form.Set("billPhone", "0123456789")
	form.Set("billExternalReferenceNo", "ORDER-1")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createBill", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.billCreateHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("create bill status = %d, body = %s", rec.Code, rec.Body.String())
	}
	for code := range p.bills.byCode {
		return p.bills.byCode[code]
	}
	t.Fatal("no bill created")
	return Bill{}
}

func withCSRF(req *http.Request) *http.Request {
	return req.WithContext(httputil.WithCSRFToken(req.Context(), "test-csrf-token"))
}
