package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/Gr3gg0r/openmuara/internal/fawry/v1"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestProviderMethods(t *testing.T) {
	p := v1.NewProvider("secret")
	if p == nil {
		t.Fatal("expected non-nil provider")
	}

	p.SetStore(nil)
	p.SetDispatcher(nil)

	h := p.WebhookHandler()
	if h == nil {
		t.Fatal("expected non-nil webhook handler")
	}

	body, err := p.PayloadBuilder()(context.Background(), provider.Transaction{Reference: "ref-1", Status: "PAID"})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}

	var payload v1.WebhookBody
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.MerchantRefNumber != "ref-1" {
		t.Errorf("merchantRefNumber: want ref-1, got %q", payload.MerchantRefNumber)
	}
	if payload.OrderStatus != "PAID" {
		t.Errorf("orderStatus: want PAID, got %q", payload.OrderStatus)
	}
	if payload.MessageSignature == "" {
		t.Error("expected non-empty message signature")
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/webhook?token=secret", bytes.NewReader([]byte("{}"))))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("webhook handler status: want 401, got %d", rec.Code)
	}
}
