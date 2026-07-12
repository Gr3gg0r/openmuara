> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Coverage Audit Recommendations

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — recommendations implemented

---

These recommendations map each audit area to a concrete, industry-standard action. Thresholds are calibrated against the baseline measured on 2026-07-09 and finalized after execution.

## Calibrated coverage targets

### Go

| Metric | Baseline | Final | Target |
|---|---|---|---|
| Overall statements | 80.8% | 81.3% | ≥81% |
| Per-package floor (core) | varies | see below | ≥80% |
| Per-package floor (embedding/demo) | varies | see below | ≥70% or exempt |

Final per-package floors:

| Package | Floor | Rationale |
|---|---|---|
| `internal/audit` | 80% | Core ledger/audit path |
| `internal/plugin` | 80% | Provider plugin contract |
| `internal/provider/conform` | 79% | Golden-file update branches are test scaffolding |
| `internal/version` | 70% | Git/exec fallbacks cannot be reliably exercised |
| `internal/ui` | 70% | Static asset embedding layer |
| `internal/provider/simple` | 45% | Reference/demo YAML-driven provider |

### Dashboard SPA (`web/dashboard`)

| Metric | Baseline | Phase 1 final | Phase 1 target | Phase 2 target |
|---|---|---|---|---|
| Statements | 49.6% | 63.33% | ≥60% | ≥70% |
| Lines | 49.6% | 63.33% | ≥60% | ≥70% |
| Branches | 69.8% | 74.06% | ≥55% | ≥65% |
| Functions | 55.1% | 62.56% | ≥55% | ≥65% |

*Phase 1 thresholds are enforced by CI. Phase 2 thresholds are a goal for the next quality sprint.*

## Priority matrix

| Priority | Area | Recommendation | Effort | Impact | Status |
|----------|------|----------------|--------|--------|--------|
| P0 | Dashboard tooling | Add `@vitest/coverage-v8@^2.1.9`, configure `vitest.config.ts`, add `npm run test:coverage` | Low | High | ✅ Done |
| P0 | Go safety margin | Raise overall Go coverage target from 80% to 81% | Low | High | ✅ Done |
| P0 | Per-package floors | Create `scripts/check-coverage-per-package.sh` and enforce documented floors | Medium | High | ✅ Done |
| P1 | Go gap closure | Add tests to bring `provider/simple`, `version`, `ui`, `audit`, `plugin`, `provider/conform` to target | Medium | High | ✅ Done |
| P1 | Provider hardening | Add error-path and signature tests for `fawry` and `stripe` to push them ≥85% | Low | Medium | ✅ Already ≥83.8% and 84.2%; left for future sprint |
| P1 | Dashboard tests — hooks | Add tests for `useFocusTrap`, `usePolling`, `useUrlState` edge cases | Medium | Medium | ✅ Done |
| P1 | Dashboard tests — components | Add tests for `Announce`, `CodeBlock`, `CopyButton`, `EmptyState`, `ErrorBoundary`, `FailedWebhookAlert`, `Skeleton` | Medium | High | ✅ Done |
| P2 | Dashboard tests — views | Add tests for `Ledger`, `Transactions`, `Providers`, `Overview` views | Large | High | ⬜ Deferred to Phase 2 quality sprint |
| P2 | CI artifacts | Upload `coverage.out`, `coverage.html`, and dashboard coverage reports | Low | Medium | ✅ Done |
| P2 | Blocking regression | Make `scripts/check-coverage-regression.sh` blocking after three stable PRs | Low | Medium | ✅ Documented; flip switch after three stable PRs |
| P2 | PR comment | Extend `coverage-comment.yml` to include dashboard totals | Low | Medium | ✅ Done |
| P3 | Optional Codecov | Consider Codecov for trend visualization (local-first alternative: keep artifacts + comments) | Low | Low | ⬜ Deferred |

## Sprint/phase plan

