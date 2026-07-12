package fawry

import (
	"context"
	"errors"
	"testing"

	"github.com/openmuara/openmuara/internal/errcode"
	v2 "github.com/openmuara/openmuara/internal/fawry/v2"
	"github.com/openmuara/openmuara/internal/provider"
)

func validConfig() map[string]any {
	// #nosec G101 -- test fixture dummy credentials
	return map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
	}
}

func validConfigV2() map[string]any {
	cfg := validConfig()
	cfg["version"] = "v2"
	return cfg
}

func validChargeRequest() ChargeRequest {
	req := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerEmail:     "test@example.com",
		CustomerName:      "Test User",
		CustomerProfileID: "user-456",
		PaymentExpiry:     1234567890000,
		Language:          "en-gb",
		ChargeItems: []ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
		ReturnURL: "http://127.0.0.1:9999/callback",
	}
	req.Signature = Sign(req, "muara-fawry-secret")
	return req
}

func assertErrcode(t *testing.T, err error, want errcode.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	var ec *errcode.Error
	if !errors.As(err, &ec) {
		t.Fatalf("expected *errcode.Error, got %T", err)
	}
	if ec.Code != want {
		t.Errorf("errcode: want %q, got %q", want, ec.Code)
	}
}

func TestProviderInitMissingMerchantCodeHasErrcode(t *testing.T) {
	p := NewProvider()
	cfg := validConfig()
	delete(cfg, "merchant_code")
	err := p.Init(cfg)
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestProviderInitUnknownVersionHasErrcode(t *testing.T) {
	p := NewProvider()
	cfg := validConfig()
	cfg["version"] = "v3"
	err := p.Init(cfg)
	assertErrcode(t, err, errcode.EProviderVersionUnsupported)
}

func TestValidateChargeRequestMissingMerchantCodeHasErrcode(t *testing.T) {
	req := validChargeRequest()
	req.MerchantCode = ""
	err := validateChargeRequest(req)
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestValidateChargeRequestMissingSignatureHasErrcode(t *testing.T) {
	req := validChargeRequest()
	req.Signature = ""
	err := validateChargeRequest(req)
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestProviderVerifyWebhookSignatureInvalidJSONHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	_, err := p.VerifyWebhookSignature([]byte("not json"), nil)
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestProviderVerifyWebhookSignatureMissingSignatureHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	payload := []byte(`{"merchantRefNumber":"ref-1","orderStatus":"PAID"}`)
	_, err := p.VerifyWebhookSignature(payload, nil)
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestProviderPayloadBuilderV2MissingTransactionHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validConfigV2()); err != nil {
		t.Fatalf("init: %v", err)
	}
	builder := p.PayloadBuilder()
	_, err := builder(context.Background(), provider.Transaction{Reference: "missing"})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestPayloadBuilderV2StoreNotConfiguredHasErrcode(t *testing.T) {
	builder := v2.NewPayloadBuilder("secret", nil)
	_, err := builder(context.Background(), provider.Transaction{Reference: "ref"})
	assertErrcode(t, err, errcode.EInternal)
}
