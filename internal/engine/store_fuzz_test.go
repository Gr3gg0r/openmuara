package engine

import "testing"

func FuzzMemoryStoreIdempotency(f *testing.F) {
	seeds := []string{"idem-1", "", "key-with-unicode-✓", "a"}
	for _, key := range seeds {
		f.Add(key)
	}

	f.Fuzz(func(t *testing.T, key string) {
		store := NewMemoryStore()
		first, created, err := store.CreateOrGet(Transaction{IdempotencyKey: key, Reference: "ref-1", Amount: 10.0})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !created {
			t.Fatalf("first CreateOrGet should create")
		}

		second, created2, err := store.CreateOrGet(Transaction{IdempotencyKey: key, Reference: "ref-2", Amount: 20.0})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if key == "" {
			if !created2 {
				t.Fatalf("empty idempotency key should create a new transaction")
			}
			return
		}

		if created2 {
			t.Fatalf("duplicate idempotency key should not create a new transaction")
		}
		if second.ID != first.ID {
			t.Fatalf("duplicate idempotency key should return the same transaction: got %q and %q", first.ID, second.ID)
		}
	})
}
