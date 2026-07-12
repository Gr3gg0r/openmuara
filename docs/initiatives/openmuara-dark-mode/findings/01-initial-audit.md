> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# Dark Mode Initial Audit

> **Date:** 2026-07-03
> **Auditor:** AI Agent (Kimi Code)
> **Scope:** Dashboard, provider pay pages, example mini-apps

---

## Summary

OpenMuara currently has no dark mode. Several surfaces already use a small set of CSS custom properties, but most colors are hard-coded hex values. A lightweight, token-based approach will cover all surfaces with minimal bundle impact.

## Dashboard (`web/dashboard/`)

- **Stack:** Preact + Vite.
- **Entry points:**
  - Dev: `web/dashboard/index.html`
  - Production (embedded): `internal/ui/dashboard-dist/index.html`
- **Current tokens (in `styles.css`):**
  - `--bg`, `--card`, `--text`, `--muted`, `--border`
- **Hard-coded colors found:**
  - Badge backgrounds/text (`#dcfce7`, `#166534`, `#e0f2fe`, `#075985`)
  - Primary button (`#2563eb`, `#1d4ed8`)
  - Secondary/segment buttons (`#f1f5f9`, `#e2e8f0`)
  - Table header (`#f1f5f9`)
  - Status text (`#166534`, `#991b1b`, `#854d0e`, `#475569`, `#4338ca`)
  - Modal overlay (`rgba(15, 23, 42, 0.5)`)
  - Help box `kbd` background (`#f1f5f9`)
  - Alert background/border/text (`#fef2f2`, `#fecaca`, `#991b1b`)
  - Error banner (`#fef2f2`, `#fecaca`, `#991b1b`)
  - Copy button copied state (`#16a34a`)
- **Toggle location:** Add to `Shell.tsx`, next to the Help button.
- **Tests:** `web/dashboard/tests/Shell.test.tsx` should be extended.

## Provider pay pages (`internal/ui/*.html`)

- **Files:**
  - `stripe-checkout.html`
  - `stripe-payment-intent.html`
  - `stripe-webhooks.html`
  - `fawry-escape.html`
  - `billplz-pay.html`
  - `toyyibpay-pay.html`
  - `ipay88-pay.html`
- **Findings:**
  - Most pages use a small `:root` token set (`--bg`, `--card`, `--text`, `--muted`, `--border`) but hard-code button/input colors.
  - `fawry-escape.html` has no tokens at all and uses hard-coded success/cancel button colors.
  - No `<meta name="color-scheme">` tags.

## Example mini-apps

- **Files:**
  - `examples/ecommerce-single-buy/index.html`
  - `examples/prepaid-topup/index.html`
- **Findings:**
  - Both use hard-coded colors (`#ddd`, `#0a0`, `#666`).
  - No custom properties.
  - No `color-scheme` meta.

## Recommendations

1. Define a shared semantic token set in the dashboard first; reuse the same names in pay pages and examples where possible.
2. Convert hard-coded colors to tokens before adding dark values.
3. Use a blocking `<head>` script in the dashboard to prevent theme flash.
4. Add `<meta name="color-scheme" content="light dark">` to all HTML pages.
5. Add cross-tab sync and `prefers-reduced-motion` handling.
6. Verify WCAG AA contrast for every status badge and button state.
