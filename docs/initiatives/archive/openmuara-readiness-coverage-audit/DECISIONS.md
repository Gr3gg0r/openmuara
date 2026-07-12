> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Coverage Audit Decision Log

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete

---

## Accepted decisions

| ID | Decision | Context | Rationale | Status |
|----|----------|---------|-----------|--------|
| D01 | Go overall coverage target is **81%** | Baseline was 80.8%; existing CI gate was 80% | A 1% safety margin prevents cliff-edge CI failures from trivial refactors | ✅ Accepted |
| D02 | Dashboard coverage provider is **@vitest/coverage-v8@^2.1.9** | Vitest 2.x is already in use | Same major version as Vitest; v8 is the default and well-maintained provider | ✅ Accepted |
| D03 | Initial dashboard thresholds: statements/lines **60%**, branches/functions **55%** | Dashboard baseline was 49.6% overall, with many uncovered views | Achievable in one execution pass; avoids blocking legitimate feature work while raising the bar | ✅ Accepted |
| D04 | Target dashboard thresholds after phase 2: statements/lines **70%**, branches/functions **65%** | Once initial gaps are closed, continue raising quality | Industry good-practice target for a Preact SPA; reviewed after initial phase | ✅ Accepted |
| D05 | Enforce **per-package Go floors** in addition to overall coverage | Overall coverage can hide weak packages | Ensures no critical package is undertested without explicit sign-off | ✅ Accepted |
| D06 | Exempt `examples/checkout-store` from Go coverage enforcement | Example application, not part of distributed library | Examples demonstrate usage; their quality is validated by build/run, not unit tests | ✅ Accepted |
| D07 | Exempt `internal/provider/factory` from floor enforcement (smoke test only) | Thin generated/registry package with no current tests | A smoke test confirming registration is sufficient; full coverage not valuable | ✅ Accepted |
| D08 | Treat `internal/ui` as an embedding layer with a **70% floor** | Embeds generated dashboard assets and serves static files | Most code is file-server wiring; 70% covers error paths and asset serving | ✅ Accepted |
| D09 | Exclude `web/dashboard` build/config files and entry points from dashboard threshold | `app.tsx`, `main.tsx`, `*.config.ts`, `scripts/**`, `e2e/**` | These are wiring/build harness, not user-facing logic; still covered in report but not gated | ✅ Accepted |
| D10 | Coverage regression gate becomes **blocking after three stable PRs** | Currently `continue-on-error: true` in `coverage-comment.yml` | Allows baseline to stabilize while giving teams time to adapt | ✅ Accepted |
| D11 | Coverage artifacts retained for **30 days** | GitHub Actions default artifact retention | Sufficient for PR review and audit; longer retention can be set at org level | ✅ Accepted |
| D12 | Calibrate `internal/provider/simple` floor to **45%** | Reference/demo YAML-driven provider with large handler surface | Focus coverage on public API and runtime wiring; full handler matrix is low-value for a demo provider | ✅ Accepted |
| D13 | Calibrate `internal/provider/conform` floor to **79%** | Uncovered lines are `t.Fatalf` branches in golden-file update path | Error branches in test scaffolding are not worth covering | ✅ Accepted |
| D14 | Calibrate `internal/version` floor to **70%** | Uncovered lines are `init()` git/exec fallbacks | Build environment variability makes these paths impossible to exercise reliably | ✅ Accepted |

## Open decisions — resolved during planning/execution

| ID | Question | Resolution | Rationale |
|----|----------|------------|-----------|
| OD01 | Which dashboard coverage provider? | `@vitest/coverage-v8@^2.1.9` | Matches Vitest major, default provider, no extra config |
| OD02 | Should `app.tsx`/`main.tsx` be included in dashboard threshold? | Excluded from threshold, included in report | Entry points are mostly wiring; gate would be noisy |
| OD03 | Should coverage regression be blocking immediately? | Non-blocking for three stable PRs, then blocking | Prevents churn while baseline stabilizes |
| OD04 | What is the right dashboard starting threshold? | 60% statements / 55% branches/functions | Calibrated against baseline; reachable without blocking delivery |
| OD05 | How to handle packages that cannot reach 80%? | Documented exemption with rationale and review date | Maintains quality bar without forcing low-value tests |

## Decisions ratified by maintainer

| ID | Question | Options | Ratified choice | Owner | Status |
|----|----------|---------|-----------------|-------|--------|
| DS01 | Confirm dashboard threshold start values | 60/55/55/55 or 50/50/50/50 | **60/55/55/55** | Maintainer | ✅ Ratified |
| DS02 | Confirm `internal/provider/simple` treatment | Add tests to 70% or formally exempt at 45% | **45% floor** with focused public-API tests | Maintainer | ✅ Ratified |
| DS03 | Confirm regression gate PR count | 3, 5, or after next release | **3 stable PRs** | Maintainer | ✅ Ratified |

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — gaps and exemptions
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones
- [`HANDOFF.md`](HANDOFF.md) — final state
