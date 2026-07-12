package ipay88

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/ui"
)

var paymentMethods = []ui.IPay88PaymentMethod{
	{ID: "1", Name: "Credit / Debit Card"},
	{ID: "2", Name: "FPX"},
	{ID: "33", Name: "Touch 'n Go eWallet"},
	{ID: "34", Name: "Boost"},
	{ID: "35", Name: "GrabPay"},
}

func (p *Provider) adminPayPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ref := r.PathValue("refNo")
		req, ok := p.getRequest(ref)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "payment request not found")
			return
		}

		amount, err := parseAmount(req.Amount)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "invalid stored amount")
			return
		}

		data := ui.IPay88PayPageData{
			RefNo:         req.RefNo,
			Amount:        amountInCents(amount),
			AmountDisplay: fmt.Sprintf("%.2f", amount),
			Currency:      req.Currency,
			Description:   req.ProdDesc,
			Methods:       paymentMethods,
		}
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := ui.RenderIPay88PayPage(w, data); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to render payment page")
		}
	}
}

func (p *Provider) adminPayActionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		ref := r.PathValue("refNo")
		req, ok := p.getRequest(ref)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "payment request not found")
			return
		}

		outcome := r.FormValue("outcome")
		paymentID := r.FormValue("payment_method")
		if paymentID == "" {
			paymentID = req.PaymentID
		}

		status := string(PaymentStatusFailure)
		targetStatus := engine.TransactionStatusUnpaid
		if outcome == "pay" {
			status = string(PaymentStatusSuccess)
			targetStatus = engine.TransactionStatusPaid
		}

		if err := p.transitionTransaction(ref, targetStatus); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, errcode.Message(err))
			return
		}

		req.SelectedPaymentID = paymentID
		req.Status = status
		p.saveRequest(req)

		if err := p.postBackendCallback(context.Background(), req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "backend post failed")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderResponseForward(w, req)
	}
}

func (p *Provider) transitionTransaction(ref string, target engine.TransactionStatus) error {
	tx, ok, err := p.store.GetByReference(ref)
	if err != nil {
		return errcode.Wrap(errcode.EInternal, "failed to lookup transaction", err)
	}
	if !ok {
		return errcode.Wrap(errcode.ETransactionTransitionInvalid, "invalid transaction transition", engine.ErrInvalidTransition)
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

func (p *Provider) postBackendCallback(ctx context.Context, req PaymentRequest) error {
	values := responseValues(req, p.merchantCode, req.SelectedPaymentID, req.Amount, req.Currency, req.Status, p.merchantKey)
	body := strings.NewReader(values.Encode())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.BackendURL, body)
	if err != nil {
		return errcode.Wrap(errcode.EWebhookBuildFailed, "build backend request", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return errcode.Wrap(errcode.EWebhookDeliveryFailed, "post backend callback", err)
	}
	defer func() { _ = resp.Body.Close() }()
	ack, err := io.ReadAll(resp.Body)
	if err != nil {
		return errcode.Wrap(errcode.EWebhookDeliveryFailed, "read backend response", err)
	}
	if !strings.EqualFold(strings.TrimSpace(string(ack)), "RECEIVEOK") {
		return errcode.New(errcode.EWebhookDeliveryFailed, "backend did not acknowledge with RECEIVEOK")
	}
	return nil
}

func renderResponseForward(w io.Writer, req PaymentRequest) {
	_, _ = fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>iPay88 Response</title></head>
<body onload="document.forms[0].submit()">
<form method="POST" action="/ipay88/response">
<input type="hidden" name="RefNo" value="%s">
<input type="hidden" name="Status" value="%s">
<input type="hidden" name="PaymentId" value="%s">
<noscript><button type="submit">Continue</button></noscript>
</form>
</body>
</html>`, req.RefNo, req.Status, req.SelectedPaymentID)
}

func responseValues(req PaymentRequest, merchantCode, paymentID, amount, currency, status, merchantKey string) url.Values {
	if paymentID == "" {
		paymentID = req.SelectedPaymentID
		if paymentID == "" {
			paymentID = req.PaymentID
		}
	}
	v := url.Values{}
	v.Set("MerchantCode", merchantCode)
	v.Set("PaymentId", paymentID)
	v.Set("RefNo", req.RefNo)
	v.Set("Amount", req.Amount)
	v.Set("Currency", currency)
	v.Set("Remark", req.Remark)
	v.Set("Status", status)
	v.Set("SignatureType", "SHA256")
	v.Set("Signature", SignResponse(merchantKey, merchantCode, paymentID, req.RefNo, amount, currency, status))
	return v
}
