// Package factory provides a name-based registry for Go provider factories.
// Providers with runtime.type: go register a constructor here; the manifest
// controls whether that constructor is invoked.
package factory

import (
	"fmt"
	"sort"
	"sync"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

// Factory creates a new provider instance from its provider-specific config.
// The returned provider is not yet initialized; the caller must invoke Init.
type Factory func(cfg map[string]any) (provider.Provider, error)

// Registry holds registered provider factories.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry creates an empty factory registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
	}
}

// Register adds a factory for the given provider name. It panics if name is
// empty, factory is nil, or a factory for that name is already registered.
func (r *Registry) Register(name string, factory Factory) {
	if name == "" {
		panic("factory.Register: provider name is empty")
	}
	if factory == nil {
		panic("factory.Register: factory is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[name]; exists {
		panic("factory.Register: provider already registered: " + name)
	}
	r.factories[name] = factory
}

// Get returns the factory registered for the given provider name.
func (r *Registry) Get(name string) (Factory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[name]
	return f, ok
}

// Names returns all registered provider names sorted alphabetically.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

var defaultRegistry = NewRegistry()

// Default returns the package-level default factory registry.
func Default() *Registry {
	return defaultRegistry
}

// Register adds a factory to the default registry.
func Register(name string, factory Factory) {
	defaultRegistry.Register(name, factory)
}

// Get returns a factory from the default registry.
func Get(name string) (Factory, bool) {
	return defaultRegistry.Get(name)
}

// Names returns all names from the default registry.
func Names() []string {
	return defaultRegistry.Names()
}

// MustRegister registers a factory on the default registry and panics on
// duplicate registration. It is intended for use in provider package init
// functions.
func MustRegister(name string, factory Factory) {
	if err := validateRegistration(name, factory); err != nil {
		panic(err)
	}
	defaultRegistry.Register(name, factory)
}

func validateRegistration(name string, factory Factory) error {
	if name == "" {
		return fmt.Errorf("factory.Register: provider name is empty")
	}
	if factory == nil {
		return fmt.Errorf("factory.Register: factory is nil")
	}
	return nil
}
