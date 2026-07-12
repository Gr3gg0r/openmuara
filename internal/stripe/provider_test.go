package stripe

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestProviderName(t *testing.T) {
	p := NewProvider()
	if got := p.Name(); got != "stripe" {
		t.Errorf("name: want stripe, got %q", got)
	}
}

func TestProviderInitValidConfig(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	})
	if err != nil {
		t.Fatalf("init: expected no error, got %v", err)
	}
}

func TestProviderInitMissingPublishableKey(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	err := p.Init(map[string]any{
		"secret_key":     "sk_test_muara",
		"webhook_secret": "whsec_muara",
	})
	if err == nil {
		t.Fatal("expected error for missing publishable_key")
	}
}

func TestProviderInitMissingSecretKey(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"webhook_secret":  "whsec_muara",
	})
	if err == nil {
		t.Fatal("expected error for missing secret_key")
	}
}

func TestProviderInitMissingWebhookSecret(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
	})
	if err == nil {
		t.Fatal("expected error for missing webhook_secret")
	}
}

func TestProviderRoutes(t *testing.T) {
	p := NewProvider()
	routes := p.Routes()
	if len(routes) != 13 {
		t.Fatalf("routes: want 13, got %d", len(routes))
	}

	paths := make(map[string][]string)
	for _, r := range routes {
		paths[r.Path] = append(paths[r.Path], r.Method)
	}
	if paths["/v1/checkout/sessions"][0] != http.MethodPost {
		t.Errorf("create route missing or wrong method")
	}
	payMethods := paths["/v1/checkout/sessions/{id}/pay"]
	if len(payMethods) != 2 {
		t.Fatalf("checkout pay routes: want 2 methods, got %d", len(payMethods))
	}
	payMethodSet := make(map[string]bool)
	for _, m := range payMethods {
		payMethodSet[m] = true
	}
	if !payMethodSet[http.MethodGet] {
		t.Errorf("checkout pay GET route missing")
	}
	if !payMethodSet[http.MethodPost] {
		t.Errorf("checkout pay POST route missing")
	}
	if paths["/v1/checkout/sessions/"][0] != http.MethodGet {
		t.Errorf("checkout retrieve route missing or wrong method")
	}
	if paths["/v1/payment_intents"][0] != http.MethodPost {
		t.Errorf("payment intent create route missing or wrong method")
	}
	if paths["/v1/payment_intents/"][0] != http.MethodGet {
		t.Errorf("payment intent retrieve route missing or wrong method")
	}
	if paths["/v1/payment_intents/{id}/confirm"][0] != http.MethodPost {
		t.Errorf("payment intent confirm route missing or wrong method")
	}
	if paths["/v1/payment_intents/{id}/cancel"][0] != http.MethodPost {
		t.Errorf("payment intent cancel route missing or wrong method")
	}
	adminMethods := paths["/_admin/stripe/payment_intent/{id}"]
	if len(adminMethods) != 2 {
		t.Fatalf("payment intent admin routes: want 2 methods, got %d", len(adminMethods))
	}
	adminMethodSet := make(map[string]bool)
	for _, m := range adminMethods {
		adminMethodSet[m] = true
	}
	if !adminMethodSet[http.MethodGet] {
		t.Errorf("payment intent admin GET route missing")
	}
	if !adminMethodSet[http.MethodPost] {
		t.Errorf("payment intent admin POST route missing")
	}
	if paths["/_admin/stripe/success"][0] != http.MethodPost {
		t.Errorf("success simulation route missing or wrong method")
	}
	if paths["/_admin/stripe/fail"][0] != http.MethodPost {
		t.Errorf("fail simulation route missing or wrong method")
	}
	if paths["/_admin/stripe/cancel"][0] != http.MethodPost {
		t.Errorf("cancel simulation route missing or wrong method")
	}
}

