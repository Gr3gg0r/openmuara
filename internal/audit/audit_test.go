package audit

import (
	"context"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/store/migrations"
)

func TestMemoryStoreSaveAndList(t *testing.T) {
	store := NewMemoryStore()
	ev := Event{
		Actor:        "test",
		Action:       "charge.created",
		ResourceType: "transaction",
		ResourceID:   "r1",
		Result:       "ok",
	}
	if err := store.Save(ev); err != nil {
		t.Fatalf("save: %v", err)
	}

	events, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	if events[0].Action != "charge.created" {
		t.Errorf("action: want charge.created, got %q", events[0].Action)
	}
}

func TestMemoryStoreListSince(t *testing.T) {
	store := NewMemoryStore()
	old := Event{Action: "old", Timestamp: time.Now().Add(-time.Hour)}
	newEv := Event{Action: "new", Timestamp: time.Now()}
	_ = store.Save(old)
	_ = store.Save(newEv)

	events, err := store.ListSince(10, 0, time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("list since: %v", err)
	}
	if len(events) != 1 || events[0].Action != "new" {
		t.Fatalf("want only new event, got %+v", events)
	}
}

func TestContextLogger(t *testing.T) {
	store := NewMemoryStore()
	logger := &StoreLogger{Store: store, Actor: "ctx-test", Synchronous: true}
	ctx := NewContext(context.Background(), logger)

	FromContext(ctx).Log(ctx, "action", "resource", "id", "payload", "ok")

	events, _ := store.List(10, 0)
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	if events[0].Actor != "ctx-test" {
		t.Errorf("actor: want ctx-test, got %q", events[0].Actor)
	}
}

func TestNoopLogger(_ *testing.T) {
	FromContext(context.Background()).Log(context.Background(), "action", "resource", "id", "", "")
}

func TestStoreLoggerAsync(t *testing.T) {
	store := NewMemoryStore()
	logger := &StoreLogger{Store: store, Actor: "async-test"}

	logger.Log(context.Background(), "async.action", "resource", "id", "payload", "ok")

	// Wait for the goroutine to persist the event.
	deadline := time.Now().Add(2 * time.Second)
	for {
		events, _ := store.List(10, 0)
		if len(events) == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("async audit event was not persisted")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStoreLoggerNilStore(_ *testing.T) {
	logger := &StoreLogger{Store: nil, Actor: "nil-test", Synchronous: true}
	logger.Log(context.Background(), "action", "resource", "id", "payload", "ok")
}

func TestJSON(t *testing.T) {
	got := JSON(map[string]string{"key": "value"})
	if got != `{"key":"value"}` {
		t.Errorf("JSON: want compact JSON, got %q", got)
	}

	if JSON(make(chan int)) != "" {
		t.Error("expected empty string for unsupported value")
	}
}

func TestSQLiteStore(t *testing.T) {
	if _, err := migrations.Read("003_audit_logs.sql"); err != nil {
		t.Fatalf("read migration: %v", err)
	}

	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("open sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Save(Event{Action: "test", Actor: "a", ResourceType: "r", ResourceID: "1"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	events, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
}
