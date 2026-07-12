package senangpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/server"
)

// #nosec G101 -- test fixture secret, not a real credential
const testSecret = "muara-senangpay-secret"

func newTestStatusQuery(orderID string) (string, string) {
	hash := SignStatusQuery(testSecret, orderID)
	return orderID, hash
}

func setupStatusTestStore(t *testing.T, status engine.TransactionStatus, ref string) engine.TransactionStore {
	t.Helper()
	store := engine.NewMemoryStore()
	tx := engine.NewTransaction(engine.Transaction{
		Provider:  ProviderName,
		Type:      "charge",
		Amount:    55.5,
		Currency:  "MYR",
		Status:    status,
		Reference: ref,
	})
	if _, _, err := store.CreateOrGet(tx); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}
	return store
}

func TestStatusHandlerSuccess(t *testing.T) {
	store := setupStatusTestStore(t, engine.TransactionStatusPaid, "order-status-1")
	handler := NewStatusHandler(testSecret, store)

	orderID, hash := newTestStatusQuery("order-status-1")
	url := fmt.Sprintf("/senangpay/query?order_id=%s&hash=%s", orderID, hash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp QueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.OrderID != orderID {
		t.Errorf("order_id: want %q, got %q", orderID, resp.OrderID)
	}
	if resp.StatusID != "1" {
		t.Errorf("status_id: want 1, got %q", resp.StatusID)
	}
	if resp.Status != "paid" {
		t.Errorf("status: want paid, got %q", resp.Status)
	}
	if resp.Amount != 55.5 {
		t.Errorf("amount: want 55.5, got %f", resp.Amount)
	}
	if resp.Currency != "MYR" {
		t.Errorf("currency: want MYR, got %q", resp.Currency)
	}
}

func TestStatusHandlerUnpaid(t *testing.T) {
	store := setupStatusTestStore(t, engine.TransactionStatusUnpaid, "order-status-0")
	handler := NewStatusHandler(testSecret, store)

	orderID, hash := newTestStatusQuery("order-status-0")
	url := fmt.Sprintf("/senangpay/query?order_id=%s&hash=%s", orderID, hash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp QueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.StatusID != "0" {
		t.Errorf("status_id: want 0, got %q", resp.StatusID)
	}
	if resp.Status != "unpaid" {
		t.Errorf("status: want unpaid, got %q", resp.Status)
	}
}

func TestStatusHandlerNotFound(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewStatusHandler(testSecret, store)

	orderID, hash := newTestStatusQuery("missing-order")
	url := fmt.Sprintf("/senangpay/query?order_id=%s&hash=%s", orderID, hash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerInvalidHash(t *testing.T) {
	store := setupStatusTestStore(t, engine.TransactionStatusPaid, "order-status-bad")
	handler := NewStatusHandler(testSecret, store)

	url := "/senangpay/query?order_id=order-status-bad&hash=invalid"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp server.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Error.Code != server.ErrInvalidSignature {
		t.Errorf("code: want %q, got %q", server.ErrInvalidSignature, resp.Error.Code)
	}
}

func TestStatusHandlerMissingOrderID(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewStatusHandler(testSecret, store)

	url := fmt.Sprintf("/senangpay/query?hash=%s", SignStatusQuery(testSecret, "x"))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerMissingHash(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewStatusHandler(testSecret, store)

	url := "/senangpay/query?order_id=order-status-missing-hash"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewStatusHandler(testSecret, store)

	req := httptest.NewRequest(http.MethodPost, "/senangpay/query", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestStatusHandlerStoreError(t *testing.T) {
	store := &errTransactionStore{getErr: errors.New("boom")}
	handler := NewStatusHandler(testSecret, store)

	orderID, hash := newTestStatusQuery("order-status-err")
	url := fmt.Sprintf("/senangpay/query?order_id=%s&hash=%s", orderID, hash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
