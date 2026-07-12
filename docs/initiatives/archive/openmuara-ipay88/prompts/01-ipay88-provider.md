> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# Prompt 01 — iPay88 Provider

## Goal
Implement a faithful iPay88 Malaysia classic ePayment provider in OpenMuara, including redirect payment request, local payment page, response/backend callbacks, and signature verification.

## Acceptance Criteria

### Provider scaffold

- [ ] Create `internal/ipay88/provider.go` with a `Provider` struct implementing `provider.Provider`.
- [ ] Register the provider under the name `ipay88`.
- [ ] Provider `Init` requires:
  - `merchant_code` (string)
  - `merchant_key` (string) — used for SHA256 signature

### Payment request

- [ ] Create `internal/ipay88/payment_types.go`.
- [ ] Implement `POST /ePayment/entry.asp`:
  - Accept form-encoded fields: `MerchantCode`, `PaymentId`, `RefNo`, `Amount`, `Currency`, `ProdDesc`, `UserName`, `UserEmail`, `UserContact`, `Remark`, `Lang`, `Signature`, `SignatureType`, `ResponseURL`, `BackendURL`.
  - Validate required fields and signature.
  - `SignatureType` must be `SHA256`.
  - Store payment request in an in-memory store keyed by `RefNo`.
  - Redirect customer to local payment page `/_admin/ipay88/pay/{refNo}`.
- [ ] Amount in the signature must have decimal/thousand separators stripped before hashing (e.g., `1,278.99` → `127899`).
- [ ] Request signature:
  ```
  SHA256(MerchantKey + MerchantCode + RefNo + Amount(stripped) + Currency)
  ```

### Payment method IDs

- [ ] Document a mapping of common `PaymentId` values in `DECISIONS.md`:
  - Credit/debit card: typically `1` or provider-specific
  - FPX bank IDs: bank-specific numeric codes
  - E-wallets: provider-specific numeric codes
- [ ] Local payment page may present a simplified selector that maps to these IDs.

### Local payment page

- [ ] Create `internal/ui/ipay88-pay.html`.
- [ ] `GET /_admin/ipay88/pay/{refNo}` renders payment method selector + amount + pay/cancel buttons.
- [ ] `POST /_admin/ipay88/pay/{refNo}` processes outcome:
  - Pay → `Status=1` (success), record ledger transaction as `PAID`, dispatch backend post.
  - Cancel → `Status=0` (fail), dispatch backend post.
  - Include CSRF token.

### Response and backend callbacks

- [ ] Implement `POST /ipay88/response`:
  - Build response signature.
  - POST form fields to the original `ResponseURL`.
  - Response signature:
    ```
    SHA256(MerchantKey + MerchantCode + PaymentId + RefNo + Amount(stripped) + Currency + Status)
    ```
- [ ] Implement `POST /ipay88/backend`:
  - Validate response/backend signature.
  - Return plain text `RECEIVEOK` to acknowledge receipt (real iPay88 expects this).
  - Dispatch webhook via `internal/webhook` using provider payload builder.
- [ ] Create `internal/ipay88/signature.go` with SHA256 signature helpers.
- [ ] Provider `PayloadBuilder` returns form-encoded callback body bytes.
- [ ] Provider `PayloadHeaders` returns `Content-Type: application/x-www-form-urlencoded`.

### Requery

- [ ] Implement `POST /ePayment/enquiry.asp`:
  - Accept `MerchantCode`, `RefNo`, `Amount`.
  - No signature required.
  - Return current payment status:
    - `00` = success
    - Others = failure/pending

### Security

- [ ] Validate `ResponseURL` and `BackendURL` are HTTP(S) and reject private/internal hosts to avoid SSRF.
- [ ] Never return or log `MerchantKey`.

### Tests

- [ ] Create `internal/ipay88/*_test.go` covering:
  - Payment request validation and signature (including amount normalization)
  - Payment page render and outcome
  - Response/backend signature validation
  - `RECEIVEOK` acknowledgement
  - Requery status
  - Invalid merchant key errors
  - SSRF protection on ResponseURL/BackendURL
- [ ] Update `scripts/smoke-test.sh` with an iPay88 happy path.

### OpenAPI

- [ ] Update `docs/openapi.yaml` and `internal/server/openapi.yaml` with iPay88 endpoints.

## Files to Create/Change

- `internal/ipay88/provider.go`
- `internal/ipay88/payment.go`
- `internal/ipay88/payment_types.go`
- `internal/ipay88/signature.go`
- `internal/ipay88/webhook.go`
- `internal/ipay88/*_test.go`
- `internal/ui/ipay88-pay.html`
- `internal/ui/embed.go`
- `internal/server/router.go`
- `docs/openapi.yaml`
- `internal/server/openapi.yaml`
- `scripts/smoke-test.sh`

## Response / Webhook Shape

Return:
1. Payment request redirect behavior
2. Local payment page HTML shape
3. Response URL form POST fields and signature
4. Backend URL form POST fields, signature, and `RECEIVEOK` acknowledgement
5. Requery response shape

## Test Notes

- `go test ./internal/ipay88/...`
- `./scripts/smoke-test.sh`
- Verify existing provider tests still pass.
