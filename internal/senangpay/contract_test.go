package senangpay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/testutil"
)

func TestChargeContract(t *testing.T) {
	reqBody := testutil.GoldenFile(t, "contract/charge_request.json")
	wantBody := testutil.GoldenFile(t, "contract/charge_response.json")

	store := engine.NewMemoryStore()
	handler := NewChargeHandler("muara-senangpay-secret", store, "http://localhost")

	req := httptest.NewRequest(http.MethodPost, "/senangpay/charge", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	assertJSONEqual(t, rec.Body.Bytes(), wantBody)

	tx, ok, err := store.GetByReference("senangpay-contract-ref")
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != ProviderName {
		t.Errorf("provider: want %q, got %q", ProviderName, tx.Provider)
	}
}

func assertJSONEqual(t *testing.T, got, want []byte) {
	t.Helper()

	var gotMap, wantMap map[string]any
	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatalf("decode got: %v", err)
	}
	if err := json.Unmarshal(want, &wantMap); err != nil {
		t.Fatalf("decode want: %v", err)
	}

	if !reflect.DeepEqual(gotMap, wantMap) {
		t.Errorf("response mismatch\ngot:  %s\nwant: %s", got, want)
	}
}
