# Prompt 14 — RevenueCat Emulation

> **Status:** ❄️ Moved to v2. This prompt is superseded by `docs/initiatives/openmuara-v2-revenuecat/prompts/01-revenuecat-emulation.md`.
> **Reason:** v1 focuses on single charge item emulation. RevenueCat (subscriptions, offerings, entitlements, mobile receipts) belongs to v2.

## Goal
Emulate RevenueCat subscriber status and entitlement endpoints.

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
