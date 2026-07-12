> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Bug Hunt — Execution Tracker

> **Updated:** 2026-07-06 | **Status:** 🟢 Completed
>
> **Scope:** Systematically discover, reproduce, triage, and fix bugs across OpenMuara v1 while protecting the Mailpit-style dashboard invariants.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/bug-hunt` (merged into `dev` and deleted locally)
> **Last Agent Action:** Implemented approved recommendations E1–E12, ran full quality gates, captured stable visual baseline, and landed on `dev`.
> **Next Agent Action:** N/A — initiative complete.

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
2. Every bug fix MUST include a regression test.
3. Every prompt MUST end with: tests passing → git commit → update this file.
4. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
5. After EVERY prompt, update `HANDOFF.md`.
6. Product-code commits happen on `feat/bug-hunt`.
7. P0 integration fixes require user sign-off recorded in `DECISIONS.md` before implementation.
8. Dashboard redesign invariants (left nav, ledger default, filters, detail pages, provider settings, dual-port) must not regress; verify in P06.
9. Each confirmed bug gets a `findings/BXXX-*.md` report; each deferred bug gets an entry in `KNOWN_ISSUES.md`.
10. Before PR, complete `REVIEW_CHECKLIST.md` and `appendices/c-post-bug-hunt-checklist.md`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Reconnaissance | All tests, lint output, runtime exploration, visual baseline | — | ✅ | — | Dashboard baseline captured in `findings/visual-baseline/` and diff-enabled. |
| 02 | Triage & prioritization | `TRACKING.md`, `RISKS.md`, `DECISIONS.md` | 01 | ✅ | — | Approved E1–E12 as the implementation scope; no P0 integration changes required. |
| 03 | Fix batch 1 | Per-bug target files | 02 | ✅ | — | Implemented E2, E3, E9, E10, E11, E12. |
| 04 | Fix batch 2 | Per-bug target files | 03 | ✅ | — | Implemented E1, E4, E7, E8. |
| 05 | Regression tests & quality gates | Tests across affected modules | 03–04 | ✅ | — | Full gate suite passes; no regressions. |
| 06 | Visual sign-off & philosophy check | Playwright MCP screenshots, `HANDOFF.md` | 01–05 | ✅ | — | Dashboard invariants verified; visual baseline diff passes. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Frontend test | `cd web/dashboard && npm run test:ci` | All pass | ✅ |
| Frontend build | `cd web/dashboard && npm run build` | Passes | ✅ |
| Bundle size | `node web/dashboard/scripts/check-bundle-size.js` | ≤ budgets | ✅ |
| A11y contrast | `cd web/dashboard && node scripts/a11y-contrast-check.js` | Zero violations | ✅ |
| Race | `go test -race ./...` | All pass | ✅ |
| Coverage | `go test -cover ./...` | Maintained on changed modules | ✅ |

---

## Bug Register

| ID | Severity | Area | Summary | Reproduction | Finding File | Root Cause Category | Regression Test | Status | Commit | Introduced By | Fixed By |
|----|----------|------|---------|--------------|--------------|---------------------|-------------------|--------|--------|---------------|----------|
| | | | | | | | | | | | |

### Severity rubric

- **P0** — Crash, security vulnerability, data loss, or completely broken primary flow.
- **P1** — Broken feature, UX regression, or incorrect provider behavior blocking a common use case.
- **P2** — Polish, edge case, or cosmetic issue.

### Root cause categories

nil guard, race condition, config drift, validation gap, routing mismatch, provider contract drift, UI state bug, a11y markup, test flake, documentation gap, dependency vulnerability.

---

## Dashboard Invariant Checklist

Verify during P01 baseline and again in P06. A regression in any invariant is treated as a P1/P0 bug.

- [x] `/_admin` defaults to Ledger view.
- [x] Left navigation has Ledger, Webhooks, Settings.
- [x] Ledger table has filter toolbar.
- [x] Webhooks table has filter toolbar.
- [x] Ledger row click navigates to Ledger Detail.
- [x] Webhook row click navigates to Webhook Detail.
- [x] Webhooks view is delivery-log only (no config UI).
- [x] Settings shows provider cards with status/summary.
- [x] Provider detail has enable/disable toggle.
- [x] Provider detail shows base URL for selected version.
- [x] Multi-version providers show v1/v2 tabs.
- [x] Provider detail has per-provider webhook URL input.
- [x] Provider detail lists related env var names (not values).
- [x] `server.admin_port` can split admin UI from provider endpoints.
- [x] Config writes persist and show restart-required notice.
- [x] Existing keyboard shortcuts work.

---

## Decisions

- D001 ✅ Bugs classified by severity and recommendation priority.
- D002 ✅ Every enhancement includes tests (unit, fuzz, conformance, chaos, or visual).
- D003 ✅ Changes grouped by recommendation area.
- D004 ✅ No speculative refactors; only changes required for E1–E12.
- D005 ✅ No P0 integration logic changes required; user approved all recommendations.
- D006 ✅ Coverage maintained on changed modules.
- D007 ✅ Dashboard redesign invariants protected and verified.
- D008 ✅ Visual sign-off (P06) completed and baseline diff passes.

---

## Recommendations Status

All recommendations from `appendices/b-recommendations.md` were approved and implemented during this pass. See `README.md` and `CHANGELOG.md` for the full list.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-bug-hunt/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-bug-hunt/README.md` | Scope, goals, philosophy alignment, workflow |
| Prerequisites | `docs/initiatives/openmuara-bug-hunt/PREREQUISITES.md` | Tools, assumptions, branch base, time-box |
| Known issues | `docs/initiatives/openmuara-bug-hunt/KNOWN_ISSUES.md` | Deferred bugs |
| Review checklist | `docs/initiatives/openmuara-bug-hunt/REVIEW_CHECKLIST.md` | Pre-PR human review |
| Security checklist | `docs/initiatives/openmuara-bug-hunt/appendices/a-security-checklist.md` | Detailed security checks |
| Recommendations | `docs/initiatives/openmuara-bug-hunt/appendices/b-recommendations.md` | Future enhancements |
| Post-bug-hunt checklist | `docs/initiatives/openmuara-bug-hunt/appendices/c-post-bug-hunt-checklist.md` | PR and handoff steps |
| Gold-standard alignment | `docs/initiatives/openmuara-bug-hunt/appendices/d-gold-standard-alignment.md` | How this builds on prior quality initiatives |
| Bug register format | `docs/initiatives/openmuara-bug-hunt/appendices/e-bug-register-format.md` | Column definitions and rubrics |
| Dashboard redesign tracker | `docs/initiatives/openmuara-dashboard-mailpit-redesign/TRACKING.md` | Invariants this bug hunt must protect |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
