package webhook

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestEventEnabled(t *testing.T) {
	cases := []struct {
		name     string
		event    string
		enabled  []string
		expected bool
	}{
		{"enabled", "charge.paid", []string{"charge.paid", "charge.failed"}, true},
		{"disabled", "charge.paid", []string{"charge.failed"}, false},
		{"empty list", "charge.paid", []string{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := eventEnabled(tc.event, tc.enabled); got != tc.expected {
				t.Errorf("eventEnabled(%q, %v) = %v, want %v", tc.event, tc.enabled, got, tc.expected)
			}
		})
	}
}

func TestNewDispatcherFromBuilder_NegativeRetries(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", -5,
		func(context.Context, provider.Transaction) ([]byte, error) { return []byte("{}"), nil }, nil)

	if d.MaxRetries != 0 {
		t.Errorf("MaxRetries: got %d, want 0", d.MaxRetries)
	}
}

func TestDispatcherFiltersDisabledEvents(t *testing.T) {
	d := NewDispatcherFromBuilder("http://localhost", 0,
		func(context.Context, provider.Transaction) ([]byte, error) {
			t.Error("builder should not be called when event is disabled")
			return nil, nil
		}, nil)
	d.EventTypeFor = func(_ string, _ PaymentStatus) string { return "charge.paid" }
	d.EnabledEvents = []string{"charge.failed"}

	_, err := d.Dispatch(context.Background(), "ref-1", PaymentStatusPaid)
	if err == nil {
		t.Fatal("expected error for disabled event")
	}

	var ec *errcode.Error
	if !errors.As(err, &ec) || ec.Code != errcode.EProviderDisabled {
		t.Errorf("expected EProviderDisabled, got %v", err)
	}
}

func TestMemoryStoreList_NegativeOffset(t *testing.T) {
	store := NewMemoryStore()
	for i := 0; i < 2; i++ {
		if err := store.Save(&Attempt{Ref: string(rune('a' + i)), ID: "id"}); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	got, err := store.List(0, -10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("want 2 attempts, got %d", len(got))
	}
}

func TestMemoryStoreList_OffsetOverflow(t *testing.T) {
	store := NewMemoryStore()
	if err := store.Save(&Attempt{Ref: "ref-1", ID: "id"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := store.List(0, 100)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("want 0 attempts, got %d", len(got))
	}
}

// errCloseReadCloser is an io.ReadCloser that fails on Close.
type errCloseReadCloser struct{ data []byte }

func (e *errCloseReadCloser) Read(p []byte) (int, error) {
	if len(e.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, e.data)
	e.data = e.data[n:]
	return n, nil
}

func (e *errCloseReadCloser) Close() error { return errors.New("close failed") }

// fakeRoundTripper returns a controlled HTTP response.
type fakeRoundTripper struct {
	status int
	body   io.ReadCloser
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.status,
		Body:       f.body,
		Request:    req,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func TestDeliveryWorker_LogsBodyCloseError(t *testing.T) {
	var buf bytes.Buffer
	orig := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})))
	t.Cleanup(func() { slog.SetDefault(orig) })

	worker := &DeliveryWorker{
		Client: &http.Client{Transport: &fakeRoundTripper{
			status: http.StatusOK,
			body:   &errCloseReadCloser{data: []byte(`{}`)},
		}},
	}

	status, _, err := worker.Deliver(context.Background(), "http://localhost", []byte(`{}`), nil, "trace-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("status: want 200, got %d", status)
	}

	if !strings.Contains(buf.String(), "failed to close response body") {
		t.Errorf("expected body-close warning, got %q", buf.String())
	}
}
