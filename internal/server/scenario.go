package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
)

// ScenarioAdminHandlers registers the scenario simulation endpoints.
func ScenarioAdminHandlers(mux *http.ServeMux, store engine.TransactionStore) {
	mux.HandleFunc("POST /_admin/scenario/", newScenarioHandler(store))
}

func newScenarioHandler(store engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		outcome := strings.TrimPrefix(r.URL.Path, "/_admin/scenario/")
		if outcome == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "scenario outcome is required")
			return
		}

		switch outcome {
		case "success", "fail", "timeout":
			// valid
		default:
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "unknown scenario outcome")
			return
		}

		ref := r.URL.Query().Get("ref")
		if ref == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "ref query parameter is required")
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

		var targetStatus engine.TransactionStatus
		switch outcome {
		case "success":
			targetStatus = engine.TransactionStatusPaid
		case "fail", "timeout":
			targetStatus = engine.TransactionStatusUnpaid
		}

		if outcome == "timeout" {
			time.Sleep(2 * time.Second)
		}

		if err := engine.Transition(&tx, targetStatus); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, err.Error())
			return
		}

		stored, _, err := store.CreateOrGet(tx)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to update transaction")
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.scenario", "transaction", stored.Reference, outcome, string(stored.Status))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"reference": stored.Reference,
			"status":    stored.Status,
			"scenario":  outcome,
		})
	}
}
