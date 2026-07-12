package v1

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
)

// WebhookBody is the legacy Fawry v1 notification shape.
type WebhookBody struct {
	MerchantRefNumber string `json:"merchantRefNumber"`
	OrderStatus       string `json:"orderStatus"`
	MessageSignature  string `json:"messageSignature"`
}

// NewWebhookHandler returns a handler that accepts Fawry v1-style callbacks.
func NewWebhookHandler(webhookSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		token := r.URL.Query().Get("token")
		if token != webhookSecret {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, "invalid token")
			return
		}

		var payload WebhookBody
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if webhookSecret != "" && !VerifySignature(payload, webhookSecret) {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, "invalid message signature")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// NewPayloadBuilder returns a builder that produces signed v1 webhook payloads.
func NewPayloadBuilder(webhookSecret string) func(context.Context, provider.Transaction) ([]byte, error) {
	return func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		body := WebhookBody{
			MerchantRefNumber: tx.Reference,
			OrderStatus:       tx.Status,
		}
		body.MessageSignature = sign(body, webhookSecret)
		return json.Marshal(body)
	}
}

func sign(body WebhookBody, secret string) string {
	text := body.MerchantRefNumber + body.OrderStatus
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature checks whether the messageSignature on a v1 webhook body matches
// the secret.
func VerifySignature(body WebhookBody, secret string) bool {
	return hmac.Equal([]byte(body.MessageSignature), []byte(sign(body, secret)))
}
