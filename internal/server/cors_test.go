package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/fawry"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestCORSPreflightReturnsHeaders(t *testing.T) {
	p := fawry.NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}

	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": p},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST"},
			AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		},
	})

	req := httptest.NewRequest(http.MethodOptions, "/_admin/transactions", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: want 204, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("allow-origin: want http://localhost:3000, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("allow-credentials: want true, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST" {
		t.Errorf("allow-methods: want GET, POST, got %q", got)
	}
}

func TestCORSAddsHeadersToNormalRequest(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("allow-origin: want *, got %q", got)
	}
}

func TestCORSRejectsDisallowedOrigin(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://trusted.example"},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("allow-origin: want empty, got %q", got)
	}
}
