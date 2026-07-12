> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Stripe FPX — Post-Completion Audit Findings

**Date:** 2026-07-03  
**Auditor:** AI Agent (Kimi Code)  
**Initiative:** `docs/initiatives/openmuara-stripe-fpx/`  
**Status at audit:** ❄️ Archived / Superseded  
**Successor:** `docs/initiatives/openmuara-stripe-checkout-sessions/`

---

## Executive summary

The Stripe FPX initiative delivered a working Fawry-style charge + escape flow for both FPX and card payments. However, it was later superseded because the custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes broke Stripe SDK compatibility. This audit documents the gaps that existed at completion and the lessons learned for future provider emulation work.

**Pre-audit solidity rating:** 6/10  
**Post-audit solidity rating:** 8.5/10 (after archival, cross-linking, and recording lessons learned)

---

## What was implemented

- `POST /v1/stripe/fpx/charge`, `GET /v1/stripe/fpx/escape`, `POST /v1/stripe/fpx/escape`
- `POST /v1/stripe/card/charge`, `GET /v1/stripe/card/escape`, `POST /v1/stripe/card/escape`
- Stripe-compatible webhook signature headers
- `checkout.session.completed` and `payment_intent.canceled` webhook events
- Unit tests with `internal/stripe` coverage at 94.7%
- Smoke test coverage for the card happy path

Implementation commits: `9ed0ac4`, `69f0f6b`, `81c0cc6`.
Removal commit: `885a14d`.

---

## Gaps and recommendations

| # | Finding | Priority | Recommendation | Status |
|---|---------|----------|----------------|--------|
| R1 | Custom routes break Stripe SDK parity | High | Prefer real Stripe API paths (`/v1/checkout/sessions`, `/v1/payment_intents`) so client code works unchanged against real Stripe. | ✅ Addressed by successor |
| R2 | No OpenAPI spec updates | High | Every new public route must be added to `docs/openapi.yaml` before the prompt is closed. | ✅ Addressed by successor |
| R3 | No example app or usage docs | Medium | Provide a minimal example (`examples/stripe-fpx/`) showing request/response flow and webhook handling. | ❌ Not addressed |
| R4 | No operational runbook | Medium | Add a runbook covering common FPX test scenarios, bank selection, and webhook signature verification. | ❌ Not addressed |
| R5 | Limited edge-case coverage | Medium | Tests for idempotency, duplicate confirm/cancel, invalid bank codes, and CSRF on escape pages. | ❌ Not addressed |
| R6 | No dashboard integration | Medium | Escape pages should be inspectable/replayable from `/_admin` alongside other transactions. | ✅ Addressed by successor |
| R7 | Currency defaults undocumented | Low | Document that FPX defaults to `myr` and card defaults to `usd` in request/response examples. | ❌ Not addressed |
| R8 | Webhook event scope narrow | Low | Explicitly document which webhook events are emitted and which are deferred to generic simulation endpoints. | ✅ Partially addressed |
| R9 | Initiative not archived | High | Mark superseded initiatives as archived, link to the replacement, and record the decision. | ✅ Addressed by this audit |
| R10 | Status inconsistency | High | Keep `README.md`, `TRACKING.md`, and `HANDOFF.md` status in sync. | ✅ Addressed by this audit |

---

## Forwarded recommendations

The following recommendations remain valuable for the successor initiative and should be considered there:

- **Examples:** Add `examples/stripe-fpx/` and `examples/stripe-card/` mini-apps that demonstrate Checkout Sessions and PaymentIntents usage against OpenMuara.
- **Runbooks:** Create operational runbooks for FPX bank selection, card confirmation, webhook signature verification, and common failure modes.
- **Edge-case tests:** Cover idempotency, duplicate confirm/cancel, invalid payment method tokens, expired sessions, and CSRF protection on local payment pages.
- **Documentation:** Document currency defaults (`myr` for FPX, `usd` for card), supported FPX banks, and the full webhook event matrix.

---

## Best practices for future provider emulation

1. **Start with real provider API paths.** OpenMuara-specific extensions should be rare and clearly documented.
2. **Update OpenAPI before closing a prompt.** Public routes that are not in `docs/openapi.yaml` are effectively undocumented.
3. **Ship an example app.** A minimal working example is the fastest way for users to understand and adopt a new flow.
4. **Write a runbook.** Operational docs reduce support burden and make onboarding easier.
5. **Test edge cases early.** Idempotency, invalid transitions, and duplicate actions are common sources of bugs.
6. **Integrate with the dashboard.** Payment pages and transactions should be visible and replayable from `/_admin`.
7. **Archive superseded work promptly.** When an initiative is replaced, mark it archived, link to the successor, and record the decision.
