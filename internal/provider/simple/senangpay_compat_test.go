package simple

import (
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/senangpay"
)

func TestSenangpaySignatureMatchesBuiltIn(t *testing.T) {
	// #nosec G101 -- test fixture dummy credential, not a real secret
	secret := "muara-senangpay-secret"
	req := senangpay.ChargeRequest{
		Detail:  "Test payment",
		Amount:  10.00,
		OrderID: "ORDER-1",
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+60123456789",
	}
	expected := senangpay.Sign(secret, req.Detail, req.Amount, req.OrderID)

	cfg := senangpayGatewayConfig()
	p := NewProvider(cfg)
	if err := p.Init(map[string]any{"secret_key": secret}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.secret = secret

	values := map[string]any{
		"detail":   req.Detail,
		"amount":   req.Amount,
		"order_id": req.OrderID,
		"hash":     expected,
	}
	if !p.verifySignature(values) {
		t.Error("expected senangpay signature to verify")
	}
}

func senangpayGatewayConfig() plugin.GatewayConfig {
	cfg := validFawryConfig()
	cfg.Metadata.Name = "senangpay"
	cfg.Runtime.Simple.ChargeRoute = "senangpay_charge"
	cfg.Runtime.Simple.ReferenceField = "order_id"
	cfg.Runtime.Simple.AmountField = "amount"
	cfg.Runtime.Simple.CustomerField = "email"
	cfg.Runtime.Simple.Currency = "MYR"
	cfg.Runtime.Simple.ResponseTemplate = map[string]any{
		"status":    "ok",
		"order_id":  "{{ .Reference }}",
		"reference": "{{ .Reference }}",
	}
	cfg.Signature = &plugin.Signature{
		Algorithm: "senangpay_md5",
		Fields:    []string{"detail", "amount", "order_id", "hash"},
		SecretKey: "secret_key",
	}
	cfg.Schemas.Requests["charge_request"] = plugin.Schema{Fields: []plugin.Field{
		{Name: "detail", JSONName: "detail", Type: "string", Required: true},
		{Name: "amount", JSONName: "amount", Type: "number", Required: true},
		{Name: "order_id", JSONName: "order_id", Type: "string", Required: true},
		{Name: "hash", JSONName: "hash", Type: "string", Required: true},
	}}
	return cfg
}
