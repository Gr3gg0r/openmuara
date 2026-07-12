---
id: hosted-testing
title: Hosted Testing Guide
---

# Hosted Testing Guide

> **Scope:** Run OpenMuara as a shared, Mailpit-style service for QA testers and CI pipelines.

---

## Overview

OpenMuara is local-first by default, but it can be hosted on a team network or exposed through a reverse proxy so testers can integrate with it like Mailpit. This guide shows the recommended network layout, config, and access controls for that scenario.

Key principles:

- **Provider emulation endpoints stay public.** Routes such as `/fawry/charge`, `/v1/checkout/sessions`, and the payment/escape pages under `/_admin/*` behave like real provider endpoints. Testers and redirected browsers must reach them without admin credentials.
- **Admin dashboard is locked.** `/_admin/*` config, replay, and inspection endpoints require admin credentials.
- **Read-only viewer accounts.** Non-developer testers can inspect the ledger and webhook log with a `viewer` account, but they cannot mutate configuration or replay webhooks.
- **Dual ports make this easy.** `server.port` carries the provider API and payment pages; `server.admin_port` carries the dashboard and admin JSON APIs.

---

## Which port does a payment redirect use?

Payment redirects and provider callbacks always use the **provider/API port** (`server.port`), never the admin port.

When a test app calls `/fawry/charge`, `/v1/checkout/sessions`, or `/senangpay/charge`, the emulated response contains a payment URL on the provider port. A browser redirected to that URL completes the payment on the provider port. Because provider payment and escape pages are part of the public emulation flow, **the browser is not asked for an admin password**.

Admin authentication applies only to the dashboard and admin JSON APIs (`/_admin/*` excluding provider simulation pages).

---

## Recommended DNS layout

For a hosted deployment at `openmuara.example.com`:

| Hostname | Points to | Purpose | Who accesses | Auth |
|---|---|---|---|---|
| `api.openmuara.example.com` | `server.port` (e.g. `0.0.0.0:9000`) | Provider API and payment/escape pages | Test apps, redirected browsers, CI | None |
| `admin.openmuara.example.com` | `server.admin_port` (e.g. `0.0.0.0:9001`) | Dashboard, ledger, config, replay | Developers and QA leads | Admin or viewer |
| Your test app webhook URL | Your test app | Receives outgoing webhooks from OpenMuara | OpenMuara dispatcher | Provider signature |

If you only need one hostname, use sub-path routing through a reverse proxy and route `/_admin/*` to `server.admin_port` while everything else goes to `server.port`. The dual-hostname layout is cleaner and easier to lock down.

---

## Example configuration

```yaml
server:
  host: 0.0.0.0
  port: 9000
  admin_port: 9001
  # Set these to the public URLs testers and browsers will use.
  public_base_url: https://api.openmuara.example.com
  admin_public_base_url: https://admin.openmuara.example.com
  tls_cert: ""
  tls_key: ""
  cors:
    allowed_origins:
      - "https://app-under-test.example.com"
    allowed_methods:
      - GET
      - POST
      - PUT
      - PATCH
      - DELETE
      - OPTIONS
    allowed_headers:
      - Content-Type
      - Authorization
      - X-CSRF-Token
      - X-Request-ID
    allow_credentials: false
  csrf:
    enabled: true
  pprof: false

# Admin account for developers. Required when hardened is true.
admin:
  enabled: true
  username: muara-admin
  password_hash: "$2a$10$..."  # muara security hash-password

# Read-only account for testers. Viewers can inspect the ledger and webhooks
# but cannot change providers, webhook targets, or replay webhooks.
viewer:
  enabled: true
  username: muara-viewer
  password_hash: "$2a$10$..."  # muara security hash-password

rate_limit:
  enabled: true
  requests_per_minute: 200

hardened: true

log:
  level: info

persistence:
  type: sqlite
  path: .muara/data/ledger.db

providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
      version: v1

webhook:
  url: "https://app-under-test.example.com/webhooks/openmuara"
  max_retries: 3
```

---

## Cloudflare Tunnel example

With [`cloudflared`](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/):

```yaml
# ~/.cloudflared/config.yml for the OpenMuara tunnel
tunnel: <TUNNEL_ID>
credentials-file: /path/to/credentials.json

ingress:
  - hostname: admin.openmuara.example.com
    service: http://localhost:9001
    originRequest:
      noTLSVerify: true
  - hostname: api.openmuara.example.com
    service: http://localhost:9000
    originRequest:
      noTLSVerify: true
  - service: http_status:404
```

Because `cloudflared` terminates TLS and forwards HTTP to OpenMuara, set `server.public_base_url` and `server.admin_public_base_url` to the public `https://` hostnames. OpenMuara will use those URLs when building payment links and redirect URLs.

If you want to restrict the admin hostname to your team, add an Access policy in Cloudflare Zero Trust for `admin.openmuara.example.com`. OpenMuara's own admin/viewer auth is still recommended as a second layer.

---

## nginx reverse-proxy example

```nginx
server {
    listen 443 ssl http2;
    server_name api.openmuara.example.com;

    location / {
        proxy_pass http://127.0.0.1:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 443 ssl http2;
    server_name admin.openmuara.example.com;

    # Restrict to office or VPN IPs, or add an auth layer here.
    allow 10.0.0.0/8;
    deny all;

    location / {
        proxy_pass http://127.0.0.1:9001;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## Caddy reverse-proxy example

```caddy
api.openmuara.example.com {
    reverse_proxy localhost:9000
}

admin.openmuara.example.com {
    @denied not remote_ip 10.0.0.0/8
    respond @denied 403

    reverse_proxy localhost:9001
}
```

---

## FAQ

### Do testers need the admin password to complete a payment?

No. Provider payment and escape pages are part of the public emulation flow. Only the dashboard and admin JSON APIs require admin or viewer credentials.

### What can a viewer account do?

A viewer can:

- View the ledger and transaction list.
- Inspect individual transactions and webhook attempts.
- View provider health and onboarding status.

A viewer cannot:

- Enable or disable providers.
- Change webhook targets or provider configuration.
- Replay webhooks or run scenario simulations.
- Access `/_admin/config/*` write endpoints.

### Should I expose OpenMuara directly to the public internet?

No. Run it behind a reverse proxy or tunnel, bind to `0.0.0.0` only when necessary, and use `hardened: true` with admin credentials. Treat the admin dashboard like any other internal tool.

### What if I do not set `public_base_url`?

Payment links and webhook payloads fall back to `http://server.host:server.port`. Behind a reverse proxy or tunnel this will produce broken `http://127.0.0.1:9000/...` links. Always set `public_base_url` (and `admin_public_base_url` when using dual-port) for hosted deployments.

### Does Cloudflare Tunnel break webhooks?

No, as long as `webhook.url` points to a URL that OpenMuara can reach. If your test app is also behind a tunnel, use its public tunnel URL as the webhook target.

---

## Hardening checklist

- [ ] Bind to `0.0.0.0` only when the host must accept external connections.
- [ ] Set `server.public_base_url` to the public provider API hostname.
- [ ] Set `server.admin_public_base_url` to the public admin hostname when using `server.admin_port`.
- [ ] Enable `hardened: true`.
- [ ] Configure `admin` credentials for developers.
- [ ] Configure a separate `viewer` account if non-developers need dashboard access.
- [ ] Restrict `admin` port/network access with a reverse proxy, firewall, or tunnel policy.
- [ ] Run `muara security audit` after changing config.
- [ ] Keep provider secrets and admin password hashes out of version control; use environment variables.
