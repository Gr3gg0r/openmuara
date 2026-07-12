> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Decision Log

| ID | Decision | Status | Date | Notes |
|----|----------|--------|------|-------|
| D000 | Core philosophy | ✅ Decided | 2026-07-02 | Security controls apply to admin/internal surface only. Provider emulation endpoints remain drop-in replacements of real provider APIs. |
| D001 | Default bind address | ⬜ Pending | — | Proposed: `127.0.0.1` with opt-in `0.0.0.0`. |
| D002 | Admin authentication mechanism | ⬜ Pending | — | Basic auth and/or bearer token for `/_admin/*` only. |
| D003 | Password storage | ⬜ Pending | — | Bcrypt hashes in config; env var injection for secrets. |
| D004 | Hardened mode preset | ⬜ Pending | — | `hardened: true` enables auth + TLS + strict CORS + admin rate limiting. Provider endpoints unchanged. |
| D005 | Provider endpoint auth | ✅ Decided | 2026-07-02 | Remain unauthenticated and contract-faithful (emulate public provider APIs). Rate limiting on provider endpoints is opt-in only under hardened mode. |
