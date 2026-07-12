// Package api provides provider-agnostic HTTP endpoints for the OpenMuara ledger.
package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/google/uuid"
)

// PaymentRequest creates a generic ledger transaction.
type PaymentRequest struct {
	Provider       string                   `json:"provider"`
	Type           string                   `json:"type"`
	Amount         float64                  `json:"amount"`
	Currency       string                   `json:"currency"`
	CustomerRef    string                   `json:"customer_ref"`
	Reference      string                   `json:"reference"`
	IdempotencyKey string                   `json:"idempotency_key"`
	Items          []engine.TransactionItem `json:"items,omitempty"`
}

// PaymentResponse is the unified payment response shape.
type PaymentResponse struct {
	ID          string                   `json:"id"`
	Provider    string                   `json:"provider"`
	Type        string                   `json:"type"`
	Amount      float64                  `json:"amount"`
	Currency    string                   `json:"currency"`
	Status      engine.TransactionStatus `json:"status"`
	CustomerRef string                   `json:"customer_ref"`
	Reference   string                   `json:"reference"`
	Items       []engine.TransactionItem `json:"items,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// NewPayHandler returns POST /v1/pay. If allowedProviders is non-empty, the
// provider field must match one of the supplied names.
func NewPayHandler(store engine.TransactionStore, allowedProviders []string) http.HandlerFunc {
	allowed := make(map[string]bool, len(allowedProviders))
	for _, name := range allowedProviders {
		allowed[name] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := validatePaymentRequest(&req, allowed); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if req.IdempotencyKey == "" {
			req.IdempotencyKey = r.Header.Get("Idempotency-Key")
		}

		tx := engine.Transaction{
			Provider:       req.Provider,
			Type:           req.Type,
			Amount:         req.Amount,
			Currency:       req.Currency,
			Status:         engine.TransactionStatusNew,
			CustomerRef:    req.CustomerRef,
			Reference:      req.Reference,
			IdempotencyKey: req.IdempotencyKey,
			TraceID:        httputil.TraceIDFromContext(r.Context()),
			Items:          req.Items,
		}

		stored, created, err := store.CreateOrGet(tx)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to record transaction")
			return
		}

		status := http.StatusCreated
		if !created {
			status = http.StatusOK
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(toPaymentResponse(stored))
	}
}

// NewGetPaymentHandler returns GET /v1/pay/{ref}.
func NewGetPaymentHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ref := strings.TrimPrefix(r.URL.Path, "/v1/pay/")
		if ref == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "reference is required")
			return
		}

		tx, ok, err := store.GetByReference(ref)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(toPaymentResponse(tx))
	}
}

// NewRefundHandler returns POST /v1/refund/{ref}.
func NewRefundHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ref := strings.TrimPrefix(r.URL.Path, "/v1/refund/")
		if ref == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "reference is required")
			return
		}

		tx, ok, err := store.GetByReference(ref)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		if err := engine.Transition(&tx, engine.TransactionStatusRefunded); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, err.Error())
			return
		}

		stored, _, err := store.CreateOrGet(tx)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to record refund")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(toPaymentResponse(stored))
	}
}

var currencyRe = regexp.MustCompile(`^[A-Z]{3}$`)

func validatePaymentRequest(req *PaymentRequest, allowed map[string]bool) error {
	if req.Provider == "" {
		return errcode.New(errcode.EInvalidRequest, "provider is required")
	}
	if len(allowed) > 0 && !allowed[req.Provider] {
		return errcode.New(errcode.EUnknownProvider, "provider is not supported")
	}
	if req.Type == "" {
		return errcode.New(errcode.EInvalidRequest, "type is required")
	}
	if req.Amount <= 0 {
		return errcode.New(errcode.EInvalidRequest, "amount must be greater than zero")
	}
	if req.Currency == "" {
		return errcode.New(errcode.EInvalidRequest, "currency is required")
	}
	if !currencyRe.MatchString(req.Currency) {
		return errcode.New(errcode.EInvalidRequest, "currency must be a 3-letter ISO code")
	}
	if req.Reference == "" {
		return errcode.New(errcode.EInvalidRequest, "reference is required")
	}
	return nil
}

func toPaymentResponse(tx engine.Transaction) PaymentResponse {
	return PaymentResponse{
		ID:          tx.ID,
		Provider:    tx.Provider,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Currency:    tx.Currency,
		Status:      tx.Status,
		CustomerRef: tx.CustomerRef,
		Reference:   tx.Reference,
		Items:       tx.Items,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}
}

// GenerateID returns a UUID string for tests or callers that need a reference ID.
func GenerateID() string { return uuid.NewString() }
