package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func TestMetricsEndpointReturnsPrometheusFormat(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	// Record at least one request so the counters appear in the scrape.
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/plain") {
		t.Errorf("content-type: want text/plain, got %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "# HELP openmuara_requests_total") {
		t.Errorf("body missing expected metric description")
	}
}

func TestMetricsMiddlewareRecordsRequests(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	path := "/healthz"
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := rec.Body.String()
	want := `openmuara_requests_total{method="GET",path="/healthz",status="200"} `
	if !strings.Contains(body, want) {
		t.Errorf("metrics missing expected request counter:\n%s", body)
	}
}

func TestMetricsEndpointIsNotCounted(t *testing.T) {
	router := NewRouter(RouterConfig{
		ActiveProvider:   "default",
		EnabledProviders: []string{"default"},
	})

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/metrics", nil))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, `path="/metrics"`) {
		t.Errorf("/metrics should not be counted as a request")
	}
}

var _ = provider.Names
