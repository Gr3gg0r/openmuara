package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/billplz"
	"github.com/openmuara/openmuara/internal/engine"
)

func TestCreateBillSuccess(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)

	b := createBill(t, p, c.ID)
	if b.ID == "" {
		t.Fatal("bill id is empty")
	}
	if b.CollectionID != c.ID {
		t.Errorf("collection_id: want %q, got %q", c.ID, b.CollectionID)
	}
	if b.Amount != 1000 {
		t.Errorf("amount: want 1000, got %d", b.Amount)
	}
	if b.State != billplz.BillStateDue {
		t.Errorf("state: want due, got %q", b.State)
	}
	if b.Paid {
		t.Error("new bill should not be paid")
	}
	if b.URL == "" {
		t.Error("bill url is empty")
	}

	// Verify the JSON response does not contain a currency field.
	var raw map[string]json.RawMessage
	body, _ := json.Marshal(billplz.BillResponse{Bill: b})
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if _, ok := raw["currency"]; ok {
		t.Error("currency field should not be present in bill response")
	}
}

func TestCreateBillMissingCollectionID(t *testing.T) {
	p := newInitializedProvider(t)
	if err := p.Init(map[string]any{
		"api_key":         "muara-billplz-api-key",
		"x_signature_key": "muara-billplz-xsig-key",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")

	body, _ := json.Marshal(map[string]string{
		"email":        "test@example.com",
		"name":         "Test User",
		"amount":       "1000",
		"callback_url": "http://localhost:9999/callback",
		"description":  "Test bill",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateBillUsesDefaultCollectionID(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	cfg := validProviderConfig()
	cfg["collection_id"] = c.ID
	if err := p.Init(cfg); err != nil {
		t.Fatalf("init: %v", err)
	}

	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")
	body, _ := json.Marshal(map[string]any{
		"email":        "test@example.com",
		"name":         "Test User",
		"amount":       1000,
		"callback_url": "http://localhost:9999/callback",
		"description":  "Test bill",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Bill.CollectionID != c.ID {
		t.Errorf("collection_id: want %q, got %q", c.ID, resp.Bill.CollectionID)
	}
}

func TestCreateBillInvalidCollectionID(t *testing.T) {
	p := newInitializedProvider(t)
	handler := routeFor(t, p, http.MethodPost, "/api/v3/bills")

	body, _ := json.Marshal(map[string]string{
		"collection_id": "missing",
		"email":         "test@example.com",
		"name":          "Test User",
		"amount":        "1000",
		"callback_url":  "http://localhost:9999/callback",
		"description":   "Test bill",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestGetBillSuccess(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodGet, "/api/v3/bills/{id}")
	req := httptest.NewRequest(http.MethodGet, "/api/v3/bills/"+b.ID, nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Bill.ID != b.ID {
		t.Errorf("id: want %q, got %q", b.ID, resp.Bill.ID)
	}
}

func TestDeleteBillSuccess(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	handler := routeFor(t, p, http.MethodDelete, "/api/v3/bills/{id}")
	req := httptest.NewRequest(http.MethodDelete, "/api/v3/bills/"+b.ID, nil)
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Bill.State != billplz.BillStateDeleted {
		t.Errorf("state: want deleted, got %q", resp.Bill.State)
	}
}

func TestCreateBillRecordsLedgerTransaction(t *testing.T) {
	store := engine.NewMemoryStore()
	p := billplz.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetBaseURL("http://localhost:9000")
	p.SetStore(store)

	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	tx, ok, err := store.GetByID(b.ID)
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if !ok {
		t.Fatal("expected ledger transaction to be recorded")
	}
	if tx.Provider != "billplz" {
		t.Errorf("provider: want billplz, got %q", tx.Provider)
	}
	if tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("status: want unpaid, got %q", tx.Status)
	}
	if tx.Amount != 10.00 {
		t.Errorf("amount: want 10.00, got %f", tx.Amount)
	}
	if tx.Currency != "MYR" {
		t.Errorf("currency: want MYR, got %q", tx.Currency)
	}
}

func routeFor(t *testing.T, p *billplz.Provider, method, path string) *http.ServeMux {
	t.Helper()
	found := false
	for _, route := range p.Routes() {
		if route.Method == method && route.Path == path {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("route not found: %s %s", method, path)
	}
	return testMux(t, p)
}
