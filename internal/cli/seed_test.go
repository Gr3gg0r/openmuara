package cli

import (
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestSeedDevDataCreatesTransactionsAndWebhooks(t *testing.T) {
	ledger := engine.NewMemoryStore()
	store := webhook.NewMemoryStore()

	if err := seedDevData(ledger, store, []string{"stripe"}); err != nil {
		t.Fatalf("seedDevData: %v", err)
	}

	txs, err := ledger.List(-1, 0)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if len(txs) != 3 {
		t.Fatalf("expected 3 seeded transactions, got %d", len(txs))
	}

	attempts, err := store.List(-1, 0)
	if err != nil {
		t.Fatalf("list attempts: %v", err)
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 seeded webhook attempts, got %d", len(attempts))
	}
}

func TestSeedDevDataPicksEnabledProvider(t *testing.T) {
	ledger := engine.NewMemoryStore()

	if err := seedDevData(ledger, nil, []string{"fawry"}); err != nil {
		t.Fatalf("seedDevData: %v", err)
	}

	txs, err := ledger.List(-1, 0)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if len(txs) == 0 {
		t.Fatal("expected at least one seeded transaction")
	}
	if txs[0].Provider != "fawry" {
		t.Fatalf("expected provider fawry, got %s", txs[0].Provider)
	}
}

func TestSeedDevDataNilStore(t *testing.T) {
	ledger := engine.NewMemoryStore()

	if err := seedDevData(ledger, nil, []string{"stripe"}); err != nil {
		t.Fatalf("seedDevData with nil store: %v", err)
	}
}
