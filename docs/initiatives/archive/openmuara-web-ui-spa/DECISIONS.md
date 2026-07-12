> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Decision Log

| ID | Decision | Status | Date | Notes |
|----|----------|--------|------|-------|
| D001 | Dashboard framework selection | ✅ Decided | 2026-07-02 | Vite + Preact (Option C). Embedded into Go binary at build time. |
| D002 | Build assets must be embeddable by Go | ✅ Decided | 2026-07-02 | No runtime Node dependency; `//go:embed dist`. |
| D003 | Preserve existing `/_admin` route stability | ✅ Decided | 2026-07-02 | Escape/pay URLs must not change. |
| D004 | SPA must stay lightweight in the Go binary | ✅ Decided | 2026-07-03 | Preact runtime is ~3 KB. Build output must stay small, avoid heavy deps, and use lazy loading/code splitting. Memory footprint must remain low because the SPA runs inside the same process as the Go emulator. |
| D005 | TypeScript for dashboard source | ✅ Decided | 2026-07-03 | Compile-time type safety with no runtime overhead. Vite transpiles TS transparently. |
| D006 | Escape/pay pages remain server-rendered | ✅ Decided | 2026-07-03 | Payment flows stay reliable without JS; `internal/ui/` keeps the existing HTML templates. Only the dashboard migrates to the SPA. |
| D007 | Source and build directory layout | ✅ Decided | 2026-07-03 | Source in `web/dashboard/`, build output in `web/dashboard/dist/`, embedded via `internal/ui/embed.go`. |
| D008 | Bundle-size budget | ✅ Decided | 2026-07-03 | Initial JS ≤ 100 KB gzipped; total `dist/` ≤ 250 KB. Enforced by `scripts/check-bundle-size.sh` in CI. |
