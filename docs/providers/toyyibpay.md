---
id: toyyibpay
title: ToyyibPay Provider
---

# ToyyibPay Provider

Emulates the [ToyyibPay API](https://toyyibpay.com/apireference/) for local development.

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
  -d "catname=Test Category" \
  -d "catdescription=Test" | jq -r '.data.categoryCode')

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
    "billCode": "...",
    "billPaymentLink": "http://127.0.0.1:9000/_admin/toyyibpay/pay/..."
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

OpenMuara validates incoming `/toyyibpay/webhook` payloads using this formula.
The `/toyyibpay/return` URL carries no hash — matching upstream, it is a plain
browser redirect with `status_id`, `billcode`, and `order_id`.

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/_admin/toyyibpay/pay/{billCode}` | Render ToyyibPay payment simulation page |
| POST | `/_admin/toyyibpay/pay/{billCode}` | Submit success or failure for the bill |

Example:

```bash
curl -X POST http://127.0.0.1:9000/_admin/toyyibpay/pay/xxx \
  -d "status=1"
```

(`status=1` = success, `status=3` = failure. The dashboard pay page submits the
CSRF token for you; raw curl needs the `openmuara_csrf` cookie + `csrf_token`
field unless CSRF is disabled.)

## Webhooks

Outgoing webhook payload (form-encoded), fields per the official callback spec:

```text
refno=MUARA-...&status=1&reason=Payment+success&billcode=...&order_id=order-1&amount=1000&transaction_time=2026-07-12+10:00:00&hash=...
```

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `missing_field` / `invalid_request` | Missing `billName`, `billAmount`, URLs, etc. | Check form fields |
| 401 | `unauthorized` | Invalid `userSecretKey` | Use `muara-toyyibpay-secret` |
| 401 | `unauthorized` | Invalid callback `hash` | Recompute MD5 with correct secret |
| 404 | `not_found` | Bill or category not found | Create it first |
| 500 | `internal` | Store error | Check server logs |

## Divergences from the official API

Request parameters, routes, status codes, and the callback hash follow the
[official reference](https://toyyibpay.com/apireference/). Response envelopes
do not:

| Official ToyyibPay | OpenMuara |
|---|---|
| `createBill` returns `[{"BillCode":"..."}]` | returns `{"status":"1","msg":"success","bill":{...}}` |
| `createCategory` returns `[{"CategoryCode":"..."}]` | returns `{"status":"1","msg":"success","data":{...}}` |
| `getCategoryDetails` returns an array | returns a single object |
| `getBillTransactions` returns `billpaymentInvoiceNo`, `billPaymentDate`, amount as decimal string (`"10.00"`) | returns `billpaymentRefNo`, `billpaymentTime`, amount in cents |
| `inactiveBill` takes `secretKey`, returns `{"status":"success","result":"..."}` | takes `userSecretKey`, returns `{"status":"1","msg":"success"}` |

`createCategory` accepts the official `catname`/`catdescription` parameters;
the legacy camelCase aliases (`categoryName`/`categoryDescription`) are still
accepted.

## Limitations

These parts of the official API are not emulated:

- Enterprise partner APIs (`createAccount`, `getBank`, `getUserStatus`, `getSettlementSummary`)
- DuitNow QR (`checkDuitNowQRStatus`, `enableDuitNowQR` bill flags)
- Dynamic bills (`billPriceSetting=0`) — `billAmount` must be a positive integer in cents
- Split payments (`billSplitPayment` / `billSplitPaymentArgs`)
- Bill expiry enforcement — `billExpiryDate` / `billExpiryDays` are stored but not enforced
- Pending simulation — the pay page offers success (`status=1`) and failure (`status=3`) only

## See also

- `runbooks/local-development.md` — running OpenMuara locally.
