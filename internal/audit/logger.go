package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type loggerKey struct{}

// Logger records audit events.
type Logger interface {
	Log(_ context.Context, action, resourceType, resourceID, payload, result string)
}

// StoreLogger writes audit events to a Store.
type StoreLogger struct {
	Store       Store
	Actor       string
	Synchronous bool
}

// Log creates and persists an audit event.
func (l *StoreLogger) Log(_ context.Context, action, resourceType, resourceID, payload, result string) {
	ev := Event{
		ID:           uuid.Must(uuid.NewRandom()).String(),
		Timestamp:    time.Now(),
		Actor:        l.Actor,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Payload:      payload,
		Result:       result,
	}

	save := func() {
		if l.Store == nil {
			return
		}
		_ = l.Store.Save(ev)
	}

	if l.Synchronous {
		save()
		return
	}
	go save()
}

// NewContext returns a context carrying the provided audit logger.
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext returns the audit logger stored in ctx, or a no-op logger if none.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok && l != nil {
		return l
	}
	return &noopLogger{}
}

type noopLogger struct{}

func (n *noopLogger) Log(_ context.Context, _, _, _, _, _ string) {
}

// JSON returns a compact JSON representation of v, or an empty string on error.
func JSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
