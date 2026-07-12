// Package fawry emulates the Fawry Express Checkout payment gateway.
package fawry

import (
	"encoding/json"
	"net/http"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

// NewChargeHandler returns a handler that validates Fawry charge requests and
// records them as transactions in the shared ledger.
func NewChargeHandler(merchantSecurityKey string, store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req ChargeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if err := validateChargeRequest(req); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if !Verify(req, merchantSecurityKey) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid signature")
			return
		}

		tx := buildTransaction(r, req)
		if _, _, err := store.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to record transaction")
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "charge.created", "transaction", req.MerchantRefNum, audit.JSON(req), "ok")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ChargeResponse{
			Status:    "ok",
			Reference: req.MerchantRefNum,
		})
	}
}

func buildTransaction(r *http.Request, req ChargeRequest) engine.Transaction {
	var amount float64
	items := make([]engine.TransactionItem, 0, len(req.ChargeItems))
	for _, item := range req.ChargeItems {
		amount += item.Price * float64(item.Quantity)
		items = append(items, engine.TransactionItem{
			ItemCode: item.ItemID,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	return engine.NewTransaction(engine.Transaction{
		Provider:       "fawry",
		Type:           "charge",
		Amount:         amount,
		Currency:       "EGP",
		Status:         engine.TransactionStatusNew,
		CustomerRef:    req.CustomerProfileID,
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
		Reference:      req.MerchantRefNum,
		TraceID:        httputil.TraceIDFromContext(r.Context()),
		Items:          items,
	})
}

func validateChargeRequest(req ChargeRequest) error {
	if req.MerchantCode == "" {
		return errcode.New(errcode.EInvalidRequest, "merchantCode is required")
	}
	if req.MerchantRefNum == "" {
		return errcode.New(errcode.EInvalidRequest, "merchantRefNum is required")
	}
	if len(req.ChargeItems) == 0 {
		return errcode.New(errcode.EInvalidRequest, "chargeItems must contain at least one item")
	}
	for _, item := range req.ChargeItems {
		if item.ItemID == "" {
			return errcode.New(errcode.EInvalidRequest, "chargeItems.itemId is required")
		}
	}
	if req.ReturnURL == "" {
		return errcode.New(errcode.EInvalidRequest, "returnUrl is required")
	}
	if req.Signature == "" {
		return errcode.New(errcode.ESignatureMissing, "signature is required")
	}
	return nil
}

// HandlerDependencies is a convenience type used during setup.
type HandlerDependencies struct {
	MerchantSecurityKey string
	WebhookSecret       string
}

// ChargeHandlerPath is the HTTP path for charge requests.
const ChargeHandlerPath = "/fawry/charge"
