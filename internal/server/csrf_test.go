package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/fawry"
	"github.com/openmuara/openmuara/internal/provider"
)

func newCSRFFawryProvider(t *testing.T) provider.Provider {
	t.Helper()
	p := fawry.NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}
	return p
}

func setupCSRFRouter(t *testing.T) http.Handler {
	t.Helper()
	p := newCSRFFawryProvider(t)

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "r", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}
	if sp, ok := p.(interface{ SetStore(engine.TransactionStore) }); ok {
		sp.SetStore(store)
	}

	return NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": p},
		TransactionStore: store,
		CSRF:             config.CSRFConfig{Enabled: true},
	})
}

func extractCSRFToken(body string) string {
	const prefix = `name="csrf_token" value="`
	start := strings.Index(body, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.IndexByte(body[start:], '"')
	if end == -1 {
		return ""
	}
	return body[start : start+end]
}

func TestCSRFSetsCookieOnSafeRequest(t *testing.T) {
	router := setupCSRFRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	cookies := rec.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == csrfCookieName && c.Value != "" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected %s cookie to be set", csrfCookieName)
	}
}

func TestCSRFRejectsAdminMutationWithoutToken(t *testing.T) {
	router := setupCSRFRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d", rec.Code)
	}
}

func TestCSRFAcceptsAdminMutationWithMatchingToken(t *testing.T) {
	router := setupCSRFRouter(t)

	// Fetch the escape page to obtain a cookie and token.
	getReq := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET status: want 200, got %d", getRec.Code)
	}

	var cookieValue string
	for _, c := range getRec.Result().Cookies() {
		if c.Name == csrfCookieName {
			cookieValue = c.Value
		}
	}
	if cookieValue == "" {
		t.Fatal("expected csrf cookie value")
	}

	token := extractCSRFToken(getRec.Body.String())
	if token == "" {
		t.Fatal("expected csrf_token form field in escape page")
	}
	if token != cookieValue {
		t.Fatalf("form token %q does not match cookie %q", token, cookieValue)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost&status=PAID&csrf_token="+token))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.AddCookie(&http.Cookie{
		Name:     csrfCookieName,
		Value:    cookieValue,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	postRec := httptest.NewRecorder()

	router.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", postRec.Code, postRec.Body.String())
	}
}

func TestCSRFAcceptsHeaderToken(t *testing.T) {
	router := setupCSRFRouter(t)

	getReq := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	var cookieValue string
	for _, c := range getRec.Result().Cookies() {
		if c.Name == csrfCookieName {
			cookieValue = c.Value
		}
	}
	if cookieValue == "" {
		t.Fatal("expected csrf cookie from escape page")
	}

	postReq := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost&status=PAID"))
	postReq.Header.Set(csrfHeaderName, cookieValue)
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.AddCookie(&http.Cookie{
		Name:     csrfCookieName,
		Value:    cookieValue,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	postRec := httptest.NewRecorder()

	router.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", postRec.Code, postRec.Body.String())
	}
}

func TestCSRFCookieFlagsWithTLS(t *testing.T) {
	p := newCSRFFawryProvider(t)
	store := engine.NewMemoryStore()
	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": p},
		TransactionStore: store,
		CSRF:             config.CSRFConfig{Enabled: true},
		SecurityHeaders:  SecurityHeadersConfig{Enabled: true, TLS: true},
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	for _, c := range rec.Result().Cookies() {
		if c.Name == csrfCookieName {
			if !c.Secure {
				t.Errorf("expected Secure flag when TLS is enabled")
			}
			if c.SameSite != http.SameSiteLaxMode {
				t.Errorf("expected SameSite=Lax without admin auth, got %v", c.SameSite)
			}
			return
		}
	}
	t.Fatalf("expected %s cookie to be set", csrfCookieName)
}

func TestCSRFCookieSameSiteStrictWithAdminAuth(t *testing.T) {
	p := newCSRFFawryProvider(t)
	store := engine.NewMemoryStore()
	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": p},
		TransactionStore: store,
		CSRF:             config.CSRFConfig{Enabled: true},
		Auth: AuthConfig{
			Enabled:  true,
			Username: "admin",
			// #nosec G101 -- test fixture bcrypt hash of "admin"
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhmM6JGKpS4G3R1G2tH9Pp5kBjFu",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	req.SetBasicAuth("admin", "admin")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	for _, c := range rec.Result().Cookies() {
		if c.Name == csrfCookieName {
			if c.SameSite != http.SameSiteStrictMode {
				t.Errorf("expected SameSite=Strict when admin auth is enabled, got %v", c.SameSite)
			}
			return
		}
	}
	t.Fatalf("expected %s cookie to be set", csrfCookieName)
}
