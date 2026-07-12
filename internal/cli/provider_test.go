package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProviderInitCreatesGatewayYAML(t *testing.T) {
	dir := t.TempDir()
	pluginsDir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	var buf bytes.Buffer
	if err := runProviderInit(&buf, "my-gateway"); err != nil {
		t.Fatalf("init: %v", err)
	}

	path := filepath.Join("plugins", "my-gateway", "gateway.yml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("gateway.yml not created: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join("plugins", "my-gateway")) })

	if !strings.Contains(buf.String(), "Created") {
		t.Errorf("output did not report creation: %q", buf.String())
	}
}

func TestProviderTestSenangpay(t *testing.T) {
	var buf bytes.Buffer
	if err := runProviderTest(&buf, "../../plugins", "senangpay"); err != nil {
		t.Fatalf("test: %v", err)
	}
	if !strings.Contains(buf.String(), "HTTP 200") {
		t.Errorf("expected 200 response, got %q", buf.String())
	}
}

func TestProviderTestNonSimpleErrors(t *testing.T) {
	var buf bytes.Buffer
	if err := runProviderTest(&buf, "../../plugins", "stripe"); err == nil {
		t.Fatal("expected error for non-simple provider")
	}
}
