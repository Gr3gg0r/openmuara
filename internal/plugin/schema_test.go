package plugin

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func unmarshalStrict(data []byte, v any) error {
	dec := yaml.NewDecoder(&strictReader{data: data})
	dec.KnownFields(true)
	return dec.Decode(v)
}

type strictReader struct {
	data []byte
	pos  int
}

func (r *strictReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func TestCurrentSchemaVersion(t *testing.T) {
	if CurrentSchemaVersion != "v1" {
		t.Errorf("CurrentSchemaVersion: want v1, got %q", CurrentSchemaVersion)
	}
}

// Given a valid gateway.yml, When unmarshaled, Then all sections populate correctly.
func TestUnmarshalValidGatewayConfig(t *testing.T) {
	yamlDoc := `
schema_version: v1
metadata:
  name: test-gateway
  version: 1.0.0
  description: A test gateway
  author: openmuara
  tags:
    - test
routes:
  - path: /test/charge
    method: POST
    action: charge
    description: Charge endpoint
    schema_ref: charge_request
schemas:
  requests:
    charge_request:
      fields:
        - name: reference
          json_name: reference
          type: string
          required: true
signature:
  algorithm: test_hmac
  fields:
    - reference
  secret_key: test.secret
webhooks:
  - name: payment_completed
    event: payment.completed
    method: POST
    template:
      reference: "{{ .Reference }}"
fixtures:
  - name: valid_charge
    route_ref: /test/charge
    request:
      reference: ref-123
    response:
      status: ok
`

	var cfg GatewayConfig
	if err := yaml.Unmarshal([]byte(yamlDoc), &cfg); err != nil {
		t.Fatalf("unmarshal valid config: %v", err)
	}

	if cfg.SchemaVersion != "v1" {
		t.Errorf("schema_version: want v1, got %q", cfg.SchemaVersion)
	}
	if cfg.Metadata.Name != "test-gateway" {
		t.Errorf("metadata.name: want test-gateway, got %q", cfg.Metadata.Name)
	}
	if len(cfg.Routes) != 1 {
		t.Fatalf("routes: want 1, got %d", len(cfg.Routes))
	}
	if cfg.Routes[0].Path != "/test/charge" {
		t.Errorf("route path: want /test/charge, got %q", cfg.Routes[0].Path)
	}
	if len(cfg.Schemas.Requests) != 1 {
		t.Fatalf("request schemas: want 1, got %d", len(cfg.Schemas.Requests))
	}
	if cfg.Signature == nil || cfg.Signature.Algorithm != "test_hmac" {
		t.Errorf("signature algorithm: want test_hmac, got %v", cfg.Signature)
	}
	if len(cfg.Webhooks) != 1 {
		t.Fatalf("webhooks: want 1, got %d", len(cfg.Webhooks))
	}
	if len(cfg.Fixtures) != 1 {
		t.Fatalf("fixtures: want 1, got %d", len(cfg.Fixtures))
	}
}

// Given a YAML with unknown top-level keys, When strict-unmarshaled, Then it returns an error.
func TestUnmarshalUnknownKeysStrict(t *testing.T) {
	yamlDoc := `
schema_version: v1
metadata:
  name: test-gateway
unknown_section:
  value: 1
`

	var cfg GatewayConfig
	if err := unmarshalStrict([]byte(yamlDoc), &cfg); err == nil {
		t.Fatal("expected strict unmarshal to error for unknown key, got nil")
	}
}

// loadPluginYAML is a small harness that loads a plugin manifest and returns its name.
func loadPluginYAML(t *testing.T, path string) *GatewayConfig {
	t.Helper()

	// #nosec G304 -- test helper loads fixture plugin YAML
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var cfg GatewayConfig
	if err := unmarshalStrict(data, &cfg); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
	return &cfg
}

func TestFawryGatewayYAMLLoads(t *testing.T) {
	path := filepath.Join("..", "..", "plugins", "fawry", "gateway.yml")
	cfg := loadPluginYAML(t, path)
	if cfg.Metadata.Name != "fawry" {
		t.Errorf("metadata.name: want fawry, got %q", cfg.Metadata.Name)
	}
	if cfg.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("schema_version: want %q, got %q", CurrentSchemaVersion, cfg.SchemaVersion)
	}
}

func TestStripeGatewayYAMLLoads(t *testing.T) {
	path := filepath.Join("..", "..", "plugins", "stripe", "gateway.yml")
	cfg := loadPluginYAML(t, path)
	if cfg.Metadata.Name != "stripe" {
		t.Errorf("metadata.name: want stripe, got %q", cfg.Metadata.Name)
	}
	if cfg.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("schema_version: want %q, got %q", CurrentSchemaVersion, cfg.SchemaVersion)
	}
}
