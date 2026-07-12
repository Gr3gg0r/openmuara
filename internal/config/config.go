// Package config loads and exposes muara configuration from YAML and environment variables.
package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/factory"
	"github.com/spf13/viper"
)

// Config is the top-level muara configuration.
type Config struct {
	Server             ServerConfig              `mapstructure:"server"`
	Log                LogConfig                 `mapstructure:"log"`
	Persistence        PersistenceConfig         `mapstructure:"persistence"`
	Admin              AdminConfig               `mapstructure:"admin"`
	Viewer             ViewerConfig              `mapstructure:"viewer"`
	RateLimit          RateLimitConfig           `mapstructure:"rate_limit"`
	Hardened           bool                      `mapstructure:"hardened"`
	DisableUpdateCheck bool                      `mapstructure:"disable_update_check"`
	Fawry              FawryConfig               `mapstructure:"fawry"`
	Stripe             StripeConfig              `mapstructure:"stripe"`
	Webhook            WebhookConfig             `mapstructure:"webhook"`
	Providers          map[string]ProviderConfig `mapstructure:"providers"`
	Dev                DevConfig                 `mapstructure:"dev"`
}

// PersistenceConfig controls how the transaction ledger is stored.
type PersistenceConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

// ProviderConfig holds generic configuration for a registered provider.
type ProviderConfig struct {
	Enabled bool           `mapstructure:"enabled"`
	Config  map[string]any `mapstructure:"config"`
}

// WebhookConfig holds outgoing webhook settings.
type WebhookConfig struct {
	URL        string              `mapstructure:"url"`
	MaxRetries int                 `mapstructure:"max_retries"`
	Targets    map[string]string   `mapstructure:"targets"`
	Events     map[string][]string `mapstructure:"events"`
}

// CORSConfig holds cross-origin resource sharing settings.
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// CSRFConfig holds cross-site request forgery protection settings.
type CSRFConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host               string     `mapstructure:"host"`
	Port               int        `mapstructure:"port"`
	AdminPort          int        `mapstructure:"admin_port"`
	PublicBaseURL      string     `mapstructure:"public_base_url"`
	AdminPublicBaseURL string     `mapstructure:"admin_public_base_url"`
	TLSCert            string     `mapstructure:"tls_cert"`
	TLSKey             string     `mapstructure:"tls_key"`
	CORS               CORSConfig `mapstructure:"cors"`
	CSRF               CSRFConfig `mapstructure:"csrf"`
	Pprof              bool       `mapstructure:"pprof"`
}

// AdminConfig holds admin authentication settings.
type AdminConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
	Token        string `mapstructure:"token"`
}

// ViewerConfig holds read-only dashboard authentication settings.
type ViewerConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
	Token        string `mapstructure:"token"`
}

// RateLimitConfig holds in-memory rate limiting settings.
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
}

