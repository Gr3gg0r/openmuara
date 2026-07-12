package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("host: want 127.0.0.1, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("port: want 9000, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("log level: want info, got %q", cfg.Log.Level)
	}
	fawryCfg, ok := cfg.Providers["fawry"]
	if !ok || !fawryCfg.Enabled {
		t.Errorf("providers.fawry: want enabled by default")
	}
	if got := fawryCfg.Config["merchant_code"]; got != "muara-merchant-code" {
		t.Errorf("providers.fawry.config.merchant_code: want muara-merchant-code, got %v", got)
	}
	if cfg.Fawry.MerchantCode != "" {
		t.Errorf("legacy fawry.merchant_code: want empty default, got %q", cfg.Fawry.MerchantCode)
	}
	if cfg.Admin.Enabled {
		t.Errorf("admin.enabled: want false by default")
	}
	if cfg.RateLimit.Enabled {
		t.Errorf("rate_limit.enabled: want false by default")
	}
	if cfg.Hardened {
		t.Errorf("hardened: want false by default")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
server:
  host: 0.0.0.0
  port: 8080
log:
  level: debug
fawry:
  merchant_code: file-merchant
  merchant_security_key: file-secret
  webhook_secret: file-webhook
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("port: want 8080, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("log level: want debug, got %q", cfg.Log.Level)
	}
	if cfg.Fawry.MerchantCode != "file-merchant" {
		t.Errorf("merchant code: want file-merchant, got %q", cfg.Fawry.MerchantCode)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("MUARA_SERVER_PORT", "7000")
	t.Setenv("MUARA_LOG_LEVEL", "warn")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 7000 {
		t.Errorf("port: want 7000, got %d", cfg.Server.Port)
	}
	if cfg.Log.Level != "warn" {
		t.Errorf("log level: want warn, got %q", cfg.Log.Level)
	}
}

func TestLoadWebhookDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Webhook.URL != "" {
		t.Errorf("webhook url default: want empty, got %q", cfg.Webhook.URL)
	}
	if cfg.Webhook.MaxRetries != 3 {
		t.Errorf("webhook max retries default: want 3, got %d", cfg.Webhook.MaxRetries)
	}
	if cfg.Webhook.Targets == nil {
		t.Error("webhook targets default: want empty map, got nil")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid defaults",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "sqlite", Path: ".muara/data/ledger.db"},
			},
		},
		{
			name: "empty host",
			cfg: Config{
				Server:      ServerConfig{Host: "", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
			},
			wantErr: true,
		},
		{
			name: "port out of range",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 70000},
				Persistence: PersistenceConfig{Type: "memory"},
			},
			wantErr: true,
		},
		{
			name: "unsupported persistence",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "postgres"},
			},
			wantErr: true,
		},
		{
			name: "sqlite missing path",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "sqlite"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateWebhookURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  WebhookConfig
		want string
	}{
		{
			name: "empty url",
			cfg:  WebhookConfig{URL: ""},
			want: "webhook.url is empty; outgoing webhooks will not be dispatched",
		},
		{
			name: "valid http url",
			cfg:  WebhookConfig{URL: "http://localhost:3000/webhook"},
			want: "",
		},
		{
			name: "valid https url",
			cfg:  WebhookConfig{URL: "https://example.com/webhook"},
			want: "",
		},
		{
			name: "unsupported scheme",
			cfg:  WebhookConfig{URL: "ftp://localhost/webhook"},
			want: `webhook.url has unsupported scheme "ftp"; must be http or https`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateWebhookURL(tt.cfg)
			if got != tt.want {
				t.Errorf("ValidateWebhookURL(): want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestLoadNewStyleProviders(t *testing.T) {
	// Given a new-style config with providers.fawry.enabled: true
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
providers:
  fawry:
    enabled: true
    config:
      merchant_code: new-merchant
      merchant_security_key: new-secret
      webhook_secret: new-webhook
  default:
    enabled: false
    config: {}
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// When loaded
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then cfg.Providers["fawry"] is populated.
	fawryCfg, ok := cfg.Providers["fawry"]
	if !ok {
		t.Fatal("expected providers.fawry to be present")
	}
	if !fawryCfg.Enabled {
		t.Error("expected providers.fawry to be enabled")
	}
	if fawryCfg.Config["merchant_code"] != "new-merchant" {
		t.Errorf("merchant_code: want new-merchant, got %v", fawryCfg.Config["merchant_code"])
	}

	defaultCfg, ok := cfg.Providers["default"]
	if !ok {
		t.Fatal("expected providers.default to be present")
	}
	if defaultCfg.Enabled {
		t.Error("expected providers.default to be disabled")
	}
}

func TestLoadOldStyleFawryNormalized(t *testing.T) {
	// Given an old-style config with top-level fawry fields
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
fawry:
  merchant_code: legacy-merchant
  merchant_security_key: legacy-secret
  webhook_secret: legacy-webhook
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// When loaded
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then it is mapped to providers.fawry.
	fawryCfg, ok := cfg.Providers["fawry"]
	if !ok {
		t.Fatal("expected providers.fawry to be present after normalization")
	}
	if !fawryCfg.Enabled {
		t.Error("expected normalized providers.fawry to be enabled")
	}
	if fawryCfg.Config["merchant_code"] != "legacy-merchant" {
		t.Errorf("merchant_code: want legacy-merchant, got %v", fawryCfg.Config["merchant_code"])
	}
}

func TestLoadStripeDefaults(t *testing.T) {
	// Given default config
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then providers.stripe is present with defaults and disabled.
	stripeCfg, ok := cfg.Providers["stripe"]
	if !ok {
		t.Fatal("expected providers.stripe to be present")
	}
	if stripeCfg.Enabled {
		t.Error("expected providers.stripe to be disabled by default")
	}
	if stripeCfg.Config["publishable_key"] != "pk_test_muara" {
		t.Errorf("publishable_key: want pk_test_muara, got %v", stripeCfg.Config["publishable_key"])
	}
	if stripeCfg.Config["secret_key"] != "sk_test_muara" {
		t.Errorf("secret_key: want sk_test_muara, got %v", stripeCfg.Config["secret_key"])
	}
}

func TestLoadSecurityFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
server:
  host: 0.0.0.0
  tls_cert: /path/to/cert.pem
  tls_key: /path/to/key.pem
admin:
  enabled: true
  username: admin
  password_hash: "$2a$10$..."
  token: tok_test
rate_limit:
  enabled: true
  requests_per_minute: 200
hardened: true
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("host: want 0.0.0.0, got %q", cfg.Server.Host)
	}
	if cfg.Server.TLSCert != "/path/to/cert.pem" {
		t.Errorf("tls_cert: want /path/to/cert.pem, got %q", cfg.Server.TLSCert)
	}
	if !cfg.Admin.Enabled {
		t.Error("admin.enabled: want true")
	}
	if cfg.Admin.Username != "admin" {
		t.Errorf("admin.username: want admin, got %q", cfg.Admin.Username)
	}
	if cfg.Admin.PasswordHash != "$2a$10$..." {
		t.Errorf("admin.password_hash mismatch")
	}
	if cfg.Admin.Token != "tok_test" {
		t.Errorf("admin.token: want tok_test, got %q", cfg.Admin.Token)
	}
	if !cfg.RateLimit.Enabled {
		t.Error("rate_limit.enabled: want true")
	}
	if cfg.RateLimit.RequestsPerMinute != 200 {
		t.Errorf("rate_limit.requests_per_minute: want 200, got %d", cfg.RateLimit.RequestsPerMinute)
	}
	if !cfg.Hardened {
		t.Error("hardened: want true")
	}
}

func TestEnvOverrideSecurity(t *testing.T) {
	t.Setenv("MUARA_ADMIN_ENABLED", "true")
	t.Setenv("MUARA_ADMIN_USERNAME", "env-admin")
	t.Setenv("MUARA_ADMIN_TOKEN", "tok_env")
	t.Setenv("MUARA_RATE_LIMIT_ENABLED", "true")
	t.Setenv("MUARA_HARDENED", "true")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Admin.Enabled {
		t.Error("admin.enabled: want true from env")
	}
	if cfg.Admin.Username != "env-admin" {
		t.Errorf("admin.username: want env-admin, got %q", cfg.Admin.Username)
	}
	if cfg.Admin.Token != "tok_env" {
		t.Errorf("admin.token: want tok_env, got %q", cfg.Admin.Token)
	}
	if !cfg.RateLimit.Enabled {
		t.Error("rate_limit.enabled: want true from env")
	}
	if !cfg.Hardened {
		t.Error("hardened: want true from env")
	}
}

func TestConfigValidateSecurity(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "admin enabled without username",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Admin:       AdminConfig{Enabled: true, PasswordHash: "$2a$10$..."},
			},
			wantErr: true,
		},
		{
			name: "admin enabled without credentials",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Admin:       AdminConfig{Enabled: true, Username: "admin"},
			},
			wantErr: true,
		},
		{
			name: "hardened without admin",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Hardened:    true,
			},
			wantErr: true,
		},
		{
			name: "tls cert without key",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000, TLSCert: "cert.pem"},
				Persistence: PersistenceConfig{Type: "memory"},
			},
			wantErr: true,
		},
		{
			name: "rate limit invalid rpm",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				RateLimit:   RateLimitConfig{Enabled: true, RequestsPerMinute: 0},
			},
			wantErr: true,
		},
		{
			name: "valid hardened config",
			cfg: Config{
				Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
				Persistence: PersistenceConfig{Type: "memory"},
				Admin:       AdminConfig{Enabled: true, Username: "admin", Token: "tok"},
				Hardened:    true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvOverrideStripeSecretKey(t *testing.T) {
	// Given environment override for Stripe secret key
	t.Setenv("MUARA_STRIPE_SECRET_KEY", "sk_env_override")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stripeCfg := cfg.Providers["stripe"]
	if stripeCfg.Config["secret_key"] != "sk_env_override" {
		t.Errorf("secret_key: want sk_env_override, got %v", stripeCfg.Config["secret_key"])
	}
}

func TestDefaultYAMLContainsExpectedKeys(t *testing.T) {
	yaml := DefaultYAML()
	if len(yaml) == 0 {
		t.Fatal("expected non-empty default YAML")
	}
	want := []string{"server:", "admin:", "rate_limit:", "hardened:", "providers:", "fawry:", "persistence:", "webhook:"}
	for _, key := range want {
		if !strings.Contains(string(yaml), key) {
			t.Errorf("default YAML missing %q", key)
		}
	}
}

func TestConfigValidateEnabledProviderNotRegistered(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
		Providers: map[string]ProviderConfig{
			"unknown": {Enabled: true},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unregistered provider")
	}
}

func TestConfigValidateProviderInitFailure(t *testing.T) {
	cfg := Config{
		Server:      ServerConfig{Host: "127.0.0.1", Port: 9000},
		Persistence: PersistenceConfig{Type: "memory"},
		Providers: map[string]ProviderConfig{
			"fawry": {Enabled: true, Config: map[string]any{"merchant_code": ""}},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid provider config")
	}
}

func TestLoadNewStyleTakesPrecedenceOverOldStyle(t *testing.T) {
	// Given both old and new config
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := `
fawry:
  merchant_code: legacy-merchant
  merchant_security_key: legacy-secret
  webhook_secret: legacy-webhook
providers:
  fawry:
    enabled: true
    config:
      merchant_code: new-merchant
      merchant_security_key: new-secret
      webhook_secret: new-webhook
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// When loaded
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then new config takes precedence.
	fawryCfg := cfg.Providers["fawry"]
	if fawryCfg.Config["merchant_code"] != "new-merchant" {
		t.Errorf("merchant_code: want new-merchant, got %v", fawryCfg.Config["merchant_code"])
	}
}
