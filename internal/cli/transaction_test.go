package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransactionListCommand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_admin/transactions" {
			t.Errorf("path: want /_admin/transactions, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newTransactionListCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "results") {
		t.Errorf("output missing results: %q", buf.String())
	}
}

func TestTransactionInspectCommand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/_admin/transactions/ref-1"
		if r.URL.Path != want {
			t.Errorf("path: want %s, got %s", want, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"transaction":{"reference":"ref-1","trace_id":"trace-1"}}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newTransactionInspectCommand()
	cmd.SetArgs([]string{"ref-1"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "trace-1") {
		t.Errorf("output missing trace_id: %q", buf.String())
	}
}

func TestTransactionInspectCommandEscapesRef(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/_admin/transactions/ref%2F1"
		if r.URL.Path != want && r.URL.EscapedPath() != want {
			t.Errorf("path: want %s, got %s (escaped %s)", want, r.URL.Path, r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"transaction":{"reference":"ref/1"}}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newTransactionInspectCommand()
	cmd.SetArgs([]string{"ref/1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
}
