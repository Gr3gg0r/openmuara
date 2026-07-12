package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	_ "github.com/Gr3gg0r/openmuara/internal/fawry" // register fawry factory for health/version tests
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/defaultplugin"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

type errTransactionStore struct {
	engine.TransactionStore
	listErr error
}

func (e *errTransactionStore) List(_, _ int) ([]engine.Transaction, error) {
	return nil, e.listErr
}

func TestListTransactionsHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions", nil)
	rec := httptest.NewRecorder()

	listTransactionsHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	txs, ok := body["results"].([]any)
	if !ok || len(txs) != 1 {
		t.Fatalf("want 1 transaction, got %+v", body)
	}
}

func TestListTransactionsHandlerPagination(t *testing.T) {
	store := engine.NewMemoryStore()
	for i := 1; i <= 3; i++ {
		if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-" + string(rune('a'+i-1)), Amount: float64(i), Currency: "USD"}); err != nil {
			t.Fatalf("seed transaction: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions?limit=1&offset=1", nil)
	rec := httptest.NewRecorder()

	listTransactionsHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["limit"] != 1.0 || body["offset"] != 1.0 {
		t.Errorf("want limit=1 offset=1, got %+v", body)
	}
	txs, ok := body["results"].([]any)
	if !ok || len(txs) != 1 {
		t.Fatalf("want 1 transaction, got %+v", body)
	}
}

func TestListProvidersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/providers", nil)
	rec := httptest.NewRecorder()

	listProvidersHandler(RouterConfig{ActiveProvider: "default", EnabledProviders: []string{"default"}, Host: "127.0.0.1", Port: 9000})(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["active"] != "default" {
		t.Errorf("active: want default, got %v", body["active"])
	}

	enabled, ok := body["enabled"].([]any)
	if !ok || len(enabled) != 1 {
		t.Errorf("enabled: want 1 entry, got %v", body["enabled"])
	}

	available, ok := body["available"].([]any)
	if !ok || len(available) == 0 {
		t.Errorf("available: want non-empty list, got %v", body["available"])
	}

	details, ok := body["providers"].(map[string]any)
	if !ok {
		t.Fatalf("providers: want details map, got %v", body["providers"])
	}
	info, ok := details["default"].(map[string]any)
	if !ok {
		t.Fatalf("default provider details missing: got %v", details)
	}
	if info["enabled"] != true {
		t.Errorf("default.enabled: want true, got %v", info["enabled"])
	}
	if info["active"] != true {
		t.Errorf("default.active: want true, got %v", info["active"])
	}
}

func TestOnboardingHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	dispatcher := webhook.NewDispatcherFromBuilder("http://localhost/hook", 3, nil, nil)
	dispatcher.Store = webhook.NewMemoryStore()

	req := httptest.NewRequest(http.MethodGet, "/_admin/onboarding", nil)
	rec := httptest.NewRecorder()

	onboardingHandler("fawry", []string{"fawry"}, store, dispatcher)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["server_ready"] != true {
		t.Errorf("server_ready want true, got %v", body["server_ready"])
	}
	if body["providers_enabled"] != true {
		t.Errorf("providers_enabled want true, got %v", body["providers_enabled"])
	}
	if body["first_transaction"] != false {
		t.Errorf("first_transaction want false, got %v", body["first_transaction"])
	}
	if body["first_webhook_received"] != false {
		t.Errorf("first_webhook_received want false, got %v", body["first_webhook_received"])
	}
	if body["webhooks_enabled"] != true {
		t.Errorf("webhooks_enabled want true, got %v", body["webhooks_enabled"])
	}
	if body["active_provider"] != "fawry" {
		t.Errorf("active_provider want fawry, got %v", body["active_provider"])
	}
	next, ok := body["next_step"].(map[string]any)
	if !ok {
		t.Fatalf("next_step missing: %v", body["next_step"])
	}
	if next["route"] != "/fawry/charge" {
		t.Errorf("next_step.route want /fawry/charge, got %v", next["route"])
	}
}

