> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Decision Register

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09

This file records the decisions required to move from the current CI/release state to a gold-standard pipeline. Each decision includes context, options, a recommended option, and the acceptance criteria that prove it is implemented correctly.

---

## D1 — Artifact provenance level

| | |
|---|---|
| **Question** | What supply-chain provenance standard should release artifacts meet? |
| **Context** | Users and contributors need confidence that release binaries and images were built from the tagged source and not tampered with. |
| **Options** | 1. None (current). 2. Basic in-toto statement generated manually. 3. SLSA Level 3 via `slsa-framework/slsa-github-generator`. |
| **Recommended** | **Option 3 — SLSA Level 3** using the official generator. It is the de-facto OSS standard and integrates cleanly with GitHub Actions OIDC. |
| **Acceptance criteria** | `.github/workflows/release.yml` produces a `.intoto.jsonl` provenance attestation attached to the GitHub Release; README links to SLSA badge/verification instructions. |

---

## D2 — Signing mechanism

| | |
|---|---|
| **Question** | How should release artifacts and container images be signed? |
| **Context** | Long-lived GPG keys add operational burden and leak risk. Sigstore cosign with OIDC is keyless and transparent. |
| **Options** | 1. GPG signing only. 2. cosign keyless signing only. 3. Both cosign (primary) and a project GPG fallback. |
| **Recommended** | **Option 2 — cosign keyless signing** as the primary mechanism. Document verification with `cosign verify-blob` and `cosign verify`. Add GPG fallback only if users explicitly request it later. |
| **Acceptance criteria** | Release workflow signs `checksums.txt` and the container image digest; verification commands are documented in `docs/install.md` and `runbooks/release.md`. |

---

## D3 — Container image UI embedding strategy

| | |
|---|---|
| **Question** | Should the Docker image build the dashboard inside the Dockerfile or copy a prebuilt artifact? |
| **Context** | The current `Dockerfile` builds only the Go binary. CI builds the dashboard into `internal/ui/dashboard-dist/`, but the image does not include it. Embedding Node in the builder stage increases image build time and supply-chain surface. |
| **Options** | 1. Build dashboard inside Dockerfile (self-contained). 2. Copy prebuilt `internal/ui/dashboard-dist/` from the build context (faster, requires CI artifact). 3. Support both: Dockerfile accepts a prebuilt dist and falls back to building if absent. |
| **Recommended** | **Option 3 — copy prebuilt dist with optional fallback build.** Default CI/release path copies `internal/ui/dashboard-dist/`; local `docker build` falls back to an embedded placeholder or a minimal build if Node is available. This keeps release builds fast and reproducible while preserving local usability. |
| **Acceptance criteria** | `docker build` in release workflow copies the artifact from the `ui-build` job; `docker run` serves the dashboard at `/_admin`; fallback behavior is documented. |

---

## D4 — Docker healthcheck implementation

| | |
|---|---|
| **Question** | How should the container healthcheck work without relying on `wget`? |
| **Context** | `docker-compose.yml` calls `wget`, but the runtime image only installs `ca-certificates`. |
| **Options** | 1. Install `wget` in the runtime stage. 2. Add a `muara health` CLI subcommand. 3. Use `HEALTHCHECK CMD-SHELL` with a Go binary subcommand. |
| **Recommended** | **Option 2 — add `muara health` subcommand** that exits 0 when `/healthz` returns `{"status":"ok"}`. This removes external binary dependencies and works in both Docker and orchestrators. |
| **Acceptance criteria** | `cmd/muara` implements `health`; `docker-compose.yml` uses `["CMD", "muara", "health"]`; Dockerfile includes `HEALTHCHECK` instruction using the same command. |

---

## D5 — Release-candidate / prerelease flow

| | |
|---|---|
| **Question** | Should prerelease tags (e.g., `v1.1.0-rc.1`) be fully automated? |
| **Context** | Currently only `v*` tags trigger releases. Prereleases need separate handling to avoid updating `latest` container tag. |
| **Options** | 1. Reject prerelease tags. 2. Allow prerelease tags but push them to GHCR only. 3. Allow prerelease tags, create GitHub pre-releases, and push versioned container tags without `latest`. |
| **Recommended** | **Option 3 — full prerelease support.** The release workflow detects semver prerelease segments and skips the `latest` tag for container pushes and creates a GitHub pre-release. |
| **Acceptance criteria** | Pushing `v1.1.0-rc.1` creates a GitHub pre-release, pushes `ghcr.io/openmuara/openmuara:1.1.0-rc.1`, and does not move `latest`. |

