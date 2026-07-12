package provider

import (
	"context"
	"net/http"
	"testing"
)

type stubProvider struct{ name string }

func (s *stubProvider) Name() string                 { return s.name }
func (s *stubProvider) Init(_ map[string]any) error  { return nil }
func (s *stubProvider) Routes() []Route              { return nil }
func (s *stubProvider) ChargeHandler() http.Handler  { return nil }
func (s *stubProvider) WebhookHandler() http.Handler { return nil }
func (s *stubProvider) PayloadBuilder() func(context.Context, Transaction) ([]byte, error) {
	return func(_ context.Context, _ Transaction) ([]byte, error) {
		return []byte("payload"), nil
	}
}
func (s *stubProvider) EscapeHandler() http.Handler { return nil }

func resetRegistry() {
	defaultRegistry.reset()
}

func TestRegisterAndGet(t *testing.T) {
	resetRegistry()

	t.Run("Given a provider named alpha is registered When the registry is queried for alpha Then the provider is returned", func(t *testing.T) {
		Register(&stubProvider{name: "alpha"})

		p, err := Get("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Name() != "alpha" {
			t.Errorf("name: want alpha, got %q", p.Name())
		}
	})

	t.Run("Given no provider named missing is registered When the registry is queried for missing Then an error is returned", func(t *testing.T) {
		_, err := Get("missing")
		if err == nil {
			t.Fatal("expected error for missing provider")
		}
	})
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	resetRegistry()

	t.Run("Given two providers with the same name are registered When the second registration occurs Then registration panics", func(t *testing.T) {
		Register(&stubProvider{name: "duplicate"})

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for duplicate registration")
			}
		}()
		Register(&stubProvider{name: "duplicate"})
	})
}

func TestRegisterPanicsOnNil(t *testing.T) {
	resetRegistry()

	t.Run("Given a nil provider is registered When registration occurs Then registration panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for nil provider")
			}
		}()
		Register(nil)
	})
}

func TestRegisterPanicsOnEmptyName(t *testing.T) {
	resetRegistry()

	t.Run("Given a provider with an empty name is registered When registration occurs Then registration panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for empty provider name")
			}
		}()
		Register(&stubProvider{name: ""})
	})
}

func TestNamesSorted(t *testing.T) {
	resetRegistry()

	t.Run("Given multiple providers are registered When Names is called Then it returns sorted names", func(t *testing.T) {
		Register(&stubProvider{name: "charlie"})
		Register(&stubProvider{name: "alpha"})
		Register(&stubProvider{name: "bravo"})

		names := Names()
		want := []string{"alpha", "bravo", "charlie"}
		if len(names) != len(want) {
			t.Fatalf("names: want %v, got %v", want, names)
		}
		for i, name := range want {
			if names[i] != name {
				t.Errorf("names[%d]: want %q, got %q", i, name, names[i])
			}
		}
	})
}

func TestPayloadBuilderWithTransaction(t *testing.T) {
	resetRegistry()

	t.Run("Given a provider.Transaction When a stub provider's PayloadBuilder is called Then it returns bytes without error", func(t *testing.T) {
		p := &stubProvider{name: "payload-test"}
		builder := p.PayloadBuilder()

		data, err := builder(context.Background(), Transaction{ID: "tx-1", Reference: "ref-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "payload" {
			t.Errorf("payload: want payload, got %q", string(data))
		}
	})
}

func TestNewRegistryIsIsolated(t *testing.T) {
	r := NewRegistry()
	r.Register(&stubProvider{name: "isolated"})

	p, err := r.Get("isolated")
	if err != nil {
		t.Fatalf("get from isolated registry: %v", err)
	}
	if p.Name() != "isolated" {
		t.Errorf("name: want isolated, got %q", p.Name())
	}

	if _, err := Default().Get("isolated"); err == nil {
		t.Error("isolated provider should not leak into default registry")
	}
}

func TestDefaultRegistryMatchesTopLevelFunctions(t *testing.T) {
	resetRegistry()

	Default().Register(&stubProvider{name: "default-test"})

	p, err := Get("default-test")
	if err != nil {
		t.Fatalf("get from top-level func: %v", err)
	}
	if p.Name() != "default-test" {
		t.Errorf("name: want default-test, got %q", p.Name())
	}
}