func TestProviderPayloadBuilder(t *testing.T) {
	p := NewProvider()
	p.SetBaseURL("http://localhost")
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	req := validCreateRequest()
	body, _ := json.Marshal(req)
	p.ChargeHandler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))

	var session CheckoutSession
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(body)))
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("decode session: %v", err)
	}

	builder := p.PayloadBuilder()
	payload, err := builder(context.Background(), provider.Transaction{
		Reference: session.ID,
		Status:    "complete",
	})
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}

	var got Event
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if got.Object != "event" {
		t.Errorf("object: want event, got %q", got.Object)
	}
	if got.Type != "checkout.session.completed" {
		t.Errorf("type: want checkout.session.completed, got %q", got.Type)
	}
	if got.Data.Object == nil {
		t.Fatal("event data object is nil")
	}
	if got.Data.Object.ID != session.ID {
		t.Errorf("id mismatch: want %q, got %q", session.ID, got.Data.Object.ID)
	}
	if got.Data.Object.Status != "complete" {
		t.Errorf("status: want complete, got %q", got.Data.Object.Status)
	}
	if got.Data.Object.PaymentStatus != "paid" {
		t.Errorf("payment_status: want paid, got %q", got.Data.Object.PaymentStatus)
	}
}

func TestProviderEscapeHandlerIsNil(t *testing.T) {
	p := NewProvider()
	if h := p.EscapeHandler(); h != nil {
		t.Fatal("expected nil escape handler")
	}
}

func TestProviderSetStoreAndWebhookHandler(t *testing.T) {
	p := NewProvider()
	store := engine.NewMemoryStore()
	p.SetStore(store)

	rec := httptest.NewRecorder()
	p.WebhookHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/webhook", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("webhook handler status: want 200, got %d", rec.Code)
	}
}

func TestProviderPayloadBuilderMissingSession(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	_, err := p.PayloadBuilder()(context.Background(), provider.Transaction{Reference: "cs_test_missing"})
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestProviderPayloadHeadersMissingSession(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	_, err := p.PayloadHeaders(context.Background(), provider.Transaction{Reference: "cs_test_missing"})
	if err == nil {
		t.Fatal("expected error for missing session")
	}
}

func TestProviderVerifyWebhookSignature(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	if _, _, err := p.ledger.CreateOrGet(engine.Transaction{Reference: "cs_test_verify", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	tx := provider.Transaction{Reference: "cs_test_verify", Status: "complete"}
	payload, err := p.PayloadBuilder()(context.Background(), tx)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	headers, err := p.PayloadHeaders(context.Background(), tx)
	if err != nil {
		t.Fatalf("build headers: %v", err)
	}

	valid, err := p.VerifyWebhookSignature(payload, headers)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !valid {
		t.Error("expected signature to be valid")
	}

	missing, err := p.VerifyWebhookSignature(payload, map[string]string{})
	if err != nil {
		t.Fatalf("verify missing header: %v", err)
	}
	if missing {
		t.Error("expected missing signature header to be invalid")
	}
}

func TestMapPaymentStatus(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"complete", "paid"},
		{"PAID", "paid"},
		{"open", "unpaid"},
		{"UNPAID", "unpaid"},
		{"canceled", "unpaid"},
		{"unknown", "unpaid"},
	}
	for _, tc := range cases {
		if got := mapPaymentStatus(tc.in); got != tc.want {
			t.Errorf("%q: want %q, got %q", tc.in, tc.want, got)
		}
	}
}

func TestMapSessionStatus(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"complete", "complete"},
		{"PAID", "complete"},
		{"UNPAID", "expired"},
		{"canceled", "expired"},
		{"expired", "expired"},
		{"open", "open"},
		{"unknown", "open"},
	}
	for _, tc := range cases {
		if got := mapSessionStatus(tc.in); got != tc.want {
			t.Errorf("%q: want %q, got %q", tc.in, tc.want, got)
		}
	}
}

func TestIsCanceledStatus(t *testing.T) {
	if !isCanceledStatus("UNPAID") {
		t.Error("UNPAID should be canceled")
	}
	if !isCanceledStatus("canceled") {
		t.Error("canceled should be canceled")
	}
	if isCanceledStatus("PAID") {
		t.Error("PAID should not be canceled")
	}
}

