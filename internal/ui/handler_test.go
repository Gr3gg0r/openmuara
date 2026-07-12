package ui

import (
	"net/http/httptest"
	"testing"
)

func TestServeEscapePage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeEscapePage(rec, EscapePageData{Ref: "r", ReturnURL: "http://localhost", Amount: "10.00"}); err != nil {
		t.Fatalf("serve escape page: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("content-type: want html, got %q", ct)
	}
}

func TestServeDashboard(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeDashboard(rec, DashboardData{ActiveProvider: "fawry"}); err != nil {
		t.Fatalf("serve dashboard: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}
