> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P03 — Ledger Default View and Detail Page

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** P01
> **Target files:** `web/dashboard/src/views/Ledger.tsx`, `web/dashboard/src/views/LedgerDetail.tsx`, `web/dashboard/src/components/FilterToolbar.tsx`, `web/dashboard/src/app.tsx`, `web/dashboard/src/styles.css`
> **Status:** ⬜

## Goal

Make the Ledger the default dashboard view, give it a reusable filter toolbar, and open ledger rows into a dedicated detail page.

## Tasks

- [ ] Create `FilterToolbar.tsx` with search input, provider select, status select, sort select, and refresh button.
- [ ] Refactor `Ledger.tsx` to use `FilterToolbar` and keep its existing filtering, sorting, and polling behavior.
- [ ] Create `LedgerDetail.tsx` that reads `ref` from URL state and fetches `/_admin/transactions/{ref}` or `/_admin/webhooks/{ref}` based on the ledger row type.
- [ ] Make ledger rows clickable; clicking a row navigates to `/_admin?view=ledger-detail&ref=...`.
- [ ] Add a back button on `LedgerDetail` that returns to the Ledger view and restores focus to the originating row.
- [ ] Ensure `/_admin` with no query params defaults to Ledger.
- [ ] Preserve existing keyboard shortcut `/` to focus the ledger search.

## Acceptance Criteria

- [ ] `/_admin` opens the Ledger table by default.
- [ ] The Ledger table has a filter toolbar with search, provider, status, and sort controls.
- [ ] Clicking a ledger row navigates to a Ledger Detail page.
- [ ] The detail page shows the full transaction or webhook payload and a replay action.
- [ ] The back button returns to the Ledger view with filters intact.
- [ ] Axe-core reports zero serious violations on Ledger and LedgerDetail.

## Quality Gates

Run before committing:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
```

## Notes

- Reuse the existing detail fetch logic from the current inline panel; do not change backend endpoints in this prompt.
- The filter toolbar will also be used by the Webhooks view in P04.
