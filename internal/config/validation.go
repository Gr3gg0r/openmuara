// Package config provides structured config validation with actionable errors.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/provider/factory"
	"gopkg.in/yaml.v3"
)

// ValidationError is a single actionable config problem.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Hint    string `json:"hint"`
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
}

// Error implements the error interface.
func (v ValidationError) Error() string {
	if v.Line > 0 {
		return fmt.Sprintf("%s:%d: %s: %s (%s)", v.File, v.Line, v.Field, v.Message, v.Hint)
	}
	if v.File != "" {
		return fmt.Sprintf("%s: %s: %s (%s)", v.File, v.Field, v.Message, v.Hint)
	}
	return fmt.Sprintf("%s: %s (%s)", v.Field, v.Message, v.Hint)
}

// ValidateWithDetails runs Validate and returns a slice of structured errors.
// The file path is used for best-effort line number lookups.
func (c *Config) ValidateWithDetails(filePath string) []ValidationError {
	var errs []ValidationError
	lineMap := make(map[string]int)
	if filePath != "" {
		lineMap = fieldLineMap(filePath)
	}

	add := func(field, message, hint string) {
		errs = append(errs, ValidationError{
			Field:   field,
			Message: message,
			Hint:    hint,
			File:    filePath,
			Line:    lineMap[field],
		})
	}

	if c.Server.Host == "" {
		add("server.host", "server host is required", "set server.host to 127.0.0.1 or 0.0.0.0")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		add("server.port", fmt.Sprintf("server.port must be between 1 and 65535, got %d", c.Server.Port), "set server.port to a valid TCP port")
	}
	if c.Server.PublicBaseURL != "" {
		if u, err := url.Parse(c.Server.PublicBaseURL); err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			add("server.public_base_url", fmt.Sprintf("server.public_base_url %q is not a valid HTTP/HTTPS URL with a host", c.Server.PublicBaseURL), "set it to the external URL testers use, e.g. https://muara.example.com, or leave it empty")
		}
	}

	switch c.Persistence.Type {
	case "memory", "sqlite", "":
		// supported
	default:
		add("persistence.type", fmt.Sprintf("persistence.type %q is not supported", c.Persistence.Type), "use sqlite or memory")
	}

	if (c.Persistence.Type == "sqlite" || c.Persistence.Type == "") && c.Persistence.Path == "" {
		add("persistence.path", "persistence.path is required when persistence.type is sqlite", "set persistence.path to a file path like .muara/data/ledger.db")
	}

	plugins, _ := plugin.LoadBuiltin("plugins", "../plugins", "../../plugins")
	pluginByName := make(map[string]*plugin.LoadedPlugin, len(plugins))
	for _, lp := range plugins {
		pluginByName[lp.Name] = lp
	}

	for name, pc := range c.Providers {
		field := "providers." + name

		p, err := provider.Get(name)
		if err != nil {
			if lp, ok := pluginByName[name]; ok {
				p, err = ProviderFromGateway(lp)
			} else if f, ok := factory.Get(name); ok {
				p, err = f(nil)
			}
		}
		if err != nil || p == nil {
			add(field, fmt.Sprintf("provider %q is not registered", name), "check the provider name or install the matching plugin")
			continue
		}
		if !pc.Enabled {
			continue
		}
		if err := p.Init(pc.Config); err != nil {
			add(field+".config", fmt.Sprintf("provider %q config is invalid: %v", name, err), "compare your config against muara.yml.example for this provider")
		}
		if name == "fawry" {
			if v, ok := pc.Config["version"].(string); ok && v != "v1" && v != "v2" {
				add(field+".config.version", fmt.Sprintf("fawry version %q is not supported", v), "use v1 (legacy) or v2 (server notification)")
			}
		}
	}

	if c.Webhook.URL != "" {
		if u, err := url.Parse(c.Webhook.URL); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			add("webhook.url", fmt.Sprintf("webhook.url %q is not a valid HTTP/HTTPS URL", c.Webhook.URL), "set webhook.url to http:// or https:// endpoint, or leave it empty")
		}
	}

	if (c.Fawry.MerchantCode != "" || c.Fawry.MerchantSecurityKey != "" || c.Fawry.WebhookSecret != "") && len(c.Providers["fawry"].Config) == 0 {
		add("fawry", "legacy top-level fawry.* keys are deprecated", "move your Fawry config under providers.fawry.config")
	}

	c.validateSecurityDetails(add)

	return errs
}

