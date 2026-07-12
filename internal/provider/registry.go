// Package provider provides a name-based registry for payment gateway plugins.
package provider

import (
	"errors"
	"sort"
	"sync"
)

// Registry holds registered providers and allows isolated instances for tests.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry. It panics if the provider is nil,
// the name is empty, or a provider with the same name is already registered.
func (r *Registry) Register(p Provider) {
	if p == nil {
		panic("provider.Register: provider is nil")
	}
	name := p.Name()
	if name == "" {
		panic("provider.Register: provider name is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.providers[name]; exists {
		panic("provider.Register: provider already registered: " + name)
	}
	r.providers[name] = p
}

// Get returns the provider registered with the given name.
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.providers[name]
	if !ok {
		return nil, errors.New("provider not found: " + name)
	}
	return p, nil
}

// Names returns all registered provider names sorted alphabetically.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// reset clears the registry. It is intended for tests only.
func (r *Registry) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]Provider)
}

// defaultRegistry is the package-level registry used by the top-level helpers.
var defaultRegistry = NewRegistry()

// Default returns the package-level default registry.
func Default() *Registry {
	return defaultRegistry
}

// Register adds a provider to the default registry.
func Register(p Provider) {
	defaultRegistry.Register(p)
}

// Get returns a provider from the default registry.
func Get(name string) (Provider, error) {
	return defaultRegistry.Get(name)
}

// Names returns all names from the default registry.
func Names() []string {
	return defaultRegistry.Names()
}
