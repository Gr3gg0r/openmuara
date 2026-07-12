> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 04 — Dashboard Usability

> **Initiative:** OpenMuara v1 Solid Gold
> **Target:** `<repo-root>/`
> **Branch:** `feat/v1-solid-gold`
> **Depends on:** Prompt 01

---

## Goal

Reduce the time from "something failed" to "I see what failed" in the web dashboard.

## Why now

The ledger and webhook debugger are useful, but testers still have to scan the
 table to notice a failed webhook. Small UX cues will make the dashboard feel
 more like Mailpit.

## Scope

### In scope

- Add a visible failed-webhook alert/notification bar in `/_admin`.
- Add a **Copy curl** button on provider cards for quick first-charge examples.
- Improve ledger table responsiveness on mobile viewports.
- Update `internal/ui/index.html` only.

### Out of scope

- Dashboard authentication.
- Major redesign or new routes.

## Acceptance criteria

- [ ] A failed webhook attempt surfaces a clear alert at the top of the dashboard.
- [ ] Each provider card has a copy-to-clipboard curl example.
- [ ] The ledger table is usable on a 375px-wide viewport.
- [ ] Existing keyboard shortcuts (`?`, `/`, `1/2/3`, `Esc`) keep working.
- [ ] Changes are verified by a Playwright/browser snapshot or documented manual QA checklist.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use the existing `/metrics` or `/webhooks` endpoint to detect failed attempts.
- Keep the alert dismissible per session.

## Deliverables

- Code changes on `feat/v1-solid-gold`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
