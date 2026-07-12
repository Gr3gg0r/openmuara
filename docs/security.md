---
id: security
title: OpenMuara Security Hardening Guide
---

> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This doc is subordinate to it.**

# OpenMuara Security Hardening Guide

> **Status:** ✅ Completed (Security Hardening initiative)  
> **Scope:** Defense-in-depth controls for the admin surface and deployment footprint. Provider emulation endpoints remain drop-in replacements.

---

## Core Philosophy

OpenMuara is a **drop-in payment emulator**. A user integrates it by changing only the base URL in their production provider client. Security controls must not violate this:

- **Provider emulation endpoints** (`/fawry/charge`, `/v1/checkout/sessions`, `/senangpay/charge`, etc.) behave exactly like the real providers. They stay public and signature-verified where the real APIs are public.
- **Security controls apply only to the admin/internal surface**: `/_admin`, admin JSON APIs, replay, simulation, and configuration endpoints.
- **Default deployment is local and safe**: bind to `127.0.0.1`, no admin auth required. Hardening is opt-in for CI/CD or shared environments.

In short: **the payment-emulation API stays a drop-in replacement; the dashboard gets a lock.**

---

## Threat Model

### Assets to protect

- Ledger data (transactions, customer refs, amounts).
- Webhook payloads and replay capability.
- Provider configuration (API keys, secrets, signatures).
- Audit logs.

### Threats

1. **Unauthorized dashboard access** — attacker views transactions and replays webhooks via `/_admin`.
2. **Unauthorized simulation** — attacker triggers escape/pay/simulation actions to mutate ledger state.
3. **Network exposure** — dashboard bound to `0.0.0.0` in CI/CD without authentication.
4. **Credential sniffing** — admin secrets or provider config sent over plaintext HTTP.
5. **CSRF / XSS** — admin actions forged or JS injected via webhook payloads rendered in the UI.
6. **DoS** — excessive requests to admin replay endpoints or provider endpoints.

---

## Configuration Reference

### `server.host`

Bind address. Default: `127.0.0.1`. Use `0.0.0.0` only in Docker or shared networks, and pair with admin auth and TLS.

### `server.tls_cert` / `server.tls_key`

Paths to TLS certificate and key. When both are set, the server serves HTTPS. Optional.

### `admin.enabled`

Enable admin authentication for `/_admin/*` and admin JSON APIs. Default: `false`.

### `admin.username`

Admin username for HTTP Basic Auth.

### `admin.password_hash`

Bcrypt hash of the admin password. Generate with `muara security hash-password`.

### `admin.token`

Static bearer token for admin API access. Alternative or complement to basic auth.

### `viewer.enabled`

Enable read-only dashboard access for testers. Default: `false`.

### `viewer.username` / `viewer.password_hash` / `viewer.token`

Viewer credentials. Viewers can inspect the ledger, transactions, and webhook log but cannot change providers, webhook targets, or replay webhooks.

### `server.public_base_url`

External URL used for provider payment links when OpenMuara sits behind a reverse proxy or tunnel (e.g. Cloudflare Tunnel). Example: `https://api.muara.example.com`.

### `server.admin_public_base_url`

External URL used for the admin dashboard when it is served on a separate hostname or port. Example: `https://admin.muara.example.com`. Optional; falls back to `server.public_base_url` when empty.

### `rate_limit.enabled`

Enable in-memory per-IP rate limiting. Default: `false`.

### `rate_limit.requests_per_minute`

Request threshold per IP per minute. Default: `100`.

### `hardened`

Single-toggle secure preset. Enables admin auth, rate limiting, and strict security headers. Requires admin credentials. Default: `false`.

Environment variables override YAML values with the `MUARA_` prefix, e.g. `MUARA_ADMIN_PASSWORD_HASH`.

---

## CLI Helpers

