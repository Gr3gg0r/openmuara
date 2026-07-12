package defaultplugin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestDefaultProviderName(t *testing.T) {
	// Given the default provider
	p := NewProvider()

	// When Name() is called
	got := p.Name()

	// Then it returns "default"
	if got != "default" {
		t.Errorf("name: want default, got %q", got)
	}
}

func TestDefaultProviderInitSucceeds(t *testing.T) {
	// Given the default provider and any config
	p := NewProvider()
	cfg := map[string]any{"foo": "bar"}

	// When Init() is called
	err := p.Init(cfg)

	// Then it succeeds
	if err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestDefaultProviderInitSucceedsWithEmptyConfig(t *testing.T) {
	// Given the default provider and no config
	p := NewProvider()

	// When Init() is called with nil config
	err := p.Init(nil)

	// Then it succeeds
	if err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestDefaultProviderRoutes(t *testing.T) {
	// Given the default provider
	p := NewProvider()

	// When Routes() is called
	routes := p.Routes()

	// Then it returns 3 routes with expected paths
	want := map[string]string{
		"/default":         http.MethodGet,
		"/default/charge":  http.MethodPost,
		"/default/webhook": http.MethodPost,
	}
	if len(routes) != len(want) {
		t.Fatalf("routes: want %d, got %d", len(want), len(routes))
	}
	for _, r := range routes {
		method, ok := want[r.Path]
		if !ok {
			t.Errorf("unexpected route %s %s", r.Method, r.Path)
			continue
		}
		if r.Method != method {
			t.Errorf("route %s method: want %s, got %s", r.Path, method, r.Method)
		}
		if r.Handler == nil {
			t.Errorf("route %s handler is nil", r.Path)
		}
	}
}

func TestDefaultProviderChargeHandler(t *testing.T) {
	// Given a POST to /default/charge
	p := NewProvider()
	req := httptest.NewRequest(http.MethodPost, "/default/charge", nil)
	rec := httptest.NewRecorder()

	// When it is handled
	p.ChargeHandler().ServeHTTP(rec, req)

	// Then it returns 200 JSON with provider: default
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("content-type: want application/json, got %q", ct)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["provider"] != "default" {
		t.Errorf("provider: want default, got %v", body["provider"])
	}
	if body["status"] != "success" {
		t.Errorf("status: want success, got %v", body["status"])
	}
	if body["transaction_id"] == "" {
		t.Error("transaction_id is empty")
	}
}

func TestDefaultProviderWebhookHandler(t *testing.T) {
	// Given a POST to /default/webhook
	p := NewProvider()
	req := httptest.NewRequest(http.MethodPost, "/default/webhook", nil)
	rec := httptest.NewRecorder()

	// When it is handled
	p.WebhookHandler().ServeHTTP(rec, req)

	// Then it returns 200 with status: acknowledged
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "acknowledged" {
		t.Errorf("status: want acknowledged, got %q", body["status"])
	}
}

func TestDefaultProviderPayloadBuilder(t *testing.T) {
	// Given a transaction
	p := NewProvider()
	builder := p.PayloadBuilder()

	// When PayloadBuilder() is called
	payload, err := builder(context.Background(), provider.Transaction{
		Reference: "ref-123",
		Status:    "PAID",
	})

	// Then it returns valid JSON
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if body["provider"] != "default" {
		t.Errorf("provider: want default, got %v", body["provider"])
	}
	if body["reference"] != "ref-123" {
		t.Errorf("reference: want ref-123, got %v", body["reference"])
	}
	if body["status"] != "PAID" {
		t.Errorf("status: want PAID, got %v", body["status"])
	}
}

func TestDefaultProviderEscapeHandler(t *testing.T) {
	// Given the default provider
	p := NewProvider()

	// When EscapeHandler() is called
	h := p.EscapeHandler()

	// Then it returns nil
	if h != nil {
		t.Fatal("expected nil escape handler")
	}
}

func TestDefaultProviderStatusHandler(t *testing.T) {
	// Given a GET to /default
	p := NewProvider()
	req := httptest.NewRequest(http.MethodGet, "/default", nil)
	rec := httptest.NewRecorder()

	// When it is handled
	routes := p.Routes()
	var statusHandler http.Handler
	for _, r := range routes {
		if r.Path == "/default" {
			statusHandler = r.Handler
			break
		}
	}
	if statusHandler == nil {
		t.Fatal("status route not found")
	}
	statusHandler.ServeHTTP(rec, req)

	// Then it returns 200 with provider: default and status: ok
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["provider"] != "default" {
		t.Errorf("provider: want default, got %q", body["provider"])
	}
	if body["status"] != "ok" {
		t.Errorf("status: want ok, got %q", body["status"])
	}
}
