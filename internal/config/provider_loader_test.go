package config

import (
	"testing"

	_ "github.com/openmuara/openmuara/internal/fawry" // register fawry factory for manifest loading
	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	_ "github.com/openmuara/openmuara/internal/provider/defaultplugin"
)

func TestLoadEnabledProviders(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"default": {Enabled: true, Config: map[string]any{}},
		},
	}

	loaded, err := LoadEnabledProviders(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(loaded))
	}
	if loaded[0].Name != "default" {
		t.Errorf("name: want default, got %q", loaded[0].Name)
	}
}

func TestLoadEnabledProvidersSkipsDisabled(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"default": {Enabled: false, Config: map[string]any{}},
		},
	}

	loaded, err := LoadEnabledProviders(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("expected 0 providers, got %d", len(loaded))
	}
}

func TestLoadEnabledProvidersUnknown(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"unknown-provider": {Enabled: true, Config: map[string]any{}},
		},
	}

	if _, err := LoadEnabledProviders(cfg); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestLoadEnabledProvidersDoesNotMutateConfig(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"default": {Enabled: true, Config: map[string]any{"key": "value"}},
		},
	}

	_, err := LoadEnabledProviders(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Providers["default"].Config["key"] != "value" {
		t.Error("provider config was mutated")
	}
}

func TestProviderRegistryAvailable(t *testing.T) {
	_, err := provider.Get("default")
	if err != nil {
		t.Fatalf("default provider should be registered: %v", err)
	}
}

func TestLoadEnabledProvidersFallbackToGatewayYAML(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			// #nosec G101 -- test fixture dummy credentials, not real secrets
			"fawry": {Enabled: true, Config: map[string]any{
				"merchant_code":         "muara-merchant-code",
				"merchant_security_key": "muara-fawry-secret",
				"webhook_secret":        "muara-webhook-secret",
			}},
		},
	}

	// Use an empty registry so fawry must be loaded from gateway.yml.
	registry := provider.NewRegistry()
	loader := func(_ ...string) ([]*plugin.LoadedPlugin, error) {
		return plugin.LoadBuiltin("../../plugins")
	}
	loaded, err := LoadEnabledProvidersWithFallback(cfg, registry, loader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(loaded))
	}
	if loaded[0].Name != "fawry" {
		t.Errorf("name: want fawry, got %q", loaded[0].Name)
	}
}

func TestLoadEnabledProvidersUnknownStillErrors(t *testing.T) {
	cfg := &Config{
		Providers: map[string]ProviderConfig{
			"unknown-provider": {Enabled: true, Config: map[string]any{}},
		},
	}

	registry := provider.NewRegistry()
	emptyLoader := func(_ ...string) ([]*plugin.LoadedPlugin, error) {
		return nil, nil
	}
	if _, err := LoadEnabledProvidersWithFallback(cfg, registry, emptyLoader); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}
