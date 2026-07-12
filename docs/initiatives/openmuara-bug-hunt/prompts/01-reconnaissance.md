> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P01 — Reconnaissance

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** —
> **Target files:** `TRACKING.md`, `RISKS.md`, `HANDOFF.md`
> **Status:** ⬜

## Goal

Discover and document bugs across the OpenMuara v1 codebase using tests, static analysis, runtime exploration, and visual inspection. Capture a visual baseline of the dashboard so later prompts can detect regressions.

## Tasks

- [ ] Run the full Go test suite (`go test ./...`) and record any failures or flakes.
- [ ] Run with race detection (`go test -race ./...`) and record any races.
- [ ] Run `golangci-lint run` and `go vet ./...`; record any warnings or suspicious patterns.
- [ ] Run frontend tests (`cd web/dashboard && npm run test:ci`) and build; record issues.
- [ ] Run the a11y contrast check and note any regressions.
- [ ] Run bundle-size check and note any budget violations.
- [ ] Run `./scripts/smoke-test.sh` if available; note failures.
- [ ] Run dependency vulnerability scan (`govulncheck ./...` or `npm audit --production` in `web/dashboard`) and record actionable findings.
- [ ] Search the codebase for `TODO`, `FIXME`, `BUG`, `HACK`, `XXX`, and `panic` usage.
- [ ] Review recent commits on `feat/dashboard-mailpit-redesign` and `dev` for potential regressions.
- [ ] Review provider contract golden files for drift.
- [ ] Exercise core provider emulation paths with `curl` or the testsdk (Fawry charge, Stripe checkout, SenangPay callback, webhook dispatch).
- [ ] Exercise config write endpoints (`PATCH /_admin/config/providers`, `PATCH /_admin/config/webhooks`) with edge cases (empty body, unknown provider, duplicate keys, missing CSRF token).
- [ ] Exercise CLI edge cases (`muara start` with invalid config, `muara doctor`, `muara transaction inspect` with missing ref).
- [ ] Capture a visual baseline of the dashboard with Playwright MCP: Ledger, Webhooks, Settings, and at least one ProviderDetail (e.g. Fawry).
- [ ] Populate the bug register in `TRACKING.md` with at least 5 findings, each with severity, area, summary, minimal reproduction, root cause category, and environment details.
- [ ] Update `RISKS.md` with any systemic risks discovered during reconnaissance.
- [ ] Update `HANDOFF.md` with what was checked and what was found.

## Acceptance Criteria

- [ ] Bug register in `TRACKING.md` has at least 5 documented findings.
- [ ] Each finding includes a severity (P0/P1/P2), area, summary, reproduction steps, root cause category, and environment.
- [ ] Visual baseline screenshots are attached or referenced in `HANDOFF.md`.
- [ ] No quality gate regressions introduced by reconnaissance-only changes (there should be no code changes in this prompt).
- [ ] `HANDOFF.md` updated with what was checked and what was found.

## Quality Gates

Run to establish baseline; this prompt should not change product code, so gates must remain as they are:

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
cd web/dashboard && node scripts/a11y-contrast-check.js
./scripts/smoke-test.sh || true
```

## Notes

- Do **not** fix bugs in this prompt. Only discover and document.
- If a bug is trivial to fix (e.g., typo in error message), still document it first; fixes happen in P03/P04.
- Prefer reproducible test cases over manual curl recipes when possible.
- Record environment details (Go version, Node version, OS) for any bug that may be environment-specific.
