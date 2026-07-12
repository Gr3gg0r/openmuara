> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 04 — Make Fawry Escape Action Update Ledger and Verify Incoming Webhook Signatures

### Context

`POST /_admin/fawry-escape` dispatches a webhook but never updates the ledger, leaving transactions in `new`. Incoming `POST /fawry/webhook` accepts any payload and ignores `messageSignature`.

Read `tasks/04-fawry-escape-ledger-and-webhook-sigs.md` for the full spec.

### Current State

- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Target Files:** `internal/fawry/escape.go`, `internal/fawry/webhook.go`, `internal/fawry/signature.go`, tests

### Scope

- **In scope:**
  - Update the ledger transaction status in `NewEscapeActionHandler` based on the simulated outcome.
  - Add signature verification for incoming Fawry V2 webhooks.
  - Make verification configurable (skip if `webhook_secret` is empty).
- **Out of scope:**
  - Changing the Fawry charge request signature flow.
  - Adding new providers.

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./internal/fawry/...
```

### Execution

1. Read `tasks/04-fawry-escape-ledger-and-webhook-sigs.md`.
2. Implement ledger update in `NewEscapeActionHandler` using `engine.Transition`.
3. Implement or reuse signature verification for Fawry V2 webhook payload.
4. Update `NewWebhookHandler` to verify `messageSignature` when secret is configured.
5. Add tests for both behaviors.

### Quality Gates

```bash
go build ./...
go test ./internal/fawry/...
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```

### Commit

```bash
git add internal/fawry/
git commit -m "feat(fawry): update ledger on escape and verify incoming webhook signatures"
```

### Post-completion

1. Update `TRACKING.md` Step 04 to ✅.
2. Log decisions in `DECISIONS.md`.
3. Update `HANDOFF.md`.
