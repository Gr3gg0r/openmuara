> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Tracking

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all phases delivered and quality gates passed

---

## Exit criteria

1. WCAG 2.1 Level AA compliance is the documented target for all dashboard views.
2. Automated a11y scans (axe-core) run in CI and pass with zero critical/serious violations.
3. Keyboard navigation is verified end-to-end for all primary workflows.
4. Screen-reader smoke tests are documented and pass on at least one screen reader.
5. Light and dark themes pass contrast requirements.
6. Mobile touch targets meet the 36×36 px minimum (44×44 px where feasible).
7. All findings are recorded in `KNOWN_ISSUES.md` with severity and remediation plan.
8. All quality gates pass.

## Metrics and targets

| Metric | Target | Measurement |
|---|---|---|
| Critical axe-core violations | 0 | `npm run test:a11y` in CI |
| Serious axe-core violations | 0 | `npm run test:a11y` in CI |
| Moderate axe-core violations | ≤ 5 with accepted rationale | `npm run test:a11y` |
| Keyboard-reachable interactive controls | 100% | Manual smoke test + Playwright tests |
| Form inputs with accessible names | 100% | axe-core + manual review |
| Color-only information cues | 0 | Manual design review |
| Touch targets < 36×36 px | 0 | CSS audit + browser devtools |
| Lighthouse a11y score | ≥ 95 | Lighthouse CI or local run |
| Views covered by a11y tests | 100% (8 views) | `web/dashboard/tests/` review |
| Components covered by a11y tests | 100% (25 components) | `web/dashboard/src/components/` review |

## Audit maturity model

| Level | Name | What it covers | Enforcement |
|---|---|---|---|
| L0 | Awareness | A11y and usability are recognized as release blockers | Initiative docs |
| L1 | Automated regression | axe-core / Lighthouse run in CI | CI gate |
| L2 | Keyboard complete | All primary workflows work without a mouse | Playwright + manual tests |
| L3 | Screen-reader viable | Core views usable with a screen reader | Manual smoke tests |
| L4 | Visual & motion polish | Contrast, focus, reduced motion meet WCAG | Design + CSS review |
| L5 | Usability heuristics | Nielsen heuristics applied; friction removed | Expert review + user feedback |
| L6 | Community validation | Real users with disabilities provide feedback | Tracked outreach |

## Phases

| Phase | Title | Goal | Acceptance criteria | Effort | Status |
|-------|-------|------|---------------------|--------|--------|
| P01 | Standards & tooling | Define WCAG target, tool stack, and component inventory | Standards doc merged; axe-core installed; baseline scans run | S | ✅ Complete |
| P02 | Automated scan baseline | Run axe-core and Lighthouse on all views; populate findings | Baseline report in `KNOWN_ISSUES.md`; no unclassified critical/serious issues | M | ✅ Complete |
| P03 | Keyboard navigation | Audit and fix focus order, indicators, shortcuts, command palette | All primary workflows pass keyboard smoke test | L | ✅ Complete |
| P04 | Screen-reader support | Add labels, ARIA roles, live regions, announcements | Core views pass VoiceOver/NVDA smoke test | L | ✅ Complete |
| P05 | Color & contrast | Fix contrast in both themes; remove color-only cues | All text/UI components meet WCAG 2.1 AA | M | ✅ Complete |
| P06 | Motion & animation | Respect `prefers-reduced-motion`; audit animations | No violations; reduced-motion tests pass | S | ✅ Complete |
| P07 | Forms & errors | Audit labels, error messaging, validation | 100% of inputs have accessible names; errors announced | M | ✅ Complete |
| P08 | Mobile & responsive | Touch targets, viewport, responsive tables, font sizes | Zero touch targets below 36×36 px | M | ✅ Complete |
| P09 | Dynamic content | Focus management for dialogs, toasts, route changes | Focus is predictable after state changes | M | ✅ Complete |
| P10 | Usability heuristics | Nielsen/Norman heuristic evaluation of key workflows | Heuristic report with prioritized fixes | M | ✅ Complete |
| P11 | Documentation & limitation registry | Document deviations and manual test procedures | `KNOWN_ISSUES.md` complete; manual test guide exists | S | ✅ Complete |
| P12 | CI enforcement & regression | Add a11y and keyboard regression gates to CI | CI fails on new critical/serious violations | S | ✅ Complete |

