> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v2 — RevenueCat Emulation — Decisions

---

## D001 — RevenueCat deferred from v1 to v2

- **Status:** Decided
- **Context:** Prompt 14 originally planned RevenueCat emulation for v1. However, v1's focus is single-charge payment emulation. RevenueCat is a subscription/entitlement platform with concepts (subscribers, offerings, trials, renewals, mobile receipts) that do not fit the single-charge model.
- **Decision:** RevenueCat emulation is explicitly out of scope for v1 and moved to v2.
- **Consequences:**
  - v1 providers remain focused on one-time charges.
  - v2 will introduce subscription and purchase emulation, including RevenueCat.
  - No RevenueCat code, migrations, or endpoints in v1.
