package httputil

import "context"

type csrfContextKey struct{}

// WithCSRFToken returns a context carrying the CSRF token for the request.
func WithCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, csrfContextKey{}, token)
}

// CSRFTokenFromContext retrieves the CSRF token injected by the CSRF guard middleware.
func CSRFTokenFromContext(ctx context.Context) (string, bool) {
	tok, ok := ctx.Value(csrfContextKey{}).(string)
	return tok, ok
}
