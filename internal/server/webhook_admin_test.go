package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/fawry"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

type errAttemptStore struct {
	webhook.AttemptStore
	listErr error
	getErr  error
	saveErr error
}

func (e *errAttemptStore) List(limit, offset int) ([]*webhook.Attempt, error) {
	if e.listErr != nil {
		return nil, e.listErr
	}
	return e.AttemptStore.List(limit, offset)
}

func (e *errAttemptStore) Get(ref string) (*webhook.Attempt, error) {
	if e.getErr != nil {
		return nil, e.getErr
	}
	return e.AttemptStore.Get(ref)
}

func (e *errAttemptStore) Save(a *webhook.Attempt) error {
	if e.saveErr != nil {
		return e.saveErr
	}
	return e.AttemptStore.Save(a)
}

func TestListWebhooksHandler(t *testing.T) {
	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{
		Ref:    "ref-1",
		URL:    "http://localhost/webhook",
		Status: webhook.AttemptStatusDelivered,
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks", nil)
	rec := httptest.NewRecorder()

	listWebhooksHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	attempts, ok := body["results"].([]any)
	if !ok || len(attempts) != 1 {
		t.Fatalf("want 1 attempt, got %+v", body)
	}
}

func TestListWebhooksHandlerWithFilter(t *testing.T) {
	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{
		Ref:          "ref-1",
		ProviderName: "fawry",
		URL:          "http://localhost/webhook",
		Status:       webhook.AttemptStatusDelivered,
	})
	_ = store.Save(&webhook.Attempt{
		Ref:          "ref-2",
		ProviderName: "stripe",
		URL:          "http://localhost/webhook",
		Status:       webhook.AttemptStatusFailed,
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks?provider=fawry&status=delivered", nil)
	rec := httptest.NewRecorder()

	listWebhooksHandler(store)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	attempts, ok := body["results"].([]any)
	if !ok || len(attempts) != 1 {
		t.Fatalf("want 1 fawry delivered attempt, got %+v", body)
	}
}

func TestInspectWebhookHandlerFound(t *testing.T) {
	store := webhook.NewMemoryStore()
	valid := true
	_ = store.Save(&webhook.Attempt{
		Ref:            "ref-1",
		ProviderName:   "stripe",
		URL:            "http://localhost/webhook",
		Status:         webhook.AttemptStatusDelivered,
		Payload:        []byte(`{"id":"evt_1"}`),
		Headers:        map[string]string{"Stripe-Signature": "t=1,v1=abc", "Content-Type": "application/json"},
		SignatureValid: &valid,
		History: []webhook.AttemptHistory{
			{Time: time.Now(), Status: 200, Error: ""},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks/ref-1", nil)
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()

	inspectWebhookHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	wh, ok := body["webhook"].(map[string]any)
	if !ok {
		t.Fatalf("expected webhook object, got %T", body["webhook"])
	}
	if wh["ref"] != "ref-1" {
		t.Errorf("reference mismatch: got %v", wh["ref"])
	}
	if wh["provider"] != "stripe" {
		t.Errorf("provider mismatch: got %v", wh["provider"])
	}

	headers, ok := wh["headers"].(map[string]any)
	if !ok {
		t.Fatalf("expected headers map, got %T", wh["headers"])
	}
	if headers["Stripe-Signature"] != "***" {
		t.Errorf("expected Stripe-Signature redacted, got %v", headers["Stripe-Signature"])
	}
	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type preserved, got %v", headers["Content-Type"])
	}

	if sv, ok := wh["signature_valid"].(bool); !ok || !sv {
		t.Errorf("expected signature_valid true, got %v", wh["signature_valid"])
	}

	attempts, ok := wh["attempt_events"].([]any)
	if !ok || len(attempts) != 1 {
		t.Errorf("expected 1 history entry, got %v", wh["attempt_events"])
	}
}

func TestInspectWebhookHandlerNotFound(t *testing.T) {
	store := webhook.NewMemoryStore()

	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks/missing", nil)
	req.SetPathValue("ref", "missing")
	rec := httptest.NewRecorder()

	inspectWebhookHandler(store)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestReplayWebhookHandler(t *testing.T) {
	var count atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := fawry.NewProvider()
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "secret",
		"webhook_secret":        "secret",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}

	txStore := engine.NewMemoryStore()
	if _, _, err := txStore.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	// Seed the provider's internal store as well so payload building succeeds.
	chargeReq := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-1",
		CustomerProfileID: "user-1",
		ReturnURL:         "http://localhost/callback",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod-1", Price: 100.0, Quantity: 1},
		},
	}
	chargeReq.Signature = fawry.Sign(chargeReq, "secret")
	body, _ := json.Marshal(chargeReq)
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	d := webhook.NewDispatcherFromProvider(ts.URL, 0, p)

	_, err := d.Dispatch(context.Background(), "ref-1", webhook.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/ref-1/replay")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()

	replayWebhookHandler(d, nil)(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status: want 202, got %d", rec.Code)
	}

	time.Sleep(200 * time.Millisecond)
	if count.Load() != 2 {
		t.Errorf("want 2 deliveries, got %d", count.Load())
	}
}

func TestReplayAllWebhooksHandler(t *testing.T) {
	var count atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := fawry.NewProvider()
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "secret",
		"webhook_secret":        "secret",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}

	txStore := engine.NewMemoryStore()
	if _, _, err := txStore.CreateOrGet(engine.Transaction{Reference: "ref-all", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	chargeReq := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "ref-all",
		CustomerProfileID: "user-1",
		ReturnURL:         "http://localhost/callback",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod-1", Price: 100.0, Quantity: 1},
		},
	}
	chargeReq.Signature = fawry.Sign(chargeReq, "secret")
	body, _ := json.Marshal(chargeReq)
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body)))

	d := webhook.NewDispatcherFromProvider(ts.URL, 0, p)
	_, err := d.Dispatch(context.Background(), "ref-all", webhook.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/replay-all")
	rec := httptest.NewRecorder()

	replayAllWebhooksHandler(d, nil)(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status: want 202, got %d", rec.Code)
	}

	time.Sleep(200 * time.Millisecond)
	if count.Load() != 2 {
		t.Errorf("want 2 deliveries, got %d", count.Load())
	}
}

