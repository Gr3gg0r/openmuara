> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Recommendations

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — recommendations implemented

---

These recommendations define a gold-standard accessibility and usability framework for the OpenMuara dashboard.

## Target standard

**WCAG 2.1 Level AA** is the minimum ship target. Where practical, adopt WCAG 2.2 criteria such as:

- 2.4.11 Focus Not Obscured (AA)
- 2.4.12 Focus Not Obscured (Enhanced) (AAA)
- 2.5.7 Dragging Movements (AA)
- 3.2.6 Consistent Help (A)
- 3.3.7 Redundant Entry (A)
- 3.3.8 Accessible Authentication (Minimum) (AA)

## Current state assessment

Based on a quick review of `web/dashboard/src/`:

| Area | State | Notes |
|---|---|---|
| Automated scans | ⬜ Not integrated | No axe-core or Lighthouse in `package.json` |
| Keyboard navigation | partial | `CommandPalette.tsx`, `useFocusTrap.ts` exist; full coverage unknown |
| Screen-reader support | partial | `Announce.tsx` live region exists; icon-only button labels unknown |
| Color & contrast | partial | Dark mode in `theme.ts`; no systematic contrast audit |
| Motion | unknown | No `prefers-reduced-motion` audit |
| Forms | partial | `Input.tsx`, `Select.tsx` exist; error association unknown |
| Mobile | partial | Responsive layout exists; touch targets not audited |
| Dynamic content | partial | Focus trap exists; route/dialog focus management unknown |
| Usability heuristics | ⬜ Not done | No heuristic evaluation yet |

## Priority matrix

| Priority | Area | Recommendation | Effort | Impact |
|---|---|---|---|---|
| P0 | Automated scans | Add `vitest-axe` or `@axe-core/react` and run scans per component/view | S | High |
| P0 | Icon-only buttons | Audit and label every icon-only control (`CopyButton`, `SidebarNav`, etc.) | S | High |
| P0 | Focus management | Standardize focus traps, focus restoration, and visible focus indicators | M | High |
| P1 | Keyboard workflows | Add Playwright tests for primary workflows without a mouse | M | High |
| P1 | Form labels & errors | Ensure every input has a label and error messages are associated/announced | M | High |
| P1 | Contrast audit | Check all text and UI components in light and dark themes | M | High |
| P2 | Screen-reader smoke tests | Document and run VoiceOver (macOS) or NVDA (Windows) tests | M | Medium |
| P2 | Mobile touch targets | Enforce min 36×36 px; prefer 44×44 px for primary actions | S | Medium |
| P2 | Reduced motion | Audit animations and add `prefers-reduced-motion` fallbacks | S | Medium |
| P2 | Skip link | Add a skip-to-main-content link | S | Medium |
| P3 | Usability heuristics | Run Nielsen heuristic evaluation on key workflows | M | Medium |
| P3 | A11y statement | Publish an accessibility statement in the docs site | S | Low |

## Recommended tool stack

| Purpose | Tool | Where |
|---|---|---|
| Component-level a11y tests | `vitest-axe` + `@testing-library/react` | `web/dashboard/tests/*.test.tsx` |
| Linting | `eslint-plugin-jsx-a11y` | `web/dashboard/eslint.config.*` |
| End-to-end keyboard tests | Playwright | `web/dashboard/tests/e2e/` |
| Visual regression | Lighthouse CI or Playwright screenshots | CI |
| Contrast checking | `axe-core` or browser devtools | Local/CI |
| Screen-reader manual tests | VoiceOver (macOS), NVDA (Windows) | Local |
| Motion audit | Browser devtools + CSS media query tests | Local/CI |

## Component-level test pattern

```tsx
import { render } from '@testing-library/react';
import { axe, toHaveNoViolations } from 'vitest-axe';
import { Button } from '../src/components/Button';

expect.extend(toHaveNoViolations);

it('has no detectable a11y violations', async () => {
  const { container } = render(<Button>Save</Button>);
  const results = await axe(container);
  expect(results).toHaveNoViolations();
});
```

## End-to-end keyboard workflow pattern

```ts
import { test, expect } from '@playwright/test';

test('command palette opens with Ctrl+K and is keyboard operable', async ({ page }) => {
  await page.goto('/');
  await page.keyboard.press('Control+k');
  await expect(page.getByRole('dialog', { name: 'Command palette' })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('dialog', { name: 'Command palette' })).not.toBeVisible();
});
```

## Standards mapping

| Recommendation | OpenSSF Scorecard | SLSA | CNCF |
|---|---|---|---|
| Automated a11y regression | — | — | Best Practice |
| Keyboard & screen-reader tests | — | — | Best Practice |
| Documented limitations | — | — | Best Practice |
| Inclusive design heuristics | — | — | Best Practice |

## WCAG 2.1 criterion mapping by dashboard pattern

| Pattern | Relevant criteria | Test approach |
|---|---|---|
| Page structure / landmarks | 1.3.1 Info and Relationships, 2.4.1 Bypass Blocks | axe-core + manual landmark review |
| Buttons and links | 2.1.1 Keyboard, 2.4.3 Focus Order, 2.4.7 Focus Visible, 4.1.2 Name, Role, Value | Keyboard test + axe-core |
| Command palette / dialogs | 1.3.1, 2.1.1, 2.1.2 No Keyboard Trap, 2.4.3, 2.4.7, 4.1.2, 4.1.3 Status Messages | Keyboard + focus trap tests + screen reader |
| Forms and inputs | 1.3.1, 1.3.5 Identify Input Purpose, 3.3.1 Error Identification, 3.3.2 Labels or Instructions, 3.3.3 Error Suggestion, 4.1.2 | axe-core + form tests + screen reader |
| Tables and data grids | 1.3.1, 1.3.2 Meaningful Sequence | axe-core + screen reader |
| Status badges and alerts | 1.4.1 Use of Color, 4.1.3 Status Messages | Visual review + `Announce.tsx` tests |
| Color and contrast | 1.4.3 Contrast (Minimum), 1.4.11 Non-text Contrast | axe-core + Lighthouse |
| Touch targets / mobile | 2.5.5 Target Size (AAA), best-practice 36×36 px minimum | CSS/devtools audit |
| Motion and animation | 2.2.2 Pause, Stop, Hide, 2.3.3 Animation from Interactions | `prefers-reduced-motion` tests |
| Notifications / toasts | 4.1.3 Status Messages | Live-region tests |

## What not to do

- Do **not** chase AAA compliance across the board; AA is the ship target.
- Do **not** redesign the UI under the guise of an audit; keep changes minimal and scoped.
- Do **not** rely solely on automated scans; manual keyboard and screen-reader tests are required.
- Do **not** leave deviations undocumented.

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — findings and deviations
- [`RISKS.md`](RISKS.md) — risk register
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones
