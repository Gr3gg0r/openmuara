// Package migrations embeds SQL schema migrations for OpenMuara persistence.
package migrations

import (
	"embed"
	"fmt"
)

//go:embed *.sql
var fs embed.FS

// Read returns the contents of the named migration file.
func Read(name string) (string, error) {
	b, err := fs.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("read migration %q: %w", name, err)
	}
	return string(b), nil
}
