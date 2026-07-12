package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/provider/simple"
)

func newPluginsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "Manage muara provider plugins",
		Long:  "List and validate declarative provider plugins.",
		Example: `  muara plugins list
  muara plugins validate`,
	}
	cmd.AddCommand(
		newPluginsListCommand(),
		newPluginsValidateCommand(),
	)
	return cmd
}

func newPluginsListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List discovered plugins",
		Example: "  muara plugins list",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPluginsList(cmd.OutOrStdout())
		},
	}
}

func runPluginsList(w io.Writer) error {
	builtins, err := plugin.LoadBuiltin("plugins")
	if err != nil {
		return fmt.Errorf("load built-in plugins: %w", err)
	}

	locals, err := plugin.LoadLocal(".muara/plugins")
	if err != nil {
		return fmt.Errorf("load local plugins: %w", err)
	}

	// Local overrides replace built-ins.
	byName := make(map[string]*plugin.LoadedPlugin)
	for _, p := range builtins {
		byName[p.Name] = p
	}
	for _, p := range locals {
		byName[p.Name] = p
	}

	if len(byName) == 0 {
		_, _ = fmt.Fprintln(w, "No plugins discovered.")
		return nil
	}

	_, _ = fmt.Fprintf(w, "%-15s %-10s %-30s %s\n", "NAME", "VERSION", "DESCRIPTION", "STATUS")
	for _, p := range byName {
		status := "ok"
		if err := plugin.Validate(p.Config); err != nil {
			status = "invalid"
		}
		desc := p.Config.Metadata.Description
		if desc == "" {
			desc = "-"
		}
		_, _ = fmt.Fprintf(w, "%-15s %-10s %-30s %s\n", p.Name, p.Config.Metadata.Version, desc, status)
	}
	return nil
}

func newPluginsValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "validate [path]",
		Short:   "Validate plugin YAML",
		Long:    "Validate a plugin directory or all discovered plugins if no path is given.",
		Args:    cobra.MaximumNArgs(1),
		Example: "  muara plugins validate\n  muara plugins validate ./my-plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginsValidate(cmd.OutOrStdout(), args)
		},
	}
}

func runPluginsValidate(w io.Writer, args []string) error {
	var plugins []*plugin.LoadedPlugin

	if len(args) == 1 {
		p, err := loadSinglePlugin(args[0])
		if err != nil {
			return err
		}
		plugins = append(plugins, p)
	} else {
		builtins, err := plugin.LoadBuiltin("plugins")
		if err != nil {
			return fmt.Errorf("load built-in plugins: %w", err)
		}
		locals, err := plugin.LoadLocal(".muara/plugins")
		if err != nil {
			return fmt.Errorf("load local plugins: %w", err)
		}
		plugins = append(plugins, builtins...)
		plugins = append(plugins, locals...)
	}

	if len(plugins) == 0 {
		_, _ = fmt.Fprintln(w, "No plugins to validate.")
		return nil
	}

	failed := false
	for _, p := range plugins {
		if err := plugin.Validate(p.Config); err != nil {
			failed = true
			_, _ = fmt.Fprintf(w, "INVALID %s: %v\n", p.Name, err)
			continue
		}
		_, _ = fmt.Fprintf(w, "OK %s\n", p.Name)
	}

	if failed {
		return fmt.Errorf("plugin validation failed")
	}
	return nil
}

func loadSinglePlugin(dir string) (*plugin.LoadedPlugin, error) {
	plugins, err := plugin.LoadBuiltin(dir)
	if err != nil {
		return nil, err
	}
	if len(plugins) == 0 {
		return nil, fmt.Errorf("no plugin found in %q", dir)
	}
	return plugins[0], nil
}

func newProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Manage simple providers",
		Long:  "Test and scaffold YAML-driven simple providers.",
		Example: `  muara provider test fawry
  muara provider init my-gateway`,
	}
	cmd.AddCommand(
		newProviderTestCommand(),
		newProviderInitCommand(),
	)
	return cmd
}

func newProviderTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "test <name>",
		Short:   "Test a simple provider charge endpoint",
		Long:    "Load a simple provider from gateway.yml and exercise its charge handler with a fixture.",
		Args:    cobra.ExactArgs(1),
		Example: "  muara provider test fawry",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProviderTest(cmd.OutOrStdout(), "plugins", args[0])
		},
	}
}

