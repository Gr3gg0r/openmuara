> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Execution Tracker

> **Updated:** 2026-07-03 | **Status:** ✅ Completed / Archived
>
> **Scope:** Migrate the `/_admin` dashboard to a Vite + Preact SPA embedded in the Go binary.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** merged to `dev`
> **Last Agent Action:** Fixed hardened-mode fetch credential bug, added `img-src 'self' data:` CSP, implemented missing `check-bundle-size.js`, updated docs, and pushed to `dev`.
> **Next Agent Action:** None — initiative is archived.

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
5. Product-code commits happen on `feat/web-ui-spa`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Framework decision | `DECISIONS.md`, `README.md` | — | ✅ | — | Vite + Preact selected; assets embedded into Go binary at build time. |
| 02 | Build pipeline scaffold | `web/dashboard/`, `internal/ui/embed.go`, `Taskfile.yml` | 01 | ✅ | — | Vite + Preact + TypeScript; output to `internal/ui/dashboard-dist/`; embedded via Go. |
| 03 | Migrate dashboard shell | `web/dashboard/src/`, `internal/ui/dashboard-dist/index.html` | 02 | ✅ | — | Shell, onboarding, providers grid, tabs, keyboard shortcuts, help modal. |
| 04 | Migrate ledger view | `web/dashboard/src/views/Ledger.tsx` | 03 | ✅ | — | Search, filter, detail panel, replay, auto-refresh. Transactions and Webhooks views also ported. |
| 05 | Keep escape/pay pages server-rendered | `internal/ui/*.html`, `internal/ui/embed.go` | 03 | ✅ | — | Escape/pay pages remain Go templates; smoke test verifies flows. |
| 06 | Tests and CI | `package.json`, `.github/workflows/ci.yml`, `scripts/` | 04, 05 | ✅ | — | Vitest tests, bundle-size budget, UI build/test jobs, `task ui:*` tasks. |
| 07 | Docs update | `runbooks/local-development.md`, `README.md` | 06 | ✅ | — | Dashboard build, dev server, and quality-gate documentation updated. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| UI Build | `cd web/dashboard && npm run build` | Produces embeddable assets | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| UI Test | `cd web/dashboard && npm run test:ci` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |
| Bundle Size | `cd web/dashboard && npm run bundle-size` | ≤ 100 KiB gzipped JS, ≤ 250 KiB total | ✅ |
| E2E | `cd web/dashboard && npm run test:e2e` | Hardened dashboard loads with credentialed URL | ✅ |

---

## Decisions

- D001 ✅ Framework: Vite + Preact.
- D002 ✅ Build assets embedded by Go; no runtime Node dependency.
- D003 ✅ Existing `/_admin` routes and provider escape/pay URLs remain stable.
- D004 ✅ SPA must stay lightweight (small bundle, low memory).
- D005 ✅ TypeScript selected for dashboard source.
- D006 ✅ Escape/pay pages remain server-rendered HTML in `internal/ui/`.
- D007 ✅ Source directory `web/dashboard/`, build output `internal/ui/dashboard-dist/`.
- D008 ✅ Bundle-size budget: initial JS ≤ 100 KiB gzipped, total `dist/` ≤ 250 KiB.

---

## Recommendations Implemented

All audit recommendations were implemented:

- TypeScript for type safety.
- Bundle-size budget and CI enforcement via `web/dashboard/scripts/check-bundle-size.js`.
- Escape/pay pages kept server-rendered.
- Vitest UI tests, including fetch credential-stripping coverage.
- Playwright E2E test for hardened-mode dashboard with credentials in the URL.
- CI jobs for UI build/test/e2e.
- Taskfile tasks for UI workflows, including `task ui:e2e`.
- Accessibility (ARIA, keyboard shortcuts) and error boundaries.
- CSP-compatible production build with `img-src 'self' data:` for the inline favicon.
- Standardized API client with CSRF token handling and embedded-credential stripping for hardened-mode URLs.
- `webhooks_enabled` onboarding flag so the failed-webhook alert only queries the webhooks API when a dispatcher is configured.

## Future Recommendations

Deferred post-v1 items:

- Dark mode toggle.
- Offline caching / service worker.
- End-to-end visual regression tests.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/archive/openmuara-web-ui-spa/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/archive/openmuara-web-ui-spa/README.md` | Goals, options, recommendation |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
