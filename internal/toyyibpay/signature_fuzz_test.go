package toyyibpay

import (
	"net/url"
	"testing"
)

func FuzzVerifyCallback(f *testing.F) {
	seeds := []struct {
		secret, status, orderID, refno string
	}{
		{"secret", "1", "ORDER-1", "REF-1"},
		{"", "", "", ""},
	}
	for _, s := range seeds {
		f.Add(s.secret, s.status, s.orderID, s.refno)
	}

	f.Fuzz(func(t *testing.T, secret, status, orderID, refno string) {
		values := url.Values{}
		values.Set("status", status)
		values.Set("order_id", orderID)
		values.Set("refno", refno)
		values.Set("hash", ComputeHash(secret, status, orderID, refno))

		if !VerifyCallback(secret, values) {
			t.Errorf("computed hash should verify")
		}

		values.Set("hash", values.Get("hash")+"x")
		if VerifyCallback(secret, values) {
			t.Errorf("tampered hash should fail verification")
		}
	})
}
