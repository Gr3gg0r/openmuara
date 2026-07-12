package stripe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
	"github.com/openmuara/openmuara/internal/testutil"
)

func TestCreateCheckoutSessionContract(t *testing.T) {
	reqBody := testutil.GoldenFile(t, "contract/create_session_request.json")
	wantBody := testutil.GoldenFile(t, "contract/create_session_response.json")

	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewCreateCheckoutSessionHandler(sessions, ledger, "http://localhost")

	req := httptest.NewRequest(http.MethodPost, "/v1/checkout/sessions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// Capture the generated session ID before normalizing it for comparison.
	sessionID, _ := got["id"].(string)

	// The session ID and URL contain a generated UUID; normalize them for comparison.
	got["id"] = "<id>"
	got["url"] = "<url>"

	var want map[string]any
	if err := json.Unmarshal(wantBody, &want); err != nil {
		t.Fatalf("decode golden response: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.MarshalIndent(got, "", "  ")
		wantJSON, _ := json.MarshalIndent(want, "", "  ")
		t.Errorf("response mismatch\ngot:\n%s\nwant:\n%s", gotJSON, wantJSON)
	}

	// The ledger should hold the transaction keyed by the generated session ID.
	tx, ok, err := ledger.GetByReference(sessionID)
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != ProviderName {
		t.Errorf("provider: want %q, got %q", ProviderName, tx.Provider)
	}
	if tx.IdempotencyKey != "stripe-contract-ref" {
		t.Errorf("client_reference_id: want stripe-contract-ref, got %q", tx.IdempotencyKey)
	}
}
