// Package config provides wizard-driven configuration generation.
package config

import (
	"fmt"
	"sort"
	"strings"
)

// WizardChoice represents one provider a user can pick during muara init.
type WizardChoice struct {
	Key            string
	DisplayName    string
	Description    string
	IsRecommended  bool
	SampleRoute    string
	SampleMethod   string
	EnvVarKeys     []string
	ProviderConfig ProviderConfig
}

// EnvVarName returns the canonical MUARA_<PROVIDER>_<KEY> environment variable
// name for a provider config key.
func EnvVarName(provider, key string) string {
	return "MUARA_" + strings.ToUpper(provider+"_"+key)
}

// ProviderNextStep gives a human-readable next action for a provider.
type ProviderNextStep struct {
	Method string `json:"method"`
	Route  string `json:"route"`
	Hint   string `json:"hint"`
}

// WizardPrompts holds the interactive prompts used by muara init.
type WizardPrompts struct {
	Choices []WizardChoice
}

// NextStep returns a suggested first API call for this provider.
func (c WizardChoice) NextStep() ProviderNextStep {
	return ProviderNextStep{
		Method: c.SampleMethod,
		Route:  c.SampleRoute,
		Hint:   fmt.Sprintf("Send %s %s to create your first charge.", c.SampleMethod, c.SampleRoute),
	}
}

// WizardTemplates returns the provider choices available in the first-run wizard.
func WizardTemplates() WizardPrompts {
	return WizardPrompts{
		Choices: []WizardChoice{
			{
				Key:           "fawry",
				DisplayName:   "Fawry",
				Description:   "Egyptian payment gateway (legacy v1 and V2 server notifications)",
				IsRecommended: true,
				SampleRoute:   "/fawry/charge",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"merchant_code", "merchant_security_key", "webhook_secret"},
				ProviderConfig: ProviderConfig{
					Enabled: true,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"merchant_code":         "muara-merchant-code",
						"merchant_security_key": "muara-fawry-secret",
						"webhook_secret":        "muara-webhook-secret",
						"version":               "v1",
					},
				},
			},
			{
				Key:           "stripe",
				DisplayName:   "Stripe",
				Description:   "Stripe Checkout Sessions and PaymentIntents emulation",
				IsRecommended: false,
				SampleRoute:   "/v1/checkout/sessions",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"publishable_key", "secret_key", "webhook_secret"},
				ProviderConfig: ProviderConfig{
					Enabled: false,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"publishable_key": "pk_test_muara",
						"secret_key":      "sk_test_muara",
						"webhook_secret":  "whsec_muara",
					},
				},
			},
			{
				Key:           "billplz",
				DisplayName:   "Billplz",
				Description:   "Malaysian payment gateway for collections and bills",
				IsRecommended: false,
				SampleRoute:   "/api/v3/bills",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"api_key", "x_signature_key", "collection_id"},
				ProviderConfig: ProviderConfig{
					Enabled: false,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"api_key":         "muara-billplz-api-key",
						"x_signature_key": "muara-billplz-x-signature",
						"collection_id":   "",
					},
				},
			},
			{
				Key:           "toyyibpay",
				DisplayName:   "ToyyibPay",
				Description:   "Malaysian payment gateway for categories and bills",
				IsRecommended: false,
				SampleRoute:   "/index.php/api/createBill",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"user_secret_key", "category_code"},
				ProviderConfig: ProviderConfig{
					Enabled: false,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"user_secret_key": "muara-toyyibpay-secret",
						"category_code":   "",
					},
				},
			},
			{
				Key:           "ipay88",
				DisplayName:   "iPay88",
				Description:   "Southeast Asian payment gateway with redirect flow",
				IsRecommended: false,
				SampleRoute:   "/ePayment/entry.asp",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"merchant_code", "merchant_key"},
				ProviderConfig: ProviderConfig{
					Enabled: false,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"merchant_code": "muara-ipay88-merchant",
						"merchant_key":  "muara-ipay88-key",
					},
				},
			},
			{
				Key:           "senangpay",
				DisplayName:   "SenangPay",
				Description:   "Malaysian payment gateway with signature verification",
				IsRecommended: false,
				SampleRoute:   "/senangpay/payment",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{"secret_key"},
				ProviderConfig: ProviderConfig{
					Enabled: false,
					// #nosec G101 -- wizard sample dummy credentials
					Config: map[string]any{
						"secret_key": "muara-senangpay-secret",
					},
				},
			},
			{
				Key:           "default",
				DisplayName:   "Default / DIY",
				Description:   "Minimal provider for custom experiments",
				IsRecommended: false,
				SampleRoute:   "/default/charge",
				SampleMethod:  "POST",
				EnvVarKeys:    []string{},
				ProviderConfig: ProviderConfig{
					Enabled: true,
					Config:  map[string]any{},
				},
			},
		},
	}
}

// WizardChoiceByKey returns the wizard choice for the given provider key.
func WizardChoiceByKey(key string) (WizardChoice, bool) {
	for _, c := range WizardTemplates().Choices {
		if c.Key == key {
			return c, true
		}
	}
	return WizardChoice{}, false
}

