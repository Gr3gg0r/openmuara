> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Stripe FPX & Card Payments — Execution Tracker

> **Updated:** 2026-06-30 | **Status:** 🟡 Active
>
> **Scope:** Emulate Stripe Checkout Session and PaymentIntents APIs for single-charge FPX/card payments, with local checkout/authentication pages and Stripe-compatible webhook configuration UI.
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
| P01 | Stripe Checkout Sessions | `internal/stripe/checkout.go`, `internal/stripe/checkout_types.go`, `internal/stripe/provider.go`, `internal/stripe/webhook.go`, `internal/ui/*`, tests, smoke test | — | ✅ | `885a14d` | Implement local checkout page for FPX/card; remove custom fpx/card routes |
| P02 | Stripe PaymentIntents | `internal/stripe/payment_intent.go`, `internal/stripe/payment_intent_types.go`, `internal/stripe/provider.go`, `internal/stripe/webhook.go`, `internal/ui/*`, tests, smoke test | P01 | ✅ | 55b3885 | Implement create/retrieve/confirm/cancel for PaymentIntents; support fpx/card; emit correct webhooks |

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

- D001 ✅ Replace the OpenMuara-native `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes with Stripe's real APIs to maintain provider fidelity and SDK compatibility.
- D002 ✅ Scope is single-charge only: `mode=payment` for Checkout, one-time PaymentIntents, inline `price_data`, no product catalog.
- D003 ✅ OpenMuara hosts the checkout/authentication pages locally so the Stripe SDK redirect flow works unchanged.
- D004 ✅ Implement both Checkout Sessions and PaymentIntents because both are common Stripe SDK entry points for FPX/card.
- D005 ✅ Webhook config UI persists to `.muara/config.yml` and updates the running dispatcher.
- D006 ✅ Stripe-compatible error JSON shape and test payment method tokens are part of the contract.
