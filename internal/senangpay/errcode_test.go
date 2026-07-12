package senangpay

import (
	"errors"
	"testing"

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

func TestProviderInitMissingSecretKeyHasErrcode(t *testing.T) {
	p := NewProvider()
	err := p.Init(map[string]any{})
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestValidateChargeRequestMissingHashHasErrcode(t *testing.T) {
	err := validateChargeRequest(ChargeRequest{Detail: "x", Amount: 10, OrderID: "123"})
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestApplyCallbackEmptyOrderIDHasErrcode(t *testing.T) {
	err := applyCallback(engine.NewMemoryStore(), CallbackQuery{})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestApplyCallbackOrderNotFoundHasErrcode(t *testing.T) {
	err := applyCallback(engine.NewMemoryStore(), CallbackQuery{OrderID: "missing"})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestApplyCallbackPreservesErrorsIs(t *testing.T) {
	err := applyCallback(engine.NewMemoryStore(), CallbackQuery{})
	if !errors.Is(err, ErrOrderNotFound) {
		t.Error("expected errors.Is(err, ErrOrderNotFound) to remain true")
	}
}
