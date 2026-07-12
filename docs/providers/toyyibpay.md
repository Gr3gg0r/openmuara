---
id: toyyibpay
title: ToyyibPay Provider
---

# ToyyibPay Provider

Emulates the ToyyibPay API for local development.

## Configuration

```yaml
providers:
  toyyibpay:
    enabled: true
    config:
      user_secret_key: muara-toyyibpay-secret
      category_code: ""
```

## First request

Create a category, then create a bill:

```bash
CATEGORY=$(curl -s -X POST http://127.0.0.1:9000/index.php/api/createCategory \
  -d "userSecretKey=muara-toyyibpay-secret" \
  -d "categoryName=Test Category" \
  -d "categoryDescription=Test" | jq -r '.data.categoryCode')

curl -X POST http://127.0.0.1:9000/index.php/api/createBill \
  -d "userSecretKey=muara-toyyibpay-secret" \
  -d "categoryCode=$CATEGORY" \
  -d "billName=Test Bill" \
  -d "billDescription=Test" \
  -d "billAmount=1000" \
  -d "billReturnUrl=http://127.0.0.1/return" \
  -d "billCallbackUrl=http://127.0.0.1/callback" \
  -d "billTo=Test User" \
  -d "billEmail=test@example.com" \
  -d "billPhone=0123456789"
```

Expected response (abbreviated):

```json
{
  "status": "1",
  "msg": "success",
  "bill": {
    "BillCode": "...",
    "BillPaymentLink": "http://127.0.0.1:9000/_admin/toyyibpay/pay/..."
  }
}
```

## Routes

| Method | Route | Purpose |
|---|---|---|
| POST | `/index.php/api/createCategory` | Create a category |
| POST | `/index.php/api/getCategoryDetails` | Get category details |
| POST | `/index.php/api/createBill` | Create a bill |
| POST | `/index.php/api/getBillTransactions` | Get bill transactions |
| POST | `/index.php/api/inactiveBill` | Deactivate a bill |
| GET | `/_admin/toyyibpay/pay/{billCode}` | Render payment simulation page |
| POST | `/_admin/toyyibpay/pay/{billCode}` | Submit payment outcome |
| GET | `/toyyibpay/return` | Return callback |
| POST | `/toyyibpay/webhook` | Incoming webhook receiver |

## Signature algorithm

ToyyibPay callbacks use an MD5 hash:

```text
hash = MD5(userSecretKey + status + order_id + refno + "ok")
```

Example:

```bash
printf '%s' "muara-toyyibpay-secret1order-1ref-1ok" | md5sum
# 3a7f...  (hex)
```

OpenMuara validates incoming `/toyyibpay/webhook` and `/toyyibpay/return` payloads using this formula.

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/_admin/toyyibpay/pay/{billCode}` | Render ToyyibPay payment simulation page |
| POST | `/_admin/toyyibpay/pay/{billCode}` | Submit success or failure for the bill |

Example:

```bash
curl -X POST http://127.0.0.1:9000/_admin/toyyibpay/pay/xxx \
  -d "outcome=success"
```

## Webhooks

Outgoing webhook payload (form-encoded):

```text
status=1&order_id=order-1&refno=ref-1&amount=1000&hash=...
```

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `missing_field` / `invalid_request` | Missing `billName`, `billAmount`, URLs, etc. | Check form fields |
| 401 | `unauthorized` | Invalid `userSecretKey` | Use `muara-toyyibpay-secret` |
| 401 | `unauthorized` | Invalid callback `hash` | Recompute MD5 with correct secret |
| 404 | `not_found` | Bill or category not found | Create it first |
| 500 | `internal` | Store error | Check server logs |

## See also

- `runbooks/local-development.md` — running OpenMuara locally.
