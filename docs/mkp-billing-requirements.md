---
id: mkp-billing-requirements
title: MKP Billing Requirements for OpenMuara
---

# MKP Billing Requirements for OpenMuara

> **Consumer:** Muslim Kids Platform (MKP) v2
> **Consumer Repo:** `<consumer-repo>/`
> **Consumer Branch:** `joyful-pony`
> **Date:** 2026-06-25
> **Status:** Requirements gathering — partially implemented in OpenMuara

---

## Purpose

This document lists what MKP needs from OpenMuara so that MKP can eventually **eject its internal
billing simulator** and rely on OpenMuara as the company-wide local payment virtualization layer.

MKP is currently building an internal simulator as a short-term measure. Once OpenMuara covers
the providers and scenarios below, MKP will remove the internal simulator and point its webhook
handlers at OpenMuara.

---

## MKP Billing Stack

MKP v2 API uses three payment gateways and two user-journey types:

| Gateway | Use Case | MKP Handler Path |
|---|---|---|
| Stripe | Web checkout, subscriptions, invoices | `services/mkp-v2-api/internal/api/billing/stripe.go` |
| Fawry | Regional payments (Egypt / MENA) | `services/mkp-v2-api/internal/api/billing/fawry.go` |
| RevenueCat | Mobile subscriptions (iOS / Android) | `services/mkp-v2-api/internal/api/billing/revenuecat.go` |

### User Journey Types

MKP `pricing_packages.billing_type` currently defines three values, but only two are product-relevant:

| Billing Type | Journey | Status |
|---|---|---|
| `recurring` | Subscription | Active — front-end package page only shows monthly/yearly subscription packages. |
| `one_time` | Prepaid one-time | Supported for future/pre-paid flows. |
| `lifetime` | Prepaid lifetime | **Deprecated / pending removal** — not used by current product. |

OpenMuara must support **subscription** and **prepaid** journeys for Stripe, Fawry, and RevenueCat.
`lifetime` emulation is not required.

### Current OpenMuara Coverage

| Gateway | OpenMuara Status | Notes |
|---|---|---|
| Fawry | ✅ Implemented | Charge + escape page + webhook receiver + payment-status query exist. |
| Stripe | ✅ Implemented | Checkout sessions, PaymentIntents, customers, subscriptions, invoices, and webhooks are available under `/v1/...`. |
| RevenueCat | ❄️ Frozen for v2 | Mobile subscription emulation deferred to v2. |

---

## Requirements by Gateway

### 1. Stripe Emulation

MKP needs high-fidelity Stripe emulation for local development and CI.

#### Implemented Endpoints

| Stripe Endpoint | OpenMuara Mirror | Purpose |
|---|---|---|
| `POST /v1/checkout/sessions` | `POST /v1/checkout/sessions` | Create checkout session |
| `GET /v1/checkout/sessions/:id` | `GET /v1/checkout/sessions/:id` | Retrieve session |
| `GET /v1/checkout/sessions/:id/pay` | `GET /v1/checkout/sessions/:id/pay` | Payment page simulation |
| `POST /v1/checkout/sessions/:id/pay` | `POST /v1/checkout/sessions/:id/pay` | Submit simulated payment |
| `POST /v1/payment_intents` | `POST /v1/payment_intents` | Create PaymentIntent |
| `GET /v1/payment_intents/:id` | `GET /v1/payment_intents/:id` | Retrieve PaymentIntent |
| `POST /v1/payment_intents/:id/confirm` | `POST /v1/payment_intents/:id/confirm` | Confirm PaymentIntent |
| `POST /v1/payment_intents/:id/cancel` | `POST /v1/payment_intents/:id/cancel` | Cancel PaymentIntent |
| `POST /v1/webhook` | `POST /v1/webhook` | Incoming Stripe webhook receiver |

#### Not Yet Implemented

| Stripe Endpoint | Status |
|---|---|
| `POST /v1/customers` | Planned for future Stripe hardening |
| `GET /v1/customers/:id` | Planned for future Stripe hardening |
| `POST /v1/subscriptions` | Planned for future Stripe hardening |
| `GET /v1/subscriptions/:id` | Planned for future Stripe hardening |
| `POST /v1/invoices` | Planned for future Stripe hardening |
| `GET /v1/invoices/:id` | Planned for future Stripe hardening |

#### Required Webhooks

MKP listens for these Stripe webhook events. The exact event set depends on the journey:

**Subscription (`recurring`):**

