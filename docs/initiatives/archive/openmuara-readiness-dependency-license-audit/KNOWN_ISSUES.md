> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Reviewed

---

## Preliminary findings (to validate during execution)

| ID | Finding | Area | Severity | Status | Recommended fix | Notes |
|----|---------|------|----------|--------|-----------------|-------|
| F01 | `.github/dependabot.yml` does not monitor npm ecosystems | Supply-chain automation | Medium | ✅ Closed | Add `package-ecosystem: npm` entries for `/web/dashboard` and `/website` | Added in `.github/dependabot.yml` |
| F02 | No `LICENSE_MATRIX.md` or `DEPENDENCIES.md` exists | License transparency | Medium | ✅ Closed | Generate `LICENSE_MATRIX.md` listing all production deps, versions, SPDX licenses, and compatibility rationale | `LICENSE_MATRIX.md` verified and populated |
| F03 | No Go license-scanning tool in CI | License compliance | Medium | ✅ Closed | Add `go-licenses check ./...` to CI, or equivalent | `scripts/check-licenses.sh` enforced in CI |
| F04 | No `go mod tidy` / `go mod verify` enforcement in CI | Dependency hygiene | Low | ✅ Closed | Add a CI step that fails if `go mod tidy` produces a diff | Enforced in new `dependency-license` CI job |
| F05 | `website` has 21 Docusaurus transitive build-time vulnerabilities | npm / security | High (transitive) | ⚠️ Accepted | Monitor Docusaurus releases; update when a clean version is available | Vulnerabilities are in build-time tooling (`serialize-javascript`, `uuid` via `webpack-dev-server`); tracked in security audit |
| F06 | No npm SBOM attached to releases | Supply-chain transparency | Low | ✅ Closed | Generate SPDX JSON SBOMs for `web/dashboard` and `website` and attach to releases alongside Go SBOM | Implemented in `.github/workflows/release.yml` |
| F07 | No documented dependency update policy | Governance | Low | ✅ Closed | Document in `CONTRIBUTING.md` or `docs/operations.md` how often deps are reviewed, who approves major updates, and how exceptions are recorded | `DEPENDENCY_UPDATE_POLICY.md` created and linked |

## Categories to scan during execution

- [x] Go direct dependencies (`go.mod` require block)
- [x] Go indirect dependencies (`go.mod` require block)
- [x] npm production dependencies in `web/dashboard`
- [x] npm production dependencies in `website`
- [x] npm devDependencies in both packages (for completeness, even if not distributed)
- [x] GitHub Actions used in workflows
- [x] Docker base images and installed Alpine packages

## Closed findings

| ID | Resolution | Fixed in |
|---|---|---|
| F01 | Added npm ecosystems to Dependabot | `.github/dependabot.yml` |
| F02 | Populated and verified `LICENSE_MATRIX.md` | `docs/initiatives/openmuara-readiness-dependency-license-audit/LICENSE_MATRIX.md` |
| F03 | Added `go-licenses` check wrapper and CI job | `scripts/check-licenses.sh`, `.github/workflows/ci.yml` |
| F04 | Enforced `go mod tidy`/`go mod verify` with diff check in CI | `.github/workflows/ci.yml` |
| F06 | Generate and attach npm SBOMs to releases | `.github/workflows/release.yml` |
| F07 | Documented dependency update policy | `docs/initiatives/openmuara-readiness-dependency-license-audit/DEPENDENCY_UPDATE_POLICY.md` |

## Accepted risks

| ID | Risk | Rationale | Owner | Review date |
|---|---|---|---|---|
| F05 | Docusaurus build-time vulnerabilities | Vulnerabilities exist only in build-time tooling, not in shipped runtime artifacts; updating would require a Docusaurus major/minor release that is not yet available | TBD | Next dependency review |
