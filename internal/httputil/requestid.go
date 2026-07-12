package httputil

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a private type for context keys.
type contextKey int

const traceIDKey contextKey = iota

// TraceIDHeader is the HTTP header used to propagate trace IDs.
const TraceIDHeader = "X-Trace-Id"

// TraceIDFromContext returns the trace ID stored in ctx, or an empty string.
func TraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}

// WithTraceID returns a new context with the given trace ID.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

// RequestIDMiddleware injects a trace ID into the request context and response header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(TraceIDHeader)
		if id == "" {
			id = uuid.New().String()
		}

		w.Header().Set(TraceIDHeader, id)
		next.ServeHTTP(w, r.WithContext(WithTraceID(r.Context(), id)))
	})
}
