package plugin

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// Given a plugin and a registry, When Register is called, Then the plugin is retrievable.
func TestRegistryRegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	impl := NewBuiltinPlugin("test", "1.0.0", func(_ context.Context, r *Registry) error {
		return r.RegisterHandler("charge", func(_ Dependencies) (http.Handler, error) {
			return http.NotFoundHandler(), nil
		})
	})

	if err := reg.Register(context.Background(), "test", validConfig(), impl); err != nil {
		t.Fatalf("register: %v", err)
	}

	loaded, plugin, ok := reg.Get("test")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if loaded.Name != "test" {
		t.Errorf("name: want test, got %q", loaded.Name)
	}
	if plugin.Name() != "test" {
		t.Errorf("plugin name: want test, got %q", plugin.Name())
	}
}

// Given two plugins with the same name, When Register is called twice, Then it returns an error.
func TestRegistryDuplicatePlugin(t *testing.T) {
	reg := NewRegistry()
	impl := NewBuiltinPlugin("test", "1.0.0", func(_ context.Context, _ *Registry) error { return nil })

	if err := reg.Register(context.Background(), "test", validConfig(), impl); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if err := reg.Register(context.Background(), "test", validConfig(), impl); err == nil {
		t.Fatal("expected duplicate plugin error, got nil")
	}
}

// Given a handler factory registered for action charge, When Handler("charge") is called, Then it returns the factory.
func TestRegistryHandlerLookup(t *testing.T) {
	reg := NewRegistry()
	factory := func(_ Dependencies) (http.Handler, error) {
		return http.NotFoundHandler(), nil
	}
	if err := reg.RegisterHandler("charge", factory); err != nil {
		t.Fatalf("register handler: %v", err)
	}

	got, ok := reg.Handler("charge")
	if !ok {
		t.Fatal("expected handler to be found")
	}
	if got == nil {
		t.Fatal("expected non-nil factory")
	}
}

// Given a handler factory already registered for an action, When RegisterHandler is called again, Then it returns an error.
func TestRegistryDuplicateHandler(t *testing.T) {
	reg := NewRegistry()
	factory := func(_ Dependencies) (http.Handler, error) { return nil, nil }
	if err := reg.RegisterHandler("charge", factory); err != nil {
		t.Fatalf("first register handler: %v", err)
	}
	if err := reg.RegisterHandler("charge", factory); err == nil {
		t.Fatal("expected duplicate handler error, got nil")
	}
}

// Given a plugin whose Register returns an error, When Register is called, Then the error is propagated.
func TestRegistryPluginRegisterError(t *testing.T) {
	reg := NewRegistry()
	impl := NewBuiltinPlugin("test", "1.0.0", func(_ context.Context, _ *Registry) error {
		return errors.New("boom")
	})

	if err := reg.Register(context.Background(), "test", validConfig(), impl); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Given multiple registered plugins, When All is called, Then it returns all plugins.
func TestRegistryAll(t *testing.T) {
	reg := NewRegistry()
	for _, name := range []string{"a", "b"} {
		name := name
		impl := NewBuiltinPlugin(name, "1.0.0", func(_ context.Context, _ *Registry) error { return nil })
		if err := reg.Register(context.Background(), name, validConfig(), impl); err != nil {
			t.Fatalf("register %s: %v", name, err)
		}
	}

	all := reg.All()
	if len(all) != 2 {
		t.Fatalf("all: want 2, got %d", len(all))
	}
}

func TestMustRegisterHandlerAndDefaultHandler(t *testing.T) {
	MustRegisterHandler("unique-action-"+t.Name(), func(_ Dependencies) (http.Handler, error) {
		return http.NotFoundHandler(), nil
	})

	factory, ok := DefaultHandler("unique-action-" + t.Name())
	if !ok {
		t.Fatal("expected default handler to be found")
	}
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
}

func TestMustRegisterHandlerPanicsOnDuplicate(t *testing.T) {
	action := "duplicate-action-" + t.Name()
	MustRegisterHandler(action, func(_ Dependencies) (http.Handler, error) {
		return http.NotFoundHandler(), nil
	})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate handler registration")
		}
	}()
	MustRegisterHandler(action, func(_ Dependencies) (http.Handler, error) {
		return http.NotFoundHandler(), nil
	})
}
