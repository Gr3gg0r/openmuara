package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddlewareInjectsHeader(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	h := rec.Header().Get("X-Trace-Id")
	if h == "" {
		t.Fatal("X-Trace-Id header missing")
	}
}

func TestRequestIDMiddlewarePreservesExistingID(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Trace-Id", "existing-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Trace-Id"); got != "existing-id" {
		t.Fatalf("X-Trace-Id: want existing-id, got %q", got)
	}
}
