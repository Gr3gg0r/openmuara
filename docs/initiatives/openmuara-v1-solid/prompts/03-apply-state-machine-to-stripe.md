> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 03 — Apply State Machine to Stripe Simulation Handlers

### Context

Stripe simulation handlers (`internal/stripe/webhook.go` success handler and `internal/stripe/simulation.go` failure/cancel handlers) directly mutate `tx.Status`. They should use `engine.Transition` for consistency with the rest of the codebase.

### Current State

- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Target Files:** `internal/stripe/webhook.go`, `internal/stripe/simulation.go`, tests

### Scope

- **In scope:**
  - Replace direct `tx.Status = ...` assignments with `engine.Transition(&tx, targetStatus)`.
  - Update test fixtures to create sessions/transactions in valid source states.
  - Handle transition errors with appropriate HTTP status (409 Conflict).
- **Out of scope:**
  - Changing the state machine rules.
  - Modifying Stripe checkout flow.

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./internal/stripe/...
```

### Execution

1. Read `internal/stripe/webhook.go` and `internal/stripe/simulation.go`.
2. In success handler, use `engine.Transition(&tx, engine.TransactionStatusPaid)`.
3. In failure/cancel handlers, use `engine.Transition(&tx, engine.TransactionStatusUnpaid)`.
4. On transition error, return `httputil.ErrInvalidState` with 409.
5. Update `internal/stripe/simulation_test.go` and `internal/stripe/webhook_test.go` fixtures if needed.

### Quality Gates

```bash
go build ./...
go test ./internal/stripe/...
go test ./...
golangci-lint run
```

### Commit

```bash
git add internal/stripe/
git commit -m "refactor(stripe): use engine.Transition in simulation handlers"
```

### Post-completion

1. Update `TRACKING.md` Step 03 to ✅.
2. Log any auto-decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
