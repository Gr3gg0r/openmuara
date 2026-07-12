> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 04 — Escape Page and Webhook Shape

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** Prompts 01, 02, 03

---

## Goal

Update the Fawry escape page so testers can choose all four statuses and the
billing journey, and see a preview of the webhook that will be dispatched.

## Why now

The escape page is the primary manual testing surface. It must expose the new
states, delay, and billing type added in previous prompts.

## Scope

### In scope

- Add **CANCELED** and **EXPIRED** action buttons to `web/fawry-escape.html`.
- Add a `billing_type` selector (recurring / one_time) that defaults to the
  value stored on the transaction.
- Show the configured `response_delay_ms` and the expected webhook status.
- Update `internal/fawry/escape.go` form handling to forward `billing_type`.
- Ensure CSRF protection still works for new form fields.
- Add UI-level tests if feasible.

### Out of scope

- Redesigning the rest of the dashboard.
- Live webhook preview from a not-yet-dispatched event.

## Acceptance criteria

- [ ] The escape page renders Paid, Unpaid, Canceled, and Expired actions.
- [ ] The billing type selector reflects the transaction's stored value and can
      be overridden.
- [ ] Submitting the form triggers the correct status and journey webhook.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Reuse the existing `ui.EscapePageData` struct; add optional fields for
  `billingType` and `delayMs`.
- Keep the page simple; use a `<select>` for billing type and small buttons for
  statuses.

## Deliverables

- Code changes on `feat/mkp-fawry`.
- Updated escape handler and page tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
