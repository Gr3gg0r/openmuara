> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 01 — Tooling Hygiene

> **Initiative:** OpenMuara v1 Solid Gold
> **Target:** `<repo-root>/`
> **Branch:** `feat/v1-solid-gold`
> **Depends on:** —

---

## Goal

Make `task quality` pass cleanly and make CI run the same full matrix.

## Why now

Currently `task quality` (and `./scripts/audit-trackers.sh`) fails locally for
reasons unrelated to feature code: `muara.yml.example` drifts from the bundled
defaults, and `scripts/smoke-test.sh` has shellcheck warnings. CI also skips
`vuln`, `forbidden`, `scripts`, `sizes`, and `audit-trackers`.

## Scope

### In scope

- Sync `muara.yml.example` with `internal/config.DefaultYAML()`.
- Fix or suppress shellcheck warnings in `scripts/smoke-test.sh`.
- Update `.github/workflows/ci.yml` with a job that runs `task quality`.
- Ensure `task quality` passes end-to-end.

### Out of scope

- Refactoring smoke-test logic.
- Adding new providers or features.

## Acceptance criteria

- [ ] `./scripts/audit-trackers.sh` passes with no errors.
- [ ] `./scripts/check-scripts.sh` passes with no warnings.
- [ ] `task quality` passes locally.
- [ ] CI has a job that runs `task quality` and passes.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- The audit script compares `muara.yml.example` to `internal/config.DefaultYAML()`
  by stripping comments; update the example file to match.
- The shellcheck warnings are for unused variables; either use them or remove them.

## Deliverables

- Code changes on `feat/v1-solid-gold`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
