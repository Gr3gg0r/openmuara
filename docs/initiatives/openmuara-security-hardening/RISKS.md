> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | Default `127.0.0.1` binding breaks Docker or remote dev workflows. | Medium | Medium | Allow explicit `server.bind: 0.0.0.0`; document Docker usage. |
| R02 | Admin auth adds friction to local development. | High | Low | Auth is opt-in; `hardened: true` or explicit `admin.enabled: true` required. |
| R03 | Password hashes in config are accidentally committed. | Medium | High | Document that config files should be ignored by git; support env vars for secrets. |
| R04 | Rate limiting breaks legitimate load tests. | Medium | Medium | Make rate limits configurable; allow disabling in local config. |
| R05 | TLS cert config is cumbersome for local dev. | Medium | Low | TLS optional; generate self-signed cert via CLI helper if desired. |
| R06 | CSRF bypass on admin simulation endpoints. | Low | High | Keep existing CSRF double-submit cookie; extend to all state-changing admin actions. |
| R07 | Webhook payloads logged in audit logs leak PII. | Medium | Medium | Exclude raw payloads from audit logs; log metadata only. |
| R08 | Provider endpoints abused as open relay. | Medium | High | Default localhost binding + rate limiting; document never exposing provider endpoints publicly. |
| R09 | Security controls accidentally added to provider emulation endpoints. | Low | High | Explicitly scope auth and admin-only middleware to `/_admin/*` routes; test provider endpoints remain contract-faithful. |
