package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

// Given a valid plugin directory, When LoadBuiltin runs, Then it returns one LoadedPlugin.
func TestLoadBuiltinValidPlugin(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "test-gateway")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := []byte("schema_version: v1\nmetadata:\n  name: test-gateway\n  version: 1.0.0\n")
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), manifest, 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	plugins, err := LoadBuiltin(dir)
	if err != nil {
		t.Fatalf("load builtin: %v", err)
	}
	if len(plugins) != 1 {
		t.Fatalf("plugins: want 1, got %d", len(plugins))
	}
	if plugins[0].Name != "test-gateway" {
		t.Errorf("name: want test-gateway, got %q", plugins[0].Name)
	}
}

// Given a plugin directory named with uppercase letters, When LoadBuiltin runs, Then it returns an error.
func TestLoadBuiltinInvalidDirectoryName(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "BadName")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), []byte("schema_version: v1\n"), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	if _, err := LoadBuiltin(dir); err == nil {
		t.Fatal("expected error for invalid directory name, got nil")
	}
}

// Given a nonexistent plugin directory, When LoadBuiltin runs, Then it returns an empty list without error.
func TestLoadBuiltinMissingDirectory(t *testing.T) {
	plugins, err := LoadBuiltin(filepath.Join(t.TempDir(), "missing"))
	if err != nil {
		t.Fatalf("load builtin: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("plugins: want 0, got %d", len(plugins))
	}
}

// Given a plugin directory without gateway.yml, When LoadBuiltin runs, Then it returns an error.
func TestLoadBuiltinMissingManifest(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "empty-plugin")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if _, err := LoadBuiltin(dir); err == nil {
		t.Fatal("expected error for missing manifest, got nil")
	}
}

// Given a plugin directory with a malformed gateway.yml, When LoadBuiltin runs, Then it returns an error.
func TestLoadBuiltinMalformedManifest(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "bad-yaml")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), []byte("not: [ valid yaml"), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	if _, err := LoadBuiltin(dir); err == nil {
		t.Fatal("expected error for malformed manifest, got nil")
	}
}

func TestLoadLocalCustomDir(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "my-plugin")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), []byte("schema_version: v1\nmetadata:\n  name: my-plugin\n  version: 1.0.0\n"), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	plugins, err := LoadLocal(dir)
	if err != nil {
		t.Fatalf("load local: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "my-plugin" {
		t.Fatalf("expected one plugin named my-plugin, got %+v", plugins)
	}
}

func TestLoadLocalDefaultDirMissing(t *testing.T) {
	plugins, err := LoadLocal(filepath.Join(t.TempDir(), "nonexistent"))
	if err != nil {
		t.Fatalf("load local: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoadBuiltinDefaultDir(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(cwd) }()

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	pluginDir := filepath.Join(dir, "plugins", "default-plugin")
	if err := os.MkdirAll(pluginDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "gateway.yml"), []byte("schema_version: v1\nmetadata:\n  name: default-plugin\n  version: 1.0.0\n"), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	plugins, err := LoadBuiltin()
	if err != nil {
		t.Fatalf("load builtin: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "default-plugin" {
		t.Fatalf("expected one default plugin, got %+v", plugins)
	}
}
