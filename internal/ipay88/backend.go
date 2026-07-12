package ipay88

import (
	"fmt"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/webhook"
)

func (p *Provider) backendHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		ref := r.FormValue("RefNo")
		status := r.FormValue("Status")
		paymentID := r.FormValue("PaymentId")
		signature := r.FormValue("Signature")
		signatureType := r.FormValue("SignatureType")

		req, ok := p.getRequest(ref)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "payment request not found")
			return
		}

		if signatureType != "" && signatureType != "SHA256" {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "signature type must be SHA256")
			return
		}

		if paymentID == "" {
			paymentID = req.SelectedPaymentID
			if paymentID == "" {
				paymentID = req.PaymentID
			}
		}

		if !VerifyResponse(p.merchantKey, p.merchantCode, paymentID, ref, req.Amount, req.Currency, status, signature) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid signature")
			return
		}

		if err := p.updateTransactionFromStatus(ref, status); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, errcode.Message(err))
			return
		}

		if p.dispatcher != nil {
			paymentStatus := mapIPay88StatusToWebhook(status)
			if _, dispatchErr := p.dispatcher.Dispatch(r.Context(), ref, paymentStatus); dispatchErr != nil {
				httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to dispatch webhook")
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "RECEIVEOK")
	}
}

func (p *Provider) updateTransactionFromStatus(ref, status string) error {
	target := engine.TransactionStatusUnpaid
	if status == string(PaymentStatusSuccess) {
		target = engine.TransactionStatusPaid
	}
	return p.transitionTransaction(ref, target)
}

func mapIPay88StatusToWebhook(status string) webhook.PaymentStatus {
	if status == string(PaymentStatusSuccess) {
		return webhook.PaymentStatusPaid
	}
	return webhook.PaymentStatusUnpaid
}