// GenerateWizardConfig builds a Config from wizard answers. The first selected
// provider is marked as the active/recommended one; all unselected providers are
// included but disabled so users can enable them later without re-running init.
func GenerateWizardConfig(choices []WizardChoice, webhookURL string, logLevel string) *Config {
	cfg := &Config{}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 9000
	cfg.Server.CORS.AllowedOrigins = []string{"*"}
	cfg.Server.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	cfg.Server.CORS.AllowedHeaders = []string{"Content-Type", "Authorization", "X-CSRF-Token", "X-Request-ID"}
	cfg.Server.CORS.AllowCredentials = false
	cfg.Server.CSRF.Enabled = true
	cfg.Log.Level = logLevel
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	cfg.Persistence.Type = "sqlite"
	cfg.Persistence.Path = ".muara/data/ledger.db"

	selected := make(map[string]bool, len(choices))
	for _, c := range choices {
		selected[c.Key] = true
	}

	cfg.Providers = make(map[string]ProviderConfig)
	for _, c := range WizardTemplates().Choices {
		pc := c.ProviderConfig
		pc.Enabled = selected[c.Key]
		cfg.Providers[c.Key] = pc
	}

	cfg.Webhook.URL = webhookURL
	cfg.Webhook.MaxRetries = 3
	return cfg
}

// RenderWizardConfig renders a wizard-generated config as YAML bytes.
// It emits every documented top-level section so the result is complete and
// valid without requiring a second editing pass.
func RenderWizardConfig(cfg *Config) []byte {
	var b strings.Builder
	b.WriteString("# OpenMuara configuration generated by muara init.\n")
	b.WriteString("# Edit providers.<name>.config with your test credentials.\n")
	b.WriteString("server:\n")
	_, _ = fmt.Fprintf(&b, "  host: %s\n", cfg.Server.Host)
	_, _ = fmt.Fprintf(&b, "  port: %d\n", cfg.Server.Port)
	_, _ = fmt.Fprintf(&b, "  tls_cert: \"%s\"\n", cfg.Server.TLSCert)
	_, _ = fmt.Fprintf(&b, "  tls_key: \"%s\"\n", cfg.Server.TLSKey)
	b.WriteString("  cors:\n")
	b.WriteString("    allowed_origins:\n")
	for _, o := range cfg.Server.CORS.AllowedOrigins {
		_, _ = fmt.Fprintf(&b, "      - \"%s\"\n", o)
	}
	b.WriteString("    allowed_methods:\n")
	for _, m := range cfg.Server.CORS.AllowedMethods {
		_, _ = fmt.Fprintf(&b, "      - %s\n", m)
	}
	b.WriteString("    allowed_headers:\n")
	for _, h := range cfg.Server.CORS.AllowedHeaders {
		_, _ = fmt.Fprintf(&b, "      - %s\n", h)
	}
	_, _ = fmt.Fprintf(&b, "    allow_credentials: %t\n", cfg.Server.CORS.AllowCredentials)
	b.WriteString("  csrf:\n")
	_, _ = fmt.Fprintf(&b, "    enabled: %t\n", cfg.Server.CSRF.Enabled)
	_, _ = fmt.Fprintf(&b, "  pprof: %t\n", cfg.Server.Pprof)

	b.WriteString("\n# Admin authentication for /_admin and admin JSON APIs only.\n")
	b.WriteString("# Provider emulation endpoints remain unauthenticated.\n")
	b.WriteString("admin:\n")
	_, _ = fmt.Fprintf(&b, "  enabled: %t\n", cfg.Admin.Enabled)
	_, _ = fmt.Fprintf(&b, "  username: \"%s\"\n", cfg.Admin.Username)
	_, _ = fmt.Fprintf(&b, "  password_hash: \"%s\"\n", cfg.Admin.PasswordHash)
	_, _ = fmt.Fprintf(&b, "  token: \"%s\"\n", cfg.Admin.Token)

	b.WriteString("\n# In-memory rate limiting. No external dependencies.\n")
	b.WriteString("rate_limit:\n")
	_, _ = fmt.Fprintf(&b, "  enabled: %t\n", cfg.RateLimit.Enabled)
	_, _ = fmt.Fprintf(&b, "  requests_per_minute: %d\n", cfg.RateLimit.RequestsPerMinute)

	b.WriteString("\n# Hardened mode enables admin auth, rate limiting, and strict security headers.\n")
	b.WriteString("# You must configure admin credentials when hardened is true.\n")
	_, _ = fmt.Fprintf(&b, "hardened: %t\n", cfg.Hardened)

	b.WriteString("\nlog:\n")
	_, _ = fmt.Fprintf(&b, "  level: %s\n", cfg.Log.Level)
	b.WriteString("\npersistence:\n")
	_, _ = fmt.Fprintf(&b, "  type: %s\n", cfg.Persistence.Type)
	_, _ = fmt.Fprintf(&b, "  path: %s\n", cfg.Persistence.Path)
	b.WriteString("\nproviders:\n")
	for _, name := range sortedProviderKeys(cfg.Providers) {
		pc := cfg.Providers[name]
		_, _ = fmt.Fprintf(&b, "  %s:\n", name)
		_, _ = fmt.Fprintf(&b, "    enabled: %t\n", pc.Enabled)
		b.WriteString("    config:\n")
		for _, k := range sortedStringKeys(pc.Config) {
			_, _ = fmt.Fprintf(&b, "      %s: %v\n", k, pc.Config[k])
		}
	}
	b.WriteString("\nwebhook:\n")
	_, _ = fmt.Fprintf(&b, "  url: \"%s\"\n", cfg.Webhook.URL)
	_, _ = fmt.Fprintf(&b, "  max_retries: %d\n", cfg.Webhook.MaxRetries)
	if len(cfg.Webhook.Targets) > 0 {
		b.WriteString("  targets:\n")
		for _, k := range sortedStringMapKeys(cfg.Webhook.Targets) {
			_, _ = fmt.Fprintf(&b, "    %s: \"%s\"\n", k, cfg.Webhook.Targets[k])
		}
	}
	return []byte(b.String())
}

func sortedProviderKeys(m map[string]ProviderConfig) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return sortedStrings(keys)
}

func sortedStringKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return sortedStrings(keys)
}

func sortedStringMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return sortedStrings(keys)
}

func sortedStrings(s []string) []string {
	sort.Strings(s)
	return s
}
