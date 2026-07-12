> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P06 — Visual Sign-off and Philosophy Check

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** P01–P05
> **Target files:** `HANDOFF.md`, `TRACKING.md`
> **Status:** ✅ Completed

## Goal

Use Playwright MCP to verify that the dashboard still matches the Mailpit-style design and the project philosophy after all bug fixes. This is the final gate before declaring the branch ready for PR.

## Tasks

- [x] Start the server with the default single-port configuration and capture:
  - [x] Ledger view: full-width table, filter toolbar, left nav active item.
  - [x] Webhooks view: delivery-log table, filter toolbar, no config UI.
  - [x] Settings view: provider card grid with status badges.
  - [x] ProviderDetail for a multi-version provider (e.g. Fawry): enable toggle, v1/v2 tabs, base URL, sample endpoint, webhook target input, env var list.
  - [x] ProviderDetail for a single-version provider (e.g. Stripe): no version tabs, base URL, env var list.
- [x] Start the server with `server.admin_port` set and verify the admin UI/API is reachable only on the admin port while provider endpoints remain on the provider port.
- [x] Verify keyboard shortcuts (`1`, `2`, `3`, `/`, `r`, `?`, `d`) still work.
- [x] Run axe-core or the existing a11y test suite and confirm zero serious violations.
- [x] Check that the bundle size remains within budget.
- [x] Compare the captured screenshots against the P01 baseline and note any unexpected visual changes.
- [x] Update `HANDOFF.md` with screenshot references and a summary of the visual/philosophy check.
- [x] Complete `REVIEW_CHECKLIST.md` and `appendices/c-post-bug-hunt-checklist.md`.
- [x] Update `TRACKING.md` to mark P06 complete and the branch ready for PR.

## Acceptance Criteria

- [x] Playwright MCP captures the Ledger, Webhooks, Settings, and ProviderDetail views.
- [x] Left navigation shows exactly Ledger, Webhooks, Settings.
- [x] Ledger is the default view at `/_admin`.
- [x] Every table view has a filter toolbar.
- [x] Ledger and webhook rows navigate to dedicated detail pages.
- [x] Webhooks view is delivery-log only; provider webhook config is in Settings.
- [x] Provider detail pages have enable toggle, base URL, version tabs when applicable, webhook URL input, and env var names.
- [x] Dual-port runtime works as documented.
- [x] Keyboard shortcuts and a11y checks pass.
- [x] Bundle size stays within budget.
- [x] `HANDOFF.md` includes final visual sign-off summary.
- [x] `REVIEW_CHECKLIST.md` and `appendices/c-post-bug-hunt-checklist.md` are complete.
- [x] `TRACKING.md` marks P06 complete.

## Quality Gates

No new product code should be needed in this prompt. Run the full suite to confirm the branch is green:

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
cd web/dashboard && node scripts/a11y-contrast-check.js
./scripts/smoke-test.sh || true
```

## Notes

- If visual sign-off reveals a bug, either fix it in a focused commit and re-run P06, or file it as a deferred item in `RISKS.md` with user sign-off.
- The priority stack for evaluation is: **UI > UX > Performance > Usability > Philosophy > Efficiency > Memory size**. A minor visual improvement that increases bundle size slightly is acceptable; a heavy dependency added only for convenience is not.
- Attach or reference screenshots in `HANDOFF.md` so the next agent or reviewer can see the final state.
