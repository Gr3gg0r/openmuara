> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1 — Risk Register & Rollback Plans

> **Purpose:** Track what could go wrong and what we do if it does.

---

## Risk Matrix

| ID | Risk | Likelihood | Impact | Status | Owner |
|----|------|------------|--------|--------|-------|
| R01 | Rebrand breaks existing tests or scripts | Medium | Major | ✅ Mitigated | AI Agent |
| R02 | SQLite migration loses in-memory test coverage | Medium | Major | ✅ Mitigated | AI Agent |
| R03 | Provider emulation drifts from real provider contract | Medium | Critical | ✅ Mitigated | AI Agent |
| R04 | Coverage drops below 50% threshold | Low | Minor | ✅ Mitigated | AI Agent |
| R05 | Direct `dev` commits destabilize branch | Medium | Major | ✅ Mitigated | AI Agent |
| R06 | Metrics endpoint exposes internal counts unauthenticated | Low | Minor | ✅ Mitigated | AI Agent |
| R07 | Audit log grows unbounded and consumes disk | Low | Major | 🟡 Mitigated | AI Agent |
| R08 | SQLite single-writer contention under load | Medium | Major | ✅ Mitigated | AI Agent |
| R09 | Webhook replay flood overwhelms consumer | Low | Major | 🟡 Mitigated | AI Agent |
| R10 | Provider config drift after release | Medium | Critical | 🟡 Mitigated | AI Agent |

---

## Detailed Risk Entries

### R01 — Rebrand breaks existing tests or scripts

- **Description:** Renaming module paths, binary, and config paths could break tests, scripts, CI, or smoke tests.
- **Trigger:** Global string replacement misses references or changes behavior.
- **Impact:** Quality gates fail; `dev` branch becomes red.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Use `go mod edit` and `goimports` for import paths; run `task check` and `task smoke` after each rebrand prompt; update scripts/CI in dedicated prompt.
- **Rollback Plan:** Revert the rebrand commit(s) or use `git checkout <last-good-commit> -- <files>`.
- **Monitoring:** `task check` and `task smoke` must pass after every rebrand step.
- **Status:** ✅ Mitigated

### R02 — SQLite migration loses in-memory test coverage

- **Description:** Replacing in-memory stores with SQLite could make unit tests slower or more complex.
- **Trigger:** SQLite store not tested adequately; tests require filesystem setup.
- **Impact:** Coverage drops; CI slows down; tests become flaky.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Keep in-memory store as an option for tests; add `:memory:` SQLite tests; ensure `task test` still passes quickly.
- **Rollback Plan:** Revert SQLite commit and keep in-memory store.
- **Monitoring:** `task coverage` and `task test` timing.
- **Status:** ✅ Mitigated

### R03 — Provider emulation drifts from real provider contract

- **Description:** Stripe/SenangPay adapters may not match real request/response shapes.
- **Trigger:** Incomplete vendor docs or incorrect assumptions.
- **Impact:** Users' code passes OpenMuara tests but fails against real providers.
- **Likelihood:** Medium
- **Impact Level:** Critical
- **Mitigation:** Reference official API docs in `REFERENCES.md`; add contract tests; document known limitations.
- **Rollback Plan:** Mark provider as experimental; fix forward.
- **Monitoring:** Contract tests and documentation.
- **Status:** ✅ Mitigated

### R04 — Coverage drops below 50% threshold

- **Description:** New features may not have tests, dropping coverage below the gate.
- **Trigger:** New code added without corresponding tests.
- **Impact:** Quality gate fails; merge blocked.
- **Likelihood:** Low
- **Impact Level:** Minor
- **Mitigation:** Every prompt includes test requirements; run `task coverage` before commit.
- **Rollback Plan:** Add missing tests before proceeding.
- **Monitoring:** `task coverage` output.
- **Status:** ✅ Mitigated

### R05 — Direct `dev` commits destabilize branch

