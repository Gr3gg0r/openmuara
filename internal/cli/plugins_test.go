package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestPlugin(t *testing.T, root string) {
	t.Helper()
	pluginDir := filepath.Join(root, "plugins", "testplugin")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("create plugin dir: %v", err)
	}
	content := `schema_version: v1
metadata:
  name: testplugin
  version: 1.0.0
  description: A test plugin
routes: []
`
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), []byte(content), 0o600); err != nil {
		t.Fatalf("write gateway.yml: %v", err)
	}
}

func TestPluginsListCommand(t *testing.T) {
	root := t.TempDir()
	writeTestPlugin(t, root)
	t.Chdir(root)

	cmd := newPluginsListCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "testplugin") {
		t.Errorf("output missing plugin name: %q", out)
	}
	if !strings.Contains(out, "A test plugin") {
		t.Errorf("output missing description: %q", out)
	}
}

func TestPluginsValidateAllCommand(t *testing.T) {
	root := t.TempDir()
	writeTestPlugin(t, root)
	t.Chdir(root)

	cmd := newPluginsValidateCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "OK testplugin") {
		t.Errorf("output missing OK: %q", buf.String())
	}
}

func TestPluginsValidateSinglePathCommand(t *testing.T) {
	root := t.TempDir()
	writeTestPlugin(t, root)
	t.Chdir(root)

	cmd := newPluginsValidateCommand()
	cmd.SetArgs([]string{"plugins"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "OK testplugin") {
		t.Errorf("output missing OK: %q", buf.String())
	}
}

func TestPluginsCommandWiring(t *testing.T) {
	root := t.TempDir()
	writeTestPlugin(t, root)
	t.Chdir(root)

	cmd := newPluginsCommand()
	cmd.SetArgs([]string{"list"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "testplugin") {
		t.Errorf("output missing plugin name: %q", buf.String())
	}
}
