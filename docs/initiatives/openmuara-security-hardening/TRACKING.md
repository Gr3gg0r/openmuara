> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Security Hardening — Execution Tracker

> **Updated:** 2026-07-08 | **Status:** ✅ Completed
>
> **Scope:** Add defense-in-depth security controls for CI/CD, shared, and port-exposed deployments.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/security-hardening`
> **Last Agent Action:** Verified existing implementation, fixed provider simulation route auth gap, and added hosted-testing docs.
> **Next Agent Action:** None.

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
| 01 | Threat model and config design | `DECISIONS.md`, `internal/config/`, `docs/security.md` | — | ✅ | — | Threat model and config schema are implemented (`admin`, `viewer`, `hardened`, `rate_limit`, TLS). |
| 02 | Admin authentication | `internal/server/auth.go`, `internal/server/router.go`, `internal/config/`, tests | 01 | ✅ | — | Basic auth and bearer token for `/_admin/*`; `viewer` role for read-only access. |
| 03 | Network binding and TLS | `internal/config/`, `internal/server/server.go`, `cmd/muara/start.go`, tests | 01 | ✅ | — | Default bind `127.0.0.1`; TLS cert/key config; `public_base_url` and `admin_public_base_url`. |
| 04 | Rate limiting and security headers | `internal/server/middleware.go`, `internal/config/`, tests | 02, 03 | ✅ | — | Per-IP token-bucket rate limiting and CSP/security headers. |
| 05 | Security audit logging | `internal/audit/`, `internal/server/`, tests | 02, 04 | ✅ | — | Failed auth, replay, config changes, and TLS state logged to audit store. |
| 06 | Docs and runbooks | `docs/security.md`, `docs/hosted-testing.md`, `docs/operations.md` | 05 | ✅ | — | Hardening guide and new hosted-testing guide added. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Security Scan | `gosec ./...` | No high/critical issues | ✅ |

---

## Decisions

- D001 ✅ Default bind address: `127.0.0.1` instead of `0.0.0.0`.
- D002 ✅ Admin authentication mechanism: HTTP Basic Auth + bearer token; separate `viewer` role for read-only access.
- D003 ✅ Password storage: bcrypt hashes in config; no plaintext passwords.
- D004 ✅ `hardened: true` preset: enables auth + rate limiting + strict security headers.
- D005 ✅ Provider endpoints remain unauthenticated by default (they emulate public gateway APIs).
- D006 ✅ Provider simulation/payment pages under `/_admin/*` are exempt from admin auth so redirected browsers can complete payments.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-security-hardening/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-security-hardening/README.md` | Threat model, goals, approach |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
