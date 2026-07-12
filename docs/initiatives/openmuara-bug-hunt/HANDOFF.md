> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Bug Hunt — Handoff

> **Updated:** 2026-07-06
> **Branch:** `feat/bug-hunt` (merged into `dev` and deleted locally)

## Current State

The bug-hunt initiative is complete. All approved recommendations (E1–E12) were implemented, all quality gates passed, the visual baseline diff is stable, and the work has been landed on `dev`.

## Pre-Flight Checklist for the Next Agent

This initiative is complete. The next quality cycle should start from `dev` and read `docs/initiatives/openmuara-bug-hunt/appendices/b-recommendations.md` for the baseline that was implemented.

- [x] All work landed on `dev`.
- [x] `feat/bug-hunt` branch deleted locally.
- [x] Root `AGENTS.md` followed throughout.
- [x] Initiative README, tracker, and prompts updated.

## What Was Last Done

- Implemented all approved bug-hunt recommendations (E1–E12):
  - E1: Deterministic Playwright visual baseline with automated diff.
  - E2: GitHub bug report and PR templates aligned with the bug register.
  - E3: `govulncheck`, `npm audit --production`, and `golangci-lint` required in CI.
  - E4: Fuzz/property tests for signatures, idempotency keys, and transaction state machine.
  - E5/E6: Mutation-testing guide and recurring bug-hunt process documentation.
  - E7: Provider contract conformance tests with golden files for every provider/version.
  - E8: Webhook dispatch chaos tests for retries, non-2xx, and timeouts.
  - E9: Root `KNOWN_ISSUES.md` synced from deferred bugs.
  - E10: `scripts/release-notes.sh` scraping fixed bug IDs from `TRACKING.md`.
  - E11: PR coverage-regression comment workflow.
  - E12: `internal/errcode` taxonomy adopted in the webhook dispatcher.
- Ran full quality gate suite (build, test, race, vet, lint, frontend tests/build, bundle size, a11y, visual baseline) — all green.
- Updated initiative docs, root `TRACKING.md`, and `CHANGELOG.md`.
- Merged `feat/bug-hunt` into `dev` and deleted the feature branch.

## Baseline Environment

Capture the actual environment when P01 starts:

```bash
# Fill in before P01
go version:
node version:
OS:
last dashboard commit on this branch:
```

## Known Baseline from Dashboard Redesign

- `feat/dashboard-mailpit-redesign` was green through P06.
- Quality gates passed: build, test, vet, lint, frontend tests/build, bundle size, a11y contrast.
- Dashboard features: left nav (Ledger/Webhooks/Settings), Ledger default, filter toolbars, detail pages, provider settings with version tabs/env vars/base URLs, optional dual-port runtime.

## Findings Workflow

During P01, use this workflow for every confirmed bug:

1. Create a new file in `findings/` from `findings/TEMPLATE.md`, named `BXXX-<short-description>.md`.
2. Attach or reference Playwright MCP screenshots in the finding file.
3. Link the finding file in the `TRACKING.md` bug register (e.g., add a `Finding File` column or note).
4. Store the P01 visual baseline under `findings/visual-baseline/` for comparison in P06.

## What Is Next

No further work on this initiative. The next quality cycle can build on the E1–E12 foundation (see `appendices/b-recommendations.md`).

## Open Questions / Blockers

- None.

## Useful Commands

```bash
# Go gates
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run

# Frontend gates
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
cd web/dashboard && node scripts/a11y-contrast-check.js

# Runtime / smoke
./scripts/smoke-test.sh || true
go run ./cmd/muara start

# Visual baseline (example Playwright MCP flow)
# 1. Start server: go run ./cmd/muara start
# 2. Navigate to http://127.0.0.1:9000/_admin with Playwright MCP.
# 3. Capture Ledger, Webhooks, Settings, and ProviderDetail screenshots.
```