## Component inventory

Dashboard components to audit:

- `Announce.tsx` — live region announcements
- `AppShell.tsx` — application landmark structure
- `Badge.tsx` — status indicators
- `Button.tsx` — buttons, focus, disabled states
- `Card.tsx` — containers
- `CodeBlock.tsx` — code content
- `CommandPalette.tsx` — dialog, search, keyboard shortcuts
- `ConfirmDialog.tsx` — modal dialog, focus trap
- `CopyButton.tsx` — icon-only button
- `DetailField.tsx` — read-only data
- `EmptyState.tsx` — empty state messaging
- `ErrorBoundary.tsx` — error messaging
- `FailedWebhookAlert.tsx` — alert, status
- `FilterChip.tsx`, `FilterToolbar.tsx` — filtering controls
- `Icon.tsx` — decorative vs. semantic icons
- `Input.tsx`, `Select.tsx` — form controls
- `Overview.tsx`, `Providers.tsx` — dashboard content
- `Shell.tsx`, `SidebarNav.tsx` — navigation, landmarks
- `Skeleton.tsx` — loading state
- `Timeline.tsx` — time-ordered content
- `WebhookConfig.tsx` — configuration form

## View inventory

Dashboard views to audit:

- `Ledger.tsx`
- `LedgerDetail.tsx`
- `ProviderDetail.tsx`
- `Settings.tsx`
- `Transactions.tsx`
- `WebhookDetail.tsx`
- `Webhooks.tsx`
- Command palette (global overlay)

## Per-component audit checklist

Use this checklist during manual review. A component is "audit-ready" when all items are checked.

| # | Check | Applies to |
|---|---|---|
| 1 | Has a logical focus order and visible focus indicator | All interactive components |
| 2 | Has an accessible name (visible text, `aria-label`, or `aria-labelledby`) | Buttons, links, inputs, icon-only controls |
| 3 | Uses semantic HTML before adding ARIA | All components |
| 4 | Does not rely on color alone to convey state | Badges, status indicators, charts |
| 5 | Meets contrast requirements in light and dark themes | Text, icons, borders, focus rings |
| 6 | Respects `prefers-reduced-motion` | Animated or transitioning components |
| 7 | Has a touch target ≥ 36×36 px (≥ 44×44 px for primary actions) | Buttons, links, tabs, chips |
| 8 | Manages focus correctly when opening/closing or updating | Dialogs, palettes, toasts, drawers |
| 9 | Announces dynamic changes via live region where needed | Alerts, toasts, async results |
| 10 | Has a matching visible label for every form control | Inputs, selects, textareas |

## Baseline matrix

Final state after the audit:

| Area | Final state | Notes |
|---|---|---|
| Automated scans | ✅ L1 regression gate | `vitest-axe`, `@axe-core/playwright`, and `eslint-plugin-jsx-a11y` run in CI |
| Keyboard nav | ✅ L2 complete | Skip link, command palette, sidebar, dialogs, and primary workflows are keyboard-operable |
| Screen reader | ✅ L3 viable | Dialog roles, labels, live regions, and announcements verified |
| Color & contrast | ✅ L4 WCAG 2.1 AA | Both themes pass axe-core contrast rules; no color-only cues |
| Motion | ✅ L4 reduced-motion safe | Animations respect `prefers-reduced-motion` |
| Forms | ✅ L4 accessible | All inputs have associated labels; errors are associated and announced |
| Mobile | ✅ L4 responsive | Touch targets meet 36×36 px minimum; viewport scaling allowed |
| Dynamic content | ✅ L4 predictable | Dialogs manage focus; route changes and notifications are announced |
| Usability heuristics | ✅ L5 applied | Key workflows reviewed; findings recorded and fixed |

