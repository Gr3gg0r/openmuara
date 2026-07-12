package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/config"
	"github.com/openmuara/openmuara/internal/plugin"
	"github.com/openmuara/openmuara/internal/provider"
	"github.com/openmuara/openmuara/internal/version"
)

// toolCheck describes a required or optional binary.
type toolCheck struct {
	name     string
	optional bool
}

var doctorTools = []toolCheck{
	{name: "go"},
	{name: "golangci-lint"},
	{name: "govulncheck", optional: true},
	{name: "task"},
}

// lookPathFunc matches exec.LookPath for test injection.
type lookPathFunc func(string) (string, error)

// doctorResult is the stable JSON schema returned by muara doctor --json.
type doctorResult struct {
	Healthy   bool             `json:"healthy"`
	Timestamp string           `json:"timestamp"`
	Tools     []doctorTool     `json:"tools"`
	Config    doctorConfig     `json:"config"`
	Providers []doctorProvider `json:"providers"`
	Webhook   doctorWebhook    `json:"webhook"`
	Version   doctorVersion    `json:"version"`
}

// doctorTool reports whether a single binary is present.
type doctorTool struct {
	Name     string `json:"name"`
	Found    bool   `json:"found"`
	Path     string `json:"path,omitempty"`
	Optional bool   `json:"optional"`
}

// doctorConfig reports config load and validation status.
type doctorConfig struct {
	OK     bool                     `json:"ok"`
	Path   string                   `json:"path"`
	Error  string                   `json:"error,omitempty"`
	Valid  bool                     `json:"valid"`
	Errors []config.ValidationError `json:"errors,omitempty"`
}

// doctorProvider reports readiness for one configured provider.
type doctorProvider struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Ready   bool   `json:"ready"`
	Error   string `json:"error,omitempty"`
}

// doctorWebhook reports whether a webhook URL is configured and reachable.
type doctorWebhook struct {
	URL        string `json:"url"`
	Configured bool   `json:"configured"`
	Reachable  *bool  `json:"reachable,omitempty"`
	Error      string `json:"error,omitempty"`
}

// doctorVersion embeds version metadata and optional update information.
type doctorVersion struct {
	Current         string `json:"current"`
	Latest          string `json:"latest,omitempty"`
	UpdateAvailable bool   `json:"update_available"`
}

func newDoctorCommand() *cobra.Command {
	var checkWebhook bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the muara environment and config",
		Example: `  muara doctor
  muara doctor --config /path/to/config.yml
  muara doctor --json
  muara doctor --check-webhook`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			result := runDoctor(cmd, checkWebhook, exec.LookPath)
			if jsonOutput {
				return respondJSON(cmd.OutOrStdout(), result)
			}
			printDoctorResult(cmd.OutOrStdout(), result)
			if !result.Healthy {
				return fmt.Errorf("doctor found problems")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&checkWebhook, "check-webhook", false, "attempt to reach the configured webhook URL")
	return cmd
}

// runDoctor collects all health signals and returns a typed result.
func runDoctor(cmd *cobra.Command, checkWebhook bool, lookPath lookPathFunc) doctorResult {
	result := doctorResult{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Config:    doctorConfig{Path: rootConfigPath, OK: true},
		Version: doctorVersion{
			Current: version.String(),
		},
	}

	result.Tools = checkTools(lookPath)
	toolsHealthy := true
	for _, t := range result.Tools {
		if !t.Found && !t.Optional {
			toolsHealthy = false
			break
		}
	}

	cfg, loadErr := config.Load(rootConfigPath)
	if loadErr != nil {
		result.Config.OK = false
		result.Config.Error = loadErr.Error()
		result.Healthy = false
		return result
	}

	valErrs := cfg.ValidateWithDetails(rootConfigPath)
	result.Config.Valid = len(valErrs) == 0
	result.Config.Errors = valErrs

	result.Providers = checkProviders(cfg)
	result.Webhook = checkWebhookStatus(cmd.Context(), cfg.Webhook, checkWebhook)

	latest, updateAvailable := maybeCheckUpdate(cmd)
	result.Version.Latest = latest
	result.Version.UpdateAvailable = updateAvailable

	providersHealthy := true
	for _, p := range result.Providers {
		if p.Enabled && !p.Ready {
			providersHealthy = false
			break
		}
	}

	result.Healthy = toolsHealthy && result.Config.Valid && providersHealthy
	return result
}

