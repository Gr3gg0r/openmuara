> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Recommendations

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09
> **Scope:** Gold-standard recommendations for CI, release engineering, container delivery, and install security.

---

## Executive summary

The current CI/release surface is already functional and covers more ground than most early-stage projects: cross-platform binaries, SHA256 checksums, Go/npm SBOMs, Trivy container scanning, multi-arch Docker pushes, and a comprehensive shell-based smoke test. However, several gaps separate it from an OSS-quality, contributor-trusted release pipeline. This document catalogs those gaps and ranks remediations by security impact, contributor trust, and maintenance cost.

### Self-rating

#### Current plan rating: 9.5/10

The plan now covers the full CI/release surface from artifact provenance and signing through container hardening, install verification, release controls, and long-term maintenance. The remaining 0.5 is reserved for execution validation (fork test release, real Scorecard score, and first production release experience), which cannot be verified until implementation begins.

**What earns the 9.5:**

- P0 security: SLSA Level 3, cosign signing, GitHub attestations, install verification.
- P1 reliability: dashboard embedding, `muara health`, distroless variant, version/changelog gates.
- P1 confidence: prerelease flow, post-release smoke, `workflow_dispatch` trigger, failure notifications.
- P1 trust: OpenSSF Scorecard tracking, minimal token permissions, full SHA pinning, reproducible builds.
- P2 polish: VEX process, scheduled CI, Dependabot grouping, branch protection docs, badges.

**What keeps it from 10/10:**

- Execution on a real fork is required to prove cosign OIDC, SLSA generator, and Scorecard ≥ 8.5 work together.
- Package-manager distribution (Homebrew, apt, Chocolatey) is deliberately out of scope.
- Multi-region registry mirrors and automated rollback orchestration are out of scope.

---

## Current-state snapshot

| Component | File(s) | Current capability | Maturity |
|-----------|---------|-------------------|----------|
| CI pipeline | `.github/workflows/ci.yml` | fmt, vet, lint, unit (race), smoke, vuln, gosec, secrets, quality, dependency-license, known-issues sync | Strong |
| Release pipeline | `.github/workflows/release.yml` | Tag-triggered cross-platform builds, SHA256, SBOMs, Trivy scan, GHCR push, GitHub Release | Good |
| Task runner | `Taskfile.yml` | `task quality`, `task release:build`, `task release:docker` | Good |
| Container | `Dockerfile`, `docker-compose.yml` | Multi-stage, non-root, alpine runtime, healthcheck | Good |
| Install script | `scripts/install.sh` | OS/arch detection, prefix override, dry-run | Basic |
| Smoke test | `scripts/smoke-test.sh` | Fawry, Stripe Checkout/PaymentIntents, Billplz, ToyyibPay, iPay88 | Strong |
| Version source | `VERSION` | `1.0.0` | Basic |
| Security policy | `.github/SECURITY.md` | Supported versions, reporting, disclosure | Good |
| Dependency updates | `.github/dependabot.yml` | Go, npm, GitHub Actions | Good |
| Scorecard tracking | `.github/workflows/scorecard.yml` | None | Missing |

---

## Gap analysis and recommendations

### G1 — Artifact provenance and signed releases

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| SLSA provenance | None | No supply-chain attestation | Generate SLSA Level 3 provenance with `slsa-framework/slsa-github-generator` | `.github/workflows/release.yml:106-116` |
| Signed checksums | None | `checksums.txt` is plaintext | Sign `checksums.txt` and release artifacts with Sigstore cosign or a project GPG key | `.github/workflows/release.yml:47-51` |
| Container signing | None | GHCR image is unsigned | Sign pushed image digest with cosign; publish public key in repo | `.github/workflows/release.yml:96-105` |
| Install verification | None | `install.sh` does not verify checksums or signatures | Add mandatory checksum verification; optional signature verification when cosign is installed | `scripts/install.sh:110-134` |

**Priority:** P0 — security-critical for any project asking users to `curl | bash`.

### G2 — Release automation and version consistency

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Tag/version alignment | Manual | No CI check that `VERSION`, git tag, and `internal/version.Version` match | Add a `release:prepare` job that fails on mismatch | `VERSION:1`, `.github/workflows/release.yml:34-43` |
| Changelog | Manual `CHANGELOG.md` | No automated enforcement of Keep a Changelog format | Add `scripts/check-changelog.sh` and a release-notes generator | `.github/workflows/release.yml:115-116` |
| Prerelease flow | None | Cannot cut `v1.1.0-rc.1` safely | Support `-rc` / `-beta` tags; publish prereleases as GitHub pre-releases and `ghcr.io/...:1.1.0-rc.1` | `.github/workflows/release.yml:4-6,102-104` |
| Post-release validation | None | No job tests the actual released artifact | Add a `release:smoke` job that downloads the release tarball and runs `scripts/smoke-test.sh` against it | `.github/workflows/release.yml:106-116` |