func TestOnboardingHandlerWithEvents(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	dispatcher := webhook.NewDispatcherFromBuilder("http://localhost/hook", 3,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)
	dispatcher.Store = webhook.NewMemoryStore()
	if _, err := dispatcher.Dispatch(t.Context(), "tx-1", webhook.PaymentStatusPaid); err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/onboarding", nil)
	rec := httptest.NewRecorder()

	onboardingHandler("default", []string{"default"}, store, dispatcher)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["first_transaction"] != true {
		t.Errorf("first_transaction want true, got %v", body["first_transaction"])
	}
	if body["first_webhook_received"] != true {
		t.Errorf("first_webhook_received want true, got %v", body["first_webhook_received"])
	}
	if body["webhooks_enabled"] != true {
		t.Errorf("webhooks_enabled want true, got %v", body["webhooks_enabled"])
	}
}

func TestAdminAPIHandlersRegisterEndpoints(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{TransactionStore: store})

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/_admin/transactions"},
		{http.MethodGet, "/_admin/transactions/tx-1"},
		{http.MethodGet, "/_admin/ledger"},
		{http.MethodGet, "/_admin/providers"},
		{http.MethodGet, "/_admin/onboarding"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		req.SetPathValue("ref", "tx-1")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("endpoint %s %s was not registered", tc.method, tc.path)
		}
	}
}

var _ = defaultplugin.NewProvider
var _ = provider.Names

func TestListTransactionsHandlerFilters(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "fawry", Status: engine.TransactionStatusPaid, Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-2", Provider: "stripe", Status: engine.TransactionStatusUnpaid, Amount: 20.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	cases := []struct {
		query     string
		wantRefs  []string
		wantTotal float64
	}{
		{"", []string{"tx-2", "tx-1"}, 2},
		{"provider=fawry", []string{"tx-1"}, 1},
		{"status=paid", []string{"tx-1"}, 1},
		{"q=tx-2", []string{"tx-2"}, 1},
		{"q=FAWRY", []string{"tx-1"}, 1},
		{"provider=stripe&status=paid", []string{}, 0},
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/_admin/transactions?"+tc.query, nil)
			rec := httptest.NewRecorder()
			listTransactionsHandler(store)(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status: want 200, got %d", rec.Code)
			}
			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if body["total"] != tc.wantTotal {
				t.Errorf("total: want %v, got %v", tc.wantTotal, body["total"])
			}
			results, _ := body["results"].([]any)
			if len(results) != len(tc.wantRefs) {
				t.Fatalf("results: want %v, got %v", tc.wantRefs, results)
			}
		})
	}
}

func TestGetTransactionHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "fawry", Status: engine.TransactionStatusPaid, Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions/tx-1", nil)
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()

	getTransactionHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["transaction"] == nil {
		t.Error("expected transaction in response")
	}
}

func TestGetTransactionHandlerNotFound(t *testing.T) {
	store := engine.NewMemoryStore()
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions/missing", nil)
	req.SetPathValue("ref", "missing")
	rec := httptest.NewRecorder()

	getTransactionHandler(store)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestReplayTransactionWebhookHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	dispatcher := webhook.NewDispatcherFromBuilder("http://localhost/hook", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)

	req := newAdminRequest(http.MethodPost, "/_admin/transactions/tx-1/replay-webhook")
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()

	replayTransactionWebhookHandler(store, dispatcher, nil)(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status: want 202, got %d", rec.Code)
	}
}

func TestReplayTransactionWebhookHandlerProviderDispatcher(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "stripe", Status: engine.TransactionStatusPaid, Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	active := webhook.NewDispatcherFromBuilder("http://localhost/hook", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)
	stripeDisp := webhook.NewDispatcherFromBuilder("http://localhost/stripe-hook", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)

	req := newAdminRequest(http.MethodPost, "/_admin/transactions/tx-1/replay-webhook")
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()

	replayTransactionWebhookHandler(store, active, map[string]*webhook.Dispatcher{"stripe": stripeDisp})(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status: want 202, got %d", rec.Code)
	}

	attempt, err := stripeDisp.Store.Get("tx-1")
	if err != nil {
		t.Fatalf("get attempt: %v", err)
	}
	if attempt == nil {
		t.Error("expected webhook attempt in stripe dispatcher store")
	}
}

