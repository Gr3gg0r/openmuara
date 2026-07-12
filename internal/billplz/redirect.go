package billplz

import (
	"net/http"

	"github.com/openmuara/openmuara/internal/httputil"
)

// NewRedirectHandler returns GET /billplz/redirect.
// It builds a signed query string for the requested bill and redirects to its redirect_url.
func NewRedirectHandler(bills *billStore, xSignatureKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		id := r.URL.Query().Get("billplz[id]")
		if id == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "billplz[id] is required")
			return
		}

		b, ok := bills.get(id)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		if b.RedirectURL == "" {
			httputil.RespondError(w, r, httputil.ErrInvalidState, http.StatusConflict, "bill has no redirect_url")
			return
		}

		http.Redirect(w, r, redirectURL(b, xSignatureKey), http.StatusFound)
	}
}

// VerifyRedirectSignature checks the x_signature on a redirect query string.
func VerifyRedirectSignature(query map[string]string, xSignatureKey string) bool {
	sig, ok := query["x_signature"]
	if !ok {
		return false
	}
	delete(query, "x_signature")
	return Verify(query, xSignatureKey, sig)
}
