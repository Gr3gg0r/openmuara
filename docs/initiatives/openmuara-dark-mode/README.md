> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dark Mode

> **Status:** ✅ COMPLETE | **Started:** 2026-07-03 | **Completed:** 2026-07-03
> **Scope:** Add a cohesive, accessible dark mode to the OpenMuara dashboard, provider pay pages, and example mini-apps, respecting OS preference and allowing a manual toggle.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/dark-mode`

---

## Initiative Structure

```
docs/initiatives/openmuara-dark-mode/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing gaps
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   ├── 01-dashboard-dark-mode.md
│   ├── 02-provider-pay-pages-dark-mode.md
│   ├── 03-example-mini-apps-dark-mode.md
│   └── 04-docs-release-notes.md
│
├── findings/              # Audit output
├── appendices/            # Reference material
│   ├── a-color-palette.md
│   ├── b-implementation-snippets.md
│   └── c-test-plan.md
│
├── GLOSSARY.md            # Shared terminology
└── .gitignore             # Ignore screenshots, logs
```

Planning docs live in `docs/initiatives/openmuara-dark-mode/` in the root repo. Product code commits to the `feat/dark-mode` branch. Do not commit directly to `main`.

---

## Why dark mode?

The OpenMuara dashboard and example landing pages currently ship with a light-only UI. Many developers run the emulator for long sessions or in dim environments, and a dark mode:

- Reduces eye strain during extended debugging.
- Matches the default appearance of most developer tools and terminal environments.
- Signals polish and attention to detail for public contributors.

Because OpenMuara is local-first and low-resource, the implementation must stay lightweight: CSS custom properties, a small theme script, and no large theming libraries.

---

## Goals

1. **Dashboard dark mode** — The Preact/Vite dashboard (`web/dashboard`) supports dark mode via `prefers-color-scheme` and a manual toggle. The choice is persisted in `localStorage` and synced across tabs.
2. **Provider pay pages** — The OpenMuara-hosted payment simulation pages (`internal/ui/*.html`, plus Go-rendered pages) follow the OS color scheme using CSS `color-scheme` and custom properties.
3. **Example mini-apps** — `examples/ecommerce-single-buy/` and `examples/prepaid-topup/` respect `prefers-color-scheme` with zero extra dependencies; no per-example toggle.
4. **Theme system** — A single, semantic CSS custom-property palette is shared where possible, with light and dark values, so future UI additions automatically support both modes.
5. **Accessibility** — All text and interactive elements meet WCAG AA contrast in both modes; focus indicators remain visible; color is not the only status indicator.
6. **No flash, no dependency** — A tiny blocking script sets the theme before first paint; no theming library is added.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. Backward compatibility
Dark mode is additive. Light mode must remain the default when no preference is set and JavaScript is disabled.

### 3. No external theming libraries
Use CSS custom properties and native `color-scheme`. Do not add Tailwind, Material-UI, or other UI frameworks just for theming.

### 4. Semantic design tokens
Name tokens by purpose, not by color value. Examples:
- `--color-bg` (not `--white`)
- `--color-surface` (not `--card`)
- `--color-text-primary` (not `--slate-900`)
- `--color-border` (not `--gray-200`)

### 5. Respect OS preference first
Default to the user's OS preference via `prefers-color-scheme`. A manual toggle overrides it and is persisted in `localStorage` under key `muara-theme`.

### 6. Low memory / low overhead
No runtime theme object, no Preact context for theme unless already present. Prefer a `<script>` that sets a class before first paint to avoid flashes.

### 7. Cross-tab sync
Listen to the `storage` event so toggling the theme in one tab updates all other OpenMuara dashboard tabs.

### 8. Accessibility
- WCAG AA contrast (4.5:1 for normal text, 3:1 for large text/UI components) in both modes.
- Visible focus indicators in both modes.
- `aria-label` on the theme toggle.
- Honor `prefers-reduced-motion` by disabling the color transition for users who request reduced motion.
- Respect `forced-colors` / Windows High Contrast Mode: do not rely solely on background-color for meaning; keep borders and focus rings visible.
- Validate contrast with browser DevTools, axe DevTools, or Lighthouse; target zero contrast errors.

### 9. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`
- `cd web/dashboard && npm run test:ci`
- Manual visual check in both light and dark modes

### 10. Definition of done
Beyond the quality gates, a prompt is done only when:

- The feature works end-to-end in both light and dark modes.
- Tests cover the theme toggle, OS-preference fallback, and persistence.
- The smoke test passes.
- `HANDOFF.md` is updated with what was built and what changed.
- `TRACKING.md` marks the prompt `✅` with the commit hash.
- User-facing changes are noted for the next release notes.

---

## Out of Scope

- Redesigning the dashboard layout or navigation.
- Adding animations or transitions beyond a simple color fade.
- Theming the provider emulation endpoints that return JSON only.
- Theming CLI terminal output.
- Native mobile apps or desktop installers.

---

## Success criteria

- Dashboard renders correctly in dark mode with no visual regressions in light mode.
- A user can toggle dark/light mode and the choice persists across reloads and syncs across tabs.
- OS preference is honored on first visit when no manual choice exists.
- Example mini-apps adapt to dark mode without extra dependencies.
- Provider pay pages do not flash white before rendering in dark mode.
- All interactive elements remain usable with keyboard and screen readers.
- Dark mode pages remain readable in print preview and under `forced-colors`.
- All quality gates pass.

## Metrics

| Metric | Target | How measured |
|--------|--------|--------------|
| No flash of wrong theme on reload | 100% | Manual test with throttled CPU |
| Light-mode visual regression | 0 | Manual comparison before/after |
| Bundle size increase | ≤5 KiB (gzipped) over current baseline (JS 11.83 KiB, total dist 158.83 KiB) | `npm run bundle-size` |
| Toggle persistence | 100% | Manual test across reloads |
| Cross-tab sync | 100% | Two tabs open, toggle in one |
| WCAG AA contrast | 100% of text/UI | Browser dev tools or axe |
| No-JS graceful degradation | Works | Disable JS, page still readable |

## Recommendations & future enhancements

These are not required for the v1 dark-mode feature but would raise the polish further. They are intentionally low-priority or out of scope for this initiative.

| # | Recommendation | Priority | Rationale |
|---|----------------|----------|-----------|
| E1 | Add a manual theme toggle to provider pay pages | Low | Currently OS-preference only. A toggle would help users whose OS is light but who want to test a dark pay page, but it adds per-page JS and complexity. |
| E2 | Add a manual theme toggle to example mini-apps | Low | Same trade-off as E1; examples are intentionally minimal. |
| E3 | Theme-aware favicon | Low | Browsers that support `media="(prefers-color-scheme: dark)"` on `<link rel="icon">` could show a dark favicon; current inline SVG favicon already adapts reasonably. |
| E4 | Use `accent-color` for native form controls | Low | Would tint checkboxes/radio buttons with the theme color; minor visual polish. |
| E5 | Automated visual regression tests for light/dark | Medium | Playwright screenshots of dashboard and pay pages in both modes would catch unintended drift. Worth adding once the UI stabilizes. |
| E6 | Automated contrast regression checks in CI | Medium | A small axe-core or Playwright step could fail builds that introduce contrast errors. |
| E7 | Theme preference scoped to OpenMuara workspace | Low | Today the choice is browser-global (`localStorage`). Per-workspace theming could be useful if a developer runs multiple OpenMuara instances. |
| E8 | CLI terminal theming | Low | Out of scope for v1; the CLI logs are structured JSON and follow terminal defaults. |
| E9 | Follow `prefers-contrast` for increased-contrast variants | Low | Gold-standard accessibility; would add a third set of tokens for high-contrast users. |

## Files changed (summary)

| Area | Key files |
|------|-----------|
| Dashboard SPA | `web/dashboard/src/styles.css`, `web/dashboard/src/theme.ts` (new), `web/dashboard/src/components/Shell.tsx`, `web/dashboard/src/app.tsx`, `web/dashboard/index.html`, `web/dashboard/tests/Shell.test.tsx`, `internal/ui/dashboard-dist/index.html` |
| Provider pages | `internal/ui/stripe-checkout.html`, `internal/ui/stripe-payment-intent.html`, `internal/ui/stripe-webhooks.html`, `internal/ui/fawry-escape.html`, `internal/ui/billplz-pay.html`, `internal/ui/toyyibpay-pay.html`, `internal/ui/ipay88-pay.html` |
| Examples | `examples/ecommerce-single-buy/index.html`, `examples/prepaid-topup/index.html` |
| Docs | `README.md`, `CHANGELOG.md` |

## Lessons learned

1. **One source of truth matters.** Using a single `data-theme` attribute avoids fights between `prefers-color-scheme` and a manual override.
2. **Blocking scripts beat flashes.** Setting the theme in a small inline `<head>` script is the cheapest way to prevent a flash of un-themed content.
3. **Provider pages can stay simple.** Transient simulation pages do not need a full theme system; `color-scheme: light dark` plus a few tokens is enough.
4. **Tests should assert behavior, not colors.** Dashboard tests check the `data-theme` attribute and toggle interaction, not specific hex values.

## Reference material

- `appendices/a-color-palette.md` — suggested semantic tokens with WCAG AA-checked values.
- `appendices/b-implementation-snippets.md` — copy-paste starting points for the blocking script, toggle helpers, pay pages, and examples.
- `appendices/c-test-plan.md` — manual QA checklist for every prompt.
