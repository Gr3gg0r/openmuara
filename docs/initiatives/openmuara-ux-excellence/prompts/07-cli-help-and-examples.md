> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P07 — CLI Help and Structured Output

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

Make the CLI self-describing for humans and parseable for AI agents.

## Why now

Developers and AI agents both start with `muara --help`. Right now Cobra lists flags but not common invocations, and there is no structured output mode for scripting.

## Scope

### In scope

- Add an `Example` field to every Cobra command in `internal/cli/`:
  - `muara init`
  - `muara init --defaults`
  - `muara start`
  - `muara start --config path/to/config.yml`
  - `muara doctor`
  - `muara scenario`
  - `muara webhook replay --ref <ref>`
  - `muara audit`
  - `muara migrate`
  - `muara version`
- Examples should be runnable and cover the most common use case.
- Add global `--json` flag where it makes sense:
  - `muara doctor --json` returns structured tool/status results.
  - `muara version --json` returns version metadata.
  - `muara scenario --json` returns scenario results.
- Document stable JSON schemas for `--json` outputs under `docs/cli-schemas/` so AI agents can rely on field names and types.
- Add global `--quiet` flag to suppress informational output (errors still go to stderr).
- Update `runbooks/local-development.md` to reference the in-CLI examples.
- Add tests for examples and structured output.

### Out of scope

- Rewording all command descriptions.
- Adding new commands beyond replay/list helpers.

## Acceptance criteria

- [ ] Every CLI command has a runnable `Example`.
- [ ] `muara <command> --help` displays the example.
- [ ] `muara doctor --json` outputs valid JSON matching the documented schema.
- [ ] `muara version --json` outputs valid JSON matching the documented schema.
- [ ] `muara --quiet` suppresses non-error stdout.
- [ ] Tests assert examples exist and JSON output parses.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Cobra supports `Example` strings natively.
- Add `--json` and `--quiet` as persistent flags on the root command, then check them in command runners.
- Keep JSON schemas stable; tests should validate shape, not exact values.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated CLI tests.
- Updated `runbooks/local-development.md`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
