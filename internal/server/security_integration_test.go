package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/audit"
)

func TestSecurityIntegrationAdminRequiresAuth(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	router := NewRouter(RouterConfig{
		Auth: AuthConfig{
			Enabled:      true,
			Username:     "admin",
			PasswordHash: hash,
		},
		AuditStore: audit.NewMemoryStore(),
	})

	// Admin route without auth is rejected.
	req := httptest.NewRequest(http.MethodGet, "/_admin/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	// Admin route with valid basic auth is accepted.
	req2 := httptest.NewRequest(http.MethodGet, "/_admin/", nil)
	req2.SetBasicAuth("admin", "secret")
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rr2.Code)
	}
}

func TestSecurityIntegrationProviderRouteUnprotected(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	router := NewRouter(RouterConfig{
		Auth: AuthConfig{
			Enabled:      true,
			Username:     "admin",
			PasswordHash: hash,
		},
		AuditStore: audit.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestSecurityIntegrationHeadersInHardenedMode(t *testing.T) {
	router := NewRouter(RouterConfig{
		Auth: AuthConfig{
			Enabled: true,
			Token:   "tok",
		},
		Hardened:        true,
		SecurityHeaders: SecurityHeadersConfig{Enabled: true, TLS: false},
		AuditStore:      audit.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/", nil)
	req.Header.Set("Authorization", "Bearer tok")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Error("expected CSP header in hardened mode")
	}
}

func TestSecurityIntegrationRateLimitInHardenedMode(t *testing.T) {
	router := NewRouter(RouterConfig{
		Auth: AuthConfig{
			Enabled: true,
			Token:   "tok",
		},
		Hardened: true,
		RateLimit: RateLimiterConfig{
			Enabled:           true,
			RequestsPerMinute: 1,
			AdminOnly:         false,
		},
		AuditStore: audit.NewMemoryStore(),
	})

	// First provider request allowed.
	req1 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("first request: want %d, got %d", http.StatusOK, rr1.Code)
	}

	// Second provider request rate-limited in hardened mode.
	req2 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: want %d, got %d", http.StatusTooManyRequests, rr2.Code)
	}
}

func TestSecurityIntegrationRateLimitLogsAuditEvent(t *testing.T) {
	store := audit.NewMemoryStore()
	router := NewRouter(RouterConfig{
		Auth: AuthConfig{
			Enabled: true,
			Token:   "tok",
		},
		Hardened: true,
		RateLimit: RateLimiterConfig{
			Enabled:           true,
			RequestsPerMinute: 1,
			AdminOnly:         false,
		},
		AuditStore: store,
	})

	// Exhaust the single request allowance.
	req1 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(httptest.NewRecorder(), req1)

	req2 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(httptest.NewRecorder(), req2)

	// Audit logging is asynchronous; poll briefly for the event.
	var found bool
	for i := 0; i < 50; i++ {
		events, err := store.List(100, 0)
		if err != nil {
			t.Fatalf("list audit events: %v", err)
		}
		for _, e := range events {
			if e.Action == SecurityEventRateLimit {
				found = true
				break
			}
		}
		if found {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !found {
		t.Errorf("expected audit event %q", SecurityEventRateLimit)
	}
}
