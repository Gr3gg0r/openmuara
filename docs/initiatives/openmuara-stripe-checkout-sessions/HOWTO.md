> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# OpenMuara Stripe FPX & Card Payments — HOWTO

## Decomposition

This initiative has two prompts:

- **P01 — Stripe Checkout Sessions**
- **P02 — Stripe PaymentIntents**

Execute P01 first. P02 can start after P01 is committed and passing all gates.

### P01 breakdown

1. **Extend Checkout Session types**
   - Add `PaymentMethodTypes` to request and session types.
   - Validate allowed values and default to `["card"]`.

2. **Remove custom FPX/card routes**
   - Delete `fpx.go`, `fpx_types.go`, `card.go`, `card_types.go`, and their tests.
   - Remove handler factory methods from `provider.go`.

3. **Implement Stripe-compatible errors**
   - Add helper to build `{"error":{"type","code","param","message"}}` responses.
   - Use it for validation and invalid-transition errors.

4. **Implement local checkout page**
   - Create `internal/ui/stripe-checkout.html`.
   - Add `GET /v1/checkout/sessions/{id}/pay` handler that renders FPX or card UI based on `payment_method_types`.
   - Add `POST /v1/checkout/sessions/{id}/pay` handler to process confirm/cancel.
   - Include CSRF tokens in forms.

5. **Update session lifecycle**
   - Confirm → `status=complete`, `payment_status=paid`, ledger `PAID`, dispatch `checkout.session.completed`, redirect to `success_url`.
   - Cancel → `status=expired`, `payment_status=unpaid`, ledger `UNPAID`, dispatch `checkout.session.expired`, redirect to `cancel_url`.

6. **Update provider wiring**
   - Register checkout page routes before/without shadowing retrieve route.
   - Keep existing admin simulation endpoints.
   - Add baseURL fallback from request host/scheme.

7. **Update webhooks**
   - Ensure event payloads include `payment_method_types`.
   - Sign with existing `SignPayload`.

8. **Add webhook configuration UI**
   - Create `internal/ui/stripe-webhooks.html`.
   - Serve at `/_admin/stripe/webhooks`.
   - Allow URL and event selection; display signing secret.
   - Persist to `.muara/config.yml` and update running dispatcher.
   - Add dashboard link.

9. **Update OpenAPI spec**
   - Add new/extended endpoints to `docs/openapi.yaml`.

10. **Test**
    - Unit tests for create/retrieve/pay/outcome for FPX and card.
    - CSRF failure test.
    - Smoke test using SDK-style HTTP calls.
    - Regression test for existing checkout behavior.

### P02 breakdown

1. **Design PaymentIntent types**
   - Create `internal/stripe/payment_intent_types.go`.
   - Mirror Stripe's PaymentIntent object shape for the implemented subset.

2. **Implement PaymentIntent store**
   - Create an in-memory store keyed by PaymentIntent ID (`pi_test_*`).
   - Define create, get, update operations.

3. **Implement create handler**
   - `POST /v1/payment_intents`.
   - Validate `amount`, `currency`, `payment_method_types`.
   - Store PaymentIntent with `status: "requires_confirmation"`.
   - Record a transaction in the shared ledger.
   - Dispatch `payment_intent.created` webhook (if enabled).

4. **Implement retrieve handler**
   - `GET /v1/payment_intents/{id}`.
   - Return stored PaymentIntent or Stripe-compatible 404 error.

5. **Implement confirm handler**
   - `POST /v1/payment_intents/{id}/confirm`.
   - Accept test payment method tokens (`pm_card_*`, `pm_fpx_*`).
   - For FPX: transition to `requires_action`, set `next_action.redirect_to_url.url` to `/_admin/stripe/payment_intent/{id}`.
   - For card: transition to `succeeded`, dispatch `payment_intent.succeeded`.

6. **Implement cancel handler**
   - `POST /v1/payment_intents/{id}/cancel`.
   - Transition to `canceled`, dispatch `payment_intent.canceled`.

7. **Implement admin authentication page**
   - `GET /_admin/stripe/payment_intent/{id}` renders FPX bank selector or card confirm page.
   - `POST /_admin/stripe/payment_intent/{id}` processes outcome.
   - Include CSRF tokens.

8. **Update provider wiring**
   - Add PaymentIntent and admin routes.
   - Add baseURL fallback from request host/scheme.

9. **Update webhooks**
   - Build `payment_intent.*` event payloads.
   - Respect enabled-events filter from webhook config UI.
   - Sign with existing `SignPayload`.

10. **Update OpenAPI spec**
    - Add PaymentIntent endpoints to `docs/openapi.yaml`.

11. **Test**
    - Unit tests for create/retrieve/confirm/cancel/error paths.
    - FPX redirect test.
    - Test payment method token test.
    - CSRF failure test.
    - Smoke test using SDK-style HTTP calls.
    - Regression test for Checkout Sessions.

## Verification

After each sub-step, run:

```bash
go build ./...
go test ./internal/stripe/...
```

After each full prompt:

```bash
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```
