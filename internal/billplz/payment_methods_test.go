package billplz_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/billplz"
)

func TestPaymentMethodsList(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)

	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}/payment_methods")
	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/"+c.ID+"/payment_methods", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp billplz.PaymentMethodsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	codes := make(map[string]bool)
	for _, m := range resp.PaymentMethods {
		codes[m.Code] = true
	}
	for _, want := range []string{"fpx", "mpgs", "boost", "touchngo"} {
		if !codes[want] {
			t.Errorf("missing payment method %q", want)
		}
	}
}

func TestPaymentMethodsCollectionNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodGet, "/api/v3/collections/{id}/payment_methods")

	req := httptest.NewRequest(http.MethodGet, "/api/v3/collections/missing/payment_methods", nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}
