package plugin

import (
	"context"
	"fmt"
)

// BuiltinPlugin is a minimal Plugin implementation for plugins that only need registration.
type BuiltinPlugin struct {
	name    string
	version string
	regFn   func(ctx context.Context, reg *Registry) error
}

// NewBuiltinPlugin creates a BuiltinPlugin.
func NewBuiltinPlugin(name, version string, regFn func(ctx context.Context, reg *Registry) error) *BuiltinPlugin {
	return &BuiltinPlugin{name: name, version: version, regFn: regFn}
}

// Name returns the plugin name.
func (p *BuiltinPlugin) Name() string { return p.name }

// Version returns the plugin version.
func (p *BuiltinPlugin) Version() string { return p.version }

// Register delegates to the provided registration function.
func (p *BuiltinPlugin) Register(ctx context.Context, reg *Registry) error {
	if p.regFn == nil {
		return fmt.Errorf("plugin %q has no registration function", p.name)
	}
	return p.regFn(ctx, reg)
}
