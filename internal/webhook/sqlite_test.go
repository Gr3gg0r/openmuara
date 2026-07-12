package webhook

import (
	"testing"
	"time"
)

func TestSQLiteStoreSaveAndGet(t *testing.T) {
	db, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = db.Close() }()

	valid := true
	a := &Attempt{
		Ref:            "ref-1",
		ProviderName:   "stripe",
		URL:            "http://localhost/webhook",
		Status:         AttemptStatusDelivered,
		Attempts:       1,
		Payload:        []byte(`{"id":"evt_1"}`),
		Headers:        map[string]string{"Content-Type": "application/json"},
		SignatureValid: &valid,
		TraceID:        "trace-1",
		History: []AttemptHistory{
			{Time: time.Now(), Status: 200},
		},
	}

	if err := db.Save(a); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := db.Get("ref-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected attempt, got nil")
	}
	if got.Ref != "ref-1" || got.ProviderName != "stripe" {
		t.Fatalf("attempt mismatch: %+v", got)
	}
	if got.SignatureValid == nil || !*got.SignatureValid {
		t.Fatal("expected signature_valid true")
	}
	if len(got.History) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(got.History))
	}
}

func TestSQLiteStoreGetMissing(t *testing.T) {
	db, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = db.Close() }()

	got, err := db.Get("missing")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestSQLiteStoreListOrdering(t *testing.T) {
	db, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = db.Close() }()

	for _, ref := range []string{"ref-1", "ref-2", "ref-3"} {
		a := &Attempt{Ref: ref, URL: "http://localhost/webhook", Status: AttemptStatusPending}
		a.CreatedAt = time.Now().Add(-time.Duration(ref[4]-'0') * time.Hour)
		if err := db.Save(a); err != nil {
			t.Fatalf("save %s: %v", ref, err)
		}
	}

	attempts, err := db.List(2, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 attempts, got %d", len(attempts))
	}
	if attempts[0].Ref != "ref-1" {
		t.Fatalf("expected newest ref-1, got %s", attempts[0].Ref)
	}
}

func TestSQLiteStoreUpdateExisting(t *testing.T) {
	db, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = db.Close() }()

	a := &Attempt{Ref: "ref-1", URL: "http://localhost/webhook", Status: AttemptStatusPending}
	if err := db.Save(a); err != nil {
		t.Fatalf("save: %v", err)
	}

	a.Status = AttemptStatusDelivered
	a.Attempts = 1
	if err := db.Save(a); err != nil {
		t.Fatalf("save update: %v", err)
	}

	got, err := db.Get("ref-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Status != AttemptStatusDelivered || got.Attempts != 1 {
		t.Fatalf("update mismatch: %+v", got)
	}
}
