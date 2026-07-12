> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt P01 — Dashboard Dark Mode

> **Initiative:** OpenMuara Dark Mode
> **Target:** `<repo-root>/`
> **Branch:** `feat/dark-mode`
> **Depends on:** —

---

## Goal

Add a cohesive, accessible dark mode to the OpenMuara Preact/Vite dashboard that follows OS preference by default, allows a manual toggle, persists the choice in `localStorage`, syncs across tabs, and avoids a theme flash on load.

## Why now

The dashboard is the primary surface developers stare at while debugging payments. A dark mode reduces eye strain and matches the rest of a typical developer toolchain.

## Current state

- The dashboard is built with Preact + Vite in `web/dashboard/`.
- `web/dashboard/src/styles.css` already defines a few custom properties (`--bg`, `--card`, `--text`, `--muted`, `--border`) but many colors are hard-coded (badges, buttons, alerts, status text, etc.).
- The built SPA is embedded into the Go binary at `internal/ui/dashboard-dist/`; changes must be rebuilt with `npm run build`.
- Existing tests live in `web/dashboard/tests/`.

## Scope

### In scope

- Replace hard-coded colors in `web/dashboard/src/styles.css` with semantic custom properties for both light and dark modes.
- Add a small, inline blocking theme-setup script in `web/dashboard/index.html` that runs before Preact mounts.
- Add a theme toggle button to `Shell.tsx` (next to the Help button) with sun/moon icon and `aria-label`.
- Persist manual choice in `localStorage` under key `muara-theme`.
- Sync theme across tabs via the `storage` event.
- Listen to `prefers-color-scheme` changes so the theme updates if the user changes OS preference while no manual override is set.
- Add `<meta name="color-scheme" content="light dark">` and a dynamic `theme-color` meta tag.
- Add a keyboard shortcut (`d`) to toggle theme; document it in the help modal.
- Honor `prefers-reduced-motion` by disabling the color transition when the user prefers reduced motion.
- Update `Shell.test.tsx` to assert the toggle exists and switches the `data-theme` attribute.
- Rebuild `internal/ui/dashboard-dist/` so the embedded dashboard reflects the changes.

### Out of scope

- Redesigning layout, navigation, or components.
- Theming provider pay pages or example mini-apps (covered in P02 and P03).
- Adding a theme preview or multi-color accent picker.

## Acceptance criteria

- [ ] Dashboard renders correctly in dark mode with no visual regressions in light mode.
- [ ] OS `prefers-color-scheme: dark` is honored on first visit when no manual choice exists.
- [ ] Manual toggle overrides OS preference and persists across reloads.
- [ ] Theme choice syncs across multiple OpenMuara tabs.
- [ ] No flash of wrong theme when reloading with dark mode selected.
- [ ] Focus indicators, buttons, links, and status badges meet WCAG AA contrast in both modes.
- [ ] `prefers-reduced-motion` disables the color transition.
- [ ] Print preview is readable (dark mode prints as light or with sufficient contrast).
- [ ] `forced-colors` / Windows High Contrast Mode does not hide critical borders or focus rings.
- [ ] Bundle size increase is ≤5 KiB gzipped.
- [ ] Embedded dashboard (`internal/ui/dashboard-dist/index.html`) is rebuilt and tested, not just the Vite dev server.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `cd web/dashboard && npm run test:ci`
  - [ ] `cd web/dashboard && npm run bundle-size`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use `data-theme="dark"` / `data-theme="light"` on `<html>`; define tokens with `:root, [data-theme="light"]` and `[data-theme="dark"]`.
- Keep the blocking script tiny (~30 lines) and inline in `index.html`.
- For the default: `const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;`
- For cross-tab sync: `window.addEventListener('storage', (e) => { if (e.key === 'muara-theme') setTheme(e.newValue); });`
- Use Unicode `☀`/`☾` or a tiny inline SVG for the toggle icon; avoid an icon dependency.
- Test the embedded build, not just the Vite dev server.

## Deliverables

- Code changes on `feat/dark-mode`.
- Updated dashboard tests.
- Rebuilt `internal/ui/dashboard-dist/`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Release-notes snippet describing the new dark mode.
- Git commit with a clear message.
