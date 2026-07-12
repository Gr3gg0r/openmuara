package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
)

func TestPayHandlerCreatesTransaction(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, []string{"test"})

	req := PaymentRequest{
		Provider:    "test",
		Type:        "charge",
		Amount:      99.99,
		Currency:    "USD",
		CustomerRef: "cust-1",
		Reference:   "ref-pay-1",
		Items:       []engine.TransactionItem{{ItemCode: "prod-1", Price: 99.99, Quantity: 1}},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp PaymentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Reference != "ref-pay-1" {
		t.Errorf("reference: want ref-pay-1, got %q", resp.Reference)
	}
	if resp.Status != engine.TransactionStatusNew {
		t.Errorf("status: want new, got %q", resp.Status)
	}

	stored, ok, err := store.GetByReference("ref-pay-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction in store")
	}
	if stored.Provider != "test" {
		t.Errorf("provider: want test, got %q", stored.Provider)
	}
}

func TestPayHandlerValidation(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, nil)

	req := PaymentRequest{Provider: "test"}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestPayHandlerIdempotency(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, []string{"test"})

	req := PaymentRequest{
		Provider:       "test",
		Type:           "charge",
		Amount:         10.0,
		Currency:       "USD",
		Reference:      "ref-idem-1",
		IdempotencyKey: "idem-pay-1",
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	httpReq.Header.Set("Idempotency-Key", "idem-pay-1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first status: want 201, got %d", rec.Code)
	}

	req.Reference = "ref-idem-2"
	body, _ = json.Marshal(req)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	httpReq2.Header.Set("Idempotency-Key", "idem-pay-1")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, httpReq2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second status: want 200, got %d", rec2.Code)
	}

	all, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(all))
	}
}

func TestGetPaymentHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "test", Type: "charge", Reference: "ref-get-1", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewGetPaymentHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/pay/ref-get-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp PaymentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Reference != "ref-get-1" {
		t.Errorf("reference: want ref-get-1, got %q", resp.Reference)
	}
}

func TestGetPaymentHandlerNotFound(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewGetPaymentHandler(store)
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/pay/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestPayHandlerRejectsInvalidCurrency(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, []string{"test"})

	cases := []struct {
		currency string
		wantOK   bool
	}{
		{"USD", true},
		{"EGP", true},
		{"MYR", true},
		{"usdx", false},
		{"US", false},
		{"US1", false},
		{"", false},
	}

	for _, tc := range cases {
		req := PaymentRequest{
			Provider:  "test",
			Type:      "charge",
			Amount:    10.0,
			Currency:  tc.currency,
			Reference: "ref-currency-" + tc.currency,
		}
		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httpReq)

		if tc.wantOK && rec.Code != http.StatusCreated {
			t.Errorf("currency %q: want 201, got %d", tc.currency, rec.Code)
		}
		if !tc.wantOK && rec.Code != http.StatusBadRequest {
			t.Errorf("currency %q: want 400, got %d", tc.currency, rec.Code)
		}
	}
}

func TestPayHandlerRejectsUnknownProvider(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, []string{"fawry", "stripe"})

	req := PaymentRequest{
		Provider:  "unknown",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "ref-provider-1",
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if id == "" {
		t.Fatal("expected non-empty id")
	}
	if id == GenerateID() {
		t.Fatal("expected unique ids")
	}
}

func TestRefundHandlerNotFound(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewRefundHandler(store)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/refund/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestRefundHandlerInvalidState(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "test", Type: "charge", Reference: "ref-refund-new", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewRefundHandler(store)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/refund/ref-refund-new", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusConflict {
		t.Errorf("status: want 409, got %d", rec.Code)
	}
}

func TestRefundHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Provider: "test", Type: "charge", Reference: "ref-refund-1", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusPaid}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewRefundHandler(store)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/refund/ref-refund-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp PaymentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != engine.TransactionStatusRefunded {
		t.Errorf("status: want refunded, got %q", resp.Status)
	}
}

func TestPayHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/v1/pay", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestPayHandlerInvalidJSON(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, nil)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader([]byte("{invalid")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestPayHandlerIdempotencyFromHeader(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewPayHandler(store, []string{"test"})

	req := PaymentRequest{
		Provider:  "test",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "ref-idem-header-1",
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	httpReq.Header.Set("Idempotency-Key", "idem-header-1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first status: want 201, got %d", rec.Code)
	}

	req.Reference = "ref-idem-header-2"
	body, _ = json.Marshal(req)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/v1/pay", bytes.NewReader(body))
	httpReq2.Header.Set("Idempotency-Key", "idem-header-1")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, httpReq2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second status: want 200, got %d", rec2.Code)
	}

	all, err := store.List(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(all))
	}
}

func TestGetPaymentHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewGetPaymentHandler(store)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/pay/ref-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestGetPaymentHandlerMissingReference(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewGetPaymentHandler(store)

	httpReq := httptest.NewRequest(http.MethodGet, "/v1/pay/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestRefundHandlerMethodNotAllowed(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewRefundHandler(store)

	httpReq := httptest.NewRequest(http.MethodGet, "/v1/refund/ref-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: want 405, got %d", rec.Code)
	}
}

func TestRefundHandlerMissingReference(t *testing.T) {
	store := engine.NewMemoryStore()
	handler := NewRefundHandler(store)

	httpReq := httptest.NewRequest(http.MethodPost, "/v1/refund/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", rec.Code)
	}
}

func TestValidatePaymentRequestMissingFields(t *testing.T) {
	allowed := map[string]bool{"stripe": true}
	valid := &PaymentRequest{
		Provider:  "stripe",
		Type:      "charge",
		Amount:    10.0,
		Currency:  "USD",
		Reference: "ref-1",
	}

	if err := validatePaymentRequest(valid, allowed); err != nil {
		t.Fatalf("valid request returned error: %v", err)
	}

	cases := []struct {
		name    string
		mutate  func(*PaymentRequest)
		wantErr string
	}{
		{"missing provider", func(r *PaymentRequest) { r.Provider = "" }, "provider is required"},
		{"provider not allowed", func(r *PaymentRequest) { r.Provider = "fawry" }, "provider is not supported"},
		{"missing type", func(r *PaymentRequest) { r.Type = "" }, "type is required"},
		{"zero amount", func(r *PaymentRequest) { r.Amount = 0 }, "amount must be greater than zero"},
		{"negative amount", func(r *PaymentRequest) { r.Amount = -1 }, "amount must be greater than zero"},
		{"missing currency", func(r *PaymentRequest) { r.Currency = "" }, "currency is required"},
		{"missing reference", func(r *PaymentRequest) { r.Reference = "" }, "reference is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := *valid
			tc.mutate(&req)
			err := validatePaymentRequest(&req, allowed)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if errcode.Message(err) != tc.wantErr {
				t.Errorf("message: want %q, got %q", tc.wantErr, errcode.Message(err))
			}
			var ec *errcode.Error
			if !errors.As(err, &ec) {
				t.Error("expected *errcode.Error")
			}
		})
	}
}
