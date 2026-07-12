> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard — Mailpit-Style Redesign — Execution Tracker

> **Updated:** 2026-07-09 | **Status:** ✅ Completed / Merged to `dev`
>
> **Scope:** Redesign the admin dashboard with a Mailpit-like left navigation, Ledger as the default view, and a Settings view for provider configuration. Top-level Webhooks view is a delivery log only; per-provider webhook targets are edited in Settings.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `dev` (feature branch merged and removed)
> **Last Agent Action:** Verified all P01–P06 dashboard redesign work is present in `dev` history and the built dashboard is served from `internal/ui/dashboard-dist`.
> **Next Agent Action:** None — initiative closed.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen |

---

## Execution Rules

1. Execute prompts in order unless marked **[PARALLEL SAFE]**.
2. Every prompt MUST end with: tests passing → git commit → update this file to `✅`.
3. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY prompt, update `HANDOFF.md`.
5. Product-code commits happen on `feat/dashboard-mailpit-redesign`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Shell layout and navigation | `web/dashboard/src/components/AppShell.tsx`, `SidebarNav.tsx`, `app.tsx`, `styles.css` | — | ✅ | a501b1d | Replace top tabs with left sidebar; three nav items: Ledger, Webhooks, Settings. Tests included. |
| 02 | Dual-port runtime | `internal/server/server.go`, `internal/config/config.go`, `internal/cli/start.go`, `web/dashboard/src/api.ts` | — | ✅ | 5ad83c9 | Optional `server.admin_port`; admin UI/API on one port, provider endpoints on another; dashboard discovers admin API base URL. Tests included. |
| 03 | Ledger default view and detail page | `web/dashboard/src/views/Ledger.tsx`, `LedgerDetail.tsx`, `components/FilterToolbar.tsx` | 01 | ✅ | 05efc4f | Promote Ledger to default outlet; reusable filter toolbar; row click navigates to detail page. Tests included. |
| 04 | Webhooks view and detail page | `web/dashboard/src/views/Webhooks.tsx`, `WebhookDetail.tsx`, `components/FilterToolbar.tsx` | 01 | ✅ | 5c4377a | Top-level Webhooks delivery-log view with filter toolbar; row click navigates to detail page. Tests included. |
| 05 | Provider settings | `web/dashboard/src/views/Settings.tsx`, `ProviderDetail.tsx`; `internal/server/admin_api.go`, `internal/config/wizard.go` | 01, 04 | ✅ | 5ad83c9 | Provider card grid, detail page, enable toggle, version tabs, base URL, per-provider webhook URL, env vars. |
| 06 | Tests and quality gates | `web/dashboard/tests/`, `internal/server/*_test.go` | 01–05 | ✅ | 5ad83c9 | Unit, integration, a11y tests; bundle-size check; lint/vet/build. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Frontend test | `cd web/dashboard && npm run test` | All pass | ✅ |
| Frontend build | `cd web/dashboard && npm run build` | Passes | ✅ |
| Bundle size | `node web/dashboard/scripts/check-bundle-size.js` | ≤ 150 KB gzipped | ✅ |
| A11y contrast | `cd web/dashboard && node scripts/a11y-contrast-check.js` | Zero WCAG AA contrast violations | ✅ |

---

## Decisions

- D001 ✅ The dashboard will use a fixed left navigation with exactly three items: Ledger, Webhooks, Settings.
- D002 ✅ Ledger becomes the default view at `/_admin`; existing `tab=` query strings redirect to the new `view=` parameter.
- D003 ✅ Provider configuration lives under Settings; webhook delivery logs remain a separate top-level view.
- D004 ✅ Environment variables are shown as read-only reference names derived from `MUARA_<PROVIDER>_<KEY>`.
- D005 ✅ Version tabs (v1/v2) appear only when the provider reports more than one version.
- D006 ✅ Design priority stack: UI > UX > performance > usability > philosophy > efficiency > memory size.
- D007 ✅ Every table view has a reusable filter toolbar.
- D008 ✅ Ledger and webhook rows navigate to dedicated detail pages.
- D009 ✅ Admin UI and provider endpoints run on separate optional ports.
- D010 ✅ Per-provider webhook targets are configured inside Settings → Provider Detail, not on the top-level Webhooks delivery log.
- D011 ✅ Light-theme muted text and shortcut colors were darkened to pass WCAG AA contrast on all backgrounds.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-dashboard-mailpit-redesign/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-dashboard-mailpit-redesign/README.md` | Scope, goals, architecture |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
