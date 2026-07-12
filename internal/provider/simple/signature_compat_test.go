package simple

import (
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/fawry"
)

func TestFawrySignatureMatchesBuiltIn(t *testing.T) {
	req := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "user-123",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
	}
	// #nosec G101 -- test fixture dummy credential, not a real secret
	secret := "muara-fawry-secret"

	expected := fawry.Sign(req, secret)

	values := map[string]any{
		"merchantCode":      req.MerchantCode,
		"merchantRefNum":    req.MerchantRefNum,
		"customerProfileId": req.CustomerProfileID,
		"returnUrl":         req.ReturnURL,
		"chargeItems": []any{
			map[string]any{"itemId": "prod_test_123", "price": 99.99, "quantity": float64(1)},
		},
	}

	p := NewProvider(validFawryConfig())
	if err := p.Init(map[string]any{"secret": secret}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.secret = secret

	got := p.sign(values)
	if got != expected {
		t.Errorf("signature mismatch\nexpected: %s\ngot:      %s", expected, got)
	}
}
