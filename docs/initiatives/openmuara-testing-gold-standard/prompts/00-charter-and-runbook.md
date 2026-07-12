# Prompt P20 — Testing Charter & Runbook

## Goal

Establish the testing standards and conventions for the OpenMuara Testing Gold Standard initiative.

## Acceptance Criteria

- [ ] `docs/initiatives/openmuara-testing-gold-standard/README.md` exists and is current.
- [ ] `docs/initiatives/openmuara-testing-gold-standard/TRACKING.md` exists with all P20–P27 prompts.
- [ ] `docs/initiatives/openmuara-testing-gold-standard/RISKS.md` exists with risk matrix and detailed entries.
- [ ] `docs/initiatives/openmuara-testing-gold-standard/KNOWN_ISSUES.md` exists with pre-existing testability issues.
- [ ] `docs/initiatives/openmuara-testing-gold-standard/HANDOFF.md` exists.
- [ ] `runbooks/testing.md` documents:
  - Test pyramid (unit → integration → contract → E2E).
  - Naming conventions for test files and functions.
  - Table-driven test pattern.
  - When to use fakes vs mocks vs real dependencies.
  - How to use `testdata/` and golden files.
  - Coverage policy (80% threshold, no trivial getters).
  - How to run the test suite and CI gates locally.
- [ ] Root `TRACKING.md` links to the new initiative.
- [ ] All gates pass: `go build ./...`, `go test ./...`, `go vet ./...`, `golangci-lint run`, `./scripts/smoke-test.sh`.

## Files to Create/Change

- `docs/initiatives/openmuara-testing-gold-standard/README.md`
- `docs/initiatives/openmuara-testing-gold-standard/TRACKING.md`
- `docs/initiatives/openmuara-testing-gold-standard/RISKS.md`
- `docs/initiatives/openmuara-testing-gold-standard/KNOWN_ISSUES.md`
- `docs/initiatives/openmuara-testing-gold-standard/HANDOFF.md`
- `runbooks/testing.md`
- `TRACKING.md`

## Response Shape

Return:
1. A link to `runbooks/testing.md`.
2. A summary of the test pyramid and conventions.
3. Quality gate results.

## Test Notes

- This prompt is docs-only; no product code changes.
- Verify all markdown links are valid.
- Run `task check` to ensure docs changes do not break anything.
