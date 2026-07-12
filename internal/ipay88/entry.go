package ipay88

import (
	"fmt"
	"net/http"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

func (p *Provider) entryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		req := parseEntryForm(r)
		if err := validateEntryRequest(req); err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if req.SignatureType != "SHA256" {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "signature type must be SHA256")
			return
		}

		if err := IsPublicURL(req.ResponseURL); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusBadRequest, "invalid response url")
			return
		}
		if err := IsPublicURL(req.BackendURL); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusBadRequest, "invalid backend url")
			return
		}

		if !VerifyRequest(req, p.merchantKey) {
			httputil.RespondError(w, r, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid signature")
			return
		}

		amount, err := parseAmount(req.Amount)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "invalid amount")
			return
		}

		tx := engine.NewTransaction(engine.Transaction{
			Provider:    ProviderName,
			Type:        "charge",
			Amount:      amount,
			Currency:    req.Currency,
			Status:      engine.TransactionStatusNew,
			CustomerRef: req.UserEmail,
			Reference:   req.RefNo,
			TraceID:     httputil.TraceIDFromContext(r.Context()),
		})
		if _, _, err := p.store.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to record transaction")
			return
		}

		p.saveRequest(req)

		http.Redirect(w, r, fmt.Sprintf("/_admin/ipay88/pay/%s", req.RefNo), http.StatusSeeOther)
	}
}

func parseEntryForm(r *http.Request) PaymentRequest {
	return PaymentRequest{
		MerchantCode:  r.FormValue("MerchantCode"),
		PaymentID:     r.FormValue("PaymentId"),
		RefNo:         r.FormValue("RefNo"),
		Amount:        r.FormValue("Amount"),
		Currency:      r.FormValue("Currency"),
		ProdDesc:      r.FormValue("ProdDesc"),
		UserName:      r.FormValue("UserName"),
		UserEmail:     r.FormValue("UserEmail"),
		UserContact:   r.FormValue("UserContact"),
		Remark:        r.FormValue("Remark"),
		Lang:          r.FormValue("Lang"),
		Signature:     r.FormValue("Signature"),
		SignatureType: r.FormValue("SignatureType"),
		ResponseURL:   r.FormValue("ResponseURL"),
		BackendURL:    r.FormValue("BackendURL"),
	}
}

func validateEntryRequest(req PaymentRequest) error {
	switch {
	case req.MerchantCode == "":
		return errcode.New(errcode.EInvalidRequest, "merchant code is required")
	case req.RefNo == "":
		return errcode.New(errcode.EInvalidRequest, "ref no is required")
	case req.Amount == "":
		return errcode.New(errcode.EInvalidRequest, "amount is required")
	case req.Currency == "":
		return errcode.New(errcode.EInvalidRequest, "currency is required")
	case req.ProdDesc == "":
		return errcode.New(errcode.EInvalidRequest, "product description is required")
	case req.UserName == "":
		return errcode.New(errcode.EInvalidRequest, "user name is required")
	case req.UserEmail == "":
		return errcode.New(errcode.EInvalidRequest, "user email is required")
	case req.UserContact == "":
		return errcode.New(errcode.EInvalidRequest, "user contact is required")
	case req.Signature == "":
		return errcode.New(errcode.ESignatureMissing, "signature is required")
	case req.SignatureType == "":
		return errcode.New(errcode.EInvalidRequest, "signature type is required")
	case req.ResponseURL == "":
		return errcode.New(errcode.EInvalidRequest, "response url is required")
	case req.BackendURL == "":
		return errcode.New(errcode.EInvalidRequest, "backend url is required")
	}
	return nil
}
