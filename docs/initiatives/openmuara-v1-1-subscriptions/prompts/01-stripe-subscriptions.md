> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 01 — Stripe Subscriptions

## Goal

Emulate Stripe Billing subscription lifecycle endpoints.

## Acceptance Criteria

- [ ] `POST /v1/products` — create a product
- [ ] `GET /v1/products/{id}` — retrieve a product
- [ ] `POST /v1/prices` — create a price
- [ ] `GET /v1/prices/{id}` — retrieve a price
- [ ] `POST /v1/customers` — create a customer
- [ ] `GET /v1/customers/{id}` — retrieve a customer
- [ ] `POST /v1/subscriptions` — create a subscription
- [ ] `GET /v1/subscriptions/{id}` — retrieve a subscription
- [ ] `POST /v1/subscriptions/{id}` — update/cancel a subscription
- [ ] `GET /v1/invoices` — list invoices
- [ ] `GET /v1/invoices/{id}` — retrieve an invoice
- [ ] Subscription state persisted in SQLite
- [ ] Webhook dispatch for subscription events (`customer.subscription.created`, `invoice.payment_succeeded`, etc.)

## Files to Create/Change

- `internal/stripe/product.go`
- `internal/stripe/price.go`
- `internal/stripe/customer.go`
- `internal/stripe/subscription.go`
- `internal/stripe/invoice.go`
- `internal/engine/subscription.go` (or new `internal/subscription/` package)
- `internal/store/migrations/005_subscriptions.sql`

## Response Shape

Return:

1. Stripe-style Product, Price, Customer, Subscription, and Invoice objects
2. Webhook event types and payloads

## Test Notes

- `go test ./internal/stripe/...`
- Verify subscription create → invoice → payment succeeded lifecycle

## v1 / v2 Boundaries

- Do not change existing Stripe Checkout single-charge behavior.
- Do not implement mobile receipt validation or RevenueCat.
