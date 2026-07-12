// Package conform provides contract conformance tests for registered providers.
package conform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

// Snapshot captures the static contract surface of a provider.
type Snapshot struct {
	Name     string   `json:"name"`
	Routes   []Route  `json:"routes"`
	Versions []string `json:"versions,omitempty"`
}

// Route is a serializable route description.
type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Capture returns a serializable snapshot of the provider.
// It attempts Init with a generic placeholder config so that providers that
// lazily build handlers in Routes() are in a safe state. If Init fails, only
// the name and versions (when safe) are captured and the error is returned.
func Capture(p provider.Provider) (Snapshot, error) {
	cfg := map[string]any{
		"merchant_code":         "TEST_MC",
		"merchant_security_key": "TEST_MSK",
		"webhook_secret":        "TEST_WS",
		"api_key":               "TEST_API_KEY",
		"secret_key":            "TEST_SECRET_KEY",
		"x_signature":           "TEST_XSIG",
		"x_signature_key":       "TEST_XSIG_KEY",
		"collection_id":         "TEST_COLLECTION",
		"merchant_key":          "TEST_MERCHANT_KEY",
		"publishable_key":       "TEST_PK",
		"user_secret_key":       "TEST_USER_SECRET",
		"version":               "v1",
	}
	initErr := p.Init(cfg)

	s := Snapshot{Name: p.Name()}
	if vp, ok := p.(provider.VersionedProvider); ok {
		s.Versions = vp.Versions()
	}
	if initErr != nil {
		return s, initErr
	}
	for _, r := range p.Routes() {
		s.Routes = append(s.Routes, Route{Method: r.Method, Path: r.Path})
	}
	return s, nil
}

// GoldenPath returns the path to the golden file for a provider.
func GoldenPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("testdata", "golden", name+".json")
}

// Compare reads the golden file for p and compares it to the current snapshot.
// Set -update to rewrite golden files.
func Compare(t *testing.T, p provider.Provider, update bool) {
	t.Helper()

	snap, err := Capture(p)
	if err != nil {
		t.Fatalf("capture provider %q: %v", p.Name(), err)
	}

	path := GoldenPath(t, snap.Name)
	if update {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}

	// #nosec G304 -- path is constructed by GoldenPath from a known base directory and a provider name.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden for %q: %v", snap.Name, err)
	}

	var expected Snapshot
	if err := json.Unmarshal(data, &expected); err != nil {
		t.Fatalf("unmarshal golden for %q: %v", snap.Name, err)
	}

	if !snapEqual(snap, expected) {
		got, _ := json.MarshalIndent(snap, "", "  ")
		t.Fatalf("provider %q contract drifted.\nexpected:\n%s\n\ngot:\n%s", snap.Name, strings.TrimSpace(string(data)), string(got))
	}
}

func snapEqual(a, b Snapshot) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}

// UpdateFlag is the conventional flag name for rewriting golden files.
const UpdateFlag = "update"

// Usage prints the -update flag usage.
func Usage() string {
	return fmt.Sprintf("Set -%s to regenerate golden files.", UpdateFlag)
}
