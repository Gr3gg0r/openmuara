> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dependency Update Policy

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ⬜ Draft

---

This policy defines how OpenMuara adds, updates, and removes dependencies. It applies to Go modules, npm packages, GitHub Actions, and container base images.

## Principles

1. **Minimize dependencies.** Only add a dependency when the cost of maintaining it is lower than the cost of implementing the capability internally.
2. **Prefer permissive licenses.** All production dependencies must be compatible with the MIT License.
3. **Pin and verify.** Use lockfiles (`go.sum`, `package-lock.json`) and pinned actions (SHA) so the dependency graph is reproducible.
4. **Automate monitoring.** Use Dependabot to detect updates and vulnerabilities.
5. **Review before merge.** Every dependency change must pass CI, including license and vulnerability scans.

## Adding a new dependency

Before adding a new production dependency:

- [ ] Confirm the dependency is necessary and no suitable standard-library alternative exists.
- [ ] Check the license against the compatibility rules in `README.md`.
- [ ] Verify the project is maintained (recent commits/releases, responsive maintainers).
- [ ] Run `go mod tidy` or `npm install` and commit the lockfile change.
- [ ] Ensure the change passes `go-licenses check ./...` (Go) or equivalent npm license review.
- [ ] Ensure the change passes `govulncheck ./...` and `npm audit --production`.
- [ ] Update `LICENSE_MATRIX.md` if it is a production dependency.

## Updating dependencies

- **Patch and minor updates** (bug fixes, features): Dependabot may open PRs automatically. Merge after CI passes.
- **Major updates** (breaking changes): Require a human review. Verify provider emulation tests and conformance tests still pass.
- **Security updates**: Treat as high priority. Apply patches as soon as CI passes; if a safe patch is not available, document the accepted risk in `KNOWN_ISSUES.md`.

## Removing dependencies

- Remove unused dependencies identified by `go mod tidy`, `depcheck`, or manual review.
- Update `LICENSE_MATRIX.md` and SBOM generation scripts after removal.

## Container images

- Pin base images by digest when feasible; let Dependabot propose digest updates.
- Scan release images for CVEs before publishing.
- Document any accepted image CVEs in `KNOWN_ISSUES.md`.

## Exceptions

Any exception to this policy (e.g., a copyleft dependency, an unpatched vulnerability, a non-reproducible lockfile change) must be:

1. Proposed in a PR with clear rationale.
2. Approved by a maintainer.
3. Recorded in `DECISIONS.md` and `RISKS.md`.

## Review cadence

- **Weekly:** Dependabot scans and PRs.
- **Monthly:** Manual review of `npm outdated` and `go list -m -u all` output.
- **Before every release:** Regenerate `LICENSE_MATRIX.md` and SBOMs; re-run vulnerability scans.
