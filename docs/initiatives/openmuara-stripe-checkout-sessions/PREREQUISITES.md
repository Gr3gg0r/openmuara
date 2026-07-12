# OpenMuara Stripe FPX & Card Payments — Prerequisites

Before starting this initiative, confirm:

- [ ] You are on the `dev` branch.
- [ ] `go test ./...` passes on `dev`.
- [ ] You understand the existing Stripe provider in `internal/stripe/`.
- [ ] You have read `REFERENCES.md` for Stripe Checkout Session and PaymentIntents documentation.
- [ ] You understand that this initiative removes the custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes.
- [ ] You understand the existing webhook dispatcher in `internal/webhook/`.
