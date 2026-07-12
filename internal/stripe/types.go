package stripe

import "strings"

// Clone returns a deep copy of the session.
func (s *CheckoutSession) Clone() *CheckoutSession {
	metadata := make(map[string]string, len(s.Metadata))
	for k, v := range s.Metadata {
		metadata[k] = v
	}
	paymentMethodTypes := make([]string, len(s.PaymentMethodTypes))
	copy(paymentMethodTypes, s.PaymentMethodTypes)

	return &CheckoutSession{
		ID:                 s.ID,
		Object:             s.Object,
		AmountTotal:        s.AmountTotal,
		Currency:           s.Currency,
		CustomerEmail:      s.CustomerEmail,
		Mode:               s.Mode,
		PaymentMethodTypes: paymentMethodTypes,
		PaymentStatus:      s.PaymentStatus,
		Status:             s.Status,
		SuccessURL:         s.SuccessURL,
		CancelURL:          s.CancelURL,
		URL:                s.URL,
		ClientReferenceID:  s.ClientReferenceID,
		Metadata:           metadata,
	}
}

// CheckoutSession mirrors a subset of the Stripe Checkout Session object.
type CheckoutSession struct {
	ID                 string            `json:"id"`
	Object             string            `json:"object"`
	AmountTotal        int64             `json:"amount_total"`
	Currency           string            `json:"currency"`
	CustomerEmail      string            `json:"customer_email,omitempty"`
	Mode               string            `json:"mode"`
	PaymentMethodTypes []string          `json:"payment_method_types"`
	PaymentStatus      string            `json:"payment_status"`
	Status             string            `json:"status"`
	SuccessURL         string            `json:"success_url"`
	CancelURL          string            `json:"cancel_url,omitempty"`
	URL                string            `json:"url"`
	ClientReferenceID  string            `json:"client_reference_id,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

// CreateCheckoutSessionRequest mirrors a subset of Stripe's create parameters.
type CreateCheckoutSessionRequest struct {
	SuccessURL         string            `json:"success_url"`
	CancelURL          string            `json:"cancel_url,omitempty"`
	LineItems          []LineItem        `json:"line_items"`
	Mode               string            `json:"mode"`
	PaymentMethodTypes []string          `json:"payment_method_types"`
	CustomerEmail      string            `json:"customer_email,omitempty"`
	ClientReferenceID  string            `json:"client_reference_id,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

// LineItem is a simplified line item.
type LineItem struct {
	PriceData *PriceData `json:"price_data"`
	Quantity  int64      `json:"quantity"`
}

// PriceData is a simplified price_data object.
type PriceData struct {
	Currency    string `json:"currency"`
	UnitAmount  int64  `json:"unit_amount"`
	ProductData struct {
		Name string `json:"name"`
	} `json:"product_data"`
}

// FPXBank represents a selectable Malaysian bank.
type FPXBank struct {
	Code string
	Name string
}

// FPXBanks is the supported Malaysian bank list for FPX emulation.
var FPXBanks = []FPXBank{
	{Code: "maybank2u", Name: "Maybank2U"},
	{Code: "cimb", Name: "CIMB Clicks"},
	{Code: "public_bank", Name: "Public Bank"},
	{Code: "rhb", Name: "RHB Now"},
	{Code: "hong_leong", Name: "Hong Leong Connect"},
	{Code: "ambank", Name: "AmBank"},
	{Code: "bank_islam", Name: "Bank Islam"},
	{Code: "affin_bank", Name: "Affin Bank"},
}

// normalizePaymentMethodTypes lowercases types and defaults to ["card"].
func normalizePaymentMethodTypes(types []string) []string {
	if len(types) == 0 {
		return []string{"card"}
	}
	out := make([]string, len(types))
	for i, t := range types {
		out[i] = strings.ToLower(t)
	}
	return out
}

func validatePaymentMethodTypes(types []string) error {
	if len(types) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(types))
	for _, t := range types {
		set[strings.ToLower(t)] = struct{}{}
	}
	if len(set) > 2 {
		return errInvalidPaymentMethodTypes
	}
	allowed := map[string]struct{}{"card": {}, "fpx": {}}
	for t := range set {
		if _, ok := allowed[t]; !ok {
			return errUnsupportedPaymentMethodType(t)
		}
	}
	return nil
}
