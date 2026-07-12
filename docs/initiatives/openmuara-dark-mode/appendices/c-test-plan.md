> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# Appendix C — Manual Test Plan

> Run these checks after implementing each prompt.

---

## Dashboard (P01)

### Light mode baseline
- [ ] Open `http://127.0.0.1:9000/_admin` with OS light mode and no `muara-theme` key.
- [ ] Confirm the page looks identical to before the change.
- [ ] Confirm Ledger, Transactions, and Webhooks tabs render correctly.
- [ ] Confirm the onboarding checklist, provider cards, and table rows are readable.

### Dark mode
- [ ] Set OS to dark mode, clear `localStorage`, reload.
- [ ] Confirm the dashboard renders in dark mode automatically.
- [ ] Confirm no white flash during reload.
- [ ] Confirm all text, buttons, badges, links, and table headers are readable.
- [ ] Confirm focus rings are visible when tabbing through the UI.

### Manual toggle
- [ ] Click the theme toggle; confirm the page switches mode.
- [ ] Reload; confirm the chosen mode persists.
- [ ] Open a second tab, toggle in one tab, confirm the other tab updates.

### OS preference change
- [ ] Clear `localStorage`, set OS to dark, reload.
- [ ] Change OS to light while the page is open; confirm the page updates.
- [ ] Set a manual override, change OS preference; confirm manual override wins.

### Reduced motion
- [ ] Enable OS reduced motion, toggle theme; confirm no color transition plays.

### Accessibility
- [ ] Run axe DevTools; confirm no contrast errors.
- [ ] Tab through all interactive elements; confirm focus indicators are visible.
- [ ] Confirm the toggle button has an `aria-label`.
- [ ] Enable Windows High Contrast Mode or use DevTools `forced-colors: active`; confirm borders and focus rings remain visible.
- [ ] Print-preview the dashboard in dark mode; confirm text remains readable.

### Bundle size
- [ ] Run `cd web/dashboard && npm run bundle-size`.
- [ ] Confirm JS bundle increase is ≤5 KiB gzipped.

### Embedded build
- [ ] Run `cd web/dashboard && npm run build`.
- [ ] Start the Go binary and open `/_admin` (not the Vite dev server).
- [ ] Confirm the embedded dashboard renders correctly in both light and dark modes.

## Provider pay pages (P02)

For each provider, create a charge and inspect the resulting pay/escape page:

- [ ] Stripe Checkout pay page (`/v1/checkout/sessions/{id}/pay`)
- [ ] Stripe Payment Intent admin page
- [ ] Fawry escape page (`/_admin/fawry-escape`)
- [ ] Billplz pay page
- [ ] ToyyibPay pay page
- [ ] iPay88 pay page

Checks for each page:
- [ ] Renders legibly in OS light mode.
- [ ] Renders legibly in OS dark mode.
- [ ] No white flash on load in dark mode.
- [ ] Form controls (inputs, selects, radios) use the dark color scheme.
- [ ] Buttons remain clickable and readable.

## Example mini-apps (P03)

- [ ] Open `http://127.0.0.1:8080/` in OS light mode; confirm the product card looks as before.
- [ ] Open `http://127.0.0.1:8080/` in OS dark mode; confirm it switches to dark.
- [ ] Toggle theme manually; confirm persistence across reload.
- [ ] Repeat for `http://127.0.0.1:8081/`.

## Docs (P04)

- [ ] `README.md` mentions dashboard dark mode.
- [ ] `CHANGELOG.md` has an unreleased entry.
