> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1.1 — Subscriptions — Handoff

> **Updated:** 2026-07-03
> **Initiative:** `docs/initiatives/openmuara-v1-1-subscriptions/`
> **Branch:** `feat/v1-1-subscriptions`
> **Status:** ⬜ Not Started

---

## Last Session Summary

Created the v1.1 subscription initiative. No product code has been written yet.

- Documented v1.1 subscription focus: Stripe first, then Malaysian recurring gateways.
- Listed Malaysian providers with subscription/recurring support.
- Created prompts for Stripe subscriptions (01) and SenangPay subscriptions (02).

---

## Next Steps

1. When v1.1 work begins, create the `feat/v1-1-subscriptions` branch from `dev`.
2. Start with `prompts/01-stripe-subscriptions.md`.
3. Design the subscription engine abstraction so later providers can reuse it.
4. Run quality gates and commit after each prompt.

---

## Open Questions

- Which Malaysian provider should follow Stripe? SenangPay is a likely candidate because it is already in OpenMuara and has documented subscription/recurring support.
- Should the subscription engine live in `internal/engine/` or a new `internal/subscription/` package?

---

## Notes

- Do not implement RevenueCat or mobile receipt validation in v1.1.
- Keep v1 single-charge behavior unchanged.
