// Package plugin defines the declarative gateway.yml schema and runtime
// support for OpenMuara payment provider plugins.
package plugin

// CurrentSchemaVersion is the supported gateway.yml schema version.
const CurrentSchemaVersion = "v1"

// GatewayConfig is the top-level gateway.yml structure.
type GatewayConfig struct {
	SchemaVersion string     `mapstructure:"schema_version" yaml:"schema_version"`
	Metadata      Metadata   `mapstructure:"metadata" yaml:"metadata"`
	Runtime       *Runtime   `mapstructure:"runtime" yaml:"runtime"`
	Routes        []Route    `mapstructure:"routes" yaml:"routes"`
	Schemas       Schemas    `mapstructure:"schemas" yaml:"schemas"`
	Signature     *Signature `mapstructure:"signature" yaml:"signature"`
	Webhooks      []Webhook  `mapstructure:"webhooks" yaml:"webhooks"`
	Fixtures      []Fixture  `mapstructure:"fixtures" yaml:"fixtures"`
}

// Runtime selects the runtime implementation for this gateway.
type Runtime struct {
	Type        string         `mapstructure:"type" yaml:"type"`
	Description string         `mapstructure:"description" yaml:"description"`
	Simple      *SimpleRuntime `mapstructure:"simple" yaml:"simple"`
}

// SimpleRuntime configures the built-in simple provider runtime.
// It is sufficient for gateways that need a single charge endpoint,
// signature verification, a templated response, an optional escape page,
// and an optional outgoing webhook event.
type SimpleRuntime struct {
	// ChargeRoute identifies which route action handles the charge request.
	ChargeRoute string `mapstructure:"charge_route" yaml:"charge_route"`
	// StatusRoute identifies which route action handles status queries.
	StatusRoute string `mapstructure:"status_route" yaml:"status_route"`
	// WebhookEvent is the event name dispatched after a successful escape action.
	WebhookEvent string `mapstructure:"webhook_event" yaml:"webhook_event"`
	// ResponseTemplate is the JSON body returned by the charge handler.
	// Values such as {{ .Reference }}, {{ .Status }}, and {{ .Amount }} are rendered.
	ResponseTemplate map[string]any `mapstructure:"response_template" yaml:"response_template"`
	// EscapePage configures the payment simulation page.
	EscapePage *EscapePageConfig `mapstructure:"escape_page" yaml:"escape_page"`
	// Currency is the default currency recorded for transactions.
	Currency string `mapstructure:"currency" yaml:"currency"`
	// ReferenceField is the JSON field used as the transaction reference.
	ReferenceField string `mapstructure:"reference_field" yaml:"reference_field"`
	// AmountField is the JSON field used as the transaction amount.
	AmountField string `mapstructure:"amount_field" yaml:"amount_field"`
	// CustomerField is the optional JSON field used as the customer reference.
	CustomerField string `mapstructure:"customer_field" yaml:"customer_field"`
}

// EscapePageConfig configures the generic escape/simulation page.
type EscapePageConfig struct {
	Enabled     bool   `mapstructure:"enabled" yaml:"enabled"`
	Route       string `mapstructure:"route" yaml:"route"`
	RefParam    string `mapstructure:"ref_param" yaml:"ref_param"`
	ReturnParam string `mapstructure:"return_param" yaml:"return_param"`
	AmountParam string `mapstructure:"amount_param" yaml:"amount_param"`
	StatusParam string `mapstructure:"status_param" yaml:"status_param"`
	RedirectURL string `mapstructure:"redirect_url" yaml:"redirect_url"`
}

// Metadata describes the plugin.
type Metadata struct {
	Name              string   `mapstructure:"name" yaml:"name"`
	Version           string   `mapstructure:"version" yaml:"version"`
	Description       string   `mapstructure:"description" yaml:"description"`
	Author            string   `mapstructure:"author" yaml:"author"`
	Tags              []string `mapstructure:"tags" yaml:"tags"`
	SupportedVersions []string `mapstructure:"supported_versions" yaml:"supported_versions"`
}

// Route maps an HTTP pattern to a handler action.
type Route struct {
	Path        string `mapstructure:"path" yaml:"path"`
	Method      string `mapstructure:"method" yaml:"method"`
	Action      string `mapstructure:"action" yaml:"action"`
	Runtime     string `mapstructure:"runtime" yaml:"runtime"`
	Description string `mapstructure:"description" yaml:"description"`
	SchemaRef   string `mapstructure:"schema_ref" yaml:"schema_ref"`
}

// Schemas holds named request/response validation rules.
type Schemas struct {
	Requests  map[string]Schema `mapstructure:"requests" yaml:"requests"`
	Responses map[string]Schema `mapstructure:"responses" yaml:"responses"`
}

// Schema is a lightweight field-level contract.
type Schema struct {
	Fields []Field `mapstructure:"fields" yaml:"fields"`
}

// Field describes one expected field.
type Field struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Type     string `mapstructure:"type" yaml:"type"`
	Required bool   `mapstructure:"required" yaml:"required"`
	JSONName string `mapstructure:"json_name" yaml:"json_name"`
}

// Signature describes how to verify or generate a provider signature.
type Signature struct {
	Algorithm string   `mapstructure:"algorithm" yaml:"algorithm"`
	Fields    []string `mapstructure:"fields" yaml:"fields"`
	SecretEnv string   `mapstructure:"secret_env" yaml:"secret_env"`
	SecretKey string   `mapstructure:"secret_key" yaml:"secret_key"`
}

// Webhook describes an outgoing provider notification.
type Webhook struct {
	Name     string         `mapstructure:"name" yaml:"name"`
	Event    string         `mapstructure:"event" yaml:"event"`
	Method   string         `mapstructure:"method" yaml:"method"`
	Template map[string]any `mapstructure:"template" yaml:"template"`
}

// Fixture is a sample request/response pair for tests.
type Fixture struct {
	Name     string         `mapstructure:"name" yaml:"name"`
	RouteRef string         `mapstructure:"route_ref" yaml:"route_ref"`
	Request  map[string]any `mapstructure:"request" yaml:"request"`
	Response map[string]any `mapstructure:"response" yaml:"response"`
}
