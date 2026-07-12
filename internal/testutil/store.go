package testutil

import (
	"database/sql"
	"testing"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	_ "modernc.org/sqlite" // register sqlite driver for in-memory test stores
)

// Stores pairs an in-memory transaction ledger with an in-memory audit store.
type Stores struct {
	Ledger engine.TransactionStore
	Audit  audit.Store
}

// NewMemoryStores returns in-memory transaction and audit stores.
func NewMemoryStores(t *testing.T) Stores {
	t.Helper()
	return Stores{
		Ledger: engine.NewMemoryStore(),
		Audit:  audit.NewMemoryStore(),
	}
}

// NewSQLiteStores returns transaction and audit stores backed by a shared
// in-memory SQLite connection. Both stores share a single *sql.DB to avoid
// SQLite locking issues in tests.
func NewSQLiteStores(t *testing.T) Stores {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ledger, err := engine.NewSQLiteStoreFromDB(db)
	if err != nil {
		t.Fatalf("init transaction store: %v", err)
	}
	auditStore, err := audit.NewSQLiteStoreFromDB(db)
	if err != nil {
		t.Fatalf("init audit store: %v", err)
	}
	return Stores{Ledger: ledger, Audit: auditStore}
}
