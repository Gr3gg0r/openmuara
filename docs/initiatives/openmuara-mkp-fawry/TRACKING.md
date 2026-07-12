> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara MKP Fawry Integration — Execution Tracker

> **Updated:** 2026-07-08 | **Status:** ⏸️ Suspended
>
> **Scope:** Close Fawry emulation gaps for MKP v2.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/mkp-fawry`
> **Last Agent Action:** Suspended implementation; uncommitted work stashed as `WIP: suspend MKP Fawry implementation`.
> **Next Agent Action:** Resume on `feat/mkp-fawry` by popping the stash and completing step 05 (docs and CI).

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
5. Product-code commits happen on `feat/mkp-fawry`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Fawry state extensions | `internal/engine/transaction.go`, `internal/fawry/escape.go`, `internal/fawry/v2/webhook.go`, tests | — | ✅ | TBD | Added `canceled` and `expired` ledger states; support `CANCELED` / `EXPIRED` order status in escape and webhook. |
| 02 | Response delay config | `internal/config/`, `internal/fawry/plugin.go`, `internal/fawry/escape.go`, tests | — | ✅ | TBD | Added `fawry.response_delay_ms` and delay outgoing webhook dispatch. |
| 03 | Billing type and journey | `internal/fawry/charge.go`, `internal/fawry/v2/webhook.go`, `internal/engine/transaction.go`, tests | 01 | ✅ | TBD | Accept `billing_type` on charge; shape subscription vs prepaid webhook payload. |
| 04 | Escape page and webhook shape | `internal/ui/`, `internal/fawry/escape.go`, `web/fawry-escape.html`, tests | 01, 02, 03 | ✅ | TBD | UI supports all statuses, billing type, and delay preview. |
| 05 | Docs and CI | `docs/providers/fawry.md`, `docs/mkp-billing-requirements.md`, `runbooks/`, smoke test | 01–04 | ⏸️ | — | Update provider docs, MKP requirements status, and add coverage. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ |
| Test | `go test ./...` | All pass | ✅ |
| Vet | `go vet ./...` | Clean | ✅ |
| Lint | `golangci-lint run` | Zero issues | ✅ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ⏸️ |

---

## Decisions

- D001 ✅ Existing Fawry behavior stays backward-compatible; new fields are optional.
- D002 ✅ `response_delay_ms` applies to outgoing webhook dispatch only, not the escape redirect.
- D003 ✅ `GET /fawry/payment-status` added under VAL01 to let clients verify payment status by `merchantRefNum`; signature required.
- D004 ✅ MKP delegates Fawry simulation to OpenMuara via `POST /fawry/charge` and `POST /fawry/simulate` when `OPENMUARA_URL` is configured.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-mkp-fawry/TRACKING.md` | Initiative execution tracker |
| Source requirements | `docs/mkp-billing-requirements.md` | MKP v2 billing requirements |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
