> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This tracker is subordinate to it.**

# OpenMuara Stripe FPX — Execution Tracker

> **Updated:** 2026-07-03 | **Status:** ❄️ Archived / Superseded
>
> **Scope:** Emulate Stripe FPX (Malaysian online bank transfer) payment flows via a charge + escape pattern.
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`
> **Superseded by:** [`docs/initiatives/openmuara-stripe-checkout-sessions/`](../openmuara-stripe-checkout-sessions/)

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
| P01 | Stripe FPX and card charge and escape | `internal/stripe/fpx.go`, `internal/stripe/fpx_types.go`, `internal/stripe/card.go`, `internal/stripe/card_types.go`, `internal/stripe/provider.go`, tests | — | ✅ / ❄️ Superseded | `69f0f6b`, `81c0cc6` | Fawry-style charge + escape for Stripe FPX and card payments; emits `checkout.session.completed` / `payment_intent.canceled`; `internal/stripe` coverage 94.7%. **Removed in `885a14d` by successor initiative.** |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ✅ (at time of completion) |
| Test | `go test ./...` | All pass | ✅ (at time of completion) |
| Vet | `go vet ./...` | Clean | ✅ (at time of completion) |
| Lint | `golangci-lint run` | Zero issues | ✅ (at time of completion) |
| Smoke | `./scripts/smoke-test.sh` | Passes | ✅ (at time of completion) |

---

## Decisions

- D001 ✅ Scope FPX as a dedicated charge + escape flow within the Stripe provider, modeled after Fawry, not as a Stripe Checkout Session payment method.
- D002 ❄️ Supersede the custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes with real Stripe Checkout Sessions and PaymentIntents APIs for SDK compatibility. Recorded in successor initiative.

---

## Archival Notes

- Implementation commits `9ed0ac4`, `69f0f6b`, and `81c0cc6` added the custom routes.
- Commit `885a14d` removed the custom routes and replaced them with Stripe Checkout Sessions.
- This tracker is frozen; active development continues in `openmuara-stripe-checkout-sessions/`.
