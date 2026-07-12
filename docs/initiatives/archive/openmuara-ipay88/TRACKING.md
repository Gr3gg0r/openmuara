> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara iPay88 — Execution Tracker

> **Updated:** 2026-07-01 | **Status:** ✅ Completed
>
> **Scope:** Emulate iPay88 payment-request flow for local Southeast Asian payment testing.
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
| P01 | iPay88 Provider | `internal/ipay88/*`, `internal/provider`, `internal/webhook`, `internal/ui`, tests | — | ✅ | — | Implemented payment request, local payment page, response/backend callbacks, requery, SSRF validation, and tests. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run ./internal/ipay88/...` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Not modified (out of scope for this task) | ⏸️ |

---

## Decisions

- D001 ✅ Add iPay88 as a first-class OpenMuara provider named `ipay88`.
