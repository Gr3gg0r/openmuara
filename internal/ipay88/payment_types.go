package ipay88

// PaymentRequest stores the original iPay88 ePayment submission keyed by RefNo.
type PaymentRequest struct {
	MerchantCode      string
	PaymentID         string
	RefNo             string
	Amount            string
	Currency          string
	ProdDesc          string
	UserName          string
	UserEmail         string
	UserContact       string
	Remark            string
	Lang              string
	Signature         string
	SignatureType     string
	ResponseURL       string
	BackendURL        string
	SelectedPaymentID string
	Status            string
}

// PaymentStatus represents iPay88 callback status values.
type PaymentStatus string

const (
	// PaymentStatusSuccess is iPay88's success status.
	PaymentStatusSuccess PaymentStatus = "1"
	// PaymentStatusFailure is iPay88's failure status.
	PaymentStatusFailure PaymentStatus = "0"
	// PaymentStatusPending is iPay88's pending status.
	PaymentStatusPending PaymentStatus = "6"
)

// RequeryStatus represents iPay88 requery status values.
type RequeryStatus string

const (
	// RequeryStatusSuccess is returned when the payment succeeded.
	RequeryStatusSuccess RequeryStatus = "00"
	// RequeryStatusFailure is returned when the payment failed or is unknown.
	RequeryStatusFailure RequeryStatus = "01"
)
