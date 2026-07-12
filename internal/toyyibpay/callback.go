package toyyibpay

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func (p *Provider) returnHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		billCode := r.URL.Query().Get("billcode")
		bill, ok := p.bills.GetByCode(billCode)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}

		redirectURL, err := buildReturnURL(bill.BillReturnURL, r.URL.Query())
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "invalid return url")
			return
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func buildReturnURL(returnURL string, q url.Values) (string, error) {
	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}
	out := u.Query()
	for _, key := range []string{"status_id", "billcode", "order_id", "transaction_id", "msg"} {
		if v := q.Get(key); v != "" {
			out.Set(key, v)
		}
	}
	u.RawQuery = out.Encode()
	return u.String(), nil
}

// WebhookHandler receives incoming ToyyibPay-style callback notifications.
func (p *Provider) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !VerifyCallback(p.secret, r.Form) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid callback hash")
			return
		}
		if err := p.applyIncomingCallback(r.Form); err != nil {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, errcode.Message(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (p *Provider) applyIncomingCallback(values url.Values) error {
	orderID := values.Get("order_id")
	if orderID == "" {
		return errcode.New(errcode.EInvalidRequest, "order_id is required")
	}
	if p.store == nil {
		return errcode.New(errcode.EInternal, "store not available")
	}
	tx, ok, err := p.store.GetByReference(orderID)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "lookup transaction", err)
	}
	if !ok {
		return errcode.New(errcode.ETransactionNotFound, "transaction not found")
	}

	target := engine.TransactionStatusUnpaid
	if values.Get("status") == "1" {
		target = engine.TransactionStatusPaid
	}
	if err := engine.Transition(&tx, target); err != nil {
		return errcode.Wrap(errcode.ETransactionTransitionInvalid, "invalid transaction transition", err)
	}
	_, _, err = p.store.CreateOrGet(tx)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "failed to update transaction", err)
	}
	return nil
}

func (p *Provider) dispatchCallback(ctx context.Context, bill Bill, status string) error {
	if p.dispatcher == nil {
		return errcode.New(errcode.EWebhookURLMissing, "dispatcher not configured")
	}

	paymentStatus := webhook.PaymentStatusUnpaid
	if status == "1" {
		paymentStatus = webhook.PaymentStatusPaid
	}

	d := webhook.NewDispatcherFromProvider(bill.BillCallbackURL, p.dispatcher.MaxRetries, p)
	d.Store = p.dispatcher.Store
	d.AuditLogger = p.dispatcher.AuditLogger
	_, err := d.Dispatch(ctx, bill.OrderID, paymentStatus)
	if err != nil {
		return errcode.Wrap(errcode.EWebhookDeliveryFailed, "dispatch callback", err)
	}
	return nil
}

func (p *Provider) buildCallbackPayload(_ context.Context, tx provider.Transaction) ([]byte, error) {
	bill, ok := p.bills.GetByOrderID(tx.Reference)
	if !ok {
		return nil, errcode.New(errcode.ETransactionNotFound, fmt.Sprintf("bill not found for order %q", tx.Reference))
	}
	refno := "MUARA-" + uuid.Must(uuid.NewRandom()).String()
	status := callbackStatusFromTransactionStatus(tx.Status)
	reason := callbackReason(status)
	transactionTime := time.Now().Format("2006-01-02 15:04:05")
	hash := ComputeHash(p.secret, status, bill.OrderID, refno)

	form := url.Values{}
	form.Set("refno", refno)
	form.Set("status", status)
	form.Set("reason", reason)
	form.Set("billcode", bill.BillCode)
	form.Set("order_id", bill.OrderID)
	form.Set("amount", strconv.Itoa(bill.BillAmount))
	form.Set("transaction_time", transactionTime)
	form.Set("hash", hash)
	return []byte(form.Encode()), nil
}

func callbackStatusFromTransactionStatus(status string) string {
	if status == string(engine.TransactionStatusPaid) || status == "PAID" {
		return "1"
	}
	return "3"
}

func callbackReason(status string) string {
	if status == "1" {
		return "Payment success"
	}
	return "Payment cancelled"
}
