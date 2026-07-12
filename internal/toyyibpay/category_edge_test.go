package toyyibpay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCategoryCreateValidation(t *testing.T) {
	p := newTestProvider(t)

	t.Run("invalid_secret", func(t *testing.T) {
		form := url.Values{}
		form.Set("userSecretKey", "wrong")
		form.Set("categoryName", "X")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/index.php/api/createCategory", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.categoryCreateHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want unauthorized", rec.Code)
		}
	})

	t.Run("invalid_form", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/index.php/api/createCategory", strings.NewReader("%zz=bad"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.categoryCreateHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want bad request", rec.Code)
		}
	})
}

func TestCategoryDetailsValidation(t *testing.T) {
	p := newTestProvider(t)
	cat := createTestCategory(t, p)

	t.Run("invalid_secret", func(t *testing.T) {
		form := url.Values{}
		form.Set("userSecretKey", "wrong")
		form.Set("categoryCode", cat.CategoryCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/index.php/api/getCategoryDetails", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.categoryDetailsHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want unauthorized", rec.Code)
		}
	})

	t.Run("missing_category_code", func(t *testing.T) {
		form := url.Values{}
		form.Set("userSecretKey", p.secret)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/index.php/api/getCategoryDetails", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.categoryDetailsHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want bad request", rec.Code)
		}
	})

	t.Run("invalid_form", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/index.php/api/getCategoryDetails", strings.NewReader("%zz=bad"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		p.categoryDetailsHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want bad request", rec.Code)
		}
	})
}
