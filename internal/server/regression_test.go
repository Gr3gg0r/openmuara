package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider/defaultplugin"
)

var _ = defaultplugin.NewProvider

type brokenTransactionStore struct {
	engine.TransactionStore
	listErr error
	getErr  error
}

func (b *brokenTransactionStore) List(limit, offset int) ([]engine.Transaction, error) {
	if b.listErr != nil {
		return nil, b.listErr
	}
	return b.TransactionStore.List(limit, offset)
}

func (b *brokenTransactionStore) GetByReference(ref string) (engine.Transaction, bool, error) {
	if b.getErr != nil {
		return engine.Transaction{}, false, b.getErr
	}
	return b.TransactionStore.GetByReference(ref)
}

type zeroUpdatedAtStore struct{ engine.TransactionStore }

func (z *zeroUpdatedAtStore) List(limit, offset int) ([]engine.Transaction, error) {
	txs, err := z.TransactionStore.List(limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range txs {
		txs[i].UpdatedAt = time.Time{}
	}
	return txs, nil
}

func TestPageParamsEnforcesMaxLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?limit=9999", nil)
	limit, _ := pageParams(req)
	if limit != maxPageLimit {
		t.Errorf("limit: want %d, got %d", maxPageLimit, limit)
	}
}

func TestListTransactionsHandlerLimitZero(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions?limit=0", nil)
	rec := httptest.NewRecorder()
	listTransactionsHandler(store)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 1 {
		t.Errorf("want 1 result with limit=0, got %d", len(results))
	}
}

func TestListTransactionsHandlerNegativeOffset(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions?offset=-5", nil)
	rec := httptest.NewRecorder()
	listTransactionsHandler(store)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["offset"] != 0.0 {
		t.Errorf("offset: want 0, got %v", body["offset"])
	}
}

func TestListTransactionsHandlerOffsetBeyondTotal(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions?offset=10", nil)
	rec := httptest.NewRecorder()
	listTransactionsHandler(store)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 0 {
		t.Errorf("want 0 results, got %d", len(results))
	}
	if body["offset"] != 1.0 {
		t.Errorf("offset should clamp to total, got %v", body["offset"])
	}
}

func TestGetTransactionHandlerStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions/tx-1", nil)
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()
	getTransactionHandler(&brokenTransactionStore{getErr: errors.New("boom")})(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestReplayTransactionWebhookHandlerNoDispatcher(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Provider: "default", Status: engine.TransactionStatusPaid, Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	req := newAdminRequest(http.MethodPost, "/_admin/transactions/tx-1/replay-webhook")
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()
	replayTransactionWebhookHandler(store, nil, nil)(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: want 503, got %d", rec.Code)
	}
}

func TestLedgerHandlerLimitZero(t *testing.T) {
	store := engine.NewMemoryStore()
	for i := 0; i < 3; i++ {
		ref := fmt.Sprintf("tx-%d", i)
		if _, _, err := store.CreateOrGet(engine.Transaction{Reference: ref, Amount: 1.0}); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger?limit=0", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 3 {
		t.Errorf("want all 3 events with limit=0, got %d", len(results))
	}
}

func TestLedgerHandlerOffsetBeyondTotal(t *testing.T) {
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger?offset=10", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 0 {
		t.Errorf("want 0 results, got %d", len(results))
	}
}

func TestLedgerHandlerStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(&brokenTransactionStore{listErr: errors.New("boom")}, nil)(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestLedgerHandlerFallsBackToCreatedAt(t *testing.T) {
	base := engine.NewMemoryStore()
	created := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, _, err := base.CreateOrGet(engine.Transaction{Reference: "tx-1", Amount: 10.0, CreatedAt: created}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	store := &zeroUpdatedAtStore{TransactionStore: base}

	req := httptest.NewRequest(http.MethodGet, "/_admin/ledger", nil)
	rec := httptest.NewRecorder()
	ledgerHandler(store, nil)(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results, _ := body["results"].([]any)
	if len(results) != 1 {
		t.Fatalf("want 1 event, got %d", len(results))
	}
	ev := results[0].(map[string]any)
	if ev["time"] == "" {
		t.Error("expected time fallback to created_at")
	}
}

func TestBuildLedgerEventsCapsMaxEvents(t *testing.T) {
	txs := make([]engine.Transaction, 1005)
	for i := range txs {
		txs[i] = engine.Transaction{Reference: fmt.Sprintf("tx-%d", i), Amount: 1.0}
	}
	events := buildLedgerEvents(txs, nil)
	if len(events) != 1000 {
		t.Errorf("want max 1000 events, got %d", len(events))
	}
}

func TestProviderBaseURLDefaults(t *testing.T) {
	if got := providerBaseURL("", 0, "", "fawry", "v1"); got != "http://127.0.0.1:9000/fawry/v1" {
		t.Errorf("fawry default base url: got %q", got)
	}
	if got := providerBaseURL("", 0, "", "stripe", "v1"); got != "http://127.0.0.1:9000/v1" {
		t.Errorf("stripe default base url: got %q", got)
	}
	if got := providerBaseURL("", 0, "", "default", "v1"); got != "http://127.0.0.1:9000" {
		t.Errorf("default default base url: got %q", got)
	}
}

func TestProviderBaseURLPublicBaseURL(t *testing.T) {
	if got := providerBaseURL("", 0, "https://muara.example.com", "stripe", "v1"); got != "https://muara.example.com/v1" {
		t.Errorf("stripe public base url: got %q", got)
	}
}

func TestProviderWebhookTargetURLEmptyConfigPath(t *testing.T) {
	if got := providerWebhookTargetURL("", "fawry"); got != "" {
		t.Errorf("want empty, got %q", got)
	}
}

func TestProviderWebhookTargetURLLoadError(t *testing.T) {
	path := t.TempDir() // a directory causes ReadInConfig to fail
	if got := providerWebhookTargetURL(path, "fawry"); got != "" {
		t.Errorf("want empty on load error, got %q", got)
	}
}

func TestProviderCategoryRedirectGateways(t *testing.T) {
	for _, name := range []string{"billplz", "toyyibpay", "senangpay", "ipay88"} {
		if got := providerCategory(name); got != "redirect" {
			t.Errorf("providerCategory(%q): want redirect, got %q", name, got)
		}
	}
}

func TestRealProvidersForRegisteredGateways(t *testing.T) {
	cases := map[string]string{
		"stripe":    "Stripe",
		"fawry":     "Fawry",
		"billplz":   "Billplz",
		"toyyibpay": "ToyyibPay",
		"senangpay": "SenangPay",
		"ipay88":    "iPay88",
		"unknown":   "OpenMuara Default",
	}
	for name, want := range cases {
		got := realProvidersFor(name)
		if len(got) == 0 {
			t.Errorf("realProvidersFor(%q): want non-empty, got %v", name, got)
			continue
		}
		if got[0] != want {
			t.Errorf("realProvidersFor(%q): want first %q, got %v", name, want, got)
		}
	}
}

func TestNextStepForProviderUnknown(t *testing.T) {
	step := nextStepForProvider("custom")
	if !strings.Contains(step["hint"].(string), "custom") {
		t.Errorf("expected hint for custom provider, got %v", step)
	}
}

func TestChoiceForProviderUnknown(t *testing.T) {
	if choice := choiceForProvider("not-a-provider"); choice.Key != "" {
		t.Errorf("expected empty choice, got %+v", choice)
	}
}

func TestBuildProviderInfoUnregistered(t *testing.T) {
	info := buildProviderInfo("not-a-provider", "not-a-provider", "127.0.0.1", 9000, "", false, nil)
	if info["name"] != "not-a-provider" {
		t.Errorf("name mismatch: %v", info["name"])
	}
}

func TestProviderHealthHandlerConfigLoadError(t *testing.T) {
	path := t.TempDir() // directory causes load error
	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{ConfigPath: path})

	req := httptest.NewRequest(http.MethodGet, "/_admin/providers/fawry/health", nil)
	req.SetPathValue("name", "fawry")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestHashPasswordEmptyRegression(t *testing.T) {
	_, err := HashPassword("")
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestCORSMiddlewareWildcardWithCredentials(t *testing.T) {
	cfg := config.CORSConfig{AllowedOrigins: []string{"*"}, AllowCredentials: true}
	handler := CORSMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://example.com" {
		t.Errorf("origin: want echoed origin, got %q", origin)
	}
	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("expected credentials header")
	}
}

func TestIsAdminPathNonAdmin(t *testing.T) {
	if isAdminPath("/healthz") {
		t.Error("expected /healthz not to be admin path")
	}
}

func TestIsCheckoutPayPath(t *testing.T) {
	if !isCheckoutPayPath("/v1/checkout/sessions/cs_test/pay") {
		t.Error("expected checkout pay path")
	}
	if isCheckoutPayPath("/v1/checkout/sessions/cs_test") {
		t.Error("expected non-pay path to be false")
	}
}

func TestCSRFRequestTokenReturnsEmptyForGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := CSRFRequestToken(req); got != "" {
		t.Errorf("want empty, got %q", got)
	}
}

func TestCSRFRejectsMismatchedToken(t *testing.T) {
	cfg := CSRFGuardConfig{Enabled: true}
	handler := CSRFGuardMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodPost, "/_admin/config/reload", nil)
	req.Header.Set(csrfHeaderName, "mismatch")
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "cookievalue", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d", rec.Code)
	}
}

