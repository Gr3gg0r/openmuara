> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara AI Slop Audit — Tracking

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

## Phases

| Phase | Title | Goal | Status |
|-------|-------|------|--------|
| P01 | Dashboard microcopy & icons | Fix placeholder/redundant copy and off-tone icons | ⬜ Not started |
| P02 | Provider metadata | Normalize `gateway.yml` descriptions and categories | ⬜ Not started |
| P03 | Documentation & prompts | Remove buzzwords and placeholder sections | ⬜ Not started |
| P04 | Code-level slop | Reduce `any`/`interface{}`, boilerplate, noise comments | ⬜ Not started |
| P05 | Test / seed data | Gate seeding and make fixtures realistic | ⬜ Not started |
| P06 | Examples & website | Replace generic store and Docusaurus copy | ⬜ Not started |

## Findings log

| ID | Finding | Area | Status | Fixed in |
|----|---------|------|--------|----------|
| F001 | Sad-face empty state | Dashboard | ⬜ Open | — |
| F002 | "Provider configuration" placeholders | Dashboard/Metadata | ⬜ Open | — |
| F003 | Redundant ACTIVE + ENABLED badges | Dashboard | ⬜ Open | — |
| F004 | Generic system-ui font | Dashboard | ⬜ Open | — |
| F005 | Preseed transactions visible on first load | Dashboard/Seed data | ⬜ Open | — |
| F006 | Header action buttons too small | Dashboard | ✅ Verified | 2026-07-08 responsiveness pass |
| F007 | Inconsistent provider descriptions | Metadata | ⬜ Open | — |
| F008 | Buzzword-heavy docs/prompts | Docs | ⬜ Open | — |
| F009 | Placeholder TODOs in docs | Docs | ⬜ Open | — |
| F010 | `any` / `interface{}` usage | Code | ⬜ Open | — |
| F011 | Copy-paste provider boilerplate | Code | ⬜ Open | — |
| F012 | Verbose comments | Code | ⬜ Open | — |
| F013 | Synthetic seed transactions | Test data | ⬜ Open | — |
| F014 | Generic example store | Examples | ⬜ Open | — |
| F015 | Generic Docusaurus copy | Website | ⬜ Open | — |

## Quality gates

Every phase must end with:

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] `golangci-lint run`
- [ ] `npm run typecheck` (in `web/dashboard/`)
- [ ] `npm run test:ci` (in `web/dashboard/`)
- [ ] Before/after screenshots for any UI change
