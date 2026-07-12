package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSecurityHashPassword(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newSecurityHashPasswordCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--password", "secret"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if out == "" {
		t.Fatal("expected hash output")
	}
	if !strings.HasPrefix(out, "$2a$") {
		t.Errorf("expected bcrypt hash, got %q", out)
	}
}

func TestSecurityHashPasswordMissing(t *testing.T) {
	cmd := newSecurityHashPasswordCommand()
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing password")
	}
}

func TestSecurityGenCert(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")

	buf := new(bytes.Buffer)
	cmd := newSecurityGenCertCommand()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--host", "localhost", "--cert-out", certPath, "--key-out", keyPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Errorf("cert file not created: %v", err)
	}
	if _, err := os.Stat(keyPath); err != nil {
		t.Errorf("key file not created: %v", err)
	}

	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("stat key: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("key file mode: want 0600, got %o", info.Mode().Perm())
	}
}

func TestSecurityAudit(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yml")
	content := `
server:
  host: 0.0.0.0
  port: 9000
`
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	rootConfigPath = configPath
	defer func() { rootConfigPath = ".muara/config.yml" }()

	buf := new(bytes.Buffer)
	cmd := newSecurityAuditCommand()
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "bind:        0.0.0.0") {
		t.Errorf("expected bind in output, got:\n%s", out)
	}
	if !strings.Contains(out, "server is bound to 0.0.0.0 without admin authentication") {
		t.Errorf("expected issue in output, got:\n%s", out)
	}
}
