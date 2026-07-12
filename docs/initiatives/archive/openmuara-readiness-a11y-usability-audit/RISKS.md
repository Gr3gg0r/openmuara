> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — risks treated and residual risks documented

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|---|---|---|---|---|---|
| R01 | Custom components (command palette, dialogs, tables) bypass native a11y | High | High | Use semantic HTML, WAI-ARIA Authoring Practices, and test with keyboard/screen reader | AI Agent |
| R02 | Automated a11y scans miss issues that only manual testing catches | High | Medium | Combine automated scans with manual keyboard and screen-reader smoke tests | AI Agent |
| R03 | Dark mode contrast issues are overlooked during light-mode fixes | Medium | Medium | Audit both themes side-by-side; use design tokens and contrast automation | AI Agent |
| R04 | Mobile layout breaks touch targets or zoom behavior | Medium | Medium | Test on small viewports; enforce min touch-target sizes in CSS | AI Agent |
| R05 | Focus management in dynamic content (dialogs, toasts, route changes) is inconsistent | Medium | High | Define focus-handling patterns; add regression tests | AI Agent |
| R06 | Icon-only buttons and complex controls lack accessible names | High | High | Audit all icon usage; require `aria-label` or visible text | AI Agent |
| R07 | Animations trigger vestibular disorders or break reduced-motion settings | Low | High | Audit `prefers-reduced-motion`; avoid auto-playing motion | AI Agent |
| R08 | A11y fixes degrade visual design or developer experience | Medium | Low | Pair a11y changes with design review; keep changes minimal | Maintainer |
| R09 | No access to users with disabilities for validation | Medium | Medium | Document the gap; use automated tools + screen-reader manual tests as proxy | Maintainer |
| R10 | Scope expands to redesign rather than audit | Medium | Medium | Bound the initiative to WCAG 2.1 AA + heuristic fixes; defer redesign | Maintainer |
| R11 | A11y CI gate becomes flaky due to timing or environment differences | Low | Medium | Run axe-core in deterministic test render; pin browser version | AI Agent |
| R12 | Third-party chart or code-editor components are inaccessible | Medium | Medium | Audit dependencies; provide accessible alternatives or labels | AI Agent |

## Risk treatment summary

- **Accept:** R09 (documented and bounded).
- **Mitigate:** R01, R02, R03, R04, R05, R06, R07, R08, R10, R11, R12.
- **Transfer:** None.
- **Avoid:** None.

## Residual risks

| ID | Residual risk | Owner | Monitoring |
|---|---|---|---|
| R09 | No users with disabilities validate the dashboard before release | Maintainer | Add outreach task to roadmap; revisit after v1.0 |
| R10 | Heuristic improvements may expand scope | Maintainer | Review each heuristic finding against WCAG 2.1 AA before fixing |
| R12 | Inaccessible third-party dependencies remain | AI Agent / Maintainer | Track upstream issues; document workarounds |