**Priority:** P1 — reduces human error and improves maintainer confidence.

### G3 — Container hardening and embedded UI

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| UI embedding | `Dockerfile` builds Go binary only | Dashboard is built in CI but not copied into image | Add a multi-stage build that copies `internal/ui/dashboard-dist/` or uses `task ui:build` inside Dockerfile | `Dockerfile:14-40`, `.github/workflows/ci.yml:67-72` |
| Healthcheck binary | `docker-compose.yml` uses `wget` | Runtime image installs only `ca-certificates`; `wget` is missing | Switch healthcheck to `muara --config /app/.muara/config.yml health` or install `wget` in runtime stage | `docker-compose.yml:15-19`, `Dockerfile:22` |
| Distroless option | Alpine runtime | Potential CVE surface | Provide a `distroless` or `scratch` variant alongside alpine | `Dockerfile:20` |
| Image labels | None | Missing OCI annotations | Add `org.opencontainers.image.*` labels in build-push-action | `.github/workflows/release.yml:96-105` |
| SBOM for image | None | `syft` only scans source dir | Generate and attach an image SBOM with `anchore/sbom-action` | `.github/workflows/release.yml:53-56` |

**Priority:** P1 — container is a primary distribution mechanism.

### G4 — Install script hardening

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Checksum verification | None | Downloads archive without verification | Verify `checksums.txt` signature and archive hash before extraction | `scripts/install.sh:110-134` |
| Signature verification | None | No cosign/GPG fallback | Add optional `cosign verify-blob` path; document GPG fallback | `scripts/install.sh:1-137` |
| Mirror / fallback | None | Hard dependency on GitHub | Document environment variable for mirror; optionally support GitHub Enterprise | `scripts/install.sh:55-59` |
| Repo case hardening | `REPO="openmuara/openmuara"` | Works, but not future-proof | Lowercase normalization already present in workflow; mirror in script | `scripts/install.sh:10`, `.github/workflows/release.yml:65-66` |
| Windows support | `muara.exe` extracted | Not tested in CI | Add a Windows install smoke test in CI or docs | `scripts/install.sh:82-88` |

**Priority:** P1 — directly impacts first-time user experience.

### G5 — CI completeness and release gates

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Docker build in CI | Only in release workflow | No PR job validates `docker build` or `docker compose up` | Add `docker-build` job to `ci.yml` | `.github/workflows/ci.yml`, `.github/workflows/release.yml:73-75` |
| Full quality matrix in release | Release runs `go test ./...` only | Does not run `task quality` or smoke test | Make release depend on the same `task quality` + smoke used in CI | `.github/workflows/release.yml:29-30`, `Taskfile.yml:135-145` |
| SBOM tool pinning | `go install syft@v1.46.0` | Tool installed at release time without checksum | Pin with `setup-syft` action or verify checksum in workflow | `.github/workflows/release.yml:53-55` |
| Trivy severity gate | Uploads SARIF, no fail | Critical vulnerabilities may be silently merged | Add `exit-code: '1'` and `severity: 'CRITICAL,HIGH'` for release scans; keep `if: always()` upload | `.github/workflows/release.yml:76-87` |
| Release permissions | `contents: write`, `packages: write` | Broad but necessary; no OIDC | Use OIDC-based Sigstore signing to avoid long-lived secrets | `.github/workflows/release.yml:8-11` |

**Priority:** P1 — aligns release pipeline with CI quality bar.

### G6 — Documentation and contributor experience

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Release runbook | None | Contributors do not know how to cut a release | Add `runbooks/release.md` with step-by-step tag/changelog/rollback instructions | `runbooks/` (no release runbook exists) |
| Install verification docs | Brief script header | No user-facing verification guide | Add `docs/install.md` with checksum/signature verification examples | `scripts/install.sh:1-7` |
| Local CI validation | None | Contributors cannot test workflows locally | Document `act` usage and provide `.actrc` / sample secrets | `docs/contributing.md` (to be updated) |
| Badges | Missing | README does not display CI/release health | Add CI, release, container, license, and OpenSSF Scorecard badges | `README.md` (to be updated) |

**Priority:** P2 — important for OSS credibility.

---

### G7 — OpenSSF Scorecard automation and branch protection

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| OpenSSF Scorecard action | None | Score is not tracked over time | Add `.github/workflows/scorecard.yml` and a README badge | `.github/workflows/` (no scorecard.yml) |
| Branch protection rules | Not documented | Required checks may not be enforced | Document required status checks for `main`/`dev`; enable rulesets on GitHub | `AGENTS.md` branch rules |
| Signed commits | Not enforced | Maintainers can push unsigned commits | Require signed commits on `main` and `dev` | GitHub repository settings |
| Private vulnerability reporting | Enabled | — | Keep enabled; verify Security Advisory form is reachable | `.github/SECURITY.md` |

