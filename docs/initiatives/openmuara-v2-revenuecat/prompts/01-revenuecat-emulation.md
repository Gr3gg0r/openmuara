> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 01 — RevenueCat Emulation (v2)

## Goal

Emulate RevenueCat subscriber status and entitlement endpoints for v2.

## Acceptance Criteria

- [ ] `POST /v1/subscribers/{app_user_id}` — get subscriber status
- [ ] `GET /v1/subscribers/{app_user_id}/offerings` — list offerings
- [ ] `POST /v1/receipts` — submit receipt, update entitlements
- [ ] Entitlement state persisted in SQLite
- [ ] Webhook dispatch for subscriber events

## Files to Create/Change

- `internal/revenuecat/subscriber.go`
- `internal/revenuecat/entitlement.go`
- `internal/revenuecat/receipt.go`
- `internal/revenuecat/provider.go`
- `internal/store/migrations/004_revenuecat.sql`

## Response Shape

Return:

1. Subscriber object shape
2. Offering/entitlement shapes
3. Webhook event types

## Test Notes

- `go test ./internal/revenuecat/...`
- Verify entitlement lifecycle

## v1 Boundary

This prompt is for v2 only. Do not add RevenueCat endpoints or state to v1.
