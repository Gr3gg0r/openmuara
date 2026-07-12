> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA

> **Status:** 🟡 In Progress | **Started:** 2026-07-02
> **Scope:** Migrate the embedded `/_admin` dashboard from vanilla HTML/JS to a Vite + Preact SPA, with built assets embedded into the Go binary.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/web-ui-spa`

---

## Initiative Structure

```
docs/initiatives/openmuara-web-ui-spa/
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

Planning docs live in `docs/initiatives/openmuara-web-ui-spa/` in the root repo.
Product-code commits to the `feat/web-ui-spa` branch. Do not commit directly to `main`.

---

## Current State

The OpenMuara dashboard is served under `/_admin` from `internal/ui/` as embedded vanilla HTML/CSS/JS:

- `internal/ui/index.html` — main dashboard (~33 KB).
- `internal/ui/embed.go` — Go embed loader.
- `internal/ui/handler.go` — HTTP handler.
- Escape/pay pages for providers (`fawry-escape.html`, `billplz-pay.html`, etc.).

This works and keeps the binary self-contained, but as the dashboard grows (ledger, search, replay, provider guides, webhook debugger) the vanilla code is becoming harder to maintain, test, and extend.

---

## Goals

1. Decide whether an SPA framework is justified for OpenMuara's dashboard.
2. If yes, choose a framework that preserves the project's values: simple, fast, local-first, single-binary deployment.
3. Migrate the dashboard incrementally without breaking existing routes or escape/pay pages.
4. Maintain zero external runtime dependencies for the Go binary (framework assets must be embedded at build time).

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

## Non-Goals

- Do not turn OpenMuara into a Node server.
- Do not require internet access at runtime.
- Do not break existing escape/pay pages or `/_admin` URL structure.
- Do not add a separate frontend repository.
