package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newScenarioCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scenario",
		Short: "Simulate payment outcomes for a transaction reference",
		Example: `  muara scenario success tx-123
  muara scenario fail tx-123`,
	}

	cmd.AddCommand(newScenarioSubCommand("success", "Mark a transaction as successfully paid"))
	cmd.AddCommand(newScenarioSubCommand("fail", "Mark a transaction as failed/unpaid"))
	cmd.AddCommand(newScenarioSubCommand("timeout", "Mark a transaction as timed out (unpaid after a short delay)"))

	return cmd
}

func newScenarioSubCommand(outcome, short string) *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("%s <ref>", outcome),
		Short:   short,
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("  muara scenario %s tx-123", outcome),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/_admin/scenario/%s?ref=%s", outcome, url.QueryEscape(args[0]))
			return postJSON(cmd, path)
		},
	}
}
