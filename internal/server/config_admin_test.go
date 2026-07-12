package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestConfigAdminHandlers_ViewerCanReadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	_ = os.WriteFile(path, []byte("server:\n  host: 127.0.0.1\n  port: 9000\n"), 0o600)

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	req := newViewerRequest(http.MethodGet, "/_admin/config")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected viewer can read config, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_ViewerCannotPatchProviders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	_ = os.WriteFile(path, []byte("providers:\n  fawry:\n    enabled: true\n"), 0o600)

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	payload := map[string]any{
		"providers": map[string]any{"fawry": map[string]any{"enabled": false}},
	}
	data, _ := json.Marshal(payload)
	req := newViewerRequest(http.MethodPatch, "/_admin/config/providers")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer patch, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_ViewerCannotReplayWebhook(t *testing.T) {
	mux := http.NewServeMux()
	AdminAPIHandlers(mux, RouterConfig{
		TransactionStore: engine.NewMemoryStore(),
	})

	req := newViewerRequest(http.MethodPost, "/_admin/transactions/tx-1/replay-webhook")
	req.SetPathValue("ref", "tx-1")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer replay, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_GetConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(`
server:
  host: 127.0.0.1
  port: 9000
providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant
      merchant_security_key: super-secret
webhook:
  url: http://example.com/webhook
  max_retries: 3
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	req := httptest.NewRequest(http.MethodGet, "/_admin/config", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	providers := body["providers"].(map[string]any)
	fawry := providers["fawry"].(map[string]any)
	if fawry["enabled"] != true {
		t.Fatalf("expected fawry enabled, got %v", fawry["enabled"])
	}
	cfg := fawry["config"].(map[string]any)
	if cfg["merchant_security_key"] != "***" {
		t.Fatalf("expected secret redacted, got %v", cfg["merchant_security_key"])
	}
	if cfg["merchant_code"] != "muara-merchant" {
		t.Fatalf("expected merchant_code preserved, got %v", cfg["merchant_code"])
	}
}

func TestConfigAdminHandlers_PatchProviders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(`
providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant
  stripe:
    enabled: false
    config:
      publishable_key: pk_test
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	payload := map[string]any{
		"providers": map[string]any{
			"fawry":  map[string]any{"enabled": false},
			"stripe": map[string]any{"enabled": true},
		},
	}
	data, _ := json.Marshal(payload)
	req := newAdminRequest(http.MethodPatch, "/_admin/config/providers")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	// #nosec G304 -- path is a temp file created by the test.
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "enabled: false") {
		t.Fatalf("expected fawry disabled in config file, got:\n%s", string(content))
	}
}

func TestConfigAdminHandlers_GetAndPatchWebhooks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(`
webhook:
  url: http://old.example.com
  max_retries: 1
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	payload := map[string]any{
		"url":         "http://new.example.com",
		"max_retries": 5,
		"targets": map[string]string{
			"stripe": "http://stripe.example.com",
		},
		"events": map[string][]string{
			"stripe": {"checkout.session.completed"},
		},
	}
	data, _ := json.Marshal(payload)
	req := newAdminRequest(http.MethodPatch, "/_admin/config/webhooks")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/_admin/config/webhooks", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["url"] != "http://new.example.com" {
		t.Fatalf("expected new url, got %v", body["url"])
	}
}

func TestConfigAdminHandlers_Reload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	_ = os.WriteFile(path, []byte("server:\n  host: 127.0.0.1\n"), 0o600)

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	req := newAdminRequest(http.MethodPost, "/_admin/config/reload")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_NoConfigPath(t *testing.T) {
	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: ""})

	req := httptest.NewRequest(http.MethodGet, "/_admin/config", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// No handler registered when ConfigPath is empty, so ServeMux returns 404.
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestConfigAdminHandlers_PatchConfigCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	original := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant
`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	payload := map[string]any{
		"providers": map[string]any{
			"fawry": map[string]any{"enabled": false},
		},
	}
	data, _ := json.Marshal(payload)
	req := newAdminRequest(http.MethodPatch, "/_admin/config/providers")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	// #nosec G304 -- path is a temp file created by the test.
	backup, err := os.ReadFile(path + ".bak")
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if !bytes.Equal(backup, original) {
		t.Fatalf("backup does not match original config:\n%s", string(backup))
	}
}

