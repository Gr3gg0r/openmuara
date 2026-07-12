package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
)

func TestCreateCollectionSuccess(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.Routes()[0].Handler

	body, _ := json.Marshal(map[string]string{"title": "Test Collection"})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.CollectionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Collection.ID == "" {
		t.Error("collection id is empty")
	}
	if resp.Collection.Title != "Test Collection" {
		t.Errorf("title: want %q, got %q", "Test Collection", resp.Collection.Title)
	}
	if resp.Collection.Status != "active" {
		t.Errorf("status: want active, got %q", resp.Collection.Status)
	}
}

func TestCreateCollectionMissingTitle(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.Routes()[0].Handler

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestCreateCollectionInvalidAPIKey(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.Routes()[0].Handler

	body, _ := json.Marshal(map[string]string{"title": "Test Collection"})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("wrong-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d", rec.Code)
	}
}

func TestGetCollectionSuccess(t *testing.T) {
	p := newInitializedProvider(t)
	collection := createCollection(t, p)

	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}")
	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/"+collection.ID, nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.CollectionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Collection.ID != collection.ID {
		t.Errorf("id: want %q, got %q", collection.ID, resp.Collection.ID)
	}
}

func TestGetCollectionNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/missing", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}
