# Prompt 11 — Pagination, CORS & CSRF

## Goal
Add standard API conveniences for admin and public endpoints.

## Acceptance Criteria
- [ ] Pagination for list endpoints:
  - Query params: `limit`, `cursor` (or `offset`)
  - Response envelope: `{ "data": [], "next_cursor": "...", "has_more": true }`
- [ ] CORS config in `.muara/config.yml`:
  - `allowed_origins`, `allowed_methods`, `allowed_headers`, `allow_credentials`
- [ ] CSRF double-submit cookie for `/_admin/` web UI:
  - Cookie `openmuara_csrf`
  - Header `X-CSRF-Token` validation on mutating requests
- [ ] Admin UI forms include CSRF token

## Files to Create/Change
- `internal/server/pagination.go`
- `internal/server/cors.go`
- `internal/server/csrf.go`
- `internal/config/config.go` — CORS/CSRF config
- `web/` — CSRF meta tag injection

## Response Shape
Return:
1. Pagination envelope shape
2. CORS config example
3. CSRF flow diagram

## Test Notes
- `go test ./internal/server/... -run 'Pagination|CORS|CSRF'`
- Verify CORS headers with preflight
