> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Accessibility & Usability Polish

> **Status:** ✅ COMPLETE | **Started:** 2026-07-03 | **Completed:** 2026-07-03
> **Scope:** Make the dashboard, provider simulation pages, and example mini-apps keyboard-friendly, screen-reader-friendly, and visually polished while keeping the bundle tiny.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/a11y-enhancements`

---

## Initiative Structure

```
docs/initiatives/openmuara-a11y-usability-polish/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing gaps from the audit
├── GLOSSARY.md            # Shared terminology
├── findings/              # Audit output
│   └── 2026-07-03-a11y-usability-audit.md
└── prompts/               # Numbered execution prompts
    ├── 01-dashboard-keyboard-navigation.md
    ├── 02-dashboard-labels-and-live-regions.md
    ├── 03-provider-pages-focus-and-landmarks.md
    ├── 04-example-apps-accessibility.md
    └── 05-shortcuts-and-theme-polish.md
```

Planning docs live in `docs/initiatives/openmuara-a11y-usability-polish/` in the root repo. Product code commits to the `feat/a11y-usability-polish` branch. Do not commit directly to `main` or `dev`.

---

## Why now?

OpenMuara's UI works well for mouse users, but an accessibility/usability audit found several keyboard and screen-reader blockers. Fixing them makes the tool easier for more developers to use and contributes to a polished, professional feel without adding heavy dependencies.

---

## Goals

1. **Keyboard parity:** every interactive element in the dashboard and provider pages is reachable and operable with a keyboard.
2. **Screen-reader clarity:** controls have meaningful labels, live regions announce state changes, and dialogs manage focus.
3. **Visual consistency:** focus indicators, color contrast, and dark-mode behavior are solid across all pages.
4. **No bundle bloat:** keep using the current Preact + vanilla CSS stack; no Tailwind or component libraries.
5. **Quality gates green:** every prompt passes build, test, lint, and smoke tests.

---

## Non-goals

- Redesign the entire UI.
- Add a component library (Daisy UI, shadcn/ui, etc.).
- Change provider emulation contracts or backend logic.

---

## Acceptance criteria for the initiative

- [x] All high-severity accessibility issues from the audit are fixed.
- [x] All medium-severity issues are fixed or explicitly accepted with a note.
- [x] Low-severity polish items are fixed where the cost is trivial.
- [x] Dashboard tests cover the new keyboard/focus behavior.
- [x] All quality gates pass:
  - [x] `go build ./...`
  - [x] `go test ./...`
  - [x] `go vet ./...`
  - [x] `golangci-lint run`
  - [x] `cd web/dashboard && npm run test:ci`
  - [x] `cd web/dashboard && npm run bundle-size`
  - [x] `./scripts/smoke-test.sh`
  - [x] `cd web/dashboard && npm run test:e2e`
  - [x] `cd web/dashboard && npm run a11y:contrast`
- [x] `TRACKING.md` and `HANDOFF.md` are updated.
- [x] Release-notes snippet added to `CHANGELOG.md`.

## Recommendations & future enhancements

| # | Recommendation | Status | Priority | Rationale |
|---|----------------|--------|----------|-----------|
| E1 | Add a "Skip to main content" link | ✅ Done | High | First item in tab order; lets keyboard users bypass the header and tabs. |
| E2 | Support `prefers-contrast: more` | ✅ Done | Medium | Adds a high-contrast token set for users who need stronger boundaries. |
| E3 | Playwright E2E accessibility smoke test | ✅ Done | Medium | Catches focus regressions and captures light/dark screenshots automatically. |
| E4 | Automated contrast regression check in CI | ✅ Done | Medium | Fail builds that introduce new WCAG contrast errors. |
| E5 | Heading hierarchy audit | ⬜ Pending | Low | Ensure every view has exactly one `<h1>` and logical heading order. |
| E6 | Persistent user preference for reduced motion | ⬜ Pending | Low | Today we only honor OS reduced-motion; a manual toggle could help testing. |
