package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Gr3gg0r/openmuara/internal/config"
	"github.com/Gr3gg0r/openmuara/internal/version"
)

const (
	releasesAPIURL = "https://api.github.com/repos/Gr3gg0r/openmuara/releases/latest"
	updateTimeout  = 2 * time.Second
)

// versionOutput is the structured output for muara version --json.
type versionOutput struct {
	Version         string `json:"version"`
	Commit          string `json:"commit"`
	BuildTime       string `json:"build_time"`
	Latest          string `json:"latest,omitempty"`
	UpdateAvailable bool   `json:"update_available"`
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print muara version",
		Example: `  muara version
  muara version --json
  muara version --quiet`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			latest, updateAvailable := maybeCheckUpdate(cmd)

			if jsonOutput {
				out := versionOutput{
					Version:         version.Version,
					Commit:          version.Commit,
					BuildTime:       version.BuildTime,
					Latest:          latest,
					UpdateAvailable: updateAvailable,
				}
				return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), version.String())
			if updateAvailable && !quietOutput {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "A newer version is available: %s\n", latest)
			}
			return nil
		},
	}
}

// maybeCheckUpdate returns the latest release tag and whether it is newer than
// the current binary. It never returns errors; network or parse failures are
// silently ignored so the CLI stays usable offline.
func maybeCheckUpdate(cmd *cobra.Command) (string, bool) {
	if quietOutput {
		return "", false
	}
	if version.IsDev() {
		return "", false
	}
	if os.Getenv("MUARA_NO_UPDATE_CHECK") != "" {
		return "", false
	}
	if cfg, err := config.Load(rootConfigPath); err == nil && cfg.DisableUpdateCheck {
		return "", false
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), updateTimeout)
	defer cancel()

	latest, err := fetchLatestRelease(ctx)
	if err != nil || latest == "" {
		return "", false
	}

	return latest, isNewerVersion(latest, version.Version)
}

// fetchLatestRelease returns the latest release tag from GitHub.
func fetchLatestRelease(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesAPIURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: updateTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return payload.TagName, nil
}

// isNewerVersion reports whether latest is a higher semantic version than current.
// It accepts optional "v" prefixes and ignores non-semver tags.
func isNewerVersion(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")
	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		ln, err1 := strconv.Atoi(latestParts[i])
		cn, err2 := strconv.Atoi(currentParts[i])
		if err1 != nil || err2 != nil {
			return false
		}
		if ln > cn {
			return true
		}
		if ln < cn {
			return false
		}
	}
	return len(latestParts) > len(currentParts)
}

func respondJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
