package toyyibpay

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/ui"
)

func (p *Provider) payPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		billCode := r.PathValue("billCode")
		bill, ok := p.bills.GetByCode(billCode)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		if bill.BillStatus != "1" {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, "bill is inactive")
			return
		}

		channel, _ := strconv.Atoi(bill.BillPaymentChannel)
		data := ui.ToyyibPayPageData{
			BillCode:      bill.BillCode,
			BillName:      bill.BillName,
			Amount:        int64(bill.BillAmount),
			AmountDisplay: fmt.Sprintf("%.2f", float64(bill.BillAmount)/100.0),
			Channel:       channel,
		}
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}
		if err := ui.RenderToyyibPayPage(w, data); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to render page")
			return
		}
	}
}

func (p *Provider) payPageActionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		billCode := r.PathValue("billCode")
		bill, ok := p.bills.GetByCode(billCode)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}

		status := r.FormValue("status")
		if status != "1" && status != "3" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "status must be 1 or 3")
			return
		}

		if err := p.recordPaymentOutcome(bill, status); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, errcode.Message(err))
			return
		}

		if err := p.dispatchCallback(r.Context(), bill, status); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "callback dispatch failed")
			return
		}

		redirectURL := fmt.Sprintf("/toyyibpay/return?status_id=%s&billcode=%s&order_id=%s&transaction_id=%s&msg=%s",
			status, bill.BillCode, bill.OrderID, "MUARA-"+bill.BillCode, urlEncodeMsg(status))
		// #nosec G710 -- redirect target is internal provider return path in emulation
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (p *Provider) recordPaymentOutcome(bill Bill, status string) error {
	if p.store == nil {
		return errcode.New(errcode.EInternal, "store not available")
	}
	tx, ok, err := p.store.GetByReference(bill.OrderID)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "lookup transaction", err)
	}
	if !ok {
		return errcode.New(errcode.ETransactionNotFound, "transaction not found")
	}

	target := engine.TransactionStatusUnpaid
	if status == "1" {
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

func urlEncodeMsg(status string) string {
	if status == "1" {
		return "Payment+success"
	}
	return "Payment+cancelled"
}
