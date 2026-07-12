> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Stripe FPX & Card Payments

> **Status:** 🟡 Active | **Started:** 2026-06-30
> **Scope:** Implement faithful Stripe Checkout Session and PaymentIntents API emulation for single-charge FPX and card payments, including a local OpenMuara-hosted checkout page and Stripe-compatible webhook configuration UI.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-stripe-checkout-sessions/
├── README.md              # This file
├── HOWTO.md               # Decomposition guide for AI
├── PREREQUISITES.md       # Human pre-flight checklist
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing bugs / out-of-scope
├── REFERENCES.md          # Links to specs, runbooks, vendor docs
├── .gitignore             # Ignore screenshots, logs, temp files
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   ├── 01-stripe-checkout-sessions.md
│   └── 02-stripe-payment-intents.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
├── findings/              # Research, audit output, analysis
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-stripe-checkout-sessions/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why both Checkout Sessions and PaymentIntents?

OpenMuara's mission is to let developers test financial infrastructure locally before they have real provider accounts. For Stripe, that means a developer should be able to write code against OpenMuara using the official Stripe SDK, then switch to real Stripe by changing only the base URL and API key.

The current custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes break that promise. They are OpenMuara-specific and require rewrites when moving to production.

Stripe has two primary SDK paths for one-time payments:

1. **Checkout Session** — the Stripe-hosted checkout page API. A developer creates a session with `line_items`, `payment_method_types`, `success_url`, and `cancel_url`; Stripe returns a `url` that the customer uses to pay. OpenMuara emulates this by hosting a local checkout page at that `url`.
2. **PaymentIntents** — the lower-level API used when the merchant builds their own payment UI. It provides `create`, `retrieve`, `confirm`, and `cancel`, and uses `next_action.redirect_to_url` for methods like FPX that require a bank redirect.

Both are legitimate, common Stripe SDK entry points. This initiative implements both so OpenMuara covers the full single-charge FPX/card surface that developers actually use.

This initiative focuses on **single-charge items only** (`mode=payment` for Checkout, one-time PaymentIntents). No product catalog, no subscription creation.

---

## Goals

1. Remove the OpenMuara-native FPX and card charge + escape routes.
2. **Checkout Sessions**
   - Extend `POST /v1/checkout/sessions` to accept `payment_method_types: ["fpx"]`, `["card"]`, or `["card","fpx"]`.
   - Implement the missing `GET /v1/checkout/sessions/{id}/pay` handler that renders a local OpenMuara checkout page.
   - Implement `POST /v1/checkout/sessions/{id}/pay` to process the customer's payment decision.
   - Return Stripe-compatible error JSON shapes.
   - Emit Stripe-compatible webhooks: `checkout.session.completed` and `checkout.session.expired`.
3. **PaymentIntents**
   - Implement `POST /v1/payment_intents`, `GET /v1/payment_intents/{id}`, `POST /v1/payment_intents/{id}/confirm`, `POST /v1/payment_intents/{id}/cancel`.
   - Support `payment_method_types: ["fpx"]`, `["card"]`, or `["card","fpx"]`.
   - Accept Stripe-style test payment method tokens (`pm_card_visa`, `pm_fpx_maybank`, etc.).
   - Render local authentication pages at `/_admin/stripe/payment_intent/{id}` for FPX bank selection and card confirmation.
   - Return Stripe-compatible error JSON shapes.
   - Emit Stripe-compatible webhooks: `payment_intent.created`, `payment_intent.succeeded`, `payment_intent.payment_failed`, `payment_intent.canceled`.
4. Add an admin UI at `/_admin/stripe/webhooks` for configuring Stripe-style webhook endpoints (URL, event selection, signing secret); persist changes to `.muara/config.yml`.
5. Update `docs/openapi.yaml` to document all new endpoints.
6. Ensure the official Stripe SDK can use both Checkout Sessions and PaymentIntents against OpenMuara.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Provider contract fidelity
Both APIs must match Stripe's documented behavior for the implemented subset, including:
- Request/response JSON shapes
- Stripe-compatible error JSON shapes
- Status values and transitions
- `url` / `next_action.redirect_to_url.url` pointing to local OpenMuara pages
- Webhook event types and payload shapes
- Test payment method tokens for PaymentIntents

