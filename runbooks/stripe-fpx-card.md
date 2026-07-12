---
id: stripe-fpx-card
title: Stripe FPX & Card Runbook — OpenMuara
---

# Stripe FPX & Card Runbook — OpenMuara

How to create, confirm, and cancel Stripe Checkout Sessions and PaymentIntents for FPX and card payments in OpenMuara.

---

## Supported payment methods

| Method | Token prefix | Flow | Default currency |
|--------|--------------|------|------------------|
| Card | `pm_card_*` (e.g. `pm_card_visa`) | Immediate success or failure | `usd` |
| FPX | `pm_fpx_*` (e.g. `pm_fpx_maybank`) | Bank redirect via `/_admin/stripe/payment_intent/{id}` | `myr` |

### Supported FPX banks

- `maybank` / `maybank2u` → Maybank2U
- `cimb` → CIMB Clicks
- `public_bank` → Public Bank
- `rhb` → RHB Now
- `hong_leong` → Hong Leong Connect
- `ambank` → AmBank
- `bank_islam` → Bank Islam
- `affin_bank` → Affin Bank

Use the bank name as the value of the `bank` form field when completing an FPX redirect.

---

## Quick start

1. Start OpenMuara:
   ```bash
   ./bin/muara start
   ```
2. Create a Checkout Session or PaymentIntent using the examples in `examples/stripe/`.
3. Complete the payment through the returned `url` or `next_action.redirect_to_url.url`.
4. Inspect the result in the dashboard at `http://127.0.0.1:9000/_admin`.

---

## Checkout Sessions

### Create a session

```bash
curl -s -X POST http://127.0.0.1:9000/v1/checkout/sessions \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method_types": ["fpx"],
    "mode": "payment",
    "success_url": "http://localhost:8080/success?session_id={CHECKOUT_SESSION_ID}",
    "cancel_url": "http://localhost:8080/cancel",
    "line_items": [{
      "price_data": {
        "currency": "myr",
        "unit_amount": 1000,
        "product_data": {"name": "Prepaid Top-up"}
      },
      "quantity": 1
    }]
  }'
```

OpenMuara returns a Stripe-compatible session object. The `url` points to a local checkout page hosted by OpenMuara.

### Complete the payment

Open the `url` in a browser, select a bank or enter card details, and submit. On success the browser redirects to `success_url`; on cancel it redirects to `cancel_url`.

### Retrieve a session

```bash
curl -s http://127.0.0.1:9000/v1/checkout/sessions/cs_test_<id>
```

### Cancel a session before payment

POST the same `url` with `action=cancel`, or let the checkout page expire.

---

## PaymentIntents

### Create a PaymentIntent

```bash
curl -s -X POST http://127.0.0.1:9000/v1/payment_intents \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000,
    "currency": "myr",
    "payment_method_types": ["fpx"],
    "metadata": {"order_id": "order-789"}
  }'
```

### Confirm a PaymentIntent

```bash
curl -s -X POST http://127.0.0.1:9000/v1/payment_intents/pi_test_<id>/confirm \
  -H "Content-Type: application/json" \
  -d '{"payment_method": "pm_fpx_maybank"}'
```

For FPX the response status is `requires_action` and `next_action.redirect_to_url.url` points to `/_admin/stripe/payment_intent/{id}`. For cards (`pm_card_visa`) the status is `succeeded` immediately.

### Complete the FPX redirect

1. Open `next_action.redirect_to_url.url` in a browser.
2. Select a bank and submit.
3. The PaymentIntent status becomes `succeeded` and a webhook is dispatched.

### Cancel a PaymentIntent

```bash
curl -s -X POST http://127.0.0.1:9000/v1/payment_intents/pi_test_<id>/cancel
```

A PaymentIntent can only be canceled while it is in `requires_confirmation` or `requires_action`. Confirmed/succeeded PaymentIntents return `payment_intent_unexpected_state`.

---

## Currency defaults

- **Checkout Session:** inherited from `line_items[].price_data.currency`. If omitted, OpenMuara returns `invalid_request_error` on the `currency` parameter.
- **PaymentIntent:** required in the create request. If omitted, OpenMuara returns `invalid_request_error` on the `currency` parameter.
- FPX requests are typically `myr`; card requests are typically `usd`. Other currencies are accepted as long as the amount is valid.

---

## Webhook event matrix

| Resource | Action | Event emitted | Notes |
|----------|--------|---------------|-------|
| PaymentIntent | create | `payment_intent.created` | Fired asynchronously after creation. |
| PaymentIntent | confirm with card | `payment_intent.succeeded` | Immediate success. |
| PaymentIntent | confirm with FPX + admin approve | `payment_intent.succeeded` | After bank redirect is approved. |
| PaymentIntent | cancel | `payment_intent.canceled` | Only from a cancellable state. |
| Checkout Session | confirm | `checkout.session.completed` | After the local pay page approves the session. |
| Checkout Session | cancel | `checkout.session.expired` | Triggered by cancel action or expiry. |

Configure the webhook URL and signing secret in `.muara/config.yml`:

```yaml
webhook:
  url: "http://localhost:8080/webhooks/stripe"
providers:
  stripe:
    webhook_secret: "whsec_muara"
```

Outgoing webhooks include a `Stripe-Signature` header. Verify it with the configured `webhook_secret` using Stripe's signature verification logic.

---

## Common failures

### `session_invalid` / `session status is complete`

The Checkout Session has already been confirmed or canceled. Retrieve it to check `status` and `payment_status`.

### `payment_intent_unexpected_state`

The requested transition is not allowed for the current PaymentIntent status. Common causes:
- Canceling a succeeded PaymentIntent.
- Confirming a canceled PaymentIntent.

### `resource_missing` on payment method

The test payment method token is not recognized. Use one of the supported prefixes and bank/card names.

### Webhook not delivered

1. Verify `webhook.url` is set.
2. Confirm the consumer returns `2xx`.
3. Inspect attempts in the dashboard or via `muara webhook list --status failed`.
4. Replay after fixing the consumer.

---

## Example scripts

Runnable examples are in `examples/stripe/`:

| Script | Purpose |
|--------|---------|
| `checkout-fpx.sh` | Create and complete an FPX Checkout Session. |
| `payment-intent-fpx.sh` | Create, confirm, and complete an FPX PaymentIntent. |
| `payment-intent-card.sh` | Create and confirm a card PaymentIntent. |

Set `MUARA_BASE_URL` to point at a different OpenMuara instance:

```bash
MUARA_BASE_URL=http://127.0.0.1:8080 examples/stripe/payment-intent-card.sh
```

---

## References

- OpenAPI spec: `docs/openapi.yaml`
- Debugging runbook: `runbooks/debugging.md`
