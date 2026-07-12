package billplz_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
)

func TestRedirectBuildsSignedQueryString(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	// Mark bill paid so paid=true appears in redirect.
	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, mustPayRequest(b.ID))

	redirectHandler := routeFor(t, p, http.MethodGet, "/billplz/redirect")
	req := httptest.NewRequest(http.MethodGet, "/billplz/redirect?billplz[id]="+b.ID, nil)
	rec = httptest.NewRecorder()

	redirectHandler.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("status: want 302, got %d, body: %s", rec.Code, rec.Body.String())
	}

	loc := rec.Header().Get("Location")
	if loc == "" {
		t.Fatal("location header missing")
	}
	if !strings.HasPrefix(loc, b.RedirectURL) {
		t.Errorf("location should start with redirect_url, got %q", loc)
	}
	if !strings.Contains(loc, "x_signature=") {
		t.Error("location missing x_signature")
	}
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("parse location: %v", err)
	}
	q := u.Query()
	if q.Get("billplz[paid]") != "true" {
		t.Errorf("billplz[paid]: want true, got %q", q.Get("billplz[paid]"))
	}
	if q.Get("x_signature") == "" {
		t.Error("x_signature is empty")
	}

	query := parseRedirectQuery(loc)
	if !billplz.VerifyRedirectSignature(query, "muara-billplz-xsig-key") {
		t.Error("redirect signature verification failed")
	}
}

func TestRedirectMissingBillID(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/billplz/redirect")
	req := httptest.NewRequest(http.MethodGet, "/billplz/redirect", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestRedirectBillNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/billplz/redirect")
	req := httptest.NewRequest(http.MethodGet, "/billplz/redirect?billplz[id]=missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}
