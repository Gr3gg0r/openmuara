package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestFakeDispatcherDispatch(t *testing.T) {
	f := &FakeDispatcher{}
	attempt, err := f.Dispatch(context.Background(), "ref-1", webhook.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if attempt.Ref != "ref-1" {
		t.Errorf("attempt ref: want ref-1, got %q", attempt.Ref)
	}
	if len(f.Calls) != 1 {
		t.Fatalf("calls: want 1, got %d", len(f.Calls))
	}
	if f.Calls[0].Ref != "ref-1" || f.Calls[0].Status != webhook.PaymentStatusPaid {
		t.Errorf("call mismatch: got %+v", f.Calls[0])
	}
}

func TestFakeProviderMethods(t *testing.T) {
	p := &FakeProvider{ProviderName: "fake"}

	if got := p.Name(); got != "fake" {
		t.Errorf("name: want fake, got %q", got)
	}
	if err := p.Init(nil); err != nil {
		t.Errorf("init: %v", err)
	}

	p.InitErr = context.Canceled
	if err := p.Init(nil); err == nil {
		t.Error("expected init error")
	}

	if p.Routes() != nil {
		t.Error("expected nil routes")
	}

	p.RoutesFunc = func() []provider.Route {
		return []provider.Route{{Method: http.MethodGet, Path: "/fake"}}
	}
	if len(p.Routes()) != 1 {
		t.Fatalf("routes: want 1, got %d", len(p.Routes()))
	}

	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/fake", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("charge handler status: want 200, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	p.WebhookHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/webhook", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("webhook handler status: want 200, got %d", rec.Code)
	}

	if p.EscapeHandler() != nil {
		t.Error("expected nil escape handler")
	}

	payload, err := p.PayloadBuilder()(context.Background(), provider.Transaction{Reference: "r", Status: "paid"})
	if err != nil {
		t.Fatalf("payload builder: %v", err)
	}
	if string(payload) != "{}" {
		t.Errorf("payload: want {}, got %s", payload)
	}

	p.SetStore(engine.NewMemoryStore())
	p.SetBaseURL("http://localhost")
	p.SetDispatcher(nil)
}
