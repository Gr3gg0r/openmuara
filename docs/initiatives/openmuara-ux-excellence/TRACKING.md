> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence — Execution Tracker

> **Updated:** 2026-07-09 | **Status:** ✅ Complete / Merged to `dev`
>
> **Scope:** Make OpenMuara the most approachable local payment emulator for developers, AI agents, testers, and contributors (ledger-style web UI inspired by Mailpit).
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev` (feature branch merged and removed)

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
5. Product-code commits happen on `feat/ux-excellence`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| P01 | First-run config wizard | `internal/cli/init.go`, `internal/config/`, `muara.yml.example`, tests | — | ✅ | a790409 | Interactive `muara init` that asks target provider and generates tailored config. |
| P02 | Dashboard onboarding checklist | `internal/server/admin_api.go`, `internal/ui/index.html`, tests | P01 | ✅ | db9c3e5 | Show getting-started steps on `/_admin` and mark them complete as the user progresses. |
| P03 | Actionable config validation | `internal/config/`, `internal/cli/doctor.go`, `internal/cli/start.go`, tests | — | ✅ | 9bea8b2 | Validate config at load time and report errors with file path, line number, and fix hint. |
| P04 | Provider selection guide | `internal/ui/index.html`, `docs/providers.md`, `docs/providers/`, tests | P02 | ✅ | 4f5ef2f | Help users pick and configure the right provider for their real gateway; include contributor checklist. |
| P05 | Webhook debugger | `internal/server/admin_api.go`, `internal/ui/index.html`, `internal/webhook/`, tests | — | ✅ | 4415629 | Expose payload, signature status, retry timeline, and replay in the dashboard; store signature verification result. |
| P06 | Transaction search and replay | `internal/server/admin_api.go`, `internal/ui/index.html`, `internal/engine/`, tests | — | ✅ | a783605 | Searchable transactions table with detail endpoint and replay actions. |
| P07 | CLI help and structured output | `internal/cli/*.go`, `docs/cli-schemas/`, `runbooks/local-development.md` | — | ✅ | — | Add runnable examples to every command's `--help`; add `--json` / `--quiet` flags and documented schemas for AI agents. |
| P08 | Ledger-style payment view | `internal/server/admin_api.go`, `internal/ui/index.html`, `internal/engine/`, `internal/webhook/`, tests | P02, P05, P06 (transaction detail endpoint) | ✅ | 7423908 | Unified, visibility-aware auto-refreshing ledger view of transactions and webhooks with search, filter, payload inspection, and replay. Endpoint `GET /_admin/ledger`; UI default tab with keyboard shortcuts (`?`, `/`, `1`/`2`/`3`). |
| P09 | Quick-start documentation | `docs/quickstart.md`, `README.md`, `runbooks/local-development.md` | P01, P04, P07, P08 | ✅ | 155702f | Single page with Developer, AI Agent, Tester, and Contributor paths from zero to first charge. |

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

- D001 ⬜ UX improvements must be additive and preserve existing CLI/config behavior.
