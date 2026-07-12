package cli

import (
	"fmt"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
	"github.com/google/uuid"
)

// seedDevData inserts a small set of synthetic transactions and webhook attempts
// so developers can explore the dashboard without running a real charge flow.
// It is only called when dev.seed is enabled.
func seedDevData(ledger engine.TransactionStore, store webhook.AttemptStore, enabled []string) error {
	if ledger == nil {
		return nil
	}

	now := time.Now().UTC()
	sampleProvider := pickEnabledProvider(enabled)

	txs := []engine.Transaction{
		{
			ID:          uuid.Must(uuid.NewRandom()).String(),
			Provider:    sampleProvider,
			Type:        "checkout_session",
			Reference:   "dev-tx-paid-001",
			Amount:      49.99,
			Currency:    "USD",
			Status:      engine.TransactionStatusPaid,
			CustomerRef: "dev-customer@example.com",
			TraceID:     "trace-dev-001",
			Items: []engine.TransactionItem{
				{ItemCode: "prod-1", Price: 49.99, Quantity: 1},
			},
			CreatedAt: now.Add(-30 * time.Minute),
			UpdatedAt: now.Add(-28 * time.Minute),
		},
		{
			ID:          uuid.Must(uuid.NewRandom()).String(),
			Provider:    sampleProvider,
			Type:        "checkout_session",
			Reference:   "dev-tx-pending-002",
			Amount:      12.50,
			Currency:    "USD",
			Status:      engine.TransactionStatusNew,
			CustomerRef: "dev-customer@example.com",
			TraceID:     "trace-dev-002",
			Items: []engine.TransactionItem{
				{ItemCode: "prod-2", Price: 12.50, Quantity: 1},
			},
			CreatedAt: now.Add(-10 * time.Minute),
			UpdatedAt: now.Add(-10 * time.Minute),
		},
		{
			ID:          uuid.Must(uuid.NewRandom()).String(),
			Provider:    sampleProvider,
			Type:        "checkout_session",
			Reference:   "dev-tx-failed-003",
			Amount:      99.00,
			Currency:    "USD",
			Status:      engine.TransactionStatusUnpaid,
			CustomerRef: "dev-customer@example.com",
			TraceID:     "trace-dev-003",
			Items: []engine.TransactionItem{
				{ItemCode: "prod-3", Price: 99.00, Quantity: 1},
			},
			CreatedAt: now.Add(-5 * time.Minute),
			UpdatedAt: now.Add(-4 * time.Minute),
		},
	}

	for _, tx := range txs {
		if _, _, err := ledger.CreateOrGet(tx); err != nil {
			return fmt.Errorf("seed transaction %s: %w", tx.Reference, err)
		}
	}

	if store == nil {
		return nil
	}

	sigValid := true
	attempts := []*webhook.Attempt{
		{
			ID:             uuid.Must(uuid.NewRandom()).String(),
			Ref:            "dev-wh-delivered-001",
			ProviderName:   sampleProvider,
			URL:            "http://localhost:9000/_admin/webhook-receiver",
			Status:         webhook.AttemptStatusDelivered,
			Payload:        []byte(`{"id":"evt_dev_001","object":"event","type":"checkout.session.completed"}`),
			Headers:        map[string]string{"Content-Type": "application/json", "Stripe-Signature": "t=1,v1=abc"},
			SignatureValid: &sigValid,
			TraceID:        "trace-dev-001",
			History: []webhook.AttemptHistory{
				{Time: now.Add(-28 * time.Minute), Status: 200},
			},
			CreatedAt: now.Add(-29 * time.Minute),
			UpdatedAt: now.Add(-28 * time.Minute),
		},
		{
			ID:           uuid.Must(uuid.NewRandom()).String(),
			Ref:          "dev-wh-failed-002",
			ProviderName: sampleProvider,
			URL:          "http://localhost:9000/_admin/webhook-receiver",
			Status:       webhook.AttemptStatusFailed,
			Payload:      []byte(`{"id":"evt_dev_002","object":"event","type":"checkout.session.expired"}`),
			Headers:      map[string]string{"Content-Type": "application/json"},
			TraceID:      "trace-dev-003",
			History: []webhook.AttemptHistory{
				{Time: now.Add(-4 * time.Minute), Status: 0, Error: "connection refused"},
			},
			CreatedAt: now.Add(-4 * time.Minute),
			UpdatedAt: now.Add(-4 * time.Minute),
		},
	}

	for _, a := range attempts {
		if err := store.Save(a); err != nil {
			return fmt.Errorf("seed webhook %s: %w", a.Ref, err)
		}
	}

	return nil
}

func pickEnabledProvider(enabled []string) string {
	for _, name := range enabled {
		if name == "stripe" || name == "fawry" {
			return name
		}
	}
	if len(enabled) > 0 {
		return enabled[0]
	}
	return "default"
}
