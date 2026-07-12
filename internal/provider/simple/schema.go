package simple

import (
	"fmt"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
)

func (p *Provider) validateRequest(values map[string]any, schemaRef string) error {
	if schemaRef == "" {
		return nil
	}
	schema, ok := p.cfg.Schemas.Requests[schemaRef]
	if !ok {
		return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("schema %q not found", schemaRef))
	}

	for _, f := range schema.Fields {
		if !f.Required {
			continue
		}
		name := f.JSONName
		if name == "" {
			name = f.Name
		}
		if _, ok := values[name]; !ok {
			return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("%s is required", name))
		}
	}
	return nil
}

func stringValue(values map[string]any, key string) (string, bool) {
	if key == "" {
		return "", false
	}
	v, ok := values[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func floatValue(values map[string]any, key string) (float64, bool) {
	if key == "" {
		return 0, false
	}
	v, ok := values[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func intValue(values map[string]any, key string) (int, bool) {
	if key == "" {
		return 0, false
	}
	v, ok := values[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return n, true
	case float64:
		return int(n), true
	case int64:
		return int(n), true
	}
	return 0, false
}
