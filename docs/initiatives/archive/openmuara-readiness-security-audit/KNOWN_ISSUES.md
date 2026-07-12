> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Reviewed

---

## Confirmed findings

| ID | Finding | Area | Severity | Status | Fixed in | Notes |
|----|---------|------|----------|--------|----------|-------|
| F01 | `gitleaks` flagged shell-script placeholders in examples and smoke tests | Secrets | Low | ✅ Closed | Added `.gitleaks.toml` allowlist | Patterns like `curl -u "${API_KEY}:"` are environment-variable placeholders, not real secrets |
| F02 | `govulncheck` reports 1 vulnerability in a required module that OpenMuara does not call | Dependencies | Medium | ✅ Accepted | N/A | Module is required but no vulnerable symbol is reachable; monitor for updates |
| F03 | `npm audit --production` reports 21 vulnerabilities in Docusaurus transitive build dependencies | Dependencies | High (transitive) | ⚠️ Accepted | N/A | Vulnerabilities are in build-time tooling (`serialize-javascript`, `uuid` via `webpack-dev-server`); the generated static site does not execute them. Monitor Docusaurus updates. |
| F04 | `muara provider init <name>` allowed path traversal in provider directory name | Input validation | High | ✅ Closed | `internal/cli/plugins.go` | Added `isSafeProviderName` validation and regression tests |
| F05 | Dashboard used `dangerouslySetInnerHTML` with provider display text | XSS | Medium | ✅ Closed | `web/dashboard/src/components/Providers.tsx` | Replaced with plain JSX text rendering |

## Candidate gaps (to verify during the audit)

These are not confirmed issues; they are areas the audit should explicitly test.

| ID | Candidate | Area | Why it matters | Verify with | Recommended fix |
|----|-----------|------|----------------|-------------|-----------------|
| C01 | Dockerfile runs as root | Container security | Root in a container increases blast radius on escape | `docker run --user`, image scan | Add `USER muara` and group; run image scan in CI |
| C02 | Missing `SECURITY.md` | Incident response | External researchers need a disclosure channel | File existence, contact process | Add `.github/SECURITY.md` with email and supported versions |
| C03 | No release artifact signing | Supply chain | Users cannot verify binary integrity | Release workflow, checksums, cosign | Generate SHA256 checksums; add cosign/GitHub attestations |
| C04 | No SBOM published | Supply chain | Users cannot audit transitive dependencies | `syft` / `gomod` SBOM | Generate SBOM in CI and attach to releases |
| C05 | Workflow actions not pinned by SHA | CI/CD security | Tag-based actions can be retagged maliciously | `.github/workflows/*.yml` | Pin every `uses:` to a SHA with a version comment |
| C06 | `GITLEAKS_LICENSE` or token exposure risk | Secrets | Enterprise scanners may need license keys | `gitleaks` history scan | Verify scanner config is not committed; use repo-level allowlist |
| C07 | Admin dashboard served from same origin as provider endpoints | Architecture | XSS on provider pages could target admin cookies | Threat model review | Keep strict CSP; consider dual-port admin in future |
| C08 | Audit logs stored in same SQLite DB as ledger | Audit integrity | Compromised DB could alter logs | Tamper-evidence review | Add monotonic event IDs + timestamps; document limitation |
| C09 | Default `dev.seed` may write PII-like demo emails | PII handling | Demo data could be mistaken for real customer data | Config/docs review | Ensure `dev.seed` defaults to false; label demo data clearly |
| C10 | Self-signed cert uses ECDSA P-256 with 1-year expiry | Cryptography | Fine for local, but docs must warn against production use | `muara security gen-cert` review | Add docs warning; recommend real certs in production |

## Closed findings

- F01 — Gitleaks false positives resolved via `.gitleaks.toml`.
- F04 — Provider-init path traversal fixed.
- F05 — Dashboard `dangerouslySetInnerHTML` removed.
- F06 — SSRF protection added for admin-configured webhook URLs (`webhook.url` and `webhook.targets`): non-HTTP(S) schemes are rejected always; loopback/link-local/private IPs are rejected in `hardened` mode. See `internal/httputil/url.go` and `internal/config/validation.go`.
- F07 — Default config validation no longer requires disabled providers to be registered, so bundled provider templates do not block startup.
- F08 — Docker entrypoint defaults `MUARA_SERVER_HOST` to `0.0.0.0`; `internal/server/server.go` binds IPv4 sockets explicitly to fix Docker Desktop port-forwarding behavior.
