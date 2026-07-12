package toyyibpay

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/google/uuid"
)

// CategoryStore is a thread-safe in-memory category repository.
type CategoryStore struct {
	mu  sync.RWMutex
	by  map[string]Category
	nam map[string]string
}

// NewCategoryStore creates an empty category store.
func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		by:  make(map[string]Category),
		nam: make(map[string]string),
	}
}

// Create stores a category and returns it with a generated code.
func (s *CategoryStore) Create(name, description string) Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	cat := Category{
		CategoryCode:        uuid.Must(uuid.NewRandom()).String(),
		CategoryName:        name,
		CategoryDescription: description,
		CategoryStatus:      "1",
	}
	s.by[cat.CategoryCode] = cat
	s.nam[cat.CategoryName] = cat.CategoryCode
	return cat
}

// Get returns a category by code.
func (s *CategoryStore) Get(code string) (Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.by[code]
	return c, ok
}

func (p *Provider) categoryCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !p.authenticate(r.FormValue("userSecretKey"), w, r) {
			return
		}

		// Accept the official ToyyibPay parameter names (catname,
		// catdescription); the camelCase aliases are kept for backward
		// compatibility with earlier OpenMuara docs.
		name := firstNonEmpty(r.FormValue("catname"), r.FormValue("categoryName"))
		desc := firstNonEmpty(r.FormValue("catdescription"), r.FormValue("categoryDescription"))
		cat := p.categories.Create(name, desc)
		writeJSON(w, CategoryResponse{Status: "1", Msg: "success", Data: cat})
	}
}

func (p *Provider) categoryDetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !p.authenticate(r.FormValue("userSecretKey"), w, r) {
			return
		}

		code := r.FormValue("categoryCode")
		if code == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "categoryCode is required")
			return
		}
		cat, ok := p.categories.Get(code)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "category not found")
			return
		}
		writeJSON(w, cat)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}

var errInvalidSecret = errcode.New(errcode.ESignatureMismatch, "toyyibpay: invalid user_secret_key")

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func (p *Provider) authenticate(secret string, w http.ResponseWriter, r *http.Request) bool {
	if secret != p.secret {
		httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(errInvalidSecret))
		return false
	}
	return true
}
