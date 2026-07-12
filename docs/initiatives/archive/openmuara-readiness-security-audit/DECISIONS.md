> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Security Audit Decision Log

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Reviewed

---

| ID | Decision | Context | Status | Date |
|----|----------|---------|--------|------|
| D01 | Provider emulation endpoints remain public | Real provider APIs are public drop-in replacements; auth would break contract fidelity | ✅ Accepted (pre-existing) | 2026-07-08 |
| D02 | Admin auth is opt-in to preserve local-first UX | Default local use on `127.0.0.1` should not require credentials | ✅ Accepted (pre-existing) | 2026-07-08 |
| D03 | In-memory rate limiter is acceptable for local/shared test use | No Redis dependency; bounded map + TTL prevents unbounded growth | ✅ Accepted (pre-existing) | 2026-07-08 |
| D04 | Webhook SSRF validation is scheme-always, IP-restrictions in hardened mode | Admin-configured URLs are trusted; local testing needs localhost webhooks, so private-IP blocking is gated on `hardened: true` | ✅ Accepted | 2026-07-09 |

## Open decisions with recommendations

| ID | Question | Options | Recommended | Rationale | Owner |
|----|----------|---------|-------------|-----------|-------|
| OD01 | How should release artifacts be signed? | (a) cosign keyless (b) GPG checksums (c) GitHub attestations | Start with SHA256 checksums + GitHub attestations; add cosign later | Low friction, no key management, verifiable | AI Agent |
| OD02 | Should the dashboard split to a separate admin port by default? | (a) Keep single port (b) Optional `admin_port` (c) Mandatory split | Add optional `server.admin_port` for high-assurance deployments; keep single-port default | Preserves local simplicity; enables stronger isolation | AI Agent |
| OD03 | Should audit logs support an append-only / external sink? | (a) SQLite + monotonic IDs (b) External syslog/webhook sink (c) Hash-chain integrity | Phase 1: monotonic IDs + timestamps; Phase 2: optional external sink | Good enough for v1.0; external sink is enterprise nice-to-have | AI Agent |
| OD04 | Which container image scanner should we use? | (a) Trivy (b) Grype (c) Both | Trivy in CI; optionally grype locally | Trivy has good SARIF/GitHub integration | AI Agent |
| OD05 | Should `gitleaks` block commits in CI? | (a) Fail CI on leak (b) Report only | Fail CI on leak in `main`/`dev` PRs; allow reports in feature branches | Prevent secrets reaching default branches | AI Agent |
| OD06 | How do we respond if a secret is found in history? | (a) Rotate + accept risk (b) Rewrite history if still private (c) Both | Rotate secret immediately; rewrite history only if repo is still private; document in `ROLLBACK_PLAN.md` | Minimize exposure without breaking public history | AI Agent |