func TestDeleteWebhookHandler(t *testing.T) {
	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{
		Ref:     "ref-1",
		URL:     "http://localhost/webhook",
		Status:  webhook.AttemptStatusDelivered,
		Payload: []byte("secret"),
	})

	req := newAdminRequest(http.MethodDelete, "/_admin/webhooks/ref-1")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()

	deleteWebhookHandler(store)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: want 204, got %d", rec.Code)
	}

	attempt, err := store.Get("ref-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if attempt == nil {
		t.Fatal("expected attempt to still exist")
	}
	if len(attempt.Payload) != 0 {
		t.Error("expected payload to be cleared")
	}
}

func TestListWebhooksHandlerStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks", nil)
	rec := httptest.NewRecorder()
	listWebhooksHandler(&errAttemptStore{listErr: errors.New("boom")})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestInspectWebhookHandlerStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks/ref-1", nil)
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()
	inspectWebhookHandler(&errAttemptStore{getErr: errors.New("boom")})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestReplayWebhookHandlerStoreError(t *testing.T) {
	d := &webhook.Dispatcher{Store: &errAttemptStore{getErr: errors.New("boom")}}
	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/ref-1/replay")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()
	replayWebhookHandler(d, nil)(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestReplayWebhookHandlerMissing(t *testing.T) {
	d := &webhook.Dispatcher{Store: webhook.NewMemoryStore()}
	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/ref-1/replay")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()
	replayWebhookHandler(d, nil)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestReplayWebhookHandlerWithProviderDispatcher(t *testing.T) {
	var count atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{
		Ref:          "ref-1",
		ProviderName: "stripe",
		URL:          ts.URL,
		Payload:      []byte("{}"),
	})

	stripeDisp := &webhook.Dispatcher{
		URL:    ts.URL,
		Store:  store,
		Worker: webhook.NewDeliveryWorker(),
	}

	d := &webhook.Dispatcher{Store: store}
	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/ref-1/replay")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()
	replayWebhookHandler(d, map[string]*webhook.Dispatcher{"stripe": stripeDisp})(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status: want 202, got %d", rec.Code)
	}

	time.Sleep(200 * time.Millisecond)
	if count.Load() != 1 {
		t.Errorf("expected stripe dispatcher delivery, got %d", count.Load())
	}
}

func TestReplayAllWebhooksHandlerStoreError(t *testing.T) {
	d := &webhook.Dispatcher{Store: &errAttemptStore{listErr: errors.New("boom")}}
	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/replay-all")
	rec := httptest.NewRecorder()
	replayAllWebhooksHandler(d, nil)(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestReplayAllWebhooksHandlerReplayError(t *testing.T) {
	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{Ref: "ref-1", ProviderName: "stripe", URL: "http://localhost/webhook", Payload: []byte("{}")})

	d := &webhook.Dispatcher{Store: store}
	stripeDisp := &webhook.Dispatcher{Store: &errAttemptStore{getErr: errors.New("boom")}}

	req := newAdminRequest(http.MethodPost, "/_admin/webhooks/replay-all")
	rec := httptest.NewRecorder()
	replayAllWebhooksHandler(d, map[string]*webhook.Dispatcher{"stripe": stripeDisp})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestDispatcherForAttempt(t *testing.T) {
	active := &webhook.Dispatcher{}
	mapped := &webhook.Dispatcher{}

	if d := dispatcherForAttempt(nil, active, nil); d != active {
		t.Error("expected active dispatcher for nil attempt")
	}
	if d := dispatcherForAttempt(&webhook.Attempt{ProviderName: ""}, active, nil); d != active {
		t.Error("expected active dispatcher for empty provider")
	}
	if d := dispatcherForAttempt(&webhook.Attempt{ProviderName: "stripe"}, active, map[string]*webhook.Dispatcher{"stripe": mapped}); d != mapped {
		t.Error("expected mapped dispatcher")
	}
	if d := dispatcherForAttempt(&webhook.Attempt{ProviderName: "missing"}, active, map[string]*webhook.Dispatcher{"stripe": mapped}); d != active {
		t.Error("expected active dispatcher fallback")
	}
}

func TestDeleteWebhookHandlerStoreErrors(t *testing.T) {
	req := newAdminRequest(http.MethodDelete, "/_admin/webhooks/ref-1")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()
	deleteWebhookHandler(&errAttemptStore{getErr: errors.New("boom")})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("get error status: want 500, got %d", rec.Code)
	}

	store := webhook.NewMemoryStore()
	_ = store.Save(&webhook.Attempt{Ref: "ref-1"})
	rec = httptest.NewRecorder()
	deleteWebhookHandler(&errAttemptStore{AttemptStore: store, saveErr: errors.New("boom")})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("save error status: want 500, got %d", rec.Code)
	}
}

func TestWebhookAdminHandlersNilDispatcher(t *testing.T) {
	mux := http.NewServeMux()
	WebhookAdminHandlers(mux, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/_admin/webhooks", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, ok := body["results"].([]any)
	if !ok || len(results) != 0 {
		t.Fatalf("expected empty results, got %v", body)
	}

	// Replay endpoints are not registered when there is no dispatcher.
	req = httptest.NewRequest(http.MethodPost, "/_admin/webhooks/ref/replay", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("replay status: want 404, got %d", rec.Code)
	}
}

func TestReplayWebhookHandlerDeniesViewer(t *testing.T) {
	d := &webhook.Dispatcher{Store: webhook.NewMemoryStore()}
	req := newViewerRequest(http.MethodPost, "/_admin/webhooks/ref-1/replay")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()

	replayWebhookHandler(d, nil)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: want 403, got %d", rec.Code)
	}
}

func TestReplayAllWebhooksHandlerDeniesViewer(t *testing.T) {
	d := &webhook.Dispatcher{Store: webhook.NewMemoryStore()}
	req := newViewerRequest(http.MethodPost, "/_admin/webhooks/replay-all")
	rec := httptest.NewRecorder()

	replayAllWebhooksHandler(d, nil)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: want 403, got %d", rec.Code)
	}
}

func TestDeleteWebhookHandlerDeniesViewer(t *testing.T) {
	store := webhook.NewMemoryStore()
	req := newViewerRequest(http.MethodDelete, "/_admin/webhooks/ref-1")
	req.SetPathValue("ref", "ref-1")
	rec := httptest.NewRecorder()

	deleteWebhookHandler(store)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: want 403, got %d", rec.Code)
	}
}
