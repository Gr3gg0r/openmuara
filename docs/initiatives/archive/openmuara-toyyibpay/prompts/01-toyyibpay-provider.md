> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# Prompt 01 — ToyyibPay Provider

## Goal
Implement a faithful ToyyibPay provider in OpenMuara, including categories, bills, a local payment page, and return/callback handling with MD5 signature verification.

## Acceptance Criteria

### Provider scaffold

- [ ] Create `internal/toyyibpay/provider.go` with a `Provider` struct implementing `provider.Provider`.
- [ ] Register the provider under the name `toyyibpay`.
- [ ] Provider `Init` requires:
  - `user_secret_key` (string) — used for MD5 callback hash and as the API secret
  - `category_code` (string, optional default)

### Categories API

- [ ] Create `internal/toyyibpay/category_types.go`.
- [ ] Implement:
  - `POST /index.php/api/createCategory`
  - `POST /index.php/api/getCategoryDetails`
- [ ] Store categories in an in-memory store keyed by category code.

### Bills API

- [ ] Create `internal/toyyibpay/bill_types.go`.
- [ ] Implement:
  - `POST /index.php/api/createBill`
  - `POST /index.php/api/getBillTransactions`
  - `POST /index.php/api/inactiveBill`
- [ ] Store bills in an in-memory store keyed by bill code.
- [ ] Bill fields must match real ToyyibPay:
  - `billCode`, `billName`, `billDescription`, `billTo`, `billEmail`, `billPhone`
  - `billAmount` (integer in sen)
  - `billStatus` (`1` active, `2` inactive, etc.)
  - `categoryCode`
  - `billReturnUrl`, `billCallbackUrl`
  - `billPaymentChannel` (`0` FPX, `1` card, `2` both)
  - `billExpiryDate`, `billExpiryDays`
  - `billPriceSetting`, `billPayorInfo`
- [ ] `billPaymentLink` points to local OpenMuara payment page `/_admin/toyyibpay/pay/{billCode}`.

### Request format

- [ ] All ToyyibPay endpoints accept **form-encoded** (`application/x-www-form-urlencoded`) bodies, not JSON.
- [ ] Responses are JSON as per real ToyyibPay.

### Local payment page

- [ ] Create `internal/ui/toyyibpay-pay.html`.
- [ ] `GET /_admin/toyyibpay/pay/{billCode}` renders payment method selector respecting `billPaymentChannel` + amount + pay/cancel buttons.
- [ ] `POST /_admin/toyyibpay/pay/{billCode}` processes outcome:
  - Pay → transaction status `1` (success), record ledger transaction as `PAID`, dispatch callback.
  - Cancel → transaction status `3` (fail), dispatch callback.
  - Include CSRF token.

### Return URL and callback

- [ ] Implement browser **return URL** handler at `GET /toyyibpay/return`:
  - Redirect to the bill's `billReturnUrl` with query params `status_id`, `billcode`, `order_id`, etc.
- [ ] Implement server-side **callback** dispatch to `billCallbackUrl`:
  - POST form-encoded fields: `refno`, `status`, `reason`, `billcode`, `order_id`, `amount`, `transaction_time`, and `hash`.
  - `hash` = `MD5(userSecretKey + status + order_id + refno + "ok")`.
- [ ] Provider `PayloadBuilder` returns the form-encoded callback body bytes.
- [ ] Provider `PayloadHeaders` returns `Content-Type: application/x-www-form-urlencoded`.

### MD5 signature helper

- [ ] Create `internal/toyyibpay/signature.go`.
- [ ] Implement `MD5(userSecretKey + status + order_id + refno + "ok")`.
- [ ] Verify hash on any incoming callback payload.

### Tests

- [ ] Create `internal/toyyibpay/*_test.go` covering:
  - Category create/getCategoryDetails
  - Bill create/transactions/inactive
  - Form-encoded request parsing
  - Payment page render and outcome
  - Return URL redirect
  - Callback hash generation and verification
  - Invalid secret key errors
- [ ] Update `scripts/smoke-test.sh` with a ToyyibPay happy path.

### OpenAPI

- [ ] Update `docs/openapi.yaml` and `internal/server/openapi.yaml` with ToyyibPay endpoints.

## Files to Create/Change

- `internal/toyyibpay/provider.go`
- `internal/toyyibpay/category.go`
- `internal/toyyibpay/category_types.go`
- `internal/toyyibpay/bill.go`
- `internal/toyyibpay/bill_types.go`
- `internal/toyyibpay/signature.go`
- `internal/toyyibpay/webhook.go`
- `internal/toyyibpay/*_test.go`
- `internal/ui/toyyibpay-pay.html`
- `internal/ui/embed.go`
- `internal/server/router.go`
- `docs/openapi.yaml`
- `internal/server/openapi.yaml`
- `scripts/smoke-test.sh`

## Response / Webhook Shape

Return:
1. Category create/getCategoryDetails JSON shapes
2. Bill create/get transactions/inactive JSON shapes
3. Local payment page HTML shape
4. Return URL query params
5. Callback form-encoded payload with MD5 `hash`

## Test Notes

- `go test ./internal/toyyibpay/...`
- `./scripts/smoke-test.sh`
- Verify existing provider tests still pass.
