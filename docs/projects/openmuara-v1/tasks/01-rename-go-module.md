> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# Step 01 — Rename Go Module

> **Purpose:** Change the Go module path from `github.com/openmuara/openmuara` to `github.com/openmuara/openmuara` and update all import paths.
> **Related Prompt:** `prompts/01-rename-go-module.md`

---

## Objective

All Go code must import the new module path. This is the foundation of the rebrand; later steps assume the module path is correct.

---

## Target Files

| # | File | Action | Repo Path |
|---|------|--------|-----------|
| 1 | `go.mod` | Modify | `<repo-root>/go.mod` |
| 2 | All `*.go` files with `github.com/openmuara/openmuara` imports | Modify | `<repo-root>/**/*.go` |
| 3 | `go.sum` | Regenerate | `<repo-root>/go.sum` |

---

## Constraints & Security

- Do not change external dependencies unless required by the rename.
- Do not alter business logic — only import paths.
- Keep commit focused on module rename only.

---

## Error Handling Requirements

- If `go mod tidy` fails, stop and investigate dependency issues.
- If `goimports` or `gofmt` reports syntax errors, fix before commit.

---

## BDD / TDD Quality Gates

- [ ] `go mod edit -module github.com/openmuara/openmuara` succeeds.
- [ ] `goimports -w .` updates all import paths.
- [ ] `go mod tidy` succeeds.
- [ ] `go build ./...` compiles.
- [ ] `go test ./...` passes.
- [ ] `task fmt` and `task vet` pass (or `gofmt` / `go vet` if `task` unavailable).
- [ ] No remaining `github.com/openmuara/openmuara` strings in `.go` files.

---

## Rollback Trigger

- If `go build ./...` fails after 2 fix attempts, STOP and consult `RISKS.md` R01.
