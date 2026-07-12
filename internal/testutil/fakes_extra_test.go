package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/webhook"
)

func TestFakeProviderSetters(t *testing.T) {
	p := &FakeProvider{ProviderName: "fake"}

	p.SetStore(engine.NewMemoryStore())
	p.SetBaseURL("http://localhost:9000")
	p.SetDispatcher(&webhook.Dispatcher{})

	rec := httptest.NewRecorder()
	p.WebhookHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/webhook", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("webhook handler status: want 200, got %d", rec.Code)
	}

	payload, err := p.PayloadBuilder()(context.Background(), provider.Transaction{})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}
	if string(payload) != "{}" {
		t.Errorf("payload: want {}, got %q", payload)
	}
}
