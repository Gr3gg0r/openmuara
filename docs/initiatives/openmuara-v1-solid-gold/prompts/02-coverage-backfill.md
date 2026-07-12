> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 02 — Coverage Backfill

> **Initiative:** OpenMuara v1 Solid Gold
> **Target:** `<repo-root>/`
> **Branch:** `feat/v1-solid-gold`
> **Depends on:** —

---

## Goal

Bring every Go package to at least 80% test coverage.

## Why now

Current package coverage varies widely. The weakest packages are:

| Package | Current | Priority | Test approach |
|---------|---------|----------|---------------|
| `internal/ui` | 21.4% | High | Playwright or jsdom dashboard test |
| `internal/fawry/v2` | 38.8% | High | Unit tests for payload builder + webhook handler |
| `internal/cli` | 73.9% | Medium | Add tests for `start`, `doctor`, `scenario` happy/error paths |
| `internal/fawry/v1` | 70.0% | Medium | Webhook handler and provider tests |
| `internal/ipay88` | 73.4% | Medium | Response/admin/paypage tests |
| `internal/toyyibpay` | 74.1% | Medium | Paypage and webhook tests |
| `internal/billplz` | 79.7% | Low | One or two extra error-path tests |
| `internal/testutil` | 78.8% | Low | Exercise fake helpers used by other tests |

Target: every package ≥80%. Do not chase 100% globally.

## Scope

### In scope

- Add unit and integration tests for the weakest packages.
- Add dashboard UI tests for `internal/ui` (e.g., Playwright or jsdom).
- Target 80% for every package; higher for critical packages (`engine`, `webhook`).

### Out of scope

- Refactoring code solely to make it testable.
- Coverage for generated code.

## Acceptance criteria

- [ ] `go test -cover ./...` reports every package ≥80%.
- [ ] New tests assert meaningful behavior, not just line hits.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- For `internal/ui`, a Playwright test that starts the server, loads `/_admin`,
  sends a charge, and asserts a ledger row is high value.
- For `internal/fawry/v2`, focus on `NewPayloadBuilder` and `NewWebhookHandler`.

## Deliverables

- Code changes on `feat/v1-solid-gold`.
- Updated `runbooks/testing.md` if coverage policy or target changes.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
