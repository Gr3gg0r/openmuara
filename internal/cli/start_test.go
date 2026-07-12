package cli

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/defaultplugin"
	"github.com/openmuara/openmuara/internal/server"
	"github.com/openmuara/openmuara/internal/testutil"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestNewStartRuntime_Memory(t *testing.T) {
	cfg := &config.Config{
		Server:      config.ServerConfig{Host: "127.0.0.1", Port: 0},
		Persistence: config.PersistenceConfig{Type: "memory"},
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
		},
	}

	registry := provider.NewRegistry()
	defaultplugin.RegisterWith(registry)

	rt, err := newStartRuntime(cfg, registry, server.New)
	if err != nil {
		t.Fatalf("newStartRuntime: %v", err)
	}
	if rt == nil {
		t.Fatal("expected runtime")
	}
	if len(rt.enabled) != 1 || rt.enabled[0] != "default" {
		t.Errorf("enabled providers: want [default], got %v", rt.enabled)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	errCh := make(chan error, 1)
	go func() { errCh <- runStart(ctx, rt, &buf) }()

	// Wait for the server to bind a random port.
	deadline := time.Now().Add(2 * time.Second)
	for rt.srv.Addr() == "" && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if rt.srv.Addr() == "" {
		t.Fatal("server did not bind")
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("runStart: %v", err)
	}
	if !strings.Contains(buf.String(), "starting muara...") {
		t.Errorf("output missing startup message: %q", buf.String())
	}
}

func TestBuildProviderDispatcher(t *testing.T) {
	p := defaultplugin.NewProvider()

	if d := buildProviderDispatcher(config.WebhookConfig{}, "default", nil, p, nil, nil); d != nil {
		t.Error("expected nil dispatcher when webhook URL is empty")
	}

	store := webhook.NewMemoryStore()
	d := buildProviderDispatcher(config.WebhookConfig{URL: "http://example.com/webhook"}, "default", nil, p, store, nil)
	if d == nil {
		t.Fatal("expected dispatcher")
	}
	if d.Store != store {
		t.Error("expected dispatcher to use provided attempt store")
	}

	d2 := buildProviderDispatcher(config.WebhookConfig{
		URL:     "http://example.com/webhook",
		Targets: map[string]string{"default": "http://target.example.com"},
	}, "default", nil, p, store, nil)
	if d2 == nil {
		t.Fatal("expected dispatcher with target override")
	}
}

func TestNewPersistenceStores(t *testing.T) {
	ledger, auditStore, _, closeStores, err := newPersistenceStores(config.PersistenceConfig{Type: "memory"})
	if err != nil {
		t.Fatalf("memory stores: %v", err)
	}
	if ledger == nil || auditStore == nil {
		t.Fatal("expected non-nil stores")
	}
	if err := closeStores(); err != nil {
		t.Errorf("close memory stores: %v", err)
	}

	_, _, _, _, err = newPersistenceStores(config.PersistenceConfig{Type: "redis"})
	if err == nil {
		t.Fatal("expected error for unsupported persistence type")
	}

	dir := testutil.TempWorkspace(t)
	ledger, auditStore, _, closeStores, err = newPersistenceStores(config.PersistenceConfig{
		Type: "sqlite",
		Path: filepath.Join(dir, "ledger.db"),
	})
	if err != nil {
		t.Fatalf("sqlite stores: %v", err)
	}
	if ledger == nil || auditStore == nil {
		t.Fatal("expected non-nil sqlite stores")
	}
	if err := closeStores(); err != nil {
		t.Errorf("close sqlite stores: %v", err)
	}
}

func TestNewStartRuntime_LoadEnabledProvidersError(t *testing.T) {
	cfg := &config.Config{
		Server:      config.ServerConfig{Host: "127.0.0.1", Port: 0},
		Persistence: config.PersistenceConfig{Type: "memory"},
		Providers: map[string]config.ProviderConfig{
			"unknown": {Enabled: true},
		},
	}

	_, err := newStartRuntime(cfg, provider.NewRegistry(), server.New)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestNewStartCommand(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfgPath := filepath.Join(dir, "config.yml")
	port := freePort(t)
	cfgBody := []byte(`server:
  host: 127.0.0.1
  port: ` + itoa(port) + `
persistence:
  type: memory
providers:
  default:
    enabled: true
`)
	if err := os.WriteFile(cfgPath, cfgBody, 0o600); err != nil {
		t.Fatal(err)
	}

	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newStartCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.ExecuteContext(ctx)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()

	if err := <-errCh; err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "starting muara...") {
		t.Errorf("missing startup message: %q", buf.String())
	}
}

