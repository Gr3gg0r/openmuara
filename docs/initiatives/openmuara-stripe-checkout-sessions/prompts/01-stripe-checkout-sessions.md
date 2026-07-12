> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# Prompt 01 — Stripe Checkout Sessions

## Goal
Replace the custom FPX/card charge + escape routes with a faithful Stripe Checkout Session flow for single-charge items, including a local OpenMuara-hosted checkout page and Stripe-compatible webhook configuration UI.

## Acceptance Criteria

### Cleanup

- [ ] Remove `internal/stripe/fpx.go`, `internal/stripe/fpx_types.go`, `internal/stripe/card.go`, `internal/stripe/card_types.go`, `internal/stripe/fpx_test.go`, and `internal/stripe/card_test.go`.

### Checkout Session API

- [ ] Extend `internal/stripe/types.go` / `checkout_types.go` to include:
  - `PaymentMethodTypes []string` on `CreateCheckoutSessionRequest`
  - `PaymentMethodTypes []string` on `CheckoutSession`
  - Validation: only `["fpx"]`, `["card"]`, or `["card","fpx"]` are accepted; default to `["card"]` when empty.
- [ ] Extend `POST /v1/checkout/sessions` (`checkout.go`) to:
  - Store `payment_method_types` on the session.
  - Keep the existing `url` pattern `/v1/checkout/sessions/{id}/pay`.
  - Remain single-charge only (`mode=payment`, at least one `line_item`, inline `price_data`).
  - Return Stripe-compatible error JSON on validation failure: `{"error":{"type":"invalid_request_error","code":"...","param":"...","message":"..."}}`.

### Local checkout page

- [ ] Create a checkout page handler at `GET /v1/checkout/sessions/{id}/pay`:
  - Load the session; 404 if missing.
  - Render HTML from `internal/ui/stripe-checkout.html`.
  - For `fpx`: show a bank selector + amount + confirm/cancel buttons.
  - For `card`: show a card form (number, expiry, cvc) + amount + confirm/cancel buttons.
  - For `card,fpx`: show a payment method toggle.
  - Display line item summary and total amount.
  - Include CSRF token in forms (read from `X-CSRF-Token` cookie / meta tag).
- [ ] Create a payment outcome handler at `POST /v1/checkout/sessions/{id}/pay`:
  - Validate CSRF token.
  - `action=confirm` → mark session `status=complete`, `payment_status=paid`, update ledger to `PAID`, dispatch `checkout.session.completed` webhook, redirect to `success_url`.
  - `action=cancel` → mark session `status=expired`, `payment_status=unpaid`, update ledger to `UNPAID`, dispatch `checkout.session.expired` webhook, redirect to `cancel_url`.
  - Invalid transition → 409 Conflict with Stripe-compatible JSON error (do not redirect).

### Provider wiring

- [ ] Update `internal/stripe/provider.go`:
  - Remove FPX/card route registrations and handler factory methods.
  - Register `GET /v1/checkout/sessions/{id}/pay` and `POST /v1/checkout/sessions/{id}/pay` **before** or alongside the retrieve route so the longer path wins.
  - Keep existing `/_admin/stripe/success`, `/_admin/stripe/fail`, `/_admin/stripe/cancel` simulation endpoints.
- [ ] Ensure `baseURL` fallback: if provider `baseURL` is empty, derive the local page URLs from the incoming request host/scheme.

### Webhooks

- [ ] Update `internal/stripe/webhook.go` / event payloads:
  - Ensure `checkout.session.completed` and `checkout.session.expired` payloads include `payment_method_types`.
  - Keep `Stripe-Signature` header generation.

### Webhook configuration UI

- [ ] Create `internal/ui/stripe-webhooks.html` served at `GET /_admin/stripe/webhooks`.
- [ ] Allow setting/configuring the target webhook URL and selecting enabled events (`checkout.session.completed`, `checkout.session.expired`, `payment_intent.*`).
- [ ] Display the configured `webhook_secret` (read-only) used for `Stripe-Signature`.
- [ ] Add `POST /_admin/stripe/webhooks` to save the configuration:
  - Validate CSRF.
  - Persist to `.muara/config.yml` under `providers.stripe.webhook_url` and `providers.stripe.enabled_events`.
  - Update the running dispatcher URL/event filter without requiring a restart.
  - Document the persistence choice in `DECISIONS.md`.
- [ ] Add a link to `/_admin/stripe/webhooks` on the main admin dashboard (`internal/ui/index.html`).

### OpenAPI

- [ ] Update `docs/openapi.yaml` with:
  - `POST /v1/checkout/sessions` extended request schema
  - `GET /v1/checkout/sessions/{id}/pay`
  - `POST /v1/checkout/sessions/{id}/pay`
  - `GET /_admin/stripe/webhooks`
  - `POST /_admin/stripe/webhooks`

### Tests

- [ ] Add tests in `internal/stripe/checkout_test.go` covering:
  - Create session with `payment_method_types=["fpx"]`, `["card"]`, and `["card","fpx"]`.
  - Validation rejects unsupported payment method types.
  - Retrieve session.
  - Checkout page renders for FPX and card.
  - Confirm/Cancel outcome updates session, ledger, and dispatches correct webhook.
  - CSRF failure returns 403.
  - Regression: existing checkout session tests still pass.
- [ ] Update `scripts/smoke-test.sh` to exercise the Checkout Session FPX and card happy paths using SDK-style HTTP calls.

## Files to Create/Change

- `internal/stripe/checkout.go` — extend create/retrieve, add pay handlers
- `internal/stripe/types.go` or `checkout_types.go` — add payment_method_types fields
- `internal/stripe/provider.go` — register new routes, remove FPX/card routes
- `internal/stripe/webhook.go` — ensure correct event payloads
- `internal/ui/stripe-checkout.html` — checkout page template
- `internal/ui/stripe-webhooks.html` — webhook config template
- `internal/ui/index.html` — dashboard link
- `internal/ui/embed.go` — register templates and data structs
- `internal/server/router.go` — mount `/_admin/stripe/webhooks`
- `docs/openapi.yaml` — new/extended endpoints
- `internal/stripe/checkout_test.go` — tests
- `scripts/smoke-test.sh` — smoke test update

## Response / Webhook Shape

Return:
1. Checkout Session create/retrieve JSON shapes
2. Stripe-compatible error JSON shape
3. Local checkout page HTML shape for FPX and card
4. Webhook payloads for `checkout.session.completed` and `checkout.session.expired`
5. Webhook configuration UI fields and persistence behavior

## Test Notes

- `go test ./internal/stripe/...`
- `./scripts/smoke-test.sh`
- Verify existing Stripe tests still pass.
