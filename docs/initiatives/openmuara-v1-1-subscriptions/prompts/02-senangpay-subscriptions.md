> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 02 — SenangPay Subscriptions

## Goal

Emulate SenangPay subscription / recurring payment flow on top of the v1.1 subscription engine.

## Acceptance Criteria

- [ ] SenangPay subscription plan creation endpoint
- [ ] Charge/subscribe endpoint that creates a subscription
- [ ] Subscription status query endpoint
- [ ] Recurring payment webhook events
- [ ] Subscription state persisted in SQLite via the shared subscription engine

## Files to Create/Change

- `internal/senangpay/subscription.go`
- `internal/senangpay/provider.go` (register subscription routes)
- Reuse `internal/engine/subscription.go` (or `internal/subscription/`)

## Response Shape

Return:

1. SenangPay subscription plan object
2. Subscription status object
3. Webhook event payloads

## Test Notes

- `go test ./internal/senangpay/...`
- Verify recurring charge lifecycle

## v1 / v2 Boundaries

- Do not change SenangPay single-charge behavior.
- Do not implement mobile receipt validation or RevenueCat.

## Reference

- SenangPay recurring/subscription docs: https://senangpay.my/
