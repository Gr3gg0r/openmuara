package webhook

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/errcode"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/prometheus/client_golang/prometheus"
)

// Sender is the minimal webhook dispatch contract.
type Sender interface {
	Dispatch(ctx context.Context, ref string, status PaymentStatus) (*Attempt, error)
}

// SignatureVerifier is an optional contract for providers that can verify the
// signature of an outgoing webhook payload and its headers.
type SignatureVerifier interface {
	VerifyWebhookSignature(payload []byte, headers map[string]string) (bool, error)
}

// Dispatcher builds and sends provider-style webhooks.
type Dispatcher struct {
	URL           string
	Secret        string
	MaxRetries    int
	ProviderName  string
	Builder       func(context.Context, provider.Transaction) ([]byte, error)
	HeaderBuilder func(context.Context, provider.Transaction) (map[string]string, error)
	EventTypeFor  func(ref string, status PaymentStatus) string
	EnabledEvents []string
	Worker        *DeliveryWorker
	Store         AttemptStore
	Ledger        engine.TransactionStore
	AuditLogger   audit.Logger
	Verifier      SignatureVerifier
}

var (
	webhookAttemptsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "openmuara_webhook_attempts_total",
		Help: "Total webhook delivery attempts by provider and final status.",
	}, []string{"provider", "status"})
)

func init() {
	prometheus.MustRegister(webhookAttemptsTotal)
}

// compile-time check that Dispatcher implements Sender.
var _ Sender = (*Dispatcher)(nil)

// NewDispatcher creates a dispatcher for the given version and config.
//
// Deprecated: Use NewDispatcherFromProvider to build a dispatcher from a
// registered provider instead of hardcoding a payload version.
func NewDispatcher(version PayloadVersion, url, secret string, maxRetries int, txStore engine.TransactionStore) (*Dispatcher, error) {
	legacyBuilder, err := PayloadFor(version, secret, txStore)
	if err != nil {
		return nil, fmt.Errorf("payload builder: %w", err)
	}

	builder := func(_ context.Context, tx provider.Transaction) ([]byte, error) {
		return legacyBuilder.Build(tx.Reference, PaymentStatus(tx.Status))
	}

	return NewDispatcherFromBuilder(url, maxRetries, builder, nil), nil
}

// NewDispatcherFromProvider creates a dispatcher that uses the provider's payload builder.
func NewDispatcherFromProvider(url string, maxRetries int, p provider.Provider) *Dispatcher {
	var headerBuilder func(context.Context, provider.Transaction) (map[string]string, error)
	if hp, ok := p.(interface {
		PayloadHeaders(context.Context, provider.Transaction) (map[string]string, error)
	}); ok {
		headerBuilder = hp.PayloadHeaders
	}
	d := NewDispatcherFromBuilder(url, maxRetries, p.PayloadBuilder(), headerBuilder)
	d.ProviderName = p.Name()

	if sv, ok := p.(SignatureVerifier); ok {
		d.Verifier = sv
	}

	if etp, ok := p.(interface {
		PayloadEventType(ref, status string) string
	}); ok {
		d.EventTypeFor = func(ref string, status PaymentStatus) string {
			return etp.PayloadEventType(ref, string(status))
		}
	}
	return d
}

// NewDispatcherFromBuilder creates a dispatcher with an explicit builder function.
func NewDispatcherFromBuilder(url string, maxRetries int, builder func(context.Context, provider.Transaction) ([]byte, error), headerBuilder func(context.Context, provider.Transaction) (map[string]string, error)) *Dispatcher {
	if maxRetries < 0 {
		maxRetries = 0
	}

	return &Dispatcher{
		URL:           url,
		MaxRetries:    maxRetries,
		Builder:       builder,
		HeaderBuilder: headerBuilder,
		Worker:        NewDeliveryWorker(),
		Store:         NewMemoryStore(),
	}
}

