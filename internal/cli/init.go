package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/openmuara/openmuara/internal/config"
)

// promptFunc reads a line of input and returns it trimmed.
type promptFunc func(string) (string, error)

// defaultPromptFunc reads from stdin using the provided reader.
func defaultPromptFunc(r io.Reader) promptFunc {
	scanner := bufio.NewScanner(r)
	return func(question string) (string, error) {
		_, _ = fmt.Fprint(os.Stdout, question)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", io.EOF
		}
		return strings.TrimSpace(scanner.Text()), nil
	}
}

func newInitCommand() *cobra.Command {
	var useDefaults bool
	var dryRun bool
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a local muara workspace",
		Example: `  # Interactive wizard (default in a terminal)
  muara init

  # Non-interactive default config for CI or scripts
  muara init --defaults

  # Preview the generated config without writing files
  muara init --dry-run

  # Overwrite an existing workspace
  muara init --force`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			workspace := filepath.Dir(rootConfigPath)
			if workspace == "" || workspace == "." {
				workspace = ".muara"
			}

			if err := os.MkdirAll(workspace, 0o750); err != nil {
				return fmt.Errorf("create workspace: %w", err)
			}

			configPath := filepath.Join(workspace, "config.yml")

			if !dryRun && !force {
				if _, err := os.Stat(configPath); err == nil {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "config already exists:", configPath)
					return nil
				}
			}

			var cfgBytes []byte
			if useDefaults || !isInteractive(cmd.InOrStdin()) {
				cfgBytes = config.DefaultYAML()
			} else {
				choices, webhookURL, logLevel, err := runWizard(cmd.OutOrStdout(), defaultPromptFunc(cmd.InOrStdin()))
				if err != nil {
					return fmt.Errorf("wizard: %w", err)
				}
				cfg := config.GenerateWizardConfig(choices, webhookURL, logLevel)
				cfgBytes = config.RenderWizardConfig(cfg)
			}

			if dryRun {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), string(cfgBytes))
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}

			if err := os.WriteFile(configPath, cfgBytes, 0o600); err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "created", configPath)

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	cmd.Flags().BoolVar(&useDefaults, "defaults", false, "write the default config without interactive prompts")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print the generated config without writing files")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing config file")
	return cmd
}

// isInteractive returns true if stdin is a terminal.
func isInteractive(stdin io.Reader) bool {
	if f, ok := stdin.(*os.File); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return false
}

// runWizard asks the user questions and returns the chosen providers, webhook URL, and log level.
func runWizard(out io.Writer, prompt promptFunc) ([]config.WizardChoice, string, string, error) {
	templates := config.WizardTemplates()

	_, _ = fmt.Fprintln(out, "Welcome to OpenMuara. Let's set up your local payment emulator.")
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Available payment providers (select by number, comma-separated):")
	for i, c := range templates.Choices {
		marker := " "
		if c.IsRecommended {
			marker = "*"
		}
		_, _ = fmt.Fprintf(out, "  %d%s %s — %s\n", i+1, marker, c.DisplayName, c.Description)
	}

	defaultIdx := 0
	for i, c := range templates.Choices {
		if c.IsRecommended {
			defaultIdx = i
			break
		}
	}

	answer, err := prompt(fmt.Sprintf("Providers [%d]: ", defaultIdx+1))
	if err != nil {
		return nil, "", "", err
	}

	indices := []int{defaultIdx}
	if answer != "" {
		indices = nil
		for _, part := range strings.Split(answer, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			n, err := strconv.Atoi(part)
			if err != nil || n < 1 || n > len(templates.Choices) {
				// Accept provider key as input too.
				found := false
				for i, c := range templates.Choices {
					if strings.EqualFold(c.Key, part) || strings.EqualFold(c.DisplayName, part) {
						indices = appendUniqIndex(indices, i)
						found = true
						break
					}
				}
				if !found {
					_, _ = fmt.Fprintf(out, "ignoring unknown selection: %q\n", part)
				}
				continue
			}
			indices = appendUniqIndex(indices, n-1)
		}
		if len(indices) == 0 {
			indices = []int{defaultIdx}
		}
	}

	choices := make([]config.WizardChoice, 0, len(indices))
	for _, idx := range indices {
		choices = append(choices, templates.Choices[idx])
	}

	webhookURL, err := prompt("Webhook URL for outgoing notifications (optional): ")
	if err != nil {
		return nil, "", "", err
	}

	logLevel, err := prompt("Log level [info]: ")
	if err != nil {
		return nil, "", "", err
	}
	if logLevel == "" {
		logLevel = "info"
	}

	return choices, webhookURL, logLevel, nil
}

func appendUniqIndex(indices []int, idx int) []int {
	for _, existing := range indices {
		if existing == idx {
			return indices
		}
	}
	return append(indices, idx)
}
