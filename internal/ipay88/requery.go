package ipay88

import (
	"fmt"
	"net/http"

	"github.com/openmuara/openmuara/internal/httputil"
)

func (p *Provider) requeryHandler() http.HandlerFunc {
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
		req, ok := p.getRequest(ref)
		if !ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, RequeryStatusFailure)
			return
		}

		status := RequeryStatusFailure
		if req.Status == string(PaymentStatusSuccess) {
			status = RequeryStatusSuccess
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, status)
	}
}
