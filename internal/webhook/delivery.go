package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Gr3gg0r/openmuara/internal/httputil"
)

// DeliveryWorker sends webhook payloads to a subscriber URL with retries.
type DeliveryWorker struct {
	Client *http.Client
}

// NewDeliveryWorker creates a delivery worker with sensible defaults.
func NewDeliveryWorker() *DeliveryWorker {
	return &DeliveryWorker{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Deliver sends the payload and returns the HTTP status code, response body, and error.
// It performs one delivery attempt. Retry logic lives in the caller.
func (w *DeliveryWorker) Deliver(ctx context.Context, url string, payload []byte, headers map[string]string, traceID string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "muara-webhook/1.0")
	if traceID != "" {
		req.Header.Set(httputil.TraceIDHeader, traceID)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("http post: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Warn("failed to close response body", "error", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read response: %w", err)
	}

	return resp.StatusCode, body, nil
}

// IsSuccess reports whether the status code indicates successful delivery.
func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
