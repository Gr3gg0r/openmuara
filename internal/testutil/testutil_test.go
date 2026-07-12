package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenFile(t *testing.T) {
	data := GoldenFile(t, "golden.json")
	want := "{\"ok\":true}\n"
	if string(data) != want {
		t.Errorf("golden contents: want %q, got %q", want, data)
	}
}

func TestGoldenPath(t *testing.T) {
	path := GoldenPath(t, "golden.json")
	if path == "" {
		t.Fatal("expected non-empty path")
	}
}

func TestNewMemoryStores(t *testing.T) {
	stores := NewMemoryStores(t)
	if stores.Ledger == nil || stores.Audit == nil {
		t.Fatal("expected non-nil stores")
	}
}

func TestNewSQLiteStores(t *testing.T) {
	stores := NewSQLiteStores(t)
	if stores.Ledger == nil || stores.Audit == nil {
		t.Fatal("expected non-nil stores")
	}
}

func TestTempWorkspace(t *testing.T) {
	dir := TempWorkspace(t)
	if _, err := os.Stat(filepath.Join(dir, "data")); err != nil {
		t.Errorf("expected data dir: %v", err)
	}
}

func TestWriteFile(t *testing.T) {
	dir := TempWorkspace(t)
	WriteFile(t, dir, "nested/file.txt", []byte("hello"))

	// #nosec G304 -- test reads file written by helper under temp dir
	got, err := os.ReadFile(filepath.Join(dir, "nested/file.txt"))
	if err != nil {
		t.Fatalf("read written file: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("contents: want hello, got %q", got)
	}
}
