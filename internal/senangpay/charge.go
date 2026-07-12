package senangpay

import (
	"encoding/json"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

// ChargeRequest is the SenangPay-style charge request.
type ChargeRequest struct {
	Detail  string  `json:"detail"`
	Amount  float64 `json:"amount"`
	OrderID string  `json:"order_id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   string  `json:"phone"`
	Hash    string  `json:"hash"`
}

// ChargeResponse is the SenangPay-style charge response.
type ChargeResponse struct {
	OrderID    string `json:"order_id"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
	Reference  string `json:"reference"`
}

// NewChargeHandler returns POST /senangpay/charge.
func NewChargeHandler(secret string, store engine.TransactionStore, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req ChargeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := validateChargeRequest(req); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if !Verify(req, secret) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid hash")
			return
		}

		tx := engine.NewTransaction(engine.Transaction{
			Provider:    ProviderName,
			Type:        "charge",
			Amount:      req.Amount,
			Currency:    "MYR",
			Status:      engine.TransactionStatusNew,
			CustomerRef: req.Email,
			Reference:   req.OrderID,
			TraceID:     httputil.TraceIDFromContext(r.Context()),
		})
		if _, _, err := store.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to record transaction")
			return
		}

		resp := ChargeResponse{
			OrderID:    req.OrderID,
			PaymentURL: baseURL + "/_admin/senangpay-escape?order_id=" + req.OrderID,
			Status:     "ok",
			Reference:  req.OrderID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func validateChargeRequest(req ChargeRequest) error {
	if req.Detail == "" {
		return errcode.New(errcode.EInvalidRequest, "detail is required")
	}
	if req.Amount <= 0 {
		return errcode.New(errcode.EInvalidRequest, "amount must be greater than zero")
	}
	if req.OrderID == "" {
		return errcode.New(errcode.EInvalidRequest, "order_id is required")
	}
	if req.Hash == "" {
		return errcode.New(errcode.ESignatureMissing, "hash is required")
	}
	return nil
}

// SignRequest is a helper for tests to compute a valid hash.
func SignRequest(req *ChargeRequest, secret string) {
	req.Hash = Sign(secret, req.Detail, req.Amount, req.OrderID)
}
