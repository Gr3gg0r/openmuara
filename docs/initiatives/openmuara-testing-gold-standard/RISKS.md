> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Testing Gold Standard — Risk Register

> **Purpose:** Track what could go wrong during the testing initiative and how to respond.

---

## Risk Matrix

| ID | Risk | Likelihood | Impact | Status | Owner |
|----|------|------------|--------|--------|-------|
| R01 | Testability refactor breaks existing tests | Medium | Major | 🟡 Mitigated | AI Agent |
| R02 | 80% coverage target creates brittle tests | Low | Minor | 🟡 Mitigated | AI Agent |
| R03 | Provider contract tests drift from real providers | Medium | Major | 🟡 Mitigated | AI Agent |
| R04 | CI time increases with split jobs | Medium | Minor | 🟡 Mitigated | AI Agent |
| R05 | Contributors resist new test conventions | Low | Minor | 🟡 Mitigated | AI Agent |

---

## Detailed Risk Entries

### R01 — Testability refactor breaks existing tests

- **Description:** Refactoring provider/plugin registries and CLI start could change behavior or break existing tests.
- **Trigger:** Global state is removed incorrectly or dependency wiring changes.
- **Impact:** `dev` becomes red; other work blocked.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Keep commits small and focused; run `task check` after every prompt; preserve existing public APIs where possible.
- **Rollback Plan:** Revert the refactor commit or wrap new behavior behind a feature flag.
- **Monitoring:** `go test ./...`, `./scripts/smoke-test.sh`.
- **Status:** 🟡 Mitigated

### R02 — 80% coverage target creates brittle tests

- **Description:** Pressure to hit 80% may lead to tests that assert implementation details.
- **Trigger:** Writing tests just to cover lines rather than behavior.
- **Impact:** Tests break on every refactor; maintainability drops.
- **Likelihood:** Low
- **Impact Level:** Minor
- **Mitigation:** Use fakes and golden files; avoid testing trivial getters; review tests for behavioral intent.
- **Rollback Plan:** Lower threshold or remove low-value tests.
- **Monitoring:** Test review checklist in `runbooks/testing.md`.
- **Status:** 🟡 Mitigated

### R03 — Provider contract tests drift from real providers

- **Description:** Golden files may not match real provider behavior if vendor contracts change.
- **Trigger:** Provider updates signature algorithm or response shape.
- **Impact:** Users' code passes OpenMuara tests but fails against real providers.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Document sources in `REFERENCES.md`; regenerate golden files deliberately; keep contract tests focused on stable shapes.
- **Rollback Plan:** Mark provider contract as experimental; update golden files.
- **Monitoring:** Provider-specific test suites and smoke tests.
- **Status:** 🟡 Mitigated

### R04 — CI time increases with split jobs

- **Description:** Splitting CI into lint/unit/integration/smoke may increase wall-clock time if jobs are not parallelized.
- **Trigger:** Integration/smoke jobs are slow or sequential.
- **Impact:** Slower feedback on PRs.
- **Likelihood:** Medium
- **Impact Level:** Minor
- **Mitigation:** Run jobs in parallel; cache Go modules; run integration/smoke only on PRs or pushes to `dev`/`main`.
- **Rollback Plan:** Recombine jobs or reduce matrix.
- **Monitoring:** CI duration on `dev`.
- **Status:** 🟡 Mitigated

### R05 — Contributors resist new test conventions

- **Description:** New testing standards may feel heavy to casual contributors.
- **Trigger:** Contributor opens PR without reading `runbooks/testing.md`.
- **Impact:** Friction in community contributions.
- **Likelihood:** Low
- **Impact Level:** Minor
- **Mitigation:** Provide clear runbook, test templates, and helpful CI error messages.
- **Rollback Plan:** Simplify conventions based on feedback.
- **Monitoring:** Contributor feedback and PR review friction.
- **Status:** 🟡 Mitigated

---

## Rollback Playbook

If a prompt introduces a critical bug on `dev`:

1. **Stop:** Do not execute additional prompts.
2. **Identify:** Find the offending commit with `git log`.
3. **Assess:** Can it be fixed forward in <30 minutes? If yes, fix. If no, rollback.
4. **Rollback:** `git revert <commit-hash>` on `dev`.
5. **Verify:** `task check` and `task smoke`.
6. **Communicate:** Update `HANDOFF.md`, `TRACKING.md`, and `RISKS.md`.
7. **Resume:** Only continue after the rollback is verified and committed.
