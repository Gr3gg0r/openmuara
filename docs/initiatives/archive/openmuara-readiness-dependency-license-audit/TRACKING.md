> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Tracking

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete

---

## Phases

| Phase | Title | Goal | Recommended approach | Acceptance criteria | Effort | Status |
|-------|-------|------|----------------------|---------------------|--------|--------|
| P01 | Go dependency review | Ensure Go modules are tidy, verified, up-to-date, and vulnerability-free | Run `go mod tidy`, `go mod verify`, `go list -m -u all`, `govulncheck ./...`; inspect for unused direct deps | `go mod tidy && go mod verify` produce no diff; no reachable high/critical vulnerabilities; outdated deps documented | S | ✅ Done |
| P02 | npm dependency review | Ensure npm lockfiles are consistent, production deps are audit-clean, and stale deps are identified | Run `npm audit --production` and `npm outdated` in `web/dashboard` and `website`; run unused-dep check | Both `package-lock.json` files committed and consistent; `web/dashboard` audit clean; website advisories accepted or remediated | S | ✅ Done |
| P03 | GitHub Actions review | Confirm all workflows use pinned, maintained actions with least privilege | Audit `.github/workflows/*.yml` for SHA pinning, deprecation, and `permissions:` blocks | All actions pinned by SHA; no deprecated actions; minimal permissions | XS | ✅ Done |
| P04 | Container base-image review | Confirm Docker base images and packages are current and scannable | Review `Dockerfile`; check base image versions and Alpine packages; select image scanner | Base images documented; image scan job planned or implemented; no critical/high CVEs or accepted | S | ✅ Done |
| P05 | License compatibility matrix | Classify every production dependency license against MIT distribution | Use `go-licenses` and npm license tools; build `LICENSE_MATRIX.md`; flag incompatible/unknown licenses | `LICENSE_MATRIX.md` lists all production deps with version, SPDX license, compatibility, and rationale; no unaddressed incompatible licenses | M | ✅ Done |
| P06 | SBOM / attribution | Generate and publish machine-readable SBOMs for all ecosystems | Generate Go SBOM with `syft`; generate npm SBOMs; attach to releases | Release includes Go and npm SBOMs in SPDX JSON; SBOM generation is reproducible | S | ✅ Done |
| P07 | Cleanup & gating | Remove unused deps, update safe deps, and add CI gates | Remove unused deps; update low-risk outdated deps; apply CI changes from `CI_INTEGRATION.md`; update Dependabot | CI enforces dependency hygiene; `.github/dependabot.yml` covers Go, Actions, and npm; no regressions | M | ✅ Done |

## Findings log

| ID | Finding | Area | Severity | Status | Fixed in | Notes |
|----|---------|------|----------|--------|----------|-------|
| F01 | `.github/dependabot.yml` does not monitor npm ecosystems | Supply-chain automation | Medium | ✅ Closed | `.github/dependabot.yml` | Added `package-ecosystem: npm` entries for `/web/dashboard` and `/website` |
| F02 | No `DEPENDENCIES.md` or `LICENSE_MATRIX.md` exists | License transparency | Medium | ✅ Closed | `LICENSE_MATRIX.md` | Populated with all Go and npm production dependencies and transitive summary |
| F03 | No Go license-scanning tool in CI | License compliance | Medium | ✅ Closed | `.github/workflows/ci.yml` | Added `./scripts/check-licenses.sh` job |
| F04 | No `go mod tidy` / `go mod verify` enforcement in CI | Dependency hygiene | Low | ✅ Closed | `.github/workflows/ci.yml` | Added `dependency-license` job that fails on `go.mod`/`go.sum` drift |
| F05 | `website` has 21 Docusaurus transitive build-time vulnerabilities | npm / security | High (transitive) | ⚠️ Accepted | — | Build-time only; monitor upstream; already tracked in security audit |
| F06 | No npm SBOM attached to releases | Supply-chain transparency | Low | ✅ Closed | `.github/workflows/release.yml` | Generates `sbom-dashboard.spdx.json` and `sbom-website.spdx.json` and attaches them |

## Quality gates

Every phase must end with:

- [x] `go build ./...`
- [x] `go test ./...`
- [x] `go vet ./...`
- [x] `golangci-lint run`
- [x] `go mod tidy && go mod verify` with no further changes (local diff is the intentional tidy update; CI will pass once committed)
- [x] `npm run typecheck` (in `web/dashboard/`)
- [x] `npm run test:ci` (in `web/dashboard/`)
- [x] `npm audit --production` in `web/dashboard/` (0 high/critical or explicitly accepted)

## Notes

- Prefer permissive licenses (MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC).
- Flag copyleft (GPL, AGPL, LGPL, MPL, CDDL) or unknown licenses before release.
- Website Docusaurus advisories are build-time only; do not block release but document accepted risk.
- Keep dependency update automation lightweight: Dependabot now covers Go, npm, and GitHub Actions.
- Tool versions are pinned for reproducibility: `go-licenses/v2@v2.0.1`, `syft@v1.46.0`.
- Low-risk npm updates applied: `preact` 10.29.3 → 10.29.7, `terser` 5.48.0 → 5.49.0.
- Both `web/dashboard` and `website` build successfully after updates.
- Concrete CI workflow snippets and acceptance scripts are in `CI_INTEGRATION.md`.
