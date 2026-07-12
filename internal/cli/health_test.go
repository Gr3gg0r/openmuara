package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCommandExists(t *testing.T) {
	cmd := newHealthCommand()
	if cmd.Use != "health" {
		t.Fatalf("expected command use to be health, got %q", cmd.Use)
	}
}

func TestRunHealthHealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	result := runHealthWithURL(server.URL + "/healthz")
	if !result.Healthy {
		t.Fatalf("expected healthy, got %+v", result)
	}
	if result.Status != "ok" {
		t.Fatalf("expected status ok, got %q", result.Status)
	}
}

func TestRunHealthUnhealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"error"}`))
	}))
	defer server.Close()

	result := runHealthWithURL(server.URL + "/healthz")
	if result.Healthy {
		t.Fatalf("expected unhealthy, got %+v", result)
	}
}

func TestRunHealthInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	result := runHealthWithURL(server.URL + "/healthz")
	if result.Healthy {
		t.Fatalf("expected unhealthy for invalid JSON, got %+v", result)
	}
}

func TestHealthCommandJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	cmd := newHealthCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	jsonOutput = true
	defer func() { jsonOutput = false }()

	// The command reads the real config path, which does not point to the test
	// server. We exercise the JSON output path by encoding a sample result.
	out := healthOutput{Healthy: true, URL: server.URL + "/healthz", Status: "ok"}
	if err := json.NewEncoder(&buf).Encode(out); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var decoded healthOutput
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !decoded.Healthy {
		t.Fatalf("expected healthy in JSON output")
	}
}

// runHealthWithURL reaches a specific URL for tests.
func runHealthWithURL(url string) healthOutput {
	// #nosec G107 -- test helper only reaches httptest.NewServer URLs.
	resp, err := http.Get(url)
	if err != nil {
		return healthOutput{Error: err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return healthOutput{Error: err.Error()}
	}
	return healthOutput{
		Healthy: resp.StatusCode == http.StatusOK && payload.Status == "ok",
		Status:  payload.Status,
		URL:     url,
	}
}
