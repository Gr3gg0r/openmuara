> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Execution Plan

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09
> **Target branch:** `feat/readiness-ci-release-audit`

---

## Goal

Transform the current functional CI/release pipeline into a gold-standard, contributor-trusted, security-hardened release system that meets OpenSSF Scorecard and SLSA Level 3 expectations without over-engineering the day-to-day maintainer workflow.

---

## Definition of Ready

Before execution starts, the following must be true:

- [ ] `AGENTS.md` branch rules reviewed (work on `dev` or `feat/readiness-ci-release-audit`).
- [ ] All existing CI jobs on `dev` are green.
- [ ] `VERSION` file and `CHANGELOG.md` current state is known.
- [ ] Decision register (`DECISIONS.md`) approved by human reviewer.
- [ ] No open P0 security findings in the release pipeline.

---

## Milestones

### M1 — Harden the release workflow (P0 security)

**Objective:** Add provenance, signing, version alignment, workflow hardening, and Scorecard tracking to `.github/workflows/release.yml` and the repository.

**Tasks:**

1. Add a `verify-version` job that fails if `refs/tags/v${VERSION}` ≠ pushed tag.
2. Extract the matching `## [X.Y.Z]` section from `CHANGELOG.md` for the release body.
3. Replace inline `go install syft` with `anchore/sbom-action` pinned to a SHA.
4. Add SLSA provenance generation with `slsa-framework/slsa-github-generator`.
5. Sign `checksums.txt` with cosign keyless signing.
6. Add GitHub artifact attestations for release tarballs and container image.
7. Add Trivy severity gate (`CRITICAL,HIGH`) while keeping SARIF upload.
8. Generate and attach an image SBOM.
9. Set top-level `permissions: {}` in `release.yml` and grant minimal job-level permissions.
10. Pin every third-party action to a full SHA and add a note about update cadence.
11. Add `-buildvcs=false` alongside `-trimpath` for reproducible builds.
12. Add `.github/workflows/scorecard.yml` and a README badge; target ≥ 8.5.
13. Create `docs/security/cve-exceptions.md` and a VEX file; make Trivy gate respect it.

**Acceptance criteria:**

- [ ] Pushing a test tag on a fork produces a GitHub Release with provenance attestation, signed checksums, image SBOM, and GitHub attestation.
- [ ] Trivy fails the build on CRITICAL vulnerabilities unless listed in VEX.
- [ ] Version mismatch fails before any artifact is published.
- [ ] Scorecard action runs on `main` and reports a score ≥ 8.5.
- [ ] No workflow uses broad top-level permissions.

**Definition of Done:**

- [ ] `release.yml` contains `verify-version` and `verify-changelog` jobs.
- [ ] `release.yml` uses `anchore/sbom-action` with a pinned SHA.
- [ ] SLSA provenance attestation is attached to the GitHub Release.
- [ ] `checksums.txt.sig` is attached to the GitHub Release.
- [ ] `actions/attest-build-provenance` attests release tarballs and image.
- [ ] Trivy step includes `severity: 'CRITICAL,HIGH'`, `exit-code: '1'`, and VEX input.
- [ ] All actions in `release.yml` and `ci.yml` are pinned to full SHAs.
- [ ] Top-level workflow permissions are minimal; job-level permissions are explicit.
- [ ] Go build flags include `-trimpath -buildvcs=false`.
- [ ] `.github/workflows/scorecard.yml` exists and passes.
- [ ] `docs/security/cve-exceptions.md` and a VEX file exist.
- [ ] A fork test release for `v0.0.0-test.1` validates all of the above.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Release artifacts signed | 0 | 2 (`checksums.txt.sig`, GitHub attestation) |
| Provenance attestations | 0 | 2 per release (SLSA + GitHub) |
| Image SBOMs | 0 | 1 per release |
| Trivy severity gate | No fail | Fails on CRITICAL/HIGH unless VEX exempt |
| Scorecard score | Unknown | ≥ 8.5 |
| Reproducible build flags | Partial | `-trimpath -buildvcs=false` |

---

### M2 — Sign and verify container images (P0 security)

**Objective:** Ensure every GHCR image is signed and verifiable.

**Tasks:**

1. Capture the built image digest from `build-push-action` outputs.
2. Add a cosign keyless signing step for the digest.
3. Publish the cosign public key/certificate verification command in `docs/install.md`.
4. Add OCI image labels (`org.opencontainers.image.source`, `.revision`, `.version`).

**Acceptance criteria:**

- [ ] `cosign verify ghcr.io/openmuara/openmuara:<tag> --certificate-identity=... --certificate-oidc-issuer=https://token.actions.githubusercontent.com` succeeds.
- [ ] Image manifest includes OCI labels.

