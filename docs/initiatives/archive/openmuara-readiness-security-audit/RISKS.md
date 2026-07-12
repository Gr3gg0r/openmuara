> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Reviewed

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|----|------|------------|--------|------------|-------|
| R01 | Secrets committed to git history | Low | Critical | Run `gitleaks` / `trufflehog` on full history; rotate any leaked credentials before going public | AI Agent |
| R02 | Weak default config exposes admin API | Low | High | Default bind to `127.0.0.1`; require explicit opt-in for `0.0.0.0` + admin auth | AI Agent |
| R03 | Signature verification bypass in provider emulator | Low | High | Add conformance tests for valid, invalid, and malformed signatures | AI Agent |
| R04 | Container escape due to root user | Low | High | Run final image as non-root user; drop capabilities | AI Agent |
| R05 | Supply-chain compromise of release artifacts | Low | High | Sign checksums or binaries; publish SBOM; pin CI actions by SHA | AI Agent |
| R06 | XSS via webhook payload rendered in dashboard | Low | High | CSP + output encoding review; sanitize rendered payloads | AI Agent |
| R07 | SSRF via admin-configured webhook URLs | Low | Medium | `httputil.ValidateWebhookURL` rejects non-HTTP(S) schemes always and loopback/link-local/private IPs when `hardened: true`; revisit if per-request callback URLs are added | AI Agent |
| R08 | DoS via unbounded payload or replay endpoints | Low | Medium | Enforce body-size limits; rate-limit replay/admin endpoints | AI Agent |
| R09 | Audit log tampering because logs share DB with ledger | Low | Medium | Document scope; consider append-only log or external sink for high-assurance deployments | AI Agent |
| R10 | No vulnerability disclosure process | Medium | Medium | Add `SECURITY.md` with contact and supported versions | AI Agent |
| R11 | Dependency with incompatible or copyleft license | Low | High | License scan of Go/npm production dependencies | AI Agent |
| R12 | PII retention in demo seed data | Low | Low | Use clearly fake demo identifiers; gate seeding behind `dev.seed` | AI Agent |
| R13 | No incident response playbook delays disclosure | Medium | Medium | Create `ROLLBACK_PLAN.md` and `SECURITY.md` before release | AI Agent |
| R14 | Compromised signing key invalidates release trust | Low | High | Use OIDC/keyless signing where possible; document key rotation in `ROLLBACK_PLAN.md` | AI Agent |
| R15 | Over-hardening breaks local-first UX | Medium | Low | Keep admin auth opt-in; document hardening as opt-in for shared environments | AI Agent |
