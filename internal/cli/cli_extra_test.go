package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/defaultplugin"
	"github.com/Gr3gg0r/openmuara/internal/testutil"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestDefaultPromptFunc(t *testing.T) {
	in := strings.NewReader("answer\n")
	prompt := defaultPromptFunc(in)
	got, err := prompt("question: ")
	if err != nil {
		t.Fatalf("prompt: %v", err)
	}
	if got != "answer" {
		t.Errorf("answer: want answer, got %q", got)
	}
}

func TestDefaultPromptFuncEOF(t *testing.T) {
	in := strings.NewReader("")
	prompt := defaultPromptFunc(in)
	_, err := prompt("question: ")
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestRunWizard(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("1\nhttp://localhost/webhook\ninfo\n")
	choices, webhookURL, logLevel, err := runWizard(&out, defaultPromptFunc(in))
	if err != nil {
		t.Fatalf("runWizard: %v", err)
	}
	if len(choices) == 0 || choices[0].Key == "" {
		t.Error("expected non-empty choice key")
	}
	if webhookURL != "http://localhost/webhook" {
		t.Errorf("webhookURL: want http://localhost/webhook, got %q", webhookURL)
	}
	if logLevel != "info" {
		t.Errorf("logLevel: want info, got %q", logLevel)
	}
}

func TestRunWizardByName(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("fawry\n\n\n")
	choices, _, logLevel, err := runWizard(&out, defaultPromptFunc(in))
	if err != nil {
		t.Fatalf("runWizard: %v", err)
	}
	if len(choices) != 1 || choices[0].Key != "fawry" {
		t.Errorf("choice: want fawry, got %v", choices)
	}
	if logLevel != "info" {
		t.Errorf("logLevel default: want info, got %q", logLevel)
	}
}

func TestRunWizardInvalidChoice(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("999\n\n\n")
	choices, _, _, err := runWizard(&out, defaultPromptFunc(in))
	if err != nil {
		t.Fatalf("runWizard: %v", err)
	}
	if len(choices) == 0 || choices[0].Key == "" {
		t.Error("expected fallback choice key")
	}
}

func TestInitCommandWizard(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, "config.yml")
	defer func() { rootConfigPath = old }()

	cmd := newInitCommand()
	cmd.SetIn(strings.NewReader("\n\n\n"))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if _, err := os.Stat(rootConfigPath); err != nil {
		t.Fatalf("config not created: %v", err)
	}
}

func TestToStringSlice(t *testing.T) {
	got := toStringSlice([]any{"a", 1, "b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestProviderWebhookURLConfigured(t *testing.T) {
	providers := map[string]config.ProviderConfig{
		"stripe": {
			Enabled: true,
			Config:  map[string]any{"webhook_url": "http://localhost/webhook"},
		},
	}
	if !providerWebhookURLConfigured(providers) {
		t.Error("expected true when stripe webhook_url is set")
	}

	providers["stripe"] = config.ProviderConfig{Enabled: true, Config: map[string]any{}}
	if providerWebhookURLConfigured(providers) {
		t.Error("expected false when webhook_url missing")
	}

	providers["default"] = config.ProviderConfig{Enabled: true, Config: map[string]any{"webhook_url": "http://localhost/webhook"}}
	if providerWebhookURLConfigured(providers) {
		t.Error("expected false for non-stripe provider")
	}
}

func TestBuildProviderDispatcherWithProviderConfig(t *testing.T) {
	p := defaultplugin.NewProvider()
	store := webhook.NewMemoryStore()
	d := buildProviderDispatcher(config.WebhookConfig{URL: "http://example.com"}, "default", map[string]any{
		"webhook_url":    "http://provider.example.com",
		"enabled_events": []any{"invoice.payment_succeeded", 123},
	}, p, store, nil)
	if d == nil {
		t.Fatal("expected dispatcher")
	}
	if d.URL != "http://provider.example.com" {
		t.Errorf("url: want provider url, got %q", d.URL)
	}
	if len(d.EnabledEvents) != 1 || d.EnabledEvents[0] != "invoice.payment_succeeded" {
		t.Errorf("enabled events: %v", d.EnabledEvents)
	}
}

func TestActiveProviderNameNoEscapeHandler(t *testing.T) {
	registry := provider.NewRegistry()
	defaultplugin.RegisterWith(registry)

	cfg := &config.Config{
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
		},
	}
	if got := activeProviderName(cfg, nil); got != "default" {
		t.Errorf("activeProviderName: want default, got %q", got)
	}
}

func TestNewPersistenceStoresSQLiteDefaultPath(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(old) }()

	ledger, auditStore, webhookStore, closeFn, err := newPersistenceStores(config.PersistenceConfig{Type: "sqlite"})
	if err != nil {
		t.Fatalf("sqlite default path: %v", err)
	}
	defer func() { _ = closeFn() }()
	if ledger == nil || auditStore == nil || webhookStore == nil {
		t.Fatal("expected non-nil stores")
	}
}

func TestNewPersistenceStoresSQLiteMkdirError(t *testing.T) {
	_, _, _, _, err := newPersistenceStores(config.PersistenceConfig{Type: "sqlite", Path: "/dev/null/ledger.db"})
	if err == nil {
		t.Fatal("expected error when data dir cannot be created")
	}
}

func TestRunPluginsListNoPlugins(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	var buf bytes.Buffer
	if err := runPluginsList(&buf); err != nil {
		t.Fatalf("runPluginsList: %v", err)
	}
	if !strings.Contains(buf.String(), "No plugins discovered") {
		t.Errorf("expected no-plugins message, got %q", buf.String())
	}
}

func TestRunPluginsValidateNoPlugins(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	var buf bytes.Buffer
	if err := runPluginsValidate(&buf, nil); err != nil {
		t.Fatalf("runPluginsValidate: %v", err)
	}
	if !strings.Contains(buf.String(), "No plugins to validate") {
		t.Errorf("expected no-plugins message, got %q", buf.String())
	}
}

func TestLoadSinglePluginEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if _, err := loadSinglePlugin(dir); err == nil {
		t.Fatal("expected error for empty plugin dir")
	}
}

func TestServerBaseURLLoadError(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("not: valid: yaml: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	old := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = old }()

	cmd := newWebhookListCommand()
	if _, err := serverBaseURL(cmd); err == nil {
		t.Fatal("expected error loading invalid config")
	}
}
