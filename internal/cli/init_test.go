package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/config"
)

func TestInitCommandDryRun(t *testing.T) {
	dir := t.TempDir()
	old := rootConfigPath
	rootConfigPath = filepath.Join(dir, ".muara", "config.yml")
	t.Cleanup(func() { rootConfigPath = old })

	cmd := newInitCommand()
	cmd.SetArgs([]string{"--dry-run", "--defaults"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if _, err := os.Stat(rootConfigPath); err == nil {
		t.Fatal("config file should not be created in dry-run mode")
	}
	out := buf.String()
	if !strings.Contains(out, "providers:") {
		t.Errorf("expected dry-run output to contain providers section, got:\n%s", out)
	}
}

func TestInitCommandForceOverwrite(t *testing.T) {
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
	cmd.SetArgs([]string{"--defaults", "--force"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	// #nosec G304
	data, err := os.ReadFile(rootConfigPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(data) == "existing" {
		t.Error("expected config to be overwritten")
	}
	if !strings.Contains(buf.String(), "created") {
		t.Errorf("expected created message, got:\n%s", buf.String())
	}
}

func TestRunWizardMultiProvider(t *testing.T) {
	inputs := []string{"1,3,5", "http://localhost:9001/hook", "debug"}
	idx := 0
	prompt := func(string) (string, error) {
		if idx >= len(inputs) {
			return "", fmt.Errorf("unexpected prompt")
		}
		v := inputs[idx]
		idx++
		return v, nil
	}

	var buf bytes.Buffer
	choices, webhook, level, err := runWizard(&buf, prompt)
	if err != nil {
		t.Fatalf("run wizard: %v", err)
	}
	if len(choices) != 3 {
		t.Fatalf("expected 3 providers, got %d", len(choices))
	}
	keys := make(map[string]bool)
	for _, c := range choices {
		keys[c.Key] = true
	}
	if !keys["fawry"] || !keys["billplz"] || !keys["ipay88"] {
		t.Errorf("unexpected provider selection: %+v", keys)
	}
	if webhook != "http://localhost:9001/hook" {
		t.Errorf("unexpected webhook: %q", webhook)
	}
	if level != "debug" {
		t.Errorf("unexpected log level: %q", level)
	}
}

func TestRunWizardDefaultProvider(t *testing.T) {
	inputs := []string{"", "", ""}
	idx := 0
	prompt := func(string) (string, error) {
		v := inputs[idx]
		idx++
		return v, nil
	}

	var buf bytes.Buffer
	choices, _, _, err := runWizard(&buf, prompt)
	if err != nil {
		t.Fatalf("run wizard: %v", err)
	}
	if len(choices) != 1 || choices[0].Key != "fawry" {
		t.Errorf("expected default fawry selection, got %+v", choices)
	}
}

func TestRunWizardProviderKeyInput(t *testing.T) {
	inputs := []string{"stripe,default", "", ""}
	idx := 0
	prompt := func(string) (string, error) {
		v := inputs[idx]
		idx++
		return v, nil
	}

	var buf bytes.Buffer
	choices, _, _, err := runWizard(&buf, prompt)
	if err != nil {
		t.Fatalf("run wizard: %v", err)
	}
	if len(choices) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(choices))
	}
	keys := make(map[string]bool)
	for _, c := range choices {
		keys[c.Key] = true
	}
	if !keys["stripe"] || !keys["default"] {
		t.Errorf("unexpected selection: %+v", keys)
	}
}

func TestGenerateMultiProviderConfig(t *testing.T) {
	fawry, _ := config.WizardChoiceByKey("fawry")
	stripe, _ := config.WizardChoiceByKey("stripe")
	cfg := config.GenerateWizardConfig([]config.WizardChoice{fawry, stripe}, "http://example.com/hook", "info")

	if !cfg.Providers["fawry"].Enabled {
		t.Error("expected fawry enabled")
	}
	if !cfg.Providers["stripe"].Enabled {
		t.Error("expected stripe enabled")
	}
	if cfg.Providers["billplz"].Enabled {
		t.Error("expected billplz disabled")
	}
}
