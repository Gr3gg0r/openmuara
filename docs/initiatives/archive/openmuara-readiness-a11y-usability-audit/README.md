> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit

> **Status:** ✅ Complete | **Started:** 2026-07-08 | **Completed:** 2026-07-09
> **Scope:** Make the OpenMuara dashboard demonstrably accessible and usable across devices, input methods, and assistive technologies before public release.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`

---

## Why this matters

OpenMuara is a developer tool used during payment integration work, often under pressure. A dashboard that is hard to use with a keyboard, screen reader, or mobile device signals low quality and excludes members of the community. Accessibility is also a hallmark of mature open-source projects and improves the experience for everyone.

This initiative treats accessibility and usability as first-class quality gates. We will audit the dashboard against WCAG 2.1 Level AA, fix violations, add regression tests, and document any intentional deviations.

## Initiative structure

```
docs/initiatives/openmuara-readiness-a11y-usability-audit/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── KNOWN_ISSUES.md        # Catalog of a11y/usability findings
├── RISKS.md               # Risk register
├── RECOMMENDATIONS.md     # Standards, tools, and prioritized actions
├── DECISIONS.md           # Decision log
├── EXECUTION_PLAN.md      # Milestones, dependencies, RACI
├── CI_INTEGRATION.md      # CI/workflow changes
├── REVIEW_CHECKLIST.md    # Sign-off checklist
├── ROLLBACK_PLAN.md       # Regression response plan
├── HANDOFF.md             # Final state and next steps
├── MANUAL_TESTING.md      # Screen-reader and keyboard test procedures
└── GLOSSARY.md            # Accessibility terms and definitions
```

## Audit areas

1. **Automated scans** — axe-core, Lighthouse, `eslint-plugin-jsx-a11y` on every dashboard view and component.
2. **Keyboard navigation** — focus order, focus indicators, Tab/Shift+Tab, Enter/Space, Escape, shortcuts, command palette, skip links.
3. **Screen-reader support** — semantic HTML, labels, `aria-*` attributes, roles, live regions, announcements, alternative text.
4. **Color & contrast** — WCAG 2.1 AA text contrast (4.5:1), UI component contrast (3:1), dark mode, focus rings, status colors.
5. **Motion & animation** — `prefers-reduced-motion`, no auto-playing content, no vestibular triggers.
6. **Forms & errors** — associated labels, error identification, required field indication, validation feedback.
7. **Mobile & responsive** — touch targets, viewport scaling, pinch zoom, responsive tables, readable font sizes.
8. **Cognitive & content** — clear language, consistent navigation, predictable behavior, headings hierarchy, breadcrumbs.
9. **Dynamic content** — focus management for dialogs, toasts, route changes, async updates, loading states.
10. **Usability heuristics** — Nielsen/Norman heuristics applied to key workflows (view ledger, replay webhook, configure provider).

## Standards mapping

| Standard | Level | How we apply it |
|---|---|---|
| WCAG 2.1 | Level AA | Primary target for all audit areas |
| WCAG 2.2 | Level AA | Adopt new success criteria where practical (focus appearance, drag alternatives, consistent help) |
| WAI-ARIA 1.2 | — | Authoring Practices for widgets (command palette, dialogs, tabs, tables) |
| EN 301 549 / Section 508 | — | Considered for public-sector adopters |
| Inclusive Design Principles | — | Guide usability enhancements beyond minimum compliance |

## Success criteria

- [x] No critical or serious automated a11y violations in any dashboard view.
- [x] All interactive controls are reachable and operable by keyboard alone.
- [x] All form inputs and icon-only controls have accessible names.
- [x] Color is never the sole means of conveying information.
- [x] Touch targets meet at least 36×36 CSS pixels (44×44 where feasible).
- [x] Both light and dark themes pass WCAG 2.1 AA contrast requirements.
- [x] Animations respect `prefers-reduced-motion`.
- [x] axe-core tests run in CI and block regressions.
- [x] Manual keyboard and screen-reader smoke tests are documented.
- [x] All findings are recorded in `KNOWN_ISSUES.md` with severity and remediation.

## Related documents

- [`TRACKING.md`](TRACKING.md) — phases, acceptance criteria, findings log
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — catalog of findings and deviations
- [`RISKS.md`](RISKS.md) — risk register
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — standards, tools, and prioritized actions
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones and dependencies
- [`CI_INTEGRATION.md`](CI_INTEGRATION.md) — concrete CI changes
- [`REVIEW_CHECKLIST.md`](REVIEW_CHECKLIST.md) — sign-off checklist
- [`ROLLBACK_PLAN.md`](ROLLBACK_PLAN.md) — regression response
- [`HANDOFF.md`](HANDOFF.md) — final state
- [`MANUAL_TESTING.md`](MANUAL_TESTING.md) — manual test procedures
- [`GLOSSARY.md`](GLOSSARY.md) — accessibility terms
