---
id: fawry
title: Fawry Provider
---

# Fawry Provider

Emulates the Fawry payment gateway for local development.

## Supported versions

- `v1` (default): legacy charge and webhook payload.
- `v2`: server notification format.

## Configuration

```yaml
providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
      version: v1
```

## First request

Create a charge. The signature is
`SHA256(merchantCode + merchantRefNum + customerProfileId + returnUrl +
sorted(itemId + quantity + price) + merchantSecurityKey)`:

```bash
SIGNATURE=$(printf '%s' "muara-merchant-coderef-1cust-1http://127.0.0.1/callbackprod-1110.00muara-fawry-secret" | sha256sum | awk '{print $1}')
curl -X POST http://127.0.0.1:9000/fawry/charge \
  -H "Content-Type: application/json" \
  -d "{
    \"merchantCode\": \"muara-merchant-code\",
    \"merchantRefNum\": \"ref-1\",
    \"customerProfileId\": \"cust-1\",
    \"returnUrl\": \"http://127.0.0.1/callback\",
    \"chargeItems\": [{\"itemId\": \"prod-1\", \"price\": 10.0, \"quantity\": 1}],
    \"signature\": \"$SIGNATURE\"
  }"
```

Expected response:

```json
{
  "status": "ok",
  "reference": "ref-1"
}
```

## Query payment status

Generate the signature as `SHA256(merchantCode + merchantRefNum + merchantSecurityKey)`:

```bash
SIGNATURE=$(printf '%s' "muara-merchant-coderef-1muara-fawry-secret" | sha256sum | awk '{print $1}')
curl "http://127.0.0.1:9000/fawry/payment-status?merchantCode=muara-merchant-code&merchantRefNum=ref-1&signature=$SIGNATURE"
```

## Signature algorithms

### Payment-status query

```text
SHA256(merchantCode + merchantRefNum + merchantSecurityKey)
```

Example:

```bash
printf '%s' "muara-merchant-coderef-1muara-fawry-secret" | sha256sum
# 8f4e2f...  (hex)
```

### Incoming webhook (`/fawry/webhook`)

The webhook receiver validates the query parameter `token` against `webhook_secret` and the JSON body's `messageSignature`.

**v1:**

```text
HMAC-SHA256(merchantRefNumber + orderStatus, webhookSecret)
```

**v2:**

The v2 payload is signed with `webhook.NewHMACSigner(webhookSecret).Sign(payload)`, which serializes
the payload fields deterministically before signing.

### Outgoing webhook

When the escape page dispatches a notification, OpenMuara signs the payload using the configured
`webhook_secret` with the same algorithm as the incoming receiver.

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/_admin/fawry-escape?ref=ref-1&returnUrl=http://127.0.0.1/callback` | Render payment simulation page |
| POST | `/_admin/fawry-escape` | Submit `status=PAID` or `status=CANCELED`; updates ledger and dispatches webhook |

Example:

```bash
curl -X POST http://127.0.0.1:9000/_admin/fawry-escape \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "ref=ref-1" \
  -d "returnUrl=http://127.0.0.1/callback" \
  -d "status=PAID"
```

## Webhooks

### Incoming v1 payload

```json
{
  "merchantRefNumber": "ref-1",
  "orderStatus": "PAID",
  "messageSignature": "8f4e2f..."
}
```

Incoming requests must include `?token=muara-webhook-secret`.

### Outgoing v1 payload

When dispatched from the escape page, OpenMuara POSTs:

```json
{
  "merchantRefNumber": "ref-1",
  "orderStatus": "PAID",
  "messageSignature": "..."
}
```

to the configured `webhook.url`.

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `invalid_request` | Missing `merchantCode`, `merchantRefNum`, or other required field | Check request body |
| 401 | `unauthorized` | Invalid `token` or `messageSignature` | Verify `webhook_secret` |
| 404 | `not_found` | Transaction not found in ledger | Create the charge first |
| 409 | `invalid_state` | Invalid status transition (e.g., `PAID` → `PAID`) | Use a fresh reference |
| 500 | `internal` | Ledger or dispatcher unavailable | Check server logs |

## See also

- `docs/mkp-billing-requirements.md` — MKP v2 Fawry integration requirements.
- `runbooks/local-development.md` — running OpenMuara locally.
