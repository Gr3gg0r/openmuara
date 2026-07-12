package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidationError aggregates schema validation failures.
type ValidationError struct {
	Issues []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("plugin validation failed (%d issues): %s", len(v.Issues), strings.Join(v.Issues, "; "))
}

// Addf appends a formatted issue to the validation error.
func (v *ValidationError) Addf(format string, args ...any) {
	v.Issues = append(v.Issues, fmt.Sprintf(format, args...))
}

var (
	validMethods = map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}
	validSigAlgs = map[string]bool{
		"fawry_sha256":        true,
		"hmac_sha256":         true,
		"md5_concat":          true,
		"senangpay_md5":       true,
		"billplz_hmac_sha256": true,
		"ipay88_sha256":       true,
		"toyyibpay_md5":       true,
		"stripe_v1":           true,
		"none":                true,
	}
	validRuntimeTypes = map[string]bool{"simple": true, "go": true, "hybrid": true}
	versionRe         = regexp.MustCompile(`^\d+\.\d+\.\d+`)
	nameRe            = regexp.MustCompile(`^[a-z0-9-]+$`)
)

// Validate checks a GatewayConfig against the canonical schema.
func Validate(cfg GatewayConfig) error {
	var ve ValidationError

	if cfg.SchemaVersion != CurrentSchemaVersion {
		ve.Addf("schema_version must be %q, got %q", CurrentSchemaVersion, cfg.SchemaVersion)
	}

	if cfg.Metadata.Name == "" {
		ve.Addf("metadata.name is required")
	} else if !nameRe.MatchString(cfg.Metadata.Name) {
		ve.Addf("metadata.name %q must match %s", cfg.Metadata.Name, nameRe.String())
	}

	if cfg.Metadata.Version == "" {
		ve.Addf("metadata.version is required")
	} else if !versionRe.MatchString(cfg.Metadata.Version) {
		ve.Addf("metadata.version %q is not semver-like", cfg.Metadata.Version)
	}

	for i, r := range cfg.Routes {
		if r.Path == "" {
			ve.Addf("routes[%d].path is required", i)
		}
		if !validMethods[r.Method] {
			ve.Addf("routes[%d].method %q is invalid", i, r.Method)
		}
		if r.Action == "" {
			ve.Addf("routes[%d].action is required", i)
		}
		if r.Runtime != "" && !validRuntimeTypes[r.Runtime] {
			ve.Addf("routes[%d].runtime %q is not supported", i, r.Runtime)
		}
		if r.SchemaRef != "" {
			if _, ok := cfg.Schemas.Requests[r.SchemaRef]; !ok {
				if _, ok := cfg.Schemas.Responses[r.SchemaRef]; !ok {
					ve.Addf("routes[%d].schema_ref %q not found in schemas", i, r.SchemaRef)
				}
			}
		}
	}

	if cfg.Signature != nil {
		if !validSigAlgs[cfg.Signature.Algorithm] {
			ve.Addf("signature.algorithm %q is not supported", cfg.Signature.Algorithm)
		}
		if cfg.Signature.SecretKey == "" {
			ve.Addf("signature.secret_key is required when signature is present")
		} else if !strings.Contains(cfg.Signature.SecretKey, ".") {
			ve.Addf("signature.secret_key %q must be a dotted config path", cfg.Signature.SecretKey)
		}
	}

	if cfg.Runtime != nil {
		if cfg.Runtime.Type != "" && !validRuntimeTypes[cfg.Runtime.Type] {
			ve.Addf("runtime.type %q is not supported", cfg.Runtime.Type)
		}
		if (cfg.Runtime.Type == "" || cfg.Runtime.Type == "simple") && cfg.Runtime.Simple != nil {
			sr := cfg.Runtime.Simple
			if sr.Currency == "" {
				ve.Addf("runtime.simple.currency is required")
			}
			if sr.ReferenceField == "" {
				ve.Addf("runtime.simple.reference_field is required")
			}
			if sr.AmountField == "" {
				ve.Addf("runtime.simple.amount_field is required")
			}
			if sr.EscapePage != nil && sr.EscapePage.Enabled {
				if sr.EscapePage.Route == "" {
					ve.Addf("runtime.simple.escape_page.route is required when enabled")
				}
				if sr.EscapePage.RefParam == "" {
					ve.Addf("runtime.simple.escape_page.ref_param is required when enabled")
				}
			}
		}
	}

	if len(ve.Issues) > 0 {
		return &ve
	}
	return nil
}

// IsValidationError reports whether err is a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
