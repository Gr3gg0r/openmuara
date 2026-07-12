> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This template is subordinate to it.**

# Example Prompt

This is an example of a good prompt. Delete this file after initializing the project.

---

## Add request body size limit to Fawry charge handler

### Context
The Fawry charge endpoint currently reads request bodies without a size limit. This creates a DoS vector and should be hardened to match the rest of the server.

### Current State
- **Repo:** `<repo-root>`
- **Branch:** `dev`
- **Target files:** `internal/fawry/charge.go`, `internal/httputil/body.go`

### Scope
- **In scope:**
  - Use `httputil.ReadJSONBody` in `internal/fawry/charge.go`.
  - Ensure `ReadJSONBody` uses `http.MaxBytesReader`.
- **Out of scope:**
  - Other providers (handled in later prompts).
  - Changing the response schema.

### Pre-flight
```bash
cd <repo-root>
git status
git branch --show-current  # must be dev
task check
```

### Execution
1. Open `internal/httputil/body.go` and verify `ReadJSONBody` wraps `r.Body` with `http.MaxBytesReader`.
2. Open `internal/fawry/charge.go` and replace direct JSON decoding with `httputil.ReadJSONBody`.
3. Add a test that sends an oversized body and expects `413 Request Entity Too Large`.

### Quality Gates
```bash
task fmt
task vet
task lint
task test ./internal/fawry/... ./internal/httputil/...
task coverage
task build
task smoke
```

### Commit
```bash
git add -A
git commit -m "feat(fawry): limit charge request body size via MaxBytesReader"
```

### Post-completion
1. Update `TRACKING.md` Step 03 → ✅, fill commit hash.
2. Log any decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
