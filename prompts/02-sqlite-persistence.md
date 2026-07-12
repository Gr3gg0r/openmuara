# Prompt 02 — SQLite Persistence Layer

## Goal
Replace in-memory stores with SQLite-backed persistence for transactions, webhooks, and provider state.

## Acceptance Criteria
- [ ] SQLite schema created for:
  - `transactions` (id, provider, type, amount, currency, status, reference, idempotency_key, payload, created_at, updated_at)
  - `webhook_attempts` (id, ref, url, status, payload, headers, attempts, last_error, created_at, updated_at)
  - `checkout_sessions` (id, provider, session_data, created_at, updated_at)
- [ ] Migration system using `golang-migrate` or simple `.sql` files in `internal/store/migrations/`
- [ ] `internal/store` package with repository pattern
- [ ] `engine.TransactionStore` backed by SQLite
- [ ] `webhook.AttemptStore` backed by SQLite
- [ ] Thread-safe access
- [ ] Build and tests pass

## Files to Create/Change
- `internal/store/db.go` — connection + migrate
- `internal/store/migrations/*.sql`
- `internal/store/transaction_repo.go`
- `internal/store/webhook_repo.go`
- `internal/engine/store.go` — delegate to SQLite
- `internal/webhook/store.go` — delegate to SQLite
- `internal/config/config.go` — `persistence.path` default `.muara/data/ledger.db`

## Response Shape
Return:
1. Schema diagram or table list
2. Migration file names
3. Repository interface signatures
4. Test coverage summary

## Test Notes
- `go test ./internal/store/...`
- Verify idempotency across restarts
- Verify concurrent writes do not race
