> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Bug Hunt — Decision Log

> **Updated:** 2026-07-06

| ID | Date | Decision | Rationale | Status |
|----|------|----------|-----------|--------|
| D001 | 2026-07-06 | Bugs classified as P0/P1/P2 by severity. | Keeps triage consistent and prevents low-priority polish from blocking high-priority fixes. | Decided |
| D002 | 2026-07-06 | Every code fix includes a regression test. | Prevents regressions and provides proof the bug is fixed. | Decided |
| D003 | 2026-07-06 | Fixes batched by area/severity, not discovery order. | Reduces context switching and keeps commits focused. | Decided |
| D004 | 2026-07-06 | No speculative refactors; minimal fix only. | Avoids scope creep and reduces risk of introducing new bugs. | Decided |
| D005 | 2026-07-06 | P0/P1 integration fixes require explicit user sign-off before implementation. | `AGENTS.md` autonomy boundaries protect provider logic, signatures, config/auth/PII, and schema contracts. | Decided |
| D006 | 2026-07-06 | Coverage must not drop on any module changed by a fix. | Ensures fixes are tested and the codebase does not silently lose coverage. | Decided |
| D007 | 2026-07-06 | Dashboard redesign invariants are protected; any proposed invariant change requires user sign-off. | The Mailpit-style layout, filters, detail pages, provider settings, and dual-port runtime are user-facing contracts. | Decided |
| D008 | 2026-07-06 | Visual sign-off with Playwright MCP (P06) is required before PR. | Confirms UI/UX alignment with the user's priority stack and catches visual regressions that automated gates miss. | Decided |
| D009 | 2026-07-06 | Deferred P0/P1 bugs require user sign-off and a target release. | Prevents high-severity issues from being forgotten. | Decided |

---

## Decision Template

```markdown
| ID | Date | Decision | Rationale | Status |
```

When adding a new decision:

1. Use the next sequential ID.
2. Reference the bug ID (e.g. `B001`) when the decision is bug-specific.
3. Record sign-off status (`approved`, `pending`, `deferred`) for integration fixes.
4. Update `TRACKING.md` decision list in the same commit.
