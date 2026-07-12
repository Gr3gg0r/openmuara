> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Accessibility & Usability Polish — Handoff

> **Last updated:** 2026-07-03
> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Branch:** `feat/a11y-enhancements` (to be merged into `dev` and deleted)
> **Status:** ✅ COMPLETE

---

## What was shipped

- P01 — Dashboard keyboard navigation: clickable table rows, help-modal focus trap, tab arrow-key navigation, detail-panel focus management.
- P02 — Dashboard labels and live regions: search input labels, unique Copy curl labels, copy announcements, onboarding aria-expanded, failed-webhook alert as button.
- P03 — Provider pages: `<main>` landmarks and button focus styles on all seven pages.
- P04 — Example mini-apps: forms with Enter-to-submit, aria-live status, disabled button during requests.
- P05 — Shortcuts/theme polish: theme toggle state sync, modifier-key guards, visible theme toggle label.
- E1 — Skip link: added as the first focusable element in `Shell.tsx`; `main` is focusable via `tabIndex={-1}`.
- E2 — High contrast: added `prefers-contrast: more` media query with black-on-white tokens and thicker borders.
- E3 — Playwright E2E: `web/dashboard/e2e/dashboard-a11y.spec.ts` covers skip link, tab navigation, theme toggle, and critical axe-core violations.
- E4 — CI contrast check: `npm run a11y:contrast` runs axe-core color-contrast checks for light and dark themes in GitHub Actions.

## Commits

| Prompt | Commit | Message |
|--------|--------|---------|
| P01, P02, P05 | `acd8a9e` | feat(a11y): dashboard keyboard navigation, labels, live regions, theme sync |
| P03 | `4c392a5` | feat(a11y): add focus indicators and main landmarks to provider pages |
| P04 | `8bb9ee3` | feat(a11y): improve example mini-app accessibility |
| E1–E4 | `f066246` | feat(a11y): skip link, prefers-contrast, E2E a11y tests, CI contrast check |
| Docs | `f39a0f3` | docs(a11y): update tracker, handoff, and changelog for E1–E4 enhancements |
| Fix | `c12554d` | fix(a11y): wait for theme application in contrast check |

## Quality gate results

All passed:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `cd web/dashboard && npm run test:ci` (16/16)
- `cd web/dashboard && npm run test:e2e` (4/4)
- `cd web/dashboard && npm run a11y:contrast` (light + dark)
- `cd web/dashboard && npm run bundle-size`
- `./scripts/smoke-test.sh`

## Notes for future maintainers

- When adding new interactive table rows, use semantic `<button>` elements rather than `onClick` on `<tr>`.
- Use the `Announce` component for any new transient status messages.
- Keep provider pages self-contained; any shared CSS/JS should be duplicated to avoid extra HTTP requests.
- Run `npm run a11y:contrast` locally after changing dashboard colors to catch contrast regressions before CI.
- The E2E test can use an externally-started server via `MUARA_URL`; otherwise it builds and starts a temporary OpenMuara binary.
