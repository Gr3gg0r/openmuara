---
id: webhooks
title: Webhooks
---

# Webhooks

OpenMuara can dispatch provider-style outgoing webhooks locally so you can test your webhook handlers without tunnels or real provider sandboxes.

## Configure

Add a `webhook` section to `.muara/config.yml`:

```yaml
webhook:
  url: "http://localhost:3000/webhook/fawry"
  max_retries: 3
```

If `url` is empty, muara will warn on startup but still run. No webhooks will be dispatched.

## Trigger a webhook

Start muara:

```bash
muara start
```

Create a charge request and open the escape page. Click **Simulate Paid**. OpenMuara will:

1. Redirect the browser to your `returnUrl`.
2. POST a Fawry V2 webhook payload to your configured `webhook.url`.

## Test without your own app

Use the built-in test receiver to capture webhooks:

```yaml
webhook:
  url: "http://127.0.0.1:9000/_admin/webhook-receiver"
```

Every dispatched webhook will be accepted and logged.

## Inspect webhooks

List recent attempts:

```bash
muara webhook list
```

Inspect one:

```bash
muara webhook inspect ref-1
```

Replay one:

```bash
muara webhook replay ref-1
```

Replay only works for attempts in the current process; the store is in-memory.

## Signature verification

The Fawry V2 `messageSignature` in muara is a **documented approximation**, not a verified Fawry formula. Real Fawry backends may use a different signing scheme.

OpenMuara computes the signature as:

```
HMAC-SHA256(secret, canonical_json(payload_without_messageSignature))
```

Your test app can verify it with the same secret:

```go
mac := hmac.New(sha256.New, []byte(secret))
mac.Write(canonicalJSON)
sig := hex.EncodeToString(mac.Sum(nil))
```

Use this only for local testing. Do not rely on muara's signature in production.
