> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Step 03 — Rename Config and Data Paths

> **Purpose:** Update default config and data directories from `.muara/` to `.muara/`.
> **Related Prompt:** `prompts/03-rename-config-paths.md`

---

## Objective

User-facing file system paths must match the shorter CLI name `muara`. This keeps the runtime behavior consistent with the command users type.

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | `internal/config/config.go` | Modify | `<repo-root>/internal/config/config.go` |
| 2 | `internal/config/config_test.go` | Modify | `<repo-root>/internal/config/config_test.go` |
| 3 | CLI init command | Modify | `<repo-root>/internal/cli/init.go` |
| 4 | Any hardcoded `.muara/` strings | Modify | `<repo-root>/**/*.go` |

---

## Constraints & Security

- Do not delete user data on disk — only change default paths.
- If a migration helper is needed, document it; do not auto-delete old `.muara/` directories.
- Tests must use temp directories, not real `~/.muara/`.

---

## Error Handling Requirements

- If `muara init` cannot create `~/.muara/`, return a clear error.
- Tests must not depend on the presence of `~/.muara/`.

---

## BDD / TDD Quality Gates

- [ ] `internal/config/config.go` uses `.muara/` as default config/data dir.
- [ ] `internal/config/config_test.go` updated and passes.
- [ ] `muara init` creates `~/.muara/`.
- [ ] No remaining `.muara/` strings in Go source.

---

## Rollback Trigger

- If config tests fail due to cross-platform path issues, STOP and consult `RISKS.md` R01.
