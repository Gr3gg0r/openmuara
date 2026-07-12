package conform

import (
	"os"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/billplz"
	"github.com/Gr3gg0r/openmuara/internal/fawry"
	"github.com/Gr3gg0r/openmuara/internal/ipay88"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/senangpay"
	"github.com/Gr3gg0r/openmuara/internal/stripe"
	"github.com/Gr3gg0r/openmuara/internal/toyyibpay"
)

func TestProviderContracts(t *testing.T) {
	update := os.Getenv("UPDATE_GOLDEN") == "1"

	registry := provider.NewRegistry()
	fawry.RegisterWith(registry)
	senangpay.RegisterWith(registry)
	ipay88.RegisterWith(registry)
	billplz.RegisterWith(registry)
	toyyibpay.RegisterWith(registry)
	stripe.RegisterWith(registry)

	for _, name := range registry.Names() {
		t.Run(name, func(t *testing.T) {
			p, err := registry.Get(name)
			if err != nil {
				t.Fatalf("get provider %q: %v", name, err)
			}
			Compare(t, p, update)
		})
	}
}
