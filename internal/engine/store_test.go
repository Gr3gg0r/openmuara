package engine

import (
	"strconv"
	"sync"
	"testing"
)

func TestMemoryStoreCreateOrGet(t *testing.T) {
	store := NewMemoryStore()
	tx := Transaction{Reference: "ref-1", Amount: 10.0}

	saved, created, err := store.CreateOrGet(tx)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !created {
		t.Error("expected created to be true")
	}
	if saved.ID == "" {
		t.Error("expected ID to be generated")
	}

	loaded, ok, err := store.GetByID(saved.ID)

	if err != nil {

		t.Fatalf("lookup transaction: %v", err)

	}
	if !ok {
		t.Fatal("expected to find saved transaction")
	}
	if loaded.Reference != "ref-1" {
		t.Errorf("reference: want ref-1, got %q", loaded.Reference)
	}
}

func TestMemoryStoreIdempotency(t *testing.T) {
	store := NewMemoryStore()
	tx := Transaction{IdempotencyKey: "idem-1", Reference: "ref-1", Amount: 10.0}

	first, created, _ := store.CreateOrGet(tx)
	if !created {
		t.Error("expected first creation")
	}

	second, created, _ := store.CreateOrGet(Transaction{IdempotencyKey: "idem-1", Reference: "ref-2", Amount: 20.0})
	if created {
		t.Error("expected second call to be a cache hit")
	}
	if second.ID != first.ID {
		t.Errorf("idempotency mismatch: want %q, got %q", first.ID, second.ID)
	}
	if second.Reference != first.Reference {
		t.Errorf("idempotency should return original reference: want %q, got %q", first.Reference, second.Reference)
	}
}

func TestMemoryStoreEmptyIdempotencyKeyCreatesNew(t *testing.T) {
	store := NewMemoryStore()

	first, created, _ := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0})
	if !created {
		t.Error("expected first creation")
	}

	second, created, _ := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0})
	if !created {
		t.Error("expected second creation because idempotency key is empty")
	}
	if second.ID == first.ID {
		t.Error("expected different IDs for empty idempotency keys")
	}
}

func TestMemoryStoreGetByReference(t *testing.T) {
	store := NewMemoryStore()
	tx := Transaction{Reference: "ref-1", Amount: 10.0}

	saved, _, _ := store.CreateOrGet(tx)

	loaded, ok, err := store.GetByReference("ref-1")

	if err != nil {

		t.Fatalf("lookup transaction: %v", err)

	}
	if !ok {
		t.Fatal("expected to find transaction by reference")
	}
	if loaded.ID != saved.ID {
		t.Errorf("ID mismatch: want %q, got %q", saved.ID, loaded.ID)
	}

	_, ok, err = store.GetByReference("missing")
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if ok {
		t.Error("expected false for missing reference")
	}
}

func TestMemoryStoreList(t *testing.T) {
	store := NewMemoryStore()

	for i := 1; i <= 5; i++ {
		ref := string(rune('a' + i - 1))
		if _, _, err := store.CreateOrGet(Transaction{Reference: ref, Amount: float64(i)}); err != nil {
			t.Fatalf("create %s: %v", ref, err)
		}
	}

	all, err := store.List(10, 0)

	if err != nil {

		t.Fatalf("list transactions: %v", err)

	}
	if len(all) != 5 {
		t.Fatalf("expected 5 transactions, got %d", len(all))
	}

	// Most recent first.
	if all[0].Reference != "e" {
		t.Errorf("expected most recent reference e, got %q", all[0].Reference)
	}
	if all[4].Reference != "a" {
		t.Errorf("expected oldest reference a, got %q", all[4].Reference)
	}

	limited, err := store.List(2, 0)

	if err != nil {

		t.Fatalf("list transactions: %v", err)

	}
	if len(limited) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(limited))
	}
	if limited[0].Reference != "e" || limited[1].Reference != "d" {
		t.Errorf("expected [e, d], got %v", limited)
	}

	zero, err := store.List(0, 0)

	if err != nil {

		t.Fatalf("list transactions: %v", err)

	}
	if len(zero) != 5 {
		t.Errorf("expected all 5 transactions when limit <= 0, got %d", len(zero))
	}
}

func TestMemoryStoreConcurrentCreateOrGet(t *testing.T) {
	store := NewMemoryStore()
	const workers = 50

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _, _ = store.CreateOrGet(Transaction{IdempotencyKey: "shared-key", Reference: "ref-1"})
		}()
	}
	wg.Wait()

	all, err := store.List(100, 0)

	if err != nil {

		t.Fatalf("list transactions: %v", err)

	}
	if len(all) != 1 {
		t.Errorf("expected exactly 1 transaction, got %d", len(all))
	}
}

func TestMemoryStoreConcurrentDifferentKeys(t *testing.T) {
	store := NewMemoryStore()
	const workers = 50

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			_, _, _ = store.CreateOrGet(Transaction{IdempotencyKey: "key-" + strconv.Itoa(i), Reference: "ref-1"})
		}(i)
	}
	wg.Wait()

	all, err := store.List(100, 0)

	if err != nil {

		t.Fatalf("list transactions: %v", err)

	}
	if len(all) != workers {
		t.Errorf("expected %d transactions, got %d", workers, len(all))
	}
}

func TestMemoryStoreReturnValueIsolation(t *testing.T) {
	store := NewMemoryStore()
	tx := Transaction{Reference: "ref-1", Amount: 10.0}

	saved, _, _ := store.CreateOrGet(tx)
	loaded, _, err := store.GetByID(saved.ID)
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	loaded.Amount = 999.0

	reloaded, _, err := store.GetByID(saved.ID)
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if reloaded.Amount != 10.0 {
		t.Errorf("mutation leaked into store: want 10.0, got %f", reloaded.Amount)
	}
}

func TestMemoryStorePreservesTraceIDOnReferenceReuse(t *testing.T) {
	store := NewMemoryStore()
	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0, TraceID: "trace-first"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	saved, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 20.0})
	if err != nil {
		t.Fatalf("reuse: %v", err)
	}
	if saved.TraceID != "trace-first" {
		t.Errorf("trace_id: want trace-first, got %q", saved.TraceID)
	}
}

func TestMemoryStoreKeepsExplicitTraceID(t *testing.T) {
	store := NewMemoryStore()
	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 10.0, TraceID: "trace-first"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	saved, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 20.0, TraceID: "trace-second"})
	if err != nil {
		t.Fatalf("reuse: %v", err)
	}
	if saved.TraceID != "trace-second" {
		t.Errorf("trace_id: want trace-second, got %q", saved.TraceID)
	}
}
