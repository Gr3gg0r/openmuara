package billplz

import (
	"encoding/json"
	"net/http"

	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
)

// NewCollectionsHandler returns the handler for /api/v3/collections.
func NewCollectionsHandler(apiKey string, store *collectionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		var req CreateCollectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.Title == "" {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, "title is required")
			return
		}

		c := store.create(req)
		writeJSON(w, CollectionResponse{Collection: c})
	}
}

// NewCollectionHandler returns the handler for /api/v3/collections/{id}.
func NewCollectionHandler(apiKey string, store *collectionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		id := r.PathValue("id")
		c, ok := store.get(id)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "collection not found")
			return
		}
		writeJSON(w, CollectionResponse{Collection: c})
	}
}

// NewPaymentMethodsHandler returns the handler for /api/v3/collections/{id}/payment_methods.
func NewPaymentMethodsHandler(apiKey string, store *collectionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := requireBasicAuth(r, apiKey); err != nil {
			httputil.RespondError(w, r, httputil.ErrUnauthorized, http.StatusUnauthorized, errcode.Message(err))
			return
		}

		id := r.PathValue("id")
		if _, ok := store.get(id); !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "collection not found")
			return
		}

		writeJSON(w, PaymentMethodsResponse{PaymentMethods: defaultPaymentMethods()})
	}
}

func defaultPaymentMethods() []PaymentMethod {
	return []PaymentMethod{
		{Code: "fpx", Name: "FPX Online Banking", Active: true},
		{Code: "mpgs", Name: "Credit / Debit Card", Active: true},
		{Code: "boost", Name: "Boost Wallet", Active: true},
		{Code: "touchngo", Name: "Touch 'n Go eWallet", Active: true},
		{Code: "twoctwopipp", Name: "Buy Now Pay Later", Active: true},
	}
}

func requireBasicAuth(r *http.Request, apiKey string) error {
	user, _, ok := r.BasicAuth()
	if !ok {
		return errcode.New(errcode.ESignatureMissing, "missing basic auth")
	}
	if user != apiKey {
		return errcode.New(errcode.ESignatureMismatch, "invalid api key")
	}
	return nil
}
