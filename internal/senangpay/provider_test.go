package senangpay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

type errTransactionStore struct {
	engine.TransactionStore
	getErr    error
	createErr error
	listErr   error
}

func (e *errTransactionStore) GetByReference(ref string) (engine.Transaction, bool, error) {
	if e.getErr != nil {
		return engine.Transaction{}, false, e.getErr
	}
	return e.TransactionStore.GetByReference(ref)
}

func (e *errTransactionStore) CreateOrGet(tx engine.Transaction) (engine.Transaction, bool, error) {
	if e.createErr != nil {
		return engine.Transaction{}, false, e.createErr
	}
	return e.TransactionStore.CreateOrGet(tx)
}

func (e *errTransactionStore) List(limit, offset int) ([]engine.Transaction, error) {
	if e.listErr != nil {
		return nil, e.listErr
	}
	return e.TransactionStore.List(limit, offset)
}

func TestProviderInit(t *testing.T) {
	p := NewProvider()
	if err := p.Init(map[string]any{"secret_key": "secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := p.Init(map[string]any{}); err == nil {
		t.Fatal("expected error for missing secret_key")
	}
}

func TestChargeHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewChargeHandler("secret", store, "http://localhost")

	req := ChargeRequest{Detail: "Test", Amount: 10.5, OrderID: "order-1", Name: "A", Email: "a@example.com", Phone: "012"}
	SignRequest(&req, "secret")
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp ChargeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.OrderID != "order-1" {
		t.Errorf("order_id: want order-1, got %q", resp.OrderID)
	}

	tx, ok, err := store.GetByReference("order-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction")
	}
	if tx.Status != engine.TransactionStatusNew {
		t.Errorf("status: want new, got %q", tx.Status)
	}
}

func TestChargeHandlerInvalidHash(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewChargeHandler("secret", store, "http://localhost")

	req := ChargeRequest{Detail: "Test", Amount: 10.0, OrderID: "order-2", Hash: "bad"}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestCallbackHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-3", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCallbackHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=1&order_id=order-3&transaction_id=tx-1&msg=ok", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	tx, ok, err := store.GetByReference("order-3")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusPaid {
		t.Errorf("status: want paid, got %q", tx.Status)
	}
}

func TestCallbackHandlerReturns404ForMissingOrder(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewCallbackHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=1&order_id=missing&transaction_id=tx-1&msg=ok", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestWebhookHandlerReturns404ForMissingOrder(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewWebhookHandler(store)
	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/webhook?status_id=1&order_id=missing&transaction_id=tx-1&msg=ok", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestProviderRoutes(t *testing.T) {
	p := NewProvider()
	routes := p.Routes()
	if len(routes) != 4 {
		t.Fatalf("routes: want 4, got %d", len(routes))
	}
}

func TestProviderInterfaceMethods(t *testing.T) {
	p := NewProvider()
	if p.Name() != ProviderName {
		t.Errorf("name: want %q, got %q", ProviderName, p.Name())
	}
	if p.ChargeHandler() == nil {
		t.Error("expected ChargeHandler to be non-nil")
	}
	if p.WebhookHandler() == nil {
		t.Error("expected WebhookHandler to be non-nil")
	}
	if p.EscapeHandler() != nil {
		t.Error("expected EscapeHandler to be nil")
	}
	if p.PayloadBuilder() == nil {
		t.Error("expected PayloadBuilder to be non-nil")
	}

	p.SetBaseURL("http://localhost")
	p.SetStore(engine.NewMemoryStore())

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: "ref-1", Status: "paid"})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}
	if !strings.Contains(string(payload), `"provider":"senangpay"`) {
		t.Errorf("payload missing provider: %s", payload)
	}
}

func TestChargeHandlerValidationErrors(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewChargeHandler("secret", store, "http://localhost")

	cases := []struct {
		name string
		req  ChargeRequest
	}{
		{"missing detail", ChargeRequest{Amount: 10, OrderID: "o1", Hash: "h"}},
		{"missing amount", ChargeRequest{Detail: "Test", OrderID: "o1", Hash: "h"}},
		{"missing order_id", ChargeRequest{Detail: "Test", Amount: 10, Hash: "h"}},
		{"missing hash", ChargeRequest{Detail: "Test", Amount: 10, OrderID: "o1"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.req)
			httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httpReq)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status: want 400, got %d", rec.Code)
			}
		})
	}
}

func TestChargeHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewChargeHandler("secret", store, "http://localhost")
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/charge", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestChargeHandlerInvalidJSON(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewChargeHandler("secret", store, "http://localhost")
	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestCallbackHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewCallbackHandler(store)
	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/callback?status_id=1&order_id=o1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestWebhookHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewWebhookHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/webhook?status_id=1&order_id=o1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestApplyCallbackIgnoresUnknownStatus(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-status", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCallbackHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=99&order_id=order-status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", rec.Code)
	}
}

func TestChargeHandlerStoreError(t *testing.T) {
	handler := NewChargeHandler("secret", &errTransactionStore{createErr: errors.New("boom")}, "http://localhost")
	req := ChargeRequest{Detail: "Test", Amount: 10.0, OrderID: "order-se", Name: "A", Email: "a@example.com", Phone: "012"}
	SignRequest(&req, "secret")
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestCallbackHandlerStoreError(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-err", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCallbackHandler(&errTransactionStore{TransactionStore: store, getErr: errors.New("boom")})
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=1&order_id=order-err", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestWebhookHandlerStoreError(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-wh", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewWebhookHandler(&errTransactionStore{TransactionStore: store, getErr: errors.New("boom")})
	httpReq := httptest.NewRequest(http.MethodPost, "/senangpay/webhook?status_id=1&order_id=order-wh", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestCallbackHandlerTransitionError(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-conflict", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusPaid}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCallbackHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=0&order_id=order-conflict", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestApplyCallbackCreateOrGetError(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "charge", Reference: "order-create-err", Amount: 10.0, Currency: "MYR", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCallbackHandler(&errTransactionStore{TransactionStore: store, createErr: errors.New("boom")})
	httpReq := httptest.NewRequest(http.MethodGet, "/senangpay/callback?status_id=1&order_id=order-create-err", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}
