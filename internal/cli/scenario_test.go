package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScenarioCommands(t *testing.T) {
	cases := []struct {
		name    string
		outcome string
		want    string
	}{
		{"success", "success", "/_admin/scenario/success"},
		{"fail", "fail", "/_admin/scenario/fail"},
		{"timeout", "timeout", "/_admin/scenario/timeout"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath, gotQuery string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotQuery = r.URL.RawQuery
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"status":"ok"}`))
			}))
			defer srv.Close()

			cfgPath := writeServerConfig(t, srv)
			oldPath := rootConfigPath
			rootConfigPath = cfgPath
			defer func() { rootConfigPath = oldPath }()

			cmd := newScenarioSubCommand(tc.outcome, "")
			cmd.SetArgs([]string{"ref-1"})
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if gotPath != tc.want {
				t.Errorf("path: want %q, got %q", tc.want, gotPath)
			}
			if gotQuery != "ref=ref-1" {
				t.Errorf("query: want ref=ref-1, got %q", gotQuery)
			}
		})
	}
}
