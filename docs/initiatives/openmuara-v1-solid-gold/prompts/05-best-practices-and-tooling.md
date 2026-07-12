> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 05 — Best Practices and Tooling

> **Initiative:** OpenMuara v1 Solid Gold
> **Target:** `<repo-root>/`
> **Branch:** `feat/v1-solid-gold`
> **Depends on:** Prompts 01, 02

---

## Goal

Strengthen static analysis, pre-commit enforcement, dependency hygiene, and
reproducible builds.

## Why now

The project already has solid tests and linting, but it can still drift. Stronger
tooling catches problems before they reach CI.

## Scope

### In scope

- Add stronger linters to `.golangci.yml` (`gosec`, `staticcheck`, `ineffassign`,
  `unparam`, `errcheck`) — enable gradually and fix or suppress findings.
- Extend the existing `.pre-commit-config.yaml` to also run `go test ./...`
  and `scripts/check-forbidden.sh`.
- Add a `govulncheck` CI job so vulnerability scanning is not only local.
- Add Dependabot config for `go.mod` and GitHub Actions.
- Add `-trimpath` to release builds and document reproducible-build verification.
- Run `deadcode`/`unused` and remove or test genuinely unused code.

### Out of scope

- Rewriting architecture.
- Adding new CI platforms (e.g., CircleCI).

## Acceptance criteria

- [ ] `golangci-lint run` passes with the expanded linter set.
- [ ] Pre-commit hooks run without errors and are documented in `runbooks/quality-gates.md`.
- [ ] CI has a `govulncheck` job.
- [ ] Dependabot config is valid.
- [ ] Release binaries are built with `-trimpath`.
- [ ] `deadcode`/`unused` findings are addressed (removed, tested, or suppressed with a reason).
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Enable new linters one at a time to avoid a giant diff.
- For `gosec`, suppress false positives with `#nosec` comments and a brief justification.

## Deliverables

- Code changes on `feat/v1-solid-gold`.
- Updated `runbooks/quality-gates.md`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