func TestConfigAdminHandlers_PatchConfigConflict(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	original := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant
`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	// Load the current checksum from the dashboard.
	req := httptest.NewRequest(http.MethodGet, "/_admin/config", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var getBody map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &getBody); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	checksum, _ := getBody["config_checksum"].(string)
	if checksum == "" {
		t.Fatal("expected config_checksum in response")
	}

	// Simulate an external edit between read and write.
	modified := []byte(`providers:
  fawry:
    enabled: false
    config:
      merchant_code: muara-merchant
`)
	if err := os.WriteFile(path, modified, 0o600); err != nil {
		t.Fatalf("write modified config: %v", err)
	}

	payload := map[string]any{
		"providers": map[string]any{
			"fawry": map[string]any{"enabled": false},
		},
		"checksum": checksum,
	}
	data, _ := json.Marshal(payload)
	req = newAdminRequest(http.MethodPatch, "/_admin/config/providers")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_PatchConfigConflictIfMatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	original := []byte(`providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant
`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	checksum := sha256Hex(original)

	// Simulate an external edit between read and write.
	modified := []byte(`providers:
  fawry:
    enabled: false
    config:
      merchant_code: muara-merchant
`)
	if err := os.WriteFile(path, modified, 0o600); err != nil {
		t.Fatalf("write modified config: %v", err)
	}

	payload := map[string]any{
		"providers": map[string]any{
			"fawry": map[string]any{"enabled": false},
		},
	}
	data, _ := json.Marshal(payload)
	req := newAdminRequest(http.MethodPatch, "/_admin/config/providers")
	req.Body = io.NopCloser(bytes.NewReader(data))
	req.Header.Set("If-Match", "\""+checksum+"\"")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestConfigAdminHandlers_TestWebhookSignatureAndLatency(t *testing.T) {
	var receivedSignature string
	var receivedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedSignature = r.Header.Get("X-Muara-Signature")
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: filepath.Join(t.TempDir(), "config.yml")})

	payload := map[string]any{
		"url":      server.URL,
		"provider": "stripe",
		"secret":   "test-secret",
	}
	data, _ := json.Marshal(payload)
	req := newAdminRequest(http.MethodPost, "/_admin/config/webhooks/test")
	req.Body = io.NopCloser(bytes.NewReader(data))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["success"] != true {
		t.Fatalf("expected success true, got %v", body["success"])
	}
	if body["status"] != 200.0 {
		t.Fatalf("expected status 200, got %v", body["status"])
	}
	latency, ok := body["latency_ms"].(float64)
	if !ok || latency < 0 {
		t.Fatalf("expected latency_ms >= 0, got %v", body["latency_ms"])
	}
	if body["signature_verified"] != true {
		t.Fatalf("expected signature_verified true, got %v", body["signature_verified"])
	}

	if receivedSignature == "" {
		t.Fatal("expected X-Muara-Signature header")
	}
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write(receivedBody)
	wantSig := hex.EncodeToString(mac.Sum(nil))
	if receivedSignature != wantSig {
		t.Fatalf("signature mismatch: got %s, want %s", receivedSignature, wantSig)
	}
}

func TestConfigAdminHandlers_GetConfigAdminPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(`
server:
  host: 127.0.0.1
  port: 9000
  admin_port: 9001
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	mux := http.NewServeMux()
	ConfigAdminHandlers(mux, RouterConfig{ConfigPath: path})

	req := httptest.NewRequest(http.MethodGet, "/_admin/config", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	server := body["server"].(map[string]any)
	if server["admin_port"] != 9001.0 {
		t.Fatalf("expected admin_port 9001, got %v", server["admin_port"])
	}
}
