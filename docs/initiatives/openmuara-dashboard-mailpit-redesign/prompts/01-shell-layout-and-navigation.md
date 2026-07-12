> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P01 — Shell Layout and Navigation

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** —
> **Target files:** `web/dashboard/src/components/AppShell.tsx`, `web/dashboard/src/components/SidebarNav.tsx`, `web/dashboard/src/app.tsx`, `web/dashboard/src/styles.css`, `web/dashboard/src/components/Shell.tsx`
> **Status:** ⬜

## Goal

Replace the current top-tabbed shell with a Mailpit-like fixed left sidebar that exposes three primary views: **Ledger**, **Webhooks**, and **Settings**.

## Tasks

- [ ] Create `AppShell.tsx` with a compact top bar and a fixed left sidebar.
- [ ] Create `SidebarNav.tsx` with three links: Ledger, Webhooks, Settings.
- [ ] Update `app.tsx` to use `AppShell` and manage `view` URL state (`ledger` | `webhooks` | `settings`).
- [ ] Preserve existing keyboard shortcuts and add `1`/`2`/`3` shortcuts for the three nav items.
- [ ] Keep `tab=` query parameter as a one-release fallback redirect.
- [ ] Update `styles.css` for the new layout without breaking existing component styles.
- [ ] Remove or repurpose the old `Shell.tsx` top tabs.

## Acceptance Criteria

- [ ] `/_admin` renders a left sidebar with Ledger, Webhooks, Settings.
- [ ] Clicking each item switches the main outlet and updates the URL.
- [ ] The active item is visually highlighted and has `aria-current="page"`.
- [ ] Keyboard shortcuts `1`, `2`, `3` switch views.
- [ ] Existing `tab=ledger` redirects to the new ledger view.
- [ ] Axe-core reports zero serious violations on the shell.

## Quality Gates

Run before committing:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
```

## Notes

- Do not implement the Settings content in this prompt; only the shell and navigation.
- Keep the existing Overview onboarding checklist accessible from the Ledger empty state or header in P02.
