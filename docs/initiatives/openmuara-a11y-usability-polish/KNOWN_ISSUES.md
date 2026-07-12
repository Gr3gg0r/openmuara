> **Initiative:** OpenMuara Accessibility & Usability Polish

# Known Issues

These issues were discovered during the accessibility/usability audit on 2026-07-03. The detailed findings live in [`findings/2026-07-03-a11y-usability-audit.md`](findings/2026-07-03-a11y-usability-audit.md).

## High severity

- Clickable table rows in Ledger, Transactions, and Webhooks views are not keyboard-accessible.
- Search inputs in Ledger and Transactions views lack persistent labels.
- Help modal lacks focus trapping and initial focus management.
- Dashboard tabs lack arrow-key navigation.
- Theme shortcut can leave the toggle button's icon/label out of sync.

## Medium severity

- "Copy curl" buttons in the Providers list are not uniquely labelled.
- Copy-to-clipboard success state is not announced to screen readers.
- Onboarding Show/Hide button lacks `aria-expanded`.
- Onboarding panel has an invalid hard-coded light background in dark mode.
- Failed-webhook alert uses a fake link and an unlabelled emoji.
- Detail panels do not manage focus on open/close.
- Provider simulation pages lack button focus styles and `<main>` landmarks.
- Example mini-apps do not announce status updates to screen readers.

## Low severity

- Theme toggle is icon-only.
- Global keyboard shortcuts do not ignore modifier keys.
- Example mini-app inputs are not wrapped in `<form>` elements.
