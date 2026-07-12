package server

import (
	"net/http"
	"strings"

	"github.com/openmuara/openmuara/internal/config"
)

// CORSMiddleware adds configurable cross-origin headers and handles preflight.
func CORSMiddleware(cfg config.CORSConfig) Middleware {
	allowedOrigins := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedOrigins[strings.ToLower(o)] = struct{}{}
	}
	wildcard := false
	if _, ok := allowedOrigins["*"]; ok {
		wildcard = true
	}

	methods := strings.Join(cfg.AllowedMethods, ", ")
	if methods == "" {
		methods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	}

	headers := strings.Join(cfg.AllowedHeaders, ", ")
	if headers == "" {
		headers = "Content-Type, Authorization, X-CSRF-Token, X-Request-ID"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowOrigin := ""

			switch {
			case wildcard && !cfg.AllowCredentials:
				allowOrigin = "*"
			case origin != "":
				if wildcard {
					allowOrigin = origin
				} else if _, ok := allowedOrigins[strings.ToLower(origin)]; ok {
					allowOrigin = origin
				}
			}

			if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", methods)
				w.Header().Set("Access-Control-Allow-Headers", headers)
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
