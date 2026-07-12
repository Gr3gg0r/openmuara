> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid — Risk Register & Rollback Plans

> **Purpose:** Track what could go wrong and what to do if it does.

---

## Risk Matrix

| ID | Risk | Likelihood | Impact | Status | Owner |
|----|------|------------|--------|--------|-------|
| R01 | OpenAPI spec drift confuses API consumers | Medium | Minor | 🟡 Mitigated | AI Agent |
| R02 | State machine change breaks existing test fixtures | Low | Minor | 🟡 Mitigated | AI Agent |
| R03 | Fawry webhook signature verification rejects legitimate payloads | Low | Major | 🟡 Mitigated | AI Agent |
| R04 | Dashboard UI remains broken after pagination fix | Low | Minor | 🟡 Mitigated | AI Agent |

**Likelihood:** Rare / Low / Medium / High / Certain
**Impact:** Negligible / Minor / Major / Critical / Catastrophic
**Status:** ⬜ Open / 🟡 Mitigated / ✅ Closed / ❌ Realized

---

## Detailed Risk Entries

### R01 — OpenAPI spec drift

- **Description:** The OpenAPI spec does not match the actual API after recent changes.
- **Trigger:** Consumers generate clients from the spec and hit mismatches.
- **Impact:** Minor — local emulator; no external consumers yet.
- **Likelihood:** Medium
- **Mitigation:** Step 02 explicitly syncs the spec and adds a sync test.
- **Rollback Plan:** Revert the spec commit; the API behavior remains unchanged.
- **Monitoring:** Run `go test ./internal/server/...` and `./scripts/smoke-test.sh`.
- **Status:** 🟡 Mitigated

### R02 — State machine breaks test fixtures

- **Description:** Stripe simulation handlers currently set status directly. Moving to `engine.Transition` may reject invalid transitions in tests.
- **Trigger:** Test creates a session in an unexpected state.
- **Impact:** Minor — tests fail until fixtures are adjusted.
- **Likelihood:** Low
- **Mitigation:** Update test fixtures to create sessions in valid source states.
- **Rollback Plan:** Revert the state-machine commit.
- **Monitoring:** `go test ./internal/stripe/...`
- **Status:** 🟡 Mitigated

### R03 — Fawry webhook signature verification rejects legitimate payloads

- **Description:** Adding signature verification to incoming Fawry webhooks may reject valid payloads if the signing logic differs from real Fawry.
- **Trigger:** Signature computation mismatch.
- **Impact:** Major — legitimate webhook payloads would fail.
- **Likelihood:** Low
- **Mitigation:** Implement verification as optional/configurable; default to verifying only when `webhook_secret` is set; document the canonical string.
- **Rollback Plan:** Disable verification via config or revert the handler change.
- **Monitoring:** `go test ./internal/fawry/...` and smoke test.
- **Status:** 🟡 Mitigated

### R04 — Dashboard UI remains broken after pagination

- **Description:** The admin dashboard still expects array responses and shows no data.
- **Trigger:** Loading `/_admin` after pagination change.
- **Impact:** Minor — dashboard unusable.
- **Likelihood:** Low
- **Mitigation:** Step 01 fixes the UI and adds a UI-level test for response shape.
- **Rollback Plan:** Revert the UI change.
- **Monitoring:** Manual dashboard load + `go test ./internal/ui/...`
- **Status:** 🟡 Mitigated

---

## Rollback Playbook

If a step introduces a critical bug:

1. **Stop:** Do not execute additional prompts.
2. **Identify:** Use `git log` to find the offending commit.
3. **Assess:** Can it be fixed forward in <30 minutes? If yes, fix. If no, rollback.
4. **Rollback:** `git revert <commit-hash>` on `dev`.
5. **Verify:** Run `go build ./...`, `go test ./...`, `./scripts/smoke-test.sh`.
6. **Communicate:** Update `HANDOFF.md`, `TRACKING.md`, and `RISKS.md`.
7. **Resume:** Continue only after rollback is verified.

---

## How to Add a New Risk

1. Pick the next number: `R##`.
2. Add a row to the Risk Matrix.
3. Add a Detailed Risk Entry.
4. If mitigation references a runbook, create a stub in `runbooks/`.
