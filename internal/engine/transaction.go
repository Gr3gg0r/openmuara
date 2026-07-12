// Package engine provides the in-memory transaction ledger for OpenMuara.
package engine

import (
	"errors"
	"fmt"
	"time"
)

// TransactionStatus represents the lifecycle state of a transaction.
type TransactionStatus string

const (
	// TransactionStatusNew means the transaction was just created.
	TransactionStatusNew TransactionStatus = "new"
	// TransactionStatusPaid means the payment succeeded.
	TransactionStatusPaid TransactionStatus = "paid"
	// TransactionStatusUnpaid means the payment did not succeed.
	TransactionStatusUnpaid TransactionStatus = "unpaid"
	// TransactionStatusRefunded means the payment was refunded.
	TransactionStatusRefunded TransactionStatus = "refunded"
)

// TransactionItem represents a single line item within a transaction.
type TransactionItem struct {
	ItemCode string  `json:"itemCode"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// Transaction is the single source of truth for a payment operation.
// It is intentionally provider-agnostic so that Stripe, Fawry, RevenueCat,
// and future providers can share the same ledger shape.
type Transaction struct {
	ID             string            `json:"id"`
	Provider       string            `json:"provider"`
	Type           string            `json:"type"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	Status         TransactionStatus `json:"status"`
	CustomerRef    string            `json:"customerRef"`
	IdempotencyKey string            `json:"idempotencyKey,omitempty"`
	Reference      string            `json:"reference"`
	TraceID        string            `json:"trace_id,omitempty"`
	Items          []TransactionItem `json:"items,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// NewTransaction returns a transaction with CreatedAt and UpdatedAt set to now
// if they are zero. The caller is responsible for setting ID and other fields.
func NewTransaction(tx Transaction) Transaction {
	now := time.Now()
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = now
	}
	if tx.UpdatedAt.IsZero() {
		tx.UpdatedAt = now
	}
	return tx
}

// ErrInvalidTransition is returned when a transaction cannot move from one
// status to another under the supported state machine.
var ErrInvalidTransition = errors.New("invalid transaction status transition")

// validTransitions maps a source status to the set of statuses it may move to.
var validTransitions = map[TransactionStatus]map[TransactionStatus]bool{
	TransactionStatusNew:      {TransactionStatusPaid: true, TransactionStatusUnpaid: true},
	TransactionStatusPaid:     {TransactionStatusRefunded: true},
	TransactionStatusUnpaid:   {},
	TransactionStatusRefunded: {},
}

// CanTransition reports whether moving from status from to status to is allowed.
func CanTransition(from, to TransactionStatus) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	return allowed[to]
}

// Transition updates tx.Status to to if the transition is valid.
// It also refreshes tx.UpdatedAt when the status changes.
func Transition(tx *Transaction, to TransactionStatus) error {
	if tx.Status == to {
		return nil
	}
	if !CanTransition(tx.Status, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, tx.Status, to)
	}
	tx.Status = to
	tx.UpdatedAt = time.Now()
	return nil
}
