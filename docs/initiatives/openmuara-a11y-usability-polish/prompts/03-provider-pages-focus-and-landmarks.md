> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 03 — Provider Pages Focus and Landmarks

> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Target:** `<repo-root>/`
> **Branch:** `feat/a11y-usability-polish`
> **Depends on:** —

---

## Goal

Add visible focus indicators and semantic landmarks to all server-rendered provider simulation pages.

## Why now

The provider pages (Stripe checkout, PaymentIntent, webhooks config, Billplz, iPay88, ToyyibPay, Fawry escape) currently show no focus outline on buttons and have no `<main>` landmark. Keyboard users cannot see where focus is.

## Scope

### In scope

- `internal/ui/stripe-checkout.html`
- `internal/ui/stripe-payment-intent.html`
- `internal/ui/stripe-webhooks.html`
- `internal/ui/billplz-pay.html`
- `internal/ui/ipay88-pay.html`
- `internal/ui/toyyibpay-pay.html`
- `internal/ui/fawry-escape.html`

### Out of scope

- Dashboard Preact app.
- Example mini-apps.

## Acceptance criteria

- [ ] Every page has a `<main>` landmark wrapping the primary content.
- [ ] Buttons on every page have a visible `:focus` style (match the existing input focus rings).
- [ ] Focus order follows the visual order.
- [ ] Form controls already have labels; verify no regressions after wrapping content in `<main>`.
- [ ] Pages still render correctly after dark-mode changes.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Keep the existing lightweight, self-contained HTML structure. Adding `<main>` and `:focus` rules is enough.
- Reuse the `--focus-ring` CSS variable already defined in each page.

## Deliverables

- Code changes on `feat/a11y-usability-polish`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
