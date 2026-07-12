> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Recommendations

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

These recommendations are planning-only. They map each audit area to a concrete, industry-standard action. Execute them in priority order once the initiative is approved.

## Priority matrix

| Priority | Area | Recommendation | Effort | Impact | Owner |
|----------|------|----------------|--------|--------|-------|
| P0 | Secrets | Run `gitleaks` and `trufflehog` on full history; rotate anything found before going public | Low | Critical | AI Agent |
| P0 | Static analysis | Fix any `gosec` high/critical findings and `govulncheck` vulnerabilities | Low-Medium | High | AI Agent |
| P0 | Auth defaults | Keep `server.host=127.0.0.1` default; require explicit opt-in for `0.0.0.0` + admin auth | Low | High | AI Agent |
| P0 | `SECURITY.md` | Add a vulnerability disclosure policy with contact email and supported versions | Low | High | AI Agent |
| P1 | Container | Add a non-root `muara` user to the Dockerfile; drop unnecessary capabilities; scan with `trivy` or `grype` | Low | High | AI Agent |
| P1 | CI hardening | Pin all third-party GitHub Actions by SHA; restrict `permissions:` to least privilege | Low | High | AI Agent |
| P1 | Release signing | Generate SHA256 checksums for release artifacts; add cosign/GitHub attestations when feasible | Low-Medium | High | AI Agent |
| P1 | SBOM | Generate and attach an SBOM (Go + npm) to releases using `syft` or `go version -m` | Low | Medium | AI Agent |
| P1 | Input validation | Audit all `r.ParseForm`, SQL queries, file paths, and external URL fetches for injection/SSRF | Medium | High | AI Agent |
| P1 | Webhook security | Add negative tests for missing/invalid signatures, enforce idempotency keys, cap payload size | Medium | High | AI Agent |
| P2 | Audit integrity | Append a monotonic event ID and timestamp to audit rows; document that logs share the ledger DB | Low | Medium | AI Agent |
| P2 | XSS / CSP | Review all places webhook payloads or provider responses are rendered; keep CSP strict | Medium | Medium | AI Agent |
| P2 | PII handling | Redact or tokenize real-looking emails in logs; ensure `dev.seed` is off by default and clearly labeled | Low | Medium | AI Agent |
| P2 | Cryptography docs | Clarify that `muara security gen-cert` is for local testing only; recommend real certs in production | Low | Low | AI Agent |
| P3 | Threat model doc | Publish a short threat model in `docs/security.md` with assets, trust boundaries, and accepted risks | Medium | Medium | AI Agent |
| P3 | Dual-port admin | Evaluate splitting admin dashboard to a separate port/host for high-assurance deployments | High | Medium | AI Agent |

## Standards mapping

| Recommendation | OWASP ASVS 4.0 | SLSA | CIS Docker | NIST SSDF |
|---|---|---|---|---|
| Secret scanning | V2.10, V6.4 | — | — | PO.3, RV.1 |
| Static analysis (`gosec`, `govulncheck`) | V1.2, V10 | Build L2 | — | PW.7, PW.8 |
| Auth / viewer role separation | V2, V4.1 | — | — | PW.5 |
| Cryptography review | V6 | — | — | PW.6 |
| Input validation | V5 | — | — | PW.5 |
| Webhook signature tests | V9, V10 | — | — | PW.6 |
| Audit logging | V7.1, V8 | — | — | RV.1 |
| Container non-root user | — | — | 4.1, 4.6 | PS.2 |
| SBOM + provenance | — | L1–L2 | — | PO.1, PS.1 |
| Release signing | — | L2–L3 | — | PS.2 |
| Pinned CI actions | — | L2 | — | PO.4 |
| `SECURITY.md` | — | — | — | RV.1, RV.2 |

## Recommended tool stack

| Purpose | Tool | Where |
|---|---|---|
| Go static analysis | `gosec`, `govulncheck` | CI + local `task security` |
| Secret scanning | `gitleaks`, `trufflehog` | CI + pre-commit hook |
| Image scanning | `trivy`, `grype` | CI release job |
| SBOM generation | `syft`, `go version -m` | Release workflow |
| Release signing | `cosign` keyless, GitHub attestations | Release workflow |
| Link / docs check | `markdown-link-check` | CI docs job |

## Copy-paste command reference

```bash
# Static analysis
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
gosec -fmt sarif -out gosec.sarif ./...
govulncheck ./...

# Secret scanning
gitleaks detect --source . --verbose
trufflehog git file://. --only-verified

# Dependency audit
cd web/dashboard && npm audit --production
cd website && npm audit --production
go mod tidy && go mod verify

# Container security
docker build -t openmuara:audit .
trivy image openmuara:audit
grype openmuara:audit

# SBOM
syft openmuara:audit -o spdx-json=sbom.spdx.json

# Release checksums
sha256sum muara_* > checksums.txt
```

## What not to do

- Do **not** add authentication to provider emulation endpoints; that would break the drop-in contract.
- Do **not** require admin auth by default; keep local-first UX on `127.0.0.1`.
- Do **not** introduce Redis or external services for rate limiting; the in-memory bounded limiter is appropriate for the target use case.
- Do **not** rewrite public git history unless a confirmed secret is found and the repo is still private.
