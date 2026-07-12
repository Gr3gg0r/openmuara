package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newTransactionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transaction",
		Short: "Inspect transactions",
		Example: `  muara transaction inspect tx-123
  muara transaction list`,
	}

	cmd.AddCommand(newTransactionInspectCommand())
	cmd.AddCommand(newTransactionListCommand())

	return cmd
}

func newTransactionListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List transactions",
		Example: "  muara transaction list",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return getJSON(cmd, "/_admin/transactions")
		},
	}
}

func newTransactionInspectCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect <ref>",
		Short:   "Inspect a transaction by reference",
		Args:    cobra.ExactArgs(1),
		Example: "  muara transaction inspect tx-123",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getJSON(cmd, "/_admin/transactions/"+url.PathEscape(args[0]))
		},
	}
}
