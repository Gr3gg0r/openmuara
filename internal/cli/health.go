package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/spf13/cobra"
)

const healthTimeout = 5 * time.Second

// healthOutput is the structured output for muara health --json.
type healthOutput struct {
	Healthy bool   `json:"healthy"`
	URL     string `json:"url"`
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
}

func newHealthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check whether the local muara server is healthy",
		Example: `  muara health
  muara health --config /path/to/config.yml
  muara health --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result := runHealth(rootConfigPath)
			if jsonOutput {
				return respondJSON(cmd.OutOrStdout(), result)
			}
			if !result.Healthy {
				if result.Error != "" {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "unhealthy: %s\n", result.Error)
				} else {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "unhealthy")
				}
				return fmt.Errorf("health check failed")
			}
			if !quietOutput {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "healthy (%s)\n", result.URL)
			}
			return nil
		},
	}
}

// runHealth queries the server's /healthz endpoint and returns a structured result.
func runHealth(configPath string) healthOutput {
	cfg, err := config.Load(configPath)
	if err != nil {
		// Fall back to defaults so the command works even when no config file exists.
		cfg, _ = config.LoadFromBytes(config.DefaultYAML())
	}

	host := cfg.Server.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := cfg.Server.Port
	if port <= 0 {
		port = 9000
	}

	url := "http://" + host + ":" + strconv.Itoa(port) + "/healthz"
	result := healthOutput{URL: url}

	client := &http.Client{Timeout: healthTimeout}
	resp, err := client.Get(url)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		result.Error = err.Error()
		return result
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		result.Error = fmt.Sprintf("invalid health response: %v", err)
		return result
	}

	result.Status = payload.Status
	result.Healthy = resp.StatusCode == http.StatusOK && payload.Status == "ok"
	if !result.Healthy {
		result.Error = fmt.Sprintf("status %q (HTTP %d)", payload.Status, resp.StatusCode)
	}
	return result
}
