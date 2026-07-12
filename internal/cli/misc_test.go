package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/testutil"
)

func TestVersionCommand(t *testing.T) {
	cmd := newVersionCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected version output")
	}
}

func TestDoctorCommandExecution(t *testing.T) {
	cmd := newDoctorCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "go") {
		t.Errorf("output missing go check: %q", buf.String())
	}
}

func TestDoctorCommandMissingRequiredTool(t *testing.T) {
	doctorTools = []toolCheck{{name: "muara-nonexistent-binary"}}
	defer func() {
		doctorTools = []toolCheck{
			{name: "go"},
			{name: "golangci-lint"},
			{name: "govulncheck", optional: true},
			{name: "task"},
		}
	}()

	cmd := newDoctorCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetContext(context.Background())
	result := runDoctor(cmd, false, func(string) (string, error) { return "", errors.New("not found") })
	if result.Healthy {
		t.Fatal("expected unhealthy result for missing required tool")
	}
	printDoctorResult(&buf, result)
	if !strings.Contains(buf.String(), "FAIL") {
		t.Errorf("output missing FAIL: %q", buf.String())
	}
}

func TestInitCommand(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	oldPath := rootConfigPath
	rootConfigPath = dir + "/config.yml"
	defer func() { rootConfigPath = oldPath }()

	cmd := newInitCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "created") {
		t.Errorf("output missing created: %q", buf.String())
	}

	// Second run should report that config already exists.
	buf.Reset()
	cmd = newInitCommand()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute second time: %v", err)
	}
	if !strings.Contains(buf.String(), "already exists") {
		t.Errorf("output missing already exists: %q", buf.String())
	}
}

func TestStartCommandInvalidConfig(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	testutil.WriteFile(t, dir, "config.yml", []byte("server: ["))

	oldPath := rootConfigPath
	rootConfigPath = dir + "/config.yml"
	defer func() { rootConfigPath = oldPath }()

	cmd := newStartCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestAuditCommandWiring(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "log:\n  level: error\npersistence:\n  type: memory\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	oldPath := rootConfigPath
	rootConfigPath = dir + "/config.yml"
	defer func() { rootConfigPath = oldPath }()

	cmd := newAuditCommand()
	cmd.SetArgs([]string{"list"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array output, got %q", buf.String())
	}
}

func TestOpenAuditStoreSQLite(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	store, err := openAuditStore(config.PersistenceConfig{Type: "sqlite", Path: dir + "/audit.db"})
	if err != nil {
		t.Fatalf("open sqlite audit store: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
	if closer, ok := store.(io.Closer); ok {
		_ = closer.Close()
	}
}

func TestOpenAuditStoreUnsupported(t *testing.T) {
	_, err := openAuditStore(config.PersistenceConfig{Type: "redis"})
	if err == nil {
		t.Fatal("expected error for unsupported persistence type")
	}
}

func TestRootCommandVersion(t *testing.T) {
	cmd := newRootCommand()
	cmd.SetArgs([]string{"version"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected version output")
	}
}

func TestScenarioCommandWiring(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"--config", cfgPath, "scenario", "success", "ref-1"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "ok") {
		t.Errorf("expected ok in output, got %q", buf.String())
	}
}

func TestWebhookCommandWiring(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)

	cmd := newRootCommand()
	cmd.SetArgs([]string{"--config", cfgPath, "webhook", "list"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "results") {
		t.Errorf("expected results in output, got %q", buf.String())
	}
}

func TestVersionJSONOutput(t *testing.T) {
	cmd := newRootCommand()
	cmd.SetArgs([]string{"--json", "version"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("parse json: %v\noutput: %s", err, buf.String())
	}
	if result["version"] == "" {
		t.Errorf("expected version field, got %v", result)
	}
}

func TestDoctorJSONOutput(t *testing.T) {
	cmd := newRootCommand()
	cmd.SetArgs([]string{"doctor", "--json"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("parse json: %v\noutput: %s", err, buf.String())
	}
	if _, ok := result["healthy"]; !ok {
		t.Errorf("expected healthy field, got %v", result)
	}
	if _, ok := result["tools"]; !ok {
		t.Errorf("expected tools field, got %v", result)
	}
	if _, ok := result["config"]; !ok {
		t.Errorf("expected config field, got %v", result)
	}
}

func TestQuietFlagSuppressesOutput(t *testing.T) {
	cmd := newRootCommand()
	cmd.SetArgs([]string{"--quiet", "version"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected no stdout output, got %q", buf.String())
	}
}

func TestCommandsHaveExamples(t *testing.T) {
	root := newRootCommand()
	if root.Example == "" {
		t.Error("root command missing Example")
	}
	for _, cmd := range root.Commands() {
		if cmd.Example == "" {
			t.Errorf("command %q missing Example", cmd.Name())
		}
		for _, sub := range cmd.Commands() {
			if sub.Example == "" {
				t.Errorf("subcommand %s %q missing Example", cmd.Name(), sub.Name())
			}
		}
	}
}
