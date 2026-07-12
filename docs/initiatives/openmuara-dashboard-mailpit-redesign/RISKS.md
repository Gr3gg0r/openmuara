> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard — Mailpit-Style Redesign — Risk Register

> **Updated:** 2026-07-06

| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
|----|------|------------|--------|------------|-------|--------|
| R001 | Provider metadata enrichment (base URL per version, env var mapping) requires backend changes that may conflict with existing wizard/provider code. | Medium | Medium | Keep changes additive; reuse `VersionedProvider` interface; derive env var names from a documented convention. | AI Agent | Open |
| R002 | Bundle size may exceed 150 KB gzipped after adding new views and components. | Low | Medium | Use native APIs, avoid new dependencies, run `check-bundle-size.js` in CI. | AI Agent | Open |
| R003 | Removing the Overview tab may regress first-time onboarding discoverability. | Medium | Low | Surface the onboarding checklist on the Ledger empty state or as a dismissable banner. | AI Agent | Open |
| R004 | Active/inactive terminology may be confused with the singular "active provider" concept in the backend. | Medium | Low | UI label uses "Enabled/Disabled"; keep backend "active provider" logic unchanged. | AI Agent | Open |
| R005 | URL state migration from `tab=` to `view=` may break existing bookmarks. | Low | Low | Support `tab=` as a fallback redirect for one release cycle. | AI Agent | Open |
| R006 | Dual-port runtime changes core server startup and config validation. | Medium | Medium | Make `admin_port` optional; preserve single-port behavior when unset; add integration tests for both modes. | AI Agent | Open |
| R007 | Detail pages require new backend detail endpoints or reuse existing ones; existing transaction/webhook detail endpoints already exist under `/_admin/transactions/{ref}` and `/_admin/webhooks/{ref}`. | Low | Low | Reuse existing detail endpoints; only the dashboard navigation layer changes. | AI Agent | Open |

---

## Risk Template

```markdown
| ID | Risk | Likelihood | Impact | Mitigation | Owner | Status |
```
