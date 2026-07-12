# Prompt 13 — Receipt Validation Framework

## Goal
Build a framework for validating mobile purchase receipts.

## Acceptance Criteria
- [ ] Receipt endpoint `POST /v1/receipts/validate`
- [ ] Supports Apple App Store and Google Play receipt shapes
- [ ] Receipts treated as lookup keys in local dataset (no real crypto validation)
- [ ] Dataset path configurable (`.muara/data/unified_matrix.json`)
- [ ] Response matches provider-specific shape
- [ ] Webhook dispatch for renewal events

## Files to Create/Change
- `internal/receipt/validator.go`
- `internal/receipt/apple.go`
- `internal/receipt/google.go`
- `internal/receipt/store.go`
- `internal/server/router.go`
- `internal/config/config.go`

## Response Shape
Return:
1. Validation request/response shapes
2. Dataset format
3. Provider shape mappings

## Test Notes
- `go test ./internal/receipt/...`
- Validate with fixture receipts
