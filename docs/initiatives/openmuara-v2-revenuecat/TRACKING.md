> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v2 — RevenueCat Emulation — Execution Tracker

> **Updated:** 2026-07-09 | **Status:** ⏸️ Suspended
>
> **Scope:** Emulate RevenueCat subscriber status, offerings, receipt submission, and entitlement webhooks for v2.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/v2-revenuecat` (no work started)
> **Last Agent Action:** User suspended v2 RevenueCat initiative on 2026-07-09.
> **Next Agent Action:** Resume only when user explicitly asks to start RevenueCat emulation.

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
5. Product-code commits happen on `feat/v2-revenuecat`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | RevenueCat emulation | `internal/revenuecat/`, `internal/store/migrations/` | — | ⬜ | — | See `prompts/01-revenuecat-emulation.md` |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ⬜ |
| Test | `go test ./internal/revenuecat/...` | All pass | ⬜ |
| Vet | `go vet ./...` | Clean | ⬜ |
| Lint | `golangci-lint run` | Zero issues | ⬜ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ⬜ |

---

## Decisions

- D001 ✅ RevenueCat deferred from v1 to v2 because it is a subscription/entitlement emulator, not a single-charge payment emulator.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-v2-revenuecat/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-v2-revenuecat/README.md` | Goals, scope, non-goals |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | v1 priority view |
| Root tracker | `TRACKING.md` | Cross-prompt and initiative status |
