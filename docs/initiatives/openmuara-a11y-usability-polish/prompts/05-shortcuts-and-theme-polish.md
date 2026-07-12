> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 05 — Shortcuts and Theme Polish

> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Target:** `<repo-root>/`
> **Branch:** `feat/a11y-usability-polish`
> **Depends on:** Prompt 01, Prompt 02

---

## Goal

Fix theme-toggle state sync, harden keyboard shortcuts, and add a visible label to the theme toggle.

## Why now

The `d` shortcut can leave the toggle button showing the wrong icon/label. Global shortcuts also fire when modifier keys are held, which can surprise users. The icon-only theme toggle may be unclear to some users.

## Scope

### In scope

- `web/dashboard/src/app.tsx`
- `web/dashboard/src/components/Shell.tsx`
- `web/dashboard/src/theme.ts`
- `web/dashboard/src/styles.css`
- Related dashboard tests.

### Out of scope

- Provider pages.
- Example mini-apps.

## Acceptance criteria

- [ ] `toggleTheme()` returns the newly selected theme and `app.tsx` notifies `Shell` so the button icon/label updates immediately.
  - Options: lift theme state to `App`, use a small pub/sub in `theme.ts`, or re-read the theme after toggle.
- [ ] Global shortcuts (`?`, `Esc`, `/`, `d`, `1`, `2`, `3`) ignore events when Ctrl, Alt, Meta, or Shift are pressed (except `?` which is Shift+/).
- [ ] The theme toggle shows a visually hidden text label (e.g., "Toggle theme") alongside the icon, or a tooltip that works for keyboard users.
- [ ] The Onboarding dark-mode background bug is already fixed in Prompt 02; verify here.
- [ ] Existing tests pass and new tests cover shortcut modifier handling and theme sync.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `cd web/dashboard && npm run test:ci`
  - [ ] `cd web/dashboard && npm run bundle-size`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- The simplest fix is to lift theme state into `App` and pass `theme` + `onToggleTheme` to `Shell`.
- For the tooltip, a `title` attribute is a lightweight first step, but a visible label is better.

## Deliverables

- Code changes on `feat/a11y-usability-polish`.
- Updated dashboard tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
