package toyyibpay

import (
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/factory"
)

func init() {
	factory.MustRegister("toyyibpay", func(_ map[string]any) (provider.Provider, error) {
		return NewProvider(), nil
	})
}
