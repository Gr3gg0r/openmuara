package webhook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestDispatcherRetriesOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 2, func(context.Context, provider.Transaction) ([]byte, error) {
		return []byte(`{}`), nil
	}, nil)

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	attempt, err := d.Store.Get("ref-1")
	if err != nil {
		t.Fatalf("get attempt: %v", err)
	}

	// Wait for retries to finish.
	deadline := time.Now().Add(5 * time.Second)
	for attempt.Status == AttemptStatusPending && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		attempt, _ = d.Store.Get("ref-1")
	}

	if attempt.Status != AttemptStatusFailed {
		t.Fatalf("expected failed status, got %s", attempt.Status)
	}
	if attempt.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempt.Attempts)
	}
}

func TestDispatcherRetriesAndSucceeds(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 3, func(context.Context, provider.Transaction) ([]byte, error) {
		return []byte(`{}`), nil
	}, nil)

	_, _ = d.Dispatch(context.Background(), "ref-2", PaymentStatusPaid)

	deadline := time.Now().Add(5 * time.Second)
	attempt, _ := d.Store.Get("ref-2")
	for attempt.Status == AttemptStatusPending && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		attempt, _ = d.Store.Get("ref-2")
	}

	if attempt.Status != AttemptStatusDelivered {
		t.Fatalf("expected delivered status, got %s", attempt.Status)
	}
}

func TestDispatcherHonoursTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 0, func(context.Context, provider.Transaction) ([]byte, error) {
		return []byte(`{}`), nil
	}, nil)
	d.Worker.Client = &http.Client{Timeout: 100 * time.Millisecond}

	_, _ = d.Dispatch(context.Background(), "ref-3", PaymentStatusPaid)

	deadline := time.Now().Add(2 * time.Second)
	attempt, _ := d.Store.Get("ref-3")
	for attempt.Status == AttemptStatusPending && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		attempt, _ = d.Store.Get("ref-3")
	}

	if attempt.Status != AttemptStatusFailed {
		t.Fatalf("expected failed status after timeout, got %s", attempt.Status)
	}
}
