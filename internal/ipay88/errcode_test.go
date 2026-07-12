package ipay88

import (
	"errors"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
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

func TestProviderInitMissingMerchantCodeHasErrcode(t *testing.T) {
	p := NewProvider()
	err := p.Init(map[string]any{"merchant_key": "key"})
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestValidateEntryRequestMissingSignatureHasErrcode(t *testing.T) {
	err := validateEntryRequest(PaymentRequest{MerchantCode: "M", RefNo: "R", Amount: "1.00", Currency: "MYR", ProdDesc: "x", UserName: "u", UserEmail: "e", UserContact: "c", SignatureType: "SHA256", ResponseURL: "http://example.com", BackendURL: "http://example.com"})
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestIsPublicURLInvalidSchemeHasErrcode(t *testing.T) {
	err := IsPublicURL("ftp://example.com")
	assertErrcode(t, err, errcode.EInvalidRequest)
}

func TestIsPublicURLUnsafeWrapsSentinel(t *testing.T) {
	err := IsPublicURL("http://127.0.0.1/callback")
	assertErrcode(t, err, errcode.EInvalidRequest)
	if !errors.Is(err, ErrUnsafeURL) {
		t.Error("expected errors.Is(err, ErrUnsafeURL) to remain true")
	}
}

func TestTransitionTransactionNotFoundHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"merchant_code": "M", "merchant_key": "K"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := p.transitionTransaction("missing", engine.TransactionStatusPaid)
	assertErrcode(t, err, errcode.ETransactionTransitionInvalid)
}
