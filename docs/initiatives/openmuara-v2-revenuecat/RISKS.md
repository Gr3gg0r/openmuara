> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v2 — RevenueCat Emulation — Risk Register

---

## R001 — Scope creep into v1

- **Likelihood:** Medium
- **Impact:** High
- **Description:** RevenueCat code or concepts accidentally leak into v1, complicating the single-charge focus.
- **Mitigation:** Keep all RevenueCat work on `feat/v2-revenuecat`. Do not merge to `dev` until v2 branching strategy is confirmed. Update root `TRACKING.md` and `DECISIONS.md` to document the boundary.

---

## R002 — Receipt validation prerequisite

- **Likelihood:** High
- **Impact:** Medium
- **Description:** RevenueCat receipt submission depends on App Store / Google Play receipt validation, which is itself a v2 capability.
- **Mitigation:** Coordinate with the v2 receipt validation initiative. Use the same `.muara/data/unified_matrix.json` lookup-key approach.

---

## R003 — Subscription state divergence from ledger

- **Likelihood:** Medium
- **Impact:** Medium
- **Description:** Subscriber/entitlement state may not map cleanly to the v1 transaction ledger.
- **Mitigation:** Design a separate subscriber/entitlement store for v2; do not force subscriptions into the charge-oriented ledger schema.
