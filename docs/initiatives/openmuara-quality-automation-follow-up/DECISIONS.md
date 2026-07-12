> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Quality Automation Follow-Up — Decision Log

> **Updated:** 2026-07-06

| ID | Date | Decision | Rationale | Status |
|----|------|----------|-----------|--------|
| D001 | 2026-07-06 | Phased rollout: gates start as non-blocking commentary and are promoted to required after proven stable. | Avoids flaky required checks disrupting the team; follows AGENTS.md reliability-before-polish order. | Approved |
| D002 | 2026-07-06 | Mutation testing initial threshold set at 70%; raise after baseline measurement. | 70% is a realistic starting point for a codebase that has not run mutation testing before. | Approved |
| D003 | 2026-07-06 | Provider errcode adoption is additive; existing error messages remain unchanged. | Prevents breaking clients that parse error messages; aligns with minimal-change philosophy. | Approved |
| D004 | 2026-07-06 | Visual baseline failures can be resolved by updating snapshots and committing only intentional changes. | Gives developers a clear escape hatch while keeping visual changes reviewable. | Approved |
| D005 | 2026-07-06 | Coverage gate compares changed modules against `main`/`dev` baseline, not global coverage. | Prevents blocking PRs for unrelated modules and focuses attention on actual changes. | Approved |
| D006 | 2026-07-06 | Initial mutation testing targets `internal/webhook`, `internal/engine`, and `internal/fawry`. | These packages changed most during bug hunt and represent high-value surface area. | Approved |
| D007 | 2026-07-06 | Visual baseline CI job uses path filters for `web/dashboard/**` and `internal/ui/**`. | Avoids running expensive Playwright tests on unrelated changes. | Approved |
| D008 | 2026-07-06 | Coverage gate ignores packages with fewer than 10 changed lines. | Prevents false positives on trivial changes. | Approved |
| D009 | 2026-07-06 | Coverage regression gate runs with `continue-on-error` until three stable PRs are observed. | Implements D001 for P03 while the parser and thresholds are proven in production. | Approved |
| D010 | 2026-07-06 | Mutation testing excludes `internal/fawry` from the initial curated list. | Gremlins times out on fawry HTTP handler/signature tests because mutations cause test servers to hang; re-evaluate after faster pure-function tests are added. | Approved |
| D011 | 2026-07-06 | Mutation testing CI job runs with `continue-on-error` during the phased rollout. | Tool behavior can be environment-sensitive (test-server timeouts under mutation); keep the job reporting scores without blocking merges until it is stable across local and CI. | Approved |
| D012 | 2026-07-06 | Visual baseline captures separate light and dark theme snapshots. | Catches theme-specific regressions and removes OS-preference dependency from CI (R19). | Approved |
| D013 | 2026-07-06 | Dynamic dashboard elements are hidden in visual tests via a shared `[data-visual-mask]` CSS rule. | Easier to stabilize future visual tests without per-test CSS injection (R21). | Approved |
| D014 | 2026-07-06 | Visual-baseline and mutation jobs live in separate workflows with path filters. | Avoids unnecessary CI minutes on unrelated changes; visual baseline filters to dashboard/UI paths, mutation filters to Go paths (R44, R45). | Approved |

---

## Decision Template

```markdown
| ID | Date | Decision | Rationale | Status |
```

When adding a new decision:

1. Use the next sequential ID.
2. Reference the prompt ID (e.g. `P01`) when the decision is prompt-specific.
3. Record sign-off status (`approved`, `pending`, `deferred`) for any change that affects public error messages or CI contracts.
