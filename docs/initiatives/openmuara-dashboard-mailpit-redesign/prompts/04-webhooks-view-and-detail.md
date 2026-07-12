> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P04 — Webhooks View and Detail Page

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** P01
> **Target files:** `web/dashboard/src/views/Webhooks.tsx`, `web/dashboard/src/views/WebhookDetail.tsx`, `web/dashboard/src/components/FilterToolbar.tsx`, `web/dashboard/src/app.tsx`, `web/dashboard/src/styles.css`
> **Status:** ✅

## Goal

Promote webhooks to a top-level navigation item as a delivery log, add the reusable filter toolbar, and open webhook rows into a dedicated detail page. Per-provider webhook configuration is intentionally moved to Settings → Provider Detail (P05).

## Tasks

- [x] Refactor `Webhooks.tsx` to be a top-level delivery-log view (no longer nested under Settings).
- [x] Add `FilterToolbar` to the Webhooks table with URL, provider, status filters and sort.
- [x] Remove `WebhookConfig` from the Webhooks view; per-provider webhook configuration moves to Settings → Provider Detail (P05).
- [x] Create `WebhookDetail.tsx` that reads `ref` from URL state and fetches `/_admin/webhooks/{ref}`.
- [x] Make webhook rows clickable; clicking a row navigates to `/_admin?view=webhook-detail&ref=...`.
- [x] Add a back button on `WebhookDetail` that returns to the Webhooks view and restores focus to the originating row.
- [x] Preserve replay functionality on both the list and detail pages.

## Acceptance Criteria

- [x] The Webhooks view is reachable from the left sidebar.
- [x] The Webhooks view is a delivery log only (no configuration UI).
- [x] The Webhooks table has a filter toolbar with URL, provider, status, and sort controls.
- [x] Clicking a webhook row navigates to a Webhook Detail page.
- [x] The detail page shows full webhook metadata, headers, payload, attempts, and a replay action.
- [x] The back button returns to the Webhooks view with filters intact.
- [x] Axe-core reports zero serious violations on Webhooks and WebhookDetail.

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

- The `WebhookConfig` component was removed from the Webhooks view; per-provider webhook configuration moves to Settings → Provider Detail in P05.
- Reuse existing backend endpoints; no new backend work in this prompt.
