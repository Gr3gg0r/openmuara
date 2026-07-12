> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# Task Spec — SenangPay Signature Emulation

## Status
Implemented. This spec documents the actual MD5-based emulation used in `internal/senangpay/`.

## Goal
Emulate the SenangPay signature scheme for charge requests so that OpenMuara can stand in for SenangPay during local integration testing.

> ⚠️ This is an emulation based on observed SenangPay-style behavior. If real integration is required, validate the exact signature algorithm against live SenangPay docs.

## Signature Algorithm

SenangPay in OpenMuara uses an MD5 hash over the concatenation:

```
hash = md5(secret_key + detail + amount + order_id)
```

Where `amount` is formatted with exactly two decimal places.

Example:

```go
msg := fmt.Sprintf("%s%s%.2f%s", secret, detail, amount, orderID)
hash := fmt.Sprintf("%x", md5.Sum([]byte(msg)))
```

## Charge Flow

1. Test app POSTs to `/senangpay/charge` with:
   - `detail`
   - `amount`
   - `order_id`
   - `name`, `email`, `phone`
   - `hash`
2. OpenMuara verifies the MD5 hash.
3. OpenMuara records transaction in SQLite ledger.
4. OpenMuara returns redirect URL to emulated payment page.

## Callback / Webhook Flow

1. After "payment", browser is redirected to `/senangpay/callback` (GET) or a backend POST is sent to `/senangpay/webhook`.
2. Query parameters:
   - `status_id=1` → paid
   - `status_id=0` → unpaid
3. OpenMuara updates transaction status via `engine.Transition`.

## Config Shape

```yaml
providers:
  senangpay:
    enabled: true
    config:
      secret_key: sp_test_secret
```

## Go Helper API

```go
package senangpay

// Sign computes the SenangPay MD5 signature.
func Sign(secret, detail string, amount float64, orderID string) string

// Verify checks the SenangPay MD5 signature.
func Verify(req ChargeRequest, secret string) bool

// SignRequest is a test helper that fills req.Hash.
func SignRequest(req *ChargeRequest, secret string)
```

## Files

- `internal/senangpay/signature.go`
- `internal/senangpay/signature_test.go`
- `internal/senangpay/charge.go`
- `internal/senangpay/callback.go`
- `internal/senangpay/provider.go`
- `plugins/senangpay/gateway.yml`

## Test Vectors

- `secret = "test_secret"`
- `detail = "Test payment"`
- `amount = 10.00`
- `order_id = "ORDER-1"`
- Expected hash: computed by `Sign("test_secret", "Test payment", 10.00, "ORDER-1")`

## Acceptance Criteria

- [x] MD5 signature computation matches test vectors
- [x] Charge endpoint rejects invalid hashes
- [x] Callback/webhook endpoint updates transaction status
