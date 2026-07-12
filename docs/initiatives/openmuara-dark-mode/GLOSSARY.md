> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Dark Mode — Glossary

---

| Term | Definition |
|------|------------|
| `prefers-color-scheme` | CSS media feature that reflects the user's OS-level light/dark preference. |
| `color-scheme` | CSS property that tells the browser which system colors to use for form controls and scrollbars. |
| `prefers-reduced-motion` | CSS media feature for users who want fewer animations; used here to disable the theme color transition. |
| Semantic token | A CSS custom property named by purpose (e.g., `--color-bg`) rather than by literal color value. |
| Theme flash / FOUC | A brief display of the wrong theme before JavaScript applies the correct one; mitigated by a blocking head script. |
| Dashboard | The Preact/Vite SPA served at `/_admin`. |
| Pay pages | OpenMuara-hosted HTML pages where users simulate payment confirmation (e.g., Stripe checkout pay page, Fawry escape page). |
| Mini-apps | The example landing pages in `examples/ecommerce-single-buy/` and `examples/prepaid-topup/`. |
| Embedded dashboard | The production build of the dashboard copied to `internal/ui/dashboard-dist/` and embedded into the Go binary. |
