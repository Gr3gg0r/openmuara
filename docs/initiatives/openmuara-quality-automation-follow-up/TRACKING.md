> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Quality Automation Follow-Up — Execution Tracker

> **Updated:** 2026-07-07 | **Status:** ✅ Completed
>
> **Scope:** Mature E1–E12 bug-hunt recommendations from "implemented" to "automated and enforced."
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `dev`
> **Last Agent Action:** Regression-coverage pass completed for `internal/server`, `internal/config`, `internal/cli`, `internal/provider/conform`, `internal/webhook`, and `internal/engine`. Also fixed a real bug in `internal/config/config.go`: `LoadFromBytes` was missing `v.SetConfigType("yaml")`. All Go and frontend gates are green.
> **Next Agent Action:** Monitor first few PRs for gate stability; no further work on this initiative.

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
2. Every gate change MUST include a local reproduction command and a CI path check.
3. Every prompt MUST end with: tests passing → git commit → update this file.
4. If a prompt fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
5. After EVERY prompt, update `HANDOFF.md`.
6. Product-code commits happen on `feat/quality-automation-follow-up`.
7. P04 provider errcode adoption must not change provider behavior; document any exception in `DECISIONS.md`.
8. A gate may start as non-blocking (commentary) and be promoted to required only after it is proven stable for at least three PRs.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | CI visual baseline | `.github/workflows/visual-baseline.yml`, `web/dashboard/e2e/visual-baseline.spec.ts` | — | ✅ | `2fe4a39`, `ec5d37a` | Path-filtered workflow; light + dark theme snapshots; generic `data-visual-mask` helper. |
| 02 | Mutation testing gate | `.github/workflows/mutation.yml`, `scripts/mutation-test.sh` | 01 | ✅ | `10db88b`, `fa1402b`, `ec5d37a` | Gremlins job for `internal/webhook`, `internal/engine`; `internal/fawry` excluded (D010). Path-filtered workflow; non-blocking during rollout (D011). Threshold 70%. |
| 03 | Coverage regression gate | `.github/workflows/coverage-comment.yml`, `scripts/check-coverage-regression.sh` | 01 | ✅ | `cf2d173`, `ce6f729`, `f665251` | Parses cached/non-cached `go test -cover` output; non-blocking during phased rollout. |
| 04 | Provider errcode adoption | `internal/*/...`, `internal/errcode/` | 01 | ✅ | `1e03921`, `22fe92b`, `6498089`, `9f439f8`, `10f2698`, `3b08d92`, `1cd9f38`, `dc800b4`, `dfdc644` | Wrap provider/config/API errors with stable codes; add tests. |
| 05 | Recurring process & KNOWN_ISSUES sync | `.github/workflows/`, `scripts/` | 02–04 | ✅ | `b5860f9` | Monthly bug-hunt prep workflow + KNOWN_ISSUES sync check (CI warning). |
| 06 | Final gates & documentation | `runbooks/quality-gates.md`, `HANDOFF.md` | 01–05 | ✅ | `4def79b`, `fa1402b`, `a3b478e`, `ec5d37a` | Full gate suite run; runbook, changelog, risks, review checklist updated. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Race | `go test -race ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Frontend test | `cd web/dashboard && npm run test:ci` | All pass | ✅ |
| Frontend build | `cd web/dashboard && npm run build` | Passes | ✅ |
| Bundle size | `node web/dashboard/scripts/check-bundle-size.js` | ≤ budgets | ✅ |
| A11y contrast | `cd web/dashboard && node scripts/a11y-contrast-check.js` | Zero violations | ✅ |
| Visual baseline | `cd web/dashboard && npm run test:visual-baseline` | Diff passes (light + dark) | ✅ |
| Mutation | `./scripts/mutation-test.sh 70` | Reports scores (non-blocking during rollout) | ✅ |
| Coverage regression | `./scripts/check-coverage-regression.sh origin/dev 10 1.0` | Reports drops (non-blocking during rollout) | ✅ |
| KNOWN_ISSUES sync | `./scripts/check-known-issues.sh` | No drift | ✅ |

---

## Regression Coverage Pass (2026-07-07)

| Package | What was covered | Commits |
|---------|------------------|---------|
| `internal/server` | Pagination limits, transaction/ledger error branches, provider info helpers, provider health config-load error, CSRF mismatch/missing, scenario validation, security audit nil logger, unknown-provider skip, admin-router pprof, test-only Stripe provider registration | `0bf73dd`, `08178c2`, `719bad9` |
| `internal/config` | `LoadFromBytes` valid/empty/invalid/env override, dual-port admin-port validation, empty host / unsupported persistence, `ValidationError.Error`, `EnvVarName`, `WizardChoice.NextStep`, `RenderWizardConfig` with targets, `sortedStringMapKeys`, `fieldLineMap` missing file. **Bug fix:** added `v.SetConfigType("yaml")` to `LoadFromBytes`. | `2aaf7a6` |
| `internal/cli` | Active-provider fallback, webhook URL helper, dispatcher event parsing, doctor webhook reachability, init defaults/existing-file/no-force, wizard invalid input/EOF, plugin validate empty path, security audit no-issues and hardened-without-admin cases, webhook CLI base-URL load error, non-JSON error response, start command load error, update-check skip, doctor timestamp, copy-file dest error | `821589a`, `77d6586` |
| `internal/provider/conform` | `Capture` Init-error path (with versions), `Capture` success, `Compare` update golden-file creation, `Usage()` | `e077087`, `fabf83c` |
| `internal/webhook` | `eventEnabled`, `NewDispatcherFromBuilder` negative retries, dispatcher event filtering, `MemoryStore.List` negative/overflow offsets, `DeliveryWorker` body-close warning path | `77bf0a2`, `9c7ecde` |
| `internal/engine` | `CanTransition` unknown source status, `MemoryStore.List` negative/overflow offsets, `SQLiteStore.CreateOrGet` marshal-items error, `NewSQLiteStoreFromDB` migrate error with closed db | `e278e71` |

All gates were re-run after the changes and remain green (see table below).

---

## Decisions

- D001 ✅ Gates start as non-blocking commentary and are promoted to required after proven stable.
- D002 ✅ Mutation testing threshold set at 70% initial target; can be raised once baseline is measured.
- D003 ✅ Provider errcode adoption is additive; existing error messages must remain intact unless user approves a breaking change.
- D004 ✅ Visual baseline failures can be resolved by running `npm run test:visual-baseline -- --update-snapshots` and committing only intentional changes.
- D005 ✅ Coverage gate compares changed modules against `main`/`dev` baseline, not global coverage.
- D006 ✅ Initial mutation testing targets `internal/webhook`, `internal/engine`, and `internal/fawry`.
- D007 ✅ Visual baseline CI job uses path filters for `web/dashboard/**` and `internal/ui/**`.
- D008 ✅ Coverage gate ignores packages with fewer than 10 changed lines.
- D009 ✅ Coverage regression gate runs with `continue-on-error` until three stable PRs are observed.
- D010 ✅ Mutation testing excludes `internal/fawry` from the initial curated list.
- D011 ✅ Mutation testing CI job runs with `continue-on-error` during the phased rollout.
- D012 ✅ Visual baseline captures separate light and dark theme snapshots.
- D013 ✅ Dynamic dashboard elements are hidden in visual tests via a shared `[data-visual-mask]` CSS rule.
- D014 ✅ Visual-baseline and mutation jobs live in separate workflows with path filters.

---

## Recommendations for Follow-Up

See `appendices/d-recommendations-roadmap.md` for the complete register of suggestions, enhancements, good-to-haves, and gold-standard practices. The highest-value items not in P01–P06 are:

- Visual baseline per theme (light/dark) and per viewport size.
- Mutation testing expansion to additional packages once the initial set is green.
- Coverage-regression gate override workflow documented in `DECISIONS.md`.
- Automated changelog generation from merged PR labels.
- Provider contract test auto-update guard (fail if golden files change unexpectedly).
- Nightly full quality matrix to catch environmental/tooling drift.
