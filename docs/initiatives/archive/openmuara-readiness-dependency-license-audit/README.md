> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit

> **Status:** ✅ Complete | **Started:** 2026-07-08 | **Completed:** 2026-07-09
> **Scope:** Review all Go, npm, GitHub Actions, and container-base dependencies for license compatibility, freshness, vulnerability exposure, and supply-chain hygiene before public release.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/readiness-dependency-license-audit` (to be created when work starts)
> **License Policy:** OpenMuara is distributed under the **MIT License**. Every production dependency must be compatible with MIT distribution, and its license must be reproducibly discoverable.

---

## Why this matters

Publishing as OSS means every transitive dependency becomes part of the distributed work. Incompatible licenses (e.g., GPL, AGPL, proprietary, or unknown) can create legal risk for users and contributors. Outdated or abandoned packages increase vulnerability exposure, and undetected supply-chain compromises can undermine the security work done elsewhere. This initiative makes dependency hygiene a repeatable, verifiable process.

## Initiative structure

```
docs/initiatives/openmuara-readiness-dependency-license-audit/
├── README.md              # This file
├── EXECUTION_PLAN.md      # Timeline, milestones, RACI
├── TRACKING.md            # Central execution tracker
├── KNOWN_ISSUES.md        # Catalog of dependency findings
├── RECOMMENDATIONS.md     # Recommended fixes and priorities
├── RISKS.md               # Risk register
├── DECISIONS.md           # Decision log
├── LICENSE_MATRIX.md      # Generated production dependency license matrix
├── CI_INTEGRATION.md      # Concrete CI workflow snippets and acceptance scripts
├── DEPENDENCY_UPDATE_POLICY.md  # How dependencies are updated and reviewed
├── ROLLBACK_PLAN.md       # How to respond to a bad dependency or license violation
└── HANDOFF.md             # Session continuity
```

## Standards & frameworks mapped

| Standard / Framework | How this initiative uses it |
|---|---|
| **OSI Open Source Definition** | All production dependencies must use an OSI-approved license or an explicitly accepted non-OSI license documented in `DECISIONS.md`. |
| **SPDX** | SBOMs are generated in SPDX JSON format; license identifiers follow SPDX where possible. |
| **OpenSSF Scorecard** | Dependency update tools, pinned dependencies, vulnerability scanning, and license clarity improve Scorecard's "Dependency-Update-Tool", "Vulnerabilities", and "License" checks. |
| **SLSA Level 1–2** | SBOMs and reproducible lockfiles support provenance and supply-chain transparency. |
| **CNCF Best Practices** | Track dependency freshness, avoid abandoned projects, document security contacts, and gate on vulnerability scans. |

## Audit areas

1. **Go dependencies** — `go.mod`, `go.sum`, unused modules, outdated modules, reachable vulnerabilities.
2. **npm dependencies** — `web/dashboard/package.json`, `website/package.json`, lockfile consistency, production-only advisories.
3. **GitHub Actions** — pinned actions, transitive action dependencies, deprecated actions.
4. **Container base images** — `golang:1.26-alpine`, `alpine:3.21` freshness and CVE exposure.
5. **License compatibility** — classify every production dependency license against MIT distribution.
6. **SBOM / attribution** — generate and publish machine-readable dependency lists for releases.
7. **Update cadence & automation** — Dependabot/Renovate coverage, manual review policy, stale-dependency budget.

## License compatibility rules

OpenMuara is distributed under the MIT License. The following rules determine whether a dependency may be used in production:

| Category | Allowed? | Examples | Notes |
|---|---|---|---|
| Permissive | ✅ Yes | MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, ISC | Compatible with MIT distribution |
| Weak copyleft | ⚠️ Review | MPL-2.0, CDDL-1.0, LGPL | Case-by-case; may be acceptable if linking/distribution conditions are met |
| Strong copyleft | ❌ No | GPL-3.0, AGPL-3.0 | Must be replaced or explicitly accepted by maintainer |
| Proprietary | ❌ No | Custom "all rights reserved" | Must be replaced or explicitly accepted by maintainer |
| Unknown | ❌ No | No clear SPDX identifier | Must be resolved before release |

Any exception must be recorded in `DECISIONS.md` with rationale and maintainer approval.

## Success criteria

- `go mod tidy && go mod verify` pass with no changes and are enforced in CI.
- All production Go and npm dependencies use MIT-compatible licenses (or are explicitly accepted and documented).
- No unused top-level dependencies remain in `go.mod` or any `package.json`.
- Lockfiles (`go.sum`, `package-lock.json`) are committed, consistent, and reproducible.
- `npm audit --production` reports 0 high/critical advisories for `web/dashboard`; website Docusaurus build-time advisories are explicitly accepted or remediated.
- `govulncheck` reports 0 reachable high/critical vulnerabilities.
- A `LICENSE_MATRIX.md` (or equivalent) lists every production dependency, its version, SPDX license, and compatibility rationale.
- Release artifacts include SBOMs for Go **and** npm (in addition to the existing Go SBOM).
- Dependabot monitors Go modules, GitHub Actions, **and** npm ecosystems.

## Key metrics

| Metric | Target | How to measure |
|---|---|---|
| Production deps with known-compatible licenses | 100% | `LICENSE_MATRIX.md` review |
| CI `go mod tidy` diff | 0 files | `git diff --exit-code go.mod go.sum` |
| Reachable high/critical vulnerabilities | 0 | `govulncheck ./...` |
| Dashboard production high/critical advisories | 0 | `npm audit --production` in `web/dashboard` |
| Release SBOM coverage | Go + npm | Artifacts attached to release |
| Open dependency-update PRs | ≤10 per ecosystem | Dependabot dashboard |

## Definition of done

This initiative is complete when:

1. All phases in `TRACKING.md` are marked done.
2. `LICENSE_MATRIX.md` is populated and reviewed.
3. CI changes from `CI_INTEGRATION.md` are merged and passing.
4. Dependabot is monitoring Go, GitHub Actions, and npm.
5. All quality gates pass.
6. Remaining accepted risks are documented in `RISKS.md` and `DECISIONS.md`.
7. `HANDOFF.md` is updated with final state and any follow-up work.

## RACI

| Activity | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| Run scans & classify licenses | R | A | C |
| Decide on copyleft/unknown exceptions | C | A | R |
| Approve CI gating changes | R | A | C |
| Approve `LICENSE_MATRIX.md` | C | A | R |
| Final sign-off | C | A | R |

*R = Responsible, A = Accountable, C = Consulted, I = Informed*

## Dependencies on other initiatives

| Initiative | Why it matters for dependency/license audit |
|---|---|
| [Security Audit](../openmuara-readiness-security-audit/README.md) | Vulnerability scanning (`govulncheck`, `npm audit`) and SBOM generation overlap with security controls. |
| [CI & Release Audit](../openmuara-readiness-ci-release-audit/README.md) | Dependency checks must be added to CI/release workflows without duplicating scan jobs. |
| [Docs Completeness Audit](../openmuara-readiness-docs-completeness-audit/README.md) | `LICENSE_MATRIX.md`, `DEPENDENCIES.md`, and README license sections must stay accurate. |

## Copy-paste command reference

```bash
# Go dependency hygiene
go mod tidy
go mod verify
go list -m -u all                                    # outdated modules
govulncheck ./...
go-licenses check ./...                              # license compatibility (install first)
go-licenses csv ./... > licenses-go.csv

# npm dependency hygiene (run in each package directory)
cd web/dashboard && npm audit --production && npm outdated
cd website && npm audit --production && npm outdated

# Unused dependency detection
# Go: review go.mod against imports manually or with go-mod-outdated / goda
# npm: npx depcheck

# SBOM generation
syft dir:. -o spdx-json=sbom.spdx.json
npm sbom --package-lock-only --sbom-format=spdx  # npm 10+
```

See [`TRACKING.md`](TRACKING.md) for the execution plan, [`RISKS.md`](RISKS.md) for the risk register, [`RECOMMENDATIONS.md`](RECOMMENDATIONS.md) for the prioritized action plan, and [`CI_INTEGRATION.md`](CI_INTEGRATION.md) for implementation-ready workflow snippets.
