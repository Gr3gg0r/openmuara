package fawry_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/fawry"
	"github.com/Gr3gg0r/openmuara/internal/server"
)

func newTestChargeRequest() fawry.ChargeRequest {
	req := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerEmail:     "test@example.com",
		CustomerName:      "Test User",
		CustomerProfileID: "user-456",
		PaymentExpiry:     1234567890000,
		Language:          "en-gb",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod_test_123", Price: 99.99, Quantity: 1},
		},
		ReturnURL: "http://127.0.0.1:9999/callback",
	}
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	return req
}

func TestChargeHandlerSuccess(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := server.RequestIDMiddleware(fawry.NewChargeHandler("muara-fawry-secret", store))

	req := newTestChargeRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp fawry.ChargeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status: want ok, got %q", resp.Status)
	}
	if resp.Reference != req.MerchantRefNum {
		t.Errorf("reference: want %q, got %q", req.MerchantRefNum, resp.Reference)
	}
	if rec.Header().Get("X-Trace-Id") == "" {
		t.Error("X-Trace-Id header missing")
	}

	tx, ok, err := store.GetByReference(req.MerchantRefNum)

	if err != nil {

		t.Fatalf("lookup transaction: %v", err)

	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != "fawry" {
		t.Errorf("provider: want fawry, got %q", tx.Provider)
	}
	if tx.Type != "charge" {
		t.Errorf("type: want charge, got %q", tx.Type)
	}
	if tx.Status != engine.TransactionStatusNew {
		t.Errorf("status: want %q, got %q", engine.TransactionStatusNew, tx.Status)
	}
	if tx.Amount != 99.99 {
		t.Errorf("amount: want 99.99, got %f", tx.Amount)
	}
	if tx.CustomerRef != req.CustomerProfileID {
		t.Errorf("customerRef: want %q, got %q", req.CustomerProfileID, tx.CustomerRef)
	}
	if len(tx.Items) != 1 {
		t.Fatalf("items: want 1, got %d", len(tx.Items))
	}
	if tx.Items[0].ItemCode != "prod_test_123" {
		t.Errorf("item code: want prod_test_123, got %q", tx.Items[0].ItemCode)
	}
}

func TestChargeHandlerIdempotency(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := server.RequestIDMiddleware(fawry.NewChargeHandler("muara-fawry-secret", store))

	req := newTestChargeRequest()
	body, _ := json.Marshal(req)

	httpReq1 := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq1.Header.Set("Content-Type", "application/json")
	httpReq1.Header.Set("Idempotency-Key", "idem-charge-1")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, httpReq1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request status: want 200, got %d, body: %s", rec1.Code, rec1.Body.String())
	}

	// Second request with the same idempotency key but a different reference.
	req2 := req
	req2.MerchantRefNum = "ref-456"
	req2.Signature = fawry.Sign(req2, "muara-fawry-secret")
	body2, _ := json.Marshal(req2)

	httpReq2 := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body2))
	httpReq2.Header.Set("Content-Type", "application/json")
	httpReq2.Header.Set("Idempotency-Key", "idem-charge-1")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, httpReq2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second request status: want 200, got %d, body: %s", rec2.Code, rec2.Body.String())
	}

	// The ledger should contain exactly one transaction.
	all, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(all))
	}
	if all[0].Reference != req.MerchantRefNum {
		t.Errorf("expected stored transaction to keep original reference %q, got %q", req.MerchantRefNum, all[0].Reference)
	}
}

func TestChargeHandlerInvalidSignature(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := server.RequestIDMiddleware(fawry.NewChargeHandler("muara-fawry-secret", store))

	req := newTestChargeRequest()
	req.Signature = "invalid"
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

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
	if resp.Error.TraceID == "" {
		t.Error("trace_id missing")
	}
}

func TestChargeHandlerMissingMerchantRefNum(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := server.RequestIDMiddleware(fawry.NewChargeHandler("muara-fawry-secret", store))

	req := newTestChargeRequest()
	req.MerchantRefNum = ""
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp server.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Error.Code != server.ErrMissingField {
		t.Errorf("code: want %q, got %q", server.ErrMissingField, resp.Error.Code)
	}
}

func TestChargeHandlerInvalidJSON(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := server.RequestIDMiddleware(fawry.NewChargeHandler("muara-fawry-secret", store))

	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader([]byte("not json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp server.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Error.Code != server.ErrInvalidJSON {
		t.Errorf("code: want %q, got %q", server.ErrInvalidJSON, resp.Error.Code)
	}
}

func TestChargeHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := httptest.NewRequest(http.MethodGet, "/fawry/charge", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestChargeHandlerMissingMerchantCode(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := newTestChargeRequest()
	req.MerchantCode = ""
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestChargeHandlerMissingChargeItems(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := newTestChargeRequest()
	req.ChargeItems = nil
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestChargeHandlerMissingItemID(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := newTestChargeRequest()
	req.ChargeItems[0].ItemID = ""
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestChargeHandlerMissingReturnURL(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := newTestChargeRequest()
	req.ReturnURL = ""
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestChargeHandlerMissingSignature(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewChargeHandler("muara-fawry-secret", store)

	req := newTestChargeRequest()
	req.Signature = ""
	body, _ := json.Marshal(req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
