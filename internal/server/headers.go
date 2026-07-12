package server

import "net/http"

// SecurityHeadersConfig holds security header settings.
type SecurityHeadersConfig struct {
	Enabled bool
	TLS     bool
}

// SecurityHeadersMiddleware adds security headers to admin responses.
func SecurityHeadersMiddleware(cfg SecurityHeadersConfig) Middleware {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAdminRoute(r.URL.Path) {
				h := w.Header()
				h.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:")
				h.Set("X-Frame-Options", "DENY")
				h.Set("X-Content-Type-Options", "nosniff")
				h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
				if cfg.TLS {
					h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
