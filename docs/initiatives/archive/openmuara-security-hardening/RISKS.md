> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | Default `127.0.0.1` binding breaks Docker or remote dev workflows. | Medium | Medium | Allow explicit `server.host: 0.0.0.0`; document Docker usage. |
| R02 | Admin auth adds friction to local development. | High | Low | Auth is opt-in; `hardened: true` or explicit `admin.enabled: true` required. |
| R03 | Password hashes in config are accidentally committed. | Medium | High | Document that config files should be ignored by git; support env vars for secrets. |
| R04 | Rate limiting breaks legitimate load tests. | Medium | Medium | Make rate limits configurable; allow disabling in local config. |
| R05 | TLS cert config is cumbersome for local dev. | Medium | Low | TLS optional; generate self-signed cert via CLI helper if desired. |
| R06 | CSRF bypass on admin simulation endpoints. | Low | High | Keep existing CSRF double-submit cookie; extend to all state-changing admin actions. |
| R07 | Webhook payloads logged in audit logs leak PII. | Medium | Medium | Exclude raw payloads from audit logs; log metadata only. |
| R08 | Provider endpoints abused as open relay. | Medium | High | Default localhost binding + rate limiting; document never exposing provider endpoints publicly. |
| R09 | Security controls accidentally added to provider emulation endpoints. | Low | High | Explicitly scope auth and admin-only middleware to `/_admin/*` routes; test provider endpoints remain contract-faithful. |
| R10 | Security scanning tools add CI time or local dev friction. | Medium | Low | Run heavy scans in CI only; keep `gosec` optional locally; cache tool installs. |
| R11 | Rate limiter grows unbounded in memory under many IPs. | Low | High | Cap map size + TTL entries; evict stale entries; reset on process restart is acceptable for local emulator. |
| R12 | Self-signed cert generation or bcrypt hashing is slow or blocks startup. | Low | Medium | Run cert gen and password hashing via CLI commands, not on every startup; use sensible bcrypt cost. |
| R13 | Security audit CLI leaks secrets in output. | Low | High | `muara security audit` must redact `password_hash`, `token`, and TLS key paths; only report presence/state, not values. |
| R14 | CSRF cookie lacks secure flags when TLS is enabled. | Low | High | ✅ Set `Secure` when TLS is enabled and `SameSite=Strict` when admin auth is enabled; `HttpOnly` always. |
| R15 | Large admin request bodies cause memory exhaustion. | Low | Medium | Enforce body size limits and read/write timeouts on `/_admin` routes. |
