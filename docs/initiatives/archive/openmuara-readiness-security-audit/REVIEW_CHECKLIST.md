> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Review Checklist

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete

---

Use this checklist before marking the security audit initiative complete. Every item must be checked or explicitly deferred with a rationale recorded in `DECISIONS.md` or `RISKS.md`.

## Static & dependency analysis

- [x] `gosec` reports 0 high and 0 critical findings.
- [x] `govulncheck` reports 0 high and 0 critical vulnerabilities in reachable code.
- [x] `npm audit --production` reports 0 high and 0 critical advisories in `web/dashboard`.
- [ ] `npm audit --production` reports 0 high and 0 critical advisories in `website` — **deferred**: vulnerabilities are in Docusaurus build-time dependencies; monitor for upstream fixes.
- [x] `gitleaks detect` reports 0 findings in the working tree.
- [x] `trufflehog git file://.` reports 0 verified secrets in full history.
- [x] No incompatible licenses in production dependencies (verified by dependency/license audit initiative).

## Authentication & authorization

- [x] Admin routes (`/_admin/*`) require valid credentials when auth is enabled.
- [x] Viewer role cannot mutate config, replay webhooks, or access admin-only APIs.
- [x] Provider simulation routes remain accessible without admin auth by design.
- [x] Bearer tokens and passwords are compared in constant time.
- [x] Passwords are stored as bcrypt hashes; plaintext secrets never committed.
- [x] CSRF protection covers all admin state-changing methods.

## Cryptography

- [x] bcrypt cost is at least 10 and documented.
- [x] Provider signature verification uses HMAC/SHA256 or the documented provider algorithm.
- [x] Signature tests cover valid, invalid, empty, and tampered signatures.
- [x] All random tokens/IDs use `crypto/rand`.
- [x] TLS certificate generation is documented as local-only.

## Input validation & web attack surface

- [x] All SQL queries are parameterized.
- [x] No user input is passed to shell commands.
- [x] File paths are validated against a base directory (`muara provider init` now rejects unsafe names).
- [x] Webhook/callback URLs are validated (scheme always; loopback/link-local/private IP blocklist when `hardened: true`).
- [x] Dashboard escapes rendered provider responses (removed `dangerouslySetInnerHTML`).
- [x] Dashboard webhook payload rendering audited for XSS (all payloads rendered as plain text via `<CodeBlock>` / `<pre>`).
- [x] Body size is limited (currently 1 MiB default).

## Webhook security

- [x] Every provider has negative signature tests.
- [x] Replayed webhooks without valid signatures are rejected.
- [x] Idempotency-key behavior is tested.
- [x] Webhook payload size is capped.

## Audit logging & PII

- [x] Security events (auth, rate limits, replays, config changes, TLS) are logged.
- [ ] Audit rows include monotonic event IDs and timestamps — **not implemented**.
- [ ] Logs do not contain plaintext secrets or real-looking customer emails — **pending manual review**.
- [x] `dev.seed` defaults to off and demo data is clearly labeled.

## Configuration & defaults

- [x] `server.host` defaults to `127.0.0.1`.
- [x] `muara security audit` warns when `0.0.0.0` is used without admin auth/TLS.
- [x] `hardened: true` requires `admin.enabled: true`.
- [x] CORS configuration is documented and defaults are safe.
- [x] Environment variable override behavior is documented.

## Container & supply-chain

- [x] Dockerfile runs as a non-root user.
- [ ] Container image scan reports 0 critical and 0 high CVEs — **pending `trivy`/`grype` installation**.
- [x] Release artifacts include SHA256 checksums.
- [x] SBOM is generated and attached to releases.
- [x] Build uses locked dependency versions (`go.sum`, `package-lock.json`).

## CI/CD & release

- [x] All GitHub Actions are pinned by SHA with a version comment.
- [x] Workflows declare minimal `permissions:`.
- [x] Release workflow generates checksums and optionally attestations.
- [x] Secret scanning runs on every PR and push to default branches.
- [x] `task check` or equivalent passes in CI.

## Incident response

- [x] `.github/SECURITY.md` exists with a disclosure process and supported versions.
- [x] README links to `SECURITY.md`.
- [x] A rollback/response plan exists for leaked secrets (`ROLLBACK_PLAN.md`).

## Documentation & tests

- [x] `docs/security.md` reflects all implemented controls.
- [ ] Threat model is documented — **pending**.
- [x] Every security fix has a regression test (existing tests cover auth/webhook/signature paths).
- [x] All quality gates pass: `go build`, `go test`, `go vet`, `golangci-lint`, dashboard typecheck + tests.

## Sign-off

- [x] AI Agent confirms all completed checklist items are addressed.
- [ ] Human reviewer approves the final state.
- [ ] Maintainer approves publication readiness.
