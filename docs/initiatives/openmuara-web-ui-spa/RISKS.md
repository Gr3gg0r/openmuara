> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | SPA build step complicates the local development workflow. | Medium | Medium | Provide `make dev` that runs Vite dev server and `muara start` together; document clearly. |
| R02 | Contributors without Node installed are blocked from UI changes. | Medium | Medium | Keep the UI build optional for Go-only changes; CI builds and checks in assets. |
| R03 | Bundle size bloats the Go binary. | Low | Medium | Use Preact or Svelte; tree-shake; monitor bundle size in CI. |
| R04 | Escape/pay pages break or lose CSRF protection. | Medium | High | Migrate one page at a time with full test coverage; preserve existing token plumbing. |
| R05 | SSE auto-refresh logic regresses. | Medium | High | Reuse existing `/events` endpoint; add explicit tests for live updates. |
| R06 | Framework choice becomes unsupported in 2–3 years. | Low | Low | Prefer established, lightweight frameworks (Preact/Svelte) with long-term maintenance. |
