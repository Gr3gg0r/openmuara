package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateWithDetailsValid(t *testing.T) {
	cfg := &Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
	}
	errs := cfg.ValidateWithDetails("")
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateWithDetailsInvalidPort(t *testing.T) {
	cfg := &Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 70000},
		Persistence: PersistenceConfig{Type: "memory"},
	}
	errs := cfg.ValidateWithDetails("")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "server.port" {
		t.Errorf("expected field server.port, got %q", errs[0].Field)
	}
	if !strings.Contains(errs[0].Message, "65535") {
		t.Errorf("expected port range message, got %q", errs[0].Message)
	}
}

func TestValidateWithDetailsUnknownProvider(t *testing.T) {
	cfg := &Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
		Providers: map[string]ProviderConfig{
			"not-real": {Enabled: true, Config: map[string]any{}},
		},
	}
	errs := cfg.ValidateWithDetails("")
	var found bool
	for _, e := range errs {
		if e.Field == "providers.not-real" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected error for providers.not-real, got %v", errs)
	}
}

func TestValidateWithDetailsFawryVersion(t *testing.T) {
	cfg := &Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
		Providers: map[string]ProviderConfig{
			"fawry": {
				Enabled: true,
				Config: map[string]any{
					"merchant_code":         "mc",
					"merchant_security_key": "sk",
					"webhook_secret":        "ws",
					"version":               "v3",
				},
			},
		},
	}
	errs := cfg.ValidateWithDetails("")
	var found bool
	for _, e := range errs {
		if e.Field == "providers.fawry.config.version" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected error for fawry version, got %v", errs)
	}
}

func TestValidateWithDetailsLineNumbers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `server:
  host: 127.0.0.1
  port: 70000
log:
  level: info
persistence:
  type: memory
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	errs := cfg.ValidateWithDetails(path)
	var found bool
	for _, e := range errs {
		if e.Field == "server.port" {
			found = true
			if e.Line != 3 {
				t.Errorf("expected server.port on line 3, got %d", e.Line)
			}
		}
	}
	if !found {
		t.Fatalf("expected server.port error, got %v", errs)
	}
}

func TestValidateWithDetailsSecurity(t *testing.T) {
	tests := []struct {
		name       string
		cfg        Config
		wantFields []string
	}{
		{
			name: "0.0.0.0 without auth warning",
			cfg: Config{
				Server:      ServerConfig{Host: "0.0.0.0", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
			},
			wantFields: []string{"server.host"},
		},
		{
			name: "admin enabled without credentials",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Admin:       AdminConfig{Enabled: true},
			},
			wantFields: []string{"admin.username", "admin.password_hash"},
		},
		{
			name: "hardened without admin",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Hardened:    true,
			},
			wantFields: []string{"hardened"},
		},
		{
			name: "tls cert key mismatch",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000, TLSCert: "cert.pem"},
				Persistence: PersistenceConfig{Type: "memory"},
			},
			wantFields: []string{"server.tls_cert"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.cfg.ValidateWithDetails("")
			for _, field := range tt.wantFields {
				found := false
				for _, e := range errs {
					if e.Field == field {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error for field %q, got %v", field, errs)
				}
			}
		})
	}
}

func TestFormatValidationErrors(t *testing.T) {
	errs := []ValidationError{
		{Field: "server.port", Message: "out of range", Hint: "use 1-65535", Line: 3},
	}
	out := FormatValidationErrors(errs)
	if !strings.Contains(out, "server.port (line 3)") {
		t.Errorf("expected formatted line number, got %q", out)
	}
	if !strings.Contains(out, "use 1-65535") {
		t.Errorf("expected hint, got %q", out)
	}
}
