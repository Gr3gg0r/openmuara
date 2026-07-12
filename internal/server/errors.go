// Package server provides the HTTP server, router, and middleware for OpenMuara.
package server

import (
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// Re-export error types and helpers from httputil to avoid import cycles.

// ErrorCode is a stable machine-readable error identifier.
type ErrorCode = httputil.ErrorCode

// ErrorResponse is the JSON envelope for HTTP errors.
type ErrorResponse = httputil.ErrorResponse

// ErrorDetail contains the specific error information.
type ErrorDetail = httputil.ErrorDetail

const (
	// ErrInvalidSignature indicates a rejected request signature.
	ErrInvalidSignature = httputil.ErrInvalidSignature
	// ErrMissingField indicates a required field was missing.
	ErrMissingField = httputil.ErrMissingField
	// ErrInvalidJSON indicates the request body could not be parsed.
	ErrInvalidJSON = httputil.ErrInvalidJSON
	// ErrUnauthorized indicates the caller is not authorized.
	ErrUnauthorized = httputil.ErrUnauthorized
	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound = httputil.ErrNotFound
	// ErrInternal indicates an unexpected internal error.
	ErrInternal = httputil.ErrInternal
)

// RespondError writes a structured error response.
func RespondError(w http.ResponseWriter, r *http.Request, code ErrorCode, status int, message string) {
	httputil.RespondError(w, r, code, status, message)
}
