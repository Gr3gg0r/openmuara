# Prompt 07 — Stripe Checkout Provider

## Goal
Complete Stripe Checkout session emulation.

## Acceptance Criteria
- [ ] `POST /v1/checkout/sessions` creates session, stores in SQLite, records transaction
- [ ] `GET /v1/checkout/sessions/{id}` retrieves session
- [ ] `POST /_admin/stripe/success` simulates success, updates ledger, dispatches webhook
- [ ] Stripe webhook signature header (`Stripe-Signature`) generated correctly
- [ ] Session payload shape matches Stripe API subset

## Files to Create/Change
- `internal/stripe/checkout.go`
- `internal/stripe/provider.go`
- `internal/stripe/signature.go`
- `internal/stripe/webhook.go`
- `internal/stripe/types.go`

## Response Shape
Return:
1. Session object shape
2. Create/retrieve request shapes
3. Webhook event shape
4. Signature scheme

## Test Notes
- `go test ./internal/stripe/...`
- Verify webhook signature with Stripe test secret
