> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit

> **Status:** ✅ Complete | **Started:** 2026-07-08
> **Scope:** Gold-standard security review of the OpenMuara codebase, configuration, CI/CD, container image, and release process before public release.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/readiness-security-audit` (to be created when work starts)
> **Reference:** [`docs/security.md`](../../security.md), [`.github/workflows/ci.yml`](../../../.github/workflows/ci.yml)

---

## Why this matters

OpenMuara emulates real payment providers, handles HMAC signatures, dispatches webhooks, and exposes an admin dashboard. Before publishing, we need defense-in-depth: no secrets in history, no unsafe defaults, no weak crypto, no auth bypasses, and a clear vulnerability-disclosure path.

This initiative treats security as a readiness gate, not a one-off scan. We will audit, fix, and add regression tests so the project can confidently accept external contributors and users.

## Initiative structure

```
docs/initiatives/openmuara-readiness-security-audit/
├── README.md              # This file
├── EXECUTION_PLAN.md      # Timeline, milestones, RACI
├── TRACKING.md            # Central execution tracker
├── KNOWN_ISSUES.md        # Catalog of security findings
├── RECOMMENDATIONS.md     # Recommended fixes and priorities
├── ATTACKER_SCENARIOS.md  # Scenario-based test cases
├── REVIEW_CHECKLIST.md    # Sign-off checklist
├── ROLLBACK_PLAN.md       # Incident response and rollback plan
├── RISKS.md               # Risk register
├── DECISIONS.md           # Decision log
└── HANDOFF.md             # Session continuity
```

## Standards & frameworks mapped

| Standard / Framework | How this initiative uses it |
|---|---|
| **OWASP ASVS 4.0** | Auth (V2), session management (V3), access control (V4), validation (V5), cryptography (V6), error handling (V7), data protection (V8), communications (V9), malicious code (V10), logging (V7) |
| **SLSA Level 1–2** | Signed release artifacts, SBOM, provenance, pinned CI actions, reproducible build verification |
| **CIS Docker Benchmark** | Non-root container user, minimal image, no unnecessary packages, image vulnerability scanning |
| **NIST SSDF** | Prepare the organization (PO), protect software (PS), produce well-secured software (PW), respond to vulnerabilities (RV) |
| **GitHub Supply-chain Security** | Dependency review, secret scanning, code scanning SARIF, pinned actions, OIDC attestations |

## RACI

| Activity | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| Run scans & triage findings | R | A | C |
| Define accepted risks | C | A | R |
| Approve `SECURITY.md` content | C | A | R |
| Approve Dockerfile changes | R | A | C |
| Approve CI/release hardening | R | A | C |
| Final sign-off | C | A | R |

*R = Responsible, A = Accountable, C = Consulted, I = Informed*

## Dependencies on other initiatives

| Initiative | Why it matters for security |
|---|---|
| [Dependency & License Audit](../openmuara-readiness-dependency-license-audit/README.md) | License compatibility and vulnerability scanning overlap with supply-chain security |
| [CI & Release Audit](../openmuara-readiness-ci-release-audit/README.md) | Release signing, install script, and Docker build are shared concerns |
| [Docs Completeness Audit](../openmuara-readiness-docs-completeness-audit/README.md) | `docs/security.md`, `SECURITY.md`, and runbooks must stay accurate |
| [Coverage Audit](../openmuara-readiness-coverage-audit/README.md) | Security fixes need regression tests; coverage targets must accommodate them |

## Existing controls (do not regress)

| Control | Location | Notes |
|---|---|---|
| Admin Basic Auth / bearer token | `internal/server/auth.go` | bcrypt password hashes, constant-time token compare |
| Viewer role | `internal/server/auth.go` | Read-only; cannot replay webhooks or change config |
| Security headers (CSP, HSTS, etc.) | `internal/server/headers.go` | Applied to `/_admin/*` |
| CSRF double-submit cookie | `internal/server/csrf.go` | 32-byte random token, exempts webhook receiver |
| Token-bucket rate limiter | `internal/server/ratelimit.go` | In-memory, bounded, with TTL eviction |
| Security audit CLI | `internal/cli/security.go` | `hash-password`, `gen-cert`, `audit` |
| Security event logging | `internal/server/security_audit.go` | Auth success/failure, rate limits, replays, TLS state |
| Signature fuzz tests | `internal/*/signature_fuzz_test.go` | Fuzzing for provider signatures |
| CI security jobs | `.github/workflows/ci.yml` | `govulncheck`, `gosec`, `gitleaks` |

## Audit areas

1. **Threat modeling & asset inventory** — identify what must be protected and who can attack it.
2. **Static & dependency analysis** — `gosec`, `govulncheck`, `gitleaks`, `npm audit`, license scan.
3. **Authentication & authorization** — admin/viewers, simulation routes, role enforcement, token lifecycle.
4. **Cryptography** — bcrypt cost, HMAC/SHA256 correctness, randomness, TLS cert generation.
5. **Input validation & web attack surface** — SQL injection, XSS, SSRF, path traversal, command injection, deserialization.
6. **Webhook security** — signature verification, replay protection, payload limits, idempotency.
7. **Audit logging & integrity** — tamper evidence, sensitive-data redaction, retention.
8. **Configuration & defaults** — bind address, hardened mode, secret handling, env vars, CORS.
9. **Container & supply-chain security** — Dockerfile user, minimal image, SBOM, signed releases, provenance.
10. **CI/CD & release security** — pinned actions, least-privilege workflow permissions, artifact signing.
11. **Incident response & disclosure** — `SECURITY.md`, contact process, supported versions.

## Success metrics

| Metric | Target | Measurement |
|---|---|---|
| `gosec` high/critical findings | 0 | SARIF output / CI job |
| `govulncheck` high/critical vulnerabilities | 0 | CI job output |
| `npm audit` high/critical advisories | 0 | `npm audit --production` |
| Secrets in git history | 0 | `gitleaks` + `trufflehog` |
| Container image critical/high CVEs | 0 | `trivy` or `grype` scan |
| Admin route auth test coverage | 100% | `go test -cover` on `internal/server` |
| Provider signature negative tests | ≥1 per provider | Test files |
| Release artifacts with checksums | 100% | Release workflow |
| `SECURITY.md` existence | Yes | File check |
| `muara security audit` clean on reference configs | Yes | Manual / CI |

## Success criteria

- No high or critical findings from `gosec`, `govulncheck`, or `npm audit`.
- `gitleaks` / `trufflehog` report zero secrets in the working tree and full history.
- `muara security audit` passes on all documented deployment scenarios.
- Dockerfile runs as non-root and passes an image scan (`trivy` or `grype`).
- `SECURITY.md` exists with a clear disclosure process.
- All security fixes have regression tests.
- All quality gates pass after each remediation batch.

See [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) for the prioritized action plan, [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) for timeline and RACI, and [`ATTACKER_SCENARIOS.md`](ATTACKER_SCENARIOS.md) for scenario-based tests.
