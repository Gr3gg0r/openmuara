> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# Prompt 01 — Stripe FPX and Card Charge and Escape

> **Status:** ❄️ Superseded. The custom routes implemented by this prompt were removed in commit `885a14d` and replaced with Stripe Checkout Sessions and PaymentIntents. See [`docs/initiatives/openmuara-stripe-checkout-sessions/`](../../openmuara-stripe-checkout-sessions/).

## Goal
Add dedicated Stripe FPX and card charge + escape flows, modeled on the Fawry provider pattern.

## Acceptance Criteria

### FPX

- [x] `POST /v1/stripe/fpx/charge` accepts an FPX charge request, validates it, records a transaction in the ledger, and returns a Stripe-style charge response with a `reference` and `escape_url`.
- [x] `GET /v1/stripe/fpx/escape` renders a minimal HTML bank selector page.
- [x] `POST /v1/stripe/fpx/escape` confirms or cancels the FPX payment, dispatches the correct Stripe webhook, and redirects to the caller.

### Card

- [x] `POST /v1/stripe/card/charge` accepts a card charge request, validates it, records a transaction in the ledger, and returns a Stripe-style charge response with a `reference` and `escape_url`.
- [x] `GET /v1/stripe/card/escape` renders a minimal HTML card confirmation page with Confirm and Cancel buttons.
- [x] `POST /v1/stripe/card/escape` confirms or cancels the card payment, dispatches the correct Stripe webhook, and redirects to the caller.

### Shared

- [x] Confirm transitions the transaction to `paid` and dispatches `checkout.session.completed`.
- [x] Cancel transitions the transaction to `unpaid` and dispatches `payment_intent.canceled`.
- [x] Invalid transitions return `409 Conflict` using `engine.Transition`.
- [x] Existing Stripe Checkout session flows remain unchanged.
- [x] All new code is covered by unit tests.

## Files to Create/Change

- `internal/stripe/fpx_types.go` — FPX request/response types
- `internal/stripe/fpx.go` — FPX charge, escape page, escape action
- `internal/stripe/card_types.go` — card request/response types
- `internal/stripe/card.go` — card charge, escape page, escape action
- `internal/stripe/provider.go` — register new routes
- `internal/stripe/fpx_test.go` — FPX tests
- `internal/stripe/card_test.go` — card tests
- `scripts/smoke-test.sh` — add card flow smoke step

## Response / Webhook Shape

Return:
1. FPX charge request/response JSON shape
2. Card charge request/response JSON shape
3. Webhook event payload for `checkout.session.completed`
4. Webhook event payload for `payment_intent.canceled`

## Test Notes

- `go test ./internal/stripe/...`
- `./scripts/smoke-test.sh`
- Verify existing Stripe tests still pass.
