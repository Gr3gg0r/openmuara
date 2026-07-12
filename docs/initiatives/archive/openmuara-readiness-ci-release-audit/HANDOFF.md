> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Handoff

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09

---

## What this initiative delivers

A hardened, gold-standard CI and release pipeline for OpenMuara that:

- Generates SLSA Level 3 provenance and GitHub artifact attestations for every release.
- Signs release artifacts and container images with Sigstore cosign.
- Verifies downloaded binaries in `scripts/install.sh` before installation.
- Embeds the built dashboard in Docker images, fixes the container healthcheck, and hardens the runtime.
- Supports safe prerelease tags, post-release validation, and controlled `workflow_dispatch` releases.
- Tracks OpenSSF Scorecard, uses minimal workflow permissions, pins actions to SHA, and documents exceptions via VEX.
- Documents the entire process for maintainers and users.

---

## Deliverables

### Planning documents

| File | Purpose |
|------|---------|
| `README.md` | Initiative overview, scope, success criteria |
| `TRACKING.md` | Central execution tracker and findings log |
| `KNOWN_ISSUES.md` | Catalog of current CI/release findings |
| `RISKS.md` | Risk register |
| `RECOMMENDATIONS.md` | Gold-standard gap analysis and priority matrix |
| `DECISIONS.md` | Decision register for every architectural choice |
| `EXECUTION_PLAN.md` | Milestones, acceptance criteria, RACI, timeline |
| `REVIEW_CHECKLIST.md` | Sign-off checklist for planning and execution |
| `CI_INTEGRATION.md` | Concrete workflow changes and YAML snippets |
| `HANDOFF.md` | This file — final state and next steps |
| `APPENDIX.md` | Sample configs, test matrices, commands |

### Product artifacts produced during execution

| Artifact | Location |
|----------|----------|
| `muara health` command | `cmd/muara/health.go` |
| Hardened release workflow | `.github/workflows/release.yml` |
| Docker build CI job | `.github/workflows/ci.yml` |
| Hardened `Dockerfile` | `Dockerfile` |
| Hardened `docker-compose.yml` | `docker-compose.yml` |
| Hardened `scripts/install.sh` | `scripts/install.sh` |
| Release runbook | `runbooks/release.md` |
| Install verification guide | `docs/install.md` |
| Local CI validation guide | `docs/contributing.md`, `.actrc` |
| VEX / CVE exceptions | `docs/security/cve-exceptions.md`, `vex.json` |
| OpenSSF Scorecard workflow | `.github/workflows/scorecard.yml` |
| Branch protection docs | `AGENTS.md` or `docs/contributing.md` |

---

## Final state checklist

When execution is complete, the repository should satisfy:

- [ ] Every GitHub Release includes binaries, checksums, checksum signature, SBOMs, SLSA provenance, GitHub attestation, and release notes.
- [ ] Every GHCR image is signed (cosign + GitHub attestation) and includes OCI labels.
- [ ] `docker compose up` starts a healthy, hardened container with the dashboard.
- [ ] `scripts/install.sh` verifies hashes and signatures by default.
- [ ] Prerelease tags create pre-releases and do not move `latest`.
- [ ] `workflow_dispatch` can trigger controlled releases.
- [ ] CI runs `docker-build` and `install-dry-run` on every PR, plus a weekly scheduled build.
- [ ] Post-release smoke tests run against published artifacts.
- [ ] Workflows use minimal permissions and full SHA-pinned actions.
- [ ] OpenSSF Scorecard action runs and badge is visible.
- [ ] VEX file and CVE exception process exist.
- [ ] Documentation is complete and badges are visible.

---

## Next steps after initiative completion

1. **Fork-based test release:** Cut `v0.0.0-test.1` on a personal fork and verify every artifact, signature, attestation, and image variant.
2. **OpenSSF Scorecard:** Confirm the score is ≥ 8.5 and address any remaining gaps.
3. **Package managers:** Evaluate Homebrew formula, apt repository, or Chocolatey package based on user demand.
4. **Observability:** Add release metrics (download counts, image pulls) to the dashboard or runbooks.
5. **Automated dependency updates:** Confirm Dependabot grouping is working; consider Renovate for advanced grouping rules.
6. **Quarterly pipeline review:** Re-evaluate base images, signing tools, and Scorecard recommendations.

---

## Open questions to resolve during execution

- [ ] Do we want a project GPG key as a fallback for cosign? (Decision D2)
- [ ] Do we keep the `-distroless` variant in scope for M3 or defer to a follow-up? (Decision D9)
- [ ] Should the release workflow also build and attach a Windows installer (MSI/zip) or remain tarball-only?
- [ ] What is the acceptable grace period for `HIGH` Trivy findings before they fail the build?
- [ ] Which release failure notification channel should be used? (GitHub issue, Slack webhook, email)
- [ ] Should `HIGH` Trivy findings block the release or only `CRITICAL`? (Decision D14)
- [ ] Should we require signed commits on `main`/`dev`? (Decision D12, branch protection)

---

## Contact / ownership

| Role | Owner |
|------|-------|
| Initiative lead | AI Agent (Kimi Code) |
| Human reviewer | ___________ |
| Release engineering | ___________ |
| Security review | ___________ |
