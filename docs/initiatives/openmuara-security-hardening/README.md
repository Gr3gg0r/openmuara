> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Add defense-in-depth security controls to OpenMuara's admin surface and deployment footprint without changing the behavior of provider emulation endpoints.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/security-hardening`

---

## Initiative Structure

```
docs/initiatives/openmuara-security-hardening/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    ├── 01-threat-model-and-config-design.md
    ├── 02-admin-authentication.md
    ├── 03-network-binding-and-tls.md
    ├── 04-rate-limiting-and-csp.md
    └── 05-security-audit-logging.md
```

Planning docs live in `docs/initiatives/openmuara-security-hardening/` in the root repo.
Product-code commits to the `feat/security-hardening` branch. Do not commit directly to `main`.

---

## Core Philosophy: Drop-In Emulator First

OpenMuara's product philosophy is that a user can integrate it by **only changing the base URL** in their production provider client. No code changes, no extra headers, no special auth flows.

This initiative **must not violate that philosophy**:

- **Provider emulation endpoints** (`/fawry/charge`, `/v1/checkout/sessions`, `/senangpay/charge`, etc.) must continue to behave exactly like the real provider endpoints they emulate. If the real provider endpoint is public and signature-verified, OpenMuara stays public and signature-verified. We do not add our own auth layer on top of provider routes.
- **Security controls apply to the admin/internal surface**: `/_admin`, admin JSON APIs, the dashboard, replay, simulation, and configuration endpoints.
- **Default deployment posture** is local and safe: bind to `127.0.0.1`, no auth required. Hardening is opt-in for CI/CD or shared environments.
- A production-ready OpenMuara deployment is typically on the same trusted host/network as the test application (localhost, Docker network, CI runner). The hardening in this initiative is for the case where that boundary is wider or shared.

In short: **the payment-emulation API stays a drop-in replacement; the dashboard gets a lock.**

---

## Why Security Hardening?

OpenMuara is "Mailpit for payment emulators." Mailpit ships with authentication, bind-address configuration, TLS, and other hardening options because it is often run in shared environments where the web UI could leak sensitive data. OpenMuara faces the same risks, but only on its **admin/internal surface**:

- Emulated transactions, webhook payloads, and provider secrets can be inspected through `/_admin`.
- The dashboard, replay, and simulation endpoints are powerful — an attacker could replay webhooks, mark transactions paid, or exfiltrate provider configuration.
- In CI/CD, the server may bind to `0.0.0.0` accidentally, exposing the dashboard to the network.

Provider emulation endpoints are intentionally public (like the real providers they emulate). The security boundary is the admin UI and the network binding, not the provider API contract.

This initiative makes OpenMuara safe to run in untrusted or shared environments without breaking the local-first, zero-config default experience and without changing provider drop-in behavior.

---

## Goals

1. Provide configurable authentication for `/_admin` and admin API endpoints only.
2. Default to `127.0.0.1` binding so the server is not exposed to the network unless explicitly configured.
3. Support TLS/HTTPS for environments that require it.
4. Add opt-in rate limiting to admin endpoints; provider endpoints remain unchanged by default.
5. Harden HTTP security headers for the admin UI (CSP, HSTS, etc.).
6. Add security-relevant audit logging (failed logins, config changes, replay actions).
7. Provide a `--hardened` / `hardened: true` configuration mode that turns on a secure-by-default preset.
8. Document hardening guidance for CI/CD and shared environments.
9. Preserve the drop-in provider replacement philosophy: no auth or behavior changes on provider emulation routes.

---

## Reference: Mailpit Security Model

Mailpit provides the following controls, which this initiative adapts for OpenMuara:

| Mailpit Feature | OpenMuara Equivalent |
|---|---|
| `MP_UI_AUTH` (basic auth for UI) | `admin.auth` config: basic auth or token for `/_admin` |
| `MP_UI_TLS_CERT` / `MP_UI_TLS_KEY` | `server.tls_cert` / `server.tls_key` |
| `MP_HTTP_BIND_ALL` (default false) | `server.bind` default `127.0.0.1`; opt-in to `0.0.0.0` |
| `MP_CORS_ORIGIN` | Existing CORS config; tighten defaults |
| `MP_BLOCK_REMOTE_CSS_AND_JS` | CSP `default-src 'self'` for admin UI |
| `MP_MAX_MESSAGES` / `MP_MAX_ATTACHMENT_SIZE` | Request body limits already exist; review admin endpoints |

---

## Threat Model (High Level)

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

Note: Provider endpoint abuse is mitigated primarily by binding to localhost and keeping test instances private, not by adding auth to the provider API contract.

---

## Recommended Approach

### Phase 1 — Config and network hardening

- Change default bind address to `127.0.0.1`.
- Add `server.bind`, `server.port`, `server.tls_cert`, `server.tls_key` config.
- Add `hardened: true` preset that enables auth + TLS + strict CORS + rate limiting.

### Phase 2 — Admin authentication

- Add HTTP Basic Auth or bearer token for `/_admin/*` and admin JSON APIs.
- Store password hashes with bcrypt (not plaintext).
- Support `admin.username` / `admin.password_hash` or `admin.token` in config.
- Allow disabling auth explicitly for local dev (`admin.enabled: false`).

### Phase 3 — Rate limiting and security headers

- Per-IP rate limiting on admin endpoints and replay actions.
- Optional rate limiting on provider endpoints only when `hardened: true` is enabled, with sane defaults that do not break normal test traffic.
- Stricter CSP for admin UI (`default-src 'self'`).
- Add `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`.

### Phase 4 — Security audit logging

- Log failed auth attempts, config changes, replay actions, and TLS enablement.
- Expose security events via CLI/API.

---

## Non-Goals

- Do not turn OpenMuara into a production payment gateway.
- Do not add authentication or behavior changes to provider emulation endpoints.
- Do not break the drop-in base-URL replacement philosophy.
- Do not add complex RBAC or multi-user sessions in v1.
- Do not break the local dev workflow when no auth is configured.
- Do not require internet access for security features.
