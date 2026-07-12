> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | SPA build step complicates the local development workflow. | Medium | Medium | Provide `task dev` that runs Vite dev server and `muara start` together; document clearly in `runbooks/local-development.md`. |
| R02 | Contributors without Node installed are blocked from UI changes. | Medium | Medium | Keep the UI build optional for Go-only changes; CI builds and checks in assets. Go tests embed a pre-built `dist/` fallback when present, and existing HTML escape/pay pages remain untouched. |
| R03 | Bundle size bloats the Go binary or increases memory use. | Low | High | Use Preact (~3 KB runtime), tree-shake, code-split/lazy-load routes, and enforce a bundle-size budget in CI (initial JS ≤ 100 KB gzipped, total `dist/` ≤ 250 KB). Avoid charting, state-management, or animation libraries. Profile memory during dashboard load and auto-refresh. |
| R04 | Escape/pay pages break or lose CSRF protection. | Low | High | Keep escape/pay pages as server-rendered Go templates in `internal/ui/`; do not migrate them. Add regression tests that verify each page still renders after the SPA migration. |
| R05 | Auto-refresh logic regresses. | Medium | High | Reuse the existing polling strategy (SSE `/events` is available but the current dashboard polls every 2 s); add explicit tests for live updates in the SPA. |
| R06 | Framework choice becomes unsupported in 2–3 years. | Low | Low | Prefer established, lightweight frameworks (Preact/Svelte) with long-term maintenance. Preact is React-compatible, reducing lock-in. |
| R07 | Built assets are out of sync with source. | Medium | Medium | CI fails if `web/dashboard/dist/` is missing or stale. `task build` runs `task ui:build` first. Commit built `dist/` only when cutting releases; otherwise rely on CI. |
| R08 | CSP from Security Hardening breaks the SPA. | Low | High | SPA has no inline scripts in production (Vite emits hashed/none) and all assets are `self`. Verify against hardened config in smoke tests. |
