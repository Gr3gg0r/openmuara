package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

// NewWebhookHandler returns a handler that accepts Fawry V2-style callbacks.
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

		var payload webhook.FawryV2Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if webhookSecret != "" {
			sig := payload.MessageSignature
			valid, err := webhook.NewHMACSigner(webhookSecret).Verify(payload, sig)
			if err != nil || !valid {
				httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, "invalid message signature")
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

// NewPayloadBuilder returns a builder that produces signed Fawry V2 payloads.
func NewPayloadBuilder(webhookSecret string, store engine.TransactionStore) func(context.Context, provider.Transaction) ([]byte, error) {
	return func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		if store == nil {
			return nil, errcode.New(errcode.EInternal, "store not configured")
		}
		stored, ok, err := store.GetByReference(tx.Reference)
		if err != nil {
			return nil, errcode.Wrap(errcode.EInternal, "lookup transaction", err)
		}
		if !ok {
			return nil, errcode.New(errcode.ETransactionNotFound, fmt.Sprintf("transaction not found for ref %q", tx.Reference))
		}

		amount := stored.Amount
		if amount == 0 && tx.Amount != 0 {
			amount = tx.Amount
		}

		status := tx.Status
		if status == "" {
			status = string(stored.Status)
		}

		now := time.Now().UnixMilli()
		payload := webhook.FawryV2Payload{
			RequestID:             uuid.Must(uuid.NewRandom()).String(),
			FawryRefNumber:        "muara-fawry-ref",
			MerchantRefNumber:     tx.Reference,
			CustomerMobile:        "01000000000",
			CustomerMail:          "customer@example.com",
			CustomerMerchantID:    stored.CustomerRef,
			PaymentAmount:         amount,
			OrderAmount:           amount,
			FawryFees:             0,
			OrderStatus:           status,
			PaymentMethod:         "CARD",
			PaymentTime:           now,
			PaymentRefrenceNumber: "muara-payment-ref",
			OrderExpiryDate:       now + 3600000,
			OrderItems:            mapItems(stored.Items),
		}

		sig, err := webhook.NewHMACSigner(webhookSecret).Sign(payload)
		if err != nil {
			return nil, errcode.Wrap(errcode.EInternal, "sign payload", err)
		}
		payload.MessageSignature = sig

		return json.Marshal(payload)
	}
}

func mapItems(items []engine.TransactionItem) []webhook.OrderItem {
	result := make([]webhook.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, webhook.OrderItem{
			ItemCode: item.ItemCode,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}
	return result
}
