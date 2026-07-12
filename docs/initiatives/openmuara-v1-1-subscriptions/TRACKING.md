> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1.1 — Subscriptions — Execution Tracker

> **Updated:** 2026-07-09 | **Status:** ⏸️ Suspended
>
> **Scope:** Add subscription emulation to OpenMuara, starting with Stripe Billing, then extending to Malaysian gateways that support recurring payments.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `feat/v1-1-subscriptions` (no work started)
> **Last Agent Action:** User suspended v1.1 subscriptions initiative on 2026-07-09.
> **Next Agent Action:** Resume only when user explicitly asks to start subscription emulation.

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
5. Product-code commits happen on `feat/v1-1-subscriptions`.

---

## Prompt Inventory

| Step | Title | Target Files | Depends On | Status | Commit | Notes |
|------|-------|--------------|------------|--------|--------|-------|
| 01 | Stripe subscriptions | `internal/stripe/subscription.go`, `internal/engine/subscription.go`, `internal/store/migrations/` | — | ⬜ | — | Emulate Stripe Billing products, prices, customers, subscriptions, invoices, and webhooks. |
| 02 | SenangPay subscriptions | `internal/senangpay/subscription.go` | 01 | ⬜ | — | Map SenangPay recurring/subscription plan flow onto the subscription engine. |

---

## Quality Gate Results

| Gate | Command | Target | Status |
|------|---------|--------|--------|
| Build | `go build ./...` | Compiles | ⬜ |
| Test | `go test ./...` | All pass | ⬜ |
| Vet | `go vet ./...` | Clean | ⬜ |
| Lint | `golangci-lint run` | Zero issues | ⬜ |
| Smoke | `./scripts/smoke-test.sh` | Passes | ⬜ |

---

## Decisions

- D001 ✅ v1.1 focuses on subscriptions, starting with Stripe Billing.
- D002 ✅ Existing Malaysian gateways with recurring support (SenangPay, iPay88, Billplz, ToyyibPay, Fiuu, 2C2P, eGHL, Curlec, HitPay, Xendit, Airwallex) are candidates after Stripe.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This tracker | `docs/initiatives/openmuara-v1-1-subscriptions/TRACKING.md` | Initiative execution tracker |
| Initiative README | `docs/initiatives/openmuara-v1-1-subscriptions/README.md` | Goals, provider landscape, target endpoints |
| v1 master backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | v1 priority view |
| Root tracker | `TRACKING.md` | Cross-prompt and initiative status |
