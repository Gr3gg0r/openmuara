package ipay88

import (
	"testing"
)

func TestStripAmount(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"1,278.99", "127899"},
		{"12.00", "1200"},
		{"12.5", "1250"},
		{"1000", "100000"},
		{"1,000,000.00", "100000000"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := stripAmount(tc.input)
			if got != tc.want {
				t.Errorf("stripAmount(%q): want %q, got %q", tc.input, tc.want, got)
			}
		})
	}
}

func TestSignRequest(t *testing.T) {
	key := "secret"
	code := "M00001"
	ref := "REF-001"
	amount := "1,278.99"
	currency := "MYR"

	expected := sha256Hex(key + code + ref + "127899" + currency)
	got := SignRequest(key, code, ref, amount, currency)
	if got != expected {
		t.Errorf("signature mismatch: want %q, got %q", expected, got)
	}
}

func TestSignResponse(t *testing.T) {
	key := "secret"
	code := "M00001"
	paymentID := "2"
	ref := "REF-001"
	amount := "12.50"
	currency := "MYR"
	status := "1"

	expected := sha256Hex(key + code + paymentID + ref + "1250" + currency + status)
	got := SignResponse(key, code, paymentID, ref, amount, currency, status)
	if got != expected {
		t.Errorf("signature mismatch: want %q, got %q", expected, got)
	}
}

func TestVerifyRequest(t *testing.T) {
	req := PaymentRequest{
		MerchantCode: "M00001",
		RefNo:        "REF-001",
		Amount:       "12.50",
		Currency:     "MYR",
		Signature:    SignRequest("secret", "M00001", "REF-001", "12.50", "MYR"),
	}
	if !VerifyRequest(req, "secret") {
		t.Error("expected valid request signature")
	}
	if VerifyRequest(req, "wrong") {
		t.Error("expected invalid request signature with wrong key")
	}
}

func TestVerifyResponse(t *testing.T) {
	key := "secret"
	code := "M00001"
	paymentID := "2"
	ref := "REF-001"
	amount := "12.50"
	currency := "MYR"
	status := "1"
	sig := SignResponse(key, code, paymentID, ref, amount, currency, status)

	if !VerifyResponse(key, code, paymentID, ref, amount, currency, status, sig) {
		t.Error("expected valid response signature")
	}
	if VerifyResponse("wrong", code, paymentID, ref, amount, currency, status, sig) {
		t.Error("expected invalid response signature with wrong key")
	}
}
