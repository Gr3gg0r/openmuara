> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Quality Automation Follow-Up — Review Checklist

> **Updated:** 2026-07-06

Final review completed after merging `feat/quality-automation-follow-up` into `dev`.

## Process

- [x] `TRACKING.md` shows all prompts `✅` with commit hashes.
- [x] `HANDOFF.md` has the final state and any open questions.
- [x] `DECISIONS.md` has final statuses for all decisions (`Decided` / `Approved` / `Deferred`).
- [x] `RISKS.md` has closed or accepted statuses for resolved risks.
- [x] Every gate is documented in `runbooks/quality-gates.md` or equivalent with its local command.

## Code & CI

- [x] Visual baseline workflow runs only on dashboard/UI changes and can be updated via `--update-snapshots`.
- [x] Visual baseline captures separate light and dark theme snapshots (R19).
- [x] Dynamic elements are masked generically via `[data-visual-mask]` (R21).
- [x] Mutation testing workflow runs only on Go changes and targets `internal/webhook`, `internal/engine`.
- [x] Mutation threshold (70%) is documented and achievable.
- [x] Coverage-regression gate compares changed modules and reports drops; non-blocking during rollout.
- [x] Provider errcode adoption does not alter existing error message text.
- [x] KNOWN_ISSUES sync script runs as a warning before promotion.

## Quality Gates

- [x] `go build ./...` ✅
- [x] `go test ./...` ✅
- [x] `go test -race ./...` ✅
- [x] `go vet ./...` ✅
- [x] `golangci-lint run` ✅
- [x] `cd web/dashboard && npm run test:ci` ✅
- [x] `cd web/dashboard && npm run build` ✅
- [x] `node web/dashboard/scripts/check-bundle-size.js` ✅
- [x] `cd web/dashboard && node scripts/a11y-contrast-check.js` ✅
- [x] `cd web/dashboard && npm run test:visual-baseline` ✅ (light + dark)
- [x] `./scripts/check-known-issues.sh` ✅
- [x] `./scripts/check-coverage-regression.sh origin/dev 10 1.0` reports drops correctly.
- [x] `./scripts/mutation-test.sh 70` reports scores.

## Security & Philosophy

- [x] No external SaaS dependencies were added.
- [x] No secrets or `.muara/config.yml` files are committed.
- [x] All new gates are reproducible locally.
- [x] No provider behavior changed unless explicitly signed off.
- [x] No secrets or PII appear in visual baseline screenshots.
- [x] New CI workflows do not expose sensitive values in logs.
- [x] `gosec` and `gitleaks` still pass (no new findings).

---

**Reviewer:** user (approved via prompt on 2026-07-06)  
**Merged:** `feat/quality-automation-follow-up` → `dev` at `ec5d37a`  
**Feature branch:** deleted  
**Working tree:** clean on `dev`
