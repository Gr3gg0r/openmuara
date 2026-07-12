# Prompt 06 — Fawry Provider Hardening

## Goal
Harden the Fawry provider emulation to match Fawry Express Checkout behavior.

## Acceptance Criteria
- [x] SHA256 signature verification (`Verify`) matches Fawry Express Checkout scheme
- [ ] Charge endpoint records transaction in SQLite ledger
- [ ] Webhook endpoint validates inbound signature
- [ ] Escape page + escape action work with provider dispatcher
- [ ] Payload builder reads transaction from ledger
- [ ] Error codes match documented Fawry shape

## Files to Create/Change
- `internal/fawry/charge.go`
- `internal/fawry/signature.go`
- `internal/fawry/webhook.go`
- `internal/fawry/escape.go`
- `internal/fawry/provider.go`

## Response Shape
Return:
1. Signature algorithm description
2. Request/response shapes
3. Webhook payload shape
4. Test coverage

## Test Notes
- `go test ./internal/fawry/...`
- Verify signature against known test vectors
