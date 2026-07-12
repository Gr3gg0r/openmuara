package server

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcryptCost is the default cost for hashing admin passwords.
	// It is configurable only via this constant to avoid unsafe config values.
	bcryptCost = 10
)

// AuthConfig holds admin and viewer authentication settings for the HTTP server.
type AuthConfig struct {
	Enabled      bool
	Username     string
	PasswordHash string
	Token        string
	Viewer       ViewerAuthConfig
}

// ViewerAuthConfig holds read-only authentication settings.
type ViewerAuthConfig struct {
	Enabled      bool
	Username     string
	PasswordHash string
	Token        string
}

// IsAuthEnabled reports whether admin authentication is configured.
func (c AuthConfig) IsAuthEnabled() bool {
	return c.Enabled && c.Username != "" && (c.PasswordHash != "" || c.Token != "")
}

// IsViewerEnabled reports whether viewer authentication is configured.
func (c AuthConfig) IsViewerEnabled() bool {
	return c.Viewer.Enabled && c.Viewer.Username != "" && (c.Viewer.PasswordHash != "" || c.Viewer.Token != "")
}

// AuthMiddleware returns a middleware that enforces admin authentication on
// protected routes. Provider emulation routes and provider simulation pages
// (payment/escape pages under /_admin) are never protected, because they are
// part of the payment flow that testers and redirected browsers must use.
// When no authentication is configured, admin routes are treated as fully
// authorized to preserve backward compatibility.
func AuthMiddleware(cfg AuthConfig) Middleware {
	authEnabled := cfg.IsAuthEnabled() || cfg.IsViewerEnabled()
	if !authEnabled {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if isAdminRoute(r.URL.Path) {
					ctx := WithRole(r.Context(), RoleAdmin)
					r = r.WithContext(ctx)
				}
				next.ServeHTTP(w, r)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isAdminRoute(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Provider simulation pages live under /_admin for historical reasons
			// but are part of the public payment-emulation flow.
			if isProviderSimulationRoute(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			role, ok := authenticate(r, cfg)
			if !ok {
				logSecurityEventFromRequest(r, SecurityEventAuthFailure, "admin route access denied")
				headers := w.Header()
				headers.Set("WWW-Authenticate", `Basic realm="muara-admin"`)
				headers.Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}

			ctx := WithRole(r.Context(), role)
			r = r.WithContext(ctx)
			logSecurityEventFromRequest(r, SecurityEventAuthSuccess, "admin route access granted")
			next.ServeHTTP(w, r)
		})
	}
}

// isAdminRoute reports whether the request path is part of the admin surface.
func isAdminRoute(path string) bool {
	return path == "/_admin" || strings.HasPrefix(path, "/_admin/")
}

// isProviderSimulationRoute reports whether the path is a provider payment or
// escape page that is part of the public emulation flow. These paths are kept
// under /_admin/ for backward compatibility but must remain accessible without
// admin credentials so that redirected browsers and testers can complete payments.
func isProviderSimulationRoute(path string) bool {
	switch {
	case path == "/_admin/fawry-escape":
		return true
	case path == "/_admin/senangpay-escape":
		return true
	case strings.HasPrefix(path, "/_admin/billplz/pay/"):
		return true
	case strings.HasPrefix(path, "/_admin/toyyibpay/pay/"):
		return true
	case strings.HasPrefix(path, "/_admin/ipay88/pay/"):
		return true
	case strings.HasPrefix(path, "/_admin/stripe/payment_intent/"):
		return true
	case path == "/_admin/stripe/success",
		path == "/_admin/stripe/fail",
		path == "/_admin/stripe/cancel":
		return true
	}
	return false
}

// authenticate checks basic auth or bearer token against the configured credentials.
// It returns the matching role and true on success, or RoleAnonymous and false.
func authenticate(r *http.Request, cfg AuthConfig) (Role, bool) {
	if cfg.PasswordHash != "" {
		if user, pass, ok := r.BasicAuth(); ok {
			if user == cfg.Username && verifyPassword(pass, cfg.PasswordHash) {
				return RoleAdmin, true
			}
		}
	}

	if cfg.Token != "" {
		if header := r.Header.Get("Authorization"); header != "" {
			const prefix = "Bearer "
			if strings.HasPrefix(header, prefix) {
				token := strings.TrimPrefix(header, prefix)
				if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Token)) == 1 {
					return RoleAdmin, true
				}
			}
		}
	}

	if cfg.Viewer.PasswordHash != "" {
		if user, pass, ok := r.BasicAuth(); ok {
			if user == cfg.Viewer.Username && verifyPassword(pass, cfg.Viewer.PasswordHash) {
				return RoleViewer, true
			}
		}
	}

	if cfg.Viewer.Token != "" {
		if header := r.Header.Get("Authorization"); header != "" {
			const prefix = "Bearer "
			if strings.HasPrefix(header, prefix) {
				token := strings.TrimPrefix(header, prefix)
				if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Viewer.Token)) == 1 {
					return RoleViewer, true
				}
			}
		}
	}

	return RoleAnonymous, false
}

// RequireAdmin returns 403 if the request context does not hold RoleAdmin.
func RequireAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	if IsAdmin(ctx) {
		return true
	}
	logSecurityEventFromRequest(r, SecurityEventAuthFailure, "viewer denied admin operation")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"admin access required"}`))
	return false
}

// verifyPassword compares a plaintext password with a bcrypt hash.
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword generates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsAdminRoute is exported for testing.
var IsAdminRoute = isAdminRoute
