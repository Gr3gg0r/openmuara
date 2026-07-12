> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ⬜ Draft

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|----|------|------------|--------|------------|-------|
| R01 | Copyleft dependency contaminates MIT distribution | Low | High | Review all transitive licenses with `go-licenses` and npm license tools; replace or document incompatible deps | AI Agent |
| R02 | Unknown/unlicensed dependency creates legal uncertainty | Low | High | Fail CI on missing license metadata; manually inspect any dep without a clear SPDX identifier | AI Agent |
| R03 | Abandoned dependency with unpatched CVE | Medium | Medium | Pin and monitor dependencies; use Dependabot; update where safe; document exceptions | AI Agent |
| R04 | Lockfile drift between environments | Medium | Low | Commit lockfiles; enforce `go mod verify` and `npm ci` in CI | AI Agent |
| R05 | Website Docusaurus build-time vulnerabilities misinterpreted as runtime risk | Medium | Medium | Document accepted risk; keep build-time deps out of the served static site; monitor upstream | AI Agent |
| R06 | npm packages not monitored by Dependabot | Medium | Medium | Add npm ecosystems to `.github/dependabot.yml` | AI Agent |
| R07 | Container base image CVEs reach releases | Medium | High | Scan images in CI with `trivy`/`grype`; pin base image digests where feasible | AI Agent |
| R08 | SBOM is incomplete or non-reproducible | Low | Medium | Use standard tools (`syft`, `npm sbom`) and attach artifacts to releases; verify in CI | AI Agent |
| R09 | Dependency update automation creates noisy or breaking PRs | Medium | Low | Limit open PRs per ecosystem; require tests; group minor/patch updates where supported | AI Agent |
| R10 | Major dependency update breaks provider emulation fidelity | Low | High | Pin critical parser/signature libs; test provider conformance after updates | AI Agent |
| R11 | License scan false positives block legitimate deps | Low | Medium | Maintain an allowlist in CI config; record exceptions in `DECISIONS.md` | AI Agent |
| R12 | Users cannot reproduce the dependency state | Low | Medium | Commit `go.sum` and `package-lock.json`; document exact tool versions | AI Agent |
