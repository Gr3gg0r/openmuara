> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Decision Log

| ID | Decision | Status | Date | Notes |
|----|----------|--------|------|-------|
| D000 | Core philosophy | ✅ Decided | 2026-07-02 | Security controls apply to admin/internal surface only. Provider emulation endpoints remain drop-in replacements of real provider APIs. |
| D001 | Default bind address | ✅ Decided | 2026-07-03 | Use existing `server.host` (default `127.0.0.1`) as the bind address. No separate `server.bind` field. Opt-in to `0.0.0.0` only with admin auth or TLS. |
| D002 | Admin authentication mechanism | ✅ Decided | 2026-07-03 | Support both HTTP Basic Auth (`admin.username` + `admin.password_hash`) and bearer token (`admin.token`). A request satisfies auth if either is valid. |
| D003 | Password storage | ✅ Decided | 2026-07-03 | Bcrypt hashes in config; plaintext not allowed. Secrets injectable via env vars (`MUARA_ADMIN_PASSWORD_HASH`, `MUARA_ADMIN_TOKEN`). |
| D004 | Hardened mode preset | ✅ Decided | 2026-07-03 | `hardened: true` requires `admin.enabled: true` and valid admin credentials. Enables rate limiting and strict security headers. Provider endpoints unchanged. |
| D005 | Provider endpoint auth | ✅ Decided | 2026-07-02 | Remain unauthenticated and contract-faithful (emulate public provider APIs). Rate limiting on provider endpoints is opt-in only under hardened mode. |
| D006 | Security CLI helpers | ✅ Decided | 2026-07-03 | Add `muara security hash-password`, `gen-cert`, and `audit` commands to make hardening self-service for devs and CI. |
| D007 | CI security scanning | ✅ Decided | 2026-07-03 | Add `gosec` and secret scanning (e.g., `gitleaks`) to CI/pre-commit; keep existing `govulncheck`. |
| D008 | Rate limiter implementation | ✅ Decided | 2026-07-03 | In-memory token-bucket or sliding-window with bounded map + TTL. No external dependencies. |
| D009 | Lazy initialization | ✅ Decided | 2026-07-03 | Security middleware, rate limiter, TLS, and auth state are only allocated when enabled. Baseline memory and CPU must not increase when hardening is off. |
| D010 | Bcrypt cost | ✅ Decided | 2026-07-03 | Default bcrypt cost is 10. It is configurable via a constant but not exposed in config file to avoid unsafe values. |
| D011 | Rate limiter algorithm | ✅ Decided | 2026-07-03 | Use token bucket per IP. Buckets are stored in a bounded map with TTL eviction. No external store. |
| D012 | TLS server implementation | ✅ Decided | 2026-07-03 | `internal/server.Server` switches to `http.ServeTLS` when both `server.tls_cert` and `server.tls_key` are configured; no reverse proxy required. |
| D013 | Security scan integration | ✅ Decided | 2026-07-03 | `gosec` and `gitleaks` run in CI and pre-commit. Local tasks (`task security`, `task secrets`) skip gracefully if tools are missing. |
| D014 | Self-signed cert generation | ✅ Decided | 2026-07-03 | `muara security gen-cert` uses ECDSA P-256 and writes cert/key files for local HTTPS testing only. |
| D015 | Security audit command output | ✅ Decided | 2026-07-03 | `muara security audit` prints posture states and issues; it never prints secrets or sensitive values. |
| D016 | CSRF cookie flags | ✅ Decided | 2026-07-03 | The CSRF cookie is always `HttpOnly`; `Secure` is set only when TLS is enabled; `SameSite=Strict` is set when admin auth is enabled. |
