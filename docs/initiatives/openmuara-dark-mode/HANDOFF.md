> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Dark Mode — Handoff

> **Last updated:** 2026-07-03
> **Initiative:** OpenMuara Dark Mode
> **Branch:** `feat/dark-mode` (merged into `dev` and deleted)
> **Status:** ✅ COMPLETE

---

## What was shipped

- P01 — Dashboard dark mode with semantic CSS tokens, manual toggle, OS-preference fallback, `localStorage` persistence, cross-tab sync, and `d` keyboard shortcut.
- P02 — Dark mode for all seven provider pay pages (`internal/ui/*.html`).
- P03 — Dark mode for `examples/ecommerce-single-buy/index.html` and `examples/prepaid-topup/index.html`.
- P04 — Updated `README.md` and `CHANGELOG.md`.

## Commits

| Prompt | Commit | Message |
|--------|--------|---------|
| P01 | `d1d49dc` | feat(dark-mode): add dashboard dark mode (P01) |
| P02 | `02fd044` | feat(dark-mode): add dark mode to provider pay pages (P02) |
| P03 | `9a56724` | feat(dark-mode): add dark mode to example mini-apps (P03) |
| P04 | `22fd9f1` | docs(dark-mode): document dark mode and update changelog (P04) |

## Key implementation details

- Theme storage key: `muara-theme`.
- Blocking inline script in `web/dashboard/index.html` and all provider/example pages sets `data-theme` before first paint.
- CSS custom properties use semantic names (`--color-bg`, `--color-text-primary`, etc.).
- `prefers-reduced-motion` guards the color transition.
- Bundle impact: JS +0.42 KiB, total dist +10.45 KiB (within limits).

## Quality gate results

All passed:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `cd web/dashboard && npm run test:ci` (16/16)
- `cd web/dashboard && npm run bundle-size`
- `./scripts/smoke-test.sh`

## Notes for future maintainers

- When adding new UI, use the semantic tokens in `web/dashboard/src/styles.css` or the provider-page token sets.
- Rebuild `internal/ui/dashboard-dist/` with `npm run build` after any dashboard change.
- Test both light and dark modes; verify WCAG AA contrast before committing.
