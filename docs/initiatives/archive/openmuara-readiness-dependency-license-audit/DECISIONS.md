> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Decision Log

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete

---

## Accepted decisions

| ID | Decision | Context | Status | Date |
|----|----------|---------|--------|------|
| D01 | OpenMuara is distributed under MIT | MIT is permissive, widely understood, and compatible with most OSS dependencies | ✅ Accepted (pre-existing) | 2026-07-08 |
| D02 | Dependabot is the primary dependency update automation | Already in use for Go and GitHub Actions; extended to npm rather than introduce Renovate | ✅ Accepted | 2026-07-09 |
| D03 | Website Docusaurus build-time vulnerabilities are accepted risk | Vulnerable packages are not executed in the served static site; updating Docusaurus is tracked separately | ✅ Accepted | 2026-07-09 |
| D04 | All currently known production dependencies use MIT-compatible licenses | Full verification of `go.mod`, `web/dashboard/package.json`, and `website/package.json` confirms only MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC, MIT-0, MPL-2.0, and CC-BY-4.0 (build-time only) licenses; no incompatible licenses found | ✅ Accepted | 2026-07-09 |
| D05 | Per-ecosystem SBOMs are preferred over a single combined SBOM | Easier to audit and update independently; consumers can select the relevant artifact | ✅ Accepted | 2026-07-09 |
| D06 | `github.com/hashicorp/hcl` (MPL-2.0) is accepted as a transitive dependency | MPL-2.0 is a file-level weak copyleft; using it unmodified in an MIT project is compatible | ✅ Accepted | 2026-07-09 |
| D07 | `modernc.org/mathutil` is treated as BSD-3-Clause | go-licenses cannot detect its license automatically; the LICENSE file was manually inspected and confirmed as BSD-3-Clause | ✅ Accepted | 2026-07-09 |

## Open decisions — resolved

| ID | Question | Resolution | Rationale | Outcome |
|----|----------|------------|-----------|---------|
| OD01 | Which Go license scanner should be used? | `go-licenses` | Open source, SPDX-aware, integrates cleanly with Go modules, no external service required | Implemented in `scripts/check-licenses.sh` and CI |
| OD02 | Should npm devDependencies be included in the license matrix? | Production only for distribution; dev deps reviewed for awareness | MIT compatibility primarily matters for distributed code; dev deps affect contributors and build environment | `LICENSE_MATRIX.md` lists production deps and a transitive license summary |
| OD03 | Should the release SBOM be one file or per-ecosystem? | Per-ecosystem SBOMs | Easier to audit and update independently; consumers can pick the relevant SBOM | Release attaches `sbom.spdx.json`, `sbom-dashboard.spdx.json`, `sbom-website.spdx.json` |
| OD04 | Which container scanner should be used in CI? | Trivy | Strong SARIF/GitHub integration and widely adopted | Trivy scan added to `.github/workflows/release.yml` |
| OD05 | Should base images be pinned by digest? | Optional; pin by digest when Dependabot can update | Improves reproducibility and supply-chain security; requires automation to avoid stale digests | Deferred; current `Dockerfile` uses tags |
| OD06 | How should license exceptions be approved? | Maintainer approves and records in `DECISIONS.md` / `LICENSE_MATRIX.md` | Creates an audit trail for future reviewers | Exceptions recorded in `LICENSE_MATRIX.md` with decision IDs |