func TestReplayTransactionWebhookHandlerMissing(t *testing.T) {
	store := engine.NewMemoryStore()
	req := newAdminRequest(http.MethodPost, "/_admin/transactions/missing/replay-webhook")
	req.SetPathValue("ref", "missing")
	rec := httptest.NewRecorder()

	replayTransactionWebhookHandler(store, nil, nil)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestLedgerHandler(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "fawry", Status: engine.TransactionStatusPaid, Amount: 10.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	dispatcher := webhook.NewDispatcherFromBuilder("http://localhost/hook", 0,
		func(_ context.Context, _ provider.Transaction) ([]byte, error) { return []byte(`{}`), nil },
		nil,
	)
	if _, err := dispatcher.Dispatch(context.Background(), "tx-1", webhook.PaymentStatusPaid); err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, dispatcher)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, ok := body["results"].([]any)
	if !ok || len(results) != 2 {
		t.Fatalf("want 2 ledger events, got %v", body)
	}
	if body["total"] != 2.0 {
		t.Errorf("total: want 2, got %v", body["total"])
	}
}

func TestLedgerHandlerFilters(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "fawry", Status: engine.TransactionStatusPaid, Amount: 10.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-2", Provider: "stripe", Status: engine.TransactionStatusUnpaid, Amount: 20.0, Currency: "USD"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	cases := []struct {
		query     string
		wantTotal int
	}{
		{"", 2},
		{"type=transaction", 2},
		{"type=webhook", 0},
		{"provider=fawry", 1},
		{"status=paid", 1},
		{"q=stripe", 1},
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/_admin/ledger?"+tc.query, nil)
			rec := httptest.NewRecorder()
			ledgerHandler(store, nil)(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status: want 200, got %d", rec.Code)
			}
			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if int(body["total"].(float64)) != tc.wantTotal {
				t.Errorf("total: want %d, got %v", tc.wantTotal, body["total"])
			}
		})
	}
}

func TestLedgerHandlerLimitDefault(t *testing.T) {
	store := engine.NewMemoryStore()
	for i := 0; i < 60; i++ {
		if _, _, err := store.CreateOrGet(engine.Transaction{Reference: fmt.Sprintf("tx-%d", i), Provider: "default", Status: engine.TransactionStatusPaid, Amount: 1.0}); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 50 {
		t.Errorf("default limit should be 50, got %d", len(results))
	}
}

func TestLedgerHandlerIncludesTraceID(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 10.0, TraceID: "trace-abc"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, ok := body["results"].([]any)
	if !ok || len(results) != 1 {
		t.Fatalf("want 1 event, got %v", body)
	}
	ev, ok := results[0].(map[string]any)
	if !ok {
		t.Fatalf("event not map: %v", results[0])
	}
	if ev["trace_id"] != "trace-abc" {
		t.Errorf("trace_id: want trace-abc, got %v", ev["trace_id"])
	}
}

func TestLedgerHandlerFiltersByTraceID(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 10.0, TraceID: "trace-abc"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-2", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 20.0, TraceID: "trace-xyz"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger?q=trace-abc", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["total"] != 1.0 {
		t.Errorf("total: want 1, got %v", body["total"])
	}
}

func TestGetTransactionHandlerIncludesTraceID(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 10.0, TraceID: "trace-abc"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions/tx-1", nil)
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()

	getTransactionHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tx, ok := body["transaction"].(map[string]any)
	if !ok {
		t.Fatalf("transaction not map: %v", body["transaction"])
	}
	if tx["trace_id"] != "trace-abc" {
		t.Errorf("trace_id: want trace-abc, got %v", tx["trace_id"])
	}
}

func TestListTransactionsHandlerStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions", nil)
	rec := httptest.NewRecorder()

	listTransactionsHandler(&errTransactionStore{listErr: errors.New("boom")})(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestOnboardingHandlerWebhooksDisabled(t *testing.T) {
	store := engine.NewMemoryStore()
	req := httptest.NewRequest(http.MethodGet, "/_admin/onboarding", nil)
	rec := httptest.NewRecorder()
	onboardingHandler("default", []string{"default"}, store, nil)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["webhooks_enabled"] != false {
		t.Errorf("webhooks_enabled want false, got %v", body["webhooks_enabled"])
	}
}

func TestProviderHealthHandler(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	configYAML := []byte(`providers:
  default:
    enabled: true
    config: {}
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
  stripe:
    enabled: false
    config:
      publishable_key: pk_test
`)
	if err := os.WriteFile(path, configYAML, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ConfigPath: path})

	cases := []struct {
		name       string
		wantStatus string
	}{
		{"default", "healthy"},
		{"fawry", "healthy"},
		{"stripe", "disabled"},
		{"unknown", "disabled"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/_admin/providers/"+tc.name+"/health", nil)
			req.SetPathValue("name", tc.name)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status: want 200, got %d", rec.Code)
			}
			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if body["name"] != tc.name {
				t.Errorf("name: want %s, got %v", tc.name, body["name"])
			}
			if body["status"] != tc.wantStatus {
				t.Errorf("status: want %s, got %v", tc.wantStatus, body["status"])
			}
		})
	}
}