// validateSecurityDetails adds structured validation details for security settings.
func (c *Config) validateSecurityDetails(add func(field, message, hint string)) {
	hasTLSCert := c.Server.TLSCert != ""
	hasTLSKey := c.Server.TLSKey != ""
	if hasTLSCert != hasTLSKey {
		add("server.tls_cert", "server.tls_cert and server.tls_key must both be set or both be empty", "provide both paths to enable TLS, or remove both to disable TLS")
	}

	if c.Admin.Enabled {
		if c.Admin.Username == "" {
			add("admin.username", "admin.username is required when admin authentication is enabled", "set admin.username or disable admin.enabled")
		}
		if c.Admin.PasswordHash == "" && c.Admin.Token == "" {
			add("admin.password_hash", "admin.password_hash or admin.token is required when admin authentication is enabled", "set admin.password_hash (bcrypt) or admin.token, or disable admin.enabled")
		}
	}

	if c.Viewer.Enabled {
		if c.Viewer.Username == "" {
			add("viewer.username", "viewer.username is required when viewer authentication is enabled", "set viewer.username or disable viewer.enabled")
		}
		if c.Viewer.PasswordHash == "" && c.Viewer.Token == "" {
			add("viewer.password_hash", "viewer.password_hash or viewer.token is required when viewer authentication is enabled", "set viewer.password_hash (bcrypt) or viewer.token, or disable viewer.enabled")
		}
		if c.Admin.Enabled {
			if c.Viewer.Username != "" && c.Viewer.Username == c.Admin.Username {
				add("viewer.username", "viewer.username must be different from admin.username", "choose a different username for the viewer account")
			}
			if c.Viewer.Token != "" && c.Viewer.Token == c.Admin.Token {
				add("viewer.token", "viewer.token must be different from admin.token", "choose a different token for the viewer account")
			}
		}
	}

	if c.Hardened {
		if !c.Admin.Enabled {
			add("hardened", "hardened mode requires admin authentication to be enabled", "set admin.enabled to true and configure admin credentials")
		}
		if c.Admin.PasswordHash == "" && c.Admin.Token == "" {
			add("hardened", "hardened mode requires admin credentials", "set admin.password_hash or admin.token")
		}
	}

	if c.RateLimit.Enabled && c.RateLimit.RequestsPerMinute <= 0 {
		add("rate_limit.requests_per_minute", "rate_limit.requests_per_minute must be greater than 0", "set a positive value such as 100")
	}

	if c.Server.Host == "0.0.0.0" && !c.Admin.Enabled && !c.Hardened {
		add("server.host", "server is bound to 0.0.0.0 without admin authentication", "bind to 127.0.0.1 or enable admin authentication / hardened mode")
	}
}

// fieldLineMap parses a YAML file and returns the first line number for each dotted field path.
func fieldLineMap(path string) map[string]int {
	result := make(map[string]int)
	// #nosec G304 -- reads already-loaded config file path for line mapping
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return result
	}
	if len(root.Content) == 0 {
		return result
	}
	walkYAML("", root.Content[0], result)
	return result
}

func walkYAML(prefix string, node *yaml.Node, out map[string]int) {
	if node.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		if _, ok := out[path]; !ok {
			out[path] = keyNode.Line
		}
		if valueNode.Kind == yaml.MappingNode {
			walkYAML(path, valueNode, out)
		}
	}
}

// FormatValidationErrors returns a human-readable string for a slice of errors.
func FormatValidationErrors(errs []ValidationError) string {
	if len(errs) == 0 {
		return ""
	}
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "config validation failed with %d error(s):\n", len(errs))
	for _, e := range errs {
		loc := e.Field
		if e.Line > 0 {
			loc = fmt.Sprintf("%s (line %d)", e.Field, e.Line)
		}
		_, _ = fmt.Fprintf(&b, "  - %s: %s\n    hint: %s\n", loc, e.Message, e.Hint)
	}
	return b.String()
}
