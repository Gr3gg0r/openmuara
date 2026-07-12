> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt P03 — Example Mini-Apps Dark Mode

> **Initiative:** OpenMuara Dark Mode
> **Target:** `<repo-root>/`
> **Branch:** `feat/dark-mode`
> **Depends on:** —

---

## Goal

Add dark mode support to the two example mini-apps (`ecommerce-single-buy` and `prepaid-topup`) so they respect OS preference with zero extra dependencies.

## Why now

The examples are the first thing a new user sees. A polished dark mode makes OpenMuara feel modern and shows that the theming system is easy to apply.

## Current state

- `examples/ecommerce-single-buy/index.html` — product card with hard-coded light colors.
- `examples/prepaid-topup/index.html` — top-up form with hard-coded light colors.
- Both are single-file HTML pages with inline styles and a small inline script.

## Scope

### In scope

- Convert inline styles in both example `index.html` files to CSS custom properties.
- Add `@media (prefers-color-scheme: dark)` rules.
- Add `<meta name="color-scheme" content="light dark">`.
- Add a small blocking head script that sets `data-theme` before first paint to avoid flashes.
- Keep the implementation dependency-free and minimal.

### Out of scope

- A manual theme toggle in the examples (follow OS preference only; see decision D006).
- The example Go servers (`main.go`) do not need changes unless required for a new static asset.
- Redesigning the example pages.

## Acceptance criteria

- [ ] Both example pages render correctly in OS dark mode.
- [ ] Both example pages render identically to before in OS light mode.
- [ ] No white flash when loading in OS dark mode.
- [ ] Form controls adopt the dark color scheme via `color-scheme`.
- [ ] No extra dependencies are added.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use the same token naming convention as the dashboard for consistency.
- Keep the head script under ~20 lines.

## Deliverables

- Code changes on `feat/dark-mode`.
- Updated `examples/README.md` if screenshots or descriptions are added.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
