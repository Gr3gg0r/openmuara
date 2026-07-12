package webhook

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// NewTestReceiverHandler returns an HTTP handler that accepts any JSON body and logs it
// as a synthetic webhook attempt. This lets users test dispatch without running their own app.
func NewTestReceiverHandler() http.HandlerFunc {
	store := NewMemoryStore()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		attempt := &Attempt{
			ID:      uuid.Must(uuid.NewRandom()).String(),
			Ref:     "test-" + uuid.Must(uuid.NewRandom()).String(),
			URL:     r.URL.String(),
			Status:  AttemptStatusDelivered,
			Payload: body,
		}

		if err := store.Save(attempt); err != nil {
			slog.Error("failed to save test receiver attempt", "error", err)
			http.Error(w, "failed to save attempt", http.StatusInternalServerError)
			return
		}

		slog.Info("test receiver accepted webhook",
			"ref", attempt.Ref,
			"size", len(body),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}
