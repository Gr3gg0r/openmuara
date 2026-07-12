> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Known Issues

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ✅ Delivered on dev

---

## Methodology

These findings were identified by reviewing the current CI/release artifacts in the repository. Each issue maps to a recommendation in `RECOMMENDATIONS.md` and a milestone in `EXECUTION_PLAN.md`. Severity is based on security impact, user trust, and release reliability.

---

## Security findings

### KI-001 — Release artifacts lack provenance attestations

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Binaries, checksums, and SBOMs are uploaded to GitHub Releases with no supply-chain attestation.
- **Risk:** Users cannot cryptographically verify that artifacts were built from the tagged source.
- **Severity:** High
- **Remediation:** Generate SLSA Level 3 provenance with `slsa-framework/slsa-github-generator`.
- **Tracked in:** `RECOMMENDATIONS.md` G1.1, `EXECUTION_PLAN.md` M1

### KI-002 — Release artifacts are not signed

- **Location:** `.github/workflows/release.yml`, `dist/checksums.txt`
- **Current behavior:** `checksums.txt` is generated but not signed.
- **Risk:** An attacker could replace the archive or checksum file without detection.
- **Severity:** High
- **Remediation:** Sign `checksums.txt` with Sigstore cosign keyless signing.
- **Tracked in:** `RECOMMENDATIONS.md` G1.2, `EXECUTION_PLAN.md` M1

### KI-003 — Container image is not signed

- **Location:** `.github/workflows/release.yml` (build-push-action)
- **Current behavior:** GHCR image is pushed without a signature.
- **Risk:** Users cannot verify image authenticity.
- **Severity:** High
- **Remediation:** Capture image digest and sign with cosign; document verification.
- **Tracked in:** `RECOMMENDATIONS.md` G1.3, `EXECUTION_PLAN.md` M2

### KI-004 — Install script does not verify downloads

- **Location:** `scripts/install.sh`
- **Current behavior:** Archive is downloaded and extracted without checksum or signature verification.
- **Risk:** Users executing `curl | bash` install a potentially tampered binary.
- **Severity:** High
- **Remediation:** Verify SHA256 hash against `checksums.txt` and cosign signature; add `SKIP_VERIFY=1` escape hatch.
- **Tracked in:** `RECOMMENDATIONS.md` G4.1, `EXECUTION_PLAN.md` M4

### KI-005 — Trivy scan does not fail on critical vulnerabilities

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Trivy uploads SARIF but does not block release on CRITICAL/HIGH findings.
- **Risk:** Vulnerable images can be released.
- **Severity:** Medium
- **Remediation:** Add `severity: 'CRITICAL,HIGH'` and `exit-code: '1'` while keeping SARIF upload.
- **Tracked in:** `RECOMMENDATIONS.md` G5.4, `EXECUTION_PLAN.md` M1

### KI-006 — SBOM tool installed without checksum verification

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** `go install github.com/anchore/syft/cmd/syft@v1.46.0` runs at release time.
- **Risk:** Compromised Syft binary could alter SBOM output.
- **Severity:** Low
- **Remediation:** Pin with `anchore/sbom-action` SHA or verify checksum after install.
- **Tracked in:** `RECOMMENDATIONS.md` G5.3, `EXECUTION_PLAN.md` M1

---

## Reliability findings

### KI-007 — Docker image does not include the built dashboard

- **Location:** `Dockerfile`
- **Current behavior:** Dockerfile builds only the Go binary; dashboard built in CI is not copied into the image.
- **Risk:** `/_admin` may serve a placeholder or fail to load assets when running from the container.
- **Severity:** High
- **Remediation:** Copy `internal/ui/dashboard-dist/` into the image; fallback build if missing.
- **Tracked in:** `RECOMMENDATIONS.md` G3.1, `EXECUTION_PLAN.md` M3

### KI-008 — Docker Compose healthcheck uses missing `wget` binary

- **Location:** `docker-compose.yml`
- **Current behavior:** Healthcheck command is `["CMD", "wget", ...]` but runtime image only installs `ca-certificates`.
- **Risk:** Container is permanently reported as unhealthy.
- **Severity:** High
- **Remediation:** Implement `muara health` CLI command and use it in Dockerfile/Compose healthchecks.
- **Tracked in:** `RECOMMENDATIONS.md` G3.2, `EXECUTION_PLAN.md` M3

