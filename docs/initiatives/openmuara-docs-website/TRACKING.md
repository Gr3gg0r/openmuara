> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Documentation Website — Execution Tracker

> **Updated:** 2026-07-08 | **Status:** ✅ Completed
>
> **Scope:** Build a standalone documentation website from existing Markdown docs.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/docs-website`
> **Last Agent Action:** Docusaurus site scaffolded, docs migrated with frontmatter, search and theming configured, GitHub Pages deployment workflow added.
> **Next Agent Action:** None.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen |

---

## Execution Rules

1. Execute prompts in order unless marked **[PARALLEL SAFE]**.
2. Every prompt MUST end with: tests passing → git commit → update this file to `✅`.
3. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY prompt, update `HANDOFF.md`.
5. Website build commits happen on `feat/docs-website`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Generator decision | `DECISIONS.md`, `README.md` | — | ✅ | — | Docusaurus selected; GitHub Pages deployment; separate system from Go binary. |
| 02 | Site scaffold and nav | `website/`, `docs/` frontmatter | 01 | ✅ | — | Docusaurus 3 classic preset scaffolded; sidebars configured for docs and runbooks. |
| 03 | Migrate core docs | `docs/**/*.md`, `runbooks/*.md` | 02 | ✅ | — | Frontmatter added; docs and runbooks sourced from root repo; internal links fixed. |
| 04 | Search and theming | `website/docusaurus.config.ts`, `website/src/css/custom.css` | 03 | ✅ | — | Local search plugin enabled; OpenMuara blue theme and logo applied. |
| 05 | CI/CD deployment | `.github/workflows/docs.yml` | 04 | ✅ | — | Build and deploy to GitHub Pages on push to `main`. |
| 06 | README and redirects | `README.md` | 05 | ✅ | — | README updated with link to the live docs site. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Docs Build | `cd website && npm run build` | Produces static site | ✅ |
| Link Check | Docusaurus `onBrokenLinks: 'throw'` | No broken internal links | ✅ |
| Lint | `golangci-lint run` / `go vet` | Go code clean | ✅ |
| Deploy | `.github/workflows/docs.yml` | GitHub Pages deployment configured | ✅ |

---

## Decisions

- D001 ✅ Static-site generator: Docusaurus 3 classic preset.
- D002 ✅ Docs stay in Markdown in the root repo; Docusaurus reads from `../docs` and `../runbooks`.
- D003 ✅ Deployment target: GitHub Pages at `https://gr3gg0r.github.io/openmuara/`.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-docs-website/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-docs-website/README.md` | Goals, options, recommendation |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
