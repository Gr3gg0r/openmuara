# Prompt 10 — Outgoing Webhooks

## Goal
Implement reliable outgoing webhook dispatch with retries, idempotency, and replay.

## Acceptance Criteria
- [ ] Webhook dispatcher persists attempts to SQLite
- [ ] Exponential backoff with jitter
- [ ] Max retries config-driven
- [ ] Replay via CLI `openmuara webhook replay <ref>`
- [ ] Admin endpoint `GET /_admin/webhooks` lists attempts
- [ ] Provider-specific payload builder + headers
- [ ] Idempotency: duplicate dispatch for same ref+status returns existing attempt

## Files to Create/Change
- `internal/webhook/dispatcher.go`
- `internal/webhook/delivery.go`
- `internal/webhook/store.go`
- `internal/cli/webhook.go`
- `internal/server/webhook_admin.go`

## Response Shape
Return:
1. Attempt lifecycle (pending → delivered/failed)
2. Retry policy
3. Replay CLI shape

## Test Notes
- `go test ./internal/webhook/...`
- Simulate failed delivery and verify retries
