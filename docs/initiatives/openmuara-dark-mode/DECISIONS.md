> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Dark Mode — Decision Log

---

## D001 — Lightweight, native theming only

- **Decision:** Implement dark mode with CSS custom properties and `prefers-color-scheme`. No external theming libraries.
- **Rationale:** Aligns with OpenMuara's low-memory, low-dependency philosophy. Keeps bundle size small and avoids framework lock-in.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented

## D002 — Semantic design tokens

- **Decision:** Name tokens by purpose (e.g., `--color-bg`, `--color-text-primary`) rather than by color value.
- **Rationale:** Makes the theme system easier to extend and keeps light/dark definitions co-located.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented

## D003 — Manual theme storage key

- **Decision:** Persist the user's manual choice in `localStorage` under the key `muara-theme`.
- **Rationale:** Simple, collision-resistant, and easy to inspect during debugging.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented

## D004 — No manual toggle on provider pay pages

- **Decision:** Provider pay pages follow OS preference only; no per-page toggle.
- **Rationale:** Pay pages are transient simulation screens; keeping them dependency-free and simple outweighs the value of a toggle.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented

## D005 — `data-theme` attribute as the single source of truth

- **Decision:** Apply theme by setting `data-theme="light" | "dark"` on `<html>`, not by toggling a class or relying on `prefers-color-scheme` media queries in the dashboard CSS.
- **Rationale:** Avoids conflicting rules between manual override and OS preference; the blocking script can set it before paint.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented

## D006 — No manual toggle in example mini-apps

- **Decision:** Example mini-apps respect `prefers-color-scheme` but do not include a manual toggle.
- **Rationale:** The examples are meant to be minimal, dependency-free demonstrations of a single payment flow; a toggle would add complexity without teaching the core integration pattern.
- **Date:** 2026-07-03
- **Status:** ✅ Accepted / Implemented
