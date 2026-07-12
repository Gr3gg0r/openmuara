---
id: ipay88
title: iPay88 Provider
---

# iPay88 Provider

Emulates the iPay88 Malaysia ePayment gateway for local development.

## Configuration

```yaml
providers:
  ipay88:
    enabled: true
    config:
      merchant_code: muara-ipay88-merchant
      merchant_key: muara-ipay88-key
```

## First request

Compute the request signature as `SHA256(MerchantKey + MerchantCode + RefNo + Amount + Currency)`,
where amount has all non-digits stripped after formatting to two decimals:

```bash
SIGNATURE=$(printf '%s' "muara-ipay88-keymuara-ipay88-merchantref-11000MYR" | sha256sum | awk '{print $1}')
curl -X POST http://127.0.0.1:9000/ePayment/entry.asp \
  -d "MerchantCode=muara-ipay88-merchant" \
  -d "RefNo=ref-1" \
  -d "Amount=10.00" \
  -d "Currency=MYR" \
  -d "ProdDesc=Test product" \
  -d "UserName=Test User" \
  -d "UserEmail=test@example.com" \
  -d "UserContact=0123456789" \
  -d "Signature=$SIGNATURE" \
  -d "SignatureType=SHA256" \
  -d "ResponseURL=http://example.com/response" \
  -d "BackendURL=http://example.com/backend"
```

The endpoint redirects to `/_admin/ipay88/pay/{refNo}` for payment simulation.

> **Note:** iPay88 requires `ResponseURL` and `BackendURL` to be publicly routable
> (not loopback or private IPs). For local webhook testing, point the backend URL
> at a public tunnel or inspect the admin simulation page directly.

## Routes

| Method | Route | Purpose |
|---|---|---|
| POST | `/ePayment/entry.asp` | Payment entry point |
| POST | `/ePayment/enquiry.asp` | Transaction requery |
| GET | `/_admin/ipay88/pay/{refNo}` | Render payment simulation page |
| POST | `/_admin/ipay88/pay/{refNo}` | Submit payment outcome |
| POST | `/ipay88/response` | Response endpoint (returns to merchant) |
| POST | `/ipay88/backend` | Backend post endpoint (server-to-server) |

## Signature algorithm

### Payment request

```text
SHA256(MerchantKey + MerchantCode + RefNo + Amount + Currency)
```

Amount is formatted to two decimals, then all non-digit characters are removed.

Example:

```bash
printf '%s' "muara-ipay88-keymuara-ipay88-merchantref-11000MYR" | sha256sum
# 9c2b...  (hex)
```

### Response / backend post

```text
SHA256(MerchantKey + MerchantCode + PaymentId + RefNo + Amount + Currency + Status)
```

## Simulation / escape routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/_admin/ipay88/pay/{refNo}` | Render iPay88 payment simulation page |
| POST | `/_admin/ipay88/pay/{refNo}` | Submit success or failure for the transaction |

Example:

```bash
curl -X POST http://127.0.0.1:9000/_admin/ipay88/pay/ref-1 \
  -d "outcome=success"
```

## Webhooks

Outgoing backend post payload (form-encoded):

```text
MerchantCode=muara-ipay88-merchant
PaymentId=1
RefNo=ref-1
Amount=10.00
Currency=MYR
Status=1
Signature=...
```

## Common errors

| HTTP status | Error code | Cause | Fix |
|---|---|---|---|
| 400 | `missing_field` | Missing required field | Check form fields |
| 400 | `invalid_signature` | `SignatureType` is not `SHA256` or signature mismatch | Use SHA256 and correct canonical string |
| 400 | `invalid_state` | `ResponseURL` or `BackendURL` is not a valid URL | Use `http://127.0.0.1/...` |
| 404 | `not_found` | Transaction not found on requery | Create via entry.asp first |
| 500 | `internal` | Store error | Check server logs |

## See also

- `runbooks/local-development.md` — running OpenMuara locally.
