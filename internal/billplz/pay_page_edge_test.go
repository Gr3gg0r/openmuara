package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

func TestPayPageMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodGet, "/_admin/billplz/pay/{id}")

	req := httptest.NewRequest(http.MethodPost, "/_admin/billplz/pay/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestPayPageDeletedBill(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	deleteHandler := routeFor(t, p, http.MethodDelete, "/api/v3/bills/{id}")
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v3/bills/"+b.ID, nil)
	deleteReq.SetBasicAuth("muara-billplz-api-key", "")
	deleteRec := httptest.NewRecorder()
	deleteHandler.ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("delete bill failed: %d", deleteRec.Code)
	}

	payHandler := routeFor(t, p, http.MethodGet, "/_admin/billplz/pay/{id}")
	req := httptest.NewRequest(http.MethodGet, "/_admin/billplz/pay/"+b.ID, nil)
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok"))
	rec := httptest.NewRecorder()
	payHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d", rec.Code)
	}
}

func TestPayActionMethodNotAllowed(t *testing.T) {
	p := newInitializedProvider(t)
	handler := handlerFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")

	req := httptest.NewRequest(http.MethodGet, "/_admin/billplz/pay/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestPayActionBillNotFound(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")

	req := mustPayRequest("missing")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestPayActionInvalidOutcome(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	form := url.Values{"outcome": {"invalid"}, "csrf_token": {"tok"}}
	req := httptest.NewRequest(http.MethodPost, "/_admin/billplz/pay/"+b.ID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok"))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestPayActionNoRedirectURL(t *testing.T) {
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

	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	payReq := mustPayRequest(resp.Bill.ID)
	payRec := httptest.NewRecorder()
	handler.ServeHTTP(payRec, payReq)

	if payRec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", payRec.Code)
	}
}
