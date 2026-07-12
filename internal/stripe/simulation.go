package stripe

import (
	"encoding/json"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// simulationHandler updates a session and the ledger to a simulated outcome.
type simulationHandler struct {
	sessions      SessionStore
	ledger        engine.TransactionStore
	outcome       string
	status        string
	paymentStatus string
}

func (h *simulationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	session, ok := h.sessions.Load(sessionID)
	if !ok {
		httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "session not found")
		return
	}

	session.Status = h.status
	session.PaymentStatus = h.paymentStatus
	h.sessions.Save(sessionID, session)

	if tx, ok, err := h.ledger.GetByReference(sessionID); err != nil {
		httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
		return
	} else if ok {
		if err := engine.Transition(&tx, engine.TransactionStatusUnpaid); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, err.Error())
			return
		}
		if _, _, err := h.ledger.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to update transaction")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"session_id": sessionID,
		"status":     h.outcome,
	})
}

// NewFailureSimulationHandler returns POST /_admin/stripe/fail.
func NewFailureSimulationHandler(sessions SessionStore, ledger engine.TransactionStore) http.Handler {
	return &simulationHandler{
		sessions:      sessions,
		ledger:        ledger,
		outcome:       "failed",
		status:        "open",
		paymentStatus: "unpaid",
	}
}

// NewCancelSimulationHandler returns POST /_admin/stripe/cancel.
func NewCancelSimulationHandler(sessions SessionStore, ledger engine.TransactionStore) http.Handler {
	return &simulationHandler{
		sessions:      sessions,
		ledger:        ledger,
		outcome:       "canceled",
		status:        "canceled",
		paymentStatus: "unpaid",
	}
}
