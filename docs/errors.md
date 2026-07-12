---
id: errors
title: OpenMuara Error-Code Taxonomy
---

# OpenMuara Error-Code Taxonomy

> **Updated:** 2026-07-06

OpenMuara uses stable error codes so users, contributors, and bug hunts can classify failures by symptom. Each code has the form `E<area><sequence>`.

## Code Groups

| Group | Range | Area |
|-------|-------|------|
| Generic / internal | E1000–E1099 | Unexpected or cross-cutting failures |
| Configuration | E2000–E2099 | Config loading, provider enablement, validation |
| Provider emulation | E3000–E3099 | Charge, callback, escape, version errors |
| Webhook dispatch | E4000–E4099 | Webhook build, delivery, replay errors |
| Transaction / ledger | E5000–E5099 | Store, state machine, idempotency errors |
| Signature / security | E6000–E6099 | HMAC/MD5 verification, missing signatures |

## Current Codes

| Code | Meaning |
|------|---------|
| `E1000` | Internal error |
| `E1001` | Unknown provider |
| `E1002` | Invalid request |
| `E2000` | Missing required config value |
| `E2001` | Invalid config value |
| `E2002` | Provider disabled |
| `E3000` | Provider charge failed |
| `E3001` | Provider callback failed |
| `E3002` | Provider escape/3-D Secure failed |
| `E3003` | Unsupported provider version |
| `E4000` | Webhook URL not configured |
| `E4001` | Failed to build webhook payload |
| `E4002` | Webhook delivery failed after retries |
| `E4003` | Webhook replay target not found |
| `E5000` | Transaction not found |
| `E5001` | Invalid transaction state transition |
| `E5002` | Duplicate transaction / idempotency conflict |
| `E6000` | Signature mismatch |
| `E6001` | Missing signature |

## Using Error Codes

Use `internal/errcode` to wrap errors:

```go
import "github.com/openmuara/openmuara/internal/errcode"

return errcode.Wrap(errcode.EWebhookDeliveryFailed, "webhook delivery exhausted retries", err)
```

When adding a new code, update this file and keep the grouping consistent.
