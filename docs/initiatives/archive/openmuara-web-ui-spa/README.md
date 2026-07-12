> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA

> **Status:** ✅ Completed | **Started:** 2026-07-02 | **Completed:** 2026-07-03
> **Scope:** Migrate the embedded `/_admin` dashboard from vanilla HTML/JS to a Vite + Preact SPA, with built assets embedded into the Go binary.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** merged to `dev`

---

## Initiative Structure

```
docs/initiatives/archive/openmuara-web-ui-spa/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    └── 01-framework-decision.md
```

Planning docs live in `docs/initiatives/archive/openmuara-web-ui-spa/` in the root repo.
Product-code commits to the `feat/web-ui-spa` branch. Do not commit directly to `main`.

---

## Current State

The OpenMuara dashboard is served under `/_admin` as a Vite + Preact SPA embedded into the Go binary:

- `web/dashboard/` — SPA source (TypeScript + Preact).
- `internal/ui/dashboard-dist/index.html` — tracked placeholder that lets `go build ./...` work on a fresh clone.
- `internal/ui/dashboard-dist/assets/` — generated SPA assets (ignored by `.gitignore`; produced by `npm run build`).
- `internal/ui/embed.go` — Go embed loader for `dashboard-dist`.
- Escape/pay pages for providers (`fawry-escape.html`, `billplz-pay.html`, etc.) remain server-rendered Go templates.

---

## Goals

1. Decide whether an SPA framework is justified for OpenMuara's dashboard.
2. If yes, choose a framework that preserves the project's values: simple, fast, local-first, single-binary deployment.
3. Migrate the dashboard incrementally without breaking existing routes or escape/pay pages.
4. Maintain zero external runtime dependencies for the Go binary (framework assets must be embedded at build time).
5. Keep the embedded bundle small and memory-efficient; the SPA shares the same process as the Go emulator.

---

## Options

### Option A — Keep vanilla JS, improve structure

- Split `index.html` into ES modules under `internal/ui/js/`.
- Add a lightweight state-management pattern.
- No build step, no Node dependency.

**Pros:** simplest, keeps single-binary ethos, fastest build, no new toolchain.  
**Cons:** still no components, no type safety, limited testability, manual DOM diffing.

### Option B — Lightweight reactivity without a build step

- Use [Alpine.js](https://alpinejs.dev/) or [petite-vue](https://github.com/vuejs/petite-vue).
- Keep HTML server-rendered but add declarative reactivity.

**Pros:** tiny footprint, no build pipeline, easy to embed, gentle learning curve.  
**Cons:** not a true component architecture; still limited tooling for larger features.

### Option C — Vite + Preact

- Use [Vite](https://vitejs.dev/) + [Preact](https://preactjs.com/) (React-compatible, 3 KB runtime).
- Build to static files, embed via Go `//go:embed dist`.

**Pros:** component model, hooks, TypeScript option, dev server, testing ecosystem, tiny runtime.  
**Cons:** adds Node/Vite build step; contributors need npm installed for UI work.

### Option D — Vite + Svelte

- Use [Vite](https://vitejs.dev/) + [Svelte](https://svelte.dev/).
- Build to static files, embed via Go `//go:embed dist`.

**Pros:** compiled away (no virtual DOM), excellent DX, component model, TypeScript support.  
**Cons:** smaller talent pool than React; still requires Node build step.

---

## Decision

**Use Option C — Vite + Preact.**

Preact gives us a React-compatible component model and hooks with only a ~3 KB runtime, which keeps the embedded bundle small. Vite provides a fast dev server, HMR, and a clean build pipeline that outputs static files we can embed with Go's `//go:embed`.

Options A (vanilla), B (Alpine.js), and D (Svelte) remain documented as evaluated alternatives, but the implementation will follow the Vite + Preact path.

Avoid heavy frameworks like Next.js or Nuxt — server-side rendering is unnecessary because OpenMuara already serves the UI from Go.

---

## Recommendations Implemented

The following audit recommendations were applied to make the initiative solid:

1. **TypeScript** — added compile-time safety with zero runtime cost.
2. **Bundle-size budget** — initial JS ≤ 100 KiB gzipped, total `dist/` ≤ 250 KiB; enforced by `web/dashboard/scripts/check-bundle-size.js` in CI and `task ui:check`.
3. **Server-rendered escape/pay pages** — kept as Go templates in `internal/ui/` rather than migrating them to the SPA, preserving payment-flow reliability.
4. **Source layout** — `web/dashboard/` for source, `internal/ui/dashboard-dist/` for built assets embedded by Go; a tracked placeholder `index.html` lets `go build ./...` work on a fresh clone.
5. **UI testing** — Vitest with `@testing-library/preact` for unit tests, including fetch credential-stripping coverage.
6. **CI integration** — `ui-build`, `ui-test`, and `ui-e2e` jobs; Go jobs depend on the built dashboard artifact.
7. **Dev workflow** — `task ui:build`, `task ui:test`, `task ui:check`, `task ui:e2e`, and `task dev` for concurrent Go + Vite dev server.
8. **Accessibility & error states** — ARIA labels, keyboard shortcuts, error boundary, and loading/empty states in all views.
9. **CSP compatibility** — no inline scripts in production; assets served from `/dashboard-assets/` under `default-src 'self'; img-src 'self' data:` so the inline SVG favicon loads in hardened mode.
10. **API client conventions** — CSRF token read from `<meta name="csrf-token">` with cookie fallback; `X-CSRF-Token` header on mutating requests. Fetch URLs are resolved to absolute URLs and stripped of embedded credentials so the SPA works when `/_admin/` is loaded with HTTP Basic Auth credentials in the URL.
11. **No unnecessary webhook probe** — `/_admin/onboarding` exposes `webhooks_enabled`; the failed-webhook alert only queries `/_admin/webhooks` when a dispatcher is configured, avoiding console 404s in minimal configs.
12. **E2E coverage** — Playwright test loads the hardened dashboard with credentials in the URL and asserts no `fetch` credential errors.

## Future Recommendations (Post-v1)

These items were deferred to keep the v1 scope focused and the bundle small:

- Dark mode toggle.
- Offline caching / service worker for the dashboard shell.
- End-to-end visual regression tests.
- Animated page transitions (e.g., React View Transition API once widely supported).

## Non-Goals

- Do not turn OpenMuara into a Node server.
- Do not require internet access at runtime.
- Do not break existing escape/pay pages or `/_admin` URL structure.
- Do not add a separate frontend repository.
- Do not pull in heavy UI libraries, charting frameworks, or state-management libraries that bloat the embedded bundle or increase memory use.
