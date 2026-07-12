package v1_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/openmuara/openmuara/internal/fawry/v1"
	"github.com/openmuara/openmuara/internal/httputil"
	"github.com/openmuara/openmuara/internal/server"
)

func signV1(ref, status, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(ref + status))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestWebhookHandlerValidToken(t *testing.T) {
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler("muara-webhook-secret"))

	payload := v1.WebhookBody{
		MerchantRefNumber: "ref-123",
		OrderStatus:       "PAID",
		MessageSignature:  signV1("ref-123", "PAID", "muara-webhook-secret"),
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/fawry/v1/webhook?token=muara-webhook-secret", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerInvalidSignature(t *testing.T) {
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler("muara-webhook-secret"))

	payload := v1.WebhookBody{
		MerchantRefNumber: "ref-123",
		OrderStatus:       "PAID",
		MessageSignature:  "invalid-signature",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/fawry/v1/webhook?token=muara-webhook-secret", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want 401, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerSkipsSignatureWhenSecretEmpty(t *testing.T) {
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler(""))

	body, _ := json.Marshal(map[string]any{
		"merchantRefNumber": "ref-123",
		"orderStatus":       "PAID",
	})

	req := httptest.NewRequest(http.MethodPost, "/fawry/v1/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestWebhookHandlerInvalidToken(t *testing.T) {
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodPost, "/fawry/v1/webhook?token=wrong", nil)
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
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodGet, "/fawry/v1/webhook?token=muara-webhook-secret", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestWebhookHandlerRejectsInvalidJSON(t *testing.T) {
	handler := server.RequestIDMiddleware(v1.NewWebhookHandler("muara-webhook-secret"))

	req := httptest.NewRequest(http.MethodPost, "/fawry/v1/webhook?token=muara-webhook-secret", bytes.NewReader([]byte("not json")))
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
