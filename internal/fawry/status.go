package fawry

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
)

// PaymentStatusResponse is the response shape for GET /fawry/payment-status.
type PaymentStatusResponse struct {
	Status    string  `json:"status"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

// NewStatusHandler returns GET /fawry/payment-status.
// It lets a client independently verify a transaction's current state using a
// signed query, matching the production best practice of status inquiry.
func NewStatusHandler(merchantCode, merchantSecurityKey string, store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		q := StatusQuery{
			MerchantCode:   r.URL.Query().Get("merchantCode"),
			MerchantRefNum: r.URL.Query().Get("merchantRefNum"),
			Signature:      r.URL.Query().Get("signature"),
		}

		if err := validateStatusQuery(q, merchantCode); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, err.Error())
			return
		}

		if !VerifyStatusQuery(q, merchantSecurityKey) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid signature")
			return
		}

		tx, ok, err := store.GetByReference(q.MerchantRefNum)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(PaymentStatusResponse{
			Status:    mapFawryStatus(string(tx.Status)),
			Reference: tx.Reference,
			Amount:    tx.Amount,
			Currency:  tx.Currency,
		})
	}
}

func validateStatusQuery(q StatusQuery, expectedMerchantCode string) error {
	if q.MerchantCode == "" {
		return errors.New("merchantCode is required")
	}
	if q.MerchantCode != expectedMerchantCode {
		return errors.New("invalid merchantCode")
	}
	if q.MerchantRefNum == "" {
		return errors.New("merchantRefNum is required")
	}
	if q.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

func mapFawryStatus(status string) string {
	switch engine.TransactionStatus(status) {
	case engine.TransactionStatusPaid:
		return "PAID"
	case engine.TransactionStatusUnpaid:
		return "UNPAID"
	case engine.TransactionStatusNew:
		return "NEW"
	default:
		return status
	}
}
