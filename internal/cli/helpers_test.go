package cli

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"log/slog"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/engine"
	_ "github.com/openmuara/openmuara/internal/provider/defaultplugin"
)

func TestParseLogLevel(t *testing.T) {
	cases := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo},
	}

	for _, tc := range cases {
		got := parseLogLevel(tc.input)
		if got != tc.want {
			t.Errorf("parseLogLevel(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestNewPersistenceStoresMemory(t *testing.T) {
	ledger, auditStore, _, closeFn, err := newPersistenceStores(config.PersistenceConfig{Type: "memory"})
	if err != nil {
		t.Fatalf("new memory stores: %v", err)
	}
	defer func() { _ = closeFn() }()
	if _, ok := ledger.(*engine.MemoryStore); !ok {
		t.Errorf("expected *engine.MemoryStore, got %T", ledger)
	}
	if auditStore == nil {
		t.Error("expected audit store")
	}
}

func TestNewPersistenceStoresSQLite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ledger.db")
	ledger, _, _, closeFn, err := newPersistenceStores(config.PersistenceConfig{Type: "sqlite", Path: path})
	if err != nil {
		t.Fatalf("new sqlite stores: %v", err)
	}
	defer func() { _ = closeFn() }()
	if _, ok := ledger.(*engine.SQLiteStore); !ok {
		t.Errorf("expected *engine.SQLiteStore, got %T", ledger)
	}
}

func TestNewPersistenceStoresDefault(t *testing.T) {
	ledger, _, _, closeFn, err := newPersistenceStores(config.PersistenceConfig{Type: ""})
	if err != nil {
		t.Fatalf("new default stores: %v", err)
	}
	defer func() { _ = closeFn() }()
	if _, ok := ledger.(*engine.MemoryStore); !ok {
		t.Errorf("expected *engine.MemoryStore for empty type, got %T", ledger)
	}
}

func TestNewPersistenceStoresUnsupported(t *testing.T) {
	_, _, _, _, err := newPersistenceStores(config.PersistenceConfig{Type: "postgres"})
	if err == nil {
		t.Fatal("expected error for unsupported persistence type")
	}
}

func TestPrintResponsePrettyJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"status":"ok"}`)),
	}
	var buf bytes.Buffer
	if err := printResponse(resp, &buf); err != nil {
		t.Fatalf("printResponse: %v", err)
	}
	if !strings.Contains(buf.String(), "status") {
		t.Errorf("expected pretty JSON output, got %q", buf.String())
	}
}

func TestPrintResponseErrorStatus(t *testing.T) {
	resp := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(`{"error":"bad"}`)),
	}
	var buf bytes.Buffer
	err := printResponse(resp, &buf)
	if err == nil {
		t.Fatal("expected error for 4xx response")
	}
}

func TestActiveProviderName(t *testing.T) {
	cfg := &config.Config{
		Providers: map[string]config.ProviderConfig{
			"default": {Enabled: true},
		},
	}
	got := activeProviderName(cfg, nil)
	if got != "default" {
		t.Errorf("activeProviderName: want default, got %q", got)
	}
}

func TestActiveProviderNameFallback(t *testing.T) {
	cfg := &config.Config{Providers: map[string]config.ProviderConfig{}}
	got := activeProviderName(cfg, nil)
	if got != "fawry" {
		t.Errorf("activeProviderName fallback: want fawry, got %q", got)
	}
}

func TestInitCommandCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, ".muara", "config.yml")
	t.Cleanup(func() { rootConfigPath = old })

	cmd := newInitCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if _, err := os.Stat(rootConfigPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if !strings.Contains(buf.String(), "created") {
		t.Errorf("expected created message, got %q", buf.String())
	}
}

func TestInitCommandSkipsExistingConfig(t *testing.T) {
	dir := t.TempDir()
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, ".muara", "config.yml")
	t.Cleanup(func() { rootConfigPath = old })

	if err := os.MkdirAll(filepath.Dir(rootConfigPath), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(rootConfigPath, []byte("existing"), 0o600); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	cmd := newInitCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !strings.Contains(buf.String(), "already exists") {
		t.Errorf("expected already exists message, got %q", buf.String())
	}
}
