package fawry

import "testing"

func FuzzSignVerify(f *testing.F) {
	seeds := []ChargeRequest{
		{
			MerchantCode:      "MC001",
			MerchantRefNum:    "REF001",
			CustomerEmail:     "a@example.com",
			CustomerName:      "Alice",
			CustomerProfileID: "CUST001",
			PaymentExpiry:     1234567890,
			Language:          "en",
			ChargeItems:       []ChargeItem{{ItemID: "ITEM1", Price: 10.5, Quantity: 2}},
			ReturnURL:         "http://example.com/return",
		},
		{
			MerchantCode:   "MC002",
			MerchantRefNum: "REF002",
			ChargeItems: []ChargeItem{
				{ItemID: "B", Price: 1.0, Quantity: 1},
				{ItemID: "A", Price: 2.0, Quantity: 3},
			},
			ReturnURL: "http://example.com/callback",
		},
	}
	for _, req := range seeds {
		f.Add(req.MerchantCode, req.MerchantRefNum, req.CustomerProfileID, req.ReturnURL, "secret")
	}

	f.Fuzz(func(t *testing.T, merchantCode, merchantRefNum, customerProfileID, returnURL, secret string) {
		req := ChargeRequest{
			MerchantCode:      merchantCode,
			MerchantRefNum:    merchantRefNum,
			CustomerProfileID: customerProfileID,
			ReturnURL:         returnURL,
			ChargeItems:       []ChargeItem{{ItemID: "ITEM1", Price: 9.99, Quantity: 1}},
		}
		req.Signature = Sign(req, secret)

		if !Verify(req, secret) {
			t.Errorf("sign/verify round-trip failed")
		}

		req.Signature += "x"
		if Verify(req, secret) {
			t.Errorf("tampered signature should fail verification")
		}
	})
}
