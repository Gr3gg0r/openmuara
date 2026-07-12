// Package errcode provides a stable error-code taxonomy for OpenMuara.
// Codes are grouped by area so bugs and support requests can be classified by
// symptom without parsing free-form error messages.
package errcode

import "fmt"

// Code is a stable error identifier.
type Code string

// Error codes are grouped by area:
//
//	E1xx — generic / internal
//	E2xx — configuration
//	E3xx — provider emulation
//	E4xx — webhook dispatch
//	E5xx — transaction / ledger
//	E6xx — signature / security
const (
	EInternal        Code = "E1000"
	EUnknownProvider Code = "E1001"
	EInvalidRequest  Code = "E1002"

	EConfigMissing    Code = "E2000"
	EConfigInvalid    Code = "E2001"
	EProviderDisabled Code = "E2002"

	EProviderChargeFailed       Code = "E3000"
	EProviderCallbackFailed     Code = "E3001"
	EProviderEscapeFailed       Code = "E3002"
	EProviderVersionUnsupported Code = "E3003"

	EWebhookURLMissing     Code = "E4000"
	EWebhookBuildFailed    Code = "E4001"
	EWebhookDeliveryFailed Code = "E4002"
	EWebhookReplayNotFound Code = "E4003"

	ETransactionNotFound          Code = "E5000"
	ETransactionTransitionInvalid Code = "E5001"
	ETransactionDuplicate         Code = "E5002"

	ESignatureMismatch Code = "E6000"
	ESignatureMissing  Code = "E6001"
)

// Error wraps an underlying error with a stable code.
type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Cause }

// Message returns the human-readable message of err if it is an *Error,
// otherwise it returns err.Error().
func Message(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Message
	}
	return err.Error()
}

// Wrap returns a new *Error with the given code and message, wrapping cause.
func Wrap(code Code, message string, cause error) error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// New returns a new *Error without an underlying cause.
func New(code Code, message string) error {
	return &Error{Code: code, Message: message}
}