- `checkout.session.completed`
- `checkout.session.async_payment_succeeded`
- `checkout.session.async_payment_failed`
- `invoice.paid`
- `invoice.payment_failed`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`

**Prepaid (`one_time`):**

- `checkout.session.completed`
- `payment_intent.succeeded`
- `payment_intent.payment_failed`

OpenMuara must dispatch these with a valid Stripe-Signature header so MKP's `stripe.go` webhook handler accepts them.

#### Required Escape Page

A simple `/_admin/stripe-escape?session_id=xxx` page with:

- **[Simulate Success]** → triggers `checkout.session.completed` + `invoice.paid` + `customer.subscription.created`
- **[Simulate Failure]** → triggers `checkout.session.async_payment_failed` or `invoice.payment_failed`
- **[Simulate Cancel]** → marks session status `canceled`, no subscription created

#### Stripe Required Configuration

```yaml
stripe:
  api_key: "muara-stripe-test-key"
  webhook_secret: "muara-stripe-webhook-secret"
  publishable_key: "muara-stripe-publishable-key"
```

---

### 2. Fawry Emulation

MKP already uses OpenMuara's Fawry emulation for local testing. The following gaps should be closed.

#### Fawry Implemented Endpoints

| Endpoint | Status |
|---|---|
| `POST /fawry/charge` | ✅ Implemented |
| `POST /fawry/v1/charge` | ✅ Implemented |
| `POST /fawry/v2/charge` | ✅ Implemented |
| `GET /fawry/payment-status` | ✅ Implemented (signed payment-status query) |
| `POST /fawry/webhook` | ✅ Implemented (receiver) |
| `POST /fawry/v1/webhook` | ✅ Implemented |
| `POST /fawry/v2/webhook` | ✅ Implemented |
| `GET /_admin/fawry-escape` | ✅ Implemented |
| `POST /_admin/fawry-escape` | ✅ Implemented (escape action updates ledger + dispatches webhook) |

#### Closed / Implemented

1. **Outgoing webhook dispatch** — `POST /_admin/fawry-escape` now updates the ledger and dispatches
   a signed Fawry V2 notification to the configured webhook target.
2. **Reference to webhook correlation** — outgoing webhooks include the same `merchantRefNumber` used in `/fawry/charge`.

#### Remaining Gaps

1. **Configurable response delay** — `fawry.response_delay_ms` is not yet implemented.
2. **Status values** — `PAID`, `UNPAID`, and `CANCELED` are supported; `EXPIRED` is not yet a first-class status.
3. **Subscription vs Prepaid** — `/fawry/charge` does not yet accept a `billing_type` hint; all
   charges are treated as one-time payments.
4. **Extended Fawry states** — MKP-specific state transitions (e.g., refund, partial capture) are not yet modeled.

#### Fawry Required Configuration

```yaml
fawry:
  merchant_code: "muara-merchant-code"
  merchant_security_key: "muara-fawry-secret"
  webhook_secret: "muara-webhook-secret"
  response_delay_ms: 0
```

---

### 3. RevenueCat Emulation

MKP needs a shadow RevenueCat entitlement layer for mobile subscription tests.

#### Required Endpoints

| RevenueCat Endpoint | OpenMuara Mirror | Purpose |
|---|---|---|
| `GET /v1/subscribers/:app_user_id` | `GET /v1/revenuecat/subscribers/:app_user_id` | Retrieve CustomerInfo |
| `POST /v1/receipts` | `POST /v1/revenuecat/receipts` | Submit App Store / Play Store receipt |
| `POST /v1/events` | `POST /v1/revenuecat/events` (incoming) | Receive RevenueCat server-to-server events |

#### Required Webhooks (Outgoing from OpenMuara to MKP)

MKP listens for these RevenueCat events. The event semantics depend on the journey:

**Subscription (`recurring`):**

- `INITIAL_PURCHASE`
- `RENEWAL`
- `CANCELLATION`
- `UNCANCELLATION`
- `BILLING_ISSUE`
- `EXPIRATION`

**Prepaid (`one_time`):**

- `INITIAL_PURCHASE`
- `BILLING_ISSUE` (rare)
- No `RENEWAL` or `EXPIRATION` events

OpenMuara must dispatch these with a valid RevenueCat-compatible signature (or V2 authorization
header) so MKP's `revenuecat.go` accepts them.

#### Required Receipt Behavior

- App Store receipts and Play Store purchase tokens can be arbitrary strings.
- OpenMuara maps a receipt/token to an entitlement state (`premium_access: true/false`).
- Same receipt/token submitted twice returns the same subscriber state (idempotency).
- `billing_type` hint (from payload or simulator config) determines whether the purchase is recurring or non-renewing (`one_time`).

#### Required CustomerInfo Shape

**Subscription (`recurring`):**

```json
{
  "request_date": "2026-06-25T00:00:00Z",
  "subscriber": {
    "original_app_user_id": "user_123",
    "subscriptions": {
      "com.mkp.premium.monthly": {
        "expires_date": "2026-07-25T00:00:00Z",
        "purchase_date": "2026-06-25T00:00:00Z",
        "store": "app_store",
        "is_sandbox": true
      }
    },
    "entitlements": {
      "premium_access": {
        "expires_date": "2026-07-25T00:00:00Z",
        "product_identifier": "com.mkp.premium.monthly",
        "purchase_date": "2026-06-25T00:00:00Z"
      }
    }
  }
}
```

**Prepaid (`one_time`):**

```json
{
  "request_date": "2026-06-25T00:00:00Z",
  "subscriber": {
    "original_app_user_id": "user_123",
    "non_subscriptions": {
      "com.mkp.premium.onetime": [
        {
          "id": "txn_123",
          "purchase_date": "2026-06-25T00:00:00Z",
          "store": "app_store",
          "is_sandbox": true
        }
      ]
    },
    "entitlements": {
      "premium_access": {
        "product_identifier": "com.mkp.premium.onetime",
        "purchase_date": "2026-06-25T00:00:00Z"
      }
    }
  }
}
```

For **prepaid** products, the purchase appears under `non_subscriptions` and the entitlement has no `expires_date`.

#### RevenueCat Required Configuration

```yaml
revenuecat:
  api_key: "muara-revenuecat-api-key"
  webhook_secret: "muara-revenuecat-webhook-secret"
  webhook_auth_version: "v2"
