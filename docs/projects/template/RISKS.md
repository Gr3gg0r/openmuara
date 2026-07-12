> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# <Project Name> — Risk Register & Rollback Plans

> **Purpose:** Track what could go wrong, how likely it is, and what we do if it does.
> **Rule:** Every project with >3 steps MUST have at least one risk logged before execution begins.

---

## Risk Matrix

| ID | Risk | Likelihood | Impact | Status | Owner |
|----|------|------------|--------|--------|-------|
| R01 | | | | ⬜ Open | |
| R02 | | | | ⬜ Open | |

**Likelihood:** Rare / Low / Medium / High / Certain
**Impact:** Negligible / Minor / Major / Critical / Catastrophic
**Status:** ⬜ Open / 🟡 Mitigated / ✅ Closed / ❌ Realized

---

## Detailed Risk Entries

### R01 — [Short Risk Title]

- **Description:** What could go wrong?
- **Trigger:** What event or condition would cause this?
- **Impact:** What happens to users, data, or the project?
- **Likelihood:** __________
- **Impact Level:** __________
- **Mitigation:** What are we doing to prevent this?
- **Rollback Plan:** If this happens, what is the exact recovery procedure?
  1. Step 1
  2. Step 2
  3. Step 3
- **Monitoring:** How will we know if this risk is materializing?
- **Status:** ⬜ Open / 🟡 Mitigated / ✅ Closed / ❌ Realized

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

---

## How to Add a New Risk

1. Pick the next number: `R##`.
2. Add a row to the Risk Matrix.
3. Add a Detailed Risk Entry using the R01 template.
4. If the risk has a rollback procedure, add it to the Rollback Playbook section.
5. **If the mitigation references a runbook that does not exist yet, create a stub in `runbooks/` before finishing plan creation.**
