package engine

import (
	"path/filepath"
	"testing"
)

func TestSQLiteStoreCreateOrGet(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	tx := Transaction{Provider: "fawry", Type: "charge", Reference: "ref-1", Amount: 99.99, Currency: "EGP", Status: TransactionStatusNew}
	saved, created, err := store.CreateOrGet(tx)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !created {
		t.Fatal("expected created to be true")
	}
	if saved.ID == "" {
		t.Fatal("expected id to be set")
	}

	loaded, ok, err := store.GetByID(saved.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if !ok {
		t.Fatal("expected to find transaction by id")
	}
	if loaded.Reference != "ref-1" {
		t.Errorf("reference: want ref-1, got %q", loaded.Reference)
	}
}

func TestSQLiteStoreIdempotency(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	first, created, err := store.CreateOrGet(Transaction{IdempotencyKey: "idem-1", Reference: "ref-1", Amount: 10.0})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if !created {
		t.Fatal("expected first to be created")
	}

	second, created, err := store.CreateOrGet(Transaction{IdempotencyKey: "idem-1", Reference: "ref-2", Amount: 20.0})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if created {
		t.Fatal("expected second to be a cache hit")
	}
	if second.Reference != first.Reference {
		t.Errorf("idempotency should preserve original reference: want %q, got %q", first.Reference, second.Reference)
	}
}

func TestSQLiteStoreList(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	for i := 1; i <= 3; i++ {
		ref := string(rune('a' + i - 1))
		if _, _, err := store.CreateOrGet(Transaction{Reference: ref, Amount: float64(i)}); err != nil {
			t.Fatalf("create %s: %v", ref, err)
		}
	}

	all, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 transactions, got %d", len(all))
	}
	if all[0].Reference != "c" {
		t.Errorf("expected most recent reference c, got %q", all[0].Reference)
	}
}

func TestSQLiteStorePersistsAcrossReopen(t *testing.T) {
	db := t.TempDir()
	path := filepath.Join(db, "ledger.db")

	store1, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	if _, _, err := store1.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := store1.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	store2, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("reopen sqlite store: %v", err)
	}
	defer func() { _ = store2.Close() }()

	loaded, ok, err := store2.GetByReference("ref-1")
	if err != nil {
		t.Fatalf("get by reference: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction to persist across reopen")
	}
	if loaded.Amount != 10.0 {
		t.Errorf("amount: want 10.0, got %f", loaded.Amount)
	}
}

func TestSQLiteStoreUpdatesExistingReference(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0, Status: TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0, Status: TransactionStatusPaid}); err != nil {
		t.Fatalf("update: %v", err)
	}

	loaded, ok, err := store.GetByReference("ref-1")
	if err != nil {
		t.Fatalf("get by reference: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction")
	}
	if loaded.Status != TransactionStatusPaid {
		t.Errorf("status: want paid, got %q", loaded.Status)
	}
}

func TestSQLiteStoreListOffset(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	for i := 1; i <= 3; i++ {
		ref := string(rune('a' + i - 1))
		if _, _, err := store.CreateOrGet(Transaction{Reference: ref, Amount: float64(i)}); err != nil {
			t.Fatalf("create %s: %v", ref, err)
		}
	}

	page, err := store.List(1, 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(page) != 1 || page[0].Reference != "b" {
		t.Errorf("expected offset 1 to return b, got %+v", page)
	}
}

func TestSQLiteStoreInvalidPath(t *testing.T) {
	_, err := NewSQLiteStore(filepath.Join(t.TempDir(), "does", "not", "exist", "ledger.db"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestSQLiteStoreScanCorruptItems(t *testing.T) {
	db := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(db, "ledger.db"))
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0}); err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := store.db.Exec("UPDATE transactions SET items = ? WHERE reference = ?", "not json", "ref-1"); err != nil {
		t.Fatalf("corrupt items: %v", err)
	}

	_, _, err = store.GetByReference("ref-1")
	if err == nil {
		t.Fatal("expected error scanning corrupt items")
	}

	_, err = store.List(10, 0)
	if err == nil {
		t.Fatal("expected error listing corrupt items")
	}
}
