package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/version"
	"github.com/spf13/cobra"
)

func TestVersionString(t *testing.T) {
	s := version.String()
	if s == "" {
		t.Fatal("version string must not be empty")
	}
	if !strings.Contains(s, version.Version) {
		t.Fatalf("expected version string to contain %q, got %q", version.Version, s)
	}
}

func TestDoctorCommand(t *testing.T) {
	cmd := newDoctorCommand()
	if cmd.Use != "doctor" {
		t.Fatalf("expected command use to be doctor, got %q", cmd.Use)
	}
}

func TestVersionCommandJSON(t *testing.T) {
	cmd := newVersionCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	jsonOutput = true
	defer func() { jsonOutput = false }()

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var out versionOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal JSON: %v\n%s", err, buf.String())
	}
	if out.Version != version.Version {
		t.Errorf("expected version %q, got %q", version.Version, out.Version)
	}
}

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"1.1.0", "1.0.0", true},
		{"1.0.1", "1.0.0", true},
		{"1.0.0", "1.0.0", false},
		{"0.9.0", "1.0.0", false},
		{"v1.1.0", "1.0.0", true},
		{"1.0.0", "v1.0.0", false},
		{"1.0.0-beta", "1.0.0", false},
	}

	for _, tc := range cases {
		got := isNewerVersion(tc.latest, tc.current)
		if got != tc.expected {
			t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tc.latest, tc.current, got, tc.expected)
		}
	}
}

func TestMaybeCheckUpdateSkipsDev(t *testing.T) {
	if !version.IsDev() {
		t.Skip("skipping dev-only test in release build")
	}
	latest, update := maybeCheckUpdate(testVersionCmd())
	if latest != "" || update {
		t.Errorf("expected no update check in dev, got latest=%q update=%v", latest, update)
	}
}

func TestMaybeCheckUpdateSkipsQuiet(t *testing.T) {
	quietOutput = true
	defer func() { quietOutput = false }()
	latest, update := maybeCheckUpdate(testVersionCmd())
	if latest != "" || update {
		t.Errorf("expected no update check when quiet, got latest=%q update=%v", latest, update)
	}
}

func testVersionCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	return cmd
}
