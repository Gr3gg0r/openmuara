> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P06 — Tests and Quality Gates

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** P01–P05
> **Target files:** `web/dashboard/tests/`, `web/dashboard/e2e/`, `internal/server/*_test.go`, `internal/config/*_test.go`
> **Status:** ⬜

## Goal

Verify the redesigned dashboard with unit, integration, accessibility, and backend tests, and ensure all quality gates pass.

## Tasks

- [ ] Add/update unit tests for `AppShell`, `SidebarNav`, `FilterToolbar`, `LedgerDetail`, `WebhookDetail`, `Settings`, and `ProviderDetail`.
- [ ] Add integration tests for navigation flows: Ledger → Webhooks → Settings → Provider Detail and back.
- [ ] Add Playwright tests for keyboard navigation, detail page routing, and filter persistence.
- [ ] Add axe-core assertions for every new view.
- [ ] Add backend tests for `/_admin/providers/{name}` metadata and dual-port startup.
- [ ] Run `go test ./...`, `golangci-lint run`, `npm run test`, `npm run build`, and `check-bundle-size.js`.
- [ ] Fix any failures or bundle-size regressions.

## Acceptance Criteria

- [ ] All Go tests pass.
- [ ] All frontend unit and integration tests pass.
- [ ] Axe-core reports zero serious violations on all views.
- [ ] Dashboard JS + CSS bundle remains ≤ 150 KB gzipped.
- [ ] All quality gates pass.

## Quality Gates

Run before committing:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
```

## Notes

- This is the final prompt; no new features should be added here.
- Update `HANDOFF.md` and `TRACKING.md` after this prompt is complete.