| Command | Purpose |
|---|---|
| `muara security hash-password` | Generate a bcrypt hash for `admin.password_hash`. |
| `muara security gen-cert` | Generate a self-signed TLS cert/key pair for local testing. |
| `muara security audit` | Print security posture and warn on insecure settings. |

---

## Hardening Checklist

- [x] Bind to `127.0.0.1` unless you explicitly need network exposure.
- [x] Enable `admin.enabled` and set `admin.username` + `admin.password_hash` or `admin.token`.
- [x] Create a separate `viewer` account if non-developer testers need dashboard access.
- [x] Set `server.public_base_url` when running behind a reverse proxy or tunnel.
- [x] Enable TLS with real certificates in shared environments.
- [x] Enable `hardened: true` for a secure-by-default preset.
- [x] Run `muara security audit` after changing config.
- [x] Keep provider secrets out of version control; use environment variables.

### Implemented Controls

| Control | What it does | Default |
|---|---|---|
| Bind address | `server.host` defaults to `127.0.0.1`; opt-in to `0.0.0.0` | safe |
| Admin auth | HTTP Basic Auth or bearer token on `/_admin/*` only | off |
| Viewer role | Read-only dashboard access; cannot mutate config or replay webhooks | off |
| Public base URL | External URL for provider payment links behind reverse proxies | off |
| Password storage | bcrypt hashes; no plaintext | safe |
| TLS | HTTPS when `server.tls_cert` and `server.tls_key` are set | off |
| Rate limiting | In-memory token bucket per IP with bounded map + TTL | off |
| Security headers | CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, HSTS when TLS is on | off |
| Audit logging | Failed auth, rate-limit triggers, replay actions, TLS state via existing audit store | on |
| Hardened preset | Enables admin auth, rate limiting, and strict headers in one toggle | off |
| CSRF cookie flags | `HttpOnly` always; `Secure` when TLS is on; `SameSite=Strict` when admin auth is on | automatic |

---

## CI/CD and Shared Environments

Run OpenMuara in hardened mode in CI/CD:

```yaml
# .muara/config.yml for CI
server:
  host: 0.0.0.0

admin:
  enabled: true
  username: ci-admin
  password_hash: "$2a$10$..."  # generate with muara security hash-password

rate_limit:
  enabled: true
  requests_per_minute: 200

hardened: true
```

Set secrets via env vars instead of committing them:

```bash
export MUARA_ADMIN_PASSWORD_HASH='$2a$10$...'
export MUARA_ADMIN_TOKEN='tok_ci_only'
```

Run a posture check before starting the server:

```bash
muara security audit
```

---

## Hosted Testing and Reverse Proxies

For Mailpit-style deployments where testers access OpenMuara through a reverse proxy or Cloudflare Tunnel, see the [Hosted Testing Guide](hosted-testing.md). It covers:

- dual-port (`server.port` + `server.admin_port`) layout,
- `server.public_base_url` and `server.admin_public_base_url`,
- `viewer` accounts for read-only tester access,
- Cloudflare Tunnel, nginx, and Caddy examples,
- why payment redirects do not prompt for an admin password.

---

## Local HTTPS Testing

Generate a self-signed certificate and start the server with TLS:

```bash
muara security gen-cert --host localhost --cert-out cert.pem --key-out key.pem
muara start --config .muara/config.yml --server.tls_cert cert.pem --server.tls_key key.pem
```

For local experiments only. Use real certificates in shared environments.

---

## Security Scanning

The following checks run in CI:

- `go vet ./...`
- `golangci-lint run`
- `govulncheck ./...`
- `gosec ./...`
- Secret scanning (`gitleaks` / `trufflehog`)
- `muara security audit` in smoke tests

---

## Memory and Performance Notes

Security features are designed to stay lightweight:

- Rate limiter is in-memory with bounded map size and TTL; no Redis.
- Auth, TLS, and rate-limiter state are allocated only when enabled.
- Security headers are applied by a single middleware with negligible overhead.
- Audit events reuse the existing SQLite audit store.
