package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/httputil"
)

func TestRespondError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(httputil.WithTraceID(req.Context(), "trace-123"))
	rec := httptest.NewRecorder()

	RespondError(rec, req, ErrInvalidSignature, http.StatusBadRequest, "signature mismatch")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	if body.Error.Code != ErrInvalidSignature {
		t.Errorf("code: want %q, got %q", ErrInvalidSignature, body.Error.Code)
	}
	if body.Error.TraceID != "trace-123" {
		t.Errorf("trace_id: want trace-123, got %q", body.Error.TraceID)
	}
}
