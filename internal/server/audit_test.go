package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/provider"
)

func TestAuditAdminHandlerListsEvents(t *testing.T) {
	store := audit.NewMemoryStore()
	_ = store.Save(audit.Event{Actor: "test", Action: "charge.created", ResourceType: "transaction", ResourceID: "r1"})

	mux := http.NewServeMux()
	AuditAdminHandlers(mux, store)

	req := httptest.NewRequest(http.MethodGet, "/_admin/audit", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, ok := body["results"].([]any)
	if !ok || len(results) != 1 {
		t.Fatalf("want 1 event, got %+v", body)
	}
}

func TestAuditMiddlewareInjectsLogger(t *testing.T) {
	store := audit.NewMemoryStore()
	logger := &audit.StoreLogger{Store: store, Actor: "http", Synchronous: true}

	var captured audit.Logger
	handler := AuditMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = audit.FromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if captured == nil {
		t.Fatal("expected audit logger in context")
	}
}

func TestRouterRegistersAuditEndpoint(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		AuditStore:       audit.NewMemoryStore(),
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/audit", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatal("audit endpoint was not registered")
	}
}

var _ = provider.Names
