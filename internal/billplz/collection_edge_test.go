package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateCollectionMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodPost, "/api/v3/collections")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestCreateCollectionMissingAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/collections")

	body, _ := json.Marshal(map[string]string{"title": "Test Collection"})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestCreateCollectionInvalidJSON(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/collections")

	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestGetCollectionMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodGet, "/api/v3/collections/{id}")

	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections/123", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestGetCollectionMissingAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestPaymentMethodsMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodGet, "/api/v3/collections/{id}/payment_methods")

	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections/123/payment_methods", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestPaymentMethodsMissingAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}/payment_methods")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/123/payment_methods", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}
