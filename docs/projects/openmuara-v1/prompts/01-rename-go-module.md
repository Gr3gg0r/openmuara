> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

## Rename Go module to `github.com/openmuara/openmuara`

### Context
The project is rebranding from `muara` to `OpenMuara`. The Go module path is the first thing that must change because every package import depends on it.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Current module:** `github.com/openmuara/openmuara`
- **Target module:** `github.com/openmuara/openmuara`

### Scope
- **In scope:**
  - Update `go.mod` module path.
  - Update all `*.go` import paths.
  - Regenerate `go.sum` with `go mod tidy`.
- **Out of scope:**
  - Binary name, CLI commands, config paths (Prompt 02–03).
  - Documentation and scripts (Prompt 04).

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
grep -R "github.com/openmuara/openmuara" --include="*.go" -l | head -20
```

### Execution
1. Run `go mod edit -module github.com/openmuara/openmuara`.
2. Update all import paths. If `goimports` is available:
   ```bash
   goimports -w .
   ```
   Otherwise use `sed` or `gofmt` with find/replace.
3. Run `go mod tidy`.
4. Verify no `github.com/openmuara/openmuara` strings remain in `.go` files:
   ```bash
   grep -R "github.com/openmuara/openmuara" --include="*.go" .
   ```

### Quality Gates
```bash
task fmt
task vet
task lint
task test
task coverage
task build
task smoke
```

If `task` is unavailable, run the equivalent `go fmt`, `go vet`, `golangci-lint run`, `go test ./...`, `go build ./...` commands directly.

### Commit
```bash
git add -A
git commit -m "refactor(repo): rename Go module to github.com/openmuara/openmuara"
```

### Post-completion
1. Update `TRACKING.md` Step 01 → ✅, fill commit hash.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
