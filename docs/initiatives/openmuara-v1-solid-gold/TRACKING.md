> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid Gold — Execution Tracker

> **Updated:** 2026-07-01 | **Status:** ✅ Completed
>
> **Scope:** v1 hygiene, testing, debuggability, and usability polish.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/v1-solid-gold`
> **Last Agent Action:** P05 best practices and tooling implemented.
> **Next Agent Action:** Land `feat/v1-solid-gold` via PR to `dev`.

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
5. Product-code commits happen on `feat/v1-solid-gold`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Tooling hygiene | `muara.yml.example`, `scripts/smoke-test.sh`, `.github/workflows/ci.yml`, `internal/config/config.go` | — | ✅ | 3522f57 | Sync example config, fix shellcheck warnings, run full `task quality` in CI. |
| 02 | Coverage backfill | `internal/ui/`, `internal/fawry/v2/`, `internal/cli/`, `internal/fawry/v1/`, `internal/ipay88/`, `internal/toyyibpay/`, `internal/billplz/`, `cmd/muara/`, `internal/testutil/`, tests | — | ✅ | 5f61d5d, d74f518, abdd2a5, 00cf1cc, fb86929 | Bring every package to ≥80% coverage. Total coverage 89.0%. |
| 03 | Debuggability | `internal/webhook/dispatcher.go`, `internal/server/`, `internal/cli/`, `runbooks/debugging.md` | 01 | ✅ | 259a821 | Trace-ID propagation, CLI inspect commands, optional pprof. **P0 integration — user approved.** |
| 04 | Dashboard usability | `internal/ui/index.html` | 01 | ✅ | c504cb9 | Failed-webhook alert, copy-curl buttons, responsive layout. Verified with browser snapshots. |
| 05 | Best practices and tooling | `.golangci.yml`, `.pre-commit-config.yaml`, `.github/dependabot.yml`, `Taskfile.yml`, `runbooks/quality-gates.md` | 01, 02 | ✅ | e161460 | Stronger linters (gosec, staticcheck, ineffassign, unparam, errcheck), pre-commit hooks, govulncheck CI, Dependabot, `-trimpath` releases, dead-code removal, dispatcher race fix. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Race | `go test -race ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |
| Quality | `task quality` | All checks pass | ✅ |

Advisory size/line-length warnings remain from `scripts/check-sizes.sh`; they do
not fail the gate and are tracked as accepted technical debt (D005).

---

## Decisions

- D001 ✅ All changes are additive/hygiene; no breaking config or API changes.
- D002 ✅ P03 (webhook trace IDs / pprof) needs explicit user approval before implementation.
- D003 ✅ Advisory size/line-length warnings are accepted debt for v1; do not block prompts.
- D004 ✅ Test-only backfills are committed per package group to keep diffs reviewable.
- D005 ✅ Linter findings are fixed when trivial; suppressed with `#nosec` + justification when they are test-only or false positives. Dead code removed only after confirming no callers.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-v1-solid-gold/TRACKING.md` | Initiative execution tracker |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
| Testing gold standard | `docs/initiatives/openmuara-testing-gold-standard/TRACKING.md` | Prior testing initiative |
