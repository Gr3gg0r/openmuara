> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara AI Slop Audit — Handoff

> **Created:** 2026-07-08
> **Last Updated:** 2026-07-08
> **Status:** ⬜ Draft

---

## Current context

This initiative was created after a UI/UX responsiveness pass on the dashboard surfaced several generic, placeholder, or AI-generated patterns. The goal is to systematically catalog and remove AI slop from the product so OpenMuara feels intentional and trustworthy.

## What has been done

- Initiative scaffold created under `docs/initiatives/openmuara-ai-slop-audit/`.
- `README.md`, `TRACKING.md`, `KNOWN_ISSUES.md`, `DECISIONS.md`, `RISKS.md`, and this `HANDOFF.md` written.
- 15 findings recorded across dashboard, provider metadata, docs, code, test data, examples, and website.

## What has not been done

- No code changes yet.
- No branch created.
- No decisions finalized (e.g., font choice, seed-data flag shape).

## Next step

Create the product branch `feat/ai-slop-audit` and begin **P01 — Dashboard microcopy & icons**:
1. Fix F001 (sad-face empty-state icon).
2. Fix F002 (provider placeholder descriptions) or F003 (redundant badges).
3. Capture before/after screenshots and update `TRACKING.md`.

## Open questions

- Should the seed-data flag be env-based (`MUARA_SEED_DATA=1`) or config-based (`seed_data: true` in `config.yml`)?
- Which typeface should replace the system-ui stack, and should it be self-hosted?
