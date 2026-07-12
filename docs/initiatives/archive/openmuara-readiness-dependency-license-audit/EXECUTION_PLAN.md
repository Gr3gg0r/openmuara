> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Dependency & License Audit Execution Plan

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ⬜ Draft

---

## Milestones

| Milestone | Target | Deliverable |
|---|---|---|
| M1 | Day 1–2 | Inventory: scan all ecosystems, record findings in `KNOWN_ISSUES.md`. |
| M2 | Day 3–4 | Classification: build `LICENSE_MATRIX.md` and compatibility decisions. |
| M3 | Day 5–6 | Remediation: remove unused deps, update safe outdated deps, accept/document unavoidable issues. |
| M4 | Day 7 | Gating: add CI checks, update Dependabot, generate SBOMs, run final quality matrix. |
| M5 | Day 8 | Review & handoff: human review, update docs, merge. |

## RACI by phase

| Phase | AI Agent | Human Reviewer | Maintainer |
|---|---|---|---|
| P01 Go dependency review | R | A | C |
| P02 npm dependency review | R | A | C |
| P03 GitHub Actions review | R | A | I |
| P04 Container base-image review | R | A | C |
| P05 License compatibility matrix | R | A | C |
| P06 SBOM / attribution | R | A | C |
| P07 Cleanup & gating | R | A | C |

## Phase details

### P01 — Go dependency review

- Run `go mod tidy && go mod verify`; fail if any file changes.
- Run `go list -m -u all` to identify outdated modules.
- Run `govulncheck ./...` and triage reachable vulnerabilities.
- Inspect `go.mod` for unused direct dependencies.
- Document findings in `KNOWN_ISSUES.md`.

### P02 — npm dependency review

- Verify `web/dashboard/package-lock.json` and `website/package-lock.json` are committed and consistent.
- Run `npm audit --production` in both directories.
- Run `npm outdated` in both directories.
- Run `npx depcheck` or equivalent to find unused dependencies.
- Document website Docusaurus transitive advisories with accepted-risk rationale.

### P03 — GitHub Actions review

- Verify every `uses:` in `.github/workflows/*.yml` is pinned to a SHA (already done in security audit).
- Identify deprecated or unmaintained actions.
- Ensure `permissions:` blocks are minimal.

### P04 — Container base-image review

- Review `Dockerfile` base images: `golang:1.26-alpine`, `alpine:3.21`.
- Check Alpine package versions (`ca-certificates`, `git`).
- Decide on image scanning tool (`trivy`/`grype`) for CI.

### P05 — License compatibility matrix

- Generate a CSV/matrix of every production Go and npm dependency with: name, version, SPDX license, source URL, compatibility with MIT, notes.
- Flag incompatible or unknown licenses.
- Record explicit decisions for any non-standard licenses in `DECISIONS.md`.
- Publish matrix as `LICENSE_MATRIX.md`.

### P06 — SBOM / attribution

- Generate Go SBOM (`syft` or `go version -m`) and attach to releases (already done).
- Generate npm SBOM for both npm packages.
- Decide whether to ship a combined SBOM or per-ecosystem SBOMs.
- Ensure SBOM format is SPDX JSON.

### P07 — Cleanup & gating

- Remove unused dependencies where safe.
- Update outdated dependencies that pass tests and have no license concerns.
- Add `go mod tidy` / `go mod verify` check to CI.
- Add `go-licenses check` or equivalent to CI.
- Add `npm audit --production` gating to CI for `web/dashboard`.
- Update `.github/dependabot.yml` to include npm ecosystems.
- Add a dependency update policy to `docs/operations.md` or `CONTRIBUTING.md`.

## Phase acceptance scripts

Concrete shell scripts and CI workflow snippets for each phase are in [`CI_INTEGRATION.md`](CI_INTEGRATION.md). A quick summary:

- **P01:** `go mod tidy && go mod verify && git diff --exit-code go.mod go.sum && govulncheck ./...`
- **P02:** `npm ci && npm audit --production && npm outdated` in each npm package.
- **P03:** Verify all `uses:` in `.github/workflows/*.yml` are pinned to SHA.
- **P04:** `docker build -t openmuara:scan . && trivy image openmuara:scan`.
- **P05:** `go-licenses csv ./... > licenses-go.csv` and merge with npm license data into `LICENSE_MATRIX.md`.
- **P06:** `syft dir:. -o spdx-json=sbom.spdx.json` and `npm sbom --package-lock-only --sbom-format=spdx` per package.
- **P07:** Apply CI changes from `CI_INTEGRATION.md` and run the full quality matrix.

## Quality gates

Every phase must end with:

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] `golangci-lint run`
- [ ] `go mod tidy && go mod verify` with no changes
- [ ] `npm run typecheck` (in `web/dashboard/`)
- [ ] `npm run test:ci` (in `web/dashboard/`)
- [ ] `npm audit --production` in `web/dashboard/` (0 high/critical or explicitly accepted)
- [ ] `go-licenses check ./...` passes (after P05)
- [ ] `LICENSE_MATRIX.md` is current (after P05)
