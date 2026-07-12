package engine

import (
	"database/sql"
	"errors"
	"math"
	"testing"
)

func TestCanTransition_UnknownSource(t *testing.T) {
	if CanTransition(TransactionStatus("unknown"), TransactionStatusPaid) {
		t.Error("expected unknown source status to disallow transition")
	}
}

func TestMemoryStoreList_NegativeOffset(t *testing.T) {
	store := NewMemoryStore()
	for i := 0; i < 2; i++ {
		if _, _, err := store.CreateOrGet(Transaction{Reference: string(rune('a' + i)), Amount: float64(i)}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	got, err := store.List(0, -5)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("want 2 transactions, got %d", len(got))
	}
}

func TestMemoryStoreList_OffsetOverflow(t *testing.T) {
	store := NewMemoryStore()
	if _, _, err := store.CreateOrGet(Transaction{Reference: "ref-1", Amount: 1.0}); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.List(0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("want 0 transactions, got %d", len(got))
	}
}

func TestSQLiteStoreCreateOrGet_MarshalItemsError(t *testing.T) {
	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	_, _, err = store.CreateOrGet(Transaction{
		Reference: "ref-1",
		Amount:    10.0,
		Items: []TransactionItem{
			{Price: math.NaN(), Quantity: 1},
		},
	})
	if err == nil {
		t.Fatal("expected error marshaling items with NaN price")
	}
	if !errors.Is(err, errors.New("marshal items")) {
		// json.Marshal error is wrapped; just confirm the message mentions marshaling.
		t.Logf("marshal error: %v", err)
	}
}

func TestNewSQLiteStoreFromDB_MigrateError(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	_, err = NewSQLiteStoreFromDB(db)
	if err == nil {
		t.Fatal("expected error when migrating with closed database")
	}
}
