> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1.1 — Subscriptions

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Add subscription emulation to OpenMuara, starting with Stripe Billing, then extending to Malaysian gateways that support recurring payments.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/v1-1-subscriptions`

---

## Initiative Structure

```
docs/initiatives/openmuara-v1-1-subscriptions/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    ├── 01-stripe-subscriptions.md
    └── 02-senangpay-subscriptions.md
```

Planning docs live in `docs/initiatives/openmuara-v1-1-subscriptions/` in the root repo.
Product-code commits to the `feat/v1-1-subscriptions` branch. Do not commit directly to `main`.

---

## Why Subscriptions in v1.1?

v1 is intentionally focused on **single charge item emulation**. v1.1 expands the emulator to cover **recurring payments and subscriptions** while keeping the same local-first, drop-in-replacement philosophy.

Stripe is the starting point because:

- Stripe Billing has the most mature and documented subscription API.
- OpenMuara already emulates Stripe Checkout in v1, so the provider scaffold exists.
- Stripe's subscription primitives (Products, Prices, Customers, Subscriptions, Invoices) are a good reference model for other providers.

After Stripe, the plan is to add subscription emulation for Malaysian gateways that already support recurring payments.

---

## Goals

1. Emulate Stripe Billing subscription lifecycle: create, update, cancel, pause, resume.
2. Persist subscription state in SQLite alongside the transaction ledger.
3. Dispatch Stripe-style subscription webhook events (`customer.subscription.created`, `invoice.payment_succeeded`, etc.).
4. Provide a reusable subscription engine abstraction that other providers can plug into.
5. Extend subscription emulation to Malaysian gateways that support recurring payments.
6. Keep v1's drop-in-base-URL philosophy: no extra auth or custom headers on provider routes.

---

## Malaysian Payment Providers with Subscription / Recurring Support

| Provider | Recurring Capability | Notes | Existing in OpenMuara |
|----------|----------------------|-------|----------------------|
| **Stripe** | Stripe Billing / Subscriptions | Full subscription lifecycle, invoices, trials, proration | ✅ Stripe Checkout |
| **SenangPay** | Recurring payments, subscription plans | Local Malaysian gateway with subscription plan support | ✅ |
| **iPay88** | Recurring payment | Established gateway with recurring support | ✅ |
| **Billplz** | Recurring billing / auto-billing | Invoice-based recurring | ✅ |
| **ToyyibPay** | Subscriptions and invoices | Payment links with subscription management | ✅ |
| **Fiuu** (formerly Razer Merchant Services / MOLPay / RMS) | Recurring payments | Formerly MOLPay; supports recurring | ❌ |
| **2C2P** | Recurring Payment Plan (RPP) | Card-based recurring payment schedules | ❌ |
| **eGHL** | Recurring payment | Regional gateway with recurring support | ❌ |
| **Curlec** (Razorpay Curlec) | Direct debit / subscriptions | Subscription specialists with mandate-based recurring | ❌ |
| **HitPay** | Native recurring billing | Dashboard-based subscriptions, broad e-wallet support | ❌ |
| **Xendit** | Supported | SEA payment stack with recurring | ❌ |
| **Airwallex** | Via Airwallex Bill Pay | Card/bank transfer recurring | ❌ |

Source references:

- [Best Recurring Billing Software in Malaysia (2026)](https://hitpayapp.com/blog/best-recurring-billing-software-malaysia)
- [14 Payment Gateways for e-Commerce Websites in Malaysia](https://kodedigital.expert/blog/2021/04/12/14-payment-gateways-for-e-commerce-websites-in-malaysia/)
- [The 6 Malaysia Payment Gateways for LMS and E-commerce](https://pukunui.com/best-malaysia-payment-gateways-for-lms/)
- [2C2P Recurring Payment Plan (RPP)](https://developer.2c2p.com/docs/sdk-recurring-payment-plan)
- [Curlec vs Fiuu: Fees, features & verdict for Malaysia](https://www.airwallex.com/my/blog/curlec-vs-fiuu)

---

## Target Stripe Endpoints (Phase 1)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/products` | Create a product |
| `GET`  | `/v1/products/{id}` | Retrieve a product |
| `POST` | `/v1/prices` | Create a price |
| `GET`  | `/v1/prices/{id}` | Retrieve a price |
| `POST` | `/v1/customers` | Create a customer |
| `GET`  | `/v1/customers/{id}` | Retrieve a customer |
| `POST` | `/v1/subscriptions` | Create a subscription |
| `GET`  | `/v1/subscriptions/{id}` | Retrieve a subscription |
| `POST` | `/v1/subscriptions/{id}` | Update/cancel a subscription |
| `GET`  | `/v1/invoices` | List invoices |
| `GET`  | `/v1/invoices/{id}` | Retrieve an invoice |

---

## Non-Goals

- Do not implement mobile receipt validation (App Store / Play Store) in v1.1.
- Do not implement RevenueCat in v1.1 (moved to v2).
- Do not add new Malaysian gateways in Phase 1; focus on Stripe first.
- Do not change v1 single-charge provider behavior.

---

## Reference

- Stripe Billing docs: https://stripe.com/docs/billing
- Original v1 single-charge decision: root `DECISIONS.md` D037