func TestNewStartCommandDryRun(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfgPath := filepath.Join(dir, "config.yml")
	cfgBody := []byte(`server:
  host: 127.0.0.1
  port: 9000
persistence:
  type: memory
providers:
  default:
    enabled: true
`)
	if err := os.WriteFile(cfgPath, cfgBody, 0o600); err != nil {
		t.Fatal(err)
	}

	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newStartCommand()
	cmd.SetArgs([]string{"--dry-run"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "config loaded and validated") {
		t.Errorf("missing dry-run message: %q", buf.String())
	}
}

func TestNewStartCommandDryRunJSON(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfgPath := filepath.Join(dir, "config.yml")
	cfgBody := []byte(`server:
  host: 127.0.0.1
  port: 9000
persistence:
  type: memory
providers:
  default:
    enabled: true
`)
	if err := os.WriteFile(cfgPath, cfgBody, 0o600); err != nil {
		t.Fatal(err)
	}

	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	jsonOutput = true
	defer func() { jsonOutput = false }()

	cmd := newStartCommand()
	cmd.SetArgs([]string{"--dry-run"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "\"ok\": true") {
		t.Errorf("missing JSON ok field: %q", buf.String())
	}
}

func TestNewStartRuntime_DualPort(t *testing.T) {
	cfg := &config.Config{
		Server:      config.ServerConfig{Host: "127.0.0.1", Port: 0, AdminPort: 0},
		Persistence: config.PersistenceConfig{Type: "memory"},
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
		},
	}

	registry := provider.NewRegistry()
	defaultplugin.RegisterWith(registry)

	rt, err := newStartRuntime(cfg, registry, server.New)
	if err != nil {
		t.Fatalf("newStartRuntime: %v", err)
	}
	if rt.adminSrv != nil {
		t.Fatal("expected no admin server when AdminPort is 0")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	errCh := make(chan error, 1)
	go func() { errCh <- runStart(ctx, rt, &buf) }()

	waitForServer(t, rt.srv)

	// Verify single-port mode serves admin API.
	resp, err := http.Get("http://" + rt.srv.Addr() + "/_admin/transactions")
	if err != nil {
		t.Fatalf("get transactions: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("single port should serve /_admin/transactions: want 200, got %d", resp.StatusCode)
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("runStart: %v", err)
	}
}

func TestNewStartRuntime_DualPortSplit(t *testing.T) {
	providerPort := freePort(t)
	adminPort := freePort(t)

	cfg := &config.Config{
		Server:      config.ServerConfig{Host: "127.0.0.1", Port: providerPort, AdminPort: adminPort},
		Persistence: config.PersistenceConfig{Type: "memory"},
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
		},
	}

	registry := provider.NewRegistry()
	defaultplugin.RegisterWith(registry)

	rt, err := newStartRuntime(cfg, registry, server.New)
	if err != nil {
		t.Fatalf("newStartRuntime: %v", err)
	}
	if rt.adminSrv == nil {
		t.Fatal("expected admin server when AdminPort is set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	errCh := make(chan error, 1)
	go func() { errCh <- runStart(ctx, rt, &buf) }()

	waitForServer(t, rt.srv)
	waitForServer(t, rt.adminSrv)

	// Provider port should not serve admin API.
	resp, err := http.Get("http://" + rt.srv.Addr() + "/_admin/transactions")
	if err != nil {
		t.Fatalf("get transactions from provider port: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("provider port should not serve /_admin/transactions: want 404, got %d", resp.StatusCode)
	}

	// Admin port should serve admin API.
	resp, err = http.Get("http://" + rt.adminSrv.Addr() + "/_admin/transactions")
	if err != nil {
		t.Fatalf("get transactions from admin port: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin port should serve /_admin/transactions: want 200, got %d", resp.StatusCode)
	}

	// Provider port should still serve provider emulation.
	resp, err = http.Post("http://"+rt.srv.Addr()+"/default/charge", "application/json", nil)
	if err != nil {
		t.Fatalf("post charge to provider port: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("provider port should serve /default/charge: want 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(buf.String(), "provider API:") || !strings.Contains(buf.String(), "admin API:") {
		t.Errorf("output missing dual-port URLs: %q", buf.String())
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("runStart: %v", err)
	}
}

func TestProviderBaseURLPrefersPublicBaseURL(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 9000, PublicBaseURL: "https://muara.example.com"},
	}
	if got := providerBaseURL(cfg); got != "https://muara.example.com" {
		t.Fatalf("want public base URL, got %q", got)
	}
}

func TestProviderBaseURLFallsBackToBindAddress(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 8080},
	}
	if got := providerBaseURL(cfg); got != "http://127.0.0.1:8080" {
		t.Fatalf("want bind address, got %q", got)
	}
}

func TestProviderBaseURLTrimsTrailingSlash(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 9000, PublicBaseURL: "https://muara.example.com/"},
	}
	if got := providerBaseURL(cfg); got != "https://muara.example.com" {
		t.Fatalf("want trimmed URL, got %q", got)
	}
}

func TestAdminBaseURLPrefersAdminPublicBaseURL(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:               "127.0.0.1",
			Port:               9000,
			AdminPort:          9001,
			PublicBaseURL:      "https://muara.example.com",
			AdminPublicBaseURL: "https://admin.muara.example.com",
		},
	}
	if got := adminBaseURLFromConfig(cfg); got != "https://admin.muara.example.com" {
		t.Fatalf("want admin public base URL, got %q", got)
	}
}

func TestAdminBaseURLUsesPublicBaseURLAndAdminPort(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:          "127.0.0.1",
			Port:          9000,
			AdminPort:     9001,
			PublicBaseURL: "https://muara.example.com",
		},
	}
	if got := adminBaseURLFromConfig(cfg); got != "https://muara.example.com:9001" {
		t.Fatalf("want public base URL with admin port, got %q", got)
	}
}

func TestAdminBaseURLFallsBackToBindAddress(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:      "127.0.0.1",
			Port:      9000,
			AdminPort: 9001,
		},
	}
	if got := adminBaseURLFromConfig(cfg); got != "http://127.0.0.1:9001" {
		t.Fatalf("want bind address with admin port, got %q", got)
	}
}

// waitForServer polls until srv has bound a non-zero port and is accepting connections.
func waitForServer(t *testing.T, srv *server.Server) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for (srv.Addr() == "" || strings.HasSuffix(srv.Addr(), ":0")) && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if srv.Addr() == "" || strings.HasSuffix(srv.Addr(), ":0") {
		t.Fatal("server did not bind to a usable address")
	}
	addr := srv.Addr()
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + addr + "/healthz")
		if err == nil {
			_ = resp.Body.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server at %s did not become ready", addr)
}

func TestNewStartCommand_LoadError(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfgPath := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("not: valid: yaml: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newStartCommand()
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when config cannot be loaded")
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
