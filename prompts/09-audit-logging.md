# Prompt 09 — Audit Logging

## Goal
Add structured audit logging for all financial and security-relevant events.

## Acceptance Criteria
- [ ] Audit events written to SQLite table `audit_logs`
- [ ] Logged events:
  - charge created
  - webhook dispatched/delivered/failed
  - provider initialized
  - config reloaded
  - receipt validated
  - admin action (replay, escape, etc.)
- [ ] Schema: id, timestamp, actor, action, resource_type, resource_id, payload, result
- [ ] Middleware injects audit logger into context
- [ ] CLI command `openmuara audit list --since --limit`
- [ ] API endpoint `GET /_admin/audit` with pagination

## Files to Create/Change
- `internal/audit/logger.go`
- `internal/audit/store.go`
- `internal/store/migrations/003_audit_logs.sql`
- `internal/server/middleware.go`
- `internal/server/audit_admin.go`
- `internal/cli/audit.go`

## Response Shape
Return:
1. Audit event schema
2. CLI/API response envelope
3. Middleware integration notes

## Test Notes
- `go test ./internal/audit/...`
- Trigger a charge and verify audit row
