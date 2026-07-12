> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1.1 — Subscriptions — Risk Register

---

## R001 — Subscription engine over-engineering

- **Likelihood:** Medium
- **Impact:** Medium
- **Description:** Designing a generic subscription engine too early may create abstraction overhead before concrete provider needs are understood.
- **Mitigation:** Build the engine around Stripe Billing first, then refactor for reuse as SenangPay and other providers are added.

---

## R002 — v1 single-charge behavior regression

- **Likelihood:** Medium
- **Impact:** High
- **Description:** Subscription code changes could accidentally alter v1 one-time charge behavior.
- **Mitigation:** Keep subscription state separate from the charge transaction ledger. Run the full v1 test suite after every subscription change.

---

## R003 — Scope overlap with v2 RevenueCat

- **Likelihood:** Low
- **Impact:** Medium
- **Description:** Subscription work might drift into entitlement/mobile receipt territory that belongs to v2.
- **Mitigation:** Clearly exclude mobile receipts, App Store / Play Store validation, and RevenueCat from v1.1. Reference root `DECISIONS.md` D037 and the v2 RevenueCat initiative.

---

## R004 — Provider API divergence

- **Likelihood:** High
- **Impact:** Medium
- **Description:** Malaysian gateways use different terminology and flows for recurring payments (plans, mandates, auto-billing, RPP).
- **Mitigation:** Document each provider's actual recurring API before implementation. Do not force all providers into a Stripe-shaped model if the real API differs.
