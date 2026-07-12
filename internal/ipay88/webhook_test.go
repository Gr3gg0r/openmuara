package ipay88

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestBuildPayloadSuccess(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-WEB-1")

	payload, err := p.buildPayload(context.Background(), provider.Transaction{
		Reference: "REF-WEB-1",
		Status:    string(webhook.PaymentStatusPaid),
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	values, err := url.ParseQuery(string(payload))
	if err != nil {
		t.Fatalf("parse payload: %v", err)
	}
	if values.Get("RefNo") != "REF-WEB-1" {
		t.Errorf("refno: want REF-WEB-1, got %q", values.Get("RefNo"))
	}
	if values.Get("Status") != "1" {
		t.Errorf("status: want 1, got %q", values.Get("Status"))
	}
	if values.Get("SignatureType") != "SHA256" {
		t.Errorf("signature type: want SHA256, got %q", values.Get("SignatureType"))
	}
	if values.Get("Signature") == "" {
		t.Error("expected non-empty signature")
	}
}

func TestBuildPayloadRequestNotFound(t *testing.T) {
	p := newTestProvider(t)
	_, err := p.buildPayload(context.Background(), provider.Transaction{Reference: "missing"})
	if err == nil {
		t.Fatal("expected error for missing request")
	}
}

func TestBuildPayloadPaymentIDFallback(t *testing.T) {
	p := newTestProvider(t)
	seedPaymentRequest(t, p, "REF-WEB-2")

	stored, _ := p.getRequest("REF-WEB-2")
	stored.SelectedPaymentID = ""
	p.saveRequest(stored)

	payload, err := p.buildPayload(context.Background(), provider.Transaction{
		Reference: "REF-WEB-2",
		Status:    string(webhook.PaymentStatusPaid),
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	if !strings.Contains(string(payload), "PaymentId=2") {
		t.Errorf("expected default PaymentId=2 in payload, got %q", string(payload))
	}
}

func TestMapWebhookStatusToIPay88(t *testing.T) {
	if got := mapWebhookStatusToIPay88(string(webhook.PaymentStatusPaid)); got != "1" {
		t.Errorf("paid: want 1, got %q", got)
	}
	if got := mapWebhookStatusToIPay88(string(webhook.PaymentStatusUnpaid)); got != "0" {
		t.Errorf("unpaid: want 0, got %q", got)
	}
	if got := mapWebhookStatusToIPay88("unknown"); got != "0" {
		t.Errorf("unknown: want 0, got %q", got)
	}
}
