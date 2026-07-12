// Package webhook dispatches provider-style outgoing webhooks locally.
package webhook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/engine"
)

// PayloadVersion identifies a provider webhook shape.
type PayloadVersion string

const (
	// FawryV1 is the legacy Fawry webhook shape.
	FawryV1 PayloadVersion = "fawry-v1"
	// FawryV2 is the modern Fawry webhook shape.
	FawryV2 PayloadVersion = "fawry-v2"
)

// PayloadBuilder builds a provider webhook payload.
type PayloadBuilder interface {
	Build(ref string, status PaymentStatus) ([]byte, error)
}

// PaymentStatus represents the outcome of a simulated payment.
type PaymentStatus string

const (
	// PaymentStatusPaid means the payment succeeded.
	PaymentStatusPaid PaymentStatus = "PAID"
	// PaymentStatusUnpaid means the payment did not succeed.
	PaymentStatusUnpaid PaymentStatus = "UNPAID"
	// PaymentStatusNew means the payment is newly created.
	PaymentStatusNew PaymentStatus = "NEW"
)

// OrderItem represents a line item in a webhook payload.
type OrderItem struct {
	ItemCode string  `json:"itemCode"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// FawryV2Payload matches the Fawry V2 webhook notification shape.
// See .agents/project/ready/phase-1-cli-foundation/findings/01-fawry-integration-audit.md.
type FawryV2Payload struct {
	RequestID             string      `json:"requestId"`
	FawryRefNumber        string      `json:"fawryRefNumber"`
	MerchantRefNumber     string      `json:"merchantRefNumber"`
	CustomerMobile        string      `json:"customerMobile"`
	CustomerMail          string      `json:"customerMail"`
	CustomerMerchantID    string      `json:"customerMerchantId"`
	PaymentAmount         float64     `json:"paymentAmount"`
	OrderAmount           float64     `json:"orderAmount"`
	FawryFees             float64     `json:"fawryFees"`
	OrderStatus           string      `json:"orderStatus"`
	PaymentMethod         string      `json:"paymentMethod"`
	PaymentTime           int64       `json:"paymentTime"`
	MessageSignature      string      `json:"messageSignature"`
	PaymentRefrenceNumber string      `json:"paymentRefrenceNumber"`
	OrderExpiryDate       int64       `json:"orderExpiryDate"`
	OrderItems            []OrderItem `json:"orderItems"`
}

// FawryV2Builder builds Fawry V2 webhook payloads.
//
// Deprecated: Fawry-specific payload construction has moved to the Fawry
// provider (internal/fawry). This builder remains temporarily so the legacy
// dispatcher constructor keeps compiling; it will be removed once the
// dispatcher is wired directly to provider PayloadBuilders.
type FawryV2Builder struct {
	// Secret is a local dummy key used to compute the messageSignature approximation.
	Secret string
	// Signer computes the messageSignature.
	Signer Signer
	// Store reads transaction data from the ledger.
	Store engine.TransactionStore
}

// NewFawryV2Builder creates a Fawry V2 payload builder.
func NewFawryV2Builder(secret string, store engine.TransactionStore) *FawryV2Builder {
	return &FawryV2Builder{
		Secret: secret,
		Signer: NewHMACSigner(secret),
		Store:  store,
	}
}

// Build constructs a Fawry V2 webhook payload for the given reference and status.
func (b *FawryV2Builder) Build(ref string, status PaymentStatus) ([]byte, error) {
	tx, ok, err := b.Store.GetByReference(ref)
	if err != nil {
		return nil, fmt.Errorf("lookup transaction: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("transaction not found for ref %q", ref)
	}

	now := time.Now().UnixMilli()
	payload := FawryV2Payload{
		RequestID:             uuid.Must(uuid.NewRandom()).String(),
		FawryRefNumber:        "muara-fawry-ref",
		MerchantRefNumber:     ref,
		CustomerMobile:        "01000000000",
		CustomerMail:          "customer@example.com",
		CustomerMerchantID:    tx.CustomerRef,
		PaymentAmount:         tx.Amount,
		OrderAmount:           tx.Amount,
		FawryFees:             0,
		OrderStatus:           string(status),
		PaymentMethod:         "CARD",
		PaymentTime:           now,
		PaymentRefrenceNumber: "muara-payment-ref",
		OrderExpiryDate:       now + 3600000,
		OrderItems:            mapItems(tx.Items),
	}

	sig, err := b.Signer.Sign(payload)
	if err != nil {
		return nil, fmt.Errorf("sign payload: %w", err)
	}
	payload.MessageSignature = sig

	return json.Marshal(payload)
}

func mapItems(items []engine.TransactionItem) []OrderItem {
	result := make([]OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, OrderItem{
			ItemCode: item.ItemCode,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}
	return result
}

// PayloadFor returns a builder for the requested version.
//
// Deprecated: Provider selection will be handled by the provider registry in
// Phase 04. This factory is kept only for backward compatibility with the
// legacy dispatcher constructor.
func PayloadFor(version PayloadVersion, secret string, store engine.TransactionStore) (PayloadBuilder, error) {
	switch version {
	case FawryV2:
		return NewFawryV2Builder(secret, store), nil
	default:
		return nil, fmt.Errorf("unsupported webhook payload version: %q", version)
	}
}
