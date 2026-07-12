> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# Prompt 01 — Billplz Provider

## Goal
Implement a faithful Billplz v3 provider in OpenMuara, including collections, bills, payment methods, a local payment page, and redirect/callback handling.

## Acceptance Criteria

### Provider scaffold

- [ ] Create `internal/billplz/provider.go` with a `Provider` struct implementing the `provider.Provider` interface.
- [ ] Register the provider under the name `billplz`.
- [ ] Provider `Init` requires:
  - `api_key` (string) — used for HTTP Basic Auth username
  - `x_signature_key` (string) — used for `x_signature` generation/verification
  - `collection_id` (string, optional default)

### Collections API

- [ ] Create `internal/billplz/collection_types.go` with `Collection`, `CreateCollectionRequest`, etc.
- [ ] Implement:
  - `POST /api/v3/collections` — create collection
  - `GET /api/v3/collections/{id}` — retrieve collection
- [ ] Store collections in an in-memory store keyed by `id`.

### Bills API

- [ ] Create `internal/billplz/bill_types.go` with `Bill`, `CreateBillRequest`, etc.
- [ ] Implement:
  - `POST /api/v3/bills` — create bill (linked to a collection)
  - `GET /api/v3/bills/{id}` — retrieve bill
  - `DELETE /api/v3/bills/{id}` — delete bill
- [ ] Store bills in an in-memory store keyed by `id`.
- [ ] Bill response fields must match real Billplz v3:
  - `id`, `collection_id`, `paid`, `state`, `amount`, `description`
  - `name`, `email`, `mobile`
  - `reference_1`, `reference_1_label`, `reference_2`, `reference_2_label`
  - `callback_url`, `redirect_url`
  - `url`, `paid_amount`, `due_at`, `paid_at`
- [ ] Amount is an integer in the smallest currency unit (sen).
- [ ] Default currency is implicitly MYR; do **not** include a `currency` field in the response.
- [ ] `url` points to local OpenMuara payment page `/_admin/billplz/pay/{id}`.

### Payment methods

- [ ] Implement `GET /api/v3/collections/{id}/payment_methods` returning the available method codes.
- [ ] Support at minimum: `fpx`, `mpgs` (card), `boost`, `touchngo`.
- [ ] Document the full list of Billplz v3 method codes in `DECISIONS.md`.

### Local payment page

- [ ] Create `internal/ui/billplz-pay.html`.
- [ ] `GET /_admin/billplz/pay/{id}` renders payment method selector + amount + pay/cancel buttons.
- [ ] `POST /_admin/billplz/pay/{id}` processes outcome:
  - Pay → `state=paid`, `paid=true`, `paid_amount=amount`, record ledger transaction as `PAID`, dispatch callback.
  - Cancel → keep `state=due` (or mark `deleted` if explicitly deleted), dispatch callback with `paid=false`.
  - Include CSRF token in forms.

### Redirect and callback

- [ ] Implement browser redirect handler at `GET /billplz/redirect`:
  - Build a signed query string with `billplz[id]`, `billplz[paid]`, `billplz[state]`, and `x_signature`.
  - Redirect to the bill's `redirect_url` with that query string.
- [ ] Implement server-side callback dispatch to the bill's `callback_url`:
  - POST a form-urlencoded Bill object with `x_signature` as a form field.
  - Use the provider's webhook dispatcher and `PayloadBuilder`.
- [ ] Provider `PayloadBuilder` returns the flat form-urlencoded Bill object bytes.
- [ ] Provider `PayloadHeaders` returns `Content-Type: application/x-www-form-urlencoded`; no `X-Signature` header.

### `x_signature` algorithm

- [ ] Create `internal/billplz/signature.go`.
- [ ] Implement HMAC-SHA256:
  - Sort keys case-insensitively ascending.
  - For each key, concatenate `key + value`.
  - Join pairs with `|`.
  - Sign with `x_signature_key`.
- [ ] Verify signature on any incoming callback payload and on redirect query params.

### Tests

- [ ] Create `internal/billplz/*_test.go` covering:
  - Collection create/retrieve
  - Bill create/retrieve/delete with correct response fields
  - Payment methods list
  - Payment page render and outcome
  - Redirect query string includes correct `x_signature`
  - Callback form payload signature verification
  - Invalid API key / missing collection errors
- [ ] Update `scripts/smoke-test.sh` with a Billplz happy path.

### OpenAPI

- [ ] Update `docs/openapi.yaml` and `internal/server/openapi.yaml` with Billplz endpoints.

## Files to Create/Change

- `internal/billplz/provider.go`
- `internal/billplz/collection.go`
- `internal/billplz/collection_types.go`
- `internal/billplz/bill.go`
- `internal/billplz/bill_types.go`
- `internal/billplz/payment_methods.go`
- `internal/billplz/signature.go`
- `internal/billplz/webhook.go`
- `internal/billplz/*_test.go`
- `internal/ui/billplz-pay.html`
- `internal/ui/embed.go`
- `internal/server/router.go`
- `docs/openapi.yaml`
- `internal/server/openapi.yaml`
- `scripts/smoke-test.sh`

## Response / Webhook Shape

Return:
1. Collection create/retrieve JSON shapes
2. Bill create/retrieve/delete JSON shapes (no `currency` field)
3. Payment methods list shape
4. Local payment page HTML shape
5. Redirect URL query params including `x_signature`
6. Callback form-urlencoded payload including `x_signature`

## Test Notes

- `go test ./internal/billplz/...`
- `./scripts/smoke-test.sh`
- Verify existing provider tests still pass.
