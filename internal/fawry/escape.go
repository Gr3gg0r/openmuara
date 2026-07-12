package fawry

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/ui"
	"github.com/openmuara/openmuara/internal/webhook"
)

// NewEscapeHandler returns the GET handler for the Fawry escape page.
func NewEscapeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ref := r.URL.Query().Get("ref")
		returnURL := r.URL.Query().Get("returnUrl")
		amount := r.URL.Query().Get("amount")
		if amount == "" {
			amount = "0.00"
		}

		if ref == "" || returnURL == "" {
			err := errcode.New(errcode.EInvalidRequest, "ref and returnUrl are required")
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		data := ui.EscapePageData{
			Ref:       ref,
			ReturnURL: returnURL,
			Amount:    amount,
		}
		if tok, ok := httputil.CSRFTokenFromContext(r.Context()); ok {
			data.CSRFToken = tok
		}
		if err := ui.RenderEscapePage(w, data); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to render escape page")
			return
		}
	}
}

// NewEscapeActionHandler returns the POST handler that simulates payment outcome.
// It updates the shared ledger, dispatches an outgoing webhook if configured,
// and then redirects back to the caller.
func NewEscapeActionHandler(dispatcher *webhook.Dispatcher, ledger engine.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}

		ref := r.FormValue("ref")
		returnURL := r.FormValue("returnUrl")
		status := r.FormValue("status")

		if ref == "" || returnURL == "" || status == "" {
			err := errcode.New(errcode.EInvalidRequest, "ref, returnUrl, and status are required")
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if ledger == nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "ledger not available")
			return
		}

		tx, ok, err := ledger.GetByReference(ref)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		targetStatus := engine.TransactionStatusUnpaid
		if status == "PAID" {
			targetStatus = engine.TransactionStatusPaid
		}

		if err := engine.Transition(&tx, targetStatus); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, err.Error())
			return
		}

		if _, _, err := ledger.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to update transaction")
			return
		}

		audit.FromContext(r.Context()).Log(r.Context(), "admin.escape", "transaction", ref, status, "ok")

		if dispatcher != nil {
			paymentStatus := webhook.PaymentStatus(status)
			if _, dispatchErr := dispatcher.Dispatch(r.Context(), ref, paymentStatus); dispatchErr != nil {
				slog.Warn("failed to dispatch fawry webhook", "ref", ref, "error", dispatchErr)
				httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to dispatch webhook")
				return
			}
		}

		redirectURL, err := buildCallbackURL(returnURL, status)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "invalid returnUrl")
			return
		}

		// #nosec G710 -- redirect target is built from caller-supplied returnUrl in emulation
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
