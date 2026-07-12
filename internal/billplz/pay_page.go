package billplz

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/ui"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

// NewPayPageHandler returns GET /_admin/billplz/pay/{id}.
func NewPayPageHandler(bills *billStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		b, ok := bills.get(r.PathValue("id"))
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		if b.State == BillStateDeleted {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, "bill is deleted")
			return
		}

		tok, _ := httputil.CSRFTokenFromContext(r.Context())
		data := ui.BillplzPayPageData{
			ID:          b.ID,
			Amount:      b.Amount,
			Description: b.Description,
			Methods: []ui.BillplzPaymentMethod{
				{Code: "fpx", Name: "FPX Online Banking"},
				{Code: "mpgs", Name: "Credit / Debit Card"},
				{Code: "boost", Name: "Boost Wallet"},
				{Code: "touchngo", Name: "Touch 'n Go eWallet"},
			},
			CSRFToken: tok,
		}
		if err := ui.ServeBillplzPayPage(w, data); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to render page")
		}
	}
}

// NewPayActionHandler returns POST /_admin/billplz/pay/{id}.
func NewPayActionHandler(
	bills *billStore,
	txStore engine.TransactionStore,
	dispatcher *webhook.Dispatcher,
	xSignatureKey string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		id := r.PathValue("id")
		b, ok := bills.get(id)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}

		outcome := r.FormValue("outcome")
		if outcome != "pay" && outcome != "cancel" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "outcome must be pay or cancel")
			return
		}

		if outcome == "pay" {
			b = markPaid(b)
		} else {
			b = markUnpaid(b)
		}
		bills.update(b)
		updateTransaction(r.Context(), txStore, b)
		dispatchCallback(r.Context(), dispatcher, b)

		if b.RedirectURL != "" {
			http.Redirect(w, r, redirectURL(b, xSignatureKey), http.StatusSeeOther)
			return
		}
		writeJSON(w, BillResponse{Bill: b})
	}
}

func markPaid(b Bill) Bill {
	b.Paid = true
	b.State = BillStatePaid
	now := time.Now()
	b.PaidAt = &now
	paidAmount := b.Amount
	b.PaidAmount = &paidAmount
	return b
}

func markUnpaid(b Bill) Bill {
	b.Paid = false
	if b.State != BillStateDeleted {
		b.State = BillStateDue
	}
	b.PaidAt = nil
	b.PaidAmount = nil
	return b
}

func dispatchCallback(ctx context.Context, dispatcher *webhook.Dispatcher, b Bill) {
	if dispatcher == nil || b.CallbackURL == "" {
		return
	}
	status := webhook.PaymentStatusUnpaid
	if b.State == BillStatePaid {
		status = webhook.PaymentStatusPaid
	}
	d := cloneDispatcher(dispatcher, b.CallbackURL)
	_, _ = d.Dispatch(ctx, b.ID, status)
}

func cloneDispatcher(src *webhook.Dispatcher, url string) *webhook.Dispatcher {
	return &webhook.Dispatcher{
		URL:           url,
		Secret:        src.Secret,
		MaxRetries:    src.MaxRetries,
		ProviderName:  src.ProviderName,
		Builder:       src.Builder,
		HeaderBuilder: src.HeaderBuilder,
		EventTypeFor:  src.EventTypeFor,
		EnabledEvents: src.EnabledEvents,
		Worker:        src.Worker,
		Store:         src.Store,
		AuditLogger:   src.AuditLogger,
	}
}

func redirectURL(b Bill, key string) string {
	values := map[string]string{
		"billplz[id]":    b.ID,
		"billplz[paid]":  strconv.FormatBool(b.Paid),
		"billplz[state]": string(b.State),
	}
	sig := Sign(values, key)

	u, _ := url.Parse(b.RedirectURL)
	q := u.Query()
	q.Set("billplz[id]", b.ID)
	q.Set("billplz[paid]", strconv.FormatBool(b.Paid))
	q.Set("billplz[state]", string(b.State))
	q.Set("x_signature", sig)
	u.RawQuery = q.Encode()
	return u.String()
}
