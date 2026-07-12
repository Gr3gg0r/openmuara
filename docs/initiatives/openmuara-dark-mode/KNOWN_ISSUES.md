> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Dark Mode — Known Issues

> **Status:** All original gaps resolved as of 2026-07-03.

## Resolved gaps

| # | Gap | Resolution |
|---|-----|------------|
| 1 | No dark mode existed | Dashboard, provider pay pages, and example mini-apps now support dark mode. |
| 2 | No theme toggle | Dashboard has a manual toggle next to the Help button. |
| 3 | Hard-coded dashboard colors | All dashboard colors converted to semantic CSS custom properties. |
| 4 | Provider pay pages not tokenized | All seven pages converted to CSS custom properties. |
| 5 | Example mini-apps not tokenized | Both examples converted to CSS custom properties. |
| 6 | No `color-scheme` meta tag | Added to dashboard, provider pages, and examples. |
| 7 | No reduced-motion handling | Color transitions are wrapped in `prefers-reduced-motion: no-preference`. |

## Post-ship notes for future work

- Example mini-apps follow OS preference only; if user demand arises, a manual toggle can be added per D006.
- Provider pay pages follow OS preference only; this matches the transient nature of simulation screens.
- No CLI terminal theming was added; that remains out of scope.
