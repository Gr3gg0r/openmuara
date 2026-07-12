package errcode

import (
	"errors"
	"testing"
)

func TestNewFormatsError(t *testing.T) {
	err := New(EConfigMissing, "config is required")
	want := "[E2000] config is required"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestWrapFormatsErrorWithCause(t *testing.T) {
	cause := errors.New("underlying")
	err := Wrap(EConfigInvalid, "config invalid", cause)
	want := "[E2001] config invalid: underlying"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestWrapUnwrapsCause(t *testing.T) {
	cause := errors.New("underlying")
	err := Wrap(EConfigInvalid, "config invalid", cause)
	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to match cause")
	}
}

func TestMessageReturnsHumanReadableText(t *testing.T) {
	ec := New(ETransactionNotFound, "transaction not found")
	if got := Message(ec); got != "transaction not found" {
		t.Errorf("Message = %q, want %q", got, "transaction not found")
	}
}

func TestMessageFallsBackToErrorString(t *testing.T) {
	err := errors.New("plain error")
	if got := Message(err); got != "plain error" {
		t.Errorf("Message = %q, want %q", got, "plain error")
	}
}

func TestErrorsAsIntoError(t *testing.T) {
	err := Wrap(ESignatureMismatch, "signature mismatch", errors.New("cause"))
	var ec *Error
	if !errors.As(err, &ec) {
		t.Fatal("expected errors.As to match *errcode.Error")
	}
	if ec.Code != ESignatureMismatch {
		t.Errorf("Code = %q, want %q", ec.Code, ESignatureMismatch)
	}
	if ec.Message != "signature mismatch" {
		t.Errorf("Message = %q, want %q", ec.Message, "signature mismatch")
	}
}
