package config

import (
	"fmt"

	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/factory"
	"github.com/openmuara/openmuara/internal/provider/simple"
)

// ProviderFromGateway builds a provider implementation from a loaded gateway.yml
// manifest. Simple providers use the declarative runtime; go providers use the
// factory registry.
func ProviderFromGateway(lp *plugin.LoadedPlugin) (provider.Provider, error) {
	cfg := lp.Config
	if cfg.Runtime == nil || cfg.Runtime.Type == "" {
		cfg.Runtime = &plugin.Runtime{Type: "simple"}
	}

	if err := plugin.Validate(cfg); err != nil {
		return nil, fmt.Errorf("validate gateway.yml: %w", err)
	}

	switch cfg.Runtime.Type {
	case "simple":
		return simple.NewProvider(cfg), nil
	case "go":
		f, ok := factory.Get(lp.Name)
		if !ok {
			return nil, fmt.Errorf("no Go factory registered for provider %q", lp.Name)
		}
		p, err := f(nil)
		if err != nil {
			return nil, fmt.Errorf("create provider %q from factory: %w", lp.Name, err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unsupported runtime.type %q", cfg.Runtime.Type)
	}
}