### 3. Quality gates
Every prompt must pass:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Out of Scope

- Stripe Billing / subscriptions / products / prices APIs
- SetupIntents, Customers, Charges, Refunds
- Payment methods other than FPX and card in this initiative
- Real 3-D Secure cryptography
- Payment method saving / reuse
- Stripe Connect
- Multi-item carts or `mode=setup`

## Implemented enhancements from `openmuara-stripe-fpx`

| # | Recommendation | Status | Location |
|---|----------------|--------|----------|
| F1 | Add runnable Stripe FPX/card example scripts | ✅ Done | `examples/stripe/checkout-fpx.sh`, `examples/stripe/payment-intent-fpx.sh`, `examples/stripe/payment-intent-card.sh` |
| F2 | Add operational runbook for FPX/card flows | ✅ Done | `runbooks/stripe-fpx-card.md` |
| F3 | Expand edge-case test coverage | ✅ Done | `internal/stripe/edge_cases_test.go` |
| F4 | Document currency defaults and supported FPX banks | ✅ Done | This section and `runbooks/stripe-fpx-card.md` |
| F5 | Document the full webhook event matrix | ✅ Done | This section and `runbooks/stripe-fpx-card.md` |

---

## Currency defaults

OpenMuara does not invent a currency when one is omitted. The Stripe-compatible create endpoints return `invalid_request_error` with the missing parameter name:

- **Checkout Session:** `currency` is taken from `line_items[0].price_data.currency`.
- **PaymentIntent:** `currency` is required in the create body.

Typical test choices:

| Method | Typical currency | Notes |
|--------|------------------|-------|
| Card | `usd` | Any valid ISO currency code is accepted. |
| FPX | `myr` | Malaysian ringgit; other codes are accepted for validation but FPX is realistically MYR-only. |

---

## Supported FPX banks

The local checkout page and PaymentIntent admin page expose the following Malaysian banks for FPX emulation:

| Code (form / UI value) | Display name | Test token |
|------------------------|--------------|------------|
| `maybank2u` | Maybank2U | `pm_fpx_maybank` |
| `cimb` | CIMB Clicks | `pm_fpx_cimb` |
| `public_bank` | Public Bank | `pm_fpx_publicbank` |
| `rhb` | RHB Now | `pm_fpx_rhb` |
| `hong_leong` | Hong Leong Connect | `pm_fpx_hongleong` |
| `ambank` | AmBank | `pm_fpx_ambank` |
| `bank_islam` | Bank Islam | `pm_fpx_bankislam` |
| `affin_bank` | Affin Bank | `pm_fpx_affinbank` |

Pass the test token to `POST /v1/payment_intents/{id}/confirm` to place the PaymentIntent into `requires_action`. The admin redirect page then lets you complete the flow with the bank selector.

---

## Webhook event matrix

OpenMuara emits Stripe-compatible webhook events for the Checkout Session and PaymentIntents flows. Configure `webhook.url` and `providers.stripe.webhook_secret` in `.muara/config.yml` to receive signed payloads.

| Resource | Transition | Event | Webhook `payment_status` |
|----------|------------|-------|--------------------------|
| PaymentIntent | Created | `payment_intent.created` | `new` |
| PaymentIntent | Card confirmed | `payment_intent.succeeded` | `paid` |
| PaymentIntent | FPX confirmed via admin page | `payment_intent.succeeded` | `paid` |
| PaymentIntent | Canceled | `payment_intent.canceled` | `unpaid` |
| Checkout Session | Confirmed on pay page | `checkout.session.completed` | `paid` |
| Checkout Session | Canceled on pay page | `checkout.session.expired` | `unpaid` |

Events not implemented in this initiative remain out of scope (see Out of Scope above). All outgoing webhooks include a `Stripe-Signature` header computed with the configured `webhook_secret`.

---

## Recommendations & future enhancements

Remaining ideas for later versions:

- SetupIntents and saved payment method emulation.
- More comprehensive 3-D Secure redirect simulation for cards.
- Stripe Customer / Price / Product catalog APIs.
- Subscription (`mode=subscription`) support for Checkout Sessions.
