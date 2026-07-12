> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Appendix

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09

---

## A. Sample GitHub Release asset layout

After the hardened release workflow runs for `v1.1.0`, the release assets should look like:

```text
muara-linux-amd64.tar.gz
muara-linux-arm64.tar.gz
muara-darwin-amd64.tar.gz
muara-darwin-arm64.tar.gz
muara-windows-amd64.tar.gz
muara-windows-arm64.tar.gz
checksums.txt
checksums.txt.sig
sbom.spdx.json
sbom-dashboard.spdx.json
sbom-website.spdx.json
openmuara-1.1.0.intoto.jsonl
```

---

## B. Test matrix

### Build targets

| GOOS | GOARCH | Archive | Smoke tested |
|------|--------|---------|--------------|
| linux | amd64 | ✅ | ✅ post-release |
| linux | arm64 | ✅ | ❌ (cross-compiled, tested via CI build) |
| darwin | amd64 | ✅ | ❌ |
| darwin | arm64 | ✅ | ❌ |
| windows | amd64 | ✅ | ❌ |
| windows | arm64 | ✅ | ❌ |

### Container variants

| Variant | Base | When to use | Smoke tested |
|---------|------|-------------|--------------|
| Default | `alpine:3.21` | General use, debugging | ✅ |
| Distroless | `gcr.io/distroless/static-debian12` | Security-sensitive deployments | ✅ if implemented |

### Install script environments

| Environment | Method | Verified |
|-------------|--------|----------|
| Ubuntu 24.04 | `curl ... \| bash` | ✅ |
| macOS 14 (arm64) | `curl ... \| bash` | ✅ |
| Windows 11 (WSL2) | `curl ... \| bash` | ✅ |
| Clean container (no Go/npm) | `install.sh --dry-run` | ✅ in CI |

---

## C. Verification commands

### Verify release checksums

```bash
curl -LO https://github.com/openmuara/openmuara/releases/download/v1.1.0/checksums.txt
curl -LO https://github.com/openmuara/openmuara/releases/download/v1.1.0/muara-linux-amd64.tar.gz
sha256sum -c checksums.txt --strict --ignore-missing
```

### Verify checksum signature with cosign

```bash
curl -LO https://github.com/openmuara/openmuara/releases/download/v1.1.0/checksums.txt.sig
cosign verify-blob \
  --signature checksums.txt.sig \
  --certificate-identity-regexp 'https://github.com/openmuara/openmuara/.github/workflows/release.yml@refs/tags/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt
```

### Verify container image signature

```bash
cosign verify \
  --certificate-identity-regexp 'https://github.com/openmuara/openmuara/.github/workflows/release.yml@refs/tags/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  ghcr.io/openmuara/openmuara:1.1.0
```

### Verify SLSA provenance

```bash
slsa-verifier verify-artifact \
  --provenance-path openmuara-1.1.0.intoto.jsonl \
  --source-uri github.com/openmuara/openmuara \
  --source-tag v1.1.0 \
  muara-linux-amd64.tar.gz
```

### Verify GitHub artifact attestation

```bash
gh attestation verify muara-linux-amd64.tar.gz \
  --repo openmuara/openmuara \
  --predicate-type https://in-toto.io/attestation/release/v0.1

gh attestation verify oci://ghcr.io/openmuara/openmuara:1.1.0 \
  --repo openmuara/openmuara
```

---

## D. Sample `.actrc`

```bash
# .actrc — local GitHub Actions validation with nektos/act
-P ubuntu-latest=node:16-bullseye
--secret GITHUB_TOKEN=<your-token>
--artifact-server-path /tmp/act-artifacts
```

Run a specific job locally:

```bash
act -j docker-build
act -j install-dry-run
act -j release --eventpath .github/test-events/release.json
```

---

## E. Sample release event payload

Save as `.github/test-events/release.json` for local `act` testing:

```json
{
  "ref": "refs/tags/v0.0.0-test.1",
  "ref_name": "v0.0.0-test.1"
}
```

---

## F. `runbooks/release.md` outline

1. Pre-release checks
   - Ensure `VERSION` matches intended tag.
   - Ensure `CHANGELOG.md` has a section for the version.
   - Ensure CI is green on `dev`.
2. Create and push the tag
   - `git tag -a v1.1.0 -m "Release v1.1.0"`
   - `git push origin v1.1.0`
3. Monitor the release workflow
   - Wait for all jobs to complete.
   - Verify provenance, signatures, and SBOMs are attached.
4. Post-release validation
   - Download the linux/amd64 tarball and verify checksum/signature.
   - Pull the container image and run `docker compose up`.
   - Run `scripts/smoke-test.sh` against the release binary.
