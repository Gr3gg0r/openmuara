> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Accessibility & Usability Audit Decision Log

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — decisions accepted and recorded

---

## Accepted decisions

| ID | Decision | Context | Rationale | Status |
|----|----------|---------|-----------|--------|
| D01 | Target WCAG 2.1 Level AA | Need a clear, achievable compliance target | AA is the legal and industry baseline for most jurisdictions; AAA is aspirational and often impractical | ✅ Accepted |
| D02 | Use axe-core as the primary automated scanner | Widely adopted, integrates with React testing libraries | Provides consistent, deterministic results in CI | ✅ Accepted |
| D03 | Add component-level a11y tests in `vitest` | Dashboard already uses Vitest | Low-friction integration; fast feedback | ✅ Accepted |
| D04 | Document manual screen-reader smoke tests | No budget for professional a11y audit yet | Structured manual tests are the best available proxy | ✅ Accepted |
| D05 | Prioritize keyboard and screen-reader blockers | These user groups are most affected by custom components | Fixes here unlock the broadest set of workflows | ✅ Accepted |

## Resolved decisions

| ID | Question | Decision | Rationale | Owner | Date |
|----|----------|----------|-----------|-------|------|
| OD01 | Should we adopt WCAG 2.2 criteria as well? | Defer | WCAG 2.1 AA is the ship target; WCAG 2.2 focus-appearance and consistent-help will be revisited after v1.0 | Maintainer | 2026-07-09 |
| OD02 | Should we add a standalone accessibility statement page? | Defer | Add a concise statement on the docs site after the public release | Maintainer | 2026-07-09 |
| OD03 | Should we support high-contrast mode (`forced-colors`)? | Best-effort | Focus indicators and buttons remain visible; full `forced-colors` optimization deferred | Maintainer | 2026-07-09 |
| OD04 | Should we commission a professional a11y audit? | Post-v1.0 | Commission if community adoption justifies the cost | Maintainer | 2026-07-09 |

## Decisions requiring maintainer sign-off

| ID | Question | Options | Recommended | Owner | Due |
|----|----------|---------|-------------|-------|-----|
| DS01 | Confirm WCAG target level | AA / AAA | AA | Maintainer | Before P01 |
| DS02 | Confirm minimum browser/screen-reader matrix | Latest Chrome/Safari/Firefox + VoiceOver/NVDA | Latest Chrome + VoiceOver for initial ship; expand based on feedback | Maintainer | Before P03 |
| DS03 | Confirm touch-target minimum | 36×36 px / 44×44 px | 36×36 px minimum; 44×44 px for primary actions | Maintainer | Before P08 |

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — findings and deviations
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — standards and tools
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones
