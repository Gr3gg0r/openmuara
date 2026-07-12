package v2_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	v2 "github.com/Gr3gg0r/openmuara/internal/fawry/v2"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestProviderMethods(t *testing.T) {
	p := v2.NewProvider("secret")
	if p == nil {
		t.Fatal("expected non-nil provider")
	}

	p.SetStore(nil)
	p.SetDispatcher(nil)

	h := p.WebhookHandler()
	if h == nil {
		t.Fatal("expected non-nil webhook handler")
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/webhook?token=secret", bytes.NewReader([]byte("{}"))))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("webhook handler status: want 401, got %d", rec.Code)
	}
}

func TestPayloadBuilderNilStore(t *testing.T) {
	builder := v2.NewPayloadBuilder("secret", nil)
	_, err := builder(context.Background(), provider.Transaction{Reference: "ref-1"})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestPayloadBuilderTransactionNotFound(t *testing.T) {
	builder := v2.NewPayloadBuilder("secret", engine.NewMemoryStore())
	_, err := builder(context.Background(), provider.Transaction{Reference: "missing"})
	if err == nil {
		t.Fatal("expected error when transaction not found")
	}
}

func TestPayloadBuilderSuccess(t *testing.T) {
	store := engine.NewMemoryStore()
	_, _, err := store.CreateOrGet(engine.NewTransaction(engine.Transaction{
		ID:          "tx-1",
		Reference:   "ref-1",
		Provider:    "fawry",
		Type:        "charge",
		Amount:      99.99,
		Currency:    "EGP",
		Status:      engine.TransactionStatusPaid,
		CustomerRef: "cust-1",
		Items: []engine.TransactionItem{
			{ItemCode: "ITEM-1", Price: 49.99, Quantity: 2},
		},
	}))
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}

	builder := v2.NewPayloadBuilder("secret", store)
	body, err := builder(context.Background(), provider.Transaction{Reference: "ref-1"})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}

	var payload webhook.FawryV2Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.MerchantRefNumber != "ref-1" {
		t.Errorf("merchantRefNumber: want ref-1, got %q", payload.MerchantRefNumber)
	}
	if payload.OrderAmount != 99.99 {
		t.Errorf("orderAmount: want 99.99, got %f", payload.OrderAmount)
	}
	if payload.CustomerMerchantID != "cust-1" {
		t.Errorf("customerMerchantID: want cust-1, got %q", payload.CustomerMerchantID)
	}
	if len(payload.OrderItems) != 1 || payload.OrderItems[0].ItemCode != "ITEM-1" {
		t.Errorf("unexpected order items: %+v", payload.OrderItems)
	}
	if payload.MessageSignature == "" {
		t.Error("expected non-empty message signature")
	}
}

func TestPayloadBuilderFallsBackToTxAmountAndStatus(t *testing.T) {
	store := engine.NewMemoryStore()
	_, _, err := store.CreateOrGet(engine.NewTransaction(engine.Transaction{
		ID:        "tx-2",
		Reference: "ref-2",
		Provider:  "fawry",
		Amount:    0,
		Currency:  "EGP",
		Status:    engine.TransactionStatusUnpaid,
	}))
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}

	builder := v2.NewPayloadBuilder("secret", store)
	body, err := builder(context.Background(), provider.Transaction{Reference: "ref-2", Amount: 10.5, Status: "PAID"})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}

	var payload webhook.FawryV2Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.PaymentAmount != 10.5 {
		t.Errorf("paymentAmount: want 10.5, got %f", payload.PaymentAmount)
	}
	if payload.OrderStatus != "PAID" {
		t.Errorf("orderStatus: want PAID, got %q", payload.OrderStatus)
	}
}
