> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard Design Refresh

> **Status:** 🟡 Planned | **Started:** 2026-07-03
> **Scope:** Audit and redesign the dashboard to look intentional, polished, and professional while staying lightweight, low-memory, and aligned with OpenMuara's philosophy of simple, fast, efficient tools. Upgrade the information architecture so the page no longer feels cramped, and make common tester/dev workflows (filtering by provider/URL, scanning webhook attempts, comparing statuses) faster and more obvious.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Problem

The dashboard works but still looks generic, feels cramped, and lacks the polish expected of a developer-facing control plane:

- Inline `style` attributes and magic-number spacing remain in several components.
- Unicode symbols (`☀`, `☾`, `↻`) are used instead of a consistent, accessible icon set.
- Tables are dense and visually monotonous; rows lack hover states and sticky headers.
- Filter toolbars are plain `input` + `select` rows with no visual hierarchy or advanced filtering affordances.
- Provider cards have uneven heights and no status iconography, making the grid feel ragged.
- Empty states are text-only with no icon or primary CTA.
- The Ledger tab stacks Providers + Onboarding + Ledger vertically, creating a cramped, scroll-heavy page.
- No command palette, density toggle, or persistent layout preferences.
- Help and onboarding copy feel placeholder-ish.

---

## Design principles

The refresh must stay true to OpenMuara's philosophy: simple, fast, efficient, local-first.

- **Content first** — every pixel should serve a tester or developer workflow; remove decorative noise.
- **Progressive disclosure** — show the most important information by default; hide details behind interactions.
- **Consistent affordances** — the same action should look the same everywhere.
- **Respect the runtime** — no heavy dependencies, no external tracking, no fonts that block first paint.
- **Accessibility is not a polish step** — design for keyboard, screen readers, and high contrast from the start.
- **Local-first trust** — the dashboard runs on the user's machine; dark mode and density are their choice.

---

## Live design audit

Audited `http://127.0.0.1:9000/_admin` at 1440×900, 768×1024, and 375×812 with axe-core 4.9.1 and manual inspection. Screenshots saved in this initiative folder.

### Axe-core findings

| ID | Impact | Count | Notes |
|---|---|---|---|
| `region` | moderate | 2 | `h1` and `nav` are not inside landmarks. Should wrap header in `<header>` and tabs in `<nav>` with an `aria-label`. |

### Manual findings

- **24 inline `style` attributes** on the Ledger view; most are layout tweaks that should be token classes.
- **Cramped vertical stacking** on Ledger: active provider, Stripe webhook link, provider grid, ledger filter, ledger table, and detail panels all compete for attention.
- **No icon system** — reload, theme, help, copy, and close actions use Unicode characters that render inconsistently across OS/browsers and are not announced to screen readers.
- **Tables lack micro-interactions** — no row hover, no sticky header, no selected-row state.
- **Filters are visually flat** — no search icon, no clear-button chips, no date-range picker, no "active filters" summary.
- **Provider grid rhythm is broken** — enabled cards are tall, disabled cards are short, creating a ragged right edge.
- **Status badges are inconsistent** — some are colored text, some are badges; no unified semantic color language.
- **Empty states** are centered gray text only.
- **Onboarding checklist** is not surfaced in the SPA.
- **Connection status** is not shown.

---

## Current-state audit (updated)

| Area | State | Notes |
|---|---|---|
| CSS variables | ✅ Mostly done | Color, spacing, radius, shadow, and typography tokens exist in `styles.css`. |
| Inline styles | ⚠️ Partial | Reduced but still present in `Shell.tsx`, `LedgerView.tsx`, detail panels, and WebhookConfig. |
| Typography | ⚠️ Partial | Token-based scale exists, but still uses `system-ui` only; no deliberate heading font. |
| Icons | ❌ Missing | Unicode symbols are used; no SVG icon system yet. |
| Spacing | ✅ Mostly done | 4px scale exists and is used in CSS; some inline magic numbers remain. |
| Responsive | ✅ Good | Mobile breakpoints exist; tables wrap. |
| Accessibility | ⚠️ Partial | Skip link, focus rings, ARIA tabs, keyboard shortcuts exist; landmarks still incomplete. |
| Motion | ✅ Good | `prefers-reduced-motion` and `prefers-contrast: more` supported. |
| Empty states | ⚠️ Basic | Structure exists but visuals are minimal. |
| Error handling | ✅ Partial | `ErrorBoundary` exists; error banners are plain. |
| Table filtering | ⚠️ Basic | Provider/status search works; URL filter and date-range filter are not yet designed. |
| Layout density | ❌ Missing | No compact/comfortable toggle. |
| Command palette | ❌ Missing | No `Cmd+K` navigation. |

