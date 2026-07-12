> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Handoff

## Current Status

- Initiative created on `feat/web-ui-spa`.
- Prompt 01 completed: Vite + Preact selected.
- No product code changes yet.

## Resolved Questions

1. **Escape/pay pages:** Kept as server-rendered HTML in `internal/ui/`. Payment flows must work without JavaScript, and converting them adds complexity without benefit. The dashboard SPA links to them.
2. **TypeScript:** Yes. It provides compile-time safety with zero runtime cost and is the natural choice for Vite-based projects.
3. **Site directory:** `web/dashboard/` with build output at `web/dashboard/dist/`.

## Audit Notes (2026-07-03)

- Low-memory/efficiency posture confirmed: Preact (~3 KB runtime) keeps bundle and runtime footprint small.
- Decision D004 added: SPA must stay lightweight because it is embedded in the Go binary and shares its process.
- Risk R03 tightened: enforce bundle-size budget, tree-shaking, lazy loading, and avoid heavy transitive dependencies.
- Added bundle-size budget: initial JS ≤ 100 KB gzipped, total `dist/` ≤ 250 KB.
- Added UI testing strategy: Vitest.
- Added CI strategy: separate `ui-build` and `ui-test` jobs; Go build depends on built UI assets.
- No changes to framework choice.

## Next Actions

- [x] Scaffold the Vite + Preact build pipeline in prompt 02.
- [x] Decide whether escape/pay pages move into the SPA or stay server-rendered.
- [x] Set bundle-size budget and CI check when build pipeline is scaffolded.
- [ ] Implement dashboard shell and view migration.
- [ ] Wire UI build into Taskfile, CI, and Go embed.