// DevConfig holds development-only convenience flags. These should never be
// enabled in production because they mutate state with synthetic data.
type DevConfig struct {
	Seed bool `mapstructure:"seed"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// FawryConfig holds default Fawry-style gateway credentials.
type FawryConfig struct {
	MerchantCode        string `mapstructure:"merchant_code"`
	MerchantSecurityKey string `mapstructure:"merchant_security_key"`
	WebhookSecret       string `mapstructure:"webhook_secret"`
}

// StripeConfig holds default Stripe-style credentials.
type StripeConfig struct {
	APIKey        string `mapstructure:"api_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

// bindSecurityEnvVars binds security-related config keys to optional env vars.
func bindSecurityEnvVars(v *viper.Viper) {
	_ = v.BindEnv("admin.enabled", "MUARA_ADMIN_ENABLED")
	_ = v.BindEnv("admin.username", "MUARA_ADMIN_USERNAME")
	_ = v.BindEnv("admin.password_hash", "MUARA_ADMIN_PASSWORD_HASH")
	_ = v.BindEnv("admin.token", "MUARA_ADMIN_TOKEN")
	_ = v.BindEnv("viewer.enabled", "MUARA_VIEWER_ENABLED")
	_ = v.BindEnv("viewer.username", "MUARA_VIEWER_USERNAME")
	_ = v.BindEnv("viewer.password_hash", "MUARA_VIEWER_PASSWORD_HASH")
	_ = v.BindEnv("viewer.token", "MUARA_VIEWER_TOKEN")
	_ = v.BindEnv("rate_limit.enabled", "MUARA_RATE_LIMIT_ENABLED")
	_ = v.BindEnv("rate_limit.requests_per_minute", "MUARA_RATE_LIMIT_REQUESTS_PER_MINUTE")
	_ = v.BindEnv("hardened", "MUARA_HARDENED")
}

// bindProviderEnvVars binds known provider config keys to optional env vars.
func bindProviderEnvVars(v *viper.Viper) {
	_ = v.BindEnv("providers.stripe.config.publishable_key", "MUARA_STRIPE_PUBLISHABLE_KEY")
	_ = v.BindEnv("providers.stripe.config.secret_key", "MUARA_STRIPE_SECRET_KEY")
	_ = v.BindEnv("providers.stripe.config.webhook_secret", "MUARA_STRIPE_WEBHOOK_SECRET")
}

// setDefaults registers safe offline defaults.
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", 9000)
	v.SetDefault("server.admin_port", "")
	v.SetDefault("server.public_base_url", "")
	v.SetDefault("server.tls_cert", "")
	v.SetDefault("server.tls_key", "")
	v.SetDefault("server.cors.allowed_origins", []string{"*"})
	v.SetDefault("server.cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("server.cors.allowed_headers", []string{"Content-Type", "Authorization", "X-CSRF-Token", "X-Request-ID"})
	v.SetDefault("server.cors.allow_credentials", false)
	v.SetDefault("server.csrf.enabled", true)
	v.SetDefault("server.pprof", false)
	v.SetDefault("admin.enabled", false)
	v.SetDefault("admin.username", "")
	v.SetDefault("admin.password_hash", "")
	v.SetDefault("admin.token", "")
	v.SetDefault("viewer.enabled", false)
	v.SetDefault("viewer.username", "")
	v.SetDefault("viewer.password_hash", "")
	v.SetDefault("viewer.token", "")
	v.SetDefault("rate_limit.enabled", false)
	v.SetDefault("rate_limit.requests_per_minute", 100)
	v.SetDefault("hardened", false)
	v.SetDefault("log.level", "info")
	v.SetDefault("persistence.type", "sqlite")
	v.SetDefault("persistence.path", ".muara/data/ledger.db")
	v.SetDefault("webhook.url", "")
	v.SetDefault("webhook.max_retries", 3)
	v.SetDefault("webhook.targets", map[string]string{})
	v.SetDefault("webhook.events", map[string][]string{})
	v.SetDefault("dev.seed", false)
	v.SetDefault("providers.stripe.enabled", false)
	v.SetDefault("providers.stripe.config.publishable_key", "pk_test_muara")
	v.SetDefault("providers.stripe.config.secret_key", "sk_test_muara")
	v.SetDefault("providers.stripe.config.webhook_secret", "whsec_muara")
	v.SetDefault("providers.senangpay.enabled", false)
	v.SetDefault("providers.senangpay.config.secret_key", "muara-senangpay-secret")
	v.SetDefault("providers.billplz.enabled", false)
	v.SetDefault("providers.billplz.config.api_key", "muara-billplz-api-key")
	v.SetDefault("providers.billplz.config.x_signature_key", "muara-billplz-x-signature")
	v.SetDefault("providers.billplz.config.collection_id", "")
	v.SetDefault("providers.toyyibpay.enabled", false)
	v.SetDefault("providers.toyyibpay.config.user_secret_key", "muara-toyyibpay-secret")
	v.SetDefault("providers.toyyibpay.config.category_code", "")
	v.SetDefault("providers.ipay88.enabled", false)
	v.SetDefault("providers.ipay88.config.merchant_code", "muara-ipay88-merchant")
	v.SetDefault("providers.ipay88.config.merchant_key", "muara-ipay88-key")
}

// DefaultYAML returns the bundled default configuration as YAML bytes.
func DefaultYAML() []byte {
	return []byte(`server:
  host: 127.0.0.1
  port: 9000
  # Optional second port for the admin UI; when set, provider endpoints stay on port.
  # admin_port: 9001
  # External URL used for payment links when behind a reverse proxy or tunnel.
  public_base_url: ""
  # Optional separate external URL for the admin UI. Use this when the admin
  # dashboard is served on a different subdomain (e.g. admin.muara.example.com)
  # instead of a different port on the same host.
  admin_public_base_url: ""
  tls_cert: ""
  tls_key: ""
  cors:
    allowed_origins:
      - "*"
    allowed_methods:
      - GET
      - POST
      - PUT
      - PATCH
      - DELETE
      - OPTIONS
    allowed_headers:
      - Content-Type
      - Authorization
      - X-CSRF-Token
      - X-Request-ID
    allow_credentials: false
  csrf:
    enabled: true
  pprof: false

# Admin authentication for /_admin and admin JSON APIs only.
# Provider emulation endpoints remain unauthenticated.
admin:
  enabled: false
  username: ""
  password_hash: ""
  token: ""

# Read-only dashboard access for testers. Viewers can inspect the ledger,
# transactions, and webhook log, but cannot change providers or webhooks.
viewer:
  enabled: false
  username: ""
  password_hash: ""
  token: ""

# In-memory rate limiting. No external dependencies.
rate_limit:
  enabled: false
  requests_per_minute: 100

# Hardened mode enables admin auth, rate limiting, and strict security headers.
# You must configure admin credentials when hardened is true.
hardened: false

log:
  level: info

persistence:
  type: sqlite
  path: .muara/data/ledger.db

# Generic provider configuration. Each key must match a registered provider name.
providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
      # Fawry supports v1 (legacy payload) and v2 (server notification).
      # Defaults to v1. Use v2 to emulate the Fawry V2 notification format.
      version: v1
  default:
    enabled: true
    config: {}
  stripe:
    enabled: false
    config:
      publishable_key: pk_test_muara
      secret_key: sk_test_muara
      webhook_secret: whsec_muara
  senangpay:
    enabled: false
    config:
      secret_key: muara-senangpay-secret
  billplz:
    enabled: false
    config:
      api_key: muara-billplz-api-key
      x_signature_key: muara-billplz-x-signature
      collection_id: ""
  toyyibpay:
    enabled: false
    config:
      user_secret_key: muara-toyyibpay-secret
      category_code: ""
  ipay88:
    enabled: false
    config:
      merchant_code: muara-ipay88-merchant
      merchant_key: muara-ipay88-key

webhook:
  url: ""
  max_retries: 3
`)
}

// ValidateWebhookURL returns a warning message if the webhook URL is invalid or empty.
// An empty URL is allowed; the dispatcher will simply skip deliveries.
func ValidateWebhookURL(cfg WebhookConfig) string {
	if cfg.URL == "" {
		return "webhook.url is empty; outgoing webhooks will not be dispatched"
	}
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return fmt.Sprintf("webhook.url is invalid: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Sprintf("webhook.url has unsupported scheme %q; must be http or https", u.Scheme)
	}
	return ""
}

// Load reads configuration from the given path and environment variables.
// If path is empty or the file does not exist, it falls back to defaults.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("MUARA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindSecurityEnvVars(v)
	bindProviderEnvVars(v)
	_ = v.BindEnv("disable_update_check", "MUARA_NO_UPDATE_CHECK")
	_ = v.BindEnv("dev.seed", "MUARA_DEV_SEED")
	setDefaults(v)

	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("read config: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.normalizeProviders()

	return &cfg, nil
}

// LoadFromBytes reads configuration from YAML bytes and environment variables.
// It applies the same defaults and validation preparation as Load.
func LoadFromBytes(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("MUARA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindSecurityEnvVars(v)
	bindProviderEnvVars(v)
	_ = v.BindEnv("disable_update_check", "MUARA_NO_UPDATE_CHECK")
	_ = v.BindEnv("dev.seed", "MUARA_DEV_SEED")
	setDefaults(v)

	if len(data) > 0 {
		if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.normalizeProviders()

	return &cfg, nil
}

// Validate returns an error if the configuration is inconsistent or unsupported.
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port)
	}
	if c.Server.AdminPort != 0 {
		if c.Server.AdminPort < 1 || c.Server.AdminPort > 65535 {
			return fmt.Errorf("server.admin_port must be between 1 and 65535, got %d", c.Server.AdminPort)
		}
		if c.Server.AdminPort == c.Server.Port {
			return fmt.Errorf("server.admin_port must be different from server.port")
		}
	}
	if c.Server.PublicBaseURL != "" {
		u, err := url.Parse(c.Server.PublicBaseURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("server.public_base_url must be a valid http or https URL, got %q", c.Server.PublicBaseURL)
		}
		if u.Host == "" {
			return fmt.Errorf("server.public_base_url must include a host")
		}
	}

	if err := c.validateSecurity(); err != nil {
		return err
	}

	switch c.Persistence.Type {
	case "memory", "sqlite", "":
		// supported
	default:
		return fmt.Errorf("persistence.type %q is not supported", c.Persistence.Type)
	}

	if c.Persistence.Type == "sqlite" && c.Persistence.Path == "" {
		return fmt.Errorf("persistence.path is required when persistence.type is sqlite")
	}

	for name, pc := range c.Providers {
		if !pc.Enabled {
			continue
		}
		p, err := providerForName(name)
		if err != nil {
			return fmt.Errorf("providers.%s is enabled but not registered", name)
		}
		if err := p.Init(pc.Config); err != nil {
			return fmt.Errorf("providers.%s config invalid: %w", name, err)
		}
	}

	return nil
}

