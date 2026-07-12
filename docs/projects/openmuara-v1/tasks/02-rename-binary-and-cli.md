> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Step 02 — Rename Binary and CLI

> **Purpose:** Rename the binary from `muara` to `muara` and update CLI command names, defaults, and error codes.
> **Related Prompt:** `prompts/02-rename-binary-and-cli.md`

---

## Objective

Developers run `muara start` instead of `openmuara start`. All CLI help text, version output, and internal references must match. The project name remains `OpenMuara` and the module path remains `github.com/openmuara/openmuara`.

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | `cmd/openmuara/` | Rename / Modify | `<repo-root>/cmd/muara/` |
| 2 | `Taskfile.yml` | Modify | `<repo-root>/Taskfile.yml` |
| 3 | `internal/cli/*.go` | Modify | `<repo-root>/internal/cli/` |
| 4 | Error code constants | Modify | `<repo-root>/internal/httputil/errors.go` |

---

## Constraints & Security

- Do not change command semantics — only names.
- Error codes that start with the legacy prefix should become `OPENMUARA_`.
- Keep backward compatibility only if explicitly requested later.
- Module path stays `github.com/openmuara/openmuara`.

---

## Error Handling Requirements

- CLI help text must show `muara` consistently.
- If a command name collision occurs, log and resolve before commit.

---

## BDD / TDD Quality Gates

- [ ] `go build ./cmd/muara` produces `muara` binary.
- [ ] `./muara --help` shows `muara` in usage.
- [ ] `./muara version` reports `muara` (and OpenMuara version).
- [ ] `go test ./internal/cli/...` passes.
- [ ] No remaining `muara` strings in `cmd/`, `internal/cli/`, or error code constants.

---

## Rollback Trigger

- If CLI tests fail after renaming and cannot be fixed in <30 minutes, STOP and consult `RISKS.md` R01.
