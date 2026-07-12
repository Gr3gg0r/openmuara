> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P01 — First-Run Config Wizard

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

Make `muara init` guide a first-time user to a working config in under one minute, while keeping non-interactive usage intact.

## Why now

Right now `muara init` silently writes a generic YAML file and exits. New users must read docs to know which providers to enable and what credentials to set. This is the highest-friction first step.

## Scope

### In scope

- Add an interactive flow to `muara init` when stdin is a TTY and no `--defaults` flag is passed.
- Ask ≤5 questions:
  1. What real payment provider are you emulating? (Stripe, Fawry, Billplz, ToyyibPay, iPay88, SenangPay, Default/DIY)
  2. Do you want webhooks delivered to a local test URL? (optional)
  3. Do you want FPX / redirect-style methods enabled? (only if relevant)
  4. Preferred log level? (`info` default)
  5. Confirm write.
- Generate a config that enables the chosen provider with sensible sample credentials and comments.
- Keep `--defaults` (or `--non-interactive`) flag to write the existing generic config for CI/scripts.
- Update `muara init` tests to cover both interactive and non-interactive paths.
- Update `muara.yml.example` with provider-specific comments if missing.

### Out of scope

- Persisting wizard answers beyond the generated YAML.
- Provider-specific deep configuration (e.g., full Stripe product catalog).
- Web UI wizard.

## Acceptance criteria

- [ ] `muara init` in a TTY asks the questions and writes a tailored `config.yml`.
- [ ] `muara init --defaults` still writes the generic config and works in CI.
- [ ] Generated config passes validation and starts successfully.
- [ ] Tests cover both interactive and non-interactive paths.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use `isatty` check on stdin (`golang.org/x/term` is already common; check `go.mod`).
- Keep the interactive reader testable by abstracting a `promptFunc` interface.
- Provider sample configs should live in `internal/config` as a map keyed by provider name (e.g., `WizardTemplates`) rather than string-concatenating YAML.
- Default answers should let a user press Enter to accept the recommended path.
- The generated config must still be loadable by `config.Load` and pass `Config.Validate`.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated `internal/cli/init_test.go`.
- Updated `muara.yml.example` if needed.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
