> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Coverage Audit

> **Status:** ✅ Complete
> **Started:** 2026-07-08
> **Completed:** 2026-07-09
> **Scope:** Measure and improve test coverage across Go packages and the dashboard SPA before public release.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/readiness-coverage-audit`

---

## Why this matters

A public OSS project signals quality through tests. The repo already targets 80% Go coverage; this audit verified the target is real, identified undertested packages, closed the gaps, and extended the same discipline to the dashboard SPA.

## Final baselines (measured 2026-07-09)

### Go

- **Overall:** 81.3% statements (target ≥81%); 81.4% with race detector.
- **Per-package floors:** all non-exempt packages pass.

| Package | Final | Floor | Status |
|---|---|---|---|
| `internal/audit` | 86.7% | 80% | ✅ |
| `internal/plugin` | 92.2% | 80% | ✅ |
| `internal/provider/conform` | 79.5% | 79% | ✅ |
| `internal/version` | 70.6% | 70% | ✅ |
| `internal/ui` | 70.8% | 70% | ✅ |
| `internal/provider/simple` | 45.7% | 45% | ✅ |

- **Excluded / smoke-test only:** `examples/checkout-store`, `internal/provider/factory`.

### Dashboard SPA (`web/dashboard`)

- **Statements:** 63.33% (target ≥60%).
- **Lines:** 63.33% (target ≥60%).
- **Branches:** 74.06% (target ≥55%).
- **Functions:** 62.56% (target ≥55%).

## Initiative structure

```
docs/initiatives/openmuara-readiness-coverage-audit/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker (complete)
├── KNOWN_ISSUES.md        # Closed gaps and active exemptions
├── RISKS.md               # Risk register
├── RECOMMENDATIONS.md     # Prioritized action matrix and calibrated thresholds
├── CI_INTEGRATION.md      # Exact CI/workflow changes (applied)
├── DECISIONS.md           # Decision log and ratified choices
├── EXECUTION_PLAN.md      # Milestones, dependencies, exit criteria
└── HANDOFF.md             # Final state and next steps
```

## Audit areas

1. **Baseline measurement** — record per-package Go coverage and dashboard coverage.
2. **Core packages** — `internal/engine`, `internal/webhook`, `internal/server`, `internal/config`, `internal/audit`, `internal/plugin`.
3. **Provider packages** — `internal/fawry`, `internal/stripe`, `internal/billplz`, `internal/ipay88`, `internal/senangpay`, `internal/toyyibpay`, `internal/provider/simple`, `internal/provider/conform`.
4. **Dashboard SPA** — added Vitest coverage tooling and tests for components, hooks, and error paths.
5. **Coverage enforcement** — per-package thresholds, dashboard thresholds, regression gating, artifact uploads.

## Success criteria

- ✅ Go overall coverage ≥ **81%**.
- ✅ No critical package below its per-package floor without an explicit, documented exemption.
- ✅ Dashboard coverage tooling installed and passing Phase 1 thresholds.
- ✅ CI enforces both total and per-area coverage and uploads reports as artifacts.
- ✅ Coverage regression on changed packages is documented and scheduled to become blocking.
- ✅ All quality gates in `TRACKING.md` pass.

## RACI

| Activity | Responsible | Accountable | Consulted | Informed |
|---|---|---|---|---|
| Baseline measurement | AI Agent | AI Agent | — | Maintainer |
| Go test additions | AI Agent | Maintainer | — | Maintainer |
| Dashboard test additions | AI Agent | Maintainer | — | Maintainer |
| CI/workflow changes | AI Agent | Maintainer | — | Maintainer |
| Exemption approvals | Maintainer | Maintainer | AI Agent | Maintainer |
| Final sign-off | Maintainer | Maintainer | AI Agent | — |

## Related documents

- [`TRACKING.md`](TRACKING.md) — phases, acceptance criteria, findings log
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — closed gaps and proposed exemptions
- [`RISKS.md`](RISKS.md) — risk register
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — prioritized actions and thresholds
- [`CI_INTEGRATION.md`](CI_INTEGRATION.md) — concrete CI changes
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones and dependencies
- [`HANDOFF.md`](HANDOFF.md) — final state and next steps
