> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid Gold — Handoff

> **Created:** 2026-07-01
> **Last Updated:** 2026-07-01
> **Status:** ✅ Completed

---

## Current Context

v1 feature work is complete. The project already passes `go build`, `go test`,
`go vet`, `golangci-lint`, and the smoke test. However, running the full local
quality matrix (`task quality`) reveals a few hygiene issues, and several
packages are below the 80% coverage target. This initiative bundles the
remaining polish into five prompts.

## Branch

`feat/v1-solid-gold` (created from `dev` and pushed to origin).

## What has been done

- Initiative docs created in `docs/initiatives/openmuara-v1-solid-gold/`.
- P01 tooling hygiene completed:
  - Synced `muara.yml.example` with `internal/config.DefaultYAML()`.
  - Fixed shellcheck warnings in `scripts/smoke-test.sh`.
  - Added a `quality` CI job that runs `task quality`.
  - Verified `task quality` passes locally.
- P02 coverage backfill completed:
  - Backfilled tests for `internal/billplz`, `internal/ipay88`, `internal/toyyibpay`.
  - Backfilled tests for `internal/fawry/v1`, `internal/fawry/v2`, `internal/testutil`.
  - Backfilled tests for `internal/cli`.
  - Backfilled tests for `internal/ui`.
  - Made `cmd/muara/main.go` injectable for full error-path coverage.
  - All packages now ≥80% coverage; total coverage 89.0%.
- P03 debuggability completed (commit `259a821`):
  - Added `TraceID` to `engine.Transaction` and both memory/SQLite stores (additive migration).
  - Captured trace IDs at all transaction creation sites.
  - Added `TraceID` to `webhook.Attempt`; outgoing webhooks now include `X-Trace-Id` header.
  - Dispatcher falls back to ledger trace ID when context has none.
  - Wired request contexts through provider dispatch callbacks.
  - Exposed `trace_id` in `/_admin/ledger`, `/_admin/transactions/{ref}`, and `/_admin/webhooks/{ref}`.
  - Rendered trace IDs in dashboard transaction and webhook detail panels.
  - Added `muara transaction list` and `muara transaction inspect <ref>` CLI commands.
  - Added optional pprof endpoints gated by `server.pprof: true`.
  - Updated `runbooks/debugging.md` with trace-ID correlation and profiling sections.
  - Added tests; all quality gates pass; coverage remains 89.0%.

## Prompt Inventory

| Step | Title | Status |
|------|-------|--------|
| 01 | Tooling hygiene | ✅ |
| 02 | Coverage backfill | ✅ |
| 03 | Debuggability | ✅ |
| 04 | Dashboard usability | ✅ |
| 05 | Best practices and tooling | ✅ |

## Decisions already made

See `DECISIONS.md`. Preliminary:

- All changes are additive/hygiene; no breaking changes.
- P03 user approved; trace IDs propagated via `X-Trace-Id` header only, no request/response body logging.
- P04 dashboard usability changes verified with Playwright/browser snapshots at desktop and 375px widths.

## What has been done (P04)

- Added a dismissible failed-webhook alert bar at the top of `/_admin`.
  - Polls `/_admin/webhooks?status=failed&limit=1` every refresh cycle.
  - Dismiss state stored in `sessionStorage`.
  - Link switches to the Webhooks tab.
- Added a **Copy curl** button to every provider card that has a `sample_route`.
  - Uses `navigator.clipboard` with a `textarea` fallback.
  - Button briefly shows "Copied!" feedback.
- Improved ledger/transactions/webhooks table responsiveness on 375px viewports.
  - Wrapped tables in `.table-wrap` for horizontal scroll.
  - Hid lower-priority columns (`Summary`, `Created`, `URL`, `Last Error`) via `.hide-narrow`.
  - Stacked the ledger toolbar and switched provider grid to a single column on narrow screens.
- Verified existing keyboard shortcuts (`?`, `/`, `1/2/3`, `Esc`) still work.
- All quality gates pass; coverage remains 89.0%.

## What has been done (P05)

- Expanded `.golangci.yml` with `gosec`, `staticcheck`, `ineffassign`, `unparam`,
  `errcheck`, `misspell`, `revive`, and `unused`; fixed or suppressed all 31
  findings across the codebase.
- Extended `.pre-commit-config.yaml` to run `go test ./...` and
  `scripts/check-forbidden.sh`; documented the hooks in
  `runbooks/quality-gates.md`.
- Added a `govulncheck` CI job in `.github/workflows/ci.yml`.
- Added `.github/dependabot.yml` for `go.mod` and GitHub Actions.
- Added `-trimpath` to `task release:build` and the release workflow; documented
  reproducible-build verification in `runbooks/quality-gates.md`.
- Removed dead code `ParseValidationPort` from `internal/config/validation.go`
  after confirming no callers.
- Fixed a race in `internal/webhook/dispatcher_test.go` trace-ID tests by using
  buffered channels instead of unsynchronized shared state.
- Synced `muara.yml.example` with `internal/config.DefaultYAML()` (`server.pprof`).
- Added `scripts/check-forbidden.sh`, `scripts/check-scripts.sh`, and
  `scripts/check-sizes.sh` as tracked quality helpers.
- All quality gates pass, including `go test -race ./...`; total coverage remains
  89.0%.

## Next step

Open a pull request from `feat/v1-solid-gold` to `dev` and run the full CI matrix.

## Open questions

- None.
