package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/fawry"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func newTestFawryProvider(t *testing.T) provider.Provider {
	t.Helper()
	p := fawry.NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"merchant_code":         "muara-merchant-code",
		"merchant_security_key": "muara-fawry-secret",
		"webhook_secret":        "muara-webhook-secret",
	}); err != nil {
		t.Fatalf("init fawry: %v", err)
	}
	return p
}

func TestRouterRegistersFawryChargeRoute(t *testing.T) {
	// Given Fawry is enabled and initialized
	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": newTestFawryProvider(t)},
	})

	req := fawry.ChargeRequest{
		MerchantCode:      "muara-merchant-code",
		MerchantRefNum:    "router-ref-1",
		CustomerProfileID: "user-1",
		ReturnURL:         "http://localhost/callback",
		ChargeItems: []fawry.ChargeItem{
			{ItemID: "prod-1", Price: 10.00, Quantity: 1},
		},
	}
	req.Signature = fawry.Sign(req, "muara-fawry-secret")
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When the route is handled
	router.ServeHTTP(rec, httpReq)

	// Then it returns 200 for a valid request.
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestRouterRegistersDefaultChargeRoute(t *testing.T) {
	// Given default is enabled
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	httpReq := httptest.NewRequest(http.MethodPost, "/default/charge", nil)
	rec := httptest.NewRecorder()

	// When the route is handled
	router.ServeHTTP(rec, httpReq)

	// Then it returns 200.
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["provider"] != "default" {
		t.Errorf("provider: want default, got %v", body["provider"])
	}
}

func TestRouterDoesNotRegisterDisabledProvider(t *testing.T) {
	// Given default is enabled but fawry is disabled
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	httpReq := httptest.NewRequest(http.MethodPost, "/fawry/charge", nil)
	rec := httptest.NewRecorder()

	// When a Fawry request is made
	router.ServeHTTP(rec, httpReq)

	// Then it returns 404 because the route is not mounted.
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestRouterRegistersActiveProviderEscapePage(t *testing.T) {
	// Given the active provider exposes an escape handler
	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": newTestFawryProvider(t)},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	rec := httptest.NewRecorder()

	// When the escape page is requested
	router.ServeHTTP(rec, httpReq)

	// Then it returns 200.
	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestRouterDoesNotRegisterEscapePageForProviderWithoutEscape(t *testing.T) {
	// Given the active provider has no escape handler
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/default-escape", nil)
	rec := httptest.NewRecorder()

	// When the escape page is requested
	router.ServeHTTP(rec, httpReq)

	// Then it returns 404.
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestRouterRegistersReadyzEndpoint(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default", "fawry"},
		Providers:        map[string]provider.Provider{"fawry": newTestFawryProvider(t)},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status: want ok, got %v", body["status"])
	}
	enabled, ok := body["providers"].([]any)
	if !ok || len(enabled) != 2 {
		t.Errorf("providers: want 2 entries, got %v", body["providers"])
	}
}

func TestRouterPprofDisabledByDefault(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", rec.Code)
	}
}

func TestRouterPprofEnabled(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		Pprof:            true,
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("Profiles")) && !bytes.Contains(rec.Body.Bytes(), []byte("profile")) {
		t.Errorf("pprof index missing expected content: %s", rec.Body.String())
	}
}

func TestRouterPprofCmdlineEnabled(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		Pprof:            true,
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/_admin/debug/pprof/cmdline", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httpReq)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
}

func TestProviderRouterDoesNotServeAdminAPI(t *testing.T) {
	store := engine.NewMemoryStore()
	providerRouter := NewProviderRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		TransactionStore: store,
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions", nil)
	rec := httptest.NewRecorder()
	providerRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("provider router should not serve /_admin/transactions: want 404, got %d", rec.Code)
	}
}

func TestAdminRouterServesTransactions(t *testing.T) {
	store := engine.NewMemoryStore()
	adminRouter := NewAdminRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		TransactionStore: store,
	})

	req := httptest.NewRequest(http.MethodGet, "/_admin/transactions", nil)
	rec := httptest.NewRecorder()
	adminRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("admin router should serve /_admin/transactions: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestFawryEscapeActionPreservesProviderFlow(t *testing.T) {
	// The Fawry escape page/action is part of the customer payment flow, not an
	// admin-only endpoint, so it must remain reachable without admin elevation.
	p := newTestFawryProvider(t)
	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "flow-ref-1", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}
	if sp, ok := p.(interface{ SetStore(engine.TransactionStore) }); ok {
		sp.SetStore(store)
	}

	router := NewRouter(RouterConfig{
		ActiveProvider:   "fawry",
		EnabledProviders: []string{"fawry"},
		Providers:        map[string]provider.Provider{"fawry": p},
		TransactionStore: store,
	})

	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=flow-ref-1&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestIsMutatingAdminRoute(t *testing.T) {
	cases := []struct {
		method string
		path   string
		want   bool
	}{
		// Stripe dashboard simulation routes are part of the customer payment
		// flow, so they must not be wrapped with requireAdmin.
		{http.MethodPost, "/_admin/stripe/success", false},
		{http.MethodPost, "/_admin/stripe/fail", false},
		{http.MethodPost, "/_admin/stripe/cancel", false},
		{http.MethodPost, "/_admin/stripe/payment_intent/pi_123", false},
		{http.MethodGet, "/_admin/stripe/payment_intent/pi_123", false},

		// Provider pay/escape pages are part of the customer flow and must not
		// be wrapped with requireAdmin.
		{http.MethodPost, "/_admin/fawry-escape", false},
		{http.MethodPost, "/_admin/billplz/pay/123", false},
		{http.MethodPost, "/_admin/toyyibpay/pay/abc", false},
		{http.MethodPost, "/_admin/ipay88/pay/ref", false},

		// Non-admin provider API routes are never wrapped.
		{http.MethodPost, "/v1/checkout/sessions", false},
		{http.MethodPost, "/fawry/charge", false},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			if got := isMutatingAdminRoute(tc.method, tc.path); got != tc.want {
				t.Errorf("isMutatingAdminRoute(%q, %q) = %v, want %v", tc.method, tc.path, got, tc.want)
			}
		})
	}
}

func TestRequireAdminWrapperAllowsAdmin(t *testing.T) {
	handler := requireAdmin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := newAdminRequest(http.MethodPost, "/_admin/stripe/success")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
}

func TestRequireAdminWrapperDeniesViewer(t *testing.T) {
	handler := requireAdmin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := newViewerRequest(http.MethodPost, "/_admin/stripe/success")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status: want 403, got %d", rec.Code)
	}
}
