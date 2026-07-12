> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt P02 — Provider Pay Pages Dark Mode

> **Initiative:** OpenMuara Dark Mode
> **Target:** `<repo-root>/`
> **Branch:** `feat/dark-mode`
> **Depends on:** —

---

## Goal

Apply a consistent dark mode to all OpenMuara-hosted provider payment simulation pages so they follow the OS color scheme without flashing white on load.

## Why now

During end-to-end testing, users are redirected from their app to OpenMuara pay pages. A bright white flash in an otherwise dark workflow is jarring and looks unpolished.

## Current state

Provider pay pages are static HTML templates in `internal/ui/`:
- `stripe-checkout.html`
- `stripe-payment-intent.html`
- `stripe-webhooks.html`
- `fawry-escape.html`
- `billplz-pay.html`
- `toyyibpay-pay.html`
- `ipay88-pay.html`

Most use hard-coded light colors and do not define a dark palette.

## Scope

### In scope

- Convert each pay page to CSS custom properties (reuse the dashboard token naming where practical).
- Add `@media (prefers-color-scheme: dark)` rules for each page.
- Add `<meta name="color-scheme" content="light dark">` to each page.
- Ensure form controls, buttons, and text meet WCAG AA contrast in both modes.
- Keep pages dependency-free: no external CSS or JS files.
- Update any Go-rendered inline styles in `internal/stripe/*.go`, `internal/fawry/*.go`, `internal/billplz/*.go`, `internal/ipay88/*.go`, `internal/toyyibpay/*.go` to use the same custom properties if they inline styles.
- Verify that no provider page references hard-coded light colors after tokenization.

### Out of scope

- Adding a manual toggle on pay pages (they should just follow OS preference).
- Redesigning the pay-page layout or copy.

## Acceptance criteria

- [ ] All seven pay pages render legibly in dark mode.
- [ ] All seven pay pages render identically to before in light mode (no regressions).
- [ ] No white flash when loading a pay page with OS dark mode enabled.
- [ ] Form controls adopt the dark color scheme via `color-scheme`.
- [ ] WCAG AA contrast is met for all text and interactive elements in both modes.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Define tokens inside a `<style>` block in each page; no external CSS to avoid extra HTTP requests.
- Use `color-scheme: light dark;` on `:root` so native form controls theme automatically.
- Test each provider end-to-end by creating a charge and inspecting the pay/escape page.

## Deliverables

- Code changes on `feat/dark-mode`.
- Updated tests if any pay-page tests assert on HTML structure or inline styles.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
