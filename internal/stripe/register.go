package stripe

import (
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/factory"
)

func init() {
	factory.MustRegister("stripe", func(_ map[string]any) (provider.Provider, error) {
		return NewProvider(), nil
	})
}
