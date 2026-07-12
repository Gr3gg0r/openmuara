package simple

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func validFawryConfig() plugin.GatewayConfig {
	return plugin.GatewayConfig{
		SchemaVersion: "v1",
		Metadata:      plugin.Metadata{Name: "test-fawry", Version: "1.0.0"},
		Runtime: &plugin.Runtime{
			Type: "simple",
			Simple: &plugin.SimpleRuntime{
				ChargeRoute:    "fawry_charge",
				WebhookEvent:   "payment.completed",
				Currency:       "EGP",
				ReferenceField: "merchantRefNum",
				AmountField:    "charge_items",
				CustomerField:  "customerProfileId",
				ResponseTemplate: map[string]any{
					"status":    "ok",
					"reference": "{{ .Reference }}",
				},
			},
		},
		Routes: []plugin.Route{
			{Path: "/fawry/charge", Method: "POST", Action: "fawry_charge", SchemaRef: "charge_request"},
			{Path: "/fawry/webhook", Method: "POST", Action: "fawry_v1_webhook"},
		},
		Schemas: plugin.Schemas{
			Requests: map[string]plugin.Schema{
				"charge_request": {Fields: []plugin.Field{
					{Name: "merchant_code", JSONName: "merchantCode", Type: "string", Required: true},
					{Name: "merchant_ref_num", JSONName: "merchantRefNum", Type: "string", Required: true},
					{Name: "return_url", JSONName: "returnUrl", Type: "string", Required: true},
					{Name: "charge_items", JSONName: "chargeItems", Type: "array", Required: true},
					{Name: "signature", JSONName: "signature", Type: "string", Required: true},
				}},
			},
		},
		Signature: &plugin.Signature{
			Algorithm: "fawry_sha256",
			Fields:    []string{"merchant_code", "merchant_ref_num", "customer_profile_id", "return_url", "charge_items", "signature"},
			SecretKey: "secret",
		},
		Webhooks: []plugin.Webhook{
			{Name: "payment", Event: "payment.completed", Method: "POST", Template: map[string]any{
				"reference": "{{ .Reference }}",
				"status":    "{{ .Status }}",
			}},
		},
	}
}

func TestProviderName(t *testing.T) {
	p := NewProvider(validFawryConfig())
	if got := p.Name(); got != "test-fawry" {
		t.Errorf("name: want test-fawry, got %q", got)
	}
}

func TestProviderInitRequiresSimpleRuntime(t *testing.T) {
	cfg := validFawryConfig()
	cfg.Runtime = nil
	p := NewProvider(cfg)
	if err := p.Init(map[string]any{"secret": "x"}); err == nil {
		t.Fatal("expected error when runtime is missing")
	}
}

func TestProviderRoutes(t *testing.T) {
	p := NewProvider(validFawryConfig())
	if err := p.Init(map[string]any{"secret": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	routes := p.Routes()
	if len(routes) != 2 {
		t.Fatalf("routes: want 2, got %d", len(routes))
	}
}

func TestChargeHandlerMissingField(t *testing.T) {
	p := NewProvider(validFawryConfig())
	if err := p.Init(map[string]any{"secret": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	body, _ := json.Marshal(map[string]any{"merchantCode": "x"})
	req := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestChargeHandlerInvalidSignature(t *testing.T) {
	p := NewProvider(validFawryConfig())
	if err := p.Init(map[string]any{"secret": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newFawryChargeRequest(t, "wrong-sig")
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestChargeHandlerSuccess(t *testing.T) {
	p := NewProvider(validFawryConfig())
	// #nosec G101 -- test fixture dummy credential, not a real secret
	if err := p.Init(map[string]any{"secret": "muara-fawry-secret"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := newFawryChargeRequest(t, "")
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["reference"] != "ref-123" {
		t.Errorf("reference: want ref-123, got %v", resp["reference"])
	}
	if resp["status"] != "ok" {
		t.Errorf("status: want ok, got %v", resp["status"])
	}

	tx, ok, err := p.store.GetByReference("ref-123")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok {
		t.Fatal("transaction not found")
	}
	if tx.Amount != 99.99 {
		t.Errorf("amount: want 99.99, got %f", tx.Amount)
	}
	if tx.Currency != "EGP" {
		t.Errorf("currency: want EGP, got %q", tx.Currency)
	}
}

func TestPayloadBuilder(t *testing.T) {
	p := NewProvider(validFawryConfig())
	if err := p.Init(map[string]any{"secret": "x"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	p.store = engine.NewMemoryStore()
	_, _, _ = p.store.CreateOrGet(engine.NewTransaction(engine.Transaction{
		Provider:  "test-fawry",
		Reference: "ref-123",
		Amount:    99.99,
		Currency:  "EGP",
		Status:    engine.TransactionStatusNew,
	}))

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{Reference: "ref-123", Status: "PAID"})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["reference"] != "ref-123" {
		t.Errorf("reference: want ref-123, got %v", got["reference"])
	}
	if got["status"] != "PAID" {
		t.Errorf("status: want PAID, got %v", got["status"])
	}
}

func newFawryChargeRequest(t *testing.T, overrideSig string) *http.Request {
	t.Helper()
	body := map[string]any{
		"merchantCode":      "muara-merchant-code",
		"merchantRefNum":    "ref-123",
		"customerProfileId": "user-123",
		"returnUrl":         "http://127.0.0.1:9999/callback",
		"chargeItems": []any{
			map[string]any{"itemId": "prod_test_123", "price": 99.99, "quantity": float64(1)},
		},
	}

	if overrideSig == "" {
		body["signature"] = computeFawrySignature(t, body, "muara-fawry-secret")
	} else {
		body["signature"] = overrideSig
	}

	b, _ := json.Marshal(body)
	return httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(b))
}
