> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Rollback Plan

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — rollback plan published

---

This plan describes how to respond if an accessibility or usability change causes regression, CI instability, or incorrect behavior.

## 1. A11y test regression

**Scenario:** A new a11y test fails in CI after a dashboard change.

**Response:**
1. Determine if the failure is due to:
   - A real a11y regression.
   - A false positive from `vitest-axe` or `eslint-plugin-jsx-a11y`.
   - A flaky browser/screen-reader test.
2. If it is a real regression, fix the component and add a regression test.
3. If it is a false positive, document the exception in `KNOWN_ISSUES.md` and pin/suppress with a comment.
4. If it is flaky, stabilize the test before merging.

## 2. Visual regression from a11y fix

**Scenario:** A contrast or focus change degrades the visual design.

**Response:**
1. Revert the specific CSS change if the regression is worse than the original issue.
2. Work with the design owner to find a compliant alternative.
3. Add a visual regression test or screenshot comparison if feasible.

## 3. Scope creep into redesign

**Scenario:** An a11y finding suggests a larger UI redesign.

**Response:**
1. Document the finding in `KNOWN_ISSUES.md`.
2. Implement the minimum WCAG-compliant fix within the audit scope.
3. Move broader redesign ideas to a separate UX initiative.

## 4. Manual screen-reader test disagreement

**Scenario:** Different screen readers report different behavior for the same control.

**Response:**
1. Test with the agreed primary screen reader (VoiceOver on macOS for initial ship).
2. Document the discrepancy in `KNOWN_ISSUES.md`.
3. Prioritize fixes that improve the primary experience without breaking others.

## 5. CI performance regression

**Scenario:** A11y tests slow CI significantly.

**Response:**
1. Profile the slow tests.
2. Run component-level axe tests in parallel.
3. Move heavy end-to-end keyboard tests to a dedicated job.

## Communication template

For significant regressions, open a GitHub issue with:

```markdown
## A11y/usability regression: <area>

- **Introduced in:** commit/PR
- **Area:** 
- **WCAG criterion:** 
- **Expected behavior:** 
- **Actual behavior:** 
- **Impact:** 
- **Proposed fix:** 
```