**Definition of Done:**

- [ ] `build-push-action` outputs the image digest.
- [ ] A cosign signing step signs `${IMAGE}@${DIGEST}` using keyless OIDC.
- [ ] `docs/install.md` includes the exact `cosign verify` command.
- [ ] Dockerfile/runtime image includes OCI labels for source, revision, and version.
- [ ] Fork test release verifies the image signature successfully.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Container image signed | No | Yes |
| OCI labels present | No | Yes |
| Verifiable by users | No | Yes |

---

### M3 — Embed dashboard and harden container (P1 reliability)

**Objective:** Make the Docker image serve the built dashboard, report healthy, and run with defense-in-depth hardening.

**Tasks:**

1. Implement `muara health` CLI subcommand (`cmd/muara/health.go`).
2. Update `Dockerfile` to copy `internal/ui/dashboard-dist/` from the build context; fallback to a minimal embedded placeholder if missing.
3. Add `HEALTHCHECK` instruction to `Dockerfile` using `muara health`.
4. Update `docker-compose.yml` healthcheck to use `muara health`.
5. Add read-only root filesystem, drop all Linux capabilities, and bind only required volumes in `docker-compose.yml`.
6. Provide a `-distroless` image variant in `release.yml`.
7. Add a `docker-build` job in `ci.yml` that runs on every PR.

**Acceptance criteria:**

- [ ] `docker compose up` starts a healthy container and serves the dashboard at `/_admin`.
- [ ] `muara health` exits 0 when server is healthy and non-zero otherwise.
- [ ] Container runs with read-only rootfs and dropped capabilities.
- [ ] CI `docker-build` job passes on PRs.
- [ ] Distroless image builds and passes smoke test.

**Definition of Done:**

- [ ] `cmd/muara health` subcommand exists and queries `/healthz`.
- [ ] `Dockerfile` copies `internal/ui/dashboard-dist/` and includes a fallback for missing dist.
- [ ] `Dockerfile` includes `HEALTHCHECK CMD ["muara", "health"]`.
- [ ] `docker-compose.yml` uses `["CMD", "muara", "health"]` as the healthcheck test.
- [ ] `docker-compose.yml` sets `read_only: true`, `cap_drop: [ALL]`, and minimal writable volumes.
- [ ] `release.yml` pushes `:<version>-distroless` and `:distroless` tags.
- [ ] `.github/workflows/ci.yml` has a `docker-build` job that builds and health-checks the image.
- [ ] `task release:docker` passes locally and produces a healthy container.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Dashboard in container | No | Yes |
| Container healthcheck | Broken (`wget` missing) | Working (`muara health`) |
| PR Docker validation | No | Yes |
| Read-only rootfs | No | Yes |
| Distroless variant | No | Yes |

---

### M4 — Harden install script (P1 trust)

**Objective:** Make `scripts/install.sh` verify what it downloads.

**Tasks:**

1. Download `checksums.txt` and the archive.
2. Verify SHA256 hash of the archive against `checksums.txt`.
3. If cosign is installed, verify the signature of `checksums.txt`.
4. Add `SKIP_VERIFY=1` escape hatch with a printed warning.
5. Add lowercase repo normalization.
6. Add a CI job that runs `install.sh --dry-run` for latest and a pinned version.

**Acceptance criteria:**

- [ ] `install.sh` fails if the archive hash does not match `checksums.txt`.
- [ ] `SKIP_VERIFY=1` bypasses verification and warns the user.
- [ ] Dry-run CI job passes for linux/amd64, darwin/arm64, and windows/amd64.

**Definition of Done:**

- [ ] `install.sh` downloads `checksums.txt` and verifies the archive hash with `sha256sum`.
- [ ] `install.sh` verifies `checksums.txt.sig` with cosign when cosign is installed.
- [ ] `SKIP_VERIFY=1` prints a warning and bypasses verification.
- [ ] `install.sh` normalizes repo name to lowercase.
- [ ] `.github/workflows/ci.yml` includes an `install-dry-run` job with an OS/arch matrix.
- [ ] A tampered archive is rejected by the install script in CI testing.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Hash verification | No | Yes |
| Signature verification | No | Yes (when cosign available) |
| Escape hatch | No | `SKIP_VERIFY=1` |
| CI dry-run matrix | No | linux/amd64, darwin/arm64, windows/amd64 |

---

### M5 — Prerelease, post-release validation, and release controls (P1 confidence)

**Objective:** Support safe prereleases, validate artifacts after publication, and prevent accidental or silent releases.

**Tasks:**

