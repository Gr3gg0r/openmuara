> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Step 05 — Run Rebrand Quality Gates

> **Purpose:** Verify the entire rebrand is complete and the `dev` branch is green.
> **Related Prompt:** `prompts/05-run-rebrand-quality-gates.md`

---

## Objective

Catch any remaining `muara` references and confirm all quality gates pass before moving to v1 feature work.

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | Entire repo | Verify | `<repo-root>` |

---

## Constraints & Security

- Do not fix unrelated pre-existing bugs.
- If a gate fails due to the rebrand, fix it. If it fails for an unrelated reason, log in `KNOWN_ISSUES.md`.

---

## Error Handling Requirements

- Any remaining `muara` string must be reviewed and either renamed or explicitly documented as intentional.

---

## BDD / TDD Quality Gates

- [ ] `task check` passes (or equivalent `fmt`, `vet`, `lint`, `test`, `coverage`, `build`, `smoke`).
- [ ] No unintentional `muara` strings remain in source, docs, scripts, or CI.
- [ ] Binary builds as `muara`.
- [ ] `muara --help` and `muara version` work.

---

## Rollback Trigger

- If any gate fails and cannot be fixed in <30 minutes, STOP and consult `RISKS.md` R01.
