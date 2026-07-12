> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara CLI & TUI Polish

> **Status:** ✅ Completed | **Started:** 2026-07-03 | **Completed:** 2026-07-03
> **Scope:** Bring the `muara` CLI up to the standard of modern developer tools like `gh`, `stripe`, `vercel`, and `flyctl`: rich interactive prompts, proper versioning, shell completions, clear errors, and a polished first-run experience.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Current gaps

1. `muara --version` / `muara -v` do not work.
2. `muara version` prints `dev (unknown) built unknown` because ldflags are not injected in dev builds.
3. `muara init` uses plain text prompts; it only supports one provider and has no visual hierarchy.
4. No shell completion generation.
5. No progress indicators for long operations.
6. Error messages are plain and do not suggest next steps.
7. `muara doctor` only checks for binaries; it does not validate config health or provider readiness.
8. No update check or version notification.
9. No `--dry-run` or config preview mode.
10. Output formatting is inconsistent across commands (some support `--json`, some do not).

---

## Inspiration from best-in-class CLIs

| Tool | Pattern to adopt |
|---|---|
| `gh` (GitHub CLI) | Rich forms, clear error hints, `gh completion` generation. |
| `stripe` | Clean `--version`, rich interactive setup, per-command `--api-version` hints. |
| `vercel` | Friendly first-run wizard, progress spinners, actionable error messages. |
| `flyctl` | `flyctl doctor` checks service health, not just binaries. |
| `ngrok` | Config validation preview, onboarding tips after first run. |
| `kubectl` | `kubectl completion`, structured output, consistent flag behavior. |

---

## Goals

### 1. Version and build metadata
- Support `muara --version` and `muara -v`.
- Inject `Version`, `Commit`, and `BuildTime` via ldflags for dev, CI, and release builds.
- Add `muara version --json` with structured fields.

### 2. Rich interactive init wizard
- Use a TUI form library (evaluate `charmbracelet/huh`) for the wizard.
- Multi-select providers with checkboxes.
- Per-provider config step (enter keys, secrets, version) only when a provider is selected.
- Webhook URL input with validation hint.
- Log level selector.
- Workspace path selector with default and validation.
- Final config preview with confirm/deny before writing.
- `--defaults` remains non-interactive.

### 3. Shell completions
- Add `muara completion bash|zsh|fish|powershell` command.
- Document one-line install in README.

### 4. Progress and status feedback
- Add spinners for `init`, `migrate`, and long-running operations.
- Use `charmbracelet/bubbles` spinner or a lightweight equivalent.

### 5. Actionable errors
- Wrap common errors with "Did you mean...?" hints.
- Suggest `muara doctor` when config is invalid.
- Suggest `muara init` when workspace is missing.

### 6. Enhanced `muara doctor`
- Check tool presence (existing).
- Check config loads and validates.
- Check provider readiness (keys configured, enabled providers valid).
- Check webhook URL reachability (optional, with timeout).
- Return structured JSON with `--json`.

### 7. Update notification
- On `version` or `doctor`, optionally check GitHub Releases for a newer version.
- Respect `--quiet` and a config flag to disable checks.
- No blocking network calls; fail silently offline.

### 8. Consistent output
- Every command that produces data supports `--json` and `--quiet` where meaningful.
- Table output uses aligned columns (evaluate `github.com/olekukonko/tablewriter` or a small custom formatter).

### 9. Dry-run and preview
- `muara init --dry-run` prints the generated config without writing it.
- `muara start --dry-run` loads config and reports validation errors without starting the server.

### 10. Tests and quality
- Unit tests for all new CLI commands and helpers.
- TUI prompts tested via injectable input streams.
- All quality gates pass.

---

## Proposed dependencies

Evaluate adding:

- `github.com/charmbracelet/huh` — rich forms and multi-select prompts.
- `github.com/charmbracelet/bubbles` — spinners and progress bars.
- `github.com/charmbracelet/lipgloss` — styled terminal output.
- `github.com/olekukonko/tablewriter` — aligned table output for list commands.

All additions must be justified against OpenMuara's low-memory philosophy. If a feature can be built with stdlib + small helpers, prefer that. Set a binary-size budget (e.g. +≤2 MB stripped binary, measured via `go build -ldflags="-s -w"` and `scripts/check-sizes.sh`); if the budget is exceeded, fall back to a lightweight stdlib form renderer.

---

## Architecture & constraints