---

## D6 — Install-script verification behavior

| | |
|---|---|
| **Question** | Should `install.sh` fail closed if checksum/signature verification cannot be performed? |
| **Context** | Requiring verification by default improves security but may break air-gapped or mirror deployments. |
| **Options** | 1. Verify if files exist, warn otherwise. 2. Verify if files exist, fail otherwise. 3. Always verify; provide `SKIP_VERIFY=1` escape hatch. |
| **Recommended** | **Option 3 — verify by default with an explicit escape hatch.** This matches security best practice while preserving flexibility for advanced users. |
| **Acceptance criteria** | `install.sh` downloads `checksums.txt` and a detached signature, verifies the archive hash, and verifies the signature when cosign is available; `SKIP_VERIFY=1` bypasses verification and prints a warning. |

---

## D7 — Version source of truth

| | |
|---|---|
| **Question** | What is the single source of truth for the project version? |
| **Context** | Today `VERSION` file holds `1.0.0`; git tags trigger releases; `internal/version.Version` is injected at build time. |
| **Options** | 1. Keep `VERSION` file as source of truth. 2. Derive version from git tags at build time. 3. Keep `VERSION` as source of truth but validate it against the git tag in CI. |
| **Recommended** | **Option 3 — `VERSION` file as source of truth with CI alignment gate.** It is simple, deterministic, and easy to bump; the release workflow fails if the tag does not match `VERSION`. |
| **Acceptance criteria** | Release workflow has a `verify-version` job that fails when `refs/tags/v${VERSION}` does not match the pushed tag; dev builds use `dev-<short-sha>` via `Taskfile.yml`. |

---

## D8 — Release notes generation

| | |
|---|---|
| **Question** | Should release notes be generated from `CHANGELOG.md` or from commit history? |
| **Context** | Current release uses `body_path: CHANGELOG.md`. This is simple but may include unreleased changes if not carefully edited. |
| **Options** | 1. Continue manual `CHANGELOG.md`. 2. Generate from conventional commits. 3. Extract the relevant version section from `CHANGELOG.md` automatically. |
| **Recommended** | **Option 3 — extract version section from `CHANGELOG.md`.** It preserves human-curated release notes while removing the risk of publishing the whole file. |
| **Acceptance criteria** | Release workflow extracts the `## [X.Y.Z]` section matching the tag and uses it as the release body; a CI check ensures the section exists before a release can be cut. |

---

## D9 — Container registry variants

| | |
|---|---|
| **Question** | Should the release produce only the alpine-based image or additional variants? |
| **Context** | Alpine is small but still has a package manager and shell. Distroless/scratch reduces CVE surface but complicates debugging. |
| **Options** | 1. Alpine only. 2. Alpine + distroless tags. 3. Alpine + scratch (static binary only). |
| **Recommended** | **Option 2 — alpine default with a `-distroless` variant.** Keep the default image familiar and debuggable; offer distroless for security-sensitive deployments. |
| **Acceptance criteria** | Release workflow pushes `:<version>`, `:latest`, `:<version>-distroless`, and `:distroless` tags; README documents when to use each. |

---

## D10 — Local CI validation strategy

| | |
|---|---|
| **Question** | How should contributors validate workflow changes locally? |
| **Context** | Pushing to GitHub to test workflows is slow and pollutes the git history. |
| **Options** | 1. Document `act` usage. 2. Provide a fork-based validation guide. 3. Both. |
| **Recommended** | **Option 3 — both `act` and fork-based validation.** `act` covers fast feedback; fork-based testing covers secrets and release-path behavior. |
| **Acceptance criteria** | `docs/contributing.md` includes an `act` section and a `.actrc` file; `runbooks/release.md` includes a fork-based dry-run checklist. |

---

## D11 — Provenance mechanism mix

