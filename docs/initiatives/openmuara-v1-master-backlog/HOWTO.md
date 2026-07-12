> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# How to Use and Maintain the Master Backlog

> **Purpose:** Guide for AI agents using this backlog as a single source of truth.

---

## What This Initiative Is

This is a **planning tracker** with one maintenance prompt. It aggregates v1 work from three sources:

- Root `TRACKING.md`
- `docs/projects/openmuara-v1/TRACKING.md`
- `docs/initiatives/openmuara-v1-solid/TRACKING.md`

It does not replace those trackers; it prioritizes them.

---

## Daily Workflow

1. Read this `HOWTO.md`.
2. Read `TRACKING.md`.
3. Pick the highest-priority item that is not ✅, ❄️, or ⏸️.
4. Open the **Entry Point** in the row — it is the executable prompt or task spec.
5. After completing the work, run the update ritual in `prompts/01-maintain-backlog.md`.

---

## Adding a New Backlog Item

1. Choose the next ID:
   - Use `P##` for root-prompt-level items.
   - Use `S##` for v1-solid items.
   - Use `T##` for task-level items.
   - Use `B##` for brand-new items that do not fit the above.
2. Add a row to `TRACKING.md` with Priority, Status, and Source.
3. If the item is non-trivial, create a matching `prompts/##-*.md` and `tasks/##-*.md`.
4. Update `RISKS.md` if the item introduces new risk.

---

## Updating Status

When an item changes status:

1. Update the row in `TRACKING.md`.
2. Update the source tracker it came from.
3. Fill the commit hash in the **Notes** column when done.
4. Update `HANDOFF.md`.

---

## Conventions

- One prompt = one logical executable step.
- Prompts must be self-contained.
- Every non-trivial prompt must have a matching `tasks/##-*.md`.
- Use `<repo-root>` for repo paths.
- Never commit directly to `main`.

---

## Project Closeout Checklist

When v1 is feature-complete and this backlog is no longer needed:

- [ ] Delete `INIT.md`.
- [ ] Verify all High/Medium items are ✅.
- [ ] Update `TRACKING.md` status to ✅ Complete.
- [ ] Update `HANDOFF.md` with final state.
- [ ] Move initiative to `docs/initiatives/archive/done/openmuara-v1-master-backlog/`.