### KI-009 — No version alignment gate

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** A tag can be pushed that does not match the `VERSION` file or the embedded `internal/version.Version`.
- **Risk:** Release metadata is inconsistent and confusing.
- **Severity:** Medium
- **Remediation:** Add a `verify-version` job that fails on mismatch.
- **Tracked in:** `RECOMMENDATIONS.md` G2.1, `EXECUTION_PLAN.md` M1, M5

### KI-010 — No prerelease flow

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Any `v*` tag pushes `latest` and creates a full release.
- **Risk:** Prerelease tags (e.g., `v1.1.0-rc.1`) incorrectly promote `latest`.
- **Severity:** Medium
- **Remediation:** Detect semver prerelease and skip `latest` tag; mark GitHub Release as prerelease.
- **Tracked in:** `RECOMMENDATIONS.md` G2.3, `EXECUTION_PLAN.md` M5

### KI-011 — No post-release validation

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Workflow ends after release creation; no job tests the published tarball or image.
- **Risk:** Broken artifacts are discovered by users instead of CI.
- **Severity:** Medium
- **Remediation:** Add `release-smoke` and `release-container-smoke` jobs.
- **Tracked in:** `RECOMMENDATIONS.md` G2.4, `EXECUTION_PLAN.md` M5

### KI-012 — No Docker build validation on PRs

- **Location:** `.github/workflows/ci.yml`
- **Current behavior:** `docker build` is only exercised during release.
- **Risk:** Dockerfile regressions are caught late.
- **Severity:** Medium
- **Remediation:** Add a `docker-build` job to `ci.yml`.
- **Tracked in:** `RECOMMENDATIONS.md` G5.1, `EXECUTION_PLAN.md` M3

---

## Completeness / polish findings

### KI-013 — Container image lacks OCI labels

- **Location:** `Dockerfile`, `.github/workflows/release.yml`
- **Current behavior:** No `org.opencontainers.image.*` labels are set.
- **Risk:** Image provenance and tooling metadata are missing.
- **Severity:** Low
- **Remediation:** Add OCI labels via build-push-action.
- **Tracked in:** `RECOMMENDATIONS.md` G3.4, `EXECUTION_PLAN.md` M2

### KI-014 — No local CI validation guidance

- **Location:** `docs/contributing.md`
- **Current behavior:** Contributors must push to GitHub to test workflow changes.
- **Risk:** Slow feedback and noisy git history.
- **Severity:** Low
- **Remediation:** Document `act` usage and fork-based validation.
- **Tracked in:** `RECOMMENDATIONS.md` G6.3, `EXECUTION_PLAN.md` M6

### KI-015 — Release notes use full CHANGELOG.md

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** `body_path: CHANGELOG.md` publishes the entire file.
- **Risk:** Unreleased or unrelated changes appear in release notes.
- **Severity:** Low
- **Remediation:** Extract the matching `## [X.Y.Z]` section automatically.
- **Tracked in:** `RECOMMENDATIONS.md` G2.2, `EXECUTION_PLAN.md` M1

---

## Expanded findings (gold-standard gaps)

### KI-016 — No OpenSSF Scorecard tracking

- **Location:** `.github/workflows/`
- **Current behavior:** Scorecard is not run in CI; score is unknown.
- **Risk:** Security posture regresses silently; OSS users cannot see trust signals.
- **Severity:** Medium
- **Remediation:** Add `.github/workflows/scorecard.yml` and a README badge; target ≥ 8.5.
- **Tracked in:** `RECOMMENDATIONS.md` G7.1, `EXECUTION_PLAN.md` M1/M6

### KI-017 — Branch protection rules are not documented

- **Location:** GitHub repository settings / `AGENTS.md`
- **Current behavior:** No documented required status checks or rulesets for `main`/`dev`.
- **Risk:** Broken code can be merged before CI passes; force-pushes to protected branches are possible.
- **Severity:** Medium
- **Remediation:** Document required checks and ruleset configuration in `AGENTS.md` or `docs/contributing.md`.
- **Tracked in:** `RECOMMENDATIONS.md` G7.2, `EXECUTION_PLAN.md` M6

### KI-018 — Workflow token permissions are overly broad

