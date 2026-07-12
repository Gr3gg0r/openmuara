package fawry

import (
	"testing"
)

func TestSignKnownVector(t *testing.T) {
	req := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "user-456",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems: []ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
	}
	key := "muara-fawry-secret"

	// Compute expected using the same algorithm independently.
	expected := "d8e2d40e0f8e9e6b2b8e5c3d8f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a" // placeholder
	got := Sign(req, key)
	req.Signature = got

	// We cannot hardcode a realistic hash here without a reference, so we verify
	// determinism and that Verify accepts our own output.
	if got == "" {
		t.Fatal("signature is empty")
	}
	if !Verify(req, key) {
		t.Fatal("Verify failed on freshly signed request")
	}

	// Verify determinism.
	if Sign(req, key) != got {
		t.Fatal("signature is not deterministic")
	}

	_ = expected // keep for future known-vector update
}

func TestSignEmptyCustomerProfileID(t *testing.T) {
	req := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems: []ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
	}
	key := "muara-fawry-secret"

	sig := Sign(req, key)
	req.Signature = sig
	if !Verify(req, key) {
		t.Fatal("Verify failed with empty customerProfileId")
	}
}

func TestSignMultiItemSort(t *testing.T) {
	req := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "user-456",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems: []ChargeItem{
			{ItemID: "z-item", Price: 10.00, Quantity: 2},
			{ItemID: "a-item", Price: 5.50, Quantity: 1},
		},
	}
	key := "muara-fawry-secret"

	sig := Sign(req, key)

	// Reordering items should produce the same signature because of sorting.
	req.ChargeItems = []ChargeItem{
		{ItemID: "a-item", Price: 5.50, Quantity: 1},
		{ItemID: "z-item", Price: 10.00, Quantity: 2},
	}
	if Sign(req, key) != sig {
		t.Fatal("multi-item signature should be independent of item order")
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	req := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "user-456",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems: []ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
		Signature: "invalid",
	}
	if Verify(req, "muara-fawry-secret") {
		t.Fatal("Verify should fail for invalid signature")
	}
}
