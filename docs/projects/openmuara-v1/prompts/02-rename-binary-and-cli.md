> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

## Rename binary and CLI from `muara` to `muara`

### Context
After the module rename, the executable and CLI commands must use the shorter command name. The project name remains `OpenMuara` and the module path remains `github.com/openmuara/openmuara`, but users will type `muara`.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Last step:** Prompt 01 completed (module rename).

### Scope
- **In scope:**
  - Rename `cmd/openmuara/` → `cmd/muara/`.
  - Update CLI command names and help text in `internal/cli/`.
  - Update `Taskfile.yml` to build `./cmd/muara` and name the binary `muara`.
  - Rename error code prefix to `OPENMUARA_`.
- **Out of scope:**
  - Module path (already done in Prompt 01).
  - Config paths (Prompt 03).
  - Documentation (Prompt 04).

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
grep -R "openmuara" --include="*.go" cmd/ internal/cli/ internal/httputil/errors.go | head -30
```

### Execution
1. Rename directory: `git mv cmd/openmuara cmd/muara`.
2. Update package declaration and imports inside `cmd/muara/`.
3. Update `Taskfile.yml` to build `./cmd/muara` and name the binary `muara`.
4. Update CLI strings in `internal/cli/` (command names, help, version).
5. Update error code constants in `internal/httputil/errors.go` to use the `OPENMUARA_` prefix.
6. Run tests for affected packages.

### Quality Gates
```bash
task fmt
task vet
task lint
task test ./cmd/muara/... ./internal/cli/... ./internal/httputil/...
task coverage
task build
task smoke
```

### Commit
```bash
git add -A
git commit -m "refactor(cli): rename binary and CLI commands to muara"
```

### Post-completion
1. Update `TRACKING.md` Step 02 → ✅, fill commit hash.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
