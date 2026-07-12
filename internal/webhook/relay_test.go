package webhook

import (
	"context"
	"testing"
)

func TestRelayDispatchesToAllSenders(t *testing.T) {
	called := 0
	sender := SenderFunc(func(_ context.Context, _ string, _ PaymentStatus) (*Attempt, error) {
		called++
		return &Attempt{ID: "attempt-1"}, nil
	})

	relay := NewRelay(sender, sender)
	_, err := relay.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if called != 2 {
		t.Errorf("expected both senders to be called, got %d", called)
	}
}

func TestRelayReturnsErrorWhenEmpty(t *testing.T) {
	relay := NewRelay()
	if _, err := relay.Dispatch(context.Background(), "ref-1", PaymentStatusPaid); err == nil {
		t.Fatal("expected error for empty relay")
	}
}

// SenderFunc is a test helper that implements Sender.
type SenderFunc func(context.Context, string, PaymentStatus) (*Attempt, error)

func (f SenderFunc) Dispatch(ctx context.Context, ref string, status PaymentStatus) (*Attempt, error) {
	return f(ctx, ref, status)
}
