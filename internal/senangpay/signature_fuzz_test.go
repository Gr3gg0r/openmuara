package senangpay

import "testing"

func FuzzSignVerify(f *testing.F) {
	seeds := []struct {
		secret, detail string
		amount         float64
		orderID        string
	}{
		{"secret", "Test product", 12.34, "ORDER-1"},
		{"", "", 0, ""},
	}
	for _, s := range seeds {
		f.Add(s.secret, s.detail, s.amount, s.orderID)
	}

	f.Fuzz(func(t *testing.T, secret, detail string, amount float64, orderID string) {
		req := ChargeRequest{
			Detail:  detail,
			Amount:  amount,
			OrderID: orderID,
			Hash:    Sign(secret, detail, amount, orderID),
		}

		if !Verify(req, secret) {
			t.Errorf("sign/verify round-trip failed")
		}

		req.Hash += "x"
		if Verify(req, secret) {
			t.Errorf("tampered hash should fail verification")
		}
	})
}
