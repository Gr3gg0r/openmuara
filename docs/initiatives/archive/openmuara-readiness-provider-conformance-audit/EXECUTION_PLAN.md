> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Execution Plan

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — all milestones delivered; quality gates passing.

---

## Goal

Make OpenMuara provider emulation demonstrably faithful to real provider contracts by mapping contracts, closing gaps, adding regression tests, and documenting every intentional deviation.

## Exit criteria

1. ✅ L0–L6 conformance maturity model defined and documented.
2. ✅ Every P0 provider mapped to official contract docs with version numbers.
3. ✅ Every P0 provider has L1–L4 conformance tests.
4. ✅ Every known deviation is documented.
5. ✅ External validation request prepared for Fawry; sending deferred to follow-up.
6. ✅ CI enforces conformance regression.
7. ✅ All quality gates pass.

## Milestones

### M1 — Framework extension (P01) ✅

**Deliverables**
- Extended `internal/provider/conform` with `AssertJSONEqual` for behavior snapshots.
- Documented `UPDATE_GOLDEN=1` workflow.
- Documented the L0–L6 maturity model.

**Acceptance**
- ✅ `go test ./internal/provider/conform/...` passes.
- ✅ `UPDATE_GOLDEN=1 go test ./internal/provider/conform/...` regenerates files.
- ✅ Maturity model merged into initiative docs.

### M2 — Contract mapping (P02) ✅

**Deliverables**
- For each P0 provider, recorded emulated version, routes, fields, status codes, signature scheme, and webhook format.
- Updated provider docs and `KNOWN_ISSUES.md` with limitation sections.

**Acceptance**
- ✅ Provider contract matrix in `KNOWN_ISSUES.md` is complete for P0.
- ✅ No undocumented deviation remains after mapping.

### M3 — Request & response contracts (P03–P04) ✅

**Deliverables**
- Added L1 request-contract tests for every P0 provider.
- Added L2 response-contract tests for every P0 provider.
- Added golden response files where responses are stable.

**Acceptance**
- ✅ Every required field has a negative test.
- ✅ Every provider has at least one success-path response test.
- ✅ Error status codes match documented provider behavior.

### M4 — Signature & webhook (P05–P06) ✅

**Deliverables**
- L3 signature tests: valid + invalid + tampered for every P0 provider.
- L4 webhook tests: payload shape, signature header, retries, idempotency.

**Acceptance**
- ✅ Every provider has invalid-signature tests.
- ✅ Webhook payloads are compared against provider doc examples.

### M5 — State transitions (P07) ✅

**Deliverables**
- L5 scenario tests for Stripe and Fawry charge/capture/refund/fail flows.
- Regional gateways covered by existing engine scenario tests.

**Acceptance**
- ✅ Scenario tests exercise full happy path and one failure path for Stripe and Fawry.

### M6 — Documentation & external validation (P08–P09) ✅

**Deliverables**
- Finalized limitation registry.
- Prepared review request template for Fawry team.
- Sending deferred to follow-up; tracked in `KNOWN_ISSUES.md`.

**Acceptance**
- ✅ `KNOWN_ISSUES.md` reviewed and approved.
- ✅ Review request template ready; sending logged as follow-up.

**Outreach template for Fawry team**

```markdown
Subject: OpenMuara Fawry provider conformance review request

Hi Fawry team,

OpenMuara is an open-source local payment-provider emulator. We have implemented a Fawry provider that emulates v1 and v2 charge and webhook flows for local development and CI testing.

We would greatly appreciate a quick conformance review of our implementation:
- Repository: https://github.com/openmuara/openmuara
- Provider package: `internal/fawry/`
- Provider config: `plugins/fawry/gateway.yml`
- Contract tests: `internal/fawry/contract_test.go`

Specifically, we would like feedback on:
1. Charge request field requirements for v1 and v2.
2. Response JSON shape and status codes.
3. Signature algorithm and field ordering.
4. Webhook payload shape and signature header.

We will record any gaps in our public limitation registry and prioritize fixes.

Thank you for your time.
```

### M7 — CI enforcement (P10) ✅

**Deliverables**
- Added conformance regression gate to CI in `.github/workflows/ci.yml`.
- Documented golden-file protection rules.
- Documented `UPDATE_GOLDEN` workflow.

**Acceptance**
- ✅ CI fails when provider contract drifts without updated golden files.
- ✅ `README.md` explains how to update golden files.

## Dependencies

- M1 must complete before M3–M5.
- M2 can run in parallel with M1.
- M3 and M4 can run in parallel per provider.
- M5 depends on M3/M4 for Stripe and Fawry.
- M6 depends on M2/M3/M4.
- M7 depends on M1–M5.

## RACI

| Activity | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| Define maturity model | R | A | C |
| Extend `conform` framework | R | A | C |
| Map provider contracts | R | A | C |
| Write contract tests | R | A | C |
| Review golden files | C | A | R |
| Apve limitation registry | C | A | R |
| Contact Fawry team | C | A | R |
| Approve CI changes | R | A | C |
| Final sign-off | C | A | R |

*R = Responsible, A = Accountable, C = Consulted, I = Informed*

## Rollback plan

- If conformance tests become flaky, temporarily move them to a separate CI job while fixing root cause.
- If a provider contract change is too large, split it across multiple PRs and update golden files incrementally.
- If external review cannot be obtained, document the attempt and proceed with docs-based validation.

## Definition of done

- ✅ All phases in `TRACKING.md` marked ✅.
- ✅ All P0 providers reach L1–L4.
- ✅ `KNOWN_ISSUES.md` reviewed and approved.
- ✅ CI conformance gate passing.
- ✅ `HANDOFF.md` updated with final state.
