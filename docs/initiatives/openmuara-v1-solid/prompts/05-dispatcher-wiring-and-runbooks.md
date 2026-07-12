> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 05 — Improve Dispatcher Wiring and Update Runbooks

### Context

`cli/start.go` builds a dispatcher for every enabled provider but only wires the active provider's dispatcher into the router. This means provider-specific webhook targets may not be honored consistently. Additionally, runbooks and README need updating for new features.

### Current State

- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Target Files:** `internal/cli/start.go`, `runbooks/local-development.md`, `runbooks/quality-gates.md`, `README.md`

### Scope

- **In scope:**
  - Review dispatcher wiring in `cli/start.go`; ensure each provider uses its own dispatcher for its own webhooks.
  - Update `runbooks/local-development.md` with `/readyz`, docker-compose, and pagination notes.
  - Update `runbooks/quality-gates.md` if needed.
  - Update `README.md` to mention `/readyz`, paginated admin APIs, and body limit.
- **Out of scope:**
  - New features.
  - Changes to provider payload builders.

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./internal/cli/...
./scripts/smoke-test.sh
```

### Execution

1. Read `internal/cli/start.go` dispatcher wiring.
2. Decide whether to pass the full `providerDispatchers` map into `RouterConfig` or keep the current active-only pattern.
3. Update runbooks and README.
4. Run smoke test to confirm no regression.

### Quality Gates

```bash
go build ./...
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```

### Commit

```bash
git add internal/cli/start.go runbooks/ README.md
git commit -m "docs(cli): improve dispatcher wiring and update runbooks"
```

### Post-completion

1. Update `TRACKING.md` Step 05 to ✅.
2. Log decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