// providerForName resolves a provider by checking the default registry, the Go
// factory registry, and then built-in gateway.yml manifests. The returned
// provider is not initialized.
func providerForName(name string) (provider.Provider, error) {
	if p, err := provider.Get(name); err == nil {
		return p, nil
	}
	if f, ok := factory.Get(name); ok {
		return f(nil)
	}
	plugins, err := plugin.LoadBuiltin("plugins", "../plugins", "../../plugins")
	if err != nil {
		return nil, err
	}
	for _, lp := range plugins {
		if lp.Name == name {
			return ProviderFromGateway(lp)
		}
	}
	return nil, fmt.Errorf("provider %q not found", name)
}

// validateSecurity checks admin, TLS, and hardened-mode settings.
func (c *Config) validateSecurity() error {
	hasTLSCert := c.Server.TLSCert != ""
	hasTLSKey := c.Server.TLSKey != ""
	if hasTLSCert != hasTLSKey {
		return fmt.Errorf("server.tls_cert and server.tls_key must both be set or both be empty")
	}

	if c.Admin.Enabled {
		if c.Admin.Username == "" {
			return fmt.Errorf("admin.username is required when admin.enabled is true")
		}
		if c.Admin.PasswordHash == "" && c.Admin.Token == "" {
			return fmt.Errorf("admin.password_hash or admin.token is required when admin.enabled is true")
		}
	}

	if c.Viewer.Enabled {
		if c.Viewer.Username == "" {
			return fmt.Errorf("viewer.username is required when viewer.enabled is true")
		}
		if c.Viewer.PasswordHash == "" && c.Viewer.Token == "" {
			return fmt.Errorf("viewer.password_hash or admin.token is required when viewer.enabled is true")
		}
		if c.Admin.Enabled {
			if c.Viewer.Username != "" && c.Viewer.Username == c.Admin.Username {
				return fmt.Errorf("viewer.username must be different from admin.username")
			}
			if c.Viewer.Token != "" && c.Viewer.Token == c.Admin.Token {
				return fmt.Errorf("viewer.token must be different from admin.token")
			}
		}
	}

	if c.Hardened {
		if !c.Admin.Enabled {
			return fmt.Errorf("hardened mode requires admin.enabled to be true")
		}
		if c.Admin.PasswordHash == "" && c.Admin.Token == "" {
			return fmt.Errorf("hardened mode requires admin.password_hash or admin.token")
		}
	}

	if c.RateLimit.Enabled && c.RateLimit.RequestsPerMinute <= 0 {
		return fmt.Errorf("rate_limit.requests_per_minute must be greater than 0")
	}

	return nil
}

