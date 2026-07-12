package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("categoryName", "Donations")
	form.Set("categoryDescription", "Monthly donations")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createCategory", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.categoryCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if len(p.categories.by) != 1 {
		t.Fatal("category was not stored")
	}
}

func TestCreateCategoryInvalidSecret(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("userSecretKey", "wrong")
	form.Set("categoryName", "Donations")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/createCategory", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.categoryCreateHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want unauthorized", rec.Code)
	}
}

func TestGetCategoryDetails(t *testing.T) {
	p := newTestProvider(t)
	cat := createTestCategory(t, p)

	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("categoryCode", cat.CategoryCode)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/getCategoryDetails", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.categoryDetailsHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), cat.CategoryCode) {
		t.Fatal("response did not contain category code")
	}
}

func TestGetCategoryDetailsNotFound(t *testing.T) {
	p := newTestProvider(t)
	form := url.Values{}
	form.Set("userSecretKey", p.secret)
	form.Set("categoryCode", "missing")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/index.php/api/getCategoryDetails", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	p.categoryDetailsHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want not found", rec.Code)
	}
}
