package stripe

import (
	"strings"
	"testing"
	"time"
)

func TestSignPayloadProducesHeader(t *testing.T) {
	// Given a payload and secret
	payload := []byte(`{"id":"cs_test_123"}`)
	// #nosec G101 -- test fixture dummy webhook secret
	secret := "whsec_muara"
	timestamp := time.Unix(1700000000, 0)

	// When SignPayload is called
	header := SignPayload(payload, secret, timestamp)

	// Then it returns a header with t= and v1=
	if !strings.Contains(header, "t=1700000000") {
		t.Errorf("header missing timestamp: %q", header)
	}
	if !strings.Contains(header, "v1=") {
		t.Errorf("header missing v1 signature: %q", header)
	}
}

func TestVerifySignatureValid(t *testing.T) {
	// Given a signed payload
	payload := []byte(`{"id":"cs_test_123"}`)
	// #nosec G101 -- test fixture dummy webhook secret
	secret := "whsec_muara"
	header := SignPayload(payload, secret, time.Now())

	// When VerifySignature is called
	err := VerifySignature(payload, header, secret)

	// Then it returns nil
	if err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
}

func TestVerifySignatureTamperedPayload(t *testing.T) {
	// Given a signed payload
	payload := []byte(`{"id":"cs_test_123"}`)
	// #nosec G101 -- test fixture dummy webhook secret
	secret := "whsec_muara"
	header := SignPayload(payload, secret, time.Now())

	// When the payload is tampered
	err := VerifySignature([]byte(`{"id":"cs_test_456"}`), header, secret)

	// Then verification fails
	if err == nil {
		t.Fatal("expected error for tampered payload")
	}
}

func TestVerifySignatureWrongSecret(t *testing.T) {
	// Given a signed payload
	payload := []byte(`{"id":"cs_test_123"}`)
	header := SignPayload(payload, "whsec_correct", time.Now())

	// When verified with the wrong secret
	err := VerifySignature(payload, header, "whsec_wrong")

	// Then verification fails
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestVerifySignatureExpiredTimestamp(t *testing.T) {
	// Given a signed payload with an old timestamp
	payload := []byte(`{"id":"cs_test_123"}`)
	// #nosec G101 -- test fixture dummy webhook secret
	secret := "whsec_muara"
	old := time.Now().Add(-10 * time.Minute)
	header := SignPayload(payload, secret, old)

	// When VerifySignature is called
	err := VerifySignature(payload, header, secret)

	// Then it fails due to timestamp tolerance
	if err == nil {
		t.Fatal("expected error for expired timestamp")
	}
}
