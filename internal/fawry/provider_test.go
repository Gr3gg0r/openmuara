package fawry_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/fawry"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func validProviderConfig() map[string]any {
	// #nosec G101 -- test fixture dummy credentials
	return map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
	}
}

func validProviderConfigV2() map[string]any {
	cfg := validProviderConfig()
	cfg["version"] = "v2"
	return cfg
}

func TestProviderNameReturnsFawry(t *testing.T) {
	p := fawry.NewProvider()
	if got := p.Name(); got != "fawry" {
		t.Errorf("name: want fawry, got %q", got)
	}
}

func TestProviderInitWithValidConfigSucceeds(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestProviderInitMissingMerchantCode(t *testing.T) {
	p := fawry.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "merchant_code")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing merchant_code")
	}
}

func TestProviderInitMissingMerchantSecurityKey(t *testing.T) {
	p := fawry.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "merchant_security_key")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing merchant_security_key")
	}
}

func TestProviderInitMissingWebhookSecret(t *testing.T) {
	p := fawry.NewProvider()
	cfg := validProviderConfig()
	delete(cfg, "webhook_secret")
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for missing webhook_secret")
	}
}

func TestProviderInitUnknownVersion(t *testing.T) {
	p := fawry.NewProvider()
	cfg := validProviderConfig()
	cfg["version"] = "v3"
	if err := p.Init(cfg); err == nil {
		t.Fatal("expected error for unknown version")
	}
}

func TestProviderRoutes(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}

	routes := p.Routes()
	want := map[string]string{
		"POST /fawry/charge":        http.MethodPost,
		"GET /fawry/payment-status": http.MethodGet,
		"POST /fawry/webhook":       http.MethodPost,
		"POST /fawry/v1/charge":     http.MethodPost,
		"POST /fawry/v1/webhook":    http.MethodPost,
		"POST /fawry/v2/charge":     http.MethodPost,
		"POST /fawry/v2/webhook":    http.MethodPost,
		"GET /_admin/fawry-escape":  http.MethodGet,
		"POST /_admin/fawry-escape": http.MethodPost,
	}
	if len(routes) != len(want) {
		t.Fatalf("routes: want %d, got %d", len(want), len(routes))
	}
	got := make(map[string]string, len(routes))
	for _, r := range routes {
		key := r.Method + " " + r.Path
		got[key] = r.Method
		if r.Handler == nil {
			t.Errorf("route %s handler is nil", key)
		}
	}
	for key, method := range want {
		if got[key] != method {
			t.Errorf("route %s: want method %s, got %s", key, method, got[key])
		}
	}
}

func TestProviderPayloadBuilderDefaultV1(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newTestChargeRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("charge failed: status %d, body %s", rec.Code, rec.Body.String())
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{
		Reference: req.MerchantRefNum,
		Status:    string(webhook.PaymentStatusPaid),
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if got["merchantRefNumber"] != req.MerchantRefNum {
		t.Errorf("merchant ref: want %q, got %q", req.MerchantRefNum, got["merchantRefNumber"])
	}
	if got["orderStatus"] != string(webhook.PaymentStatusPaid) {
		t.Errorf("order status: want PAID, got %q", got["orderStatus"])
	}
	if got["messageSignature"] == "" {
		t.Error("message signature is empty")
	}
}

func TestProviderPayloadBuilderV2(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfigV2()); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newTestChargeRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httpReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("charge failed: status %d, body %s", rec.Code, rec.Body.String())
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{
		Reference: req.MerchantRefNum,
		Status:    string(webhook.PaymentStatusPaid),
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var got webhook.FawryV2Payload
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if got.MerchantRefNumber != req.MerchantRefNum {
		t.Errorf("merchant ref: want %q, got %q", req.MerchantRefNum, got.MerchantRefNumber)
	}
	if got.OrderStatus != string(webhook.PaymentStatusPaid) {
		t.Errorf("order status: want PAID, got %q", got.OrderStatus)
	}
	if got.MessageSignature == "" {
		t.Error("message signature is empty")
	}
	if got.PaymentAmount != 99.99 {
		t.Errorf("payment amount: want 99.99, got %f", got.PaymentAmount)
	}
	if got.CustomerMerchantID != req.CustomerProfileID {
		t.Errorf("customer merchant id: want %q, got %q", req.CustomerProfileID, got.CustomerMerchantID)
	}
	if len(got.OrderItems) != 1 {
		t.Fatalf("order items: want 1, got %d", len(got.OrderItems))
	}
	if got.OrderItems[0].ItemCode != "prod_test_123" {
		t.Errorf("item code: want prod_test_123, got %q", got.OrderItems[0].ItemCode)
	}
}

func TestProviderPayloadBuilderV1MissingTransaction(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: "missing", Status: "PAID"})
	if err != nil {
		t.Fatalf("v1 payload builder should not require stored transaction: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("expected non-empty v1 payload")
	}
}

func TestProviderPayloadBuilderV2MissingTransaction(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfigV2()); err != nil {
		t.Fatalf("init: %v", err)
	}

	builder := p.PayloadBuilder()
	_, err := builder(context.Background(), provider.Transaction{Reference: "missing"})
	if err == nil {
		t.Fatal("expected error for missing transaction")
	}
}

func TestProviderChargeHandlerMatchesExistingHandler(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	req := newTestChargeRequest()
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	p.ChargeHandler().ServeHTTP(rec, httpReq)

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
}

func TestProviderEscapeHandlerReturnsNonNil(t *testing.T) {
	p := fawry.NewProvider()
	if h := p.EscapeHandler(); h == nil {
		t.Fatal("expected non-nil escape handler")
	}
}

func TestProviderSetters(t *testing.T) {
	p := fawry.NewProvider()
	store := engine.NewMemoryStore()
	p.SetStore(store)

	d := webhook.NewDispatcherFromBuilder("http://127.0.0.1:1", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	p.SetDispatcher(d)

	if p.ChargeHandler() == nil {
		t.Fatal("expected charge handler after setters")
	}
}

func TestProviderWebhookHandlerReturnsNonNil(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	if h := p.WebhookHandler(); h == nil {
		t.Fatal("expected non-nil webhook handler")
	}
}

func TestProviderVerifyWebhookSignatureV1(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newTestChargeRequest()
	body, _ := json.Marshal(req)
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	tx := provider.Transaction{Reference: req.MerchantRefNum, Status: "PAID"}
	payload, err := p.PayloadBuilder()(context.Background(), tx)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	valid, err := p.VerifyWebhookSignature(payload, nil)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !valid {
		t.Error("expected v1 signature to be valid")
	}
}

func TestProviderVerifyWebhookSignatureV2(t *testing.T) {
	p := fawry.NewProvider()
	if err := p.Init(validProviderConfigV2()); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newTestChargeRequest()
	req.MerchantRefNum = "ref-v2-123"
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	tx := provider.Transaction{Reference: req.MerchantRefNum, Status: "PAID"}
	payload, err := p.PayloadBuilder()(context.Background(), tx)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	valid, err := p.VerifyWebhookSignature(payload, nil)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !valid {
		t.Error("expected v2 signature to be valid")
	}
}
