---
id: default
title: Default Provider
---

# Default Provider

A minimal provider for custom experiments and quick smoke tests. It requires no
provider-specific configuration and is always available as a fallback.

## Configuration

```yaml
providers:
  default:
    enabled: true
    config: {}
```

## Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/default` | Hello endpoint |
| POST | `/default/charge` | Create a charge |
| POST | `/default/webhook` | Webhook endpoint |

## First request

```bash
curl -X POST http://127.0.0.1:9000/default/charge
```

Expected response:

```json
{
  "provider": "default",
  "transaction_id": "...",
  "status": "success"
}
```

## Webhooks

Incoming webhooks are acknowledged with:

```json
{
  "status": "acknowledged"
}
```

Outgoing webhook payloads have the shape:

```json
{
  "provider": "default",
  "reference": "ref-1",
  "status": "PAID"
}
```

## Common errors

| HTTP status | Cause | Fix |
|---|---|---|
| 405 | Method not allowed | Use the supported method for the route |
| 500 | Internal error | Check server logs |