```

---

## Cross-Cutting Requirements

### 1. Idempotency

OpenMuara must respect `Idempotency-Key` headers and return the same response for duplicate keys.
This is critical for MKP's webhook handlers.

### 2. Deterministic References

OpenMuara should accept an optional `reference` or `client_reference_id` and use it in returned
objects. MKP tests assert against known reference values.

### 3. Webhook Dispatch to MKP

OpenMuara must be able to POST outgoing webhooks to MKP's local webhook endpoints:

```yaml
webhook:
  url: "http://127.0.0.1:8080/api/v1/webhooks/stripe"
  max_retries: 3
```

Support for **multiple webhook targets** (one per gateway) would be ideal:

```yaml
webhook:
  targets:
    stripe: "http://127.0.0.1:8080/api/v1/webhooks/stripe"
    fawry: "http://127.0.0.1:8080/api/v1/webhooks/fawry"
    revenuecat: "http://127.0.0.1:8080/api/v1/webhooks/revenuecat"
```

### 4. Chaos / Failure Injection

MKP wants to test unhappy paths:

| Behavior | Header / Config |
|---|---|
| Delayed webhook | `X-Muara-Delay: 5000` |
| Duplicate webhook | `X-Muara-Behavior: double_spend` |
| Failed webhook | `X-Muara-Behavior: webhook_500` |
| Invalid signature | `X-Muara-Behavior: bad_signature` |
| Expired subscription | `X-Muara-Behavior: expired_subscription` |

### 5. Admin Inspection

A web dashboard or CLI to:

- List recent checkout sessions / charges / receipts.
- Inspect dispatched webhooks and responses.
- Replay a webhook by reference.
- Fast-forward subscription renewal/expiration for time-sensitive tests.

### 6. Docker / CI Support

MKP CI should be able to run OpenMuara as a service:

```yaml
services:
  muara:
    image: muara:latest
    ports:
      - "9000:9000"
    environment:
      - MUARA_STRIPE_API_KEY=muara-stripe-test-key
      - MUARA_STRIPE_WEBHOOK_SECRET=muara-stripe-webhook-secret
      - MUARA_FAWRY_MERCHANT_CODE=muara-merchant-code
      - MUARA_FAWRY_MERCHANT_SECURITY_KEY=muara-fawry-secret
      - MUARA_REVENUECAT_API_KEY=muara-revenuecat-api-key
```

---

## MKP → OpenMuara Migration Plan

Once OpenMuara meets the requirements above, MKP will:

1. Remove `services/mkp-v2-api/internal/billing/simulator/`.
2. Update MKP dev/CI config to point webhook URLs at OpenMuara.
3. Replace simulator tests with OpenMuara-backed integration tests.
4. Document the new local testing flow in MKP's runbooks.

---

## Open Questions for OpenMuara Team

1. What is the target timeline for Stripe and RevenueCat emulation?
2. Should OpenMuara expose a single `/v1/*` prefix per provider or keep provider-specific paths like `/fawry/*`?
3. Does OpenMuara plan to support subscription lifecycle simulation (renewal, expiration) out of the box?
4. Will OpenMuara provide an SDK or just HTTP endpoints?

---

## References

- MKP Billing Project: `<consumer-repo>/docs/projects/billing-cleanup-and-strategy/`
- MKP Billing Handlers: `<consumer-repo>/services/mkp-v2-api/internal/api/billing/`
- OpenMuara Feature Roadmap: `.agents/feature/FEATURE_LIST.md`
