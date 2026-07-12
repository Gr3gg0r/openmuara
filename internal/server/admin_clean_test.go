package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestCleanHandlerRequiresAdmin(t *testing.T) {
	store := engine.NewMemoryStore()
	cfg := RouterConfig{TransactionStore: store}

	req := newViewerRequest(http.MethodPost, "/_admin/clean")
	rec := httptest.NewRecorder()

	cleanHandler(cfg)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d", rec.Code)
	}
}

func TestCleanHandlerClearsData(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	auditStore := audit.NewMemoryStore()
	_ = auditStore.Save(audit.Event{Action: "test"})

	dispatcher := webhook.NewDispatcherFromBuilder("http://localhost/hook", 3,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)
	dispatcher.Store = webhook.NewMemoryStore()

	cfg := RouterConfig{
		TransactionStore: store,
		AuditStore:       auditStore,
		Dispatcher:       dispatcher,
	}

	req := newAdminRequest(http.MethodPost, "/_admin/clean")
	rec := httptest.NewRecorder()

	cleanHandler(cfg)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d: %s", rec.Code, rec.Body.String())
	}

	txs, err := store.List(-1, 0)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if len(txs) != 0 {
		t.Errorf("want 0 transactions after clean, got %d", len(txs))
	}

	events, err := auditStore.List(-1, 0)
	if err != nil {
		t.Fatalf("list audit events: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("want 0 audit events after clean, got %d", len(events))
	}

	attempts, err := dispatcher.Store.List(-1, 0)
	if err != nil {
		t.Fatalf("list webhook attempts: %v", err)
	}
	if len(attempts) != 0 {
		t.Errorf("want 0 webhook attempts after clean, got %d", len(attempts))
	}
}
