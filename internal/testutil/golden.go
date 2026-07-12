package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// GoldenFile reads a file from the package's testdata directory. The path is
// relative to testdata/, e.g., GoldenFile(t, "charge_response.json").
func GoldenFile(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	// #nosec G304 -- test helper loads golden files from testdata
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file %q: %v", path, err)
	}
	return data
}

// GoldenPath returns the absolute path to a testdata file.
func GoldenPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("testdata", name)
}
