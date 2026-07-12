package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/fawry"
	"github.com/openmuara/openmuara/internal/provider/defaultplugin"
	"github.com/openmuara/openmuara/internal/testutil"
	"github.com/openmuara/openmuara/internal/webhook"
	"github.com/spf13/cobra"
)

var _ = defaultplugin.NewProvider

func TestActiveProviderNameFallbackRegression(t *testing.T) {
	cfg := &config.Config{Providers: map[string]config.ProviderConfig{}}
	if got := activeProviderName(cfg, nil); got != "fawry" {
		t.Errorf("want fawry fallback, got %q", got)
	}
}

func TestActiveProviderNamePicksEscapeProvider(t *testing.T) {
	cfg := &config.Config{
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
			"fawry":   {Enabled: true},
		},
	}
	loaded := []config.LoadedProvider{
		{Name: "default", Provider: defaultplugin.NewProvider()},
		{Name: "fawry", Provider: fawry.NewProvider()},
	}
	if got := activeProviderName(cfg, loaded); got != "fawry" {
		t.Errorf("want fawry (has escape handler), got %q", got)
	}
}

func TestProviderWebhookURLConfiguredRegression(t *testing.T) {
	providers := map[string]config.ProviderConfig{
		"stripe": {Enabled: true, Config: map[string]any{"webhook_url": "http://stripe.example.com"}},
	}
	if !providerWebhookURLConfigured(providers) {
		t.Error("expected configured when stripe has webhook_url")
	}
	if providerWebhookURLConfigured(map[string]config.ProviderConfig{"default": {Enabled: true}}) {
		t.Error("expected not configured for default provider")
	}
}

func TestToStringSliceRegression(t *testing.T) {
	got := toStringSlice([]any{"a", 1, "b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("want [a b], got %v", got)
	}
}

func TestBuildProviderDispatcherWithEvents(t *testing.T) {
	p := defaultplugin.NewProvider()
	d := buildProviderDispatcher(config.WebhookConfig{
		URL: "http://example.com",
	}, "default", map[string]any{
		"enabled_events": []any{"a", "b"},
	}, p, webhook.NewMemoryStore(), engine.NewMemoryStore())
	if d == nil {
		t.Fatal("expected dispatcher")
	}
	if len(d.EnabledEvents) != 2 {
		t.Errorf("want 2 enabled events, got %v", d.EnabledEvents)
	}
}

func writeConfigWithWebhook(t *testing.T, url string) string {
	t.Helper()
	dir := testutil.TempWorkspace(t)
	body := fmt.Sprintf("server:\n  host: 127.0.0.1\n  port: 9000\npersistence:\n  type: memory\nwebhook:\n  url: %s\n", url)
	testutil.WriteFile(t, dir, "config.yml", []byte(body))
	return filepath.Join(dir, "config.yml")
}

func TestRunDoctorWebhookReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer srv.Close()

	oldPath := rootConfigPath
	rootConfigPath = writeConfigWithWebhook(t, srv.URL)
	defer func() { rootConfigPath = oldPath }()

	result := runDoctor(testDoctorCmd(), true, func(string) (string, error) { return "/usr/bin/tool", nil })
	if !result.Webhook.Configured {
		t.Error("expected webhook configured")
	}
	if result.Webhook.Reachable == nil || !*result.Webhook.Reachable {
		t.Errorf("expected webhook reachable, got %v", result.Webhook.Reachable)
	}
}

func TestRunDoctorWebhookUnreachable(t *testing.T) {
	oldPath := rootConfigPath
	rootConfigPath = writeConfigWithWebhook(t, "http://localhost:1")
	defer func() { rootConfigPath = oldPath }()

	result := runDoctor(testDoctorCmd(), true, func(string) (string, error) { return "/usr/bin/tool", nil })
	if result.Webhook.Reachable == nil || *result.Webhook.Reachable {
		t.Errorf("expected webhook unreachable, got %v", result.Webhook.Reachable)
	}
}

func TestIsWebhookReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer srv.Close()
	if !isWebhookReachable(context.Background(), srv.URL) {
		t.Error("expected reachable")
	}
	if isWebhookReachable(context.Background(), "http://localhost:1") {
		t.Error("expected unreachable")
	}
}

func TestPrintDoctorResultBranches(t *testing.T) {
	var buf bytes.Buffer
	reachable := true
	result := doctorResult{
		Healthy: true,
		Tools:   []doctorTool{{Name: "go", Found: true}},
		Config:  doctorConfig{OK: true, Valid: true},
		Webhook: doctorWebhook{URL: "http://example.com", Configured: true, Reachable: &reachable},
	}
	printDoctorResult(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "config is valid") {
		t.Errorf("expected valid config message, got:\n%s", out)
	}
	if !strings.Contains(out, "webhook: reachable") {
		t.Errorf("expected reachable message, got:\n%s", out)
	}

	buf.Reset()
	unreachable := false
	result.Webhook.Reachable = &unreachable
	printDoctorResult(&buf, result)
	if !strings.Contains(buf.String(), "webhook: unreachable") {
		t.Errorf("expected unreachable message, got:\n%s", buf.String())
	}
}

