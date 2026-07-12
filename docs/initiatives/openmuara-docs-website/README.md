> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Documentation Website

> **Status:** 🟡 In Progress | **Started:** 2026-07-02
> **Scope:** Create a standalone, searchable documentation website for OpenMuara from the existing Markdown docs in `docs/`. The site runs as a separate system (not embedded in the Go binary) and is deployed to GitHub Pages.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/docs-website`

---

## Initiative Structure

```
docs/initiatives/openmuara-docs-website/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    └── 01-generator-decision.md
```

Planning docs live in `docs/initiatives/openmuara-docs-website/` in the root repo.
Website build commits to the `feat/docs-website` branch. Do not commit directly to `main`.

---

## Current State

Documentation is authored in Markdown under `docs/`:

- `docs/architecture.md`
- `docs/providers/*.md`
- `docs/migration/*.md`
- `runbooks/*.md`
- `README.md`
- `CHANGELOG.md`

These render well in GitHub and IDEs but are not a cohesive, searchable web experience for end users.

---

## Goals

1. Choose a static-site generator suitable for technical docs.
2. Produce a deployable documentation site from existing Markdown with minimal rewrite.
3. Add search, navigation, and version-aware URLs.
4. Automate deployment via GitHub Actions (e.g., GitHub Pages or Vercel).
5. Keep docs source in Markdown so non-engineers can contribute.

---

## Options

### Option A — VitePress

- [VitePress](https://vitepress.dev/) (Vue/Vite-based).

**Pros:** fast, modern, great default theme, full-text search plugin, Vue ecosystem.  
**Cons:** requires Node; Vue-specific if custom components are needed.

### Option B — Docusaurus

- [Docusaurus](https://docusaurus.io/) (React-based).

**Pros:** mature, versioning support, blog, plugins, large community.  
**Cons:** heavier than VitePress; slower build; more config.

### Option C — MkDocs Material

- [MkDocs](https://www.mkdocs.org/) with [Material theme](https://squidfunk.github.io/mkdocs-material/).

**Pros:** Python-based (no Node), excellent Material Design, search built-in, very stable.  
**Cons:** less "modern" JS interactivity; custom theming is Jinja2/Python.

### Option D — Hugo + Docsy

- [Hugo](https://gohugo.io/) with [Docsy](https://www.docsy.dev/) theme.

**Pros:** extremely fast builds, Go-based (matches project stack), deploys anywhere.  
**Cons:** Docsy setup can be finicky; Go templating learning curve.

### Option E — Keep GitHub-rendered Markdown

- Do nothing; rely on GitHub's Markdown rendering and `docs/` structure.

**Pros:** zero work, zero tooling.  
**Cons:** no search, no versioning, no branded experience, poor navigation.

---

## Recommendation

**Use Option B — Docusaurus.**

Docusaurus is the chosen generator because it provides:

- mature React-based documentation site out of the box,
- built-in versioning and i18n support for future v1/v2 docs,
- full-text search via Algolia DocSearch,
- blog support for release notes and tutorials,
- large community and plugin ecosystem.

The docs website is a **separate system** from the OpenMuara Go binary. It does not ship inside the binary and does not add to its bundle size. It is built and deployed independently to **GitHub Pages** via GitHub Actions.

Options A, C, and D remain documented as alternatives in case Docusaurus proves too heavy during implementation.

---

## Non-Goals

- Do not move docs out of the root repo.
- Do not require docs authors to learn JSX, Vue, or Python templating.
- Do not make the docs website a runtime dependency of the OpenMuara binary.
- Do not embed the docs website into the Go binary; it is a separate build/deploy system.
