package webhook

import (
	"encoding/json"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestFawryV2BuilderProducesValidJSON(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{
		Provider:    "fawry",
		Type:        "charge",
		Amount:      100.00,
		Currency:    "EGP",
		Status:      engine.TransactionStatusNew,
		CustomerRef: "user-123",
		Reference:   "ref-123",
		Items: []engine.TransactionItem{
			{ItemCode: "prod-1", Price: 100.00, Quantity: 1},
		},
	}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	builder := NewFawryV2Builder("secret", store)
	payload, err := builder.Build("ref-123", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var v FawryV2Payload
	if err := json.Unmarshal(payload, &v); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if v.MerchantRefNumber != "ref-123" {
		t.Errorf("merchant ref: want ref-123, got %q", v.MerchantRefNumber)
	}
	if v.OrderStatus != "PAID" {
		t.Errorf("order status: want PAID, got %q", v.OrderStatus)
	}
	if v.MessageSignature == "" {
		t.Error("message signature is empty")
	}
	if v.PaymentAmount != 100.00 {
		t.Errorf("payment amount: want 100.00, got %f", v.PaymentAmount)
	}
	if v.CustomerMerchantID != "user-123" {
		t.Errorf("customer merchant id: want user-123, got %q", v.CustomerMerchantID)
	}
	if len(v.OrderItems) != 1 {
		t.Fatalf("order items: want 1, got %d", len(v.OrderItems))
	}
	if v.OrderItems[0].ItemCode != "prod-1" {
		t.Errorf("item code: want prod-1, got %q", v.OrderItems[0].ItemCode)
	}
}

func TestFawryV2BuilderStatusUnpaid(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{
		Provider:  "fawry",
		Type:      "charge",
		Amount:    50.00,
		Reference: "ref-456",
	}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	builder := NewFawryV2Builder("secret", store)
	payload, err := builder.Build("ref-456", PaymentStatusUnpaid)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var v FawryV2Payload
	if err := json.Unmarshal(payload, &v); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if v.OrderStatus != "UNPAID" {
		t.Errorf("order status: want UNPAID, got %q", v.OrderStatus)
	}
}

func TestFawryV2BuilderMissingTransaction(t *testing.T) {
	store := engine.NewMemoryStore()
	builder := NewFawryV2Builder("secret", store)

	_, err := builder.Build("missing", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error for missing transaction")
	}
}

func TestPayloadForUnsupportedVersion(t *testing.T) {
	store := engine.NewMemoryStore()
	_, err := PayloadFor("unknown", "secret", store)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}
