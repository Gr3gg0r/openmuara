> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dark Mode — Execution Tracker

> **Updated:** 2026-07-09 | **Status:** ✅ COMPLETE / Merged to `dev`
>
> **Scope:** Add a cohesive, accessible dark mode to the OpenMuara dashboard, provider pay pages, and example mini-apps, respecting OS preference and allowing a manual toggle.
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev` (feature branch merged and removed)

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
5. Product-code commits happen on `feat/dark-mode`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| P01 | Dashboard dark mode | `web/dashboard/src/styles.css`, `web/dashboard/src/components/Shell.tsx`, `web/dashboard/src/app.tsx`, `web/dashboard/index.html`, `web/dashboard/tests/`, `internal/ui/dashboard-dist/` (rebuild) | — | ✅ | d1d49dc | Semantic tokens, OS-preference fallback, manual toggle, localStorage, cross-tab sync, no flash, WCAG AA. |
| P02 | Provider pay pages dark mode | `internal/ui/stripe-checkout.html`, `internal/ui/stripe-payment-intent.html`, `internal/ui/stripe-webhooks.html`, `internal/ui/fawry-escape.html`, `internal/ui/billplz-pay.html`, `internal/ui/toyyibpay-pay.html`, `internal/ui/ipay88-pay.html`, Go-rendered pages if any | — | ✅ | 02fd044 | Follow `prefers-color-scheme`; convert hard-coded colors to custom properties; no flash. |
| P03 | Example mini-apps dark mode | `examples/ecommerce-single-buy/index.html`, `examples/prepaid-topup/index.html` | — | ✅ | 9a56724 | Respect OS preference with zero dependencies; no per-example toggle (see D006). |
| P04 | Docs and release notes | `README.md`, `docs/initiatives/openmuara-dark-mode/README.md`, `CHANGELOG.md` | P01–P03 | ✅ | 22fd9f1 | Document the toggle, keyboard shortcut, and how contributors should use theme tokens. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| UI Tests | `cd web/dashboard && npm run test:ci` | 16/16 pass | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |
| Bundle size | `cd web/dashboard && npm run bundle-size` | JS 12.25 KiB / total dist 169.28 KiB (within limits) | ✅ |
| Visual QA | Manual | Light & dark modes | ✅ |

---

## Decisions

- D001 ✅ Dark mode must be additive, use semantic CSS custom properties, and not pull in a theming library.
- D002 ✅ Manual choice is stored in `localStorage` key `muara-theme`; OS preference is the fallback.
- D003 ✅ Semantic design tokens named by purpose.
- D004 ✅ No manual toggle on provider pay pages (OS preference only).
- D005 ✅ `data-theme` attribute is the single source of truth.
- D006 ✅ No manual toggle in example mini-apps (OS preference only).
