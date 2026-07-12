package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/billplz"
)

func TestRedirectMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodGet, "/billplz/redirect")

	req := httptest.NewRequest(http.MethodPost, "/billplz/redirect", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestRedirectBillHasNoRedirectURL(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)

	mux := testMux(t, p)
	reqBody := map[string]any{
		"collection_id": c.ID,
		"email":         "test@example.com",
		"name":          "Test User",
		"amount":        1000,
		"callback_url":  "http://localhost:9999/callback",
		"description":   "Test bill",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("create bill failed: %d, %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode bill: %v", err)
	}

	redirectHandler := routeFor(t, p, http.MethodGet, "/billplz/redirect")
	redirectReq := httptest.NewRequest(http.MethodGet, "/billplz/redirect?billplz[id]="+resp.Bill.ID, nil)
	redirectRec := httptest.NewRecorder()
	redirectHandler.ServeHTTP(redirectRec, redirectReq)

	if redirectRec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d", redirectRec.Code)
	}
}

func TestVerifyRedirectSignatureMissingSignature(t *testing.T) {
	query := map[string]string{"billplz[id]": "bill-123"}
	if billplz.VerifyRedirectSignature(query, "secret") {
		t.Fatal("expected verification to fail when x_signature is missing")
	}
}
