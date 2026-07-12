package ui

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderStripeCheckoutPage(t *testing.T) {
	var buf bytes.Buffer
	data := StripeCheckoutPageData{
		ID:                 "cs_123",
		AmountTotal:        1000,
		Currency:           "myr",
		AmountTotalDisplay: "RM10.00",
		LineItems: []StripeCheckoutLineItem{
			{Name: "Item", Quantity: 1, Amount: 1000, Currency: "myr"},
		},
		PaymentMethodTypes:      []string{"card", "fpx"},
		ShowCard:                true,
		ShowFPX:                 true,
		ShowPaymentMethodToggle: true,
		FPXBanks:                []FPXBank{{Code: "TEST", Name: "Test Bank"}},
		CSRFToken:               "tok",
	}
	if err := RenderStripeCheckoutPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderStripeWebhooksPage(t *testing.T) {
	var buf bytes.Buffer
	data := StripeWebhooksPageData{
		URL:           "http://localhost/webhook",
		WebhookSecret: "whsec_123",
		Events:        []StripeWebhookEventOption{{Name: "invoice.payment_succeeded", Checked: true}},
		CSRFToken:     "tok",
	}
	if err := RenderStripeWebhooksPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderStripePaymentIntentPage(t *testing.T) {
	var buf bytes.Buffer
	data := StripePaymentIntentPageData{
		ID:            "pi_123",
		Amount:        1000,
		Currency:      "myr",
		AmountDisplay: "RM10.00",
		Status:        "requires_action",
		ShowCard:      true,
		ShowFPX:       true,
		FPXBanks:      []FPXBank{{Code: "TEST", Name: "Test Bank"}},
		CSRFToken:     "tok",
	}
	if err := RenderStripePaymentIntentPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderBillplzPayPage(t *testing.T) {
	var buf bytes.Buffer
	data := BillplzPayPageData{
		ID:          "bill-1",
		Amount:      1000,
		Description: "Test bill",
		Methods:     []BillplzPaymentMethod{{Code: "fpx", Name: "FPX"}},
		CSRFToken:   "tok",
	}
	if err := RenderBillplzPayPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderToyyibPayPage(t *testing.T) {
	var buf bytes.Buffer
	data := ToyyibPayPageData{
		BillCode:      "BILL-1",
		BillName:      "Test",
		Amount:        1000,
		AmountDisplay: "RM10.00",
		Channel:       2,
		CSRFToken:     "tok",
	}
	if err := RenderToyyibPayPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderIPay88PayPage(t *testing.T) {
	var buf bytes.Buffer
	data := IPay88PayPageData{
		RefNo:         "ORDER-1",
		Amount:        1000,
		AmountDisplay: "RM10.00",
		Currency:      "MYR",
		Description:   "Test",
		Methods:       []IPay88PaymentMethod{{ID: "fpx", Name: "FPX"}},
		CSRFToken:     "tok",
	}
	if err := RenderIPay88PayPage(&buf, data); err != nil {
		t.Fatalf("render: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestServeStripeCheckoutPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeStripeCheckoutPage(rec, StripeCheckoutPageData{ID: "cs_123", Currency: "myr"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Errorf("content-type: want html, got %q", rec.Header().Get("Content-Type"))
	}
}

func TestServeStripeWebhooksPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeStripeWebhooksPage(rec, StripeWebhooksPageData{URL: "http://localhost/webhook"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestServeStripePaymentIntentPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeStripePaymentIntentPage(rec, StripePaymentIntentPageData{ID: "pi_123"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestServeBillplzPayPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeBillplzPayPage(rec, BillplzPayPageData{ID: "bill-1"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestServeToyyibPayPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeToyyibPayPage(rec, ToyyibPayPageData{BillCode: "BILL-1"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestServeIPay88PayPage(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := ServeIPay88PayPage(rec, IPay88PayPageData{RefNo: "ORDER-1"}); err != nil {
		t.Fatalf("serve: %v", err)
	}
	if rec.Code != 200 {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}
