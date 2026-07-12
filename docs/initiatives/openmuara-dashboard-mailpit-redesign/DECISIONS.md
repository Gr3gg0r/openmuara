> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard — Mailpit-Style Redesign — Decision Log

> **Updated:** 2026-07-06

| ID | Date | Decision | Rationale | Status |
|----|------|----------|-----------|--------|
| D001 | 2026-07-06 | Left navigation with Ledger, Webhooks, Settings. | User explicitly requested Mailpit-like layout with three primary nav items. | Decided |
| D002 | 2026-07-06 | Ledger is the default landing view. | User wants the outlet to default to the full ledger table. | Decided |
| D003 | 2026-07-06 | Provider config lives under Settings, separate from Webhooks delivery log. | Keeps configuration (stateless edits) distinct from operational logs. | Decided |
| D004 | 2026-07-06 | Environment variables shown as read-only reference names. | Avoids exposing secrets; provides copy-paste convenience. | Decided |
| D005 | 2026-07-06 | Version tabs appear only when provider reports multiple versions. | Matches existing `provider.VersionedProvider` contract; keeps UI simple for single-version providers. | Decided |
| D006 | 2026-07-06 | Design priority stack: UI > UX > performance > usability > philosophy > efficiency > memory size. | Explicitly captures how to resolve trade-offs during implementation; UI density and clarity lead, memory size is a constraint not a driver. | Decided |
| D007 | 2026-07-06 | Every table view has a reusable filter toolbar. | User explicitly requested filters on every table; keeps UX consistent across Ledger and Webhooks. | Decided |
| D008 | 2026-07-06 | Ledger and webhook rows navigate to dedicated detail pages. | User requested detail pages rather than inline panels; improves deep-linking and focus management. | Decided |
| D009 | 2026-07-06 | Admin UI and provider endpoints run on separate optional ports. | User requested Mailpit-style two-port setup for easy external exposure of the API only. | Decided |
| D010 | 2026-07-06 | Per-provider webhook targets live in Settings → Provider Detail. | User clarified the top-level Webhooks view should be a delivery log only, keeping configuration out of operational views. | Decided |
| D011 | 2026-07-06 | Light-theme muted text and shortcut colors darkened to pass WCAG AA contrast. | axe-core reported contrast failures on zebra table rows and command shortcuts; fixing improves readability without harming visual hierarchy. | Decided |

---

## Decision Template

```markdown
| ID | Date | Decision | Rationale | Status |
```
