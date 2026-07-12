> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Review Checklist

> **Status:** ⬜ Draft | **Created:** 2026-07-08 | **Last Updated:** 2026-07-09

Use this checklist to sign off the initiative before merging to `dev` and before cutting the next release.

---

## Pre-implementation review (planning sign-off)

- [ ] `README.md` accurately describes initiative scope and structure.
- [ ] `RECOMMENDATIONS.md` covers all major CI/release gaps with priorities, including SLSA, GitHub attestations, Scorecard, token hardening, reproducible builds, and VEX.
- [ ] `DECISIONS.md` records a decision for every architectural choice (D1–D15).
- [ ] `EXECUTION_PLAN.md` has clear milestones, acceptance criteria, RACI, and updated timeline.
- [ ] `KNOWN_ISSUES.md` is populated with real findings from current files (25 total after expansion).
- [ ] `RISKS.md` includes likelihood, impact, mitigation, and owner for each risk.
- [ ] `CI_INTEGRATION.md` contains concrete workflow snippets for every new capability.
- [ ] All documents use consistent terminology (`OpenMuara`, `muara`, `VERSION`, etc.).
- [ ] No speculative product changes are embedded in planning docs.

---

## Post-implementation review (execution sign-off)

### Security

- [ ] SLSA provenance attestation is attached to every GitHub Release.
- [ ] GitHub artifact attestation is produced for release tarballs and container image.
- [ ] `checksums.txt` is signed with cosign; verification command documented.
- [ ] GHCR image digest is signed with cosign; verification command documented.
- [ ] Trivy scan fails on `CRITICAL` vulnerabilities unless exempted by VEX.
- [ ] `install.sh` verifies SHA256 hash before extraction.
- [ ] `install.sh` verifies cosign signature when cosign is available.
- [ ] `SKIP_VERIFY=1` exists and prints a clear warning.
- [ ] Workflows use minimal top-level permissions and explicit job-level permissions.
- [ ] All third-party actions are pinned to full SHA references.
- [ ] Go builds use `-trimpath -buildvcs=false` for reproducibility.
- [ ] OpenSSF Scorecard action runs and badge is visible.
- [ ] VEX file and `docs/security/cve-exceptions.md` exist.

### Reliability

- [ ] `docker compose up` starts a healthy container.
- [ ] `muara health` returns 0 when healthy and non-zero when unhealthy.
- [ ] Docker image includes the built dashboard at `/_admin`.
- [ ] Container runs with read-only rootfs and dropped capabilities.
- [ ] Distroless image variant builds and passes smoke test.
- [ ] Release workflow fails if `VERSION` does not match the pushed tag.
- [ ] Prerelease tags create GitHub pre-releases and do not move `latest`.
- [ ] `workflow_dispatch` can trigger a release from a chosen tag/ref.
- [ ] Release failures notify maintainers (issue or webhook).

### Validation

- [ ] Post-release smoke test runs against the published linux/amd64 tarball.
- [ ] Post-release container smoke test runs against the published image.
- [ ] CI includes a `docker-build` job on every PR.
- [ ] CI includes an `install-dry-run` job on every PR.
- [ ] CI includes a weekly scheduled build.
- [ ] Full `task quality` passes locally and in CI.
- [ ] Fork-based test release passes every validation step.

### Documentation

- [ ] `runbooks/release.md` exists and is accurate.
- [ ] `docs/install.md` exists and includes verification examples.
- [ ] `.actrc` or equivalent local CI validation guidance exists.
- [ ] Branch protection rules and required status checks are documented.
- [ ] Root `README.md` displays CI, release, container, license, and OpenSSF badges.
- [ ] `CHANGELOG.md` has a CI format check.

---

## Release-readiness sign-off

Before cutting the next release after this initiative:

- [ ] All items above are checked.
- [ ] A test release was performed on a fork using `v0.0.0-test.1`.
- [ ] The test release artifacts were verified with cosign and checksum commands.
- [ ] The test container image was pulled and smoke-tested.
- [ ] `install.sh` was tested in a clean Ubuntu container and a clean macOS environment.
- [ ] Rollback plan was reviewed and understood by the release cutter.

---

## Signatures

| Role | Name | Date | Signature / Approval |
|------|------|------|---------------------|
| Implementer | AI Agent (Kimi Code) | | |
| Reviewer | ___________ | | |
| Release cutter | ___________ | | |
