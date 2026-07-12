> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara ToyyibPay — Execution Tracker

> **Updated:** 2026-07-01 | **Status:** ✅ Completed
>
> **Scope:** Emulate ToyyibPay API for local Malaysian payment testing.
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

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

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| P01 | ToyyibPay Provider | `internal/toyyibpay/*`, `internal/provider`, `internal/webhook`, `internal/ui`, tests, smoke test | — | ✅ | e6d02f4 | Implemented categories, bills, local payment page, return/callback, and tests. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ |

---

## Decisions

- D001 ✅ Add ToyyibPay as a first-class OpenMuara provider named `toyyibpay`.
