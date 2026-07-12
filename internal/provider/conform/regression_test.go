package conform

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/provider"
)

// fakeProvider is a minimal provider used for regression tests.
type fakeProvider struct {
	name     string
	initErr  error
	routes   []provider.Route
	versions []string
}

func (f *fakeProvider) Name() string                 { return f.name }
func (f *fakeProvider) Init(_ map[string]any) error  { return f.initErr }
func (f *fakeProvider) Routes() []provider.Route     { return f.routes }
func (f *fakeProvider) ChargeHandler() http.Handler  { return nil }
func (f *fakeProvider) WebhookHandler() http.Handler { return nil }
func (f *fakeProvider) EscapeHandler() http.Handler  { return nil }
func (f *fakeProvider) PayloadBuilder() func(context.Context, provider.Transaction) ([]byte, error) {
	return func(context.Context, provider.Transaction) ([]byte, error) { return nil, nil }
}

// fakeVersionedProvider extends fakeProvider with version support.
type fakeVersionedProvider struct {
	fakeProvider
}

func (f *fakeVersionedProvider) Versions() []string     { return f.versions }
func (f *fakeVersionedProvider) CurrentVersion() string { return "v1" }

func TestCapture_InitErrorCapturesVersions(t *testing.T) {
	wantErr := errors.New("init failed")
	p := &fakeVersionedProvider{fakeProvider{
		name:     "capture-failing",
		initErr:  wantErr,
		versions: []string{"v1", "v2"},
	}}

	snap, err := Capture(p)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected init error %v, got %v", wantErr, err)
	}
	if snap.Name != p.name {
		t.Errorf("name: got %q, want %q", snap.Name, p.name)
	}
	if len(snap.Versions) != 2 {
		t.Errorf("expected versions to be captured on init error, got %v", snap.Versions)
	}
	if len(snap.Routes) != 0 {
		t.Errorf("expected no routes on init error, got %v", snap.Routes)
	}
}

func TestCapture_Success(t *testing.T) {
	p := &fakeProvider{
		name: "capture-ok",
		routes: []provider.Route{
			{Method: "GET", Path: "/test"},
		},
	}

	snap, err := Capture(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Routes) != 1 || snap.Routes[0].Method != "GET" || snap.Routes[0].Path != "/test" {
		t.Errorf("unexpected routes: %v", snap.Routes)
	}
}

func TestUsage(t *testing.T) {
	got := Usage()
	if got == "" {
		t.Fatal("Usage() returned empty string")
	}
	if !strings.Contains(got, "-"+UpdateFlag) {
		t.Errorf("Usage() %q does not reference -%s", got, UpdateFlag)
	}
}

func TestCompare_Update(t *testing.T) {
	name := "regression-fake"
	p := &fakeProvider{
		name: name,
		routes: []provider.Route{
			{Method: "POST", Path: "/regression"},
		},
	}

	path := GoldenPath(t, name)
	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	Compare(t, p, true)

	// #nosec G304 -- path is produced by GoldenPath from a known base directory and provider name.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden file not written: %v", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("golden file is invalid JSON: %v", err)
	}
	if snap.Name != name {
		t.Errorf("golden name: got %q, want %q", snap.Name, name)
	}
	if len(snap.Routes) != 1 {
		t.Errorf("golden routes: got %v, want 1 route", snap.Routes)
	}
}