## Target matrix

State at initiative completion:

| Area | Target |
|---|---|
| Automated scans | L1 — CI gate with zero critical/serious violations |
| Keyboard nav | L2 — all primary workflows keyboard-operable |
| Screen reader | L3 — core views screen-reader viable |
| Color & contrast | L4 — WCAG 2.1 AA in both themes |
| Motion | L4 — `prefers-reduced-motion` respected |
| Forms | L4 — fully labeled and error-accessible |
| Mobile | L4 — touch targets ≥ 36×36 px |
| Dynamic content | L4 — predictable focus management |
| Usability heuristics | L5 — heuristic evaluation complete |

## Findings log

| ID | Finding | Area | Severity | Status | Fixed in |
|----|---------|------|----------|--------|----------|
| A001 | Command palette container had no dialog role or accessible name | Screen reader / Dynamic content | Serious | ✅ Fixed | `CommandPalette.tsx` |
| A002 | Confirm dialog backdrop used `role="presentation"` and closed on any inner click | Keyboard / Dynamic content | Serious | ✅ Fixed | `ConfirmDialog.tsx` |
| A003 | Form labels in `WebhookConfig.tsx` and `ProviderDetail.tsx` used JSX `for=` instead of `htmlFor=` | Forms / Screen reader | Serious | ✅ Fixed | `WebhookConfig.tsx`, `ProviderDetail.tsx` |
| A004 | Timeline `<ol>` had redundant `role="list"` | Screen reader | Minor | ✅ Fixed | `Timeline.tsx` |
| A005 | Unused `Shell.tsx` component and dead test added maintenance noise | Usability heuristic | Minor | ✅ Fixed | Deleted `Shell.tsx` and `tests/Shell.test.tsx` |
| A006 | `Providers.tsx` used `dangerouslySetInnerHTML` for provider display text | Usability / Security | Moderate | ✅ Fixed | `Providers.tsx` |

*All findings were triaged, fixed, and verified by automated scans and Playwright tests.*

## Quality gates

Every phase must end with:

- [x] `go build ./...`
- [x] `go test ./...`
- [x] `go vet ./...`
- [x] `golangci-lint run`
- [x] `npm run typecheck` (in `web/dashboard/`)
- [x] `npm run test:ci` (in `web/dashboard/`)
- [x] `npm run test:a11y` (in `web/dashboard/`)
- [x] `npm run lint:a11y` (in `web/dashboard/`) — added during P12

## Definition of Ready per phase

| Phase | Ready when |
|---|---|
| P01 | Standards and tool stack agreed; component/view inventory complete |
| P02 | axe-core can run against rendered views; baseline report template ready |
| P03 | Primary user journeys documented; keyboard test script ready |
| P04 | Screen reader test environment available; annotation guide ready |
| P05 | Design tokens for light/dark themes identified; contrast checker chosen |
| P06 | Animation inventory complete; reduced-motion test plan ready |
| P07 | All form views identified; error-state audit template ready |
| P08 | Breakpoints and smallest target viewport agreed |
| P09 | Dynamic UI inventory (dialogs, toasts, routes) complete |
| P10 | Heuristic evaluation checklist selected |
| P11 | All findings triaged and documented |
| P12 | CI can run headless browser and axe-core |

## Notes

- Capture before/after screenshots and screen-reader recordings where possible.
- Prioritize findings that block keyboard or screen-reader users.
- Treat undocumented deviations as bugs.
- Use real users with disabilities for L6 validation if possible; otherwise document the attempt.
- Manual test procedures are documented in [`MANUAL_TESTING.md`](MANUAL_TESTING.md).
- Accessibility terms are defined in [`GLOSSARY.md`](GLOSSARY.md).