func TestProviderHealthHandlerMisconfigured(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	configYAML := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: ""
`)
	if err := os.WriteFile(path, configYAML, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ConfigPath: path})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/fawry/health", nil)
	req.SetPathValue("name", "fawry")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "misconfigured" {
		t.Errorf("status: want misconfigured, got %v", body["status"])
	}
	if body["reason"] == "" {
		t.Error("expected non-empty reason for misconfigured provider")
	}
}

func TestListProvidersHandlerEnriched(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/providers", nil)
	rec := httptest.NewRecorder()

	listProvidersHandler(RouterConfig{ActiveProvider: "fawry", EnabledProviders: []string{"fawry"}, Host: "127.0.0.1", Port: 9000})(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	details, ok := body["providers"].(map[string]any)
	if !ok {
		t.Fatalf("providers: want details map, got %v", body["providers"])
	}
	info, ok := details["fawry"].(map[string]any)
	if !ok {
		t.Fatalf("fawry provider details missing: got %v", details)
	}

	if info["base_url"] != "http://127.0.0.1:9000/fawry/v1" {
		t.Errorf("base_url: want http://127.0.0.1:9000/fawry/v1, got %v", info["base_url"])
	}

	envVars, ok := info["env_vars"].([]any)
	if !ok || len(envVars) == 0 {
		t.Fatalf("env_vars: want non-empty list, got %v", info["env_vars"])
	}
	if envVars[0] != "MUARA_FAWRY_MERCHANT_CODE" {
		t.Errorf("env_vars[0]: want MUARA_FAWRY_MERCHANT_CODE, got %v", envVars[0])
	}

	versionDetails, ok := info["version_details"].(map[string]any)
	if !ok {
		t.Fatalf("version_details: want map, got %v", info["version_details"])
	}
	v1, ok := versionDetails["v1"].(map[string]any)
	if !ok {
		t.Fatalf("version_details.v1: want map, got %v", versionDetails["v1"])
	}
	if v1["base_url"] != "http://127.0.0.1:9000/fawry/v1" {
		t.Errorf("version_details.v1.base_url: want http://127.0.0.1:9000/fawry/v1, got %v", v1["base_url"])
	}
	if v1["sample_route"] != "/fawry/v1/charge" {
		t.Errorf("version_details.v1.sample_route: want /fawry/v1/charge, got %v", v1["sample_route"])
	}
}

func TestGetProviderHandler(t *testing.T) {
	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ActiveProvider: "fawry", EnabledProviders: []string{"fawry"}, Host: "127.0.0.1", Port: 9000})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/fawry", nil)
	req.SetPathValue("name", "fawry")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if body["name"] != "fawry" {
		t.Errorf("name: want fawry, got %v", body["name"])
	}
	if body["display_name"] != "Fawry" {
		t.Errorf("display_name: want Fawry, got %v", body["display_name"])
	}
	if body["enabled"] != true {
		t.Errorf("enabled: want true, got %v", body["enabled"])
	}
	if body["active"] != true {
		t.Errorf("active: want true, got %v", body["active"])
	}
	if body["base_url"] != "http://127.0.0.1:9000/fawry/v1" {
		t.Errorf("base_url: want http://127.0.0.1:9000/fawry/v1, got %v", body["base_url"])
	}
	if body["webhook_target_url"] != "" {
		t.Errorf("webhook_target_url: want empty, got %v", body["webhook_target_url"])
	}

	envVars, ok := body["env_vars"].([]any)
	if !ok || len(envVars) != 3 {
		t.Fatalf("env_vars: want 3 entries, got %v", body["env_vars"])
	}
	wantEnvVars := []string{"MUARA_FAWRY_MERCHANT_CODE", "MUARA_FAWRY_MERCHANT_SECURITY_KEY", "MUARA_FAWRY_WEBHOOK_SECRET"}
	for i, want := range wantEnvVars {
		if envVars[i] != want {
			t.Errorf("env_vars[%d]: want %s, got %v", i, want, envVars[i])
		}
	}

	versionDetails, ok := body["version_details"].(map[string]any)
	if !ok {
		t.Fatalf("version_details: want map, got %v", body["version_details"])
	}
	if len(versionDetails) != 2 {
		t.Errorf("version_details: want 2 versions, got %d", len(versionDetails))
	}
}

func TestGetProviderHandlerDisabled(t *testing.T) {
	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ActiveProvider: "default", EnabledProviders: []string{"default"}, Host: "127.0.0.1", Port: 9000})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/stripe", nil)
	req.SetPathValue("name", "stripe")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if body["name"] != "stripe" {
		t.Errorf("name: want stripe, got %v", body["name"])
	}
	if body["enabled"] != false {
		t.Errorf("enabled: want false, got %v", body["enabled"])
	}
	if body["active"] != false {
		t.Errorf("active: want false, got %v", body["active"])
	}
	if body["base_url"] != "http://127.0.0.1:9000/v1" {
		t.Errorf("base_url: want http://127.0.0.1:9000/v1, got %v", body["base_url"])
	}
}

func TestGetProviderHandlerNotFound(t *testing.T) {
	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ActiveProvider: "default", EnabledProviders: []string{"default"}})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/unknown", nil)
	req.SetPathValue("name", "unknown")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestGetProviderHandlerWebhookTargetURL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	configYAML := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
webhook:
  targets:
    fawry: http://fawry.example.com/webhook
`)
	if err := os.WriteFile(path, configYAML, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ActiveProvider: "fawry", EnabledProviders: []string{"fawry"}, ConfigPath: path, Host: "127.0.0.1", Port: 9000})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/fawry", nil)
	req.SetPathValue("name", "fawry")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if body["webhook_target_url"] != "http://fawry.example.com/webhook" {
		t.Errorf("webhook_target_url: want http://fawry.example.com/webhook, got %v", body["webhook_target_url"])
	}
}

func TestGetProviderHandlerWebhookTargetURLFromConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	configYAML := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
      webhook_url: http://legacy.example.com/webhook
`)
	if err := os.WriteFile(path, configYAML, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ActiveProvider: "fawry", EnabledProviders: []string{"fawry"}, ConfigPath: path, Host: "127.0.0.1", Port: 9000})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/fawry", nil)
	req.SetPathValue("name", "fawry")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if body["webhook_target_url"] != "http://legacy.example.com/webhook" {
		t.Errorf("webhook_target_url: want http://legacy.example.com/webhook, got %v", body["webhook_target_url"])
	}
}
