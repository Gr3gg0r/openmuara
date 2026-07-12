package plugin

import (
	"context"
	"fmt"
)

// Registry holds validated plugins and their runtime handlers.
type Registry struct {
	plugins  map[string]*entry
	handlers map[string]HandlerFactory
}

type entry struct {
	config GatewayConfig
	impl   Plugin
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins:  make(map[string]*entry),
		handlers: make(map[string]HandlerFactory),
	}
}

// Register adds a validated plugin and its implementation to the registry.
func (r *Registry) Register(ctx context.Context, name string, cfg GatewayConfig, impl Plugin) error {
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %q already registered", name)
	}
	if err := impl.Register(ctx, r); err != nil {
		return fmt.Errorf("register plugin %q: %w", name, err)
	}
	r.plugins[name] = &entry{config: cfg, impl: impl}
	return nil
}

// Get returns a loaded plugin by name.
func (r *Registry) Get(name string) (*LoadedPlugin, Plugin, bool) {
	e, ok := r.plugins[name]
	if !ok {
		return nil, nil, false
	}
	return &LoadedPlugin{Name: name, Config: e.config}, e.impl, true
}

// All returns all registered plugins.
func (r *Registry) All() []*LoadedPlugin {
	out := make([]*LoadedPlugin, 0, len(r.plugins))
	for name, e := range r.plugins {
		out = append(out, &LoadedPlugin{Name: name, Config: e.config})
	}
	return out
}

// RegisterHandler registers a factory for a route action.
func (r *Registry) RegisterHandler(action string, factory HandlerFactory) error {
	if _, exists := r.handlers[action]; exists {
		return fmt.Errorf("handler for action %q already registered", action)
	}
	r.handlers[action] = factory
	return nil
}

// Handler returns the factory registered for an action.
func (r *Registry) Handler(action string) (HandlerFactory, bool) {
	f, ok := r.handlers[action]
	return f, ok
}

// defaultRegistry is the package-level registry used by MustRegisterHandler.
var defaultRegistry = NewRegistry()

// MustRegisterHandler registers a factory on the package-level default registry
// and panics if the action is already registered. It is intended for use in
// built-in plugin init functions.
func MustRegisterHandler(action string, factory HandlerFactory) {
	if err := defaultRegistry.RegisterHandler(action, factory); err != nil {
		panic(err)
	}
}

// DefaultHandler returns the factory registered on the package-level registry.
func DefaultHandler(action string) (HandlerFactory, bool) {
	return defaultRegistry.Handler(action)
}
