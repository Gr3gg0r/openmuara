package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempWorkspace creates a temporary directory that mimics a muara workspace.
// It returns the workspace path. The directory is automatically removed when
// the test ends.
func TempWorkspace(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o750); err != nil {
		t.Fatalf("create data dir: %v", err)
	}
	return dir
}

// WriteFile writes data to a path under the workspace and fails the test on error.
func WriteFile(t *testing.T, dir, relPath string, data []byte) {
	t.Helper()
	path := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("create parent dirs: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
