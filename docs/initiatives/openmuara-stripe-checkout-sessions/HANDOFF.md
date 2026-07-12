# OpenMuara Stripe FPX & Card Payments — Handoff

> Update this file at the end of every session working on this initiative.

## Last Session

- Date: 2026-07-01
- Agent: Kimi Code
- Branch: dev
- Action: Completed P02 — Stripe PaymentIntents for FPX/card with local authentication page, Stripe-compatible errors, webhook payloads, and updated smoke tests.

## Status

- P01 ✅ Completed.
- P02 ✅ Completed. All quality gates passed (build, vet, test, lint, smoke).

## Key Changes

- Added `internal/stripe/payment_intent_types.go` with `PaymentIntent`, `PaymentIntentRequest`, `PaymentIntentConfirmRequest`, `PaymentIntentNextAction`, and `PaymentIntentRedirectToURL`.
- Added `internal/stripe/payment_intent.go` implementing:
  - `POST /v1/payment_intents` (create, ledger record, dispatch `payment_intent.created`)
  - `GET /v1/payment_intents/{id}`
  - `POST /v1/payment_intents/{id}/confirm` (card → `succeeded`, FPX → `requires_action` with redirect URL, unknown token → `resource_missing`)
  - `POST /v1/payment_intents/{id}/cancel` (valid from `requires_confirmation`/`requires_action`)
- Added `internal/stripe/payment_intent_admin.go` for `GET/POST /_admin/stripe/payment_intent/{id}` with CSRF-protected FPX bank selector / card confirm page.
- Added `internal/ui/stripe-payment-intent.html`, `StripePaymentIntentPageData`, and render/serve helpers.
- Updated `internal/stripe/provider.go` to register PaymentIntent routes, build PaymentIntent webhook payloads, and route event types by reference prefix (`pi_test_*` vs `cs_test_*`).
- Updated `docs/openapi.yaml` and `internal/server/openapi.yaml` with PaymentIntent endpoints and schemas.
- Added `internal/stripe/payment_intent_test.go` covering create/retrieve/confirm/cancel, FPX/card flows, invalid token shape, CSRF failure, and webhook event type selection.
- Updated `scripts/smoke-test.sh` to exercise PaymentIntents FPX and card flows with Stripe-Signature verification.

## Blockers

None.

## Next Action

P03 is not yet defined in the initiative prompt inventory.
