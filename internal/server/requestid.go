package server

import (
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// Re-export request ID middleware from httputil.

// TraceIDHeader is the HTTP header used to propagate trace IDs.
const TraceIDHeader = httputil.TraceIDHeader

// RequestIDMiddleware injects a trace ID into the request context and response header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return httputil.RequestIDMiddleware(next)
}
