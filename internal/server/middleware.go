package server

import (
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

const maxRequestBodySize = 1 << 20 // 1 MiB

// Middleware is an alias for functions that wrap an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain composes multiple middleware functions into one.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// MaxBodySizeMiddleware rejects requests whose body exceeds maxRequestBodySize.
func MaxBodySizeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxRequestBodySize {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		h.ServeHTTP(w, r)
	})
}
