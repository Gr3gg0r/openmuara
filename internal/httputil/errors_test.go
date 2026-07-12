package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondErrorWritesJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	RespondError(rec, req, ErrNotFound, http.StatusNotFound, "not found")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("content-type: want application/json, got %q", ct)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error.Code != ErrNotFound {
		t.Errorf("code: want %q, got %q", ErrNotFound, resp.Error.Code)
	}
	if resp.Error.Message != "not found" {
		t.Errorf("message: want not found, got %q", resp.Error.Message)
	}
}
