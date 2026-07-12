> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v2 — RevenueCat Emulation — Handoff

> **Updated:** 2026-07-03
> **Initiative:** `docs/initiatives/openmuara-v2-revenuecat/`
> **Branch:** `feat/v2-revenuecat`
> **Status:** ⬜ Not Started

---

## Last Session Summary

Created the v2 RevenueCat initiative. No product code has been written yet.

- Moved RevenueCat out of v1 scope (Prompt 14).
- Documented v1 single-charge philosophy in root `DECISIONS.md` (D037).
- Created initiative docs under `docs/initiatives/openmuara-v2-revenuecat/`.

---

## Next Steps

1. When v2 work begins, start with `prompts/01-revenuecat-emulation.md`.
2. Create the `feat/v2-revenuecat` branch from the v2 base branch.
3. Implement the RevenueCat provider package and migrations.
4. Run quality gates and commit.

---

## Open Questions

- What is the v2 base branch? (Likely `dev` at the time v2 begins, or a dedicated `v2` branch.)
- Does v2 reuse the existing SQLite ledger schema or introduce a separate subscriber/entitlement store?

---

## Notes

- Do not implement RevenueCat in v1.
- Keep v1 provider emulation untouched.
