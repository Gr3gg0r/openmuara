# Prompt 12 — SenangPay Provider

## Goal
Add SenangPay gateway emulation (Malaysian payment provider).

## Acceptance Criteria
- [x] `POST /senangpay/charge` endpoint
- [x] MD5 signature verification matching the emulated SenangPay scheme
- [x] Callback handling (`GET /senangpay/callback`)
- [x] Backend webhook handling (`POST /senangpay/webhook`)
- [x] Provider registered in provider registry
- [x] Gateway YAML in `plugins/senangpay/gateway.yml`

## Files Changed
- `internal/senangpay/charge.go`
- `internal/senangpay/signature.go`
- `internal/senangpay/callback.go`
- `internal/senangpay/provider.go`
- `plugins/senangpay/gateway.yml`
- `internal/config/config.go` — defaults

## Response Shape
Return:
1. SenangPay signature algorithm (MD5 concatenation)
2. Charge request/response shapes
3. Callback/webhook query parameter shapes

## Test Notes
- `go test ./internal/senangpay/...`
- Verify MD5 hash against test vectors

## Notes
- This is an emulation. If real SenangPay integration is required, validate signature scheme against live docs.