1. Detect semver prerelease in `release.yml` and skip `latest` tag + mark GitHub Release as prerelease.
2. Add a `release:smoke` job that downloads the release tarball for linux/amd64 and runs `scripts/smoke-test.sh` against it.
3. Add a `release:container-smoke` job that pulls the published image and runs smoke tests inside a container.
4. Add a version-consistency check between `VERSION`, `CHANGELOG.md`, and embedded version.
5. Add `workflow_dispatch` trigger as the primary release mechanism with tag push as fallback.
6. Add a release failure notification step (GitHub issue, Slack/Discord webhook, or maintainer email).
7. Document release retention policy (last 10 stable, last 5 prereleases) in `runbooks/release.md`.

**Acceptance criteria:**

- [ ] `v1.1.0-rc.1` creates a pre-release and pushes only `1.1.0-rc.1` tags.
- [ ] Post-release smoke jobs pass before the release is considered complete.
- [ ] Version mismatch blocks the release.
- [ ] `workflow_dispatch` can cut a release from a chosen branch/tag.
- [ ] Failed releases notify maintainers within the workflow.

**Definition of Done:**

- [ ] `release.yml` detects semver prerelease and skips the `latest` container tag.
- [ ] `release.yml` marks GitHub Releases as prerelease for RC tags.
- [ ] `release.yml` includes a `release-smoke` job that downloads and tests the linux/amd64 tarball.
- [ ] `release.yml` includes a `release-container-smoke` job that pulls and tests the published image.
- [ ] `release.yml` supports `workflow_dispatch` with an input for the target tag/ref.
- [ ] A release failure creates a GitHub issue or sends a webhook notification.
- [ ] `runbooks/release.md` documents retention policy and release trigger choice.
- [ ] A fork test of `v0.0.0-test.1` validates prerelease behavior.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Prerelease support | None | Full (no `latest` promotion) |
| Post-release binary smoke | No | Yes |
| Post-release container smoke | No | Yes |
| Version mismatch blocks release | No | Yes |
| Manual release trigger | No | `workflow_dispatch` |
| Release failure alerting | Silent | Notified |

---

### M6 — Documentation, contributor experience, and continuous validation (P2 polish)

**Objective:** Document the new pipeline, make it reproducible for contributors, and keep it healthy over time.

**Tasks:**

1. Create `runbooks/release.md` with tag, changelog, rollback, verification, and retention steps.
2. Create `docs/install.md` with checksum/signature verification examples.
3. Add `.actrc` and `docs/contributing.md` section for local workflow validation.
4. Document branch protection rules and required status checks in `AGENTS.md` or `docs/contributing.md`.
5. Add CI/release/container/OpenSSF badges to root `README.md`.
6. Add a CI check that validates `CHANGELOG.md` format and that the version section exists.
7. Add a weekly `schedule` trigger to `ci.yml` for continuous validation.
8. Configure Dependabot grouping for minor/patch Go and npm updates.

**Acceptance criteria:**

- [ ] A new maintainer can cut a release by following `runbooks/release.md`.
- [ ] A user can verify a downloaded binary by following `docs/install.md`.
- [ ] Badges are visible on the repository landing page.
- [ ] Branch protection rules and required checks are documented.
- [ ] CI runs weekly even without commits.

**Definition of Done:**

- [ ] `runbooks/release.md` exists with tag, changelog, rollback, verification, and retention steps.
- [ ] `docs/install.md` exists with checksum and signature verification examples.
- [ ] `.actrc` exists and `docs/contributing.md` has a local CI validation section.
- [ ] Branch protection rules and required status checks are documented.
- [ ] Root `README.md` includes CI, release, container, license, and OpenSSF Scorecard badges.
- [ ] `.github/workflows/ci.yml` includes a `changelog-check` job and a weekly `schedule` trigger.
- [ ] `.github/dependabot.yml` groups minor/patch Go and npm updates.
- [ ] Documentation links are validated in CI.

**Metrics:**

| Metric | Before | After |
|--------|--------|-------|
| Release runbook | None | Yes |
| Install verification guide | None | Yes |
| Local CI validation guide | None | Yes (`act` + fork) |
| Branch protection docs | None | Yes |
| README badges | None | Yes |
| Changelog CI check | None | Yes |
| Scheduled CI | None | Weekly |
| Dependabot grouping | None | Minor/patch groups |

---

## RACI

| Activity | Responsible | Accountable | Consulted | Informed |
|----------|-------------|-------------|-----------|----------|
| Workflow changes | AI Agent | Human reviewer | — | Contributors |
| Signing key/cosign setup | AI Agent | Human reviewer | Security-aware reviewer | Users |
| VERSION/CHANGELOG bumps | Human reviewer | Human reviewer | AI Agent | Contributors |
| Release cut | Human reviewer | Human reviewer | AI Agent | Users |
| Post-release validation | AI Agent (CI) | Human reviewer | — | Users |

