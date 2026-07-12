> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Execution Plan

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all milestones delivered

---

## Goal

Make the OpenMuara dashboard accessible and usable for the widest possible audience by auditing against WCAG 2.1 Level AA, fixing violations, adding regression tests, and documenting intentional deviations.

## Exit criteria

1. WCAG 2.1 Level AA target is documented and agreed.
2. Automated a11y scans (axe-core) run in CI with zero critical/serious violations.
3. Keyboard navigation is verified for all primary workflows.
4. Screen-reader smoke tests are documented and pass on one screen reader.
5. Light and dark themes pass contrast requirements.
6. Mobile touch targets meet the agreed minimum.
7. All findings are recorded in `KNOWN_ISSUES.md`.
8. All quality gates pass.

## Milestones

### M1 — Standards & tooling (P01)

**Deliverables**
- Confirm WCAG 2.1 AA target (with optional WCAG 2.2 enhancements).
- Install and configure `vitest-axe` and `eslint-plugin-jsx-a11y`.
- Create component and view inventory.

**Acceptance**
- `npm run typecheck` passes after adding dependencies.
- A sample axe test passes in CI.
- Inventory lists all 25 components and 8 views.

### M2 — Automated baseline scan (P02)

**Deliverables**
- Run axe-core on every component and view.
- Run Lighthouse a11y audit on key views.
- Populate `KNOWN_ISSUES.md` with baseline findings.

**Acceptance**
- Baseline report classifies every finding by severity and WCAG criterion.
- No unclassified critical/serious issues remain.

### M3 — Keyboard navigation (P03)

**Deliverables**
- Audit focus order and focus indicators across all views.
- Verify command palette, dialogs, and sidebar are keyboard-operable.
- Add Playwright keyboard regression tests for primary workflows.

**Acceptance**
- All primary workflows pass without a mouse.
- Keyboard tests run in CI.

### M4 — Screen-reader support (P04)

**Deliverables**
- Audit labels, ARIA roles, and live regions.
- Ensure icon-only buttons have accessible names.
- Document VoiceOver smoke-test steps.

**Acceptance**
- Core views pass VoiceOver smoke test documented in [`MANUAL_TESTING.md`](MANUAL_TESTING.md).
- `Announce.tsx` covers all critical status changes.

### M5 — Color & contrast (P05)

**Deliverables**
- Audit all text and UI component contrast in both themes.
- Remove color-only information cues.
- Ensure focus indicators are visible.

**Acceptance**
- Zero contrast failures in axe-core.
- Every status indicator has a non-color equivalent (text, icon shape, label).

### M6 — Motion & animation (P06)

**Deliverables**
- Audit all animations and transitions.
- Add `prefers-reduced-motion` fallbacks.
- Remove or pause auto-playing motion.

**Acceptance**
- No motion violations in axe-core.
- Reduced-motion tests pass.

### M7 — Forms & errors (P07)

**Deliverables**
- Ensure all inputs have associated labels.
- Improve error messaging and association.
- Announce validation errors via live region.

**Acceptance**
- 100% of inputs have accessible names.
- Form error states are keyboard-focusable and screen-reader announced.

### M8 — Mobile & responsive (P08)

**Deliverables**
- Audit touch targets across breakpoints.
- Ensure viewport scaling and pinch zoom are not blocked.
- Improve responsive tables and font sizes.

**Acceptance**
- Zero touch targets below 36×36 px.
- Viewport meta tag allows user scaling.

### M9 — Dynamic content (P09)

**Deliverables**
- Standardize focus management for dialogs, toasts, and route changes.
- Restore focus after modal close.
- Manage focus on route navigation.

**Acceptance**
- Focus behavior is predictable in all dynamic UI patterns.
- Regression tests cover focus changes.

### M10 — Usability heuristics (P10)

**Deliverables**
- Apply Nielsen heuristics to key workflows.
- Prioritize fixes that reduce friction.
- Add skip-to-main-content link.

**Acceptance**
- Heuristic report with prioritized recommendations.
- Skip link is keyboard-reachable.

### M11 — Documentation & limitation registry (P11)

**Deliverables**
- Finalize `KNOWN_ISSUES.md` with all findings and accepted deviations.
- Finalize [`MANUAL_TESTING.md`](MANUAL_TESTING.md).
- Update provider/dashboard docs with a11y notes.
- Prepare outreach template for users with disabilities (L6 validation).

**Acceptance**
- No undocumented deviations remain.
- Manual test guide is in `runbooks/` or `docs/`.
- Outreach template is ready to send.

**Outreach template for users with disabilities**

```markdown
Subject: Help us make OpenMuara accessible

Hi,

OpenMuara is an open-source local payment emulator. We are preparing for our v1.0 release and want to ensure the dashboard is accessible to everyone.

If you use a screen reader, keyboard navigation, or other assistive technology, we would love your feedback on the dashboard at:
http://localhost:9000/_admin

We are especially interested in:
1. Whether you can navigate the sidebar and command palette with a keyboard.
2. Whether form labels and error messages are clear.
3. Whether status badges and alerts make sense without relying on color.
4. Any pain points or blockers you encounter.

You can reply to this email or open an issue at https://github.com/openmuara/openmuara/issues.

Thank you for helping us improve OpenMuara.
```

### M12 — CI enforcement (P12)

**Deliverables**
- Add `npm run test:a11y` to CI.
- Block merges on new critical/serious violations.
- Document how to update baselines.

**Acceptance**
- CI fails when a PR introduces a critical/serious a11y violation.
- README explains how to run and update a11y tests.

## Timeline

Estimated calendar assuming one focused sprint (approximately 3 weeks). Adjust based on findings.

| Week | Milestones | Key deliverables |
|---|---|---|
| Week 1 | M1, M2 | Standards confirmed; `vitest-axe` installed; baseline scan complete; `KNOWN_ISSUES.md` populated |
| Week 2 | M3–M8 (in parallel) | Keyboard fixes; screen-reader labels; contrast fixes; reduced-motion; form and mobile fixes |
| Week 3 | M9–M12 | Dynamic focus management; heuristic evaluation; manual test guide; CI enforcement; final sign-off |

## Dependencies

- M1 must complete before M2–M12.
- M2 informs M3–M10 prioritization.
- M3–M10 can run in parallel once M2 is done.
- M11 depends on M2–M10.
- M12 depends on M1–M11.

## RACI

| Activity | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| Define standards | R | A | C |
| Install tooling | R | A | C |
| Run baseline scans | R | A | C |
| Fix keyboard issues | R | A | C |
| Fix screen-reader issues | R | A | C |
| Fix contrast/color issues | R | A | C |
| Review design impact | C | A | R |
| Write manual test guide | R | A | C |
| Approve CI changes | R | A | C |
| Final sign-off | C | A | R |

*R = Responsible, A = Accountable, C = Consulted, I = Informed*

## Rollback plan

- If a change causes visual regression, revert the specific CSS/JS change.
- If the a11y CI gate becomes flaky, pin the browser/tool version and re-run.
- If scope expands beyond audit, defer non-WCAG items to a separate UX initiative.

## Definition of done

- All phases in `TRACKING.md` marked ✅.
- Zero critical/serious axe-core violations.
- Keyboard and screen-reader smoke tests pass.
- `KNOWN_ISSUES.md` reviewed and approved.
- CI a11y gate passing.
- `HANDOFF.md` updated with final state.
