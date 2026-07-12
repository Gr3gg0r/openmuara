# Prompt 05 — Core HTTP Router

## Goal
Build the central HTTP router that dispatches provider-matched requests with middleware.

## Acceptance Criteria
- [ ] Router registers provider routes dynamically
- [ ] Middleware stack: request ID, logging, recovery, CORS (config-driven)
- [ ] `GET /healthz` returns `{"status":"ok"}`
- [x] `GET /readyz` returns 200 only when DB and required providers are ready
- [ ] Provider paths matched exactly (Stripe `/v1/checkout/sessions`, Fawry `/fawry/charge`, etc.)
- [ ] Admin routes under `/_admin/` preserved

## Files to Create/Change
- `internal/server/router.go`
- `internal/server/middleware.go`
- `internal/server/health.go`
- `internal/server/router_test.go`

## Response Shape
Return:
1. Middleware order diagram
2. Route registration snippet
3. Health/readiness criteria

## Test Notes
- `go test ./internal/server/...`
- `curl http://localhost:9000/healthz`