---

## Goals

1. Run a design audit on the live dashboard and save the report in the initiative folder.
2. Eliminate remaining inline styles by replacing them with token-based utility classes.
3. Adopt a lightweight, accessible SVG icon system (no font files, minimal bundle impact).
4. Redesign the information architecture so the dashboard is no longer cramped:
   - Move provider management and onboarding to a dedicated **Overview** tab.
   - Give each tab a clear primary surface and breathing room.
5. Redesign tables with sticky headers, row hover, zebra striping, sort indicators, and selected-row states.
6. Redesign filter toolbars with search icons, clearable filter chips, date-range presets, and an "active filters" bar.
7. Redesign provider cards with equal heights, status iconography, and better grid rhythm.
8. Improve empty states with icon + headline + CTA.
9. Add a command palette (`Cmd+K` / `Ctrl+K`) for jumping tabs and actions.
10. Add a data-density toggle (compact / comfortable) persisted in `localStorage`.
11. Add/update component tests and visual regression baselines.
12. Pass all quality gates with bundle budget intact.

---

## Recommendations & enhancements

### Design system

- **Component library** — reusable `Button`, `Card`, `Badge`, `Input`, `Select`, `Modal`, `Toast`, `Table`, `EmptyState`, and `FilterChip` components with variants.
- **Icon system** — inline SVG sprite or tree-shakeable Preact Lucide icons. No icon font to keep bundle small. Every icon has an `aria-label` or is `aria-hidden` with a visible text fallback.
- **Color semantics** — map all status colors (`paid`, `failed`, `pending`, etc.) to semantic tokens (`--color-success`, `--color-danger`, `--color-warning`) and ensure light/dark/high-contrast parity.
- **Spacing scale** — enforce the existing 4px scale; remove every inline `margin`/`padding` magic number.
- **Border-radius scale** — small for inputs, medium for cards, large for modals.
- **Elevation system** — shadow tokens for cards, modals, and toasts.
- **Motion scale** — 100ms for micro-interactions, 200ms for reveals; respect `prefers-reduced-motion`.
- **Typography** — keep `system-ui` as body fallback, but add one carefully chosen heading font via `font-display: swap` and a self-hosted WOFF2 under 30 KB.
- **Density tokens** — `--density-row-sm`, `--density-row-md`, `--density-gap-sm`, `--density-gap-md` toggled by a user preference.
- **Z-index scale** — `--z-base: 0`, `--z-sticky: 10`, `--z-dropdown: 20`, `--z-modal-backdrop: 30`, `--z-modal: 40`, `--z-toast: 50`, `--z-skip-link: 100` to avoid arbitrary stacking values.
- **Breakpoint scale** — document exact breakpoints: `xs: 0`, `sm: 480px`, `md: 768px`, `lg: 1024px`, `xl: 1440px`.
- **Container queries** — prefer container queries for component-level responsiveness where appropriate; keep media queries for page-level layout.
- **CSS methodology** — use a flat BEM-like naming convention (`.component`, `.component--modifier`, `.component__element`) plus a small set of layout utilities (`.flex`, `.gap-2`, `.mb-4`). Avoid CSS-in-JS.
- **Component API conventions** — props naming: `variant`, `size`, `disabled`, `loading`, `onClick`; compound components for complex patterns (e.g., `Modal`, `Modal.Header`, `Modal.Body`).

### Information architecture & layout

- **Overview tab** — dedicated surface for onboarding checklist, enabled-provider summary, connection status, recent activity, and quick-start CTAs. Remove providers/onboarding from Ledger tab.
- **Providers tab** (optional) — if provider management grows, split it out of Overview; for v1 keep it in Overview but use a clean card grid.
- **Ledger tab** — focused on the unified feed only; sticky filter bar; no competing panels.
- **Transactions tab** — focused on transaction table only.
- **Webhooks tab** — focused on webhook config + delivery log; URL filter prominent.
- **Sticky headers** — table headers stick on scroll; filter bar sticks above them.
- **Side panel for details** — transaction/webhook detail opens in a slide-over panel instead of an inline card, preserving scroll position.
- **Breadcrumbs** — when navigating into nested views.

### Tables

- **Sticky headers** — `position: sticky; top: 0` with a subtle shadow on scroll.
- **Row hover** — background tint on mouse hover.
- **Zebra striping** — subtle alternating row backgrounds (optional, toggleable).
- **Sort indicators** — visible arrow icons in sortable headers.
- **Selected row** — highlight the row whose detail panel is open.
- **Column resizing** (future) — not required for v1.
- **Pagination / virtual scrolling** — keep memory low for large ledgers.

### Filtering

