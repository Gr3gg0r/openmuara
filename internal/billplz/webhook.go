package billplz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

// PayloadHeaders returns the Content-Type header for Billplz callbacks.
func (p *Provider) PayloadHeaders(_ context.Context, _ provider.Transaction) (map[string]string, error) {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, nil
}

// PayloadBuilder returns a form-urlencoded Bill object with x_signature.
func (p *Provider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return p.buildPayload
}

func (p *Provider) buildPayload(_ context.Context, tx provider.Transaction) ([]byte, error) {
	b, ok := p.bills.get(tx.Reference)
	if !ok {
		return nil, errcode.New(errcode.ETransactionNotFound, fmt.Sprintf("bill not found for ref %q", tx.Reference))
	}
	return buildCallbackPayload(b, p.xSignatureKey), nil
}

// NewWebhookHandler returns the handler for POST /billplz/webhook.
func NewWebhookHandler(xSignatureKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form payload")
			return
		}

		values := make(map[string]string)
		for k, v := range r.Form {
			if len(v) > 0 {
				values[k] = v[0]
			}
		}
		if !VerifyCallback(values, xSignatureKey) {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, "invalid x_signature")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// VerifyCallback checks the x_signature in a form-urlencoded callback payload.
func VerifyCallback(values map[string]string, xSignatureKey string) bool {
	sig, ok := values["x_signature"]
	if !ok {
		return false
	}
	delete(values, "x_signature")
	return Verify(values, xSignatureKey, sig)
}

func buildCallbackPayload(b Bill, key string) []byte {
	values := billToValues(b)
	sig := Sign(values, key)
	values["x_signature"] = sig
	return []byte(urlValues(values).Encode())
}

func billToValues(b Bill) map[string]string {
	values := map[string]string{
		"id":            b.ID,
		"collection_id": b.CollectionID,
		"paid":          strconv.FormatBool(b.Paid),
		"state":         string(b.State),
		"amount":        strconv.FormatInt(b.Amount, 10),
		"description":   b.Description,
		"name":          b.Name,
		"email":         b.Email,
		"callback_url":  b.CallbackURL,
		"url":           b.URL,
	}
	addIfNonEmpty(values, "mobile", b.Mobile)
	addIfNonEmpty(values, "reference_1", b.Reference1)
	addIfNonEmpty(values, "reference_1_label", b.Reference1Label)
	addIfNonEmpty(values, "reference_2", b.Reference2)
	addIfNonEmpty(values, "reference_2_label", b.Reference2Label)
	addIfNonEmpty(values, "redirect_url", b.RedirectURL)
	if b.PaidAmount != nil {
		values["paid_amount"] = strconv.FormatInt(*b.PaidAmount, 10)
	}
	if b.DueAt != nil {
		values["due_at"] = b.DueAt.Format(time.RFC3339)
	}
	if b.PaidAt != nil {
		values["paid_at"] = b.PaidAt.Format(time.RFC3339)
	}
	return values
}

func addIfNonEmpty(m map[string]string, key, value string) {
	if value != "" {
		m[key] = value
	}
}

func urlValues(m map[string]string) url.Values {
	v := make(url.Values, len(m))
	for k, val := range m {
		v.Set(k, val)
	}
	return v
}
