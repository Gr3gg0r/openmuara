package stripe

import "strings"

const (
	// PaymentMethodTypeCard is the Stripe card payment method type.
	PaymentMethodTypeCard = "card"
	// PaymentMethodTypeFPX is the Stripe FPX payment method type.
	PaymentMethodTypeFPX = "fpx"
)

// PaymentIntent mirrors a subset of the Stripe PaymentIntent object.
type PaymentIntent struct {
	ID                 string                   `json:"id"`
	Object             string                   `json:"object"`
	Amount             int64                    `json:"amount"`
	Currency           string                   `json:"currency"`
	Status             string                   `json:"status"`
	ClientSecret       string                   `json:"client_secret"`
	PaymentMethodTypes []string                 `json:"payment_method_types"`
	NextAction         *PaymentIntentNextAction `json:"next_action,omitempty"`
	Metadata           map[string]string        `json:"metadata,omitempty"`
	Livemode           bool                     `json:"livemode"`
	ConfirmationMethod string                   `json:"confirmation_method"`
	PaymentMethod      string                   `json:"payment_method,omitempty"`
}

// Clone returns a deep copy of the PaymentIntent.
func (p *PaymentIntent) Clone() *PaymentIntent {
	metadata := make(map[string]string, len(p.Metadata))
	for k, v := range p.Metadata {
		metadata[k] = v
	}
	types := make([]string, len(p.PaymentMethodTypes))
	copy(types, p.PaymentMethodTypes)

	var next *PaymentIntentNextAction
	if p.NextAction != nil {
		next = &PaymentIntentNextAction{
			Type:          p.NextAction.Type,
			RedirectToURL: p.NextAction.RedirectToURL,
		}
	}

	return &PaymentIntent{
		ID:                 p.ID,
		Object:             p.Object,
		Amount:             p.Amount,
		Currency:           p.Currency,
		Status:             p.Status,
		ClientSecret:       p.ClientSecret,
		PaymentMethodTypes: types,
		NextAction:         next,
		Metadata:           metadata,
		Livemode:           p.Livemode,
		ConfirmationMethod: p.ConfirmationMethod,
		PaymentMethod:      p.PaymentMethod,
	}
}

// PaymentIntentRequest mirrors the create parameters for a PaymentIntent.
type PaymentIntentRequest struct {
	Amount             int64             `json:"amount"`
	Currency           string            `json:"currency"`
	PaymentMethodTypes []string          `json:"payment_method_types"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	ReceiptEmail       string            `json:"receipt_email,omitempty"`
}

// PaymentIntentConfirmRequest mirrors the confirm parameters for a PaymentIntent.
type PaymentIntentConfirmRequest struct {
	PaymentMethod string `json:"payment_method"`
}

// PaymentIntentNextAction represents the next_action field on a PaymentIntent.
type PaymentIntentNextAction struct {
	Type          string                      `json:"type"`
	RedirectToURL *PaymentIntentRedirectToURL `json:"redirect_to_url,omitempty"`
}

// PaymentIntentRedirectToURL contains the redirect URL for FPX authentication.
type PaymentIntentRedirectToURL struct {
	URL       string `json:"url"`
	ReturnURL string `json:"return_url,omitempty"`
}

// cardPaymentMethods are the test card tokens accepted by confirm.
var cardPaymentMethods = map[string]struct{}{
	"pm_card_visa":       {},
	"pm_card_mastercard": {},
}

// fpxPaymentMethods are the test FPX tokens accepted by confirm.
var fpxPaymentMethods = map[string]struct{}{
	"pm_fpx_maybank":    {},
	"pm_fpx_cimb":       {},
	"pm_fpx_publicbank": {},
	"pm_fpx_rhb":        {},
	"pm_fpx_hongleong":  {},
	"pm_fpx_ambank":     {},
	"pm_fpx_bankislam":  {},
	"pm_fpx_affinbank":  {},
}

func normalizePaymentIntentMethodTypes(types []string) []string {
	if len(types) == 0 {
		return []string{PaymentMethodTypeCard}
	}
	out := make([]string, len(types))
	for i, t := range types {
		out[i] = strings.ToLower(t)
	}
	return out
}

func validatePaymentIntentRequest(req PaymentIntentRequest) error {
	if req.Amount <= 0 {
		return errInvalidParam("amount", "must be greater than 0")
	}
	if req.Currency == "" {
		return errMissingParam("currency")
	}
	if err := validatePaymentMethodTypes(req.PaymentMethodTypes); err != nil {
		return errInvalidParam("payment_method_types", err.Error())
	}
	return nil
}

func isCardPaymentMethod(token string) bool {
	_, ok := cardPaymentMethods[token]
	return ok
}

func isFPXPaymentMethod(token string) bool {
	_, ok := fpxPaymentMethods[token]
	return ok
}
