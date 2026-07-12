package billplz

import (
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/provider/factory"
)

func init() {
	factory.MustRegister("billplz", func(_ map[string]any) (provider.Provider, error) {
		return NewProvider(), nil
	})
}
