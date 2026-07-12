---
id: accessibility
title: Accessibility
---

# Accessibility

OpenMuara aims to meet WCAG 2.1 AA in the admin dashboard and provider
simulation pages. This document summarizes the current state and how we verify
it.

## Current features

### Admin dashboard (`/_admin`)

- **Skip link:** keyboard users can jump to the main content area.
- **Landmarks and labels:** navigation, main content, and overview sections use
  semantic elements and `aria-label` attributes.
- **Live region:** status announcements use `aria-live="polite"` so screen-reader
  users are notified of async updates.
- **Focus management:** focus-visible styles and modal `aria-modal` attributes
  keep keyboard focus predictable.
- **Color contrast:** text and interactive elements target WCAG AA contrast
  ratios in both light and dark mode.
- **High contrast:** `prefers-contrast: more` increases border visibility.

### Provider simulation pages

- Decorative icons use `aria-hidden="true"`.
- Form inputs have associated labels and focus indicators.

## Running accessibility checks

The dashboard test suite includes Playwright smoke tests for keyboard
navigation, theme toggle, and axe-core critical violations:

```bash
cd web/dashboard
npm install
npm run test:a11y
```

A dedicated contrast regression check is also available:

```bash
npm run a11y:contrast
```

## Reporting issues

Open a docs issue using the **Documentation** issue template and include:

- The page URL (`/_admin/...` or provider simulation page).
- Browser / screen reader combination.
- The accessibility guideline violated (e.g., WCAG 1.4.3 contrast).
- Steps to reproduce.

## See also

- `web/dashboard/README.md` — dashboard development guide.
- `docs/operations.md` — running the admin dashboard in production-like modes.
