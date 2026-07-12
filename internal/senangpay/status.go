package senangpay

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// QueryResponse is the response shape for GET /senangpay/query.
type QueryResponse struct {
	OrderID       string  `json:"order_id"`
	StatusID      string  `json:"status_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	TransactionID string  `json:"transaction_id,omitempty"`
}

// NewStatusHandler returns GET /senangpay/query.
// It lets a client independently verify a transaction's current state using a
// signed query, matching the production best practice of status inquiry.
func NewStatusHandler(secret string, store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		orderID := r.URL.Query().Get("order_id")
		hash := r.URL.Query().Get("hash")

		if err := validateStatusQuery(orderID, hash); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, err.Error())
			return
		}

		if !VerifyStatusQuery(orderID, hash, secret) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid hash")
			return
		}

		tx, ok, err := store.GetByReference(orderID)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "order not found")
			return
		}

		statusID, status := mapSenangPayStatus(tx.Status)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(QueryResponse{
			OrderID:  orderID,
			StatusID: statusID,
			Status:   status,
			Amount:   tx.Amount,
			Currency: tx.Currency,
		})
	}
}

func validateStatusQuery(orderID, hash string) error {
	if orderID == "" {
		return errors.New("order_id is required")
	}
	if hash == "" {
		return errors.New("hash is required")
	}
	return nil
}

func mapSenangPayStatus(status engine.TransactionStatus) (string, string) {
	switch status {
	case engine.TransactionStatusPaid:
		return "1", "paid"
	default:
		return "0", "unpaid"
	}
}
