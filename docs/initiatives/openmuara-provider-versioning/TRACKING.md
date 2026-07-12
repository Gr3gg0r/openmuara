> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider API Versioning — Execution Tracker

> **Updated:** 2026-07-01 | **Status:** ✅ Completed
>
> **Scope:** Introduce a versioning convention for provider packages and migrate Fawry as the reference implementation.
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/provider-versioning`

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
5. Product-code commits happen on `feat/provider-versioning`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| P01 | Versioned provider layout and Fawry reference migration | `internal/fawry/`, `internal/provider`, `plugins/fawry/`, tests, smoke test, config defaults | — | ✅ | TBD | Restructure Fawry into `v1/` and `v2/`; keep unversioned routes as aliases to the configured default version. |

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

- D001 ✅ Add a provider-versioning convention that supports both single-version and multi-version providers without forcing every provider to adopt it immediately.
