package fawry_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/fawry"
	"github.com/openmuara/openmuara/internal/server"
)

const (
	testMerchantCode        = "muara-merchant-code"
	testMerchantSecurityKey = "muara-fawry-secret"
	testReference           = "ref-status-1"
)

func newTestStatusQuery(ref string) fawry.StatusQuery {
	q := fawry.StatusQuery{
		MerchantCode:   testMerchantCode,
		MerchantRefNum: ref,
	}
	q.Signature = fawry.SignStatusQuery(q, testMerchantSecurityKey)
	return q
}

func setupStatusTestStore(t *testing.T) engine.TransactionStore {
	t.Helper()
	store := engine.NewMemoryStore()
	tx := engine.NewTransaction(engine.Transaction{
		Provider:  "fawry",
		Type:      "charge",
		Amount:    99.99,
		Currency:  "EGP",
		Status:    engine.TransactionStatusPaid,
		Reference: testReference,
	})
	if _, _, err := store.CreateOrGet(tx); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}
	return store
}

func TestStatusHandlerSuccess(t *testing.T) {
	store := setupStatusTestStore(t)
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	q := newTestStatusQuery(testReference)
	url := fmt.Sprintf("/fawry/payment-status?merchantCode=%s&merchantRefNum=%s&signature=%s", q.MerchantCode, q.MerchantRefNum, q.Signature)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp fawry.PaymentStatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Status != "PAID" {
		t.Errorf("status: want PAID, got %q", resp.Status)
	}
	if resp.Reference != testReference {
		t.Errorf("reference: want %q, got %q", testReference, resp.Reference)
	}
	if resp.Amount != 99.99 {
		t.Errorf("amount: want 99.99, got %f", resp.Amount)
	}
	if resp.Currency != "EGP" {
		t.Errorf("currency: want EGP, got %q", resp.Currency)
	}
}

func TestStatusHandlerNotFound(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	q := newTestStatusQuery("missing-ref")
	url := fmt.Sprintf("/fawry/payment-status?merchantCode=%s&merchantRefNum=%s&signature=%s", q.MerchantCode, q.MerchantRefNum, q.Signature)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerInvalidSignature(t *testing.T) {
	store := setupStatusTestStore(t)
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	url := fmt.Sprintf("/fawry/payment-status?merchantCode=%s&merchantRefNum=%s&signature=invalid", testMerchantCode, testReference)
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

func TestStatusHandlerMissingMerchantRefNum(t *testing.T) {
	store := setupStatusTestStore(t)
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	url := fmt.Sprintf("/fawry/payment-status?merchantCode=%s", testMerchantCode)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerInvalidMerchantCode(t *testing.T) {
	store := setupStatusTestStore(t)
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	q := newTestStatusQuery(testReference)
	url := fmt.Sprintf("/fawry/payment-status?merchantCode=wrong-code&merchantRefNum=%s&signature=%s", q.MerchantRefNum, q.Signature)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)

	req := httptest.NewRequest(http.MethodPost, "/fawry/payment-status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestStatusHandlerMapsNewStatus(t *testing.T) {
	store := engine.NewMemoryStore()
	tx := engine.NewTransaction(engine.Transaction{
		Provider:  "fawry",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "EGP",
		Status:    engine.TransactionStatusNew,
		Reference: testReference,
	})
	if _, _, err := store.CreateOrGet(tx); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	handler := fawry.NewStatusHandler(testMerchantCode, testMerchantSecurityKey, store)
	q := newTestStatusQuery(testReference)
	url := fmt.Sprintf("/fawry/payment-status?merchantCode=%s&merchantRefNum=%s&signature=%s", q.MerchantCode, q.MerchantRefNum, q.Signature)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp fawry.PaymentStatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Status != "NEW" {
		t.Errorf("status: want NEW, got %q", resp.Status)
	}
}