**Priority:** P1 — directly impacts OSS trust and Scorecard score.

### G8 — Workflow hardening and reproducible builds

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Token permissions | `contents: read`/`write` at workflow top level | Overly broad default permissions | Set top-level `permissions: {}` and grant minimal job-level permissions | `.github/workflows/ci.yml:9-11`, `.github/workflows/release.yml:8-11` |
| Third-party action pinning | Some actions pinned to major versions | Supply-chain risk from mutable tags | Pin every third-party action to a full SHA; document update cadence | `.github/workflows/*.yml` |
| Reproducible builds | `-trimpath` used in release | `-buildvcs=false` and consistent flags missing | Add `-trimpath -buildvcs=false` to all build paths; record build metadata in SBOM | `Taskfile.yml:157`, `.github/workflows/release.yml:43` |
| Cache poisoning | `actions/setup-go`/`setup-node` cache enabled | Caches are scoped but not auditable | Document cache key strategy; avoid caching `dist/` | `.github/workflows/*.yml` |

**Priority:** P1 — hardens the supply chain without user-visible friction.

### G9 — Native GitHub artifact attestations

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Artifact attestations | None | GitHub's native build provenance is unused | Add `actions/attest-build-provenance` for release tarballs and images | `.github/workflows/release.yml:106-116` |

**Priority:** P1 — complements SLSA with GitHub-native, user-friendly verification via `gh attestation verify`.

### G10 — Container runtime hardening and vulnerability management

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Read-only root filesystem | Not set | Container can mutate its own filesystem | Add `read_only: true` in Compose; document writable volumes | `docker-compose.yml:11` |
| Capability dropping | Not set | Container retains unnecessary Linux capabilities | Drop all capabilities in Compose and Kubernetes examples | `docker-compose.yml` |
| Distroless variant | Not produced | Single alpine image is the only option | Provide `-distroless` tag as decided in D9 | `.github/workflows/release.yml:96-105` |
| VEX / accepted CVEs | No process | Trivy gate may block on unfixable upstream CVEs | Create `docs/security/cve-exceptions.md` and a VEX document | `docs/security.md` |

**Priority:** P2 — defense in depth for container deployments.

### G11 — Release automation, notifications, and retention

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Manual tag-only trigger | `push: tags: ['v*']` | Accidental tag push can release | Add `workflow_dispatch` for controlled releases; keep tag push as fallback | `.github/workflows/release.yml:3-6` |
| Release failure alerts | None | Failed releases may go unnoticed | Add a failure notification step (GitHub Issues, Slack webhook, or email) | `.github/workflows/release.yml` |
| Release retention | No policy | Old prereleases accumulate | Document retention: keep last 10 stable releases and last 5 prereleases | `runbooks/release.md` |

**Priority:** P2 — improves maintainer experience and reduces accidents.

### G12 — Continuous validation and dependency freshness

| Item | Current | Gap | Recommendation | Source |
|------|---------|-----|----------------|--------|
| Scheduled builds | None | Nightlies catch upstream breakage | Add weekly `schedule` trigger to `ci.yml` | `.github/workflows/ci.yml:3-7` |
| Dependency update tool | Dependabot configured | Go module grouping not configured | Group minor/patch Go and npm updates; enable auto-merge for patch | `.github/dependabot.yml` |
| Stale release artifacts | Not checked | Old SBOMs/checksums left in `dist/` | Clean `dist/` at start of release job; use `actions/upload-artifact` only for debugging | `.github/workflows/release.yml` |

**Priority:** P2 — keeps the pipeline healthy over time.

---

## Priority matrix