// Dispatch builds and delivers a webhook for the given reference and status.
func (d *Dispatcher) Dispatch(ctx context.Context, ref string, status PaymentStatus) (*Attempt, error) {
	if d.URL == "" {
		return nil, errcode.New(errcode.EWebhookURLMissing, "webhook url is not configured")
	}

	if d.EventTypeFor != nil && len(d.EnabledEvents) > 0 {
		eventType := d.EventTypeFor(ref, status)
		if !eventEnabled(eventType, d.EnabledEvents) {
			return nil, errcode.New(errcode.EProviderDisabled, fmt.Sprintf("event type %q is not enabled", eventType))
		}
	}

	payload, err := d.Builder(ctx, provider.Transaction{Reference: ref, Status: string(status)})
	if err != nil {
		return nil, errcode.Wrap(errcode.EWebhookBuildFailed, "build payload", err)
	}

	var headers map[string]string
	if d.HeaderBuilder != nil {
		headers, err = d.HeaderBuilder(ctx, provider.Transaction{Reference: ref, Status: string(status)})
		if err != nil {
			return nil, errcode.Wrap(errcode.EWebhookBuildFailed, "build headers", err)
		}
	}

	traceID := httputil.TraceIDFromContext(ctx)
	if traceID == "" && d.Ledger != nil {
		if tx, ok, _ := d.Ledger.GetByReference(ref); ok && tx.TraceID != "" {
			traceID = tx.TraceID
		}
	}

	attempt := &Attempt{
		ID:           uuid.Must(uuid.NewRandom()).String(),
		Ref:          ref,
		ProviderName: d.ProviderName,
		URL:          d.URL,
		Status:       AttemptStatusPending,
		Payload:      payload,
		Headers:      headers,
		TraceID:      traceID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if d.Verifier != nil {
		valid, err := d.Verifier.VerifyWebhookSignature(payload, headers)
		if err == nil {
			attempt.SignatureValid = &valid
		}
	}

	if err := d.Store.Save(attempt); err != nil {
		return nil, errcode.Wrap(errcode.EInternal, "save attempt", err)
	}

	result := snapshot(attempt)
	// Detach from the caller context so async delivery survives request
	// cancellation, but keep the trace ID for the outgoing request.
	deliver(context.Background(), d, attempt)

	return result, nil
}

// traceIDFromAttempt returns the attempt's trace ID, falling back to a new
// UUID only when empty so tests and legacy callers still get a header.
func traceIDFromAttempt(a *Attempt) string {
	if a.TraceID != "" {
		return a.TraceID
	}
	return uuid.Must(uuid.NewRandom()).String()
}

// Replay re-sends the stored payload for a reference.
func (d *Dispatcher) Replay(_ context.Context, ref string) (*Attempt, error) {
	attempt, err := d.Store.Get(ref)
	if err != nil {
		return nil, errcode.Wrap(errcode.EInternal, "get attempt", err)
	}
	if attempt == nil {
		return nil, errcode.Wrap(errcode.EWebhookReplayNotFound, "no webhook found for ref", fmt.Errorf("ref=%s", ref))
	}

	newAttempt := &Attempt{
		ID:             uuid.Must(uuid.NewRandom()).String(),
		Ref:            attempt.Ref,
		URL:            attempt.URL,
		Status:         AttemptStatusPending,
		Payload:        attempt.Payload,
		Headers:        attempt.Headers,
		TraceID:        attempt.TraceID,
		SignatureValid: attempt.SignatureValid,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := d.Store.Save(newAttempt); err != nil {
		return nil, errcode.Wrap(errcode.EInternal, "save replay attempt", err)
	}

	result := snapshot(newAttempt)
	deliver(context.Background(), d, newAttempt)

	return result, nil
}

func deliver(ctx context.Context, d *Dispatcher, attempt *Attempt) {
	// #nosec G118 -- intentional background delivery detached from request context
	go func(ctx context.Context) {
		var history []AttemptHistory
		var lastErr string
		for i := 0; i <= d.MaxRetries; i++ {
			status, _, err := d.Worker.Deliver(ctx, attempt.URL, attempt.Payload, attempt.Headers, traceIDFromAttempt(attempt))
			attempt.Attempts++

			if err == nil && IsSuccess(status) {
				attempt.Status = AttemptStatusDelivered
				attempt.LastError = ""
				history = append(history, AttemptHistory{Time: time.Now(), Status: status, Error: ""})
				break
			}

			if err != nil {
				lastErr = err.Error()
			} else {
				lastErr = fmt.Sprintf("non-2xx response: %d", status)
			}
			attempt.LastError = lastErr
			history = append(history, AttemptHistory{Time: time.Now(), Status: status, Error: lastErr})

			if i < d.MaxRetries {
				slog.Warn("webhook delivery failed, retrying",
					"ref", attempt.Ref,
					"attempt", attempt.Attempts,
					"error", lastErr,
				)
				time.Sleep(retryDelay(i))
			}
		}

		if attempt.Status != AttemptStatusDelivered {
			attempt.Status = AttemptStatusFailed
		}
		attempt.History = history

		providerName := d.ProviderName
		if providerName == "" {
			providerName = "unknown"
		}
		webhookAttemptsTotal.WithLabelValues(providerName, string(attempt.Status)).Inc()

		action := "webhook.delivered"
		if attempt.Status != AttemptStatusDelivered {
			action = "webhook.failed"
		}
		if d.AuditLogger != nil {
			d.AuditLogger.Log(context.Background(), action, "webhook", attempt.Ref, "", string(attempt.Status))
		}

		attempt.UpdatedAt = time.Now()
		if err := d.Store.Save(attempt); err != nil {
			slog.Error("failed to save webhook attempt", "error", err, "ref", attempt.Ref)
		}
	}(ctx)
}

func snapshot(a *Attempt) *Attempt {
	clone := *a
	if a.Headers != nil {
		clone.Headers = make(map[string]string, len(a.Headers))
		for k, v := range a.Headers {
			clone.Headers[k] = v
		}
	}
	return &clone
}

func retryDelay(attempt int) time.Duration {
	// Exponential backoff: 100ms, 200ms, 400ms, ... capped at 5 seconds.
	delay := time.Duration(100*(1<<attempt)) * time.Millisecond
	if delay > 5*time.Second {
		delay = 5 * time.Second
	}
	return delay
}

func eventEnabled(eventType string, enabled []string) bool {
	for _, e := range enabled {
		if e == eventType {
			return true
		}
	}
	return false
}
