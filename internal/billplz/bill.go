package billplz

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

// NewBillsHandler returns the handler for POST /api/v3/bills.
func NewBillsHandler(
	apiKey, defaultCollectionID string,
	bills *billStore,
	collections *collectionStore,
	txStore engine.TransactionStore,
	baseURL string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		var req CreateBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.CollectionID == "" {
			req.CollectionID = defaultCollectionID
		}
		if err := validateCreateBillRequest(req, collections); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		b := bills.create(req, baseURL)
		recordTransaction(r.Context(), txStore, b)

		writeJSON(w, BillResponse{Bill: b})
	}
}

// NewBillHandler returns the handler for GET /api/v3/bills/{id}.
func NewBillHandler(apiKey string, bills *billStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		b, ok := bills.get(r.PathValue("id"))
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		writeJSON(w, BillResponse{Bill: b})
	}
}

// NewDeleteBillHandler returns the handler for DELETE /api/v3/bills/{id}.
func NewDeleteBillHandler(apiKey string, bills *billStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		b, ok := bills.delete(r.PathValue("id"))
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		writeJSON(w, BillResponse{Bill: b})
	}
}

func validateCreateBillRequest(req CreateBillRequest, collections *collectionStore) error {
	if req.CollectionID == "" {
		return errRequiredField("collection_id")
	}
	if _, ok := collections.get(req.CollectionID); !ok {
		return errCollectionNotFound
	}
	if req.Email == "" {
		return errRequiredField("email")
	}
	if req.Name == "" {
		return errRequiredField("name")
	}
	if req.Amount <= 0 {
		return errInvalidAmount
	}
	if req.CallbackURL == "" {
		return errRequiredField("callback_url")
	}
	if req.Description == "" {
		return errRequiredField("description")
	}
	return nil
}

func recordTransaction(ctx context.Context, txStore engine.TransactionStore, b Bill) {
	if txStore == nil {
		return
	}
	status := engine.TransactionStatusUnpaid
	if b.State == BillStatePaid {
		status = engine.TransactionStatusPaid
	}
	_, _, _ = txStore.CreateOrGet(engine.NewTransaction(engine.Transaction{
		ID:          b.ID,
		Provider:    ProviderName,
		Type:        "bill",
		Amount:      float64(b.Amount) / 100.0,
		Currency:    "MYR",
		Status:      status,
		CustomerRef: b.Email,
		Reference:   b.ID,
		TraceID:     httputil.TraceIDFromContext(ctx),
	}))
}

func updateTransaction(ctx context.Context, txStore engine.TransactionStore, b Bill) {
	if txStore == nil {
		return
	}
	stored, ok, err := txStore.GetByID(b.ID)
	if err != nil || !ok {
		recordTransaction(ctx, txStore, b)
		return
	}
	status := engine.TransactionStatusUnpaid
	if b.State == BillStatePaid {
		status = engine.TransactionStatusPaid
	}
	if stored.Status == status {
		return
	}
	stored.Status = status
	stored.UpdatedAt = time.Now()
	_, _, _ = txStore.CreateOrGet(stored)
}
