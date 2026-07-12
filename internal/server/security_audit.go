package server

import (
	"context"
	"net/http"

	"github.com/openmuara/openmuara/internal/audit"
)

// SecurityEvent types.
const (
	SecurityEventAuthFailure  = "security.auth.failure"
	SecurityEventAuthSuccess  = "security.auth.success"
	SecurityEventConfigChange = "security.config.change"
	SecurityEventRateLimit    = "security.rate_limit.triggered"
	SecurityEventTLSEnabled   = "security.tls.enabled"
	SecurityEventReplayAction = "security.replay.action"
)

// logSecurityEvent attempts to log a security event to the audit store in the
// request context. It silently ignores missing loggers so middleware stays
// composable.
func logSecurityEvent(ctx context.Context, action, detail string) {
	logger := audit.FromContext(ctx)
	if logger == nil {
		return
	}
	logger.Log(ctx, action, "security", "", detail, "ok")
}

// logSecurityEventFromRequest extracts the logger from the request context.
func logSecurityEventFromRequest(r *http.Request, event, detail string) {
	logSecurityEvent(r.Context(), event, detail)
}
