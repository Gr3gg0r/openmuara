> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — CI & Release Audit Risk Register

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-09
> **Status:** ⬜ Draft

---

## Active risks

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|----|------|------------|--------|------------|-------|
| R01 | CI passes locally but fails on GitHub due to workflow syntax or environment differences | Medium | High | Lint workflows with `actionlint`; test with `act` and on a fork | AI Agent |
| R02 | Docker image fails without prebuilt UI or reports unhealthy | Medium | High | Add `docker-build` CI job; implement `muara health`; validate `docker compose up` | AI Agent |
| R03 | Install script breaks on new OS version or architecture | Medium | Medium | Test in clean containers for latest Ubuntu/macOS; add dry-run CI matrix | AI Agent |
| R04 | cosign/Sigstore signing fails in release workflow due to OIDC or network issues | Low | High | Test on a fork first; document GPG fallback path; keep workflow idempotent | AI Agent |
| R05 | SLSA provenance generator produces incompatible attestation format | Low | Medium | Use official generator; validate with `slsa-verifier` during fork test | AI Agent |
| R06 | Hardened install script breaks existing `curl \| bash` users | Low | High | Maintain backward-compatible invocation; `SKIP_VERIFY=1` for edge cases; document changes | AI Agent |
| R07 | Trivy severity gate blocks releases on upstream CVEs in base image | Medium | Medium | Pin base image digest; maintain a documented exception process; offer distroless variant | AI Agent |
| R08 | Version mismatch gate blocks legitimate releases due to human error | Medium | Medium | Add clear error messages; create `scripts/release-prep.sh` to bump `VERSION` and `CHANGELOG.md` | AI Agent |
| R09 | Release artifact signing increases build time beyond acceptable threshold | Low | Low | Benchmark on fork; cosign keyless signing is typically fast | AI Agent |
| R10 | Contributor confusion over new local validation workflow | Low | Low | Document `act` setup and provide copy-paste commands in `docs/contributing.md` | AI Agent |
| R11 | GitHub artifact attestations fail due to permissions or digest mismatch | Low | High | Grant `attestations: write`; use exact digest from `build-push-action` outputs; test on fork | AI Agent |
| R12 | Scorecard score stays below 8.5 due to inherited repo settings | Medium | Medium | Review Scorecard results; document required branch protection and signed-commit settings | AI Agent |
| R13 | Read-only rootfs or dropped capabilities break runtime behavior | Medium | High | Validate `docker compose up` in CI with full smoke test; document writable volumes | AI Agent |
| R14 | Distroless image lacks shell for debugging and entrypoint scripts | Low | Medium | Use a static binary entrypoint; test distroless smoke path separately | AI Agent |
| R15 | `workflow_dispatch` release uses wrong ref and bypasses version gate | Low | High | Validate input ref against `VERSION`; require environment protection rules for release job | AI Agent |

---

## Risk treatment plan

| ID | Treatment | Target state |
|----|-----------|--------------|
| R01 | Mitigate | All workflows linted and tested on a fork before merge |
| R02 | Mitigate | `docker-build` job green on every PR; container health passes locally and in CI |
| R03 | Mitigate | Install script tested across OS/arch matrix in CI |
| R04 | Mitigate + Contingency | Fork-tested; GPG fallback documented if cosign becomes unavailable |
| R05 | Mitigate | Attestation validated with `slsa-verifier` in fork test |
| R06 | Mitigate | Default invocation remains `curl \| bash`; verification is transparent |
| R07 | Mitigate + Accept | Pin base image; accept rare upstream exceptions via documented process |
| R08 | Mitigate | Add release-prep helper script and clear CI error output |
| R09 | Accept | Accept small increase in release time for security gain |
| R10 | Mitigate | Clear documentation and onboarding checklist |
| R11 | Mitigate | Attestation step passes on fork with correct permissions |
| R12 | Mitigate | Scorecard ≥ 8.5 or documented gap accepted by reviewer |
| R13 | Mitigate | Smoke test passes on hardened container before merge |
| R14 | Mitigate | Distroless image has its own smoke test in CI |
| R15 | Mitigate + Contingency | Input ref verified; environment protection rules enabled |

---

## Resolved risks

None yet. This section will be updated as milestones close.

---

## Risk monitoring

- Review this register at the end of every milestone.
- Escalate any new risk with Likelihood × Impact ≥ Medium-High to the human reviewer immediately.
