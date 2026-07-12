package plugin

import (
	"context"
	"net/http"
)

// Plugin is the programmatic half of a provider definition.
// Each plugin directory may optionally provide a Go package that implements Plugin
// and registers handlers via the registry.
type Plugin interface {
	// Name returns the plugin identifier. It must match metadata.name in gateway.yml.
	Name() string
	// Version returns the plugin version. It should match metadata.version.
	Version() string
	// Register is called after the declarative config is loaded and validated.
	// Implementations use reg to register HTTP handlers and other runtime behavior.
	Register(ctx context.Context, reg *Registry) error
}

// HandlerFactory creates an http.Handler for a route action.
type HandlerFactory func(deps Dependencies) (http.Handler, error)

// Dependencies carries runtime dependencies available to plugin handlers.
type Dependencies struct {
	Config map[string]any
}
