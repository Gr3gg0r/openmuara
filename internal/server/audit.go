package server

import (
	"net/http"

	"github.com/openmuara/openmuara/internal/audit"
)

// AuditMiddleware injects the audit logger into the request context.
func AuditMiddleware(logger audit.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := audit.NewContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
