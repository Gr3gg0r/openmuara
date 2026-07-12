package webhook

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/provider"
)

type errAttemptStore struct {
	AttemptStore
	getErr  error
	saveErr error
	listErr error
}

func (e *errAttemptStore) Get(ref string) (*Attempt, error) {
	if e.getErr != nil {
		return nil, e.getErr
	}
	return e.AttemptStore.Get(ref)
}

func (e *errAttemptStore) Save(a *Attempt) error {
	if e.saveErr != nil {
		return e.saveErr
	}
	return e.AttemptStore.Save(a)
}

func (e *errAttemptStore) List(limit, offset int) ([]*Attempt, error) {
	if e.listErr != nil {
		return nil, e.listErr
	}
	return e.AttemptStore.List(limit, offset)
}

func TestDispatcherSkipsWhenURLMissing(t *testing.T) {
	d, err := NewDispatcher(FawryV2, "", "secret", 0, engine.NewMemoryStore())
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}

	_, err = d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error when url missing")
	}
}

func TestDispatcherDeliversSuccessfully(t *testing.T) {
	var received atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		received.Store(true)
		body, _ := io.ReadAll(req.Body)
		if len(body) == 0 {
			t.Error("expected non-empty body")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	d, err := NewDispatcher(FawryV2, server.URL, "secret", 0, store)
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}

	attempt, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	// Wait for async delivery.
	time.Sleep(200 * time.Millisecond)

	if !received.Load() {
		t.Error("server did not receive webhook")
	}

	updated, err := d.Store.Get(attempt.Ref)
	if err != nil {
		t.Fatalf("get attempt: %v", err)
	}
	if updated.Status != AttemptStatusDelivered {
		t.Errorf("status: want delivered, got %q", updated.Status)
	}
}

func TestDispatcherRetriesOnFailure(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	d, err := NewDispatcher(FawryV2, server.URL, "secret", 2, store)
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}

	_, err = d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	// Wait for async delivery + retries.
	time.Sleep(500 * time.Millisecond)

	if attempts.Load() != 3 {
		t.Errorf("want 3 delivery attempts, got %d", attempts.Load())
	}
}

func TestDispatcherReplay(t *testing.T) {
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	d, err := NewDispatcher(FawryV2, server.URL, "secret", 0, store)
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}

	_, err = d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	_, err = d.Replay(context.Background(), "ref-1")
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	if count.Load() != 2 {
		t.Errorf("want 2 deliveries, got %d", count.Load())
	}
}

func TestDispatcherReplayMissingRef(t *testing.T) {
	d, err := NewDispatcher(FawryV2, "http://localhost", "secret", 0, engine.NewMemoryStore())
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}

	_, err = d.Replay(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing ref")
	}
}

func TestDispatcherBuilderError(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", 0, func(context.Context, provider.Transaction) ([]byte, error) {
		return nil, errors.New("build failed")
	}, nil)

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error from builder")
	}
}

func TestDispatcherHeaderBuilderError(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil },
		func(context.Context, provider.Transaction) (map[string]string, error) {
			return nil, errors.New("header failed")
		},
	)

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error from header builder")
	}
}

func TestDispatcherStoreSaveError(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Store = &errAttemptStore{saveErr: errors.New("save failed")}

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error saving attempt")
	}
}

func TestDispatcherReplayStoreSaveError(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	_ = d.Store.Save(&Attempt{Ref: "ref-1"})

	d.Store = &errAttemptStore{AttemptStore: d.Store, saveErr: errors.New("save failed")}
	_, err := d.Replay(context.Background(), "ref-1")
	if err == nil {
		t.Fatal("expected error saving replay attempt")
	}
}

func TestDispatcherSnapshotCopiesHeaders(t *testing.T) {
	a := &Attempt{Ref: "r", Headers: map[string]string{"k": "v"}}
	clone := snapshot(a)
	clone.Headers["k"] = "changed"
	if a.Headers["k"] != "v" {
		t.Error("snapshot should copy headers")
	}
}

func TestRetryDelay(t *testing.T) {
	if retryDelay(0) != 100*time.Millisecond {
		t.Errorf("retryDelay(0): want 100ms, got %v", retryDelay(0))
	}
	if retryDelay(10) != 5*time.Second {
		t.Errorf("retryDelay(10): want 5s cap, got %v", retryDelay(10))
	}
}

type fakeVerifier struct {
	valid bool
	err   error
}

func (f fakeVerifier) VerifyWebhookSignature([]byte, map[string]string) (bool, error) {
	return f.valid, f.err
}

func TestDispatcherSetsSignatureValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Verifier = fakeVerifier{valid: true}

	attempt, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if attempt.SignatureValid == nil || !*attempt.SignatureValid {
		t.Errorf("expected signature_valid true, got %v", attempt.SignatureValid)
	}
}

func TestDispatcherReplayCopiesSignatureValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Verifier = fakeVerifier{valid: true}

	if _, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	replay, err := d.Replay(context.Background(), "ref-1")
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if replay.SignatureValid == nil || !*replay.SignatureValid {
		t.Errorf("expected replay signature_valid copied, got %v", replay.SignatureValid)
	}
}

func TestDispatcherRecordsHistory(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 3,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)

	if _, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	time.Sleep(800 * time.Millisecond)

	updated, err := d.Store.Get("ref-1")
	if err != nil {
		t.Fatalf("get attempt: %v", err)
	}
	if len(updated.History) < 3 {
		t.Errorf("expected at least 3 history entries, got %d", len(updated.History))
	}
	if updated.Status != AttemptStatusDelivered {
		t.Errorf("expected delivered after retries, got %q", updated.Status)
	}
}

func TestDispatcherPropagatesTraceIDFromContext(t *testing.T) {
	headerValue := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		headerValue <- req.Header.Get(httputil.TraceIDHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := NewDispatcherFromBuilder(server.URL, 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)

	ctx := httputil.WithTraceID(context.Background(), "trace-ctx")
	attempt, err := d.Dispatch(ctx, "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	select {
	case got := <-headerValue:
		if got != "trace-ctx" {
			t.Errorf("X-Trace-Id header: want trace-ctx, got %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for webhook delivery")
	}
	if attempt.TraceID != "trace-ctx" {
		t.Errorf("attempt.TraceID: want trace-ctx, got %q", attempt.TraceID)
	}
}

func TestDispatcherFallsBackToLedgerTraceID(t *testing.T) {
	headerValue := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		headerValue <- req.Header.Get(httputil.TraceIDHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0, TraceID: "trace-ledger"}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	d := NewDispatcherFromBuilder(server.URL, 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Ledger = store

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	select {
	case got := <-headerValue:
		if got != "trace-ledger" {
			t.Errorf("X-Trace-Id header: want trace-ledger, got %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for webhook delivery")
	}
}

func TestDispatcherContextTraceIDOverridesLedger(t *testing.T) {
	headerValue := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		headerValue <- req.Header.Get(httputil.TraceIDHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	store := engine.NewMemoryStore()
	if _, _, err := store.CreateOrGet(engine.Transaction{Reference: "ref-1", Amount: 100.0, TraceID: "trace-ledger"}); err != nil {
		t.Fatalf("seed transaction: %v", err)
	}

	d := NewDispatcherFromBuilder(server.URL, 0,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)
	d.Ledger = store

	ctx := httputil.WithTraceID(context.Background(), "trace-ctx")
	_, err := d.Dispatch(ctx, "ref-1", PaymentStatusPaid)
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	select {
	case got := <-headerValue:
		if got != "trace-ctx" {
			t.Errorf("X-Trace-Id header: want trace-ctx, got %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for webhook delivery")
	}
}
