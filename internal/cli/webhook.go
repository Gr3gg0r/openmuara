package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/config"
)

func newWebhookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Inspect and replay outgoing webhooks",
		Example: `  muara webhook list
  muara webhook inspect tx-123
  muara webhook replay tx-123`,
	}

	cmd.AddCommand(newWebhookListCommand())
	cmd.AddCommand(newWebhookInspectCommand())
	cmd.AddCommand(newWebhookReplayCommand())

	return cmd
}

func newWebhookListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List recent webhook attempts",
		Example: "  muara webhook list",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return getJSON(cmd, "/_admin/webhooks")
		},
	}
}

func newWebhookInspectCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect <ref>",
		Short:   "Inspect a webhook attempt by reference",
		Args:    cobra.ExactArgs(1),
		Example: "  muara webhook inspect tx-123",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getJSON(cmd, "/_admin/webhooks/"+url.PathEscape(args[0]))
		},
	}
}

func newWebhookReplayCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "replay <ref>",
		Short:   "Replay a webhook by reference",
		Args:    cobra.ExactArgs(1),
		Example: "  muara webhook replay --ref tx-123\n  muara webhook replay tx-123",
		RunE: func(cmd *cobra.Command, args []string) error {
			return postJSON(cmd, "/_admin/webhooks/"+url.PathEscape(args[0])+"/replay")
		},
	}
}

func getJSON(cmd *cobra.Command, path string) error {
	baseURL, err := serverBaseURL(cmd)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + path)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return printResponse(resp, cmd.OutOrStdout())
}

func postJSON(cmd *cobra.Command, path string) error {
	baseURL, err := serverBaseURL(cmd)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(baseURL+path, "application/json", nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return printResponse(resp, cmd.OutOrStdout())
}

func serverBaseURL(_ *cobra.Command) (string, error) {
	cfg, err := config.Load(rootConfigPath)
	if err != nil {
		return "", fmt.Errorf("load config: %w", err)
	}

	return fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port), nil
}

func printResponse(resp *http.Response, out io.Writer) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var pretty any
	if err := json.Unmarshal(body, &pretty); err == nil {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(pretty)
	} else {
		_, _ = out.Write(body)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	return nil
}
