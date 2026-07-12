> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1 — Execution Tracker

> **Updated:** 2026-06-29 | **Total Prompts:** 18 | **Status:** ✅ COMPLETE
>
> **Scope:** Rebrand `muara` to `OpenMuara` and ship v1 core features on `dev`.
> **Repo:** `<repo-root>`
> **Product Branch:** `dev`
> **Last Agent Action:** Committed and pushed v1 close-out to `dev` (`5ef1f0d`).
> **Next Agent Action:** Create release tag or merge `dev` → `main` per release workflow.

---

## ⚠️ Repo & Branch Boundary Rules (READ BEFORE EVERY PROMPT)

- `docs/projects/openmuara-v1/` is for **planning docs only**.
- All implementation code commits to the **`dev`** branch in `<repo-root>`.
- **Never commit directly to `main` or `master`.**
- **Do not mix planning-doc commits with product-code commits.**
- Pre-flight check: `git branch --show-current` must return `dev` before any product-code commit.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen for v2 |
| 🔀 | Parallel Safe |

---

## Execution Rules

1. Execute steps in order unless marked **[PARALLEL SAFE]**.
2. Every step MUST end with: tests passing → git commit → update this file to `✅`.
3. If a step fails a quality gate, STOP. Do not proceed. Log the blocker in `RISKS.md`.
4. After EVERY step, update `HANDOFF.md`.
5. Product-code commits happen on `dev`.

---

## Phase 0 — Rebrand

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 01 | Rename Go module | `go.mod`, all import paths | — | ✅ | `32c53e0` | — | See `prompts/01-rename-go-module.md` |
| 02 | Rename binary and CLI | `cmd/openmuara/`, `Taskfile.yml`, CLI code | 01 | ✅ | `bf4edd4` | D005 | See `prompts/02-rename-binary-and-cli.md` |
| 03 | Rename config/data paths | `internal/config/`, `.muara/` refs | 02 | ✅ | `5a09e40` | — | See `prompts/03-rename-config-paths.md` |
| 04 | Update docs and tooling | `README.md`, `AGENTS.md`, scripts, CI | 03 | ✅ | `23241b6` | — | See `prompts/04-update-docs-and-tooling.md` |
| 05 | Rebrand quality gates | all | 04 | ✅ | `766367d` | — | Run full `task check` + `task smoke`; final `muara` sweep in tests |

## Phase 1 — Core Runtime

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 06 | Add SQLite persistence | `internal/engine/`, schema, migrations | 05 | ✅ | `27b13d9` | D006 | SQLite default; MemoryStore retained |
| 07 | Universal payment API | `internal/server/router.go`, `internal/api/pay.go` | 06 | ✅ | `1b2f9e6` | — | `POST /v1/pay`, `GET /v1/pay/{ref}`, `POST /v1/refund/{ref}` |
| 08 | Scenario commands | `internal/cli/scenario.go` | 07 | ✅ | `594693c` | — | `muara scenario success/fail/timeout` |

## Phase 2 — Providers

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 09 | Stripe provider adapter | `internal/stripe/` or `plugins/stripe/` | 05 | ✅ | `64bbbe2` | — | Added fail/cancel simulation endpoints |
| 10 | SenangPay provider adapter | `internal/senangpay/` or `plugins/senangpay/` | 05 | ✅ | `6af3bd4` | — | Charge + callback + webhook with MD5 signature |
| 11 | Provider config loader | `internal/config/`, `internal/provider/` | 07, 09 | ✅ | `53db749` | — | Declarative provider activation |

## Phase 3 — Webhook Relay & UI

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 12 | Webhook relay core | `internal/webhook/relay.go` | 06 | ✅ | `ce3cf63` | — | Multi-destination forwarding; per-provider targets |
| 13 | Webhook replay API | `internal/server/webhook_admin.go` | 12 | ✅ | `d074bd6` | — | Filters, replay-all, delete |
| 14 | Basic web UI | `internal/ui/`, `web/` | 07, 13 | ✅ | `dbb490a` | — | Dashboard + admin JSON endpoints |

## Phase 4 — Packaging & Docs

| Step | Title | Target Files | Depends On | Status | Commit | Decisions | Notes |
|------|-------|--------------|------------|--------|--------|-----------|-------|
| 15 | Docker image | `Dockerfile`, `docker-compose.yml` | 05 | ✅ | `981cf86` | — | Local image and compose support; CI badge added in README via `prompts/15-docker-ci.md`. |
| 16 | OpenAPI spec | `docs/openapi.yaml` | 07, 09 | ✅ | `6e99349` | — | API contract + live endpoint |
| 17 | Test SDK | `internal/testsdk/` | 07 | ✅ | `3eb945c` | — | Provider-agnostic client |
| 18 | Migration guide | `docs/migration/openmuara-to-openmuara.md`, `scripts/migrate-openmuara.sh`, `internal/cli/migrate.go` | all | ✅ | `5ef1f0d` | — | `muara migrate` CLI and shell script; does not delete old workspace. |

---

## Quality Gate Results

### Stack: Go (`<repo-root>`)

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Format | `task fmt` | Clean | ✅ |
| Vet | `task vet` | Clean | ✅ |
| Lint | `task lint` | Zero issues | ✅ |
| Test | `task test` | All pass | ✅ |
| Race | `task test -race` | All pass | ✅ |
| Coverage | `task coverage` | 72.8% (≥ 50%) | ✅ |
| Build | `task build` | Compiles | ✅ |
| Smoke | `task smoke` | Passes | ✅ |

---

## Findings Inventory

| Step | Finding File | Status | Summary |
|------|-------------|--------|---------|
| | | ⬜ | |

---

## Runbooks Inventory

| Runbook | File | Status | Purpose |
|---------|------|--------|---------|
| Local development | `runbooks/local-development.md` | ✅ | Build, run, and configure locally |
| Quality gates | `runbooks/quality-gates.md` | ✅ | Test, lint, coverage, and smoke workflow |
| On-call | `runbooks/on-call.md` | ✅ | Alerts and first response |
| Debugging | `runbooks/debugging.md` | ✅ | Inspect state, replay webhooks, fix issues |
| Operations | `docs/operations.md` | ✅ | Deployment, metrics, logs, backup |

---

## Completion Checklist

- [x] All prompts in this tracker are ✅.
- [x] All quality gates show ✅ Pass.
- [x] `DECISIONS.md` is up to date.
- [x] `RISKS.md` shows all risks as mitigated or accepted.
- [x] `KNOWN_ISSUES.md` is accurate.
- [x] `HANDOFF.md` reflects final state.
- [x] `dev` branch is green.
- [x] `dev` branch is pushed (`5ef1f0d`).
