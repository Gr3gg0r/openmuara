package stripe

import (
	"errors"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
)

func assertErrcode(t *testing.T, err error, want errcode.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	var ec *errcode.Error
	if !errors.As(err, &ec) {
		t.Fatalf("expected *errcode.Error, got %T", err)
	}
	if ec.Code != want {
		t.Errorf("errcode: want %q, got %q", want, ec.Code)
	}
}

func TestProviderInitMissingPublishableKeyHasErrcode(t *testing.T) {
	p := NewProvider()
	err := p.Init(map[string]any{
		"secret_key":     "sk_test_xxx",
		"webhook_secret": "whsec_xxx",
	})
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestBuildPayloadMissingPaymentIntentHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	_, err := p.buildPayload("pi_test_missing", "PAID")
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestBuildPayloadMissingSessionHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	_, err := p.buildPayload("cs_test_missing", "PAID")
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestVerifySignatureMismatchHasErrcode(t *testing.T) {
	secret := "whsec_test"
	payload := []byte(`{"id":"evt_1"}`)
	header := SignPayload(payload, secret, time.Now())
	err := VerifySignature(payload, header, "wrong_secret")
	assertErrcode(t, err, errcode.ESignatureMismatch)
}

func TestVerifySignatureMissingTimestampHasErrcode(t *testing.T) {
	err := VerifySignature([]byte(`{}`), "v1=abc", "secret")
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestUpdateLedgerStatusNotFoundHasErrcode(t *testing.T) {
	ledger := engine.NewMemoryStore()
	err := updateLedgerStatus(ledger, "missing", engine.TransactionStatusPaid)
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestValidateCreateRequestMissingSuccessURLHasErrcode(t *testing.T) {
	err := validateCreateRequest(CreateCheckoutSessionRequest{LineItems: []LineItem{{Quantity: 1, PriceData: &PriceData{Currency: "myr", UnitAmount: 100, ProductData: struct {
		Name string `json:"name"`
	}{Name: "x"}}}}})
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestValidatePaymentIntentRequestInvalidAmountHasErrcode(t *testing.T) {
	err := validatePaymentIntentRequest(PaymentIntentRequest{Amount: 0, Currency: "myr"})
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestErrMissingParamMessage(t *testing.T) {
	err := errMissingParam("currency")
	if got := errcode.Message(err); got != "missing required param: currency" {
		t.Errorf("Message = %q, want %q", got, "missing required param: currency")
	}
}

func validProviderConfig() map[string]any {
	return map[string]any{
		"publishable_key": "pk_test_xxx",
		"secret_key":      "sk_test_xxx",
		"webhook_secret":  "whsec_xxx",
	}
}
