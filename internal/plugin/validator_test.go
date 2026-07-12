package plugin

import (
	"strings"
	"testing"
)

func validConfig() GatewayConfig {
	return GatewayConfig{
		SchemaVersion: CurrentSchemaVersion,
		Metadata: Metadata{
			Name:    "test-gateway",
			Version: "1.0.0",
		},
		Routes: []Route{
			{Path: "/test", Method: "POST", Action: "charge", SchemaRef: "charge_request"},
		},
		Schemas: Schemas{
			Requests: map[string]Schema{
				"charge_request": {Fields: []Field{{Name: "amount", Type: "integer", Required: true}}},
			},
		},
		Signature: &Signature{Algorithm: "fawry_sha256", SecretKey: "fawry.secret"},
	}
}

// Given a valid GatewayConfig, When Validate runs, Then it returns nil.
func TestValidateValidConfig(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// Given a GatewayConfig with the wrong schema_version, When Validate runs, Then it returns a ValidationError.
func TestValidateWrongSchemaVersion(t *testing.T) {
	cfg := validConfig()
	cfg.SchemaVersion = "v0"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if !strings.Contains(err.Error(), "schema_version") {
		t.Errorf("expected schema_version issue, got %q", err.Error())
	}
}

// Given a GatewayConfig with an invalid plugin name, When Validate runs, Then it returns a ValidationError.
func TestValidateInvalidName(t *testing.T) {
	cfg := validConfig()
	cfg.Metadata.Name = "Bad Name"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with a non-semver version, When Validate runs, Then it returns a ValidationError.
func TestValidateInvalidVersion(t *testing.T) {
	cfg := validConfig()
	cfg.Metadata.Version = "not-a-version"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with an invalid route method, When Validate runs, Then it returns a ValidationError.
func TestValidateInvalidMethod(t *testing.T) {
	cfg := validConfig()
	cfg.Routes[0].Method = "TRACE"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with a missing schema_ref, When Validate runs, Then it returns a ValidationError.
func TestValidateMissingSchemaRef(t *testing.T) {
	cfg := validConfig()
	cfg.Routes[0].SchemaRef = "missing"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with an unsupported signature algorithm, When Validate runs, Then it returns a ValidationError.
func TestValidateUnsupportedSignatureAlgorithm(t *testing.T) {
	cfg := validConfig()
	cfg.Signature.Algorithm = "unknown_alg"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with a non-dotted secret_key, When Validate runs, Then it returns a ValidationError.
func TestValidateDottedSecretKey(t *testing.T) {
	cfg := validConfig()
	cfg.Signature.SecretKey = "secret"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

// Given a GatewayConfig with multiple issues, When Validate runs, Then it returns all issues aggregated.
func TestValidateAggregatesIssues(t *testing.T) {
	cfg := validConfig()
	cfg.SchemaVersion = "v0"
	cfg.Metadata.Name = ""
	cfg.Routes[0].Method = "TRACE"

	err := Validate(cfg)
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 3 {
		t.Errorf("expected at least 3 issues, got %d: %v", len(ve.Issues), ve.Issues)
	}
}
