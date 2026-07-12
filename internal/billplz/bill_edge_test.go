package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
)

func TestCreateBillMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodPost, "/api/v3/bills")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/bills", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestCreateBillInvalidJSON(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")

	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestCreateBillMissingBasicAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")

	body, _ := json.Marshal(map[string]any{
		"collection_id": "missing",
		"email":         "test@example.com",
		"name":          "Test User",
		"amount":        1000,
		"callback_url":  "http://localhost:9999/callback",
		"description":   "Test bill",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestCreateBillValidationErrors(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")

	cases := []struct {
		name string
		body map[string]any
	}{
		{"missing email", map[string]any{
			"collection_id": c.ID, "name": "Test User", "amount": 1000,
			"callback_url": "http://localhost:9999/callback", "description": "Test bill",
		}},
		{"missing name", map[string]any{
			"collection_id": c.ID, "email": "test@example.com", "amount": 1000,
			"callback_url": "http://localhost:9999/callback", "description": "Test bill",
		}},
		{"invalid amount", map[string]any{
			"collection_id": c.ID, "email": "test@example.com", "name": "Test User", "amount": 0,
			"callback_url": "http://localhost:9999/callback", "description": "Test bill",
		}},
		{"missing callback_url", map[string]any{
			"collection_id": c.ID, "email": "test@example.com", "name": "Test User", "amount": 1000,
			"description": "Test bill",
		}},
		{"missing description", map[string]any{
			"collection_id": c.ID, "email": "test@example.com", "name": "Test User", "amount": 1000,
			"callback_url": "http://localhost:9999/callback",
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth("muara-billplz-api-key", "")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestCreateBillValidationMissingCollectionID(t *testing.T) {
	p := billplz.NewProvider()
	if err := p.Init(map[string]any{
		"api_key":         "muara-billplz-api-key",
		"x_signature_key": "muara-billplz-xsig-key",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetBaseURL("http://localhost:9000")

	body, _ := json.Marshal(map[string]any{
		"email": "test@example.com", "name": "Test User", "amount": 1000,
		"callback_url": "http://localhost:9999/callback", "description": "Test bill",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	mux := testMux(t, p)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetBillMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodGet, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills/123", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestGetBillMissingAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/bills/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestGetBillNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/bills/missing", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestDeleteBillMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodDelete, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/bills/123", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestDeleteBillMissingAuth(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodDelete, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodDelete, "/api/v3/bills/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestDeleteBillNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodDelete, "/api/v3/bills/{id}")

	req := httptest.NewRequest(http.MethodDelete, "/api/v3/bills/missing", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}