// normalizeProviders maps legacy cfg.Fawry values into the providers map
// if providers.fawry is not already present.
func (c *Config) normalizeProviders() {
	if c.Providers == nil {
		c.Providers = make(map[string]ProviderConfig)
	}

	// Migrate legacy top-level fawry.* keys into providers.fawry.
	if _, ok := c.Providers["fawry"]; !ok && (c.Fawry.MerchantCode != "" || c.Fawry.MerchantSecurityKey != "" || c.Fawry.WebhookSecret != "") {
		slog.Warn("legacy fawry.* config keys are deprecated; use providers.fawry instead")
		c.Providers["fawry"] = ProviderConfig{
			Enabled: true,
			Config: map[string]any{
				"merchant_code":         c.Fawry.MerchantCode,
				"merchant_security_key": c.Fawry.MerchantSecurityKey,
				"webhook_secret":        c.Fawry.WebhookSecret,
				"version":               "v1",
			},
		}
	}

	// Apply built-in defaults for providers that ship with OpenMuara.
	if _, ok := c.Providers["fawry"]; !ok {
		c.Providers["fawry"] = ProviderConfig{
			Enabled: true,
			// #nosec G101 -- safe offline placeholder credentials, not real secrets
			Config: map[string]any{
				"merchant_code":         "muara-merchant-code",
				"merchant_security_key": "muara-fawry-secret",
				"webhook_secret":        "muara-webhook-secret",
				"version":               "v1",
			},
		}
	}
	if _, ok := c.Providers["default"]; !ok {
		c.Providers["default"] = ProviderConfig{Enabled: true, Config: map[string]any{}}
	}
}
