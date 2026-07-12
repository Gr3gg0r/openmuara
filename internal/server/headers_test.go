package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	cfg := SecurityHeadersConfig{Enabled: true, TLS: true}
	handler := SecurityHeadersMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want %d, got %d", http.StatusOK, rr.Code)
	}

	headers := rr.Header()
	if got := headers.Get("Content-Security-Policy"); got != "default-src 'self'; img-src 'self' data:" {
		t.Errorf("CSP: want %q, got %q", "default-src 'self'; img-src 'self' data:", got)
	}
	if got := headers.Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("X-Frame-Options: want DENY, got %q", got)
	}
	if got := headers.Get("Strict-Transport-Security"); got == "" {
		t.Error("expected HSTS header when TLS is enabled")
	}
}

func TestSecurityHeadersMiddlewareProviderRoute(t *testing.T) {
	cfg := SecurityHeadersConfig{Enabled: true, TLS: true}
	handler := SecurityHeadersMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/fawry/charge", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Security-Policy") != "" {
		t.Error("provider route should not get admin security headers")
	}
}

func TestSecurityHeadersMiddlewareDisabled(t *testing.T) {
	cfg := SecurityHeadersConfig{Enabled: false}
	handler := SecurityHeadersMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Security-Policy") != "" {
		t.Error("headers should not be set when disabled")
	}
}
