# OpenMuara Stripe FPX — Known Issues

## Pre-existing

None.

## Supersession

- The custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes implemented by this initiative were removed in commit `885a14d`.
- Functionality is replaced by Stripe Checkout Sessions and PaymentIntents in `docs/initiatives/openmuara-stripe-checkout-sessions/`.

## Out of Scope

- FPX via PaymentIntents (only Checkout Sessions) — now handled by successor.
- Bank-specific redirect simulation (each bank has its own real login flow).
- FPX refunds, disputes, and chargebacks (reuse generic Stripe simulation endpoints).
- Real-time FX or MYR amount validation.
