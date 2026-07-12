package config

import (
	"os"
	"path/filepath"
	"testing"
)

func FuzzLoad(f *testing.F) {
	f.Add(DefaultYAML())
	f.Add([]byte("not: valid: yaml: ["))

	f.Fuzz(func(t *testing.T, data []byte) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yml")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}

		cfg, err := Load(path)
		if err != nil {
			// Errors are fine; panics are not.
			return
		}

		_ = cfg.Validate()
		_ = ValidateWebhookURL(cfg.Webhook)
	})
}
