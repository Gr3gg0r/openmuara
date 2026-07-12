> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Handoff

## Current Status

- Initiative created on `feat/security-hardening`.
- No product code changes yet.
- Prompt 01 (threat model and config design) is the next step.

## Open Questions

1. Should auth be basic auth, bearer token, or both?
2. Should hardened mode require TLS, or just strongly recommend it?
3. Do we need a CLI command to generate a bcrypt password hash or self-signed TLS cert?
4. How do we ensure middleware only protects admin routes and never provider routes?

## Next Actions

- [ ] Complete prompt 01 and lock down the config schema in `DECISIONS.md`.
- [ ] Implement admin authentication in prompt 02.
