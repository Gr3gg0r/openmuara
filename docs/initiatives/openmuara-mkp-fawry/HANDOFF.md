> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara MKP Fawry Integration — Handoff

> **Created:** 2026-07-01
> **Last Updated:** 2026-07-08
> **Status:** ⏸️ Suspended

---

## Current Context

MKP v2 needs a Fawry emulator for local development and CI. OpenMuara already
supports charge creation, an escape page, and webhook dispatch, but the MKP
handler expects additional order statuses, configurable delays, and a
subscription/prepaid journey hint. This initiative was created to track that
work.

Implementation was started on `feat/mkp-fawry` and progressed through prompts
01–04. Step 05 (docs and CI) was in progress when the initiative was suspended.

## Branch

`feat/mkp-fawry` — uncommitted work is stashed as `WIP: suspend MKP Fawry implementation`.

## What has been done

- Prompts 01–04 implemented and passing quality gates on `feat/mkp-fawry`.
- Step 05 (docs/runbooks and smoke test) started but not completed.
- Initiative docs created in `docs/initiatives/openmuara-mkp-fawry/`.

## Prompt Inventory

| Step | Title | Status |
|------|-------|--------|
| 01 | Fawry state extensions | ✅ |
| 02 | Response delay config | ✅ |
| 03 | Billing type and journey | ✅ |
| 04 | Escape page and webhook shape | ✅ |
| 05 | Docs and CI | ⏸️ |

## Decisions already made

See `DECISIONS.md` and `docs/initiatives/openmuara-mkp-fawry/TRACKING.md`:

- Existing Fawry behavior stays backward-compatible; new fields are optional.
- `response_delay_ms` applies to outgoing webhook dispatch only, not the escape redirect.
- `GET /fawry/payment-status` added under VAL01 to let clients verify payment status by `merchantRefNum`; signature required.
- MKP delegates Fawry simulation to OpenMuara via `POST /fawry/charge` and `POST /fawry/simulate` when `OPENMUARA_URL` is configured.

## Next step

Resume on `feat/mkp-fawry`:

1. `git checkout feat/mkp-fawry`
2. `git stash pop` (stash name: `WIP: suspend MKP Fawry implementation`)
3. Complete step 05: update provider docs, MKP requirements status, runbook, and smoke test.
4. Run quality gates and commit.

## Open questions

- None remaining from prompts 01–04.
- Step 05: decide whether to add a dedicated MKP Fawry runbook or fold the guidance into `docs/providers/fawry.md`.
