> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1.1 — Subscriptions — Decisions

---

## D001 — v1.1 focuses on subscriptions, starting with Stripe Billing

- **Status:** Decided
- **Context:** v1 covers single-charge payment emulation. The next natural expansion is subscriptions/recurring payments.
- **Decision:**
  - v1.1 adds subscription emulation.
  - Stripe Billing is the first provider because it has the most mature API and OpenMuara already has a Stripe Checkout provider.
  - Malaysian gateways with recurring support will follow.
- **Consequences:**
  - v1.1 introduces products, prices, customers, subscriptions, invoices, and subscription webhooks.
  - A reusable subscription engine abstraction will be designed so other providers can plug in.

---

## D002 — Malaysian recurring gateway candidates

- **Status:** Decided
- **Context:** Several Malaysian gateways support recurring payments or subscriptions. We need to decide which to emulate after Stripe.
- **Decision:** Candidates are SenangPay, iPay88, Billplz, ToyyibPay, Fiuu, 2C2P, eGHL, Curlec, HitPay, Xendit, and Airwallex. Priority will be given to providers already in OpenMuara (SenangPay, iPay88, Billplz, ToyyibPay) before adding new gateways.
- **Consequences:**
  - Existing providers get subscription extensions first.
  - New gateways (Fiuu, 2C2P, eGHL, Curlec, HitPay, Xendit, Airwallex) are backlog items.
