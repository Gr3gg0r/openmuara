> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P06 — Transaction Search and Replay

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

Make the transactions table searchable and let users replay a charge or webhook for any transaction.

## Why now

As soon as a user starts testing, the transactions table fills up. Finding a specific record and re-triggering its webhook currently requires manual API calls.

## Scope

### In scope

- Add `GET /_admin/transactions/{ref}` to return the full transaction record, including status timeline if available.
- Add query parameters to `GET /_admin/transactions`:
  - `q` — free-text search across reference, provider, and status.
  - `provider` — filter by provider name.
  - `status` — filter by status.
- Add `POST /_admin/transactions/{ref}/replay-webhook` to re-emit the webhook for that transaction.
- Update `internal/ui/index.html` with a search bar and provider/status filters.
- Add a "Replay webhook" action per transaction row and a click-to-inspect detail panel.
- Add tests for detail, search, filter, and replay.

### Out of scope

- Full pagination redesign.
- Bulk operations.

## Acceptance criteria

- [ ] `GET /_admin/transactions` supports `q`, `provider`, and `status` filters.
- [ ] `POST /_admin/transactions/{ref}/replay-webhook` re-emits the webhook.
- [ ] Dashboard has search/filter controls and a replay action.
- [ ] Tests cover filters and replay.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Search can be done in-memory against the existing transaction store for now.
- Replay can reuse the webhook dispatcher; ensure idempotency if the original webhook already succeeded.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated `internal/server/admin_api_test.go`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