func TestInitCommandDefaultsCreatesFile(t *testing.T) {
	dir := t.TempDir()
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, ".muara", "config.yml")
	defer func() { rootConfigPath = old }()

	cmd := newInitCommand()
	cmd.SetArgs([]string{"--defaults"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if _, err := os.Stat(rootConfigPath); err != nil {
		t.Errorf("config file should be created: %v", err)
	}
}

func TestInitCommandExistingNoForce(t *testing.T) {
	dir := t.TempDir()
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, ".muara", "config.yml")
	defer func() { rootConfigPath = old }()

	if err := os.MkdirAll(filepath.Dir(rootConfigPath), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(rootConfigPath, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := newInitCommand()
	cmd.SetArgs([]string{"--defaults"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	// #nosec G304 -- rootConfigPath is set by the test to a known temp file.
	data, _ := os.ReadFile(rootConfigPath)
	if string(data) != "existing" {
		t.Error("expected existing config to be preserved")
	}
}

func TestRunWizardIgnoresInvalidInput(t *testing.T) {
	inputs := []string{"invalid,999,default", "", ""}
	idx := 0
	prompt := func(string) (string, error) {
		v := inputs[idx]
		idx++
		return v, nil
	}
	var buf bytes.Buffer
	choices, _, _, err := runWizard(&buf, prompt)
	if err != nil {
		t.Fatalf("run wizard: %v", err)
	}
	keys := make(map[string]bool)
	for _, c := range choices {
		keys[c.Key] = true
	}
	if !keys["default"] {
		t.Errorf("expected default selection after invalid input, got %v", keys)
	}
}

func TestRunWizardEOF(t *testing.T) {
	prompt := func(string) (string, error) { return "", io.EOF }
	_, _, _, err := runWizard(&bytes.Buffer{}, prompt)
	if err == nil {
		t.Fatal("expected EOF error")
	}
}

func TestPluginsValidateSinglePathNoPlugin(t *testing.T) {
	dir := t.TempDir()
	cmd := newPluginsValidateCommand()
	cmd.SetArgs([]string{dir})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for empty plugin directory")
	}
}

func TestSecurityAuditNoIssues(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "server:\n  host: 127.0.0.1\n  port: 9000\npersistence:\n  type: memory\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, "config.yml")
	defer func() { rootConfigPath = old }()

	cmd := newSecurityAuditCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "no issues detected") {
		t.Errorf("expected no issues, got:\n%s", buf.String())
	}
}

func TestSecurityAuditHardenedWithoutAdmin(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "server:\n  host: 127.0.0.1\n  port: 9000\npersistence:\n  type: memory\nhardened: true\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, "config.yml")
	defer func() { rootConfigPath = old }()

	cmd := newSecurityAuditCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "hardened mode requires admin.enabled=true") {
		t.Errorf("expected hardened issue, got:\n%s", buf.String())
	}
}

func TestWebhookServerBaseURLError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte("not: valid: yaml: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	old := rootConfigPath
	rootConfigPath = path
	defer func() { rootConfigPath = old }()

	if _, err := serverBaseURL(&cobra.Command{}); err == nil {
		t.Fatal("expected error when config cannot be loaded")
	}
}

func TestWebhookPrintResponseNonJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	var buf bytes.Buffer
	if err := printResponse(resp, &buf); err == nil {
		t.Fatal("expected error for 500 response")
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected body in output, got:\n%s", buf.String())
	}
}

func TestNewStartCommandLoadError(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, "config.yml")
	defer func() { rootConfigPath = old }()

	if err := os.WriteFile(rootConfigPath, []byte("not: valid: yaml: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := newStartCommand()
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestMaybeCheckUpdateSkipsConfigDisabled(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "server:\n  host: 127.0.0.1\n  port: 9000\npersistence:\n  type: memory\ndisable_update_check: true\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, "config.yml")
	defer func() { rootConfigPath = old }()

	latest, update := maybeCheckUpdate(testVersionCmd())
	if latest != "" || update {
		t.Errorf("expected update check skipped, got latest=%q update=%v", latest, update)
	}
}

func TestDoctorTimestamp(t *testing.T) {
	result := runDoctor(testDoctorCmd(), false, func(string) (string, error) { return "", fmt.Errorf("missing") })
	if result.Timestamp == "" {
		t.Error("expected timestamp")
	}
	if _, err := time.Parse(time.RFC3339, result.Timestamp); err != nil {
		t.Errorf("timestamp not RFC3339: %v", err)
	}
}
