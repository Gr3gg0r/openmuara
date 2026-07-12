package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestPayPageRenders(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodGet, "/_admin/billplz/pay/{id}")
	req := httptest.NewRequest(http.MethodGet, "/_admin/billplz/pay/"+b.ID, nil)
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "test-csrf"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, b.ID) {
		t.Error("page body missing bill id")
	}
	if !strings.Contains(body, "test-csrf") {
		t.Error("page body missing csrf token")
	}
	if !strings.Contains(body, "fpx") {
		t.Error("page body missing fpx payment method")
	}
}

func TestPayPagePayOutcome(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	form := url.Values{"outcome": {"pay"}, "csrf_token": {"tok"}}
	req := httptest.NewRequest(http.MethodPost, "/_admin/billplz/pay/"+b.ID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK && rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 200 or 303, got %d, body: %s", rec.Code, rec.Body.String())
	}

	// Verify the bill is now paid.
	getHandler := routeFor(t, p, http.MethodGet, "/api/v3/bills/{id}")
	getReq := httptest.NewRequest(http.MethodGet, "/api/v3/bills/"+b.ID, nil)
	getReq.SetBasicAuth("muara-billplz-api-key", "")
	getRec := httptest.NewRecorder()
	getHandler.ServeHTTP(getRec, getReq)

	var resp billplz.BillResponse
	if err := json.Unmarshal(getRec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode bill: %v", err)
	}
	if !resp.Bill.Paid {
		t.Error("bill should be paid")
	}
	if resp.Bill.State != billplz.BillStatePaid {
		t.Errorf("state: want paid, got %q", resp.Bill.State)
	}
	if resp.Bill.PaidAmount == nil || *resp.Bill.PaidAmount != b.Amount {
		t.Error("paid_amount mismatch")
	}
}

func TestPayPagePayDispatchesCallback(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)

	received := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		received <- r.Form.Get("id") + "|" + r.Form.Get("state") + "|" + r.Form.Get("x_signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b := createBillWithCallback(t, p, c.ID, server.URL)

	d := webhook.NewDispatcherFromBuilder(server.URL, 0, p.PayloadBuilder(), p.PayloadHeaders)
	p.SetDispatcher(d)

	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	req := mustPayRequest(b.ID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	select {
	case payload := <-received:
		parts := strings.Split(payload, "|")
		if parts[0] != b.ID {
			t.Errorf("callback id: want %q, got %q", b.ID, parts[0])
		}
		if parts[1] != string(billplz.BillStatePaid) {
			t.Errorf("callback state: want paid, got %q", parts[1])
		}
		if parts[2] == "" {
			t.Error("callback x_signature is empty")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("callback was not received")
	}
}

func createBillWithCallback(t *testing.T, p *billplz.Provider, collectionID, callbackURL string) billplz.Bill {
	t.Helper()
	mux := testMux(t, p)
	reqBody := map[string]any{
		"collection_id": collectionID,
		"email":         "test@example.com",
		"mobile":        "+60123456789",
		"name":          "Test User",
		"amount":        1000,
		"callback_url":  callbackURL,
		"description":   "Test bill",
		"redirect_url":  "http://localhost:9999/redirect",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create bill failed: status %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode bill: %v", err)
	}
	return resp.Bill
}

func TestPayPageCancelOutcome(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	form := url.Values{"outcome": {"cancel"}, "csrf_token": {"tok"}}
	req := httptest.NewRequest(http.MethodPost, "/_admin/billplz/pay/"+b.ID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK && rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 200 or 303, got %d", rec.Code)
	}

	getHandler := routeFor(t, p, http.MethodGet, "/api/v3/bills/{id}")
	getReq := httptest.NewRequest(http.MethodGet, "/api/v3/bills/"+b.ID, nil)
	getReq.SetBasicAuth("muara-billplz-api-key", "")
	getRec := httptest.NewRecorder()
	getHandler.ServeHTTP(getRec, getReq)

	var resp billplz.BillResponse
	if err := json.Unmarshal(getRec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode bill: %v", err)
	}
	if resp.Bill.Paid {
		t.Error("bill should not be paid after cancel")
	}
}
