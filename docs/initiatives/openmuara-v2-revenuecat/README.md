> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v2 — RevenueCat Emulation

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Emulate RevenueCat subscriber status, offerings, receipt submission, and entitlement webhooks for v2.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/v2-revenuecat`

---

## Initiative Structure

```
docs/initiatives/openmuara-v2-revenuecat/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    └── 01-revenuecat-emulation.md
```

Planning docs live in `docs/initiatives/openmuara-v2-revenuecat/` in the root repo.
Product-code commits to the `feat/v2-revenuecat` branch. Do not commit directly to `main`.

---

## Why v2?

RevenueCat is a subscription and entitlement platform, not a single-charge payment gateway. v1's explicit philosophy is to focus the emulator on **single charge items** — one-time payments through providers like Stripe Checkout, Fawry, SenangPay, iPay88, Billplz, and ToyyibPay.

Subscriptions, mobile receipts, and entitlement lifecycle are a different emulation surface:

- They require persistent subscriber state separate from the transaction ledger.
- They introduce time-based concepts (trials, renewals, expirations, grace periods).
- They need mobile receipt validation (App Store / Google Play) as a prerequisite.
- Their webhook event vocabulary is different from one-time payment webhooks.

These capabilities belong in v2, where OpenMuara can expand from "payment emulator" to "subscription and purchase emulator."

---

## Goals

1. Emulate RevenueCat REST endpoints faithfully enough for local integration testing.
2. Persist subscriber, offering, and entitlement state in SQLite.
3. Support receipt submission and entitlement updates.
4. Dispatch RevenueCat-style subscriber webhook events.
5. Align with v2's broader subscription/purchase emulation architecture.

---

## Target Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/subscribers/{app_user_id}` | Get or create subscriber status |
| `GET`  | `/v1/subscribers/{app_user_id}/offerings` | List available offerings |
| `POST` | `/v1/receipts` | Submit a receipt and update entitlements |

---

## Non-Goals

- Do not implement real App Store / Google Play receipt cryptography. Use `.muara/data/unified_matrix.json` lookup keys, consistent with v2 receipt validation design.
- Do not add RevenueCat to v1. This initiative is explicitly out of scope for v1.
- Do not change v1 provider emulation behavior.

---

## Reference

- Original v1 prompt: `prompts/14-revenuecat-emulation.md` (superseded by this initiative).
- RevenueCat API docs: https://docs.revenuecat.com