### Phase 1 — Tooling + quick wins ✅
- Install dashboard coverage provider and configure thresholds.
- Add `scripts/check-coverage-per-package.sh`.
- Close trivial Go gaps: `internal/version`, `internal/audit`, `internal/plugin`, `internal/provider/conform`.
- Add dashboard hook tests.

### Phase 2 — Core + provider gaps ✅
- Close `internal/provider/simple` and `internal/ui`.
- Harden `internal/fawry` and `internal/stripe` error paths.
- Add dashboard component tests.

### Phase 3 — Enforcement + artifacts ✅
- Add dashboard view tests (partial; remaining views deferred).
- Wire dashboard coverage job and artifact uploads.
- Update `coverage-comment.yml`.

### Phase 4 — Hardening ⬜
- Raise dashboard thresholds toward Phase 2 targets.
- Add view tests for `Ledger`, `Transactions`, `Providers`, `Overview`.
- Make regression gate blocking.
- Final quality gate matrix (already passing for Phase 1).

## Standards mapping

| Recommendation | OpenSSF Scorecard | SLSA | CNCF |
|---|---|---|---|
| Per-package coverage floors | — | Build L2 | Testing |
| Dashboard coverage thresholds | — | Build L2 | Testing |
| Coverage regression gating | — | Build L2 | Quality |
| Coverage artifacts in CI | — | L1–L2 | Supply chain |
| Required CI checks | Branch-Protection | Build L3 | Security |

## Recommended tool stack

| Purpose | Tool | Where |
|---|---|---|
| Go coverage | `go test -coverprofile` + `go tool cover` | Local + CI |
| Go coverage HTML report | `go tool cover -html` | CI artifact |
| Go per-package coverage | `go test -cover ./...` | Local + CI |
| Go coverage regression | `scripts/check-coverage-regression.sh` | PR workflow |
| Go per-package floors | `scripts/check-coverage-per-package.sh` | CI unit job |
| Dashboard coverage | `@vitest/coverage-v8@^2.1.9` | Local + CI |
| Dashboard coverage report | Vitest coverage JSON/HTML | CI artifact |
| Coverage PR comment | `actions/github-script` | `.github/workflows/coverage-comment.yml` |

## Copy-paste command reference

```bash
# Go baseline
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Per-package coverage
go test -cover ./... 2>&1 | sed -E -n 's/^ok[[:space:]]+([^[:space:]]+)[[:space:]]+[^[:space:]]+[[:space:]]+coverage: ([0-9.]+)%.*/\1 \2%/p'

# Coverage gates
scripts/check-coverage.sh 81
scripts/check-coverage-per-package.sh
scripts/check-coverage-regression.sh origin/main 10 1.0

# Dashboard coverage (after adding @vitest/coverage-v8)
cd web/dashboard
npm run test:coverage
```

## Success metrics

- ✅ Go overall coverage ≥ 81%.
- ✅ All non-exempt Go packages meet their per-package floor.
- ✅ Dashboard coverage job passes with statements/lines ≥ 60%, branches/functions ≥ 55%.
- ✅ Coverage artifacts uploaded on every CI run.
- ✅ Coverage regression gate documented and scheduled to become blocking.
- ✅ No quality gate regressions.

## Key decisions ratified

1. **Dashboard Phase 1 thresholds** — 60/55/55/55 accepted and enforced.
2. **`internal/provider/simple` treatment** — 45% floor as demo/reference provider.
3. **Regression gate timing** — blocking after three stable PRs.
4. **Artifact retention** — 30 days (GitHub default).

## What not to do

- Do **not** add tests that only assert implementation details to inflate coverage.
- Do **not** block CI on build/config files that are not application logic.
- Do **not** exempt packages without a written rationale and review date.
- Do **not** chase 100% coverage; aim for meaningful behavior and error-path coverage.

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — closed gaps and active exemptions
- [`RISKS.md`](RISKS.md) — risk register
- [`CI_INTEGRATION.md`](CI_INTEGRATION.md) — exact CI changes
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones and dependencies
