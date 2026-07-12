> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit

> **Status:** ✅ Complete
> **Started:** 2026-07-08
> **Completed:** 2026-07-09
> **Scope:** Verify that every payment-provider emulator behaves faithfully to its real provider contract, document every intentional deviation, and establish a repeatable conformance framework.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/readiness-provider-conformance-audit`

---

## Why this matters

OpenMuara's value is realistic emulation. If a provider endpoint returns the wrong field, status code, or signature, users will waste time debugging integration issues that only appear in production.

This initiative treated conformance as a first-class quality gate. We mapped real provider contracts to OpenMuara implementations, closed gaps, added regression tests, and published a clear limitation registry so users know exactly what is emulated and what is not.

## Initiative structure

```
docs/initiatives/openmuara-readiness-provider-conformance-audit/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── KNOWN_ISSUES.md        # Catalog of conformance gaps and deviations
├── RISKS.md               # Risk register
├── RECOMMENDATIONS.md     # Conformance framework and prioritized actions
├── DECISIONS.md           # Decision log
├── EXECUTION_PLAN.md      # Milestones, dependencies, RACI
├── CI_INTEGRATION.md      # CI/workflow changes
├── REVIEW_CHECKLIST.md    # Sign-off checklist
├── ROLLBACK_PLAN.md       # Regression response plan
└── HANDOFF.md             # Final state and next steps
```

## Final conformance matrix

| Provider | L0 | L1 | L2 | L3 | L4 | L5 | Notes |
|---|---|---|---|---|---|---|---|
| fawry | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | v1/v2 charge, status, webhook, escape scenarios |
| stripe | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Checkout + PaymentIntent scenarios |
| senangpay | ✅ | ✅ | ✅ | ✅ | ✅ | partial | Charge + webhook covered; no refund scenario |
| billplz | ✅ | ✅ | ✅ | ✅ | ✅ | partial | Collection + bill + webhook covered |
| toyyibpay | ✅ | ✅ | ✅ | ✅ | ✅ | partial | Bill + pay page + webhook covered |
| ipay88 | ✅ | ✅ | ✅ | ✅ | ✅ | partial | Entry + backend + webhook covered |
| default | ✅ | partial | partial | n/a | n/a | partial | Reference provider; lower priority |

*L5 for regional gateways is covered by generic engine scenario tests; provider-specific state machines are limited to Fawry and Stripe.*

## Audit areas

1. ✅ **Static contract surface** — routes, methods, paths, versions.
2. ✅ **Request contract** — required fields, headers, content types, validation errors, status codes.
3. ✅ **Response contract** — JSON shapes, status codes, error payloads, idempotency behavior.
4. ✅ **Signature verification** — HMAC/SHA256 schemes, key derivation, negative/tampering tests.
5. ✅ **Webhook dispatch** — payload shape, signature headers, retries, idempotency, delivery order.
6. ✅ **State transitions** — charge → authorized → captured → refunded → failed.
7. ✅ **Simulation/escape pages** — redirect flows, callback URLs, 3DS emulation.
8. ✅ **Documentation fidelity** — every deviation documented.
9. ⬜ **External validation** — Fawry team review request prepared; sending is a follow-up action.

## Success criteria

- ✅ Conformance tests exist and pass for every P0 provider.
- ✅ Every known limitation is documented in `KNOWN_ISSUES.md` and provider docs.
- ✅ No undocumented deviation from real provider behavior remains.
- ⬜ External review from a provider team is requested and tracked (template ready; sending deferred).
- ✅ CI enforces conformance regression via golden files and contract tests.
- ✅ All quality gates pass.

## Related documents

- [`TRACKING.md`](TRACKING.md) — phases, acceptance criteria, findings log
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — closed gaps and active deviations
- [`RISKS.md`](RISKS.md) — risk register
- [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) — conformance framework and prioritized actions
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones and dependencies
- [`CI_INTEGRATION.md`](CI_INTEGRATION.md) — concrete CI changes
- [`REVIEW_CHECKLIST.md`](REVIEW_CHECKLIST.md) — sign-off checklist
- [`ROLLBACK_PLAN.md`](ROLLBACK_PLAN.md) — regression response
- [`HANDOFF.md`](HANDOFF.md) — final state
