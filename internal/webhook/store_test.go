package webhook

import (
	"testing"
	"time"
)

func TestMemoryStoreSaveAndGet(t *testing.T) {
	store := NewMemoryStore()
	attempt := &Attempt{
		ID:      "id-1",
		Ref:     "ref-1",
		URL:     "http://localhost/webhook",
		Status:  AttemptStatusPending,
		Payload: []byte(`{}`),
	}

	if err := store.Save(attempt); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := store.Get("ref-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected attempt")
	}
	if got.ID != "id-1" {
		t.Errorf("id: want id-1, got %q", got.ID)
	}
}

func TestMemoryStoreGetMissing(t *testing.T) {
	store := NewMemoryStore()
	got, err := store.Get("missing")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != nil {
		t.Error("expected nil for missing ref")
	}
}

func TestMemoryStoreListOrder(t *testing.T) {
	store := NewMemoryStore()
	for i := 0; i < 5; i++ {
		attempt := &Attempt{
			ID:      "id",
			Ref:     "ref-" + string(rune('a'+i)),
			URL:     "http://localhost/webhook",
			Status:  AttemptStatusPending,
			Payload: []byte(`{}`),
		}
		// Stagger created times so ordering is deterministic.
		attempt.CreatedAt = time.Now().Add(time.Duration(i) * time.Second)
		if err := store.Save(attempt); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	got, err := store.List(3, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 attempts, got %d", len(got))
	}
	// Most recent first.
	if got[0].Ref != "ref-e" {
		t.Errorf("most recent: want ref-e, got %q", got[0].Ref)
	}
}
