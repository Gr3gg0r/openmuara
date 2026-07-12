package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPIHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	OpenAPIHandler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/yaml" {
		t.Errorf("content-type: want application/yaml, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "openapi: 3.0.3") {
		t.Error("response missing openapi version")
	}

	var doc map[string]any
	if err := yaml.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid yaml: %v", err)
	}
	if doc["openapi"] != "3.0.3" {
		t.Errorf("unexpected openapi version: %v", doc["openapi"])
	}
}

func TestOpenAPISpecMatchesDocs(t *testing.T) {
	embedded, err := os.ReadFile("openapi.yaml")
	if err != nil {
		t.Fatalf("read embedded openapi.yaml: %v", err)
	}
	docs, err := os.ReadFile("../../docs/openapi.yaml")
	if err != nil {
		t.Fatalf("read docs/openapi.yaml: %v", err)
	}
	if string(embedded) != string(docs) {
		t.Error("internal/server/openapi.yaml and docs/openapi.yaml are out of sync")
	}
}
