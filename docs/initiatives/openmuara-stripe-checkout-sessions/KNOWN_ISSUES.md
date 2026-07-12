# OpenMuara Stripe FPX & Card Payments — Known Issues

## Pre-existing

- The OpenMuara-native `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes were added as a temporary test API. This initiative removes them.
- The existing Checkout Session response sets `url` to `/v1/checkout/sessions/{id}/pay`, but no handler is registered for that path. P01 implements it.

## Out of Scope

- Payment methods other than FPX and card.
- SetupIntents, Customers, Charges, Refunds.
- Real 3-D Secure cryptography.
- Payment method saving / reuse.
- Stripe Connect.
- Multi-item carts or `mode=setup`.
- Exact HTML styling parity with Stripe Checkout or Stripe.js.
