> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P02 — Triage and Prioritization

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** P01
> **Target files:** `TRACKING.md`, `RISKS.md`, `DECISIONS.md`
> **Status:** ⬜

## Goal

Classify every finding from P01, group related bugs into fix batches, and explicitly defer anything too risky or out of scope.

## Tasks

- [ ] Validate each P01 finding is reproducible; remove or downgrade false positives with rationale.
- [ ] Assign severity using the rubric:
  - **P0** — crash, security vulnerability, data loss, or completely broken primary flow.
  - **P1** — broken feature, UX regression, or incorrect provider behavior blocking a common use case.
  - **P2** — polish, edge case, or cosmetic issue.
- [ ] Mark any bug that breaks a dashboard redesign invariant (left nav, ledger default, filters, detail pages, provider settings, dual-port) as at least P1.
- [ ] Identify P0/P1 bugs that touch provider emulation logic, webhook signature verification, config persistence, auth/billing/PII, or the provider plugin schema contract; these require explicit user sign-off before fixing.
- [ ] Group bugs by area (provider, webhook, config, UI, CLI, build/test, docs) and by fix risk.
- [ ] Identify 2–3 quick wins (safe, high-impact) for P03.
- [ ] Identify risky or cross-cutting bugs for P04 or deferred handling.
- [ ] Update `TRACKING.md` bug register with severity, assigned batch, root cause category, and sign-off status.
- [ ] Update `RISKS.md` with any new systemic risks and deferred items.
- [ ] Update `DECISIONS.md` with triage outcomes and any sign-off requests.
- [ ] Update `HANDOFF.md` with triage rationale.

## Acceptance Criteria

- [ ] Every P01 finding is either confirmed with severity or marked false positive with rationale.
- [ ] All P0/P1 integration fixes have a sign-off decision in `DECISIONS.md` (approved, pending, or deferred).
- [ ] P03 batch contains 2–3 safe, high-impact fixes.
- [ ] P04 batch contains remaining fixes or a plan for deferred items.
- [ ] `HANDOFF.md` updated with triage rationale and open questions.

## Quality Gates

No product code changes in this prompt.

## Notes

- Do **not** fix bugs in this prompt.
- If a bug is actually a missing feature, mark it as deferred and link to the v1.1/v2 backlog.
- If deferring a P0 or P1 bug, get user sign-off and record the rationale.
