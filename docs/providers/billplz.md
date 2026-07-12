---
id: billplz
title: Billplz Provider
---

# Billplz Provider

Emulates the Billplz v3 API for local development.

## Configuration

```yaml
providers:
  billplz:
    enabled: true
    config:
      api_key: muara-billplz-api-key
      x_signature_key: muara-billplz-x-signature
      collection_id: ""
```

## First request

Create a collection, then create a bill:

```bash
COLLECTION=$(curl -s -X POST http://127.0.0.1:9000/api/v3/collections \
  -u muara-billplz-api-key: \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Collection"}' | jq -r '.id')

curl -X POST http://127.0.0.1:9000/api/v3/bills \
  -u muara-billplz-api-key: \
  -H "Content-Type: application/json" \
  -d "{
    \"collection_id\": \"$COLLECTION\",
    \"description\": \"Test bill\",
    \"amount\": 1000,
    \"name\": \"Test User\",
    \"email\": \"test@example.com\",
    \"callback_url\": \"http://127.0.0.1/callback\"
  }"
```

Expected response (abbreviated):

```json
{
  "id": "bill_...",
  "collection_id": "col_...",
  "description": "Test bill",
  "amount": 1000,
  "state": "due",
  "paid": false,
  "url": "http://127.0.0.1:9000/_admin/billplz/pay/bill_..."
}
```

## Routes

| Method | Route | Purpose |
|---|---|---|
| POST | `/api/v3/collections` | Create a collection |
| GET | `/api/v3/collections/{id}` | Retrieve a collection |
| GET | `/api/v3/collections/{id}/payment_methods` | List payment methods |
| POST | `/api/v3/bills` | Create a bill |
| GET | `/api/v3/bills/{id}` | Retrieve a bill |
| DELETE | `/api/v3/bills/{id}` | Delete a bill |
| GET | `/_admin/billplz/pay/{id}` | Render payment simulation page |
| POST | `/_admin/billplz/pay/{id}` | Submit payment outcome |
| GET | `/billplz/redirect` | Redirect callback with signature |
| POST | `/billplz/webhook` | Incoming webhook receiver |

## Signature algorithm

Billplz uses HMAC-SHA256 over sorted key-value pairs:

1. Sort keys case-insensitively ascending.
2. For each key, append `key` + `value`.
3. Join pairs with `|`.
4. Sign with `x_signature_key`.

```go
msg := "amount1000|callback_urlhttp://127.0.0.1/callback|..."
sig := HMAC-SHA256(msg, xSignatureKey)
```

### Incoming webhook / redirect

OpenMuara validates the `x_signature` in form-encoded callbacks and webhooks by recomputing the
signature over all fields except `x_signature`.

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/_admin/billplz/pay/{id}` | Render Billplz payment simulation page |
| POST | `/_admin/billplz/pay/{id}` | Submit success or failure for the bill |

Example:

```bash
curl -X POST http://127.0.0.1:9000/_admin/billplz/pay/bill_xxx \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "outcome=success"
```

## Webhooks

Outgoing webhook payload (form-encoded):

```text
id=bill_xxx&collection_id=col_xxx&paid=true&state=paid&amount=1000&description=Test+bill&...
&x_signature=...
```

Headers:

```text
Content-Type: application/x-www-form-urlencoded
```

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `missing_field` | Missing `title`, `collection_id`, `email`, `name`, etc. | Check request body |
| 401 | `unauthorized` | Missing or invalid basic auth API key | Use `-u muara-billplz-api-key:` |
| 401 | `unauthorized` | Invalid `x_signature` | Recompute signature with correct `x_signature_key` |
| 404 | `not_found` | Collection or bill not found | Create it first |
| 500 | `internal` | Store error | Check server logs |

## See also

- `runbooks/local-development.md` — running OpenMuara locally.
