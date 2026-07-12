> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Prompt Template

Copy this file for each prompt in the initiative.
Replace all `{{ }}` placeholders.

---

## {{ VERB }} {{ NOUN }}

### Context
{{ 1–2 sentence motivation. Why this step? What changed upstream? }}

### Current State
- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Last Commit:** `{{ abc123 }}` (if any)
- **Modified Files:** `{{ path/to/file1.go, path/to/file2.go }}`

### Scope
- **In scope:** {{ list }}
- **Out of scope:** {{ list }}

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./...              # must pass before changes
```

If preflight fails, STOP and fix before proceeding.

### Execution

```bash
# Step 1
{{ commands }}
```

```bash
# Step 2
{{ commands }}
```

### Quality Gates

After completing, run:

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
./scripts/smoke-test.sh
```

All gates must pass before committing.

### Commit

```bash
git add -A
git commit -m "{{ type(scope): verb noun }}"
```

### Post-completion

1. Update `TRACKING.md` → mark step complete, fill commit hash.
2. If you made a non-trivial decision, log it in `DECISIONS.md`.
3. Update `HANDOFF.md` with status and next step.
