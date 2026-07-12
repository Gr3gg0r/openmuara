package stripe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/openmuara/openmuara/internal/errcode"
)

// Error mirrors the error object returned by the Stripe API.
type Error struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse is the top-level JSON wrapper used by Stripe.
type ErrorResponse struct {
	Error Error `json:"error"`
}

func writeStripeError(w http.ResponseWriter, status int, typ, code, param, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: Error{
			Type:    typ,
			Code:    code,
			Param:   param,
			Message: message,
		},
	})
}

func writeStripeInvalidRequestError(w http.ResponseWriter, status int, code, param, message string) {
	writeStripeError(w, status, "invalid_request_error", code, param, message)
}

// errMissingParam returns a validation error for a missing required parameter.
func errMissingParam(param string) error {
	return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("missing required param: %s", param))
}

// errInvalidParam returns a validation error for an invalid parameter.
func errInvalidParam(param, message string) error {
	return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("invalid %s: %s", param, message))
}

// errInvalidPaymentMethodTypes is returned when too many payment method types are supplied.
var errInvalidPaymentMethodTypes = errcode.New(errcode.EInvalidRequest, "invalid payment_method_types")

// errUnsupportedPaymentMethodType returns a validation error for an unsupported payment method type.
func errUnsupportedPaymentMethodType(t string) error {
	return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("unsupported payment_method_type: %s", t))
}

func writeStripeValidationError(w http.ResponseWriter, err error) {
	msg := errcode.Message(err)
	code := "parameter_invalid"
	param := ""
	switch {
	case strings.HasPrefix(msg, "invalid payment_method_types"):
		param = "payment_method_types"
	case strings.HasPrefix(msg, "missing required param:"):
		code = "parameter_missing"
		param = strings.TrimSpace(strings.TrimPrefix(msg, "missing required param:"))
	}
	writeStripeInvalidRequestError(w, http.StatusBadRequest, code, param, msg)
}
