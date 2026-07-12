package engine

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTransactionMarshalsAllFields(t *testing.T) {
	tx := Transaction{
		ID:             "tx-1",
		Provider:       "fawry",
		Type:           "charge",
		Amount:         100.00,
		Currency:       "EGP",
		Status:         TransactionStatusNew,
		CustomerRef:    "cust-1",
		IdempotencyKey: "idem-1",
		Reference:      "ref-1",
		Items: []TransactionItem{
			{ItemCode: "prod-1", Price: 99.99, Quantity: 1},
		},
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	want := []string{"id", "provider", "type", "amount", "currency", "status", "customerRef", "idempotencyKey", "reference", "items", "createdAt", "updatedAt"}
	for _, key := range want {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing field %q in JSON", key)
		}
	}
}

func TestTransactionOmitsEmptyIdempotencyKeyAndItems(t *testing.T) {
	tx := Transaction{
		ID:        "tx-1",
		Provider:  "fawry",
		Reference: "ref-1",
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	if _, ok := raw["idempotencyKey"]; ok {
		t.Error("expected empty idempotencyKey to be omitted")
	}
	if _, ok := raw["items"]; ok {
		t.Error("expected empty items to be omitted")
	}
}

func TestNewTransactionSetsTimestamps(t *testing.T) {
	tx := NewTransaction(Transaction{ID: "tx-1"})

	if tx.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if tx.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}
}

func TestNewTransactionPreservesExistingTimestamps(t *testing.T) {
	created := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	tx := NewTransaction(Transaction{ID: "tx-1", CreatedAt: created, UpdatedAt: updated})

	if !tx.CreatedAt.Equal(created) {
		t.Errorf("CreatedAt changed: want %v, got %v", created, tx.CreatedAt)
	}
	if !tx.UpdatedAt.Equal(updated) {
		t.Errorf("UpdatedAt changed: want %v, got %v", updated, tx.UpdatedAt)
	}
}

func TestCanTransition(t *testing.T) {
	cases := []struct {
		from TransactionStatus
		to   TransactionStatus
		want bool
	}{
		{TransactionStatusNew, TransactionStatusPaid, true},
		{TransactionStatusNew, TransactionStatusUnpaid, true},
		{TransactionStatusNew, TransactionStatusRefunded, false},
		{TransactionStatusPaid, TransactionStatusRefunded, true},
		{TransactionStatusPaid, TransactionStatusUnpaid, false},
		{TransactionStatusPaid, TransactionStatusNew, false},
		{TransactionStatusUnpaid, TransactionStatusPaid, false},
		{TransactionStatusUnpaid, TransactionStatusNew, false},
		{TransactionStatusRefunded, TransactionStatusPaid, false},
		{TransactionStatusRefunded, TransactionStatusNew, false},
	}

	for _, tc := range cases {
		got := CanTransition(tc.from, tc.to)
		if got != tc.want {
			t.Errorf("CanTransition(%q, %q) = %v, want %v", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestTransitionUpdatesStatusAndTimestamp(t *testing.T) {
	tx := NewTransaction(Transaction{ID: "tx-1", Status: TransactionStatusNew})
	before := tx.UpdatedAt

	if err := Transition(&tx, TransactionStatusPaid); err != nil {
		t.Fatalf("transition: %v", err)
	}
	if tx.Status != TransactionStatusPaid {
		t.Errorf("status: want %q, got %q", TransactionStatusPaid, tx.Status)
	}
	if !tx.UpdatedAt.After(before) {
		t.Error("UpdatedAt was not refreshed")
	}
}

func TestTransitionNoOpWhenSameStatus(t *testing.T) {
	tx := NewTransaction(Transaction{ID: "tx-1", Status: TransactionStatusPaid})
	before := tx.UpdatedAt

	if err := Transition(&tx, TransactionStatusPaid); err != nil {
		t.Fatalf("transition: %v", err)
	}
	if !tx.UpdatedAt.Equal(before) {
		t.Error("UpdatedAt should not change for no-op transition")
	}
}

func TestTransitionRejectsInvalidTransition(t *testing.T) {
	tx := NewTransaction(Transaction{ID: "tx-1", Status: TransactionStatusNew})

	if err := Transition(&tx, TransactionStatusRefunded); err == nil {
		t.Fatal("expected error for invalid transition")
	}
	if tx.Status != TransactionStatusNew {
		t.Errorf("status should remain unchanged: want %q, got %q", TransactionStatusNew, tx.Status)
	}
}

func TestTransactionStatusRoundTrips(t *testing.T) {
	original := Transaction{ID: "tx-1", Status: TransactionStatusPaid}
	data, _ := json.Marshal(original)

	var parsed Transaction
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if parsed.Status != TransactionStatusPaid {
		t.Errorf("status: want %q, got %q", TransactionStatusPaid, parsed.Status)
	}
}
