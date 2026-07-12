> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Step 04 — Update Docs and Tooling

> **Purpose:** Update all user-facing and CI-facing references from `muara` to `OpenMuara` / `muara`.
> **Related Prompt:** `prompts/04-update-docs-and-tooling.md`

---

## Objective

README, AGENTS.md, scripts, CI, pre-commit config, and any markdown docs must reflect the new name. Project name = OpenMuara, module = `github.com/openmuara/openmuara`, CLI = `muara`.

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | `README.md` | Modify | `<repo-root>/README.md` |
| 2 | `AGENTS.md` | Modify | `<repo-root>/AGENTS.md` |
| 3 | `.github/workflows/*.yml` | Modify | `<repo-root>/.github/workflows/` |
| 4 | `scripts/*.sh` | Modify | `<repo-root>/scripts/` |
| 5 | `.pre-commit-config.yaml` | Modify | `<repo-root>/.pre-commit-config.yaml` |
| 6 | Other markdown docs | Modify | `<repo-root>/docs/**/*.md` |

---

## Constraints & Security

- Do not change git history.
- Do not rename the GitHub org/repo unless explicitly told to.
- Update module references in docs to `github.com/openmuara/openmuara`.
- Update CLI references in docs/scripts to `muara`.

---

## Error Handling Requirements

- If a script fails after renaming, fix or document it.
- CI workflows must reference `./cmd/muara` and `muara` binary.

---

## BDD / TDD Quality Gates

- [ ] No remaining `muara` strings in `README.md`, `AGENTS.md`, CI, scripts, or pre-commit config.
- [ ] CI workflow still builds and tests successfully.
- [ ] `task smoke` passes.

---

## Rollback Trigger

- If CI cannot be fixed in <30 minutes, STOP and consult `RISKS.md` R01.
