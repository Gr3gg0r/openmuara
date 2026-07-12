package webhook

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTestReceiverAcceptsWebhook(t *testing.T) {
	handler := NewTestReceiverHandler()
	body := []byte(`{"orderStatus":"PAID"}`)
	req := httptest.NewRequest(http.MethodPost, "/_admin/webhook-receiver", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestTestReceiverRejectsGet(t *testing.T) {
	handler := NewTestReceiverHandler()
	req := httptest.NewRequest(http.MethodGet, "/_admin/webhook-receiver", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}