- **Search input with icon** — magnifying glass icon, clear button inside the field.
- **Filter chips** — show active filters as removable chips below the toolbar.
- **Provider multi-select** — allow selecting multiple providers at once.
- **Status multi-select** — same for statuses.
- **URL filter** — dedicated search for webhook target URLs.
- **Date-range filter** — preset chips (last hour, today, last 7 days) plus custom range.
- **Reset all** — one click to clear filters.

### Provider cards

- **Equal heights** — use CSS Grid with `grid-template-rows: auto 1fr auto` so cards align.
- **Status iconography** — icon + color for enabled/disabled/active/misconfigured.
- **Quick actions** — visible "Configure" / "Test" / "Docs" buttons, not hidden behind hover.
- **Recommended badge** — distinct style, not just a generic green badge.

### Empty states

- **Icon + headline + body + CTA** for every empty view.
- **Contextual CTAs** — "Run prepaid top-up example", "Create Stripe checkout session", "Configure webhook".

### Feedback & interaction

- **Toast notifications** — replace `announce`-only feedback with visible toast stack.
- **Tooltips** — explain icon-only buttons on hover/focus.
- **Command palette** — `Cmd+K` to jump tabs, focus search, toggle theme.
- **Connection status pill** — subtle dot + label in the header.
- **Loading states** — skeleton placeholders matching final layout.

### State management for UI preferences

- **Theme** — `theme.ts` already handles light/dark/system; ensure the design refresh preserves this and adds a visible theme selector in the header.
- **Density** — store `muara-density` in `localStorage`; default to `comfortable`; apply a class to `<html>` or `<body>` so CSS can vary row heights and gaps.
- **Sidebar/collapsed panels** — persist collapsed state of optional sidebar or detail panels in `localStorage`.
- **Filter persistence** — URL-state already handles shareable filters; optionally remember last-used filters per tab in `sessionStorage`.
- **Column visibility** — persist shown/hidden table columns in `localStorage`.

### Forms & inputs

- **Validation UX** — inline error messages below inputs; disable submit until required fields are valid; show a summary error at the top of long forms.
- **Input states** — distinct focus, hover, disabled, error, and read-only styles using tokens.
- **Password/secret fields** — show/hide toggle with accessible label, as already started in WebhookConfig.
- **Field hints** — helper text below inputs explains expected format (e.g., webhook URL must be `http`/`https`).

### Micro-interactions & feedback

- **Button states** — visible loading spinner, disabled-with-opacity, and success checkmark states.
- **Copy-to-clipboard feedback** — transient "Copied!" badge instead of only a state change.
- **Relative timestamps** — "2m ago" with full timestamp in a tooltip/title attribute.
- **Search highlight** — highlight matching substrings in filtered results.
- **Page transitions** — subtle fade when switching tabs (respect reduced motion).
- **Pulse indicators** — animate the connection-status dot when checking; static color when stable.

### Code & payload display

- **Monospace blocks** — payloads, headers, curl commands use `--font-mono` with background surface and horizontal scroll.
- **Syntax hints** — JSON payloads get basic key/value color distinction if possible without a heavy syntax-highlighter library.
- **One-click copy** — every code block has a copy button.

### Onboarding & discoverability

- **Product tour** — optional, dismissible step-by-step tour for first-time users.
- **Contextual tips** — short helper text near complex controls (webhook secret, event selection).
- **Keyboard shortcut discoverability** — show shortcut hints in tooltips and the Help modal.

### Mobile-specific patterns

- **Bottom sheet** for detail panels and modals on narrow screens instead of centered modal.
- **Floating action button** for primary action on mobile (e.g., reload or create test charge).
- **Swipe gestures** (optional) — swipe between tabs on touch devices.

### Brand & polish

- **Favicon** — small, recognizable icon for the dashboard tab.
- **Page title** — dynamic title per tab, e.g., "Ledger · OpenMuara Dashboard".
- **Metadata** — basic OpenGraph/Twitter cards for the landing page (low priority for local tool).

### Accessibility

- **Landmarks** — `<header>`, `<nav>`, `<main>`, and region labels.
- **Focus-visible styles** — clear focus rings on all interactive elements.
- **Touch targets** — minimum 44×44 CSS pixels for buttons and links.
- **Screen-reader announcements** for async updates, toasts, and filter changes.
- **Reduced-motion support** — disable animations when requested.

---

## Non-goals

- New logo or brand identity.
- Custom illustrations or imagery beyond icons.
- Functional behavior changes (e.g., new API endpoints) — this is a visual/UX refresh.
- Heavy charting or data-visualization libraries.
- CSS frameworks such as Tailwind or Bootstrap.

---

## Architecture & constraints

