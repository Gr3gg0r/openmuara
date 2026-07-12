package webhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/fawry"
	"github.com/openmuara/openmuara/internal/provider/defaultplugin"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestDispatcherFromFawryProviderProducesFawryV2Payload(t *testing.T) {
	// Given a dispatcher built from the Fawry provider with a seeded transaction
	p := fawry.NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
		"version":               "v2",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		received.Store(true)
		body, _ := io.ReadAll(req.Body)
		if len(body) == 0 {
			t.Error("expected non-empty body")
		}
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("unmarshal payload: %v", err)
		}
		if payload["fawryRefNumber"] == nil {
			t.Error("expected fawryRefNumber in payload")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Seed a transaction via the provider's charge handler.
	req := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "disp-ref-1",
		CustomerProfileID: "user-1",
		ReturnURL:         "http://localhost/callback",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod-1", Price: 25.00, Quantity: 2},
		},
	}
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)
	chargeReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", nil)
	chargeReq.Body = io.NopCloser(bytes.NewReader(body))
	chargeReq.Header.Set("Content-Type", "application/json")
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), chargeReq)

	d := webhook.NewDispatcherFromProvider(ts.URL, 0, p)

	// When Dispatch is called
	_, err := d.Dispatch(context.Background(), "disp-ref-1", webhook.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	// Then it produces a Fawry V2-shaped payload and delivers it.
	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("server did not receive webhook")
	}
}

func TestDispatcherFromDefaultProviderProducesDefaultPayload(t *testing.T) {
	// Given a dispatcher built from the default provider
	p := defaultplugin.NewProvider()

	var received atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		received.Store(true)
		body, _ := io.ReadAll(req.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("unmarshal payload: %v", err)
		}
		if payload["provider"] != "default" {
			t.Errorf("provider: want default, got %v", payload["provider"])
		}
		if payload["reference"] != "ref-1" {
			t.Errorf("reference: want ref-1, got %v", payload["reference"])
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	d := webhook.NewDispatcherFromProvider(ts.URL, 0, p)

	// When Dispatch is called
	_, err := d.Dispatch(context.Background(), "ref-1", webhook.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	// Then it produces the default payload shape.
	time.Sleep(200 * time.Millisecond)
	if !received.Load() {
		t.Error("server did not receive webhook")
	}
}
