package fawry

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
	handler := NewChargeHandler("muara-fawry-secret", store)

	req := httptest.NewRequest(http.MethodPost, ChargeHandlerPath, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	assertJSONEqual(t, rec.Body.Bytes(), wantBody)

	// The transaction should be recorded in the ledger with the contract reference.
	tx, ok, err := store.GetByReference("fawry-contract-ref")
	if err != nil {
		t.Fatalf("lookup transaction: %v", err)
	}
	if !ok {
		t.Fatal("expected transaction to be recorded")
	}
	if tx.Provider != "fawry" {
		t.Errorf("provider: want fawry, got %q", tx.Provider)
	}
	if len(tx.Items) != 1 || tx.Items[0].ItemCode != "prod-1" {
		t.Errorf("unexpected items: %+v", tx.Items)
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