func TestCSRFRejectsMissingRequestToken(t *testing.T) {
	cfg := CSRFGuardConfig{Enabled: true}
	handler := CSRFGuardMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	req := httptest.NewRequest(http.MethodPost, "/_admin/config/reload", nil)
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "cookievalue", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d", rec.Code)
	}
}

func TestServerAddrWithPortZero(t *testing.T) {
	srv := New(Config{Host: "127.0.0.1", Port: 0, Handler: http.NotFoundHandler()})
	if got := srv.Addr(); got != "" {
		t.Errorf("want empty addr before listen, got %q", got)
	}
}

func TestServerBaseURLWithTLS(t *testing.T) {
	srv := New(Config{Host: "127.0.0.1", Port: 8443, Handler: http.NotFoundHandler(), TLSCert: "cert.pem", TLSKey: "key.pem"})
	if got := srv.BaseURL(); got != "https://127.0.0.1:8443" {
		t.Errorf("want https base url, got %q", got)
	}
}

func TestScenarioHandlerEmptyOutcome(t *testing.T) {
	store := engine.NewMemoryStore()
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/")
	rec := httptest.NewRecorder()
	newScenarioHandler(store)(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestScenarioHandlerUnknownOutcomeRegression(t *testing.T) {
	store := engine.NewMemoryStore()
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/unknown?ref=r")
	rec := httptest.NewRecorder()
	newScenarioHandler(store)(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestScenarioHandlerStoreError(t *testing.T) {
	req := newAdminRequest(http.MethodPost, "/_admin/scenario/success?ref=r")
	rec := httptest.NewRecorder()
	newScenarioHandler(&brokenTransactionStore{getErr: errors.New("boom")})(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestLogSecurityEventNilLogger(_ *testing.T) {
	// Should not panic when no logger is in the context.
	logSecurityEvent(context.Background(), "test", "detail")
}

func TestNewRouterSkipsUnknownProvider(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default", "unknown-provider"},
	})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
}

func TestNewAdminRouterPprof(t *testing.T) {
	router := NewAdminRouter(RouterConfig{Pprof: true})
	req := httptest.NewRequest(http.MethodGet, "/_admin/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
}