func checkTools(lookPath lookPathFunc) []doctorTool {
	tools := make([]doctorTool, 0, len(doctorTools))
	for _, tool := range doctorTools {
		path, err := lookPath(tool.name)
		tools = append(tools, doctorTool{
			Name:     tool.name,
			Found:    err == nil,
			Path:     path,
			Optional: tool.optional,
		})
	}
	return tools
}

func checkProviders(cfg *config.Config) []doctorProvider {
	registry := provider.Default()
	plugins, _ := plugin.LoadBuiltin("plugins", "../plugins", "../../plugins")
	pluginByName := make(map[string]*plugin.LoadedPlugin, len(plugins))
	for _, lp := range plugins {
		pluginByName[lp.Name] = lp
	}

	providers := make([]doctorProvider, 0, len(cfg.Providers))
	for name, pc := range cfg.Providers {
		dp := doctorProvider{Name: name, Enabled: pc.Enabled}
		if !pc.Enabled {
			providers = append(providers, dp)
			continue
		}
		p, err := registry.Get(name)
		if err != nil {
			if lp, ok := pluginByName[name]; ok {
				p, err = config.ProviderFromGateway(lp)
			}
		}
		if err != nil || p == nil {
			dp.Error = err.Error()
			providers = append(providers, dp)
			continue
		}
		configCopy := make(map[string]any, len(pc.Config))
		for k, v := range pc.Config {
			configCopy[k] = v
		}
		if err := p.Init(configCopy); err != nil {
			dp.Error = err.Error()
		} else {
			dp.Ready = true
		}
		providers = append(providers, dp)
	}
	return providers
}

func checkWebhookStatus(ctx context.Context, wc config.WebhookConfig, checkReachability bool) doctorWebhook {
	result := doctorWebhook{URL: wc.URL}
	if wc.URL == "" {
		return result
	}
	result.Configured = true
	if !checkReachability {
		return result
	}
	reachable := isWebhookReachable(ctx, wc.URL)
	result.Reachable = &reachable
	if !reachable {
		result.Error = "webhook URL did not respond within the timeout"
	}
	return result
}

func isWebhookReachable(ctx context.Context, url string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode < http.StatusInternalServerError
}

func printDoctorResult(w io.Writer, r doctorResult) {
	for _, t := range r.Tools {
		status := "PASS"
		if !t.Found {
			status = "FAIL"
		}
		optionalNote := ""
		if t.Optional {
			optionalNote = " (optional)"
		}
		_, _ = fmt.Fprintf(w, "%s %s%s\n", status, t.Name, optionalNote)
	}

	if !r.Config.OK {
		_, _ = fmt.Fprintf(w, "\nconfig load error: %s\n", r.Config.Error)
		_, _ = fmt.Fprintln(w, "hint: run 'muara init' to create a config file")
	} else if !r.Config.Valid {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprint(w, config.FormatValidationErrors(r.Config.Errors))
	} else {
		_, _ = fmt.Fprintln(w, "config is valid")
	}

	for _, p := range r.Providers {
		status := "SKIP"
		if p.Enabled {
			if p.Ready {
				status = "READY"
			} else {
				status = "FAIL"
			}
		}
		_, _ = fmt.Fprintf(w, "provider %s: %s\n", p.Name, status)
		if p.Error != "" {
			_, _ = fmt.Fprintf(w, "  error: %s\n", p.Error)
		}
	}

	if r.Webhook.Configured {
		if r.Webhook.Reachable == nil {
			_, _ = fmt.Fprintf(w, "webhook: configured (%s, reachability not checked)\n", r.Webhook.URL)
		} else if *r.Webhook.Reachable {
			_, _ = fmt.Fprintf(w, "webhook: reachable (%s)\n", r.Webhook.URL)
		} else {
			_, _ = fmt.Fprintf(w, "webhook: unreachable (%s)\n", r.Webhook.URL)
		}
	} else {
		_, _ = fmt.Fprintln(w, "webhook: not configured")
	}

	if r.Version.UpdateAvailable && !quietOutput {
		_, _ = fmt.Fprintf(w, "version: %s (newer version available: %s)\n", r.Version.Current, r.Version.Latest)
	}
}
