package v2_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v2 "github.com/Gr3gg0r/openmuara/internal/fawry/v2"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/server"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestWebhookHandlerValidToken(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler("muara-webhook-secret"))

	payload := webhook.FawryV2Payload{
		RequestID:         "req-123",
		FawryRefNumber:    "fawry-456",
		MerchantRefNumber: "ref-123",
		OrderStatus:       "PAID",
		PaymentMethod:     "CARD",
	}
	sig, err := webhook.NewHMACSigner("muara-webhook-secret").Sign(payload)
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}
	payload.MessageSignature = sig

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/fawry/v2/webhook?token=muara-webhook-secret", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerInvalidSignature(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler("muara-webhook-secret"))

	payload := webhook.FawryV2Payload{
		RequestID:         "req-123",
		FawryRefNumber:    "fawry-456",
		MerchantRefNumber: "ref-123",
		OrderStatus:       "PAID",
		PaymentMethod:     "CARD",
		MessageSignature:  "invalid-signature",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/fawry/v2/webhook?token=muara-webhook-secret", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerSkipsSignatureWhenSecretEmpty(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler(""))

	body, _ := json.Marshal(map[string]any{
		"requestId":         "req-123",
		"fawryRefNumber":    "fawry-456",
		"merchantRefNumber": "ref-123",
		"orderStatus":       "PAID",
	})

	req := httptest.NewRequest(http.MethodPost, "/fawry/v2/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerInvalidToken(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodPost, "/fawry/v2/webhook?token=wrong", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp server.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if resp.Error.Code != server.ErrUnauthorized {
		t.Errorf("code: want %q, got %q", server.ErrUnauthorized, resp.Error.Code)
	}
}

func TestWebhookHandlerRejectsWrongMethod(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodGet, "/fawry/v2/webhook?token=muara-webhook-secret", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestWebhookHandlerRejectsInvalidJSON(t *testing.T) {
	handler := server.RequestIDMiddleware(v2.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodPost, "/fawry/v2/webhook?token=muara-webhook-secret", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}

	var resp httputil.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error.Code != httputil.ErrInvalidJSON {
		t.Errorf("code: want %q, got %q", httputil.ErrInvalidJSON, resp.Error.Code)
	}
}
