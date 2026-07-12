> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P03 — Actionable Config Validation Errors

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

When `muara start` or `muara doctor` fails because of bad config, tell the user exactly what is wrong, where it is, and how to fix it.

## Why now

Config errors currently bubble up from Viper/Go as low-level messages without line numbers or context. This is frustrating for new users.

## Scope

### In scope

- Add a config validation pass in `internal/config` that runs after load and returns a slice of structured errors.
- Each error includes:
  - `field` — dotted path, e.g., `providers.fawry.config.version`.
  - `message` — human-readable description.
  - `hint` — suggested fix.
  - `file` and `line` — source location when available (best-effort via raw YAML parse; fall back to file path only).
- Validate:
  - Unknown provider names under `providers`.
  - Missing required fields for enabled providers (e.g., `secret_key`, `api_key`).
  - Invalid `version` values for versioned providers.
  - Deprecated top-level `fawry` / `stripe` keys with a migration hint.
- Print validation results in `muara doctor`.
- Fail `muara start` on validation errors with a clear, grouped message.
- Add tests for validation scenarios.

### Out of scope

- Rewriting the config loader.
- Network validation (e.g., webhook URL reachability).

## Acceptance criteria

- [ ] `muara doctor` reports config problems with field, message, hint, and best-effort file/line.
- [ ] `muara start` fails early with grouped validation errors.
- [ ] Tests cover at least unknown provider, missing required field, and invalid version.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Keep the validator separate from the loader so tests can call it directly.
- Provider required-field rules can be a map in `internal/config` or derived from provider metadata later.

## Deliverables

- Code changes on `feat/ux-excellence`.
- New/updated tests in `internal/config`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
