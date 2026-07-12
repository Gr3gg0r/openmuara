package billplz

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/provider"
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

func TestProviderInitMissingAPIKeyHasErrcode(t *testing.T) {
	p := NewProvider()
	err := p.Init(map[string]any{"x_signature_key": "x"})
	assertErrcode(t, err, errcode.EConfigMissing)
}

func TestRequireBasicAuthMissingAuthHasErrcode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := requireBasicAuth(req, "key")
	assertErrcode(t, err, errcode.ESignatureMissing)
}

func TestRequireBasicAuthInvalidKeyHasErrcode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("wrong", "")
	err := requireBasicAuth(req, "key")
	assertErrcode(t, err, errcode.ESignatureMismatch)
}

func TestValidateCreateBillRequestMissingCollectionHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"api_key": "k", "x_signature_key": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	err := validateCreateBillRequest(CreateBillRequest{CollectionID: "missing"}, p.collections)
	assertErrcode(t, err, errcode.ETransactionNotFound)
}

func TestBuildPayloadMissingBillHasErrcode(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"api_key": "k", "x_signature_key": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	_, err := p.buildPayload(context.Background(), provider.Transaction{Reference: "missing"})
	assertErrcode(t, err, errcode.ETransactionNotFound)
}
