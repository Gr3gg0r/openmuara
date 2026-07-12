package engine

import (
	"sync"
	"testing"
)

func TestMemoryStoreStressCreateOrGet(_ *testing.T) {
	store := NewMemoryStore()
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			ref := string(rune('a' + i%26))
			_, _, _ = store.CreateOrGet(Transaction{Reference: ref, Amount: float64(i), Status: TransactionStatusNew})
			_, _, _ = store.GetByReference(ref)
		}()
	}

	wg.Wait()
}

func TestMemoryStoreStressTransition(t *testing.T) {
	store := NewMemoryStore()
	ref := "shared-ref"
	if _, _, err := store.CreateOrGet(Transaction{Reference: ref, Amount: 10.0, Status: TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			tx, ok, err := store.GetByReference(ref)
			if err != nil || !ok {
				return
			}
			target := TransactionStatusPaid
			if i%2 == 0 {
				target = TransactionStatusUnpaid
			}
			if err := Transition(&tx, target); err != nil {
				return
			}
			_, _, _ = store.CreateOrGet(tx)
		}()
	}

	wg.Wait()
}