5. If something goes wrong
   - Follow the rollback plan in `EXECUTION_PLAN.md`.
   - Document the issue in `CHANGELOG.md` under `[Unreleased]`.

---

## G. Useful references

- [SLSA GitHub Generator](https://github.com/slsa-framework/slsa-github-generator)
- [Sigstore cosign](https://docs.sigstore.dev/cosign/overview/)
- [OpenSSF Scorecard](https://github.com/ossf/scorecard)
- [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [OCI Image Spec Annotations](https://specs.opencontainers.org/image-spec/annotations/)
- [GitHub Actions OIDC](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)

---

## H. Common pitfalls and mitigations

| Pitfall | Why it happens | Mitigation |
|---------|----------------|------------|
| cosign signing fails with `token exchange` error | OIDC audience mismatch or incorrect `permissions` | Ensure `id-token: write` on the signing job and the correct `certificate-identity-regexp` |
| SLSA provenance generator cannot find artifacts | Subject hashes not passed correctly | Use base64-encoded concatenated hashes of all release artifacts in a dedicated output |
| `muara health` fails inside Docker because config is not initialized | Entrypoint initializes config only on `start` | Make `health` use default server host/port or read config path from env |
| Dashboard assets missing in container | `internal/ui/dashboard-dist/` not copied from CI artifact | Ensure `ui-build` job artifact is downloaded before `docker build` in release workflow |
| Install script breaks on macOS with `shasum` vs `sha256sum` | macOS uses `shasum -a 256` | Use `sha256sum` if available, fall back to `shasum -a 256`, or verify in Python/Go |
| Prerelease tag still moves `latest` | Semver detection regex is too permissive | Use a strict semver parser and skip `latest` for any prerelease segment |
| Trivy gate blocks release on upstream base-image CVE | New CVE in `alpine:3.21` | Pin base image digest; document exception process; consider distroless variant |
| `checksums.txt.sig` not attached to release | Artifact upload glob is too narrow | Include `dist/checksums.txt.sig` in `softprops/action-gh-release` `files` list |
| Local `act` runs fail because `ui-build` artifact is missing | `act` does not persist artifacts between jobs by default | Use `--artifact-server-path` or run dependent jobs in a single composite job locally |
| Scorecard action fails with insufficient permissions | `security-events: write` missing or `permissions: read-all` not set | Use `permissions: read-all` at workflow top level and `security-events: write` + `id-token: write` on the job |
| GitHub attestation fails for container image | `attestations: write` missing or wrong subject digest | Pass the exact digest from `build-push-action` outputs and grant `attestations: write` |
| Read-only rootfs causes runtime errors | `/tmp`, `/app/.muara`, or dashboard cache need writable paths | Mount empty `tmpfs` for `/tmp` and a volume for `/app/.muara` |
| `workflow_dispatch` input tag differs from `VERSION` | Maintainer provides a tag that does not match the source | Run `verify-version` against both `VERSION` and the dispatch input |
| Dependabot groups create giant PRs | Group includes major updates | Separate major updates from minor/patch groups |

---

## I. Post-implementation review schedule

| Review | When | Owner | Purpose |
|--------|------|-------|---------|
| M1–M3 completion review | After M3 | Human reviewer | Validate security and reliability changes before container/dashboard work is finalized |
| M4–M5 completion review | After M5 | Human reviewer | Validate install script and release validation before declaring release-ready |
| Final initiative review | After M6 | Human reviewer + AI Agent | Complete `REVIEW_CHECKLIST.md` and sign off |
| 30-day follow-up | 30 days after merge | Human reviewer | Check OpenSSF Scorecard, first real release experience, and any CI flakiness |
| Quarterly pipeline review | Every quarter | Maintainer | Re-evaluate tools, base images, and signing mechanisms |

---

## J. Sample scheduled CI trigger

```yaml
on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]
  schedule:
    - cron: '17 4 * * 1'
```

## K. Sample Dependabot grouping

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    groups:
      go-minor-patch:
        patterns: ['*']
        update-types: [minor, patch]
  - package-ecosystem: npm
    directory: /web/dashboard
    schedule:
      interval: weekly
    groups:
      npm-minor-patch:
        patterns: ['*']
        update-types: [minor, patch]
```

## L. README badges

```markdown
![CI](https://github.com/openmuara/openmuara/actions/workflows/ci.yml/badge.svg)
![Release](https://github.com/openmuara/openmuara/actions/workflows/release.yml/badge.svg)
![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/openmuara/openmuara/badge)
![Container](https://ghcr-badge.egpl.dev/openmuara/openmuara/latest_tag?trim=major)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
```

---

## M. Fork test release checklist

Use this checklist when cutting `v0.0.0-test.1` on a personal fork to validate the hardened pipeline before merging to `dev`.

### Before pushing the tag

- [ ] Fork the repository and enable GitHub Actions.
- [ ] Enable GitHub Container Registry on the fork.
- [ ] Ensure `VERSION` contains `0.0.0-test.1`.
- [ ] Ensure `CHANGELOG.md` has a `## [0.0.0-test.1]` section.
- [ ] Push the branch with the hardened workflows to the fork.

### Push and monitor

- [ ] Push tag: `git tag -a v0.0.0-test.1 -m "Test release" && git push origin v0.0.0-test.1`.
- [ ] Confirm `verify-version` and `verify-changelog` jobs pass.
- [ ] Confirm all release artifacts build and upload.

### Artifact verification

- [ ] Download `muara-linux-amd64.tar.gz`, `checksums.txt`, and `checksums.txt.sig`.
- [ ] Run `sha256sum -c checksums.txt --strict --ignore-missing`.
- [ ] Run `cosign verify-blob --signature checksums.txt.sig --certificate-identity-regexp 'https://github.com/<fork>/openmuara/.github/workflows/release.yml@refs/tags/.*' --certificate-oidc-issuer https://token.actions.githubusercontent.com checksums.txt`.
- [ ] Download the `.intoto.jsonl` provenance file.
- [ ] Run `slsa-verifier verify-artifact --provenance-path openmuara-0.0.0-test.1.intoto.jsonl --source-uri github.com/<fork>/openmuara --source-tag v0.0.0-test.1 muara-linux-amd64.tar.gz`.
- [ ] Run `gh attestation verify muara-linux-amd64.tar.gz --repo <fork>/openmuara`.

### Container verification

- [ ] Pull `ghcr.io/<fork>/openmuara:0.0.0-test.1`.
- [ ] Run `cosign verify ghcr.io/<fork>/openmuara:0.0.0-test.1 --certificate-identity-regexp 'https://github.com/<fork>/openmuara/.github/workflows/release.yml@refs/tags/.*' --certificate-oidc-issuer https://token.actions.githubusercontent.com`.
- [ ] Run `gh attestation verify oci://ghcr.io/<fork>/openmuara:0.0.0-test.1 --repo <fork>/openmuara`.
- [ ] Run `docker compose up` and confirm `/_admin` loads.
- [ ] Confirm `docker ps` shows the container as healthy.
- [ ] Pull the distroless tag (`:0.0.0-test.1-distroless`) and confirm it starts.

### Install script verification

- [ ] In a clean Ubuntu container, run `curl -sSL https://raw.githubusercontent.com/<fork>/openmuara/main/scripts/install.sh | bash`.
- [ ] Confirm the installed binary runs `muara version` and reports `0.0.0-test.1`.
- [ ] Confirm a tampered archive is rejected by the install script (CI tamper test).

### Prerelease behavior

- [ ] Push `v0.0.0-test.2-rc.1` and confirm:
  - [ ] GitHub Release is marked as prerelease.
  - [ ] Container tag `:0.0.0-test.2-rc.1` is pushed.
  - [ ] Container tag `:latest` is **not** moved.

---

## N. OpenSSF Scorecard target breakdown

Target score: **≥ 8.5/10**. Below is how each checklist item contributes.

| Check | Current | Target | How this initiative satisfies it |
|-------|---------|--------|----------------------------------|
| Code-Review | Likely partial | Pass | Branch protection docs + required checks |
| Dependency-Update-Tool | Pass | Pass | `dependabot.yml` already exists |
| Maintained | Pass | Pass | Recent commits and releases |
| Security-Policy | Pass | Pass | `.github/SECURITY.md` exists |
| Signed-Releases | Fail | Pass | cosign-signed `checksums.txt.sig` |
| SAST | Pass | Pass | gosec, govulncheck, CodeQL SARIF upload |
| Vulnerabilities | Pass | Pass | Trivy + govulncheck |
| Token-Permissions | Partial | Pass | Minimal top-level permissions |
| Pinned-Dependencies | Partial | Pass | Full SHA pinning |
| CII-Best-Practices | Unknown | Optional | Document Silver/Gold readiness if applicable |
| Fuzzing | Unknown | Optional | Out of scope for this initiative |
| License | Pass | Pass | LICENSE file exists |
| Binary-Artifacts | Pass | Pass | No checked-in binaries |
| Branch-Protection | Partial | Pass | Documented rules + required status checks |

> Scorecard scoring is not strictly additive; the table is illustrative. Run `scorecard-action` after implementation for the authoritative score.

---

## O. Sign-off log

| Review | Date | Outcome | Sign-off |
|--------|------|---------|----------|
| Planning review | | | |
| M1–M3 completion | | | |
| M4–M5 completion | | | |
| Final initiative review | | | |
| 30-day follow-up | | | |
