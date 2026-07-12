> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Documentation Website — Decision Log

| ID | Decision | Status | Date | Notes |
|----|----------|--------|------|-------|
| D001 | Static-site generator selection | ✅ Decided | 2026-07-02 | Docusaurus (Option B). Separate system, not embedded in Go binary. |
| D002 | Docs remain Markdown in root repo | ✅ Decided | 2026-07-02 | Non-negotiable for contributor accessibility. |
| D003 | Deployment target | ✅ Decided | 2026-07-02 | GitHub Pages via GitHub Actions. |
| D004 | Docusaurus build weight is acceptable | ✅ Decided | 2026-07-03 | Docusaurus is heavier than VitePress/MkDocs, but it is acceptable because the docs site is a separate system deployed to GitHub Pages. It does not ship in the Go binary and does not affect the emulator's runtime memory or bundle size. |
