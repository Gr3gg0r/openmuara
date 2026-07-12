package ipay88

import (
	"io"
	"net/http"
	"strings"

	"github.com/openmuara/openmuara/internal/httputil"
)

func (p *Provider) responseHandler() http.HandlerFunc {
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

		req, ok := p.getRequest(ref)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "payment request not found")
			return
		}

		if status == "" {
			status = req.Status
		}

		values := responseValues(req, p.merchantCode, paymentID, req.Amount, req.Currency, status, p.merchantKey)
		if err := IsPublicURL(req.ResponseURL); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusBadRequest, "invalid response url")
			return
		}

		resp, err := p.httpClient.Post(req.ResponseURL, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusInternalServerError, "failed to post response url")
			return
		}
		defer func() { _ = resp.Body.Close() }()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}
