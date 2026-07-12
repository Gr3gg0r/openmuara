package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/config"
)

func newCleanCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Reset local muara data",
		Long: `Remove the local SQLite database used for the ledger, audit log, and
webhook attempts.

This is useful when you want to start from a blank slate. For memory persistence,
use the dashboard or restart the server instead; this command only removes files
on disk.`,
		Example: `  muara clean
  muara clean --force`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runClean(cmd.OutOrStdout(), cmd.InOrStdin(), force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")
	return cmd
}

func runClean(out io.Writer, in io.Reader, force bool) error {
	cfg, err := config.Load(rootConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	switch cfg.Persistence.Type {
	case "memory", "":
		_, _ = fmt.Fprintln(out, "persistence is in-memory; nothing to clean on disk")
		return nil
	case "sqlite":
		return cleanSQLite(out, in, cfg.Persistence.Path, force)
	default:
		return fmt.Errorf("unsupported persistence type %q", cfg.Persistence.Type)
	}
}

func cleanSQLite(out io.Writer, in io.Reader, path string, force bool) error {
	if path == "" {
		path = ".muara/data/ledger.db"
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}

	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		_, _ = fmt.Fprintf(out, "nothing to clean: %s does not exist\n", abs)
		return nil
	}
	if err != nil {
		return fmt.Errorf("check data file: %w", err)
	}

	if !force {
		_, _ = fmt.Fprintf(out, "This will permanently delete:\n  %s\n\nType 'yes' to continue: ", abs)
		reader := bufio.NewReader(in)
		answer, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("read confirmation: %w", readErr)
		}
		if answer != "yes\n" {
			_, _ = fmt.Fprintln(out, "clean cancelled")
			return nil
		}
	}

	if err := os.Remove(abs); err != nil {
		return fmt.Errorf("remove %s: %w", abs, err)
	}

	_, _ = fmt.Fprintf(out, "cleaned: %s\n", abs)
	return nil
}
