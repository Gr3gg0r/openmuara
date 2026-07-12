> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Review Checklist

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — ready for sign-off

---

Use this checklist to sign off the accessibility and usability audit initiative.

## Standards & tooling

- [x] WCAG 2.1 Level AA target is documented.
- [x] Component and view inventory is complete.
- [x] `vitest-axe` is installed and configured.
- [x] `eslint-plugin-jsx-a11y` is installed and configured.
- [x] Manual screen-reader test guide exists.

## Automated scans

- [x] axe-core runs on every component and view.
- [x] Zero critical axe-core violations.
- [x] Zero serious axe-core violations.
- [x] Moderate violations are triaged and documented.

## Keyboard navigation

- [x] All interactive controls are reachable by Tab/Shift+Tab.
- [x] Focus order is logical.
- [x] Focus indicators are visible.
- [x] Command palette opens and closes via keyboard.
- [x] Dialogs trap focus and restore focus on close.
- [x] Skip-to-main-content link is present and works.

## Screen-reader support

- [x] All icon-only buttons have accessible names.
- [x] All form inputs have associated labels.
- [x] ARIA roles and properties are used correctly.
- [x] Live regions announce critical status changes.
- [x] Core views pass VoiceOver or NVDA smoke test.

## Color & contrast

- [x] Normal text meets 4.5:1 contrast in both themes.
- [x] Large text meets 3:1 contrast in both themes.
- [x] UI components and graphical objects meet 3:1 contrast.
- [x] Color is never the sole means of conveying information.
- [x] Focus indicators are visible in both themes.

## Motion & animation

- [x] Animations respect `prefers-reduced-motion`.
- [x] No auto-playing content that cannot be paused.

## Forms & errors

- [x] All inputs have accessible names.
- [x] Required fields are indicated accessibly.
- [x] Error messages are associated with inputs.
- [x] Errors are announced to screen readers.

## Mobile & responsive

- [x] Touch targets are at least 36×36 CSS pixels.
- [x] Primary touch targets are at least 44×44 CSS pixels where feasible.
- [x] Viewport scaling is not disabled.
- [x] Tables and forms are usable at 320 px width.

## Dynamic content

- [x] Dialogs manage focus on open and close.
- [x] Toasts and alerts are announced via live regions.
- [x] Route changes move focus predictably.
- [x] Loading states are announced where needed.

## Usability heuristics

- [x] Nielsen heuristics evaluated for key workflows.
- [x] Findings prioritized and documented.

## Documentation

- [x] `KNOWN_ISSUES.md` contains all findings and accepted deviations.
- [x] No undocumented deviation remains.
- [x] Manual test guide is published.

## CI & quality gates

- [x] A11y tests run in CI.
- [x] `npm run typecheck` passes.
- [x] `npm run test:ci` passes.
- [x] `npm run test:a11y` passes.
- [x] `npm run lint:a11y` passes.

## Sign-off

| Role | Name | Date | Signature |
|---|---|---|---|
| AI Agent | Kimi Code | | |
| Human Reviewer | | | |
| Maintainer | | | |
