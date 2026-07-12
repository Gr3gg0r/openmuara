> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Web UI SPA — Execution Tracker

> **Updated:** 2026-07-02 | **Status:** 🟡 In Progress
>
> **Scope:** Migrate the `/_admin` dashboard to a Vite + Preact SPA embedded in the Go binary.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/web-ui-spa`
> **Last Agent Action:** Decided on Vite + Preact; assets embedded via Go `//go:embed`.
> **Next Agent Action:** Scaffold the Vite + Preact build pipeline (prompt 02).

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen |

---

## Execution Rules

1. Execute prompts in order unless marked **[PARALLEL SAFE]**.
2. Every prompt MUST end with: tests passing → git commit → update this file to `✅`.
3. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY prompt, update `HANDOFF.md`.
5. Product-code commits happen on `feat/web-ui-spa`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Framework decision | `DECISIONS.md`, `README.md` | — | ✅ | — | Vite + Preact selected; assets embedded into Go binary at build time. |
| 02 | Build pipeline scaffold | `web/dashboard/`, `internal/ui/embed.go`, `Makefile` | 01 | ⬜ | — | Add Node/Vite build that outputs to a directory embedded by Go. |
| 03 | Migrate dashboard shell | `web/dashboard/src/`, `internal/ui/index.html` | 02 | ⬜ | — | Port the existing dashboard layout, tabs, and navigation. |
| 04 | Migrate ledger view | `web/dashboard/src/views/Ledger.jsx` (or `.svelte`) | 03 | ⬜ | — | Port transaction/webhook search, filter, replay, and SSE auto-refresh. |
| 05 | Migrate provider escape/pay pages | `web/dashboard/src/pages/`, provider handlers | 03 | ⬜ | — | Port Fawry escape, Billplz pay, iPay88 pay, ToyyibPay pay, Stripe pages. |
| 06 | Tests and CI | `package.json`, `.github/workflows/`, `scripts/` | 04, 05 | ⬜ | — | Add UI unit tests, build check in CI, and documentation. |
| 07 | Docs update | `docs/providers.md`, `runbooks/local-development.md`, `README.md` | 06 | ⬜ | — | Document the new dashboard build and contribution workflow. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ⬜ |
| UI Build | `cd web/dashboard && npm run build` | Produces embeddable assets | ⬜ |
| Test | `go test ./...` | All pass | ⬜ |
| UI Test | `cd web/dashboard && npm test` | All pass | ⬜ |
| Vet | `go vet ./...` | Clean | ⬜ |
| Lint | `golangci-lint run` | Zero issues | ⬜ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ⬜ |

---

## Decisions

- D001 ⬜ Framework decision pending: vanilla ES modules, Alpine.js, Preact, or Svelte.
- D002 ⬜ Build assets must be embeddable by Go; no runtime Node dependency.
- D003 ⬜ Existing `/_admin` routes and provider escape/pay URLs must remain stable.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-web-ui-spa/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-web-ui-spa/README.md` | Goals, options, recommendation |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
