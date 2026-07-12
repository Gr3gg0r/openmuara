package cli

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/testutil"
)

func TestWebhookListCommand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_admin/webhooks" {
			t.Errorf("path: want /_admin/webhooks, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newWebhookListCommand()
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

func TestWebhookInspectCommand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/_admin/webhooks/ref-1"
		if r.URL.Path != want {
			t.Errorf("path: want %s, got %s", want, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ref":"ref-1"}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newWebhookInspectCommand()
	cmd.SetArgs([]string{"ref-1"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "ref-1") {
		t.Errorf("output missing ref: %q", buf.String())
	}
}

func TestWebhookReplayCommand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/_admin/webhooks/ref-1/replay"
		if r.URL.Path != want {
			t.Errorf("path: want %s, got %s", want, r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method: want POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newWebhookReplayCommand()
	cmd.SetArgs([]string{"ref-1"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(buf.String(), "ok") {
		t.Errorf("output missing ok: %q", buf.String())
	}
}

func TestWebhookServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	cfgPath := writeServerConfig(t, srv)
	oldPath := rootConfigPath
	rootConfigPath = cfgPath
	defer func() { rootConfigPath = oldPath }()

	cmd := newWebhookListCommand()
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

// writeServerConfig writes a config file pointing at the httptest server.
func writeServerConfig(t *testing.T, srv *httptest.Server) string {
	t.Helper()
	host, portStr, err := splitHostPort(strings.TrimPrefix(srv.URL, "http://"))
	if err != nil {
		t.Fatalf("split server address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}

	dir := testutil.TempWorkspace(t)
	cfg := fmt.Sprintf("server:\n  host: %s\n  port: %d\nlog:\n  level: error\npersistence:\n  type: memory\n", host, port)
	testutil.WriteFile(t, dir, "config.yml", []byte(cfg))
	return filepath.Join(dir, "config.yml")
}

// splitHostPort is a tiny wrapper so tests do not need to import net.
func splitHostPort(hostport string) (string, string, error) {
	parts := strings.Split(hostport, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address %q", hostport)
	}
	return parts[0], parts[1], nil
}
