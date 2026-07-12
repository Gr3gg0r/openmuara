> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Execution Tracker

> **Updated:** 2026-07-03 | **Status:** ✅ Completed | **Archived**
>
> **Scope:** Add defense-in-depth security controls for CI/CD, shared, and port-exposed deployments.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/security-hardening` (merged to `dev`)
> **Last Agent Action:** Completed prompts 02–09: admin auth, TLS, rate limiting, security headers, audit logging, CLI helpers, CI gates, automated tests, docs/runbooks.
> **Next Agent Action:** None — initiative archived.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen |

---

## Execution Rules

1. Execute prompts in order unless marked **[PARALLEL SAFE]**.
2. Every prompt MUST end with: tests passing → git commit → update this file to `✅`.
3. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY prompt, update `HANDOFF.md`.
5. Product-code commits happen on `feat/security-hardening`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Threat model and config design | `DECISIONS.md`, `internal/config/`, `docs/security.md` | — | ✅ | `d597d34` area | Config schema added: `server.tls_cert/key`, `admin.*`, `rate_limit.*`, `hardened`. `docs/security.md` created. Decisions D001–D009 resolved. |
| 02 | Admin authentication | `internal/server/auth.go`, `internal/server/router.go`, `internal/config/`, tests | 01 | ✅ | — | Basic auth + bearer token for `/_admin/*` and admin APIs; bcrypt password storage. |
| 03 | Network binding and TLS | `internal/config/`, `internal/server/server.go`, `cmd/muara/start.go`, tests | 01 | ✅ | — | Default bind `127.0.0.1`; TLS cert/key config with HTTPS server support. |
| 04 | Rate limiting and security headers | `internal/server/ratelimit.go`, `internal/server/headers.go`, `internal/config/`, tests | 02, 03 | ✅ | — | Token-bucket per-IP rate limiter; CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, HSTS. |
| 05 | Security audit logging | `internal/server/security_audit.go`, `internal/audit/`, tests | 02, 04 | ✅ | — | Failed auth, rate-limit triggers, replay actions, TLS state logged to audit store. |
| 06 | Security CLI helpers | `cmd/muara/security.go`, `internal/cli/security.go`, tests | 01 | ✅ | — | `muara security hash-password`, `gen-cert`, and `audit` commands. |
| 07 | Security scanning & CI gates | `.github/workflows/ci.yml`, `.pre-commit-config.yaml`, `scripts/`, `Taskfile.yml` | 01 | ✅ | — | `gosec`, `gitleaks`, and `muara security audit` integrated into CI/pre-commit/local tasks. |
| 08 | Automated security tests | `internal/server/*_test.go`, `internal/cli/*_test.go` | 02, 04 | ✅ | — | Auth bypass, brute-force, rate-limit, CSRF, provider endpoint isolation, and TLS server tests. |
| 09 | Docs and runbooks | `docs/security.md`, `runbooks/on-call.md`, `runbooks/quality-gates.md`, `runbooks/local-development.md`, `README.md` | 05, 06, 07 | ✅ | — | Hardening guide, CLI reference, on-call triage, and quality-gate docs updated. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |
| Security Scan | `gosec ./...` | No high/critical issues | ✅ |
| Secret Scan | `gitleaks detect` or CI equivalent | No leaked secrets | ✅ |
| Security Audit | `muara security audit` | Reports posture, warns on insecure defaults | ✅ |

---

## Decisions

- D001 ✅ Use `server.host` as bind address (default `127.0.0.1`); no separate `server.bind`.
- D002 ✅ Admin auth supports both HTTP Basic Auth and bearer token.
- D003 ✅ Passwords stored as bcrypt hashes; env vars supported.
- D004 ✅ `hardened: true` requires admin auth + credentials; enables rate limiting and strict headers.
- D005 ✅ Provider endpoints remain unauthenticated and contract-faithful.
- D006 ✅ Add `muara security` CLI helpers: `hash-password`, `gen-cert`, `audit`.
- D007 ✅ CI must run `gosec` and secret scanning; `govulncheck` already present.
- D008 ✅ Rate limiter is in-memory, bounded, and TTL-based; no external Redis.
- D009 ✅ Security features are lazy-initialized and consume no extra resources when disabled.
- D010 ✅ Default bcrypt cost is 10; configurable via constant only.
- D011 ✅ Token-bucket per-IP rate limiter with bounded map + TTL.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/archive/openmuara-security-hardening/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/archive/openmuara-security-hardening/README.md` | Threat model, goals, approach |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
