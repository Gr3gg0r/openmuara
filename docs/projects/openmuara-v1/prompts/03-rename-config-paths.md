> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

## Rename config and data paths from `.muara/` to `.muara/`

### Context
After renaming the module and binary, the runtime config and data paths must match the shorter CLI name. This affects where `muara init` writes files and where the server looks for config.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Last step:** Prompt 02 completed (binary/CLI rename).

### Scope
- **In scope:**
  - Update default config directory in `internal/config/config.go` to `.muara/`.
  - Update `internal/config/config_test.go`.
  - Update `internal/cli/init.go` and any other CLI commands that reference paths.
  - Update any hardcoded `.muara/` strings in Go source.
- **Out of scope:**
  - Documentation and scripts (Prompt 04).
  - Migration of existing user `.muara/` directories (document only).

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
grep -R "\.openmuara" --include="*.go" . | head -30
```

### Execution
1. Update `internal/config/config.go` to use `.muara/` as default config/data dir.
2. Update `internal/config/config_test.go` accordingly (use temp dirs in tests).
3. Update `internal/cli/init.go` and any other path references.
4. Search for remaining `.muara/` strings in Go source and update.

### Quality Gates
```bash
task fmt
task vet
task lint
task test ./internal/config/... ./internal/cli/...
task coverage
task build
task smoke
```

### Commit
```bash
git add -A
git commit -m "refactor(config): rename default config paths to .muara/"
```

### Post-completion
1. Update `TRACKING.md` Step 03 → ✅, fill commit hash.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
