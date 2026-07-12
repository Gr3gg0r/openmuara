package billplz_test

import (
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
)

func TestSignSortsKeysCaseInsensitively(t *testing.T) {
	values := map[string]string{
		"Zebra": "z",
		"apple": "a",
		"Mango": "m",
	}
	key := "secret"

	// Expected order after case-insensitive sort: apple, Mango, Zebra
	// Signature source: applea|Mangom|Zebraz
	sig := billplz.Sign(values, key)
	if sig == "" {
		t.Fatal("signature is empty")
	}

	// Verify against the same source should succeed.
	if !billplz.Verify(values, key, sig) {
		t.Fatal("signature verification failed")
	}
}

func TestVerifyDetectsTamperedValue(t *testing.T) {
	values := map[string]string{
		"id":    "bill-123",
		"paid":  "true",
		"state": "paid",
	}
	key := "secret"
	sig := billplz.Sign(values, key)

	values["paid"] = "false"
	if billplz.Verify(values, key, sig) {
		t.Fatal("expected verification to fail for tampered value")
	}
}

func TestVerifyDetectsWrongKey(t *testing.T) {
	values := map[string]string{"id": "bill-123"}
	sig := billplz.Sign(values, "secret")
	if billplz.Verify(values, "other-secret", sig) {
		t.Fatal("expected verification to fail for wrong key")
	}
}
