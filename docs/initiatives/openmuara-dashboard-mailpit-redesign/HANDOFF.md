> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard — Mailpit-Style Redesign — Handoff

> **Updated:** 2026-07-06
> **Branch:** `feat/dashboard-mailpit-redesign`

## Current State

P01–P06 complete. Dashboard has Mailpit-style left navigation, Ledger as default, Webhooks delivery log, Settings provider grid, ProviderDetail with version tabs/base URL/webhook config/env vars, and optional dual-port runtime.

## What Was Last Done

- **P01** — Implemented `AppShell` and `SidebarNav`, updated `app.tsx` routing, preserved header/shortcuts, added tests. Commit `a501b1d`.
- **P02** — Added optional `server.admin_port`, split router into `NewProviderRouter`/`NewAdminRouter`, preserved single-port `NewRouter`, injected admin API base URL into dashboard HTML, added tests. Commit `5ad83c9`.
- **P03** — Promoted Ledger to default view, added `FilterToolbar`, created `LedgerDetail` page, added tests. Commit `05efc4f`.
- **P04** — Refactored `WebhooksView` to delivery-log-only with `FilterToolbar`, created `WebhookDetail` page, added tests. Commit `5c4377a`.
- **P05** — Built Settings provider card grid, `ProviderDetail` with enable toggle, version tabs, base URL, sample endpoint, per-provider webhook URL, and env var reference. Commit `5ad83c9`.
- **P06** — Added Settings/ProviderDetail/dual-port tests, fixed WCAG AA contrast, verified bundle size and quality gates. Commit `5ad83c9`.
- **Docs** — Updated initiative README, TRACKING, DECISIONS, HANDOFF, and P04/P05 prompts to reflect final scope.

## Quality Gate Results

- `go build ./...` ✅
- `go test ./...` ✅
- `go vet ./...` ✅
- `golangci-lint run` ✅ 0 issues
- `cd web/dashboard && npm run test:ci` ✅ 74 tests passed
- `cd web/dashboard && npm run build` ✅
- `node web/dashboard/scripts/check-bundle-size.js` ✅ JS 21.84 KiB ≤ 100 KiB; total 128.01 KiB ≤ 250 KiB
- `cd web/dashboard && node scripts/a11y-contrast-check.js` ✅ No WCAG AA contrast violations (light + dark)

## Visual Test Results

- Ledger view: full-width table with filter toolbar ✅
- Webhooks view: delivery log only, no config UI ✅
- Settings view: provider card grid with status badges ✅
- ProviderDetail (Fawry): version tabs, base URL, sample endpoint, webhook target, env vars ✅

## What Is Next

Await user visual sign-off; open PR from `feat/dashboard-mailpit-redesign` to `dev` when approved.

## Useful Commands

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
```
