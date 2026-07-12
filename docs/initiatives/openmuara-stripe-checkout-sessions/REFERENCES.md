# OpenMuara Stripe FPX & Card Payments — References

## Stripe Documentation

### Checkout Sessions
- Checkout Session overview: https://stripe.com/docs/payments/checkout
- Create a Checkout Session: https://stripe.com/docs/api/checkout/sessions/create
- Retrieve a Checkout Session: https://stripe.com/docs/api/checkout/sessions/retrieve
- Checkout Session object: https://stripe.com/docs/api/checkout/sessions/object
- Checkout Session line_items: https://stripe.com/docs/api/checkout/sessions/create#create_checkout_session-line_items
- Checkout Session payment_method_types: https://stripe.com/docs/api/checkout/sessions/create#create_checkout_session-payment_method_types
- FPX with Checkout: https://stripe.com/docs/payments/fpx/checkout

### PaymentIntents
- PaymentIntents overview: https://stripe.com/docs/payments/payment-intents
- Create a PaymentIntent: https://stripe.com/docs/api/payment_intents/create
- Retrieve a PaymentIntent: https://stripe.com/docs/api/payment_intents/retrieve
- Confirm a PaymentIntent: https://stripe.com/docs/api/payment_intents/confirm
- Cancel a PaymentIntent: https://stripe.com/docs/api/payment_intents/cancel
- PaymentIntent object: https://stripe.com/docs/api/payment_intents/object
- FPX with PaymentIntents: https://stripe.com/docs/payments/fpx/accept-a-payment
- Card payments with PaymentIntents: https://stripe.com/docs/payments/accept-a-payment

### Errors
- Stripe error object: https://stripe.com/docs/api/errors

### Webhooks
- Webhook events: https://stripe.com/docs/api/events/types
- Webhook signatures: https://stripe.com/docs/webhooks/signatures

## Project Docs

- `AGENTS.md` — workspace rules and quality gates
- `docs/openapi.yaml` — current API spec
- `internal/stripe/provider.go` — existing Stripe provider
- `internal/stripe/checkout.go` — existing Checkout session handler
- `internal/stripe/types.go` — existing Stripe types
- `internal/webhook/dispatcher.go` — webhook dispatcher
