> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Decision Log

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — accepted decisions implemented; open decisions resolved or deferred with owners.

---

## Accepted decisions

| ID | Decision | Context | Rationale | Status |
|----|----------|---------|-----------|--------|
| D01 | Use an L0–L6 conformance maturity model | Need a clear, incremental definition of "faithful emulation" | Maturity model lets us ship L1–L4 first and defer L5/L6 without losing rigor | ✅ Accepted |
| D02 | Keep conformance tests in provider packages | Tests already live near implementation (`internal/<provider>/*_test.go`) | Locality reduces maintenance burden; `internal/provider/conform` provides shared harness | ✅ Accepted |
| D03 | Use golden files for stable contract snapshots | `internal/provider/conform` already uses golden files for routes | Extend to request/response snapshots; review in PRs like code | ✅ Accepted |
| D04 | Document every deviation in provider docs + `KNOWN_ISSUES.md` | Undocumented deviations are treated as bugs | Dual documentation ensures users see limitations in docs and the audit trail is preserved | ✅ Accepted |
| D05 | Prioritize Fawry for external validation | User specifically mentioned inviting the Fawry team | Fawry is P0 and has v1/v2 complexity; a real review would surface the most gaps | ✅ Accepted |
| D06 | Pin emulated provider versions in `gateway.yml` and docs | Provider APIs are versioned and change | Version pinning lets users know exactly which contract is emulated | ✅ Accepted |

## Open decisions resolved or deferred

| ID | Question | Resolution | Owner | Date |
|----|----------|------------|-------|------|
| OD01 | How strictly should simple-runtime providers emulate real providers? | Strict for request/response shape; best-effort for quirks; deviations documented | Maintainer | 2026-07-09 |
| OD02 | Should we generate OpenAPI specs from `gateway.yml`? | Deferred to post-OSS roadmap | Maintainer | 2026-07-09 |
| OD03 | Should conformance levels be surfaced in the dashboard? | Deferred to UX enhancement backlog | Maintainer | 2026-07-09 |
| OD04 | How do we handle providers with no public sandbox or docs? | Accept limitation with documented rationale | Maintainer | 2026-07-09 |

## Decisions requiring maintainer sign-off

| ID | Question | Options | Chosen | Owner | Date |
|----|----------|---------|--------|-------|------|
| DS01 | Confirm P0 provider list | All 6 + default / Drop regional gateways | All 6 P0 providers + default as reference | Maintainer | 2026-07-09 |
| DS02 | Confirm minimum maturity to ship v1.0 | L1–L4 / L1–L5 / L1–L6 | L1–L4 for all P0; L5 for Stripe + Fawry; L6 deferred | Maintainer | 2026-07-09 |
| DS03 | Confirm external review scope for Fawry | Full v1+v2 / v1 only / Routes only | Full v1+v2 if team is willing; request template prepared | Maintainer | 2026-07-09 |

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — gaps and deviations
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — conformance framework
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones
