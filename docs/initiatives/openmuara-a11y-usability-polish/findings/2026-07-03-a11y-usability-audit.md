> **Initiative:** OpenMuara Accessibility & Usability Polish
> **Date:** 2026-07-03
> **Auditor:** AI Agent (Kimi Code)

# Accessibility & Usability Audit Findings

## Method

1. Started a local OpenMuara server on a temporary workspace.
2. Inspected the dashboard with the browser accessibility tree and a custom script checking for unlabelled inputs, unnamed buttons, and new-tab links.
3. Read the source of every dashboard component, provider simulation page, and example mini-app.
4. Ranked findings by severity for keyboard users, screen-reader users, and visual usability.

## Findings

### 🔴 High severity

| # | File | Issue | Why it matters |
|---|------|-------|----------------|
| H1 | `web/dashboard/src/views/Ledger.tsx` | Ledger table rows are `<tr onClick>` with `class="row-click"`. They are not focusable. | Keyboard users cannot open the detail panel. Screen readers don't know the row is interactive. |
| H2 | `web/dashboard/src/views/Transactions.tsx` | Same clickable-row pattern for transactions. | Same as H1. |
| H3 | `web/dashboard/src/views/Webhooks.tsx` | Same clickable-row pattern for webhooks. | Same as H1. |
| H4 | `web/dashboard/src/views/Ledger.tsx` | Search `<input>` has only a `placeholder`, no `<label>` or `aria-label`. | Placeholders disappear when typing; screen readers may not announce the field purpose reliably. |
| H5 | `web/dashboard/src/views/Transactions.tsx` | Search `<input>` has only a `placeholder`. | Same as H4. |
| H6 | `web/dashboard/src/components/Shell.tsx` | Help modal has no focus trap and does not move focus on open. | Keyboard users can tab behind the modal; screen-reader users are not informed it opened. |
| H7 | `web/dashboard/src/components/Shell.tsx` | Tab bar uses `role="tab"` but has no arrow-key navigation. | Expected ARIA tab-list keyboard behavior is missing. |
| H8 | `web/dashboard/src/app.tsx` | The `d` shortcut calls `toggleTheme()` but `Shell` keeps its own `theme` state, so the toggle button's icon and `aria-label` stay stale. | Visual and announced state become inconsistent. |

### 🟡 Medium severity

| # | File | Issue | Why it matters |
|---|------|-------|----------------|
| M1 | `web/dashboard/src/components/Providers.tsx` | "Copy curl" buttons repeat for every provider with identical text. | Screen-reader users can't tell which provider's curl will be copied. |
| M2 | `web/dashboard/src/components/Providers.tsx` | "Copied!" state is visual only; no `aria-live` announcement. | Screen-reader users get no confirmation. |
| M3 | `web/dashboard/src/components/Onboarding.tsx` | Show/Hide button lacks `aria-expanded`. | Users don't know whether the panel is currently expanded. |
| M4 | `web/dashboard/src/components/Onboarding.tsx` | Inline style `background:'#f8fafc'` is invalid in dark mode and keeps a light background. | Visual bug in dark mode. |
| M5 | `web/dashboard/src/components/FailedWebhookAlert.tsx` | "Webhooks" control is an `<a href="#">` with `onClick`. | Semantically a button; may behave unexpectedly with assistive tech. |
| M6 | `web/dashboard/src/components/FailedWebhookAlert.tsx` | Uses emoji `⚠️` without alternative text. | Screen-reader announcement is inconsistent across platforms. |
| M7 | All detail panels | Opening/closing a detail panel does not move focus or return focus. | Disorienting for keyboard and screen-reader users. |
| M8 | `internal/ui/*.html` | Buttons have no `:focus` styles (only inputs/selects do). | Keyboard users can't see which button is focused. |
| M9 | `internal/ui/*.html` | Provider pages lack a `<main>` landmark and skip link. | Less efficient navigation for assistive tech. |
| M10 | `examples/ecommerce-single-buy/index.html` | Status messages (`#status`) are not `aria-live` regions. | Screen-reader users miss "Creating..." / error updates. |
| M11 | `examples/prepaid-topup/index.html` | Status messages (`#status`) are not `aria-live` regions. | Same as M10. |

### 🟢 Low severity / polish

| # | File | Issue | Recommendation |
|---|------|-------|----------------|
| L1 | `web/dashboard/src/components/Shell.tsx` | Theme toggle is icon-only. | Add a visually hidden label or tooltip. |
| L2 | `web/dashboard/src/app.tsx` | Global shortcuts don't check modifier keys. | Ignore `1`/`2`/`3`/`d`/`/` when Ctrl/Alt/Meta are pressed. |
| L3 | `internal/ui/*.html` | Theme script is duplicated in every provider page. | Keep duplication for zero-runtime-dependency pages, or extract if templates support includes. |
| L4 | `examples/ecommerce-single-buy/index.html` | Inputs are not wrapped in `<form>`. | Wrap them so Enter submits naturally. |
| L5 | `examples/prepaid-topup/index.html` | Inputs are not wrapped in `<form>`. | Same as L4. |

## Recommended grouping for implementation

- **P01 — Dashboard keyboard navigation:** H1, H2, H3, H6, H7, M7
- **P02 — Dashboard labels and live regions:** H4, H5, M1, M2, M3, M5, M6
- **P03 — Provider pages focus and landmarks:** M8, M9
- **P04 — Example apps accessibility:** M10, M11, L4, L5
- **P05 — Shortcuts and theme polish:** H8, L1, L2, M4
