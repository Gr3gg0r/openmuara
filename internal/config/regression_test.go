package config

import (
	"strings"
	"testing"
)

func TestLoadFromBytesValid(t *testing.T) {
	data := []byte("server:\n  host: 127.0.0.1\n  port: 8080\n")
	cfg, err := LoadFromBytes(data)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("port: want 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadFromBytesEmpty(t *testing.T) {
	cfg, err := LoadFromBytes(nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("host default: want 127.0.0.1, got %q", cfg.Server.Host)
	}
}

func TestLoadFromBytesInvalidYAML(t *testing.T) {
	_, err := LoadFromBytes([]byte("not: valid: yaml: ["))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadFromBytesEnvOverride(t *testing.T) {
	t.Setenv("MUARA_SERVER_PORT", "7777")
	cfg, err := LoadFromBytes([]byte("server:\n  host: 127.0.0.1\n  port: 8080\n"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Server.Port != 7777 {
		t.Errorf("port: want env override 7777, got %d", cfg.Server.Port)
	}
}

func TestConfigValidateAdminPortEqualToPort(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000, AdminPort: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when admin_port equals port")
	}
}

func TestConfigValidateAdminPortOutOfRange(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000, AdminPort: 70000},
		Persistence: PersistenceConfig{Type: "memory"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for admin_port out of range")
	}
}

func TestConfigValidateEmptyHost(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestConfigValidateUnsupportedPersistence(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "postgres"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unsupported persistence")
	}
}

func TestValidationErrorString(t *testing.T) {
	e := ValidationError{Field: "server.port", Message: "out of range", Hint: "use 1-65535", Line: 3, File: "config.yml"}
	if got := e.Error(); !strings.Contains(got, "config.yml:3") {
		t.Errorf("expected line info, got %q", got)
	}
}

func TestEnvVarName(t *testing.T) {
	if got := EnvVarName("stripe", "secret_key"); got != "MUARA_STRIPE_SECRET_KEY" {
		t.Errorf("want MUARA_STRIPE_SECRET_KEY, got %q", got)
	}
}

func TestWizardChoiceNextStep(t *testing.T) {
	choice := WizardChoice{Key: "default", SampleMethod: "POST", SampleRoute: "/default/charge"}
	step := choice.NextStep()
	if step.Method != "POST" || step.Route != "/default/charge" {
		t.Errorf("unexpected step: %+v", step)
	}
}

func TestRenderWizardConfigWithTargets(t *testing.T) {
	choice, _ := WizardChoiceByKey("default")
	cfg := GenerateWizardConfig([]WizardChoice{choice}, "http://example.com/hook", "info")
	cfg.Webhook.Targets = map[string]string{"stripe": "http://stripe.example.com/hook"}
	out := RenderWizardConfig(cfg)
	if !strings.Contains(string(out), "targets:") {
		t.Error("expected targets section")
	}
	if !strings.Contains(string(out), "stripe:") {
		t.Error("expected stripe target")
	}
}

func TestSortedStringMapKeys(t *testing.T) {
	m := map[string]string{"z": "1", "a": "2", "m": "3"}
	got := sortedStringMapKeys(m)
	want := []string{"a", "m", "z"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("sortedStringMapKeys: want %v, got %v", want, got)
		}
	}
}

func TestValidateWithDetailsEmptyHost(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Host: "", Port: 9000}, Persistence: PersistenceConfig{Type: "memory"}}
	errs := cfg.ValidateWithDetails("")
	var found bool
	for _, e := range errs {
		if e.Field == "server.host" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected server.host error, got %v", errs)
	}
}

func TestValidateWithDetailsUnsupportedPersistence(t *testing.T) {
	cfg := &Config{Server: ServerConfig{Host: "127.0.0.1", Port: 9000}, Persistence: PersistenceConfig{Type: "redis"}}
	errs := cfg.ValidateWithDetails("")
	var found bool
	for _, e := range errs {
		if e.Field == "persistence.type" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected persistence.type error, got %v", errs)
	}
}

func TestFieldLineMapMissingFile(t *testing.T) {
	m := fieldLineMap("/nonexistent/config.yml")
	if len(m) != 0 {
		t.Errorf("expected empty map for missing file, got %v", m)
	}
}
