package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
)

func TestScenarioHandlerSuccess(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "test", Type: "charge", Reference: "ref-scenario-1", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := newScenarioHandler(store)
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/success?ref=ref-scenario-1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := store.GetByReference("ref-scenario-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction")
	}
	if tx.Status != engine.TransactionStatusPaid {
		t.Errorf("status: want paid, got %q", tx.Status)
	}
}

func TestScenarioHandlerFail(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "test", Type: "charge", Reference: "ref-scenario-2", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := newScenarioHandler(store)
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/fail?ref=ref-scenario-2")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := store.GetByReference("ref-scenario-2")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction")
	}
	if tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("status: want unpaid, got %q", tx.Status)
	}
}

func TestScenarioHandlerMissingRef(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := newScenarioHandler(store)
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/success")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestScenarioHandlerUnknownOutcome(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := newScenarioHandler(store)
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/unknown?ref=x")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestScenarioHandlerDeniesViewer(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := newScenarioHandler(store)
	req := newViewerRequest(http.MethodPost, "/_admin/scenario/success?ref=x")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: want 403, got %d", rec.Code)
	}
}
