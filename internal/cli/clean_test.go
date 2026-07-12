package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/config"
)

func TestCleanSQLiteRemovesFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "ledger.db")
	if err := os.WriteFile(dbPath, []byte("data"), 0o600); err != nil {
		t.Fatalf("create db file: %v", err)
	}

	var out bytes.Buffer
	if err := cleanSQLite(&out, &bytes.Buffer{}, dbPath, true); err != nil {
		t.Fatalf("cleanSQLite: %v", err)
	}

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatal("expected db file to be removed")
	}
	if !strings.Contains(out.String(), "cleaned:") {
		t.Fatalf("expected cleaned message, got %q", out.String())
	}
}

func TestCleanSQLiteMissingFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "missing.db")

	var out bytes.Buffer
	if err := cleanSQLite(&out, &bytes.Buffer{}, dbPath, true); err != nil {
		t.Fatalf("cleanSQLite: %v", err)
	}

	if !strings.Contains(out.String(), "does not exist") {
		t.Fatalf("expected does not exist message, got %q", out.String())
	}
}

func TestCleanSQLiteRequiresConfirmation(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "ledger.db")
	if err := os.WriteFile(dbPath, []byte("data"), 0o600); err != nil {
		t.Fatalf("create db file: %v", err)
	}

	var out bytes.Buffer
	in := strings.NewReader("no\n")
	if err := cleanSQLite(&out, in, dbPath, false); err != nil {
		t.Fatalf("cleanSQLite: %v", err)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected db file to remain when confirmation is denied")
	}
	if !strings.Contains(out.String(), "clean cancelled") {
		t.Fatalf("expected cancelled message, got %q", out.String())
	}
}

func TestRunCleanMemoryPersistence(t *testing.T) {
	rootConfigPath = filepath.Join(t.TempDir(), "config.yml")
	cfg := config.DefaultYAML()
	cfg = bytes.Replace(cfg, []byte("persistence:\n  type: sqlite"), []byte("persistence:\n  type: memory"), 1)
	if err := os.WriteFile(rootConfigPath, cfg, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var out bytes.Buffer
	if err := runClean(&out, &bytes.Buffer{}, true); err != nil {
		t.Fatalf("runClean: %v", err)
	}

	if !strings.Contains(out.String(), "in-memory") {
		t.Fatalf("expected in-memory message, got %q", out.String())
	}
}
