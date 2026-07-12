> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Testing Gold Standard — Execution Tracker

> **Updated:** 2026-06-29 | **Status:** 🟡 Active
>
> **Scope:** Establish OSS-grade testing practices and backfill tests.
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`
> **Selected Approach:** Option A — 80% coverage, refactor-first.
> **Current Coverage:** 88.6%

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
5. Product-code commits happen on `dev`.
6. P21 touches provider/plugin/cli/server registries — this is cross-module refactoring. Confirm with the user before starting.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| P20 | Charter & runbook | `runbooks/testing.md`, initiative docs | — | ✅ | `b5edbe9` | Define conventions and create tracker |
| P21 | Testability refactor | `internal/provider/registry.go`, `internal/plugin/registry.go`, `internal/cli/start.go`, `internal/server/server.go` | P20 | ✅ | `9c1e336` | Confirm before starting |
| P22 | Shared test utilities | `internal/testutil/*.go` | P21 | ✅ | `6048821` | Fakes, temp workspace, SQLite helpers |
| P23 | Unit-test backfill | `internal/cli/*_test.go`, `internal/server/server.go`, `internal/audit/logger_test.go`, `internal/config/config_test.go`, provider tests | P22 | ✅ | `46e6561` | Coverage 81.0% → 88.6% after backfill |
| P24 | Provider contract tests | `internal/{fawry,stripe,senangpay}/contract_test.go`, `testdata/` | P22 | ✅ | — | Golden files for charge/create-session endpoints |
| P25 | Integration & E2E hardening | `scripts/smoke-test.sh` | P23, P24 | ✅ | — | Random ports, isolated workspaces, race/shuffle pass, parallel smoke |
| P26 | Advanced testing | `*_fuzz_test.go` | P25 | ✅ | — | Fuzz signature roundtrips and state-machine transitions |
| P27 | CI & quality gates | `.github/workflows/ci.yml`, `Taskfile.yml`, `runbooks/quality-gates.md` | P25 | ✅ | — | Split lint/unit/smoke jobs; 80% coverage gate |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Race | `go test -race ./...` | All pass | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |
| Coverage | `task coverage` | 88.6% / target 80% | ✅ |

---

## Completion Checklist

- [x] All prompts in this tracker are ✅.
- [x] All quality gates show ✅ Pass.
- [x] Coverage ≥ 80%.
- [ ] `RISKS.md` shows all risks as mitigated or accepted.
- [ ] `KNOWN_ISSUES.md` is accurate.
- [x] `HANDOFF.md` reflects final state.
- [x] `dev` branch is green and pushed.