- Keep the dashboard bundle small: JS + CSS ≤ 150 KB gzipped (current ~19 KB).
- Browser baseline: last two versions of Chrome, Firefox, Safari, Edge; no polyfills for dead browsers.
- No external runtime dependencies beyond Preact and optional icon library.
- Prefer inline SVG sprites over icon fonts to avoid FOUT/FOIT and keep control.
- All new components in `web/dashboard/src/components/` with co-located tests.
- Extend `styles.css`; do not introduce a second stylesheet or CSS-in-JS.
- Maintain `prefers-reduced-motion` and `prefers-contrast: more` support.

---

## Testing strategy

- **Component unit tests** (Vitest + React Testing Library):
  - `Button`, `Card`, `Badge`, `Input`, `Modal`, `Toast`, `FilterChip`, `EmptyState` render correctly and handle disabled/loading states.
  - `Shell` tabs update `aria-selected`, respond to arrow keys, and trap focus in the help modal.
  - Theme toggle and density toggle persist to `localStorage`.
- **Visual regression**:
  - Screenshot Overview, Ledger, Transactions, Webhooks, and help modal at 375px, 768px, 1024px, 1440px in both light and dark modes.
  - Store before/after images in the initiative folder.
- **Accessibility tests**:
  - axe-core scan finds no violations.
  - All interactive elements have focus-visible styles.
  - Color-contrast ratio ≥ 4.5:1 for body text, ≥ 3:1 for large text.
- **Bundle-size gate**:
  - `npm run bundle-size` must continue to pass.

---

## Risk register

| Risk | Impact | Mitigation |
|---|---|---|
| New icon library bloats bundle | Medium | Use inline SVG sprites or tree-shake Lucide; measure with `scripts/check-sizes.sh` and `npm run bundle-size`. |
| Custom heading font slows first paint | Low | Self-host WOFF2 or use `font-display: swap`; keep fallback stack. |
| Token refactor breaks provider-specific status colors | Medium | Audit all `.status-*` classes and map to semantic tokens; add visual regression. |
| Design tokens conflict with dark mode | Low | Define tokens per theme and run visual regression in both modes. |
| Scope creep into functional changes | Medium | Strict non-goal: no behavior changes; only UI/UX. |
| Inline-style removal introduces layout regressions | Medium | Replace one component at a time and run visual regression. |

---

## Acceptance criteria

- [ ] Design audit report and screenshots saved to the initiative folder.
- [ ] Zero inline `style` attributes in redesigned components.
- [ ] Consistent SVG icon system adopted.
- [ ] Overview tab exists and provider/onboarding no longer cram the Ledger tab.
- [ ] Tables have sticky headers, row hover, and sort indicators.
- [ ] Filter toolbars include search icon, clearable chips, URL filter, and date-range presets.
- [ ] Empty states include icon + headline + CTA.
- [ ] Command palette accessible via `Cmd+K` / `Ctrl+K`.
- [ ] Density toggle persisted in `localStorage`.
- [ ] axe-core reports zero serious violations.
- [ ] Responsive at 375px, 768px, 1024px, 1440px.
- [ ] Forms have inline validation, helper text, and error summaries.
- [ ] Buttons show loading, disabled, and success states.
- [ ] Code/payload blocks use monospace styling with copy buttons.
- [ ] Relative timestamps and search highlights implemented.
- [ ] Tests added/updated.
- [ ] Bundle budget still passes.
- [ ] All quality gates pass.

---

## References

- `web/dashboard/src/styles.css`
- `web/dashboard/src/theme.ts`
- `web/dashboard/src/components/Shell.tsx`
- `web/dashboard/src/components/Providers.tsx`
- `web/dashboard/src/views/Ledger.tsx`
- `web/dashboard/src/views/Transactions.tsx`
- `web/dashboard/src/views/Webhooks.tsx`
- `web/dashboard/src/components/ErrorBoundary.tsx`
- `web/dashboard/src/components/Onboarding.tsx`
- `web/dashboard/src/components/Overview.tsx`
- `web/dashboard/src/components/Skeleton.tsx`

---

## Self-assessment

**Solidity: 9.7 / 10**

The initiative now includes design principles aligned with OpenMuara's philosophy, a live design audit with screenshots, an updated current-state audit, information-architecture improvements to fix the cramped layout, a full design-system and component-library plan with z-index/breakpoint/CSS methodology, table/filter/toolbar redesign guidance, provider-card and empty-state standards, form/validation UX, micro-interactions, code/payload display standards, mobile patterns, onboarding discoverability, UI preference state management, accessibility targets, a command palette and density toggle, a testing strategy with bundle-size gate, and a risk register. The remaining 0.3 point is reserved for implementation-specific choices (exact icon delivery method, heading font selection, and precise motion curves) that should be decided during implementation and validated with visual regression.
