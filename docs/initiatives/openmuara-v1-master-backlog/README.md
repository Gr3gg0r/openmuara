> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Master Backlog

> **Status:** Active | **Updated:** 2026-06-28
> **Scope:** One consolidated, priority-ranked view of every known OpenMuara v1 item.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** Unassigned / TBD
> **Repo:** `<repo-root>`
> **Product Branch:** `dev`

---

## Purpose

This initiative is a living backlog that aggregates work from:

- Root `TRACKING.md` (prompts 01–19, tasks T01–T02).
- `docs/projects/openmuara-v1/TRACKING.md` (execution tracker with phases and commit hashes).
- `docs/initiatives/openmuara-v1-solid/` (active regression-fix prompts).

Every item is tagged **High / Medium / Low** so the next agent or human can see what matters most without re-reading three trackers.

---

## Quick Start for Agents

1. Read `AGENTS.md` at the repo root.
2. Read `TRACKING.md` in this folder.
3. Pick the highest-priority item that is not ✅ or ❄️.
4. If the item already has a prompt/task elsewhere, execute that prompt; otherwise create one.
5. After completing work, update **this** `TRACKING.md` and the source tracker it came from.

---

## Priority Rules

| Priority | Rule |
|----------|------|
| **High** | Regression, API contract break, core runtime gap, or daily-use blocker. |
| **Medium** | Provider hardening, observability, packaging, or docs that improve solid v1. |
| **Low** | Deferred, nice-to-have, or explicitly frozen for v2. |

---

## Initiative Structure

```
docs/initiatives/openmuara-v1-master-backlog/
├── README.md              # This file
├── HOWTO.md               # How to use and maintain this backlog
├── INIT.md                # One-time setup (temporary; delete at closeout)
├── TRACKING.md            # Master backlog table
├── PREREQUISITES.md       # Lightweight pre-flight checklist
├── KNOWN_ISSUES.md        # Merged known issues and boundaries
├── RISKS.md               # Merged risk register
├── REFERENCES.md          # Links to source trackers and docs
├── DECISIONS.md           # Decision log
├── HANDOFF.md             # Session continuity
├── .gitignore             # Ignore agent artifacts
├── prompts/
│   ├── _template.md
│   └── 01-maintain-backlog.md   # Update ritual after product-code changes
└── tasks/
    ├── _template.md
    └── 01-maintain-backlog.md   # Spec for the update ritual
```

---

## Conventions

- **Read-only tracker.** Do not commit product code here.
- **Update after every product-code change.** If you finish a prompt, update the matching row in `TRACKING.md`.
- **Use `<repo-root>`** for repo paths; never hard-code absolute worktree paths.
- **Never commit directly to `main`.**

---

## Completion Criteria

This backlog is healthy when:

- [x] Every active item has a clear priority, status, owner, and entry point.
- [x] No item is duplicated across trackers without a cross-reference.
- [x] All High items have an executable prompt/task or an assigned owner.
- [x] `RISKS.md` and `KNOWN_ISSUES.md` are current.
