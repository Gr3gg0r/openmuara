package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/testutil"
)

func TestAuditListCommand(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "log:\n  level: error\npersistence:\n  type: memory\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	oldPath := rootConfigPath
	rootConfigPath = dir + "/config.yml"
	defer func() { rootConfigPath = oldPath }()

	cmd := newAuditListCommand()
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

func TestAuditListWithSinceFlag(t *testing.T) {
	dir := testutil.TempWorkspace(t)
	cfg := "log:\n  level: error\npersistence:\n  type: memory\n"
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))

	oldPath := rootConfigPath
	rootConfigPath = dir + "/config.yml"
	defer func() { rootConfigPath = oldPath }()

	cmd := newAuditListCommand()
	cmd.SetArgs([]string{"--since", "2026-01-01T00:00:00Z", "--limit", "5"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array output, got %q", buf.String())
	}
}
