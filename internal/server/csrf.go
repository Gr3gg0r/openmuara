package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// CSRFGuardConfig configures the CSRF guard middleware at runtime.
type CSRFGuardConfig struct {
	Enabled        bool
	SecureCookie   bool // set when TLS is enabled
	SameSiteStrict bool // set when admin auth is enabled
}

const (
	csrfCookieName = "openmuara_csrf"
	csrfHeaderName = "X-CSRF-Token"
	csrfFormField  = "csrf_token"
)

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func setCSRFCookie(w http.ResponseWriter, token string, cfg CSRFGuardConfig) {
	sameSite := http.SameSiteLaxMode
	if cfg.SameSiteStrict {
		sameSite = http.SameSiteStrictMode
	}
	// #nosec G124 -- Secure is set only when TLS is enabled; omitted over plain HTTP for local dev
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		SameSite: sameSite,
		HttpOnly: true,
		Secure:   cfg.SecureCookie,
	})
}

// CSRFTokenFromCookie returns the token stored in the CSRF cookie.
func CSRFTokenFromCookie(r *http.Request) (string, bool) {
	c, err := r.Cookie(csrfCookieName)
	if err != nil || c == nil || c.Value == "" {
		return "", false
	}
	return c.Value, true
}

// CSRFRequestToken returns the token submitted with the request via header or form.
func CSRFRequestToken(r *http.Request) string {
	if tok := r.Header.Get(csrfHeaderName); tok != "" {
		return tok
	}
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
		_ = r.ParseForm()
		return r.FormValue(csrfFormField)
	}
	return ""
}

var csrfExemptAdminPaths = map[string]struct{}{
	"/_admin/webhook-receiver": {},
}

func isAdminPath(path string) bool {
	if !strings.HasPrefix(path, "/_admin/") {
		return false
	}
	if _, exempt := csrfExemptAdminPaths[path]; exempt {
		return false
	}
	return true
}

func isCheckoutPayPath(path string) bool {
	return strings.HasPrefix(path, "/v1/checkout/sessions/") && strings.HasSuffix(path, "/pay")
}

func isCSRFProtectedPath(path string) bool {
	return isAdminPath(path) || isCheckoutPayPath(path)
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

// CSRFGuardMiddleware enforces a double-submit CSRF cookie for admin mutations.
func CSRFGuardMiddleware(cfg CSRFGuardConfig) Middleware {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isSafeMethod(r.Method) {
				token, ok := CSRFTokenFromCookie(r)
				if !ok {
					var err error
					token, err = generateCSRFToken()
					if err != nil {
						http.Error(w, "failed to generate csrf token", http.StatusInternalServerError)
						return
					}
					setCSRFCookie(w, token, cfg)
				}
				r = r.WithContext(httputil.WithCSRFToken(r.Context(), token))
				next.ServeHTTP(w, r)
				return
			}

			if isCSRFProtectedPath(r.URL.Path) {
				cookieToken, ok := CSRFTokenFromCookie(r)
				if !ok {
					http.Error(w, "csrf cookie missing", http.StatusForbidden)
					return
				}
				requestToken := CSRFRequestToken(r)
				if requestToken == "" || !strings.EqualFold(cookieToken, requestToken) {
					http.Error(w, "csrf token mismatch", http.StatusForbidden)
					return
				}
				r = r.WithContext(httputil.WithCSRFToken(r.Context(), cookieToken))
			}

			next.ServeHTTP(w, r)
		})
	}
}