- **Location:** `.github/workflows/ci.yml`, `.github/workflows/release.yml`
- **Current behavior:** Top-level permissions grant `contents: read`/`actions: write` (CI) and `contents: write`, `packages: write`, `security-events: write` (release) to every job.
- **Risk:** Compromised action or job can abuse unnecessary permissions.
- **Severity:** High
- **Remediation:** Set top-level `permissions: {}` and grant minimal job-level permissions.
- **Tracked in:** `RECOMMENDATIONS.md` G8.1, `EXECUTION_PLAN.md` M1

### KI-019 — Third-party actions are not all pinned to full SHAs

- **Location:** `.github/workflows/*.yml`
- **Current behavior:** Some workflow snippets reference mutable tags; although current files pin SHAs, the policy is not enforced.
- **Risk:** Mutable tags can be retagged to malicious commits.
- **Severity:** High
- **Remediation:** Pin every third-party action to a full SHA and add a CI check or Dependabot grouping rule.
- **Tracked in:** `RECOMMENDATIONS.md` G8.2, `EXECUTION_PLAN.md` M1

### KI-020 — Builds are not fully reproducible

- **Location:** `Taskfile.yml`, `.github/workflows/release.yml`
- **Current behavior:** `-trimpath` is used in release but not in local `task build`; `-buildvcs=false` is not used anywhere.
- **Risk:** Build metadata leaks paths/version-control info; binaries differ between builds.
- **Severity:** Medium
- **Remediation:** Use `-trimpath -buildvcs=false` on all build paths and record build metadata in SBOMs.
- **Tracked in:** `RECOMMENDATIONS.md` G8.3, `EXECUTION_PLAN.md` M1

### KI-021 — No GitHub artifact attestations

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Release artifacts and images have no GitHub-native attestation.
- **Risk:** Users must rely only on external cosign/SLSA tools for verification.
- **Severity:** High
- **Remediation:** Add `actions/attest-build-provenance` for tarballs and container images.
- **Tracked in:** `RECOMMENDATIONS.md` G9.1, `EXECUTION_PLAN.md` M1

### KI-022 — Container runtime hardening is incomplete

- **Location:** `Dockerfile`, `docker-compose.yml`
- **Current behavior:** Non-root user exists, but root filesystem is writable and Linux capabilities are not dropped.
- **Risk:** Container escape or post-exploitation is easier.
- **Severity:** Medium
- **Remediation:** Add read-only rootfs, drop all capabilities, and provide a distroless variant.
- **Tracked in:** `RECOMMENDATIONS.md` G10.2, `EXECUTION_PLAN.md` M3

### KI-023 — No VEX or CVE exception process

- **Location:** `docs/security.md`, `.github/workflows/release.yml`
- **Current behavior:** Trivy gate may block on unfixable upstream CVEs with no documented exception path.
- **Risk:** Releases are blocked or vulnerable images are released silently.
- **Severity:** Medium
- **Remediation:** Create `docs/security/cve-exceptions.md` and a machine-readable VEX file; gate reads it before failing.
- **Tracked in:** `RECOMMENDATIONS.md` G10.3, `EXECUTION_PLAN.md` M1

### KI-024 — Release can be triggered accidentally by tag push only

- **Location:** `.github/workflows/release.yml`
- **Current behavior:** Any `v*` tag push creates a release; no manual approval gate.
- **Risk:** Mistags or compromised tokens can publish releases.
- **Severity:** Medium
- **Remediation:** Add `workflow_dispatch` as primary trigger and require a maintainer to approve tag-based releases via environment protection rules.
- **Tracked in:** `RECOMMENDATIONS.md` G11.1, `EXECUTION_PLAN.md` M5

### KI-025 — No release failure notification or scheduled validation

- **Location:** `.github/workflows/ci.yml`, `.github/workflows/release.yml`
- **Current behavior:** Failed releases are not notified; no scheduled CI run catches upstream breakage.
- **Risk:** Silent failures and delayed detection of broken dependencies.
- **Severity:** Low
- **Remediation:** Add release failure notification step and a weekly `schedule` trigger in CI.
- **Tracked in:** `RECOMMENDATIONS.md` G11.2, G12.1, `EXECUTION_PLAN.md` M5/M6

---

## Summary

| Severity | Count |
|----------|-------|
| High | 9 |
| Medium | 11 |
| Low | 5 |
| **Total** | **25** |

> Note: Counts reflect the expanded gold-standard audit. The original 15 findings are retained; KI-016–KI-025 add 10 new gaps.

All findings are tracked in `TRACKING.md` and assigned to milestones in `EXECUTION_PLAN.md`.
