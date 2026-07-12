# Prompt 08b — Prometheus Metrics

## Goal
Expose Prometheus metrics for observability.

## Acceptance Criteria
- [ ] `GET /metrics` endpoint returns Prometheus format
- [ ] Metrics recorded:
  - `openmuara_requests_total` (method, path, status)
  - `openmuara_request_duration_seconds` (method, path)
  - `openmuara_webhook_attempts_total` (provider, status)
  - `openmuara_transactions_total` (provider, status)
- [ ] Metrics names prefixed with `openmuara_`
- [ ] Endpoint excluded from request metrics

## Files to Create/Change
- `internal/server/metrics.go`
- `internal/server/middleware.go` — request metrics
- `internal/server/router.go` — `/metrics` registration
- `internal/webhook/dispatcher.go` — webhook counter
- `internal/engine/store.go` — transaction counter

## Response Shape
Return:
1. Metric table with names, types, labels
2. Sample `/metrics` output

## Test Notes
- `go test ./internal/server/... -run Metrics`
- Scrape `/metrics` after a request
