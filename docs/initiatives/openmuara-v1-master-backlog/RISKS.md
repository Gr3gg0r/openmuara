> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Master Backlog — Risk Register

> **Purpose:** Track what could go wrong across the v1 backlog and how to respond.

---

## Risk Matrix

| ID | Risk | Priority | Likelihood | Impact | Status | Owner |
|----|------|----------|------------|--------|--------|-------|
| R01 | OpenAPI spec drift confuses API consumers | High | Medium | Minor | ✅ Mitigated | AI Agent |
| R02 | State machine change breaks existing test fixtures | High | Low | Minor | ✅ Mitigated | AI Agent |
| R03 | Fawry webhook signature verification rejects legitimate payloads | High | Low | Major | ✅ Mitigated | AI Agent |
| R04 | Dashboard UI remains broken after pagination fix | High | Low | Minor | ✅ Mitigated | AI Agent |
| R05 | `muara` rebrand left incomplete | High | Medium | Major | ✅ Mitigated | AI Agent |
| R06 | Scope creep into v2-frozen providers | Medium | Low | Major | ✅ Mitigated | Human |
| R07 | Metrics endpoint exposes internal counts unauthenticated | Low | Low | Minor | ✅ Mitigated | AI Agent |
| R08 | Audit log grows unbounded and consumes disk | Low | Low | Major | 🟡 Mitigated | AI Agent |
| R09 | SQLite single-writer contention under load | Medium | Medium | Major | ✅ Mitigated | AI Agent |
| R10 | Webhook replay flood overwhelms consumer | Low | Low | Major | 🟡 Mitigated | AI Agent |
| R11 | Provider config drift after release | Medium | Medium | Critical | 🟡 Mitigated | AI Agent |

---

## Detailed Risk Entries

### R01 — OpenAPI spec drift

- **Trigger:** Consumers generate clients from the spec and hit mismatches.
- **Mitigation:** S02 explicitly syncs the spec and adds a sync test.
- **Rollback:** Revert the spec commit; API behavior unchanged.
- **Monitoring:** `go test ./internal/server/...`, `./scripts/smoke-test.sh`.
- **Status:** ✅ Mitigated

### R02 — State machine breaks test fixtures

- **Trigger:** Test creates a session in an unexpected state.
- **Mitigation:** Update fixtures to valid source states in S03.
- **Rollback:** Revert the state-machine commit.
- **Monitoring:** `go test ./internal/stripe/...`.
- **Status:** ✅ Mitigated

### R03 — Fawry webhook signature verification rejects valid payloads

- **Trigger:** Signing logic differs from real Fawry.
- **Mitigation:** Make verification optional; verify only when `webhook_secret` is set.
- **Rollback:** Disable verification via config or revert handler change.
- **Monitoring:** `go test ./internal/fawry/...`, smoke test.
- **Status:** ✅ Mitigated

### R04 — Dashboard UI remains broken after pagination

- **Trigger:** `/_admin` loads and still expects arrays.
- **Mitigation:** S01 fixes UI and adds a response-shape test.
- **Rollback:** Revert the UI change.
- **Monitoring:** `go test ./internal/ui/...`, manual dashboard load.
- **Status:** ✅ Mitigated

### R05 — `muara` rebrand left incomplete

- **Trigger:** Module paths, binary names, or config paths still reference `muara`.
- **Mitigation:** P01 finishes the rebrand; grep for remaining `muara` refs.
- **Rollback:** Revert the rebrand commit if it breaks builds.
- **Monitoring:** `go build ./...`, `go test ./...`.
- **Status:** ✅ Mitigated

### R06 — Scope creep into v2-frozen providers

- **Trigger:** Agent starts implementing App Store / Play Store / RevenueCat.
- **Mitigation:** `KNOWN_ISSUES.md` lists hard boundaries; require human approval.
- **Rollback:** Delete out-of-scope code and update tracker.
- **Monitoring:** Review every PR for provider scope.
- **Status:** ✅ Mitigated

### R07 — Metrics endpoint exposes internal counts unauthenticated

- **Trigger:** OpenMuara is exposed beyond localhost.
- **Impact:** Internal request counts become visible on the network; no PII or payload data is exposed.
- **Mitigation:** Bind to `127.0.0.1` by default; protect with a reverse proxy in shared environments.
- **Rollback:** Disable metrics endpoint or add authentication at the proxy layer.
- **Monitoring:** Network exposure review.
- **Status:** ✅ Mitigated

### R08 — Audit log grows unbounded and consumes disk

- **Trigger:** High event volume or long-lived instances.
- **Impact:** SQLite file grows and may exhaust disk space.
- **Mitigation:** Document operational limitation; operators can archive/prune old rows.
- **Rollback:** Truncate old audit records or switch to `memory` persistence.
- **Monitoring:** Disk usage of `.muara/data/ledger.db`.
- **Status:** 🟡 Mitigated

### R09 — SQLite single-writer contention under load

- **Trigger:** Parallel load tests or bursts of concurrent transactions.
- **Impact:** `database is locked` errors and 500 responses.
- **Mitigation:** Share one `*sql.DB` between transaction and audit stores; document single-writer limits.
- **Rollback:** Switch to `memory` persistence or reduce concurrency.
- **Monitoring:** 5xx rate and logs for `database is locked`.
- **Status:** ✅ Mitigated

### R10 — Webhook replay flood overwhelms consumer

- **Trigger:** Operator clicks replay-all or runs a replay script.
- **Impact:** Consumer service is overloaded.
- **Mitigation:** Replay individual attempts first; verify consumer health before large replays.
- **Rollback:** Stop replays and restart OpenMuara if needed.
- **Monitoring:** `openmuara_webhook_attempts_total` rate.
- **Status:** 🟡 Mitigated

### R11 — Provider config drift after release

- **Trigger:** Vendor changes contracts or OpenMuara docs become stale.
- **Impact:** Users' code passes OpenMuara tests but fails against real providers.
- **Mitigation:** Maintain `docs/providers.md`, `docs/openapi.yaml`, and contract tests; version provider plugins via `plugins/*/gateway.yml`.
- **Rollback:** Pin to previous release or update provider config to match vendor docs.
- **Monitoring:** Contract tests, provider-specific test suites, smoke tests.
- **Status:** 🟡 Mitigated

---

## Rollback Playbook

1. **Stop:** Do not execute additional items.
2. **Identify:** Use `git log` to find the offending commit.
3. **Assess:** Can it be fixed forward in <30 minutes? If yes, fix. If no, rollback.
4. **Rollback:** `git revert <commit-hash>` on `dev`.
5. **Verify:** `go build ./...`, `go test ./...`, `./scripts/smoke-test.sh`.
6. **Communicate:** Update `HANDOFF.md`, `TRACKING.md`, and `RISKS.md`.
