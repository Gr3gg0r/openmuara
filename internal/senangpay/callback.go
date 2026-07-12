package senangpay

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

// ErrOrderNotFound is returned when a callback references an unknown order.
var ErrOrderNotFound = errors.New("senangpay: order not found")

// CallbackQuery represents SenangPay callback/webhook query parameters.
type CallbackQuery struct {
	StatusID      string `json:"status_id"`
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Msg           string `json:"msg"`
}

// NewCallbackHandler returns GET /senangpay/callback.
func NewCallbackHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		q := parseCallback(r)
		if err := applyCallback(store, q); err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, errcode.Message(err))
				return
			}
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, errcode.Message(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"order_id": q.OrderID,
			"status":   q.StatusID,
		})
	}
}

// NewWebhookHandler returns POST /senangpay/webhook.
func NewWebhookHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		q := parseCallback(r)
		if err := applyCallback(store, q); err != nil {
			if errors.Is(err, ErrOrderNotFound) {
				httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, errcode.Message(err))
				return
			}
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, errcode.Message(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func parseCallback(r *http.Request) CallbackQuery {
	return CallbackQuery{
		StatusID:      r.URL.Query().Get("status_id"),
		OrderID:       r.URL.Query().Get("order_id"),
		TransactionID: r.URL.Query().Get("transaction_id"),
		Msg:           r.URL.Query().Get("msg"),
	}
}

func applyCallback(store engine.TransactionStore, q CallbackQuery) error {
	if q.OrderID == "" {
		return errcode.Wrap(errcode.ETransactionNotFound, "senangpay: order not found", ErrOrderNotFound)
	}

	tx, ok, err := store.GetByReference(q.OrderID)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "lookup transaction", err)
	}
	if !ok {
		return errcode.Wrap(errcode.ETransactionNotFound, "senangpay: order not found", ErrOrderNotFound)
	}

	var targetStatus engine.TransactionStatus
	switch q.StatusID {
	case "1":
		targetStatus = engine.TransactionStatusPaid
	case "0":
		targetStatus = engine.TransactionStatusUnpaid
	default:
		return nil
	}

	if err := engine.Transition(&tx, targetStatus); err != nil {
		return errcode.Wrap(errcode.ETransactionTransitionInvalid, "invalid transaction transition", err)
	}

	_, _, err = store.CreateOrGet(tx)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "failed to update transaction", err)
	}
	return nil
}
