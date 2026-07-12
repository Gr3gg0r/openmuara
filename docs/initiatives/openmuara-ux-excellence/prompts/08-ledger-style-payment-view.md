> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P08 — Ledger-Style Payment View

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** P02, P05, P06 (transaction detail endpoint)

---

## Goal

Turn the dashboard into a ledger-style view for payment traffic: one screen where testers see transactions and webhooks arrive in real time, inspect payloads, and replay events. Inspired by Mailpit's email inbox, but named and shaped around OpenMuara's transaction ledger.

## Why now

Testers do not want to read logs or write scripts. They want to open a web page, see traffic flow in, click a row, and understand what happened. Mailpit does this for email; OpenMuara should do it for payments through its ledger.

## Scope

### In scope

- Add `GET /_admin/ledger` endpoint that returns a unified, time-ordered list of "events" derived from transactions and webhook attempts.
  - Each event includes: `id`, `type` (`transaction` | `webhook`), `time`, `provider`, `reference`, `status`, `summary`.
- Update `internal/ui/index.html` with a ledger view as the default landing tab:
  - Clear zero-data state when no transactions or webhooks exist yet (e.g., "Send your first charge to see it here" with a copy-paste example).
  - Auto-refresh every 2 seconds.
  - Search box filtering reference, provider, and status.
  - Filter tabs: All | Transactions | Webhooks.
  - Click an event to open a detail panel:
    - For transactions: full transaction JSON, provider, amount, status timeline (reuse `GET /_admin/transactions/{ref}` from P06).
    - For webhooks: payload, headers (redacted), signature status, attempts timeline, replay button (reuse `GET /_admin/webhooks/{ref}` from P05).
- Keep existing Transactions and Webhooks tables as secondary tabs or sections for users who prefer them.
- Add keyboard shortcut `?` to show help and `/` to focus search.
- Pause auto-refresh when the browser tab is hidden (`document.visibilityState`) and resume on focus.
- Add tests for `/_admin/ledger` and the detail endpoints.

### Out of scope

- Real-time WebSocket/SSE streaming (polling is fine for v1).
- Multi-select or bulk actions.
- Export to file.

## Acceptance criteria

- [ ] `GET /_admin/ledger` returns a unified event list.
- [ ] Dashboard shows a ledger as the primary view.
- [ ] Ledger shows a useful zero-data state before the first event.
- [ ] Ledger auto-refreshes and supports search/filter.
- [ ] Clicking an event opens a detail panel with payload inspection and replay.
- [ ] `?` shows keyboard shortcuts; `/` focuses search.
- [ ] Tests cover ledger endpoint and detail panels.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Derive events from existing transaction and webhook stores; no new persistence needed.
- Use `setInterval` for polling; WebSocket/SSE can be a future enhancement.
- Keep the ledger under 50 events by default with a "Load more" link to avoid overwhelming the UI.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated `internal/server/admin_api_test.go` and UI tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