func TestToStringSlice(t *testing.T) {
	got := toStringSlice([]any{"a", "b", 1, "c"})
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestSessionFromTransaction(t *testing.T) {
	tx := engine.Transaction{
		Reference:   "cs_test_fallback",
		Amount:      12.34,
		Currency:    "MYR",
		Status:      engine.TransactionStatusPaid,
		CustomerRef: "user@example.com",
	}
	session := sessionFromTransaction(tx)
	if session.ID != tx.Reference {
		t.Errorf("id: want %q, got %q", tx.Reference, session.ID)
	}
	if session.Status != "complete" {
		t.Errorf("status: want complete, got %q", session.Status)
	}
	if session.PaymentStatus != "paid" {
		t.Errorf("payment_status: want paid, got %q", session.PaymentStatus)
	}
	if session.AmountTotal != 1234 {
		t.Errorf("amount_total: want 1234, got %d", session.AmountTotal)
	}
}

func TestProviderInitEnabledEvents(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
		"enabled_events":  []any{"checkout.session.completed"},
	})
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if len(p.enabledEvents) != 1 || p.enabledEvents[0] != "checkout.session.completed" {
		t.Errorf("enabled events not set: %v", p.enabledEvents)
	}
}

func TestProviderWebhookConfigPageHandler(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetDispatcher(&webhook.Dispatcher{URL: "http://localhost/webhook"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/_admin/stripe/webhooks", nil)
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok123"))
	p.WebhookConfigPageHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "http://localhost/webhook") {
		t.Error("page missing webhook url")
	}
	if !strings.Contains(body, "whsec_muara") {
		t.Error("page missing webhook secret")
	}
	if !strings.Contains(body, "checkout.session.completed") {
		t.Error("page missing event options")
	}
}

func TestProviderWebhookConfigSaveHandler(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yml")
	cfg := `providers:
  stripe:
    enabled: true
    config:
      publishable_key: pk_test_muara
      secret_key: sk_test_muara
      webhook_secret: whsec_muara
`
	if err := os.WriteFile(configPath, []byte(cfg), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetConfigPath(configPath)
	d := &webhook.Dispatcher{URL: "http://old.url"}
	p.SetDispatcher(d)

	form := url.Values{}
	form.Set("webhook_url", "http://new.url/webhook")
	form.Add("enabled_events", "checkout.session.completed")
	form.Add("enabled_events", "checkout.session.expired")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/webhooks", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok123"))
	p.WebhookConfigSaveHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d", rec.Code)
	}
	if d.URL != "http://new.url/webhook" {
		t.Errorf("dispatcher url not updated: %q", d.URL)
	}
	if len(d.EnabledEvents) != 2 {
		t.Errorf("dispatcher events: want 2, got %d", len(d.EnabledEvents))
	}

	// #nosec G304 -- test reads config file it just wrote
	updated, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read updated config: %v", err)
	}
	if !strings.Contains(string(updated), "http://new.url/webhook") {
		t.Error("config did not persist webhook_url")
	}
}

func TestProviderWebhookConfigSaveHandlerNoConfigPath(t *testing.T) {
	p := NewProvider()
	// #nosec G101 -- test fixture dummy credentials
	if err := p.Init(map[string]any{
		"publishable_key": "pk_test_muara",
		"secret_key":      "sk_test_muara",
		"webhook_secret":  "whsec_muara",
	}); err != nil {
		t.Fatalf("init: %v", err)
	}
	d := &webhook.Dispatcher{URL: "http://old.url"}
	p.SetDispatcher(d)

	form := url.Values{}
	form.Set("webhook_url", "http://new.url/webhook")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/webhooks", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(httputil.WithCSRFToken(req.Context(), "tok123"))
	p.WebhookConfigSaveHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d", rec.Code)
	}
	if d.URL != "http://new.url/webhook" {
		t.Errorf("dispatcher url not updated: %q", d.URL)
	}
}

func TestProviderWebhookEventOptions(t *testing.T) {
	p := NewProvider()
	p.enabledEvents = []string{"checkout.session.completed"}
	options := p.webhookEventOptions()
	if len(options) != 5 {
		t.Fatalf("want 5 options, got %d", len(options))
	}
	for _, opt := range options {
		if opt.Name == "checkout.session.completed" && !opt.Checked {
			t.Error("checkout.session.completed should be checked")
		}
		if opt.Name == "payment_intent.created" && opt.Checked {
			t.Error("payment_intent.created should not be checked")
		}
	}
}
