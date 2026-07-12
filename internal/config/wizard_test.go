package config

import (
	"strings"
	"testing"
)

func TestWizardTemplatesIncludesRecommendedProvider(t *testing.T) {
	templates := WizardTemplates()
	var found bool
	for _, c := range templates.Choices {
		if c.IsRecommended {
			found = true
			if c.Key != "fawry" {
				t.Errorf("expected fawry as recommended, got %q", c.Key)
			}
		}
	}
	if !found {
		t.Error("expected at least one recommended provider")
	}
}

func TestWizardChoiceByKey(t *testing.T) {
	choice, ok := WizardChoiceByKey("stripe")
	if !ok {
		t.Fatal("expected stripe choice to exist")
	}
	if choice.Key != "stripe" {
		t.Errorf("expected key stripe, got %q", choice.Key)
	}
	if choice.ProviderConfig.Enabled {
		t.Error("expected stripe choice to be disabled by default in wizard")
	}
}

func TestGenerateWizardConfig(t *testing.T) {
	choice, _ := WizardChoiceByKey("fawry")
	cfg := GenerateWizardConfig([]WizardChoice{choice}, "http://localhost:9001/webhook", "debug")

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("host: want 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("log level: want debug, got %s", cfg.Log.Level)
	}
	if cfg.Webhook.URL != "http://localhost:9001/webhook" {
		t.Errorf("webhook url mismatch: %s", cfg.Webhook.URL)
	}
	pc, ok := cfg.Providers["fawry"]
	if !ok {
		t.Fatal("expected fawry provider config")
	}
	if !pc.Enabled {
		t.Error("expected fawry to be enabled")
	}
	if pc.Config["version"] != "v1" {
		t.Errorf("expected fawry version v1, got %v", pc.Config["version"])
	}
}

func TestGenerateWizardConfigMultiProvider(t *testing.T) {
	fawry, _ := WizardChoiceByKey("fawry")
	stripe, _ := WizardChoiceByKey("stripe")
	cfg := GenerateWizardConfig([]WizardChoice{fawry, stripe}, "", "info")

	if !cfg.Providers["fawry"].Enabled {
		t.Error("expected fawry to be enabled")
	}
	if !cfg.Providers["stripe"].Enabled {
		t.Error("expected stripe to be enabled")
	}
	if cfg.Providers["billplz"].Enabled {
		t.Error("expected billplz to be disabled")
	}
}

func TestRenderWizardConfig(t *testing.T) {
	choice, _ := WizardChoiceByKey("billplz")
	cfg := GenerateWizardConfig([]WizardChoice{choice}, "", "info")
	out := RenderWizardConfig(cfg)

	if len(out) == 0 {
		t.Fatal("expected non-empty config output")
	}
	s := string(out)
	if !strings.Contains(s, "providers:") {
		t.Error("expected providers section")
	}
	if !strings.Contains(s, "billplz:") {
		t.Error("expected billplz provider")
	}
	if !strings.Contains(s, "api_key:") {
		t.Error("expected api_key config")
	}
	if !strings.Contains(s, "webhook:") {
		t.Error("expected webhook section even when URL is empty")
	}
}

func TestRenderWizardConfigIncludesWebhook(t *testing.T) {
	choice, _ := WizardChoiceByKey("default")
	cfg := GenerateWizardConfig([]WizardChoice{choice}, "http://localhost:9001/hook", "info")
	out := RenderWizardConfig(cfg)

	if !strings.Contains(string(out), "webhook:") {
		t.Error("expected webhook section when URL is set")
	}
}
