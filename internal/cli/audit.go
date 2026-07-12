package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/audit"
	"github.com/openmuara/openmuara/internal/config"
)

func newAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Inspect the audit log",
		Example: `  muara audit list
  muara audit list --since 2026-01-01T00:00:00Z`,
	}
	cmd.AddCommand(newAuditListCommand())
	return cmd
}

func openAuditStore(cfg config.PersistenceConfig) (audit.Store, error) {
	switch cfg.Type {
	case "memory", "":
		return audit.NewMemoryStore(), nil
	case "sqlite":
		path := cfg.Path
		if path == "" {
			path = ".muara/data/ledger.db"
		}
		if dir := filepath.Dir(path); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return nil, fmt.Errorf("create data dir: %w", err)
			}
		}
		return audit.NewSQLiteStore(path)
	default:
		return nil, fmt.Errorf("unsupported persistence type %q", cfg.Type)
	}
}

func newAuditListCommand() *cobra.Command {
	var limit int
	var since string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List recent audit events",
		Example: "  muara audit list --limit 10 --since 2026-01-01T00:00:00Z",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(rootConfigPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			store, err := openAuditStore(cfg.Persistence)
			if err != nil {
				return fmt.Errorf("init audit store: %w", err)
			}
			if closer, ok := store.(interface{ Close() error }); ok {
				defer func() { _ = closer.Close() }()
			}

			var sinceTime time.Time
			if since != "" {
				t, err := time.Parse(time.RFC3339, since)
				if err != nil {
					return fmt.Errorf("parse --since: %w", err)
				}
				sinceTime = t
			}

			events, err := store.ListSince(limit, 0, sinceTime)
			if err != nil {
				return fmt.Errorf("list audit events: %w", err)
			}

			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(events)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "maximum number of events to return")
	cmd.Flags().StringVar(&since, "since", "", "only return events at or after this RFC3339 timestamp")

	return cmd
}
