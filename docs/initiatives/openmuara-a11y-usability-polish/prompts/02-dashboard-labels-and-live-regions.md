> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 02 — Dashboard Labels and Live Regions

> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Target:** `<repo-root>/`
> **Branch:** `feat/a11y-usability-polish`
> **Depends on:** Prompt 01

---

## Goal

Give every dashboard control a meaningful accessible name and announce important state changes to screen readers.

## Why now

Several controls rely on placeholder text, repeated button text, or visual state alone. Screen-reader users cannot identify the search fields, distinguish "Copy curl" buttons per provider, or hear confirmation when copy succeeds.

## Scope

### In scope

- `web/dashboard/src/views/Ledger.tsx`
- `web/dashboard/src/views/Transactions.tsx`
- `web/dashboard/src/components/Providers.tsx`
- `web/dashboard/src/components/Onboarding.tsx`
- `web/dashboard/src/components/FailedWebhookAlert.tsx`
- Related dashboard tests.

### Out of scope

- Table row click behavior (Prompt 01).
- Provider pages and examples.

## Acceptance criteria

- [ ] Search inputs in Ledger and Transactions have persistent labels (`<label>` or `aria-label`).
- [ ] Each "Copy curl" button in the Providers list has a unique `aria-label` including the provider name.
- [ ] Copy-to-clipboard success is announced with an `aria-live` region (e.g., "Copied Fawry curl to clipboard").
- [ ] The Onboarding Show/Hide button has `aria-expanded` reflecting panel state.
- [ ] The Onboarding panel's hard-coded `background:'#f8fafc'` is replaced with a theme token (fixes dark-mode visual bug).
- [ ] The failed-webhook alert uses a `<button>` styled as a link (instead of `<a href="#">`) and the warning icon has `aria-hidden` with a text label.
- [ ] Existing tests pass and new tests cover label/live-region behavior.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `cd web/dashboard && npm run test:ci`
  - [ ] `cd web/dashboard && npm run bundle-size`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- A single visually hidden `<span aria-live="polite" aria-atomic="true">` at the app root is enough for global announcements.
- For the warning icon, prefer text like "Warning:" plus `aria-hidden` on the emoji, or use an SVG with `<title>`.

## Deliverables

- Code changes on `feat/a11y-usability-polish`.
- Updated dashboard tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
