package webhook

import (
	"context"
	"fmt"
	"log/slog"
)

// Relay forwards a single webhook to multiple destinations.
type Relay struct {
	Senders []Sender
}

// NewRelay creates a relay from a list of senders.
func NewRelay(senders ...Sender) *Relay {
	return &Relay{Senders: senders}
}

// Dispatch fans out the webhook to all configured senders and returns the
// first successful attempt. Each sender still records its attempts in its
// own store; callers can inspect all attempts via the aggregate store.
func (r *Relay) Dispatch(ctx context.Context, ref string, status PaymentStatus) (*Attempt, error) {
	if len(r.Senders) == 0 {
		return nil, fmt.Errorf("no relay senders configured")
	}

	var first *Attempt
	for _, sender := range r.Senders {
		attempt, err := sender.Dispatch(ctx, ref, status)
		if err != nil {
			slog.Warn("relay dispatch failed", "ref", ref, "error", err)
			continue
		}
		if first == nil {
			first = attempt
		}
	}

	if first == nil {
		return nil, fmt.Errorf("all relay senders failed for ref %q", ref)
	}
	return first, nil
}

// compile-time check that Relay implements Sender.
var _ Sender = (*Relay)(nil)
