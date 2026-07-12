package stripe

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/webhook"
)

// Event mirrors a subset of the Stripe event object.
type Event struct {
	ID     string    `json:"id"`
	Object string    `json:"object"`
	Type   string    `json:"type"`
	Data   EventData `json:"data"`
}

// EventData wraps the object inside a Stripe event.
type EventData struct {
	Object *CheckoutSession `json:"object"`
}

// Dispatcher is the minimal interface needed to dispatch a webhook.
type Dispatcher interface {
	Dispatch(ctx context.Context, ref string, status webhook.PaymentStatus) (*webhook.Attempt, error)
}

// NewSuccessSimulationHandler returns POST /_admin/stripe/success.
// It marks the session complete, updates the ledger, and dispatches a webhook.
func NewSuccessSimulationHandler(sessions SessionStore, ledger engine.TransactionStore, dispatcher Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			sessionID = r.FormValue("session_id")
		}
		if sessionID == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "session_id is required")
			return
		}

		session, ok := sessions.Load(sessionID)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "session not found")
			return
		}

		session.Status = "complete"
		session.PaymentStatus = "paid"
		sessions.Save(sessionID, session)

		if tx, ok, err := ledger.GetByReference(sessionID); err != nil {
			slog.Warn("failed to lookup transaction", "session_id", sessionID, "error", err)
		} else if ok {
			if err := engine.Transition(&tx, engine.TransactionStatusPaid); err != nil {
				httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, err.Error())
				return
			}
			if _, _, err := ledger.CreateOrGet(tx); err != nil {
				slog.Warn("failed to update transaction", "session_id", sessionID, "error", err)
			}
		}

		if dispatcher != nil {
			if _, err := dispatcher.Dispatch(r.Context(), sessionID, webhook.PaymentStatusPaid); err != nil {
				slog.Warn("failed to dispatch stripe webhook", "session_id", sessionID, "error", err)
			}
		} else {
			slog.Warn("no webhook dispatcher configured; stripe webhook not dispatched", "session_id", sessionID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"session_id": sessionID,
			"status":     "complete",
		})
	}
}