- **Root version wiring.** Cobra serves `--version`/`-v` only when `root.Version` is set in `newRootCommand()`. Keep the `version` subcommand because it supports `--json`; `--version` cannot.
- **Build metadata injection.** Dev/CI builds should inject `Version=dev-<short-sha>`, `Commit=<sha>`, and `BuildTime=<iso>` via ldflags. Release builds use the `VERSION` file for `Version` and still inject commit/build-time.
- **TTY fallback.** The TUI wizard must detect non-TTY environments and fall back to plain prompts or `--defaults` without error.
- **Complete config emission.** `init` and `--defaults` must emit all documented config sections (`server`, `admin`, `rate_limit`, `hardened`, `log`, `persistence`, `providers`, `webhook`).
- **Exit-code contract.** Document and test exit codes: `0` success, `1` general error, `2` CLI misuse, `3` config error, `4` environment unhealthy.
- **Provider readiness.** `doctor` should reuse the same validation rules as `provider.Init()` rather than inventing ad-hoc checks.
- **Update-check privacy.** Default update checks must respect `--quiet`, a config flag, and `MUARA_NO_UPDATE_CHECK`. Cache the last-check timestamp; never fail loudly offline.

---

## Recommendations & future enhancements

Optional additions that further raise the CLI quality bar:

### Wizard polish
- **Onboarding tips** — after first `muara init`, print next steps (start server, run an example, open dashboard).
- **Template presets** — "e-commerce", "prepaid top-up", "webhook testing" presets that pre-select providers and config.
- **Input validation inline** — flag invalid webhook URLs, missing provider keys, or duplicate provider names before the final step.
- **Back navigation** — let users go back to a previous wizard step.

### Commands and aliases
- **`muara config edit`** — open `.muara/config.yml` in the user's `$EDITOR`.
- **`muara config validate`** — load and validate config without starting the server.
- **`muara config path`** — print the resolved config path.
- **Command aliases** — `muara tx` for `transaction`, `muara wh` for `webhook`.
- **`muara scenario run <ref>`** — quick alias for `muara scenario success|fail <ref>`.

### Output and scripting
- **Pipe-friendly defaults** — when stdout is not a TTY, prefer plain text or JSON automatically.
- **Man page generation** — `muara docs --man` for system man pages.
- **Structured logging control** — `--log-format json|text` and `--log-level` global flags.
- **Confirmation prompts** — `--yes` / `-y` flag to skip confirmations in scripts.

### Reliability
- **Crash-friendly TUI fallback** — if the TUI fails (e.g. non-interactive CI), fall back to plain prompts or `--defaults`.
- **Idempotent init** — running `muara init` twice is safe and reports what already exists.
- **Migration hints** — when a config version is outdated, suggest `muara migrate`.

---

## Non-goals

- A full GUI outside the terminal.
- Plugin installation from the CLI.
- Remote telemetry or crash reporting.
- Auto-updater binary replacement.

---

## Acceptance criteria

- [ ] `muara --version` and `muara -v` print version and exit 0.
- [ ] `go run ./cmd/muara version` no longer prints `dev (unknown) built unknown`.
- [ ] Dev and release builds inject real version/commit/build-time metadata.
- [ ] `muara init` supports multi-provider selection with a TUI form.
- [ ] `muara init` falls back to plain prompts or `--defaults` in non-TTY environments.
- [ ] `muara init --dry-run` emits a complete, valid config to stdout without writing files.
- [ ] `muara init` safely handles an existing workspace and supports `--force`/`-y`.
- [ ] `muara completion <shell>` generates working completions.
- [ ] `muara doctor` validates config and provider readiness, not just binaries.
- [ ] `muara doctor --json` returns a stable, typed schema.
- [ ] Common errors include actionable next-step hints.
- [ ] Update check respects `--quiet`, config flag, and `MUARA_NO_UPDATE_CHECK`; never blocks or fails loudly offline.
- [ ] List commands produce aligned table output.
- [ ] Tests cover new commands, TUI paths, and non-TTY fallback.
- [ ] All quality gates pass.

---

## References

- `internal/version/version.go`
- `internal/cli/root.go`
- `internal/cli/version.go`
- `internal/cli/init.go`
- `internal/cli/doctor.go`
- `internal/config/wizard.go`
- `internal/config/config.go`
- `Taskfile.yml`

## Self-assessment

**Solidity: 8 / 10**

The initiative now covers version wiring, build metadata, TUI forms, completions, doctor, update checks, output formatting, dry-run, and a full set of future enhancements. Remaining risks are dependency sizing and ensuring the TUI degrades cleanly in CI; these are captured in Architecture & constraints and Acceptance criteria.
