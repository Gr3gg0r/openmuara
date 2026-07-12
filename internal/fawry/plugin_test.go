package fawry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/plugin"
)

// Given a Fawry plugin with default config, When it is registered, Then all four routes are available.
func TestPluginRegistersAllRoutes(t *testing.T) {
	impl := NewPlugin("secret", "whsecret", engine.NewMemoryStore(), nil)
	reg := plugin.NewRegistry()
	if err := reg.Register(context.Background(), "fawry", plugin.GatewayConfig{}, impl); err != nil {
		t.Fatalf("register: %v", err)
	}

	for _, action := range []string{"fawry_charge", "fawry_v1_webhook", "fawry_v2_webhook", "fawry_escape_page", "fawry_escape_action"} {
		if _, ok := reg.Handler(action); !ok {
			t.Errorf("missing handler for action %q", action)
		}
	}
}

// Given a signed charge request, When posted to the plugin route, Then it returns status ok.
func TestPluginChargeRouteValidSignature(t *testing.T) {
	// #nosec G101 -- test fixture dummy secret
	const secret = "muara-fawry-secret"
	impl := NewPlugin(secret, "whsecret", engine.NewMemoryStore(), nil)
	reg := plugin.NewRegistry()
	if err := reg.Register(context.Background(), "fawry", plugin.GatewayConfig{}, impl); err != nil {
		t.Fatalf("register: %v", err)
	}

	factory, ok := reg.Handler("fawry_charge")
	if !ok {
		t.Fatal("fawry_charge handler not found")
	}
	handler, err := factory(plugin.Dependencies{})
	if err != nil {
		t.Fatalf("build handler: %v", err)
	}

	reqBody := ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-123",
		CustomerProfileID: "user-123",
		ReturnURL:         "http://127.0.0.1:9999/callback",
		ChargeItems:       []ChargeItem{{ItemID: "prod_test_123", Price: 99.99, Quantity: 1}},
	}
	reqBody.Signature = Sign(reqBody, secret)
	body, _ := json.Marshal(reqBody)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp ChargeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status: want ok, got %q", resp.Status)
	}
}

// Given an invalid signature, When posted to the plugin route, Then it returns 400 with invalid signature code.
func TestPluginChargeRouteInvalidSignature(t *testing.T) {
	impl := NewPlugin("secret", "whsecret", engine.NewMemoryStore(), nil)
	reg := plugin.NewRegistry()
	if err := reg.Register(context.Background(), "fawry", plugin.GatewayConfig{}, impl); err != nil {
		t.Fatalf("register: %v", err)
	}

	factory, ok := reg.Handler("fawry_charge")
	if !ok {
		t.Fatal("fawry_charge handler not found")
	}
	handler, err := factory(plugin.Dependencies{})
	if err != nil {
		t.Fatalf("build handler: %v", err)
	}

	body, _ := json.Marshal(ChargeRequest{
		MerchantCode:   "muara-merchant-code",
		MerchantRefNum: "ref-123",
		ReturnURL:      "http://127.0.0.1:9999/callback",
		ChargeItems:    []ChargeItem{{ItemID: "prod_test_123", Price: 99.99, Quantity: 1}},
		Signature:      "invalid",
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`OPENMUARA_INVALID_SIGNATURE`)) {
		t.Errorf("expected invalid signature code, got %s", rec.Body.Bytes())
	}
}

func TestPluginRegistrationError(t *testing.T) {
	impl := NewPlugin("secret", "whsecret", engine.NewMemoryStore(), nil)
	reg := plugin.NewRegistry()
	if err := reg.Register(context.Background(), "fawry", plugin.GatewayConfig{}, impl); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if err := reg.Register(context.Background(), "fawry", plugin.GatewayConfig{}, impl); err == nil {
		t.Fatal("expected error registering duplicate fawry plugin")
	}
}