- **Description:** Working directly on `dev` without feature branches could introduce regressions that affect other work.
- **Trigger:** A broken commit lands on `dev`.
- **Impact:** `dev` is red; other developers blocked.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Run full quality gates before every commit; keep commits small and focused; use `HANDOFF.md` to track state.
- **Rollback Plan:** `git revert <commit-hash>` on `dev`.
- **Monitoring:** CI status on `dev` after push.
- **Status:** ✅ Mitigated

### R06 — Metrics endpoint exposes internal counts unauthenticated

- **Description:** `/metrics` returns counters and histograms without authentication.
- **Trigger:** OpenMuara is exposed beyond localhost.
- **Impact:** Internal request counts become visible on the network; no PII or payload data is exposed.
- **Likelihood:** Low
- **Impact Level:** Minor
- **Mitigation:** Bind to `127.0.0.1` by default; protect with a reverse proxy in shared environments.
- **Rollback Plan:** Disable metrics endpoint or add authentication at the proxy layer.
- **Monitoring:** Network exposure review.
- **Status:** ✅ Mitigated

### R07 — Audit log grows unbounded and consumes disk

- **Description:** The `audit_logs` table appends every logged event without rotation.
- **Trigger:** High event volume or long-lived instances.
- **Impact:** SQLite file grows and may exhaust disk space.
- **Likelihood:** Low
- **Impact Level:** Major
- **Mitigation:** Document operational limitation; operators can archive/prune old rows.
- **Rollback Plan:** Truncate old audit records or switch to `memory` persistence.
- **Monitoring:** Disk usage of `.muara/data/ledger.db`.
- **Status:** 🟡 Mitigated

### R08 — SQLite single-writer contention under load

- **Description:** SQLite serializes writes through a single shared connection.
- **Trigger:** Parallel load tests or bursts of concurrent transactions.
- **Impact:** `database is locked` errors and 500 responses.
- **Likelihood:** Medium
- **Impact Level:** Major
- **Mitigation:** Share one `*sql.DB` between transaction and audit stores; document single-writer limits.
- **Rollback Plan:** Switch to `memory` persistence or reduce concurrency.
- **Monitoring:** 5xx rate and logs for `database is locked`.
- **Status:** ✅ Mitigated

### R09 — Webhook replay flood overwhelms consumer

- **Description:** Bulk replay of failed webhooks can generate a sudden burst of traffic.
- **Trigger:** Operator clicks replay-all or runs a replay script.
- **Impact:** Consumer service is overloaded.
- **Likelihood:** Low
- **Impact Level:** Major
- **Mitigation:** Replay individual attempts first; verify consumer health before large replays.
- **Rollback Plan:** Stop replays and restart OpenMuara if needed.
- **Monitoring:** `openmuara_webhook_attempts_total` rate.
- **Status:** 🟡 Mitigated

### R10 — Provider config drift after release

- **Description:** Provider signatures, payloads, or config keys may diverge from real providers over time.
- **Trigger:** Vendor changes contracts or OpenMuara docs become stale.
- **Impact:** Users' code passes OpenMuara tests but fails against real providers.
- **Likelihood:** Medium
- **Impact Level:** Critical
- **Mitigation:** Maintain `docs/providers.md`, `docs/openapi.yaml`, and contract tests; version provider plugins via `plugins/*/gateway.yml`.
- **Rollback Plan:** Pin to previous release or update provider config to match vendor docs.
- **Monitoring:** Contract tests, provider-specific test suites, smoke tests.
- **Status:** 🟡 Mitigated

---

## Rollback Playbook

If a step introduces a critical bug on `dev`:

1. **Stop:** Do not execute additional prompts until the issue is contained.
2. **Identify:** Determine which commit introduced the bug (use `git log` and `git bisect` if needed).
3. **Assess:** Can the bug be fixed forward in <30 minutes? If yes, fix. If no, rollback.
4. **Rollback:** On `dev`, run `git revert <commit-hash>` or `git checkout <last-good-commit> -- <files>`.
5. **Verify:** Run `task check` and `task smoke`.
6. **Communicate:** Update `HANDOFF.md`, `TRACKING.md`, and `RISKS.md` with what happened.
7. **Resume:** Only continue after the rollback is verified and committed.
