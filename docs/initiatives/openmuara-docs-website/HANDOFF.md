> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Documentation Website — Handoff

## Current Status

- Initiative created on `feat/docs-website`.
- Prompt 01 completed: Docusaurus + GitHub Pages selected.
- No product code changes yet.

## Open Questions

1. Docusaurus site directory name (`website/`, `docs-site/`, or `docusaurus/`)?
2. Should the docs site have versioned docs for v1/v2 later?
3. Custom domain or `gr3gg0r.github.io/openmuara`?

## Audit Notes (2026-07-03)

- Docusaurus weight audited and accepted: the docs site is a separate system deployed to GitHub Pages, so it does not affect the Go binary size or runtime memory.
- Decision D004 added: Docusaurus build weight is acceptable because it is decoupled from the emulator.
- Risk R07 added: monitor Docusaurus build time and memory in CI; keep VitePress/MkDocs as fallback options if CI becomes a bottleneck.
- No changes to generator or deployment choice.

## Next Actions

- [ ] Scaffold the Docusaurus site in prompt 02.
- [ ] Map existing `docs/` Markdown into the Docusaurus nav/sidebar.
- [ ] Configure CI caching for `node_modules` and build artifacts to mitigate slow/heavy builds.
