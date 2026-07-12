// Package httputil provides shared HTTP helpers used by both server and provider packages.
package httputil

import (
	"encoding/json"
	"net/http"
)

// ErrorCode is a stable machine-readable error identifier.
type ErrorCode string

// Error codes returned by OpenMuara HTTP handlers.
const (
	ErrInvalidSignature ErrorCode = "OPENMUARA_INVALID_SIGNATURE"
	ErrMissingField     ErrorCode = "OPENMUARA_MISSING_FIELD"
	ErrInvalidJSON      ErrorCode = "OPENMUARA_INVALID_JSON"
	ErrUnauthorized     ErrorCode = "OPENMUARA_UNAUTHORIZED"
	ErrNotFound         ErrorCode = "OPENMUARA_NOT_FOUND"
	ErrInvalidState     ErrorCode = "OPENMUARA_INVALID_STATE"
	ErrInternal         ErrorCode = "OPENMUARA_INTERNAL_ERROR"
)

// ErrorResponse is the JSON envelope for HTTP errors.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the specific error information.
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	TraceID string    `json:"trace_id"`
}

// RespondError writes a structured error response with the trace ID from r.Context().
func RespondError(w http.ResponseWriter, r *http.Request, code ErrorCode, status int, message string) {
	traceID := TraceIDFromContext(r.Context())
	if traceID == "" {
		traceID = r.Header.Get(TraceIDHeader)
	}

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			TraceID: traceID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
