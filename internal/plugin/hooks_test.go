package plugin

import (
	"context"
	"testing"
)

func TestBuiltinPluginNameAndVersion(t *testing.T) {
	p := NewBuiltinPlugin("test", "1.2.3", func(_ context.Context, _ *Registry) error { return nil })

	if p.Name() != "test" {
		t.Errorf("name: want test, got %q", p.Name())
	}
	if p.Version() != "1.2.3" {
		t.Errorf("version: want 1.2.3, got %q", p.Version())
	}
}

func TestBuiltinPluginRegisterNilFunction(t *testing.T) {
	p := NewBuiltinPlugin("nil", "1.0.0", nil)
	if err := p.Register(context.Background(), NewRegistry()); err == nil {
		t.Fatal("expected error for nil registration function")
	}
}