| ID | Recommendation | Security | Trust | Maintenance | Priority |
|----|----------------|----------|-------|-------------|----------|
| G1.1 | SLSA provenance generation | High | High | Low | P0 |
| G1.2 | cosign-signed checksums | High | High | Low | P0 |
| G1.3 | cosign-signed container images | High | High | Low | P0 |
| G4.1 | Install-script checksum verification | High | High | Low | P0 |
| G2.1 | Version/tag alignment gate | Medium | High | Low | P1 |
| G2.3 | Prerelease flow | Medium | Medium | Low | P1 |
| G3.1 | Embed dashboard in Docker image | Medium | High | Medium | P1 |
| G3.2 | Fix Docker healthcheck | Medium | High | Low | P1 |
| G5.1 | Docker build CI job | Medium | High | Low | P1 |
| G5.4 | Trivy severity gate | Medium | High | Low | P1 |
| G2.4 | Post-release artifact smoke test | Medium | High | Medium | P1 |
| G7.1 | OpenSSF Scorecard action | Medium | High | Low | P1 |
| G7.2 | Branch protection rules documented | Medium | High | Low | P1 |
| G8.1 | Minimal workflow token permissions | High | High | Low | P1 |
| G8.2 | Full SHA action pinning | High | High | Low | P1 |
| G8.3 | Reproducible build flags | Medium | Medium | Low | P1 |
| G9.1 | GitHub artifact attestations | High | High | Low | P1 |
| G10.1 | Distroless container variant | Medium | Medium | Medium | P2 |
| G10.2 | Read-only root filesystem and dropped capabilities | Medium | Medium | Low | P2 |
| G10.3 | VEX / CVE exception process | Low | High | Low | P2 |
| G11.1 | workflow_dispatch release trigger | Low | High | Low | P2 |
| G11.2 | Release failure notification | Low | Medium | Low | P2 |
| G12.1 | Scheduled CI builds | Low | Medium | Low | P2 |
| G4.3 | Mirror/fallback support | Low | Medium | Low | P2 |
| G6.1 | Release runbook | Low | Medium | Low | P2 |
| G3.4 | OCI image labels | Low | Medium | Low | P2 |

---

## Recommended tool stack

| Concern | Primary tool | Alternative | Notes |
|---------|--------------|-------------|-------|
| SLSA provenance | `slsa-framework/slsa-github-generator` | Hand-rolled in-toto | Use official generator for Level 3 |
| Artifact signing | Sigstore cosign | GPG + keyserver | cosign is keyless via OIDC |
| Container signing | Sigstore cosign | Notary v2 | Attach signature to GHCR digest |
| SBOM generation | `anchore/sbom-action` (Syft) | `ko` + `spdx-sbom-generator` | Already using Syft |
| Vulnerability scan | Trivy + SARIF | Grype | Already using Trivy |
| Secret scanning | gitleaks-action | truffleHog | Already in CI |
| Local CI | `nektos/act` | Fork-based testing | Document both |

---

## Standards mapping

| Standard / Framework | How we satisfy it after implementation |
|----------------------|----------------------------------------|
| **OpenSSF Scorecard** | Signed releases, SLSA provenance, dependency update tool, token permissions, branch protection, security policy, Scorecard action, pinned dependencies |
| **SLSA v1.0 Build L3** | Provenance generator, isolated build runners, hermetic-ish builds with `-trimpath -buildvcs=false`, artifact attestations |
| **SemVer 2.0.0** | `VERSION` file, tag validation, prerelease tag support |
| **Keep a Changelog 1.1** | Enforced format check, release-notes generator |
| **OCI Image Spec** | Labels, multi-arch index, signed manifests, distroless variant |
| **GitHub Artifact Attestations** | `actions/attest-build-provenance` for tarballs and container images |
| **VEX** | Documented CVE exception process and machine-readable VEX file |

---

## Before / after summary

| Capability | Before | After |
|------------|--------|-------|
| Release artifact signing | Unsigned `checksums.txt` | cosign-signed `checksums.txt.sig` + GitHub attestation |
| Supply-chain provenance | None | SLSA Level 3 `.intoto.jsonl` per release |
| Container signing | Unsigned GHCR image | cosign-signed image digest + GitHub attestation |
| Install verification | None | SHA256 + optional cosign verification |
| Dashboard in container | Not included | Embedded from `internal/ui/dashboard-dist/` |
| Container healthcheck | Broken (`wget` missing) | `muara health` |
| PR Docker validation | None | `docker-build` job in CI |
| Version alignment | Manual | Automated CI gate |
| Prerelease flow | None | Safe RC flow, no `latest` promotion |
| Post-release validation | None | Smoke tests on published binary + image |
| Release documentation | Minimal | `runbooks/release.md`, `docs/install.md`, badges |
| OpenSSF Scorecard | Not tracked | Action + badge; target ≥ 8.5 |
| Workflow token hardening | Broad top-level permissions | Minimal job-level permissions |
| Reproducible builds | Partial | `-trimpath -buildvcs=false` on all paths |
| Container hardening | Basic non-root user | Read-only rootfs, dropped caps, distroless option |
| Vulnerability exceptions | No process | Documented VEX + CVE exception file |
| Release trigger | Tag push only | `workflow_dispatch` + tag push |
| Release failure alerting | Silent | Notification step on failure |

## Non-goals

The following are intentionally out of scope for this initiative to keep the plan focused:

- Replacing the custom release workflow with GoReleaser (can be a follow-up).
- Automated package-manager distribution (Homebrew, apt, Chocolatey).
- Multi-region container registry mirrors.
- Automated rollback orchestration beyond documented runbooks.
