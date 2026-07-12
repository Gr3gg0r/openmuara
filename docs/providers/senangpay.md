---
id: senangpay
title: SenangPay Provider
---

# SenangPay Provider

Emulates the SenangPay payment flow for local development.

## Configuration

```yaml
providers:
  senangpay:
    enabled: true
    config:
      secret_key: muara-senangpay-secret
```

## First request

Compute the charge hash as `MD5(secret_key + detail + amount + order_id)`, then create the charge:

```bash
HASH=$(printf '%s' "muara-senangpay-secretTest payment10.00order-1" | md5sum | awk '{print $1}')
curl -X POST http://127.0.0.1:9000/senangpay/charge \
  -H "Content-Type: application/json" \
  -d "{
    \"detail\": \"Test payment\",
    \"amount\": 10.00,
    \"order_id\": \"order-1\",
    \"name\": \"Test User\",
    \"email\": \"test@example.com\",
    \"phone\": \"0123456789\",
    \"hash\": \"$HASH\"
  }"
```

Expected response:

```json
{
  "order_id": "order-1",
  "payment_url": "http://127.0.0.1:9000/_admin/senangpay-escape?order_id=order-1",
  "status": "ok",
  "reference": "order-1"
}
```

## Routes

| Method | Route | Purpose |
|---|---|---|
| POST | `/senangpay/charge` | Create a charge |
| GET | `/senangpay/callback` | Payment callback from gateway |
| GET | `/senangpay/query` | Query payment status by order ID |
| POST | `/senangpay/webhook` | Webhook notification |

## Signature algorithm

### Charge request

```text
hash = MD5(secret_key + detail + amount + order_id)
```

Amount is formatted with exactly two decimal places.

Example:

```bash
printf '%s' "muara-senangpay-secretTest payment10.00order-1" | md5sum
# 7a3b9c...  (hex)
```

### Status query

```text
hash = MD5(secret_key + order_id)
```

Example:

```bash
printf '%s' "muara-senangpay-secretorder-1" | md5sum | awk '{print $1}'
```

### Callback / webhook

SenangPay callbacks and webhooks use query parameters. OpenMuara does not sign these; it validates
the `order_id` against the ledger and applies the `status_id`.

## Query payment status

```bash
HASH=$(printf '%s' "muara-senangpay-secretorder-1" | md5sum | awk '{print $1}')
curl "http://127.0.0.1:9000/senangpay/query?order_id=order-1&hash=$HASH"
```

## Simulation / escape routes

SenangPay status is updated through the callback or webhook endpoints:

| Method | Route | Purpose |
|---|---|---|
| GET | `/senangpay/callback?order_id=order-1&status_id=1&transaction_id=txn-1&msg=Paid` | Mark order as paid |
| GET | `/senangpay/callback?order_id=order-1&status_id=0&transaction_id=txn-1&msg=Failed` | Mark order as unpaid |
| POST | `/senangpay/webhook?order_id=order-1&status_id=1` | Webhook notification |

Example success callback:

```bash
curl "http://127.0.0.1:9000/senangpay/callback?order_id=order-1&status_id=1&transaction_id=txn-1&msg=Paid"
```

## Webhooks

Outgoing webhook payload (dispatched when status changes):

```json
{
  "provider": "senangpay",
  "reference": "order-1",
  "status": "PAID"
}
```

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `invalid_signature` | Charge `hash` does not match | Recompute `MD5(secret + detail + amount + order_id)` |
| 400 | `missing_field` | Missing `detail`, `amount`, `order_id`, or `hash` | Check request body |
| 404 | `not_found` | Order not found in ledger | Create the charge first |
| 500 | `internal` | Ledger error | Check server logs |

## See also

- `tasks/senangpay-signature.md` — detailed signature specification.
- `runbooks/local-development.md` — running OpenMuara locally.
