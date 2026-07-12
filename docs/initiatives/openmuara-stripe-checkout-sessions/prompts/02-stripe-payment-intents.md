> **⚠️ AI AGENT: Read `AGENTS.md`, the initiative `README.md`, and complete P01 first.**

# Prompt 02 — Stripe PaymentIntents API

## Goal
Add a faithful Stripe PaymentIntents API emulation for single-charge FPX and card payments, alongside the Checkout Session flow completed in P01.

## Acceptance Criteria

### PaymentIntent API

- [ ] Create `internal/stripe/payment_intent_types.go` with `PaymentIntent`, `PaymentIntentRequest`, `PaymentIntentNextAction`, and related types matching Stripe's documented shape for the implemented subset.
- [ ] Create `internal/stripe/payment_intent.go` implementing:
  - `POST /v1/payment_intents` — create a PaymentIntent
  - `GET /v1/payment_intents/{id}` — retrieve a PaymentIntent
  - `POST /v1/payment_intents/{id}/confirm` — confirm and produce `next_action`
  - `POST /v1/payment_intents/{id}/cancel` — cancel
- [ ] Support `payment_method_types: ["fpx"]`, `["card"]`, or `["card","fpx"]`.
- [ ] Return Stripe-compatible error JSON on validation failure: `{"error":{"type":"invalid_request_error","code":"...","param":"...","message":"..."}}`.
- [ ] Store PaymentIntents in a new in-memory store keyed by ID (`pi_test_*`).
- [ ] Record transactions in the shared ledger with reference equal to PaymentIntent ID.

### Confirm behavior

- [ ] `confirm` for FPX returns `status: "requires_action"` and `next_action.redirect_to_url.url` pointing to `GET /_admin/stripe/payment_intent/{id}`.
- [ ] `confirm` for card returns `status: "succeeded"` and dispatches `payment_intent.succeeded` webhook.
- [ ] Accept test payment method tokens:
  - `pm_card_visa`, `pm_card_mastercard` → success
  - `pm_fpx_maybank`, `pm_fpx_cimb`, `pm_fpx_publicbank` → success (used with FPX)
  - Unknown tokens → `invalid_request_error` with `code: resource_missing`

### Local authentication page

- [ ] `GET /_admin/stripe/payment_intent/{id}` renders a local page:
  - For FPX: bank selector with Confirm/Cancel
  - For card: Confirm/Cancel only
  - Include CSRF token in forms.
- [ ] `POST /_admin/stripe/payment_intent/{id}` processes the outcome:
  - Validate CSRF.
  - Confirm → `status: "succeeded"`, dispatch `payment_intent.succeeded`
  - Cancel → `status: "canceled"`, dispatch `payment_intent.canceled`
  - Invalid transition → `409 Conflict` with Stripe-compatible JSON error

### Cancel behavior

- [ ] `POST /v1/payment_intents/{id}/cancel` dispatches `payment_intent.canceled`.
- [ ] Cancel is only valid from `requires_confirmation` or `requires_action`.

### Webhooks

- [ ] All webhooks include a valid `Stripe-Signature` header.
- [ ] Dispatch `payment_intent.created` on create.
- [ ] Dispatch `payment_intent.succeeded` on confirm success or FPX confirm.
- [ ] Dispatch `payment_intent.canceled` on cancel or authentication page cancel.
- [ ] Respect the enabled-events filter configured in `/_admin/stripe/webhooks`.

### Provider wiring

- [ ] Update `internal/stripe/provider.go` to register PaymentIntent and admin routes.
- [ ] Ensure `baseURL` fallback: if provider `baseURL` is empty, derive local page URLs from the incoming request host/scheme.
- [ ] Existing Stripe Checkout Session behavior from P01 remains unchanged.

### OpenAPI

- [ ] Update `docs/openapi.yaml` with:
  - `POST /v1/payment_intents`
  - `GET /v1/payment_intents/{id}`
  - `POST /v1/payment_intents/{id}/confirm`
  - `POST /v1/payment_intents/{id}/cancel`
  - `GET /_admin/stripe/payment_intent/{id}`
  - `POST /_admin/stripe/payment_intent/{id}`

### Tests

- [ ] Add tests in `internal/stripe/payment_intent_test.go` covering:
  - Create, retrieve, confirm, cancel.
  - FPX confirm produces `requires_action` with redirect URL.
  - Card confirm succeeds immediately.
  - Invalid payment method token error shape.
  - CSRF failure on admin page.
  - Webhook event filtering.
- [ ] Update `scripts/smoke-test.sh` to exercise the PaymentIntents flow using Stripe SDK-style HTTP calls.

## Files to Create/Change

- `internal/stripe/payment_intent_types.go` — PaymentIntent types
- `internal/stripe/payment_intent.go` — PaymentIntent handlers and store
- `internal/stripe/provider.go` — register PaymentIntent routes
- `internal/stripe/webhook.go` — add PaymentIntent event payloads
- `internal/stripe/payment_intent_test.go` — tests
- `docs/openapi.yaml` — new endpoints
- `scripts/smoke-test.sh` — update smoke test

## Response / Webhook Shape

Return:
1. PaymentIntent create/retrieve/confirm/cancel JSON shapes
2. Stripe-compatible error JSON shape
3. FPX and card authentication page HTML shapes
4. Webhook payloads for `payment_intent.created`, `payment_intent.succeeded`, `payment_intent.canceled`

## Test Notes

- `go test ./internal/stripe/...`
- `./scripts/smoke-test.sh`
- Verify existing Stripe tests still pass.