---

## Timeline estimate

| Milestone | Estimated effort | Dependencies |
|-----------|------------------|--------------|
| M1 | 1.5 days | — |
| M2 | 0.5 day | M1 |
| M3 | 1.5 days | — |
| M4 | 0.5 day | M1 |
| M5 | 1 day | M1, M2, M4 |
| M6 | 0.5 day | M1–M5 |
| **Total** | **~5.5 days** | — |

---

## Milestone / decision dependency matrix

| Milestone | Depends on decisions | Blocks | Why |
|-----------|---------------------|--------|-----|
| M1 | D1, D2, D7, D8, D11, D12, D13, D14 | M2, M4, M5 | Signing, provenance, permissions, and version gates must be in place first. |
| M2 | D2 | M5 | Image signing must work before post-release container smoke tests can verify it. |
| M3 | D3, D4, D9 | M5 | Dashboard embedding and `muara health` must be stable before container smoke tests. |
| M4 | D2, D6 | M5 | Hardened install script must be ready before post-release install tests. |
| M5 | D5, D7, D15 | M6 | Prerelease behavior and release controls must be validated before documenting them. |
| M6 | D10 | — | Documentation depends on all implementation choices being finalized. |

---

## Quality gates at every milestone

Each milestone must end with:

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] `golangci-lint run`
- [ ] `task quality` (or equivalent) passes locally.
- [ ] Relevant workflow file validates with `actionlint` if available.
- [ ] Documentation updated to reflect changes.

---

## Rollback plan

If a release causes breakage:

1. Delete the GitHub Release and tag (requires admin).
2. Re-tag the last known-good commit and re-run the release workflow if the artifacts themselves are clean.
3. If the container image is broken, retag `latest` to the previous digest using `crane` or `docker pull/tag/push`.
4. Update `CHANGELOG.md` with a regression note under `[Unreleased]`.

---

## Cross-initiative dependencies

| Initiative | Relationship | Action required |
|------------|--------------|-----------------|
| `openmuara-readiness-security-audit` | Overlaps on signing, SBOM, VEX, Scorecard | Coordinate so signing keys, SBOM generation, and VEX format are consistent |
| `openmuara-readiness-docs-completeness-audit` | Produces `docs/install.md` and `runbooks/release.md` | Ensure docs audit covers new CI/release documentation |
| `openmuara-readiness-a11y-usability-audit` | Dashboard UI is embedded in Docker image | Verify dashboard build is accessible before embedding |
| `openmuara-readiness-coverage-audit` | Smoke tests and release validation depend on test coverage | Ensure coverage floors do not regress when adding release-only code paths |
| `openmuara-readiness-dependency-license-audit` | SBOMs and license checks feed into release artifacts | Align license scanning output with SBOM formats and release bundle |

---

## Success metrics and KPIs

| KPI | Baseline | Target | Measurement |
|-----|----------|--------|-------------|
| Release artifact verification | Not possible | 100% of releases signed + provenance + attestation | Check every GitHub Release |
| Container image verification | Not possible | 100% of images signed + attested | `cosign verify` + `gh attestation verify` on every release tag |
| Install script verification | No verification | Hash verified; signature verified when cosign available | CI install-dry-run + tamper test |
| PR Docker validation | None | Every PR | `docker-build` job in CI |
| Release confidence | Smoke runs on source only | Smoke runs on published binary and image | `release-smoke` and `release-container-smoke` jobs |
| OpenSSF Scorecard | Unknown | ≥ 8.5 | Scorecard action runs after implementation |

---

## Before / after summary

| Area | Before | After |
|------|--------|-------|
| **Security** | Unsigned binaries/images, no provenance | Signed artifacts, SLSA provenance, signed images, GitHub attestations |
| **Reliability** | Missing dashboard in image, broken healthcheck | Dashboard embedded, `muara health` works, hardened runtime |
| **Trust** | `curl \| bash` with no verification | Verified install by default with escape hatch |
| **Quality gates** | Docker/build not tested on PRs | Docker build and install dry-run on every PR |
| **Release discipline** | Manual version/changelog alignment | Automated gates for version and changelog, `workflow_dispatch` controls |
| **Documentation** | No release/install runbooks | Complete runbooks, verification guides, badges, branch protection docs |
| **Observability** | No Scorecard tracking | Scorecard action + badge, release failure notifications |
| **Supply-chain hygiene** | Broad permissions, mixed pinning | Minimal permissions, full SHA pinning, reproducible builds |
