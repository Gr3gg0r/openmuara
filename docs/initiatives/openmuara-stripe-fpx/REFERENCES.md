# OpenMuara Stripe FPX — References

> **Status:** ❄️ Archived / Superseded. For current Stripe FPX/card emulation, see [`docs/initiatives/openmuara-stripe-checkout-sessions/REFERENCES.md`](../openmuara-stripe-checkout-sessions/REFERENCES.md).

## Stripe Documentation

- Stripe FPX payment method: https://stripe.com/docs/payments/fpx
- Stripe FPX bank list: https://stripe.com/docs/payments/fpx/accept-a-payment
- Stripe Checkout Sessions API: https://stripe.com/docs/api/checkout/sessions
- Stripe `checkout.session.completed` webhook: https://stripe.com/docs/api/events/types#event_types-checkout.session.completed
- Stripe `payment_intent.canceled` webhook: https://stripe.com/docs/api/events/types#event_types-payment_intent.canceled

## Project Docs

- `AGENTS.md` — workspace rules and quality gates
- `docs/openapi.yaml` — current API spec
- `internal/stripe/provider.go` — existing Stripe provider
- `internal/stripe/checkout.go` — existing Checkout session handler
- `internal/stripe/types.go` — existing Stripe types
- [`docs/initiatives/openmuara-stripe-checkout-sessions/`](../openmuara-stripe-checkout-sessions/) — successor initiative
