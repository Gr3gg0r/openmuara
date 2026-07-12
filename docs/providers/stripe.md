---
id: stripe
title: Stripe Provider
---

# Stripe Provider

Emulates Stripe Checkout Sessions and PaymentIntents for local development.

## Configuration

```yaml
providers:
  stripe:
    enabled: true
    config:
      publishable_key: pk_test_muara
      secret_key: sk_test_muara
      webhook_secret: whsec_muara
```

## First request

Create a Checkout Session:

```bash
curl -X POST http://127.0.0.1:9000/v1/checkout/sessions \
  -u sk_test_muara: \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method_types": ["card"],
    "mode": "payment",
    "success_url": "http://127.0.0.1/success",
    "cancel_url": "http://127.0.0.1/cancel",
    "line_items": [{
      "price_data": {
        "currency": "usd",
        "unit_amount": 1000,
        "product_data": {"name": "Test product"}
      },
      "quantity": 1
    }]
  }'
```

Expected response (abbreviated):

```json
{
  "id": "cs_test_...",
  "object": "checkout.session",
  "mode": "payment",
  "status": "open",
  "url": "http://127.0.0.1:9000/v1/checkout/sessions/cs_test_.../pay"
}
```

## Routes

| Method | Route | Purpose |
|---|---|---|
| POST | `/v1/checkout/sessions` | Create a Checkout Session |
| GET | `/v1/checkout/sessions/{id}` | Retrieve a session |
| GET | `/v1/checkout/sessions/{id}/pay` | Open local checkout page |
| POST | `/v1/checkout/sessions/{id}/pay` | Complete payment |
| POST | `/v1/payment_intents` | Create a PaymentIntent |
| GET | `/v1/payment_intents/{id}` | Retrieve a PaymentIntent |
| POST | `/v1/payment_intents/{id}/confirm` | Confirm a PaymentIntent |
| POST | `/v1/payment_intents/{id}/cancel` | Cancel a PaymentIntent |
| POST | `/v1/webhook` | Incoming Stripe webhook receiver |

## Signature algorithm

### Incoming webhook (`/v1/webhook`)

OpenMuara validates the `Stripe-Signature` header:

```text
Stripe-Signature: t=<unix-timestamp>,v1=<hex-hmac>
hex-hmac = HMAC-SHA256("<timestamp>.<payload>", webhookSecret)
```

Example verification in Go:

```go
err := stripe.VerifySignature(payload, header, "whsec_muara")
```

The timestamp must be within 5 minutes of now.

### Outgoing webhook

When OpenMuara dispatches a Stripe event, it generates the same `Stripe-Signature` header using the configured `webhook_secret`.

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/v1/checkout/sessions/{id}/pay` | Render payment page for a session |
| POST | `/v1/checkout/sessions/{id}/pay` | Submit success/fail for the session |
| GET | `/_admin/stripe/payment_intent/{id}` | Render PaymentIntent simulation page |
| POST | `/_admin/stripe/payment_intent/{id}` | Submit success/fail for the PaymentIntent |
| POST | `/_admin/stripe/success` | Admin success simulation endpoint |
| POST | `/_admin/stripe/fail` | Admin failure simulation endpoint |
| POST | `/_admin/stripe/cancel` | Admin cancel simulation endpoint |

Complete a session payment:

```bash
curl -X POST http://127.0.0.1:9000/v1/checkout/sessions/cs_test_xxx/pay \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "outcome=success"
```

## Webhooks

Outgoing webhook payload example:

```json
{
  "id": "evt_...",
  "object": "event",
  "type": "checkout.session.completed",
  "data": {
    "object": {
      "id": "cs_test_...",
      "status": "complete"
    }
  }
}
```

Headers:

```text
Stripe-Signature: t=1234567890,v1=...
```

Enabled event types include:

- `checkout.session.completed`
- `checkout.session.expired`
- `payment_intent.created`
- `payment_intent.succeeded`
- `payment_intent.canceled`

## Common errors

| HTTP status | Stripe error type | Cause | Fix |
|---|---|---|---|
| 400 | `invalid_request_error` | Missing required parameter | Check form fields |
| 401 | `authentication_error` | Invalid secret key | Use `sk_test_muara` |
| 404 | `resource_missing` | Session or PaymentIntent not found | Use the ID returned by create |
| 409 | `idempotency_error` / `invalid_state` | Duplicate idempotency key or invalid transition | Use a fresh key or reference |
| 500 | `api_error` | Internal simulation error | Check server logs |

## See also

- `docs/mkp-billing-requirements.md` — MKP v2 Stripe integration requirements.
- `runbooks/stripe-fpx-card.md` — provider-specific runbook.