| | |
|---|---|
| **Question** | Should the project use SLSA, GitHub artifact attestations, or both? |
| **Context** | SLSA Level 3 is the industry standard, but GitHub artifact attestations are native and easy for users to verify with `gh attestation verify`. |
| **Options** | 1. SLSA only. 2. GitHub attestations only. 3. Both. |
| **Recommended** | **Option 3 — both.** SLSA covers the formal supply-chain standard; GitHub attestations provide low-friction verification for GitHub-centric users. |
| **Acceptance criteria** | Every release includes both `.intoto.jsonl` and GitHub attestations; both verification methods are documented. |

---

## D12 — Workflow token permission model

| | |
|---|---|
| **Question** | How should GitHub Actions token permissions be scoped? |
| **Context** | Broad top-level permissions violate least-privilege and reduce OpenSSF Scorecard score. |
| **Options** | 1. Keep current top-level permissions. 2. Set `permissions: {}` at workflow level and grant per-job permissions. 3. Use a mix. |
| **Recommended** | **Option 2 — minimal top-level permissions with explicit job-level grants.** This is the OpenSSF recommendation and prevents permission inheritance. |
| **Acceptance criteria** | Both `ci.yml` and `release.yml` set top-level `permissions: {}`; every job declares only the permissions it needs. |

---

## D13 — Action pinning policy

| | |
|---|---|
| **Question** | How should third-party GitHub Actions be pinned? |
| **Context** | Pinning to mutable tags exposes the project to supply-chain attacks; full SHA pinning is safest. |
| **Options** | 1. Pin to major version tags. 2. Pin to full SHAs with a version comment. 3. Use GitHub's immutable actions where available. |
| **Recommended** | **Option 2 — full SHA pinning with version comments.** Add a note in `docs/contributing.md` and rely on Dependabot for SHA update PRs. |
| **Acceptance criteria** | All third-party actions in `.github/workflows/*.yml` use full SHA references with `# vX.Y.Z` comments. |

---

## D14 — Vulnerability exception process

| | |
|---|---|
| **Question** | How should the project handle unfixable upstream CVEs that would otherwise block release? |
| **Context** | Trivy may flag CVEs in base images or dependencies with no available fix; a documented exception process prevents silent overrides. |
| **Options** | 1. Ignore findings. 2. Document exceptions in a markdown file. 3. Maintain a machine-readable VEX file and feed it to Trivy. |
| **Recommended** | **Option 3 — VEX file plus markdown rationale.** Use Trivy's `--vex` flag and keep a human-readable `docs/security/cve-exceptions.md`. |
| **Acceptance criteria** | A VEX file exists; Trivy in `release.yml` consumes it; exceptions are reviewed quarterly. |

---

## D15 — Release trigger control

| | |
|---|---|
| **Question** | How should releases be triggered to balance automation and safety? |
| **Context** | Tag-push-only releases are convenient but allow accidental or compromised releases. |
| **Options** | 1. Tag push only. 2. `workflow_dispatch` only. 3. `workflow_dispatch` primary with tag push as fallback. |
| **Recommended** | **Option 3 — `workflow_dispatch` primary with tag-push fallback.** Maintainers can run a controlled release from the UI; tag push still works for automation and CI integrations. |
| **Acceptance criteria** | `release.yml` supports `workflow_dispatch` with a tag/ref input; tag push still triggers the workflow; both paths pass fork tests. |

---

## Decision summary table

| ID | Decision | Recommended option |
|----|----------|-------------------|
| D1 | Provenance level | SLSA Level 3 |
| D2 | Signing mechanism | cosign keyless (primary) |
| D3 | UI embedding | Copy prebuilt dist with optional fallback |
| D4 | Healthcheck | New `muara health` subcommand |
| D5 | Prerelease flow | Full automated prerelease support |
| D6 | Install verification | Verify by default, `SKIP_VERIFY=1` escape hatch |
| D7 | Version source | `VERSION` file + CI alignment gate |
| D8 | Release notes | Extract section from `CHANGELOG.md` |
| D9 | Image variants | Alpine default + distroless variant |
| D10 | Local validation | `act` + fork-based testing |
| D11 | Provenance mix | SLSA + GitHub artifact attestations |
| D12 | Token permissions | Minimal top-level, explicit job-level |
| D13 | Action pinning | Full SHA with version comments |
| D14 | Vulnerability exceptions | VEX file + markdown rationale |
| D15 | Release trigger | `workflow_dispatch` primary, tag push fallback |
