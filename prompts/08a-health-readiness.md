# Prompt 08a — Health & Readiness

## Goal
Implement health and readiness endpoints.

## Acceptance Criteria
- [x] `GET /healthz` — always returns 200 `{"status":"ok"}`
- [x] `GET /readyz` — returns 200 only when:
  - SQLite connection is alive
  - All required providers initialized
- [x] Readiness failures return 503 with reason
- [x] Endpoints bypass auth/logging verbosity

## Files Changed
- `internal/server/router.go`
- `internal/server/health.go`

## Response Shape
Return:
1. `/healthz` response shape
2. `/readyz` response shape (success and failure)
3. Readiness checks list

## Test Notes
- `go test ./internal/server/... -run Health`
- `curl /readyz` with DB closed → 503
