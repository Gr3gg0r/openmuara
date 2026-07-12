package billplz_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/openmuara/openmuara/internal/billplz"
	"github.com/openmuara/openmuara/internal/httputil"
)

func newInitializedProvider(t *testing.T) *billplz.Provider {
	t.Helper()
	p := billplz.NewProvider()
	if err := p.Init(validProviderConfig()); err != nil {
		t.Fatalf("init: %v", err)
	}
	p.SetBaseURL("http://localhost:9000")
	return p
}

func createCollection(t *testing.T, p *billplz.Provider) billplz.Collection {
	t.Helper()
	mux := testMux(t, p)
	body, _ := json.Marshal(map[string]string{"title": "Test Collection"})
	req := httptest.NewRequest(http.MethodPost, "/api/v3/collections", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create collection failed: status %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.CollectionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode collection: %v", err)
	}
	return resp.Collection
}

func createBill(t *testing.T, p *billplz.Provider, collectionID string) billplz.Bill {
	t.Helper()
	mux := testMux(t, p)
	reqBody := map[string]any{
		"collection_id": collectionID,
		"email":         "test@example.com",
		"mobile":        "+60123456789",
		"name":          "Test User",
		"amount":        1000,
		"callback_url":  "http://localhost:9999/callback",
		"description":   "Test bill",
		"redirect_url":  "http://localhost:9999/redirect",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v3/bills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("muara-billplz-api-key", "")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create bill failed: status %d, body: %s", rec.Code, rec.Body.String())
	}
	var resp billplz.BillResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode bill: %v", err)
	}
	return resp.Bill
}

func testMux(t *testing.T, p *billplz.Provider) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	for _, route := range p.Routes() {
		mux.Handle(route.Method+" "+route.Path, route.Handler)
	}
	return mux
}

func handlerFor(t *testing.T, p *billplz.Provider, method, path string) http.Handler {
	t.Helper()
	for _, route := range p.Routes() {
		if route.Method == method && route.Path == path {
			return route.Handler
		}
	}
	t.Fatalf("route not found: %s %s", method, path)
	return nil
}

func mustPayRequest(billID string) *http.Request {
	form := url.Values{"outcome": {"pay"}, "csrf_token": {"tok"}}
	req := httptest.NewRequest(http.MethodPost, "/_admin/billplz/pay/"+billID, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req.WithContext(httputil.WithCSRFToken(req.Context(), "tok"))
}

func parseRedirectQuery(loc string) map[string]string {
	u, _ := url.Parse(loc)
	result := make(map[string]string)
	for k, v := range u.Query() {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
