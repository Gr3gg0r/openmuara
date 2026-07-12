> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 04 — Example Mini-Apps Accessibility

> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Target:** `<repo-root>/`
> **Branch:** `feat/a11y-usability-polish`
> **Depends on:** —

---

## Goal

Make the example mini-apps announce status updates and submit naturally with a keyboard.

## Why now

The ecommerce and prepaid-topup examples update a status `<div>` visually but do not inform screen readers. Inputs are also not wrapped in `<form>` elements, so pressing Enter does not submit.

## Scope

### In scope

- `examples/ecommerce-single-buy/index.html`
- `examples/prepaid-topup/index.html`

### Out of scope

- Dashboard.
- Provider pages.

## Acceptance criteria

- [ ] Each example wraps its inputs in a `<form>` and the submit button is `type="submit"`.
- [ ] Pressing Enter inside any field submits the form.
- [ ] The status `<div>` has `aria-live="polite"` and `aria-atomic="true"` so screen readers announce "Creating checkout session...", errors, and success redirects.
- [ ] The submit button is disabled while the request is in flight to prevent double submission.
- [ ] Existing behavior (redirect to provider page on success) is preserved.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use `event.preventDefault()` in the form submit handler and keep the existing fetch logic.
- Reuse the existing `--focus-ring` and `--primary` CSS variables for the disabled state.

## Deliverables

- Code changes on `feat/a11y-usability-polish`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
