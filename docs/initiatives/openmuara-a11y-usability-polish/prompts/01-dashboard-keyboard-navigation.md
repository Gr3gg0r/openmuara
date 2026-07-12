> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 01 — Dashboard Keyboard Navigation

> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Target:** `<repo-root>/`
> **Branch:** `feat/a11y-usability-polish`
> **Depends on:** —

---

## Goal

Make every interactive part of the dashboard reachable and operable with a keyboard, and manage focus correctly for tabs, modals, and detail panels.

## Why now

The ledger, transactions, and webhooks tables rely on `onClick` attached to `<tr>` elements, which are not focusable. The help modal and detail panels also lack focus management. These are blockers for keyboard and screen-reader users.

## Scope

### In scope

- `web/dashboard/src/views/Ledger.tsx`
- `web/dashboard/src/views/Transactions.tsx`
- `web/dashboard/src/views/Webhooks.tsx`
- `web/dashboard/src/components/Shell.tsx`
- `web/dashboard/src/styles.css`
- Related dashboard tests.

### Out of scope

- Provider simulation pages.
- Example mini-apps.
- Theme/shortcut logic.

## Acceptance criteria

- [ ] Table rows in Ledger, Transactions, and Webhooks are keyboard-actionable. Options:
  - Add a real `<button>` or `<a>` inside each row for the primary action, **or**
  - Make the row itself focusable with `tabIndex`, `role="button"`, and `onKeyDown` for Enter/Space.
- [ ] Focus moves into a detail panel when it opens, and the Close button returns focus to the triggering row/control.
- [ ] The help modal traps focus while open and returns focus to the Help button on close.
- [ ] The tab bar supports left/right arrow keys, Home, and End (roving `tabIndex` or equivalent).
- [ ] Escape closes the help modal and detail panels.
- [ ] Existing tests pass and new tests cover the keyboard behavior.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `cd web/dashboard && npm run test:ci`
  - [ ] `cd web/dashboard && npm run bundle-size`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Keep changes minimal: prefer a small `<button class="row-action">` inside the first cell over re-engineering the whole table.
- For the help modal, consider using a small `useFocusTrap` hook.
- Use `data-testid` for new test targets.

## Deliverables

- Code changes on `feat/a11y-usability-polish`.
- Updated dashboard tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
