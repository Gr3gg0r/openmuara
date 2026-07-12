package server

import (
	"net/http"

	"github.com/openmuara/openmuara/internal/audit"
)

// AuditAdminHandlers registers the audit log read endpoint.
func AuditAdminHandlers(mux *http.ServeMux, store audit.Store) {
	if store == nil {
		return
	}
	mux.HandleFunc("GET /_admin/audit", listAuditHandler(store))
}

func listAuditHandler(store audit.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := pageParams(r)
		events, err := store.List(limit, offset)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"limit":   limit,
			"offset":  offset,
			"results": events,
		})
	}
}
