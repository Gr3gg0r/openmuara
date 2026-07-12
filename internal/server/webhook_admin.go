package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/webhook"
)

// WebhookAdminHandlers registers admin endpoints for inspecting and replaying webhooks.
// When dispatcher is nil, list/inspect/delete handlers are still registered against an
// empty in-memory store so the dashboard can render an empty webhook list. Replay
// endpoints are only registered when a dispatcher is available.
// dispatchers maps provider name to its own dispatcher so replays use the correct
// payload builder and target URL.
func WebhookAdminHandlers(mux *http.ServeMux, dispatcher *webhook.Dispatcher, dispatchers map[string]*webhook.Dispatcher) {
	var store webhook.AttemptStore
	if dispatcher != nil {
		store = dispatcher.Store
	} else {
		store = webhook.NewMemoryStore()
	}

	mux.HandleFunc("GET /_admin/webhooks", listWebhooksHandler(store))
	mux.HandleFunc("GET /_admin/webhooks/{ref}", inspectWebhookHandler(store))
	mux.HandleFunc("DELETE /_admin/webhooks/{ref}", deleteWebhookHandler(store))

	if dispatcher != nil {
		mux.HandleFunc("POST /_admin/webhooks/{ref}/replay", replayWebhookHandler(dispatcher, dispatchers))
		mux.HandleFunc("POST /_admin/webhooks/replay-all", replayAllWebhooksHandler(dispatcher, dispatchers))
	}
}

func listWebhooksHandler(store webhook.AttemptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := pageParams(r)
		attempts, err := store.List(limit, offset)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		statusFilter := r.URL.Query().Get("status")
		providerFilter := r.URL.Query().Get("provider")
		if statusFilter != "" || providerFilter != "" {
			filtered := make([]*webhook.Attempt, 0, len(attempts))
			for _, a := range attempts {
				if statusFilter != "" && string(a.Status) != statusFilter {
					continue
				}
				if providerFilter != "" && a.ProviderName != providerFilter {
					continue
				}
				filtered = append(filtered, a)
			}
			attempts = filtered
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"limit":   limit,
			"offset":  offset,
			"results": attempts,
		})
	}
}

func inspectWebhookHandler(store webhook.AttemptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.PathValue("ref")
		attempt, err := store.Get(ref)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if attempt == nil {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		var payload any = string(attempt.Payload)
		if json.Valid(attempt.Payload) {
			payload = json.RawMessage(attempt.Payload)
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"webhook": map[string]any{
				"ref":             attempt.Ref,
				"provider":        attempt.ProviderName,
				"provider_name":   attempt.ProviderName,
				"url":             attempt.URL,
				"payload":         payload,
				"headers":         redactHeaders(attempt.Headers),
				"signature_valid": attempt.SignatureValid,
				"status":          attempt.Status,
				"attempt_events":  attempt.History,
				"trace_id":        attempt.TraceID,
				"created_at":      attempt.CreatedAt,
				"updated_at":      attempt.UpdatedAt,
			},
		})
	}
}

func redactHeaders(headers map[string]string) map[string]string {
	redacted := make(map[string]string, len(headers))
	for k, v := range headers {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "signature") ||
			strings.Contains(lower, "authorization") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "secret") {
			redacted[k] = "***"
			continue
		}
		redacted[k] = v
	}
	return redacted
}

func replayWebhookHandler(dispatcher *webhook.Dispatcher, dispatchers map[string]*webhook.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		ref := r.PathValue("ref")
		attempt, err := dispatcher.Store.Get(ref)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if attempt == nil {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		d := dispatcherForAttempt(attempt, dispatcher, dispatchers)
		attempt, err = d.Replay(r.Context(), ref)
		if err != nil {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.webhook_replay", "webhook", ref, "", "accepted")
		respondJSON(w, http.StatusAccepted, attempt)
	}
}

func replayAllWebhooksHandler(dispatcher *webhook.Dispatcher, dispatchers map[string]*webhook.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		attempts, err := dispatcher.Store.List(1000, 0)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		statusFilter := r.URL.Query().Get("status")
		providerFilter := r.URL.Query().Get("provider")
		replayed := 0
		for _, a := range attempts {
			if statusFilter != "" && string(a.Status) != statusFilter {
				continue
			}
			if providerFilter != "" && a.ProviderName != providerFilter {
				continue
			}
			d := dispatcherForAttempt(a, dispatcher, dispatchers)
			if _, err := d.Replay(r.Context(), a.Ref); err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			replayed++
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.webhook_replay_all", "webhook", "*", "", fmt.Sprintf("replayed=%d", replayed))
		respondJSON(w, http.StatusAccepted, map[string]int{"replayed": replayed})
	}
}

// dispatcherForAttempt returns the dispatcher that matches the attempt's
// provider, falling back to the active dispatcher when no mapping exists.
func dispatcherForAttempt(a *webhook.Attempt, active *webhook.Dispatcher, dispatchers map[string]*webhook.Dispatcher) *webhook.Dispatcher {
	if a == nil || a.ProviderName == "" {
		return active
	}
	if d, ok := dispatchers[a.ProviderName]; ok && d != nil {
		return d
	}
	return active
}

func deleteWebhookHandler(store webhook.AttemptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !RequireAdmin(r.Context(), w, r) {
			return
		}
		// AttemptStore does not support deletion; clear the payload/headers to drop sensitive data.
		ref := r.PathValue("ref")
		attempt, err := store.Get(ref)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if attempt == nil {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		attempt.Payload = nil
		attempt.Headers = nil
		if err := store.Save(attempt); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
