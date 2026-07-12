package httputil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTraceIDFromContextMissing(t *testing.T) {
	if got := TraceIDFromContext(context.Background()); got != "" {
		t.Errorf("expected empty trace id, got %q", got)
	}
}

func TestTraceIDFromContextPresent(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")
	if got := TraceIDFromContext(ctx); got != "trace-123" {
		t.Errorf("expected trace id trace-123, got %q", got)
	}
}

func TestRequestIDMiddlewareGeneratesID(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if got := rec.Header().Get(TraceIDHeader); got == "" {
		t.Fatal("expected trace id header to be set")
	}
}

func TestRequestIDMiddlewarePreservesHeader(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(TraceIDHeader, "existing")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(TraceIDHeader); got != "existing" {
		t.Errorf("trace id header: want existing, got %q", got)
	}
}
