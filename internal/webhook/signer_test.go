package webhook

import (
	"testing"
)

func TestHMACSignerDeterministic(t *testing.T) {
	signer := NewHMACSigner("dummy-secret")
	payload := FawryV2Payload{
		RequestID:         "req-1",
		MerchantRefNumber: "ref-1",
		OrderStatus:       "PAID",
		PaymentAmount:     100.00,
	}

	sig1, err := signer.Sign(payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	sig2, err := signer.Sign(payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	if sig1 != sig2 {
		t.Errorf("signatures not deterministic: %q vs %q", sig1, sig2)
	}
}

func TestHMACSignerVerify(t *testing.T) {
	signer := NewHMACSigner("dummy-secret")
	payload := FawryV2Payload{
		RequestID:         "req-1",
		MerchantRefNumber: "ref-1",
		OrderStatus:       "PAID",
		PaymentAmount:     100.00,
	}

	sig, err := signer.Sign(payload)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	ok, err := signer.Verify(payload, sig)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Error("expected signature to verify")
	}

	ok, err = signer.Verify(payload, "bad-signature")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if ok {
		t.Error("expected bad signature to fail verification")
	}
}

func TestHMACSignerDifferentSecrets(t *testing.T) {
	s1 := NewHMACSigner("secret-a")
	s2 := NewHMACSigner("secret-b")
	payload := FawryV2Payload{
		RequestID:         "req-1",
		MerchantRefNumber: "ref-1",
		OrderStatus:       "PAID",
		PaymentAmount:     100.00,
	}

	sig1, _ := s1.Sign(payload)
	sig2, _ := s2.Sign(payload)

	if sig1 == sig2 {
		t.Error("signatures from different secrets should differ")
	}
}
