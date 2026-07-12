package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddlewareDisabled(t *testing.T) {
	handler := AuthMiddleware(AuthConfig{})(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthMiddlewareProviderRouteUnprotected(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "tok"}
	handler := AuthMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodPost, "/fawry/charge", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("provider route must not require auth: want %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthMiddlewareAdminRouteRequiresAuth(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "tok"}
	handler := AuthMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddlewareBearerToken(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "tok"}
	handler := AuthMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.Header.Set("Authorization", "Bearer tok")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthMiddlewareBasicAuth(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	cfg := AuthConfig{Enabled: true, Username: "admin", PasswordHash: hash}
	handler := AuthMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.SetBasicAuth("admin", "secret")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthMiddlewareInvalidCredentials(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "tok"}
	handler := AuthMiddleware(cfg)(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestIsAdminRoute(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/_admin", true},
		{"/_admin/", true},
		{"/_admin/api/transactions", true},
		{"/fawry/charge", false},
		{"/v1/checkout/sessions", false},
		{"/healthz", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := IsAdminRoute(tt.path); got != tt.want {
				t.Errorf("IsAdminRoute(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsProviderSimulationRoute(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/_admin/fawry-escape", true},
		{"/_admin/senangpay-escape", true},
		{"/_admin/billplz/pay/abc-123", true},
		{"/_admin/toyyibpay/pay/abc-123", true},
		{"/_admin/ipay88/pay/REF-1", true},
		{"/_admin/stripe/payment_intent/pi_test", true},
		{"/_admin/stripe/success", true},
		{"/_admin/stripe/fail", true},
		{"/_admin/stripe/cancel", true},
		{"/_admin", false},
		{"/_admin/config", false},
		{"/_admin/transactions", false},
		{"/_admin/stripe/webhooks", false},
		{"/fawry/charge", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isProviderSimulationRoute(tt.path); got != tt.want {
				t.Errorf("isProviderSimulationRoute(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestAuthMiddlewareProviderSimulationRouteUnprotected(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "admintok"}
	handler := AuthMiddleware(cfg)(okHandler())

	paths := []string{
		"/_admin/fawry-escape",
		"/_admin/senangpay-escape",
		"/_admin/billplz/pay/abc-123",
		"/_admin/toyyibpay/pay/abc-123",
		"/_admin/ipay88/pay/REF-1",
		"/_admin/stripe/payment_intent/pi_test",
		"/_admin/stripe/success",
		"/_admin/stripe/fail",
		"/_admin/stripe/cancel",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("simulation route %q must not require auth: want %d, got %d", path, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if !verifyPassword("secret", hash) {
		t.Error("expected password to verify")
	}
	if verifyPassword("wrong", hash) {
		t.Error("expected wrong password to fail")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	if _, err := HashPassword(""); err == nil {
		t.Error("expected error for empty password")
	}
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func roleRecordingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Role", string(RoleFromContext(r.Context())))
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthMiddlewareAdminRole(t *testing.T) {
	cfg := AuthConfig{Enabled: true, Username: "admin", Token: "admintok"}
	handler := AuthMiddleware(cfg)(roleRecordingHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.Header.Set("Authorization", "Bearer admintok")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("X-Role") != string(RoleAdmin) {
		t.Errorf("want role %q, got %q", RoleAdmin, rr.Header().Get("X-Role"))
	}
}

func TestAuthMiddlewareViewerRole(t *testing.T) {
	cfg := AuthConfig{
		Viewer: ViewerAuthConfig{Enabled: true, Username: "viewer", Token: "viewertok"},
	}
	handler := AuthMiddleware(cfg)(roleRecordingHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.Header.Set("Authorization", "Bearer viewertok")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("X-Role") != string(RoleViewer) {
		t.Errorf("want role %q, got %q", RoleViewer, rr.Header().Get("X-Role"))
	}
}

func TestAuthMiddlewareNoAuthGrantsAdminRole(t *testing.T) {
	handler := AuthMiddleware(AuthConfig{})(roleRecordingHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("X-Role") != string(RoleAdmin) {
		t.Errorf("want role %q, got %q", RoleAdmin, rr.Header().Get("X-Role"))
	}
}

func TestAuthMiddlewareViewerCannotUseAdminToken(t *testing.T) {
	cfg := AuthConfig{
		Enabled:  true,
		Username: "admin",
		Token:    "admintok",
		Viewer:   ViewerAuthConfig{Enabled: true, Username: "viewer", Token: "viewertok"},
	}
	handler := AuthMiddleware(cfg)(roleRecordingHandler())

	req := httptest.NewRequest(http.MethodGet, "/_admin", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	req := newAdminRequest(http.MethodPost, "/_admin/config/providers")
	rec := httptest.NewRecorder()
	if !RequireAdmin(req.Context(), rec, req) {
		t.Error("expected RequireAdmin to allow admin")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected no response written, got status %d", rec.Code)
	}
}

func TestRequireAdminDeniesViewer(t *testing.T) {
	req := newViewerRequest(http.MethodPost, "/_admin/config/providers")
	rec := httptest.NewRecorder()
	if RequireAdmin(req.Context(), rec, req) {
		t.Error("expected RequireAdmin to deny viewer")
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestRequireAdminDeniesAnonymous(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/_admin/config/providers", nil)
	rec := httptest.NewRecorder()
	if RequireAdmin(req.Context(), rec, req) {
		t.Error("expected RequireAdmin to deny anonymous")
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}
