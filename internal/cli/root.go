// Package cli implements the muara command-line interface using Cobra.
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/Gr3gg0r/openmuara/internal/version"
	"github.com/spf13/cobra"
)

// rootConfigPath holds the value of the persistent --config flag.
var rootConfigPath string

// jsonOutput and quietOutput are persistent flags set on the root command.
var jsonOutput bool
var quietOutput bool

// newRootCommand builds the top-level cobra command.
func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "muara",
		Short:         "Local-first payment virtualization layer",
		Long:          "OpenMuara emulates payment providers for local development and testing.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `  muara --help
  muara --quiet doctor
  muara --json version`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if quietOutput {
				cmd.SetOut(io.Discard)
			}
			return nil
		},
	}

	root.Version = version.String()
	root.SetVersionTemplate("{{.Version}}\n")

	root.PersistentFlags().StringVar(&rootConfigPath, "config", ".muara/config.yml", "path to config file")
	root.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output results as JSON where supported")
	root.PersistentFlags().BoolVar(&quietOutput, "quiet", false, "suppress non-error output")

	root.AddCommand(
		newVersionCommand(),
		newInitCommand(),
		newStartCommand(),
		newDoctorCommand(),
		newHealthCommand(),
		newCompletionCommand(),
		newWebhookCommand(),
		newTransactionCommand(),
		newScenarioCommand(),
		newPluginsCommand(),
		newProviderCommand(),
		newAuditCommand(),
		newSecurityCommand(),
		newCleanCommand(),
	)

	return root
}

// Execute runs the root command and propagates context cancellation.
func Execute(ctx context.Context) error {
	root := newRootCommand()
	root.SetContext(ctx)

	if err := root.Execute(); err != nil {
		return fmt.Errorf("muara: %w", err)
	}
	return nil
}
