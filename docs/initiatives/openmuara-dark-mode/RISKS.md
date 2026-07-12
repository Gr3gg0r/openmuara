> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first.**

# OpenMuara Dark Mode — Risk Register

---

| ID | Risk | Likelihood | Impact | Mitigation | Owner |
|----|------|------------|--------|------------|-------|
| R01 | Flash of un-themed content before JavaScript applies dark mode | Medium | Medium | Inject a small blocking script in `<head>` to set the theme class before first paint. | AI Agent |
| R02 | Hard-coded colors scattered across dashboard CSS make theming inconsistent | High | Medium | Audit and convert every hard-coded color to a semantic custom property before adding dark values. | AI Agent |
| R03 | Provider pay pages use hard-coded light colors and no custom properties | Medium | Medium | Convert pay pages to CSS custom properties; add `color-scheme: light dark`. | AI Agent |
| R04 | Dashboard React tests fail because they expect light-mode classes or colors | Low | Low | Update tests to assert theme class presence and behavior, not specific colors. | AI Agent |
| R05 | Bundle size exceeds 5 KiB target | Low | Low | Avoid images/icons for the toggle; use inline SVG or Unicode; keep the theme script small. | AI Agent |
| R06 | Embedded dashboard (`internal/ui/dashboard-dist/`) drifts from `web/dashboard/` | Medium | Medium | Rebuild the dashboard and commit `internal/ui/dashboard-dist/` as part of P01. | AI Agent |
| R07 | Contrast failures in dark mode for status badges, alerts, or buttons | Medium | High | Use a contrast checker during implementation; add a manual visual QA gate. | AI Agent |
| R08 | Theme toggle state not synced across tabs | Low | Low | Listen to `storage` event and re-apply theme on change. | AI Agent |
| R09 | Pay pages flash white because they have no head script | Medium | Medium | Add `<meta name="color-scheme" content="light dark">` and dark background via custom properties in the first paint. | AI Agent |
| R10 | `prefers-reduced-motion` users get jarring color transitions | Low | Low | Wrap transition in `@media (prefers-reduced-motion: no-preference)`. | AI Agent |
