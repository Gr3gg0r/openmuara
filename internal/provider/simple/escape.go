package simple

import (
	"net/http"
	"net/url"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/ui"
	"github.com/openmuara/openmuara/internal/webhook"
)

func (p *Provider) escapePageHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		cfg := p.runtime.EscapePage
		ref := req.URL.Query().Get(cfg.RefParam)
		returnURL := req.URL.Query().Get(cfg.ReturnParam)
		amount := req.URL.Query().Get(cfg.AmountParam)
		if amount == "" {
			amount = "0.00"
		}

		if ref == "" || returnURL == "" {
			err := errcode.New(errcode.EInvalidRequest, "ref and returnUrl are required")
			httputil.RespondError(w, req, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		data := ui.EscapePageData{
			Ref:       ref,
			ReturnURL: returnURL,
			Amount:    amount,
		}
		if tok, ok := httputil.CSRFTokenFromContext(req.Context()); ok {
			data.CSRFToken = tok
		}
		if err := ui.RenderEscapePage(w, data); err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to render escape page")
			return
		}
	}
}

func (p *Provider) escapeActionHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := req.ParseForm(); err != nil {
			httputil.RespondError(w, req, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		cfg := p.runtime.EscapePage
		ref := req.FormValue(cfg.RefParam)
		returnURL := req.FormValue(cfg.ReturnParam)
		status := req.FormValue(cfg.StatusParam)
		if status == "" {
			status = req.FormValue("status")
		}

		if ref == "" || returnURL == "" || status == "" {
			err := errcode.New(errcode.EInvalidRequest, "ref, returnUrl, and status are required")
			httputil.RespondError(w, req, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		tx, ok, err := p.store.GetByReference(ref)
		if err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, req, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		targetStatus := engine.TransactionStatusUnpaid
		if status == "PAID" {
			targetStatus = engine.TransactionStatusPaid
		}

		if err := engine.Transition(&tx, targetStatus); err != nil {
			httputil.RespondError(w, req, httputil.ErrInvalidState, http.StatusConflict, err.Error())
			return
		}

		if _, _, err := p.store.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to update transaction")
			return
		}

		audit.FromContext(req.Context()).Log(req.Context(), "admin.escape", "transaction", ref, status, "ok")

		if p.dispatcher != nil {
			paymentStatus := webhook.PaymentStatus(status)
			if _, dispatchErr := p.dispatcher.Dispatch(req.Context(), ref, paymentStatus); dispatchErr != nil {
				httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to dispatch webhook")
				return
			}
		}

		redirectURL, err := buildCallbackURL(returnURL, status)
		if err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "invalid returnUrl")
			return
		}

		// #nosec G710 -- redirect target is built from caller-supplied returnUrl in emulation
		http.Redirect(w, req, redirectURL, http.StatusSeeOther)
	}
}

func buildCallbackURL(returnURL, status string) (string, error) {
	u, err := url.Parse(returnURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("orderStatus", status)
	q.Set("statusCode", "200")
	u.RawQuery = q.Encode()
	return u.String(), nil
}
