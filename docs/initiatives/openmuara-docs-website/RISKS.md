> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Documentation Website — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | Generator becomes unmaintained. | Low | Medium | Choose established tools (VitePress, Docusaurus, MkDocs Material). |
| R02 | Docs drift between Markdown source and live site. | Medium | High | Automate deployment on every merged PR; make live site the canonical link. |
| R03 | Internal links break during migration. | Medium | Medium | Run a link checker in CI; use generator's link validation. |
| R04 | Authors avoid contributing because of new frontmatter/tooling. | Medium | Medium | Keep Markdown as the source; document the minimal frontmatter required. |
| R05 | Deployment secrets or tokens complicate CI. | Low | Low | GitHub Pages deployment uses built-in `GITHUB_TOKEN`; no extra secrets. |
| R06 | Search indexing excludes important provider docs. | Low | Medium | Configure search to include all `docs/` content; test queries manually. |
| R07 | Docusaurus build is slow or memory-heavy in CI. | Medium | Low | Cache `node_modules` and build artifacts; monitor build time and memory in GitHub Actions; consider VitePress or MkDocs if builds become a bottleneck. |
