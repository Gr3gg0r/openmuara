> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Recommendations

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete

---

These recommendations are planning-only. They map each audit area to a concrete, industry-standard action. Execute them in priority order once the initiative is approved.

## Priority matrix

| Priority | Area | Recommendation | Effort | Impact | Owner |
|----------|------|----------------|--------|--------|-------|
| P0 | npm automation | Add npm ecosystems to `.github/dependabot.yml` for `web/dashboard` and `website` | Low | High | AI Agent |
| P0 | Go hygiene | Enforce `go mod tidy && go mod verify` in CI and add `govulncheck` gating | Low | High | AI Agent |
| P0 | License matrix | Generate `LICENSE_MATRIX.md` with all production Go and npm dependencies, versions, and SPDX licenses | Medium | High | AI Agent |
| P1 | License gating | Add `go-licenses check ./...` (or `fossa-cli`/`snyk`) to CI to block incompatible licenses | Low | High | AI Agent |
| P1 | npm audit gating | Keep `npm audit --production` gating in CI for `web/dashboard`; document website exceptions | Low | High | AI Agent |
| P1 | SBOM completeness | Generate npm SBOMs and attach Go + npm SBOMs to releases | Low | Medium | AI Agent |
| P1 | Container scanning | Add `trivy` or `grype` image scan to CI release job | Low | High | AI Agent |
| P2 | Unused deps | Run `depcheck` on npm packages and review `go.mod` for unused direct dependencies | Low | Medium | AI Agent |
| P2 | Update policy | Document dependency update cadence and approval process in `CONTRIBUTING.md` | Low | Low | AI Agent |
| P2 | GitHub Actions | Verify all actions are pinned and maintained; replace any deprecated actions | Low | Medium | AI Agent |
| P3 | Reproducible builds | Document and verify reproducible build steps for release binaries and Docker image | Medium | Medium | AI Agent |

## Standards mapping

| Recommendation | OSI / SPDX | OpenSSF Scorecard | SLSA | CNCF |
|---|---|---|---|---|
| Dependabot npm coverage | — | Dependency-Update-Tool | — | Dependency mgmt |
| `go mod tidy` / `govulncheck` gating | — | Vulnerabilities | Build L2 | Security |
| License matrix / `go-licenses` | SPDX | License | — | Compliance |
| npm SBOM | SPDX | — | L1–L2 | Supply chain |
| Container image scanning | — | Vulnerabilities | — | Security |
| Pinned GitHub Actions | — | Pinned-Dependencies | L2–L3 | Supply chain |

## Recommended tool stack

| Purpose | Tool | Where |
|---|---|---|
| Go vulnerability scan | `govulncheck` | CI + local |
| Go license scan | `go-licenses` | CI + local |
| Go outdated deps | `go list -m -u all` | Local / Dependabot |
| npm vulnerability scan | `npm audit --production` | CI + local |
| npm unused deps | `depcheck` | Local |
| npm outdated deps | `npm outdated` / Dependabot | Local / CI |
| SBOM generation | `syft`, `npm sbom` | Release workflow |
| Container image scan | `trivy`, `grype` | CI release job |
| Dependency update automation | GitHub Dependabot | `.github/dependabot.yml` |

## Copy-paste command reference

```bash
# Go
go mod tidy
go mod verify
go list -m -u all
govulncheck ./...
go install github.com/google/go-licenses/v2@v2.0.1
./scripts/check-licenses.sh
go-licenses csv ./... > licenses-go.csv

# npm (use --omit=dev; --production is deprecated)
cd web/dashboard && npm audit --omit=dev && npm outdated
npx depcheck
cd ../../website && npm audit --omit=dev && npm outdated
npx depcheck

# SBOM
go install github.com/anchore/syft/cmd/syft@v1.46.0
syft dir:. -o spdx-json=sbom.spdx.json
# In each npm package directory:
npm sbom --package-lock-only --sbom-format=spdx

# Container
docker build -t openmuara:audit .
docker save openmuara:audit -o openmuara-audit.tar
trivy image --input openmuara-audit.tar --severity HIGH,CRITICAL
grype openmuara:audit
```

## CI integration

Concrete workflow snippets and per-phase acceptance scripts are provided in [`CI_INTEGRATION.md`](CI_INTEGRATION.md). Highlights:

- Add npm ecosystems to `.github/dependabot.yml`.
- Add a `dependency-license` CI job that runs `go mod tidy`, `go mod verify`, `go-licenses check`, and `npm audit --production`.
- Attach Go and npm SBOMs to releases.
- Scan the release Docker image with Trivy before pushing.
- Pin Dockerfile base images by digest for reproducible builds.

## License compatibility rules

See `README.md` for the canonical rules. In short:

- ✅ Allowed: MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC.
- ⚠️ Review: MPL-2.0, CDDL-1.0, LGPL (case-by-case).
- ❌ Blocked: GPL-3.0, AGPL-3.0, proprietary, unknown.

Record any exception in `DECISIONS.md`.

## Related documents

- `DEPENDENCY_UPDATE_POLICY.md` — how to add, update, and remove dependencies.
- `ROLLBACK_PLAN.md` — incident response for dependency/license issues.
- `LICENSE_MATRIX.md` — populated template with known license values.

## What not to do

- Do **not** add license scanning to provider emulation endpoints; it is irrelevant to runtime behavior.
- Do **not** block CI on website Docusaurus build-time advisories unless a reachable exploit path is demonstrated; document the accepted risk instead.
- Do **not** rewrite git history to remove a dependency; normal `go mod tidy` and `npm uninstall` are sufficient.
- Do **not** introduce proprietary or copyleft dependencies without explicit maintainer approval recorded in `DECISIONS.md`.