func runProviderTest(w io.Writer, pluginsDir, name string) error {
	plugins, err := plugin.LoadBuiltin(pluginsDir)
	if err != nil {
		return fmt.Errorf("load plugins: %w", err)
	}

	var lp *plugin.LoadedPlugin
	for _, p := range plugins {
		if p.Name == name {
			lp = p
			break
		}
	}
	if lp == nil {
		return fmt.Errorf("provider %q not found in plugins/", name)
	}
	if lp.Config.Runtime == nil || lp.Config.Runtime.Type != "simple" {
		return fmt.Errorf("provider %q does not use runtime.type=simple", name)
	}
	if err := plugin.Validate(lp.Config); err != nil {
		return fmt.Errorf("validate provider %q: %w", name, err)
	}

	cfg := map[string]any{}
	if lp.Config.Signature != nil && lp.Config.Signature.SecretKey != "" {
		cfg[secretConfigKey(lp.Config.Signature.SecretKey)] = "muara-secret"
	}

	// Test mode disables signature verification so fixtures with placeholder
	// signatures can still exercise the charge handler.
	testCfg := lp.Config
	if testCfg.Signature != nil {
		testCfg.Signature = &plugin.Signature{
			Algorithm: "none",
			SecretKey: lp.Config.Signature.SecretKey,
		}
	}

	p := simple.NewProvider(testCfg)
	if err := p.Init(cfg); err != nil {
		return fmt.Errorf("init provider %q: %w", name, err)
	}

	chargePath := routePathForAction(lp.Config.Routes, lp.Config.Runtime.Simple.ChargeRoute)
	if chargePath == "" {
		return fmt.Errorf("no route path found for charge action %q", lp.Config.Runtime.Simple.ChargeRoute)
	}

	fixture := findFixture(lp.Config, chargePath)
	if fixture == nil {
		return fmt.Errorf("no fixture found for charge route of provider %q", name)
	}

	body, err := json.Marshal(fixture.Request)
	if err != nil {
		return fmt.Errorf("marshal fixture: %w", err)
	}

	req := httptest.NewRequest(http.MethodPost, chargePath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	p.ChargeHandler().ServeHTTP(rec, req)

	_, _ = fmt.Fprintf(w, "HTTP %d\n%s\n", rec.Code, rec.Body.String())
	if rec.Code != http.StatusOK {
		return fmt.Errorf("charge handler returned status %d", rec.Code)
	}
	return nil
}

func secretConfigKey(dotted string) string {
	parts := filepath.SplitList(dotted)
	if len(parts) == 0 {
		return dotted
	}
	return parts[len(parts)-1]
}

func routePathForAction(routes []plugin.Route, action string) string {
	for _, r := range routes {
		if r.Action == action {
			return r.Path
		}
	}
	return ""
}

func findFixture(cfg plugin.GatewayConfig, route string) *plugin.Fixture {
	for i := range cfg.Fixtures {
		if cfg.Fixtures[i].RouteRef == route {
			return &cfg.Fixtures[i]
		}
	}
	return nil
}

func newProviderInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "init <name>",
		Short:   "Scaffold a new simple provider",
		Long:    "Create a new provider directory with a starter gateway.yml.",
		Args:    cobra.ExactArgs(1),
		Example: "  muara provider init my-gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProviderInit(cmd.OutOrStdout(), args[0])
		},
	}
}

func runProviderInit(w io.Writer, name string) error {
	dir := filepath.Join("plugins", name)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create directory %q: %w", dir, err)
	}
	path := filepath.Join(dir, "gateway.yml")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("gateway.yml already exists at %q", path)
	}

	// #nosec G306 -- starter manifest is non-sensitive
	if err := os.WriteFile(path, []byte(providerStarterYAML(name)), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	_, _ = fmt.Fprintf(w, "Created %s\n", path)
	return nil
}

func providerStarterYAML(name string) string {
	return fmt.Sprintf(`schema_version: v1

metadata:
  name: %s
  version: 1.0.0
  description: %s provider
  author: openmuara

runtime:
  type: simple
  simple:
    charge_route: %s_charge
    webhook_event: payment.completed
    currency: MYR
    reference_field: order_id
    amount_field: amount
    response_template:
      status: "ok"
      reference: "{{ .Reference }}"

routes:
  - path: /%s/charge
    method: POST
    action: %s_charge
    description: Accept a signed charge request
    schema_ref: charge_request

schemas:
  requests:
    charge_request:
      fields:
        - name: order_id
          json_name: order_id
          type: string
          required: true
        - name: amount
          json_name: amount
          type: number
          required: true
        - name: signature
          json_name: signature
          type: string
          required: true

signature:
  algorithm: hmac_sha256
  fields:
    - order_id
    - amount
    - signature
  secret_key: providers.%s.config.secret_key

webhooks:
  - name: %s_payment
    event: payment.completed
    method: POST
    template:
      order_id: "{{ .Reference }}"
      status: "{{ .Status }}"

fixtures:
  - name: valid_charge
    route_ref: /%s/charge
    request:
      order_id: ORDER-1
      amount: 10.00
      signature: placeholder
    response:
      status: ok
      reference: ORDER-1
`, name, name, name, name, name, name, name, name)
}
