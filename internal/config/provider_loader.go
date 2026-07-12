// Package config loads and exposes muara configuration from YAML and environment variables.
package config

import (
	"fmt"
	"log/slog"

	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/factory"
)

// LoadedProvider pairs a provider with its runtime identifier.
type LoadedProvider struct {
	Name     string
	Provider provider.Provider
	Config   ProviderConfig
}

// LoadEnabledProviders returns all enabled providers from cfg that are registered
// in the default provider registry and whose Init call succeeds.
func LoadEnabledProviders(cfg *Config) ([]LoadedProvider, error) {
	return LoadEnabledProvidersWithRegistry(cfg, provider.Default())
}

// LoadEnabledProvidersWithRegistry returns all enabled providers from cfg that
// are registered in the given registry and whose Init call succeeds. Tests can
// pass an isolated registry to avoid package-level state.
func LoadEnabledProvidersWithRegistry(cfg *Config, registry *provider.Registry) ([]LoadedProvider, error) {
	return LoadEnabledProvidersWithFallback(cfg, registry, plugin.LoadBuiltin)
}

// LoadEnabledProvidersWithFallback loads enabled providers from gateway.yml when
// available. Providers without a manifest fall back to the Go factory registry
// with a deprecation warning (D007). The pluginLoader discovers plugin manifests.
func LoadEnabledProvidersWithFallback(
	cfg *Config,
	registry *provider.Registry,
	pluginLoader func(...string) ([]*plugin.LoadedPlugin, error),
) ([]LoadedProvider, error) {
	if cfg.Providers == nil {
		return nil, nil
	}

	plugins, err := pluginLoader("plugins", "../plugins", "../../plugins")
	if err != nil {
		return nil, fmt.Errorf("load plugin manifests: %w", err)
	}
	byName := make(map[string]*plugin.LoadedPlugin, len(plugins))
	for _, p := range plugins {
		byName[p.Name] = p
	}

	var enabled []LoadedProvider
	for name, pc := range cfg.Providers {
		if !pc.Enabled {
			continue
		}

		var p provider.Provider
		if lp, ok := byName[name]; ok {
			p, err = ProviderFromGateway(lp)
			if err != nil {
				return nil, fmt.Errorf("load provider %q from gateway.yml: %w", name, err)
			}
			slog.Info("loaded provider from gateway.yml", "provider", name, "runtime", lp.Config.Runtime.Type)
		} else {
			// D007 soft landing: configured provider has no manifest. Warn and
			// still load from the factory registry when possible.
			slog.Warn(fmt.Sprintf("provider %q is configured but has no gateway.yml manifest; auto-loading built-in providers is deprecated. See docs/migration/provider-manifests.md", name))
			if f, ok := factory.Get(name); ok {
				p, err = f(nil)
				if err != nil {
					return nil, fmt.Errorf("load provider %q from factory: %w", name, err)
				}
			} else {
				p, err = registry.Get(name)
				if err != nil {
					return nil, fmt.Errorf("unknown provider %q", name)
				}
			}
		}

		// Clone the config map so providers cannot mutate the shared config.
		configCopy := make(map[string]any, len(pc.Config))
		for k, v := range pc.Config {
			configCopy[k] = v
		}

		if err := p.Init(configCopy); err != nil {
			return nil, fmt.Errorf("init provider %q: %w", name, err)
		}

		enabled = append(enabled, LoadedProvider{
			Name:     name,
			Provider: p,
			Config:   pc,
		})
	}

	return enabled, nil
}
