package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider/defaultplugin"
)

func TestRouterConcurrentDefaultProvider(_ *testing.T) {
	cfg := RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
		TransactionStore: engine.NewMemoryStore(),
		CORS:             config.CORSConfig{AllowedOrigins: []string{"*"}},
		CSRF:             config.CSRFConfig{Enabled: false},
	}
	// Ensure provider is registered (init in defaultplugin registers it globally).
	_ = defaultplugin.NewProvider()
	handler := NewRouter(cfg)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()

			body, _ := json.Marshal(map[string]any{
				"amount":    10.0,
				"currency":  "MYR",
				"reference": string(rune('a' + i%26)),
				"provider":  "default",
			})
			req := httptest.NewRequest(http.MethodPost, "/default/charge", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			req2 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			rec2 := httptest.NewRecorder()
			handler.ServeHTTP(rec2, req2)
		}()
	}

	wg.Wait()
}
