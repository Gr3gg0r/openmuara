> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Testing Gold Standard

> **Status:** 🟡 Active  
> **Scope:** Establish OSS-grade testing practices and backfill tests until every contributor has a clear, passing quality bar.  
> **Product Branch:** `dev`  
> **Selected Approach:** Option A — 80% coverage, refactor-first.

---

## Why this initiative

OpenMuara v1 shipped with ~73% coverage and a 50% gate. Several CLI commands, lifecycle helpers, and provider setters are untested. Global state and hardcoded ports make some tests flaky or impossible. This initiative fixes the root causes and raises the bar.

## Goals

1. Document test conventions in `runbooks/testing.md`.
2. Refactor global state into injectable dependencies.
3. Provide reusable test utilities in `internal/testutil`.
4. Backfill unit tests until all exported code is covered.
5. Add provider contract tests with golden files.
6. Harden integration/E2E tests with random ports and race detection.
7. Add fuzz/property tests for signatures and the state machine.
8. Split CI into lint/unit/integration/smoke jobs and raise coverage to 80%.

## Entry points

| Phase | Prompt | What it covers |
|-------|--------|----------------|
| 0 | [`prompts/00-charter-and-runbook.md`](prompts/00-charter-and-runbook.md) | Standards, runbook, initiative docs |
| 1 | [`prompts/01-testability-refactor.md`](prompts/01-testability-refactor.md) | DI refactor for provider/plugin/cli/server |
| 2 | [`prompts/02-shared-test-utilities.md`](prompts/02-shared-test-utilities.md) | `internal/testutil` package |
| 3 | [`prompts/03-unit-test-backfill.md`](prompts/03-unit-test-backfill.md) | Backfill CLI, lifecycle, audit, provider tests |
| 4 | [`prompts/04-provider-contract-tests.md`](prompts/04-provider-contract-tests.md) | Golden-file contract tests |
| 5 | [`prompts/05-integration-e2e-hardening.md`](prompts/05-integration-e2e-hardening.md) | Random-port smoke, parallel race tests |
| 6 | [`prompts/06-advanced-testing.md`](prompts/06-advanced-testing.md) | Fuzz, property, shuffle |
| 7 | [`prompts/07-ci-quality-gates.md`](prompts/07-ci-quality-gates.md) | CI split, coverage gate, Taskfile |

## Success criteria

- `go test ./...` passes with `-race -shuffle=on`.
- Coverage ≥ 80%.
- Smoke test uses random ports and passes reliably.
- CI is green and gives feedback in <10 minutes.
- A new contributor can read `runbooks/testing.md` and write a passing test.
