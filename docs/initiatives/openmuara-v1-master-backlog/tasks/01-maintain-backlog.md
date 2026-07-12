> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# 01 — Maintain the Master Backlog

## Objective

After any product-code change, the master backlog and its source tracker accurately reflect the new status, commit hash, and any newly discovered risks or issues.

## Background

The master backlog is a consolidated view. If it is not kept current, agents will duplicate work or start from stale assumptions. This task defines the minimal update ritual.

## Constraints

- Do not modify product code in this task.
- Do not reprioritize v2-frozen items without human approval.
- Use `<repo-root>` for repo paths; never hard-code absolute worktree paths.

## Acceptance Criteria

- [ ] The changed backlog item has an updated status in `TRACKING.md`.
- [ ] The source tracker (`root`, `project`, or `v1-solid`) is also updated.
- [ ] Any new risk or known issue is logged in `RISKS.md` or `KNOWN_ISSUES.md`.
- [ ] `HANDOFF.md` reflects the latest session state.
- [ ] No absolute filesystem paths are introduced into backlog files.

## Test Expectations

- Search backlog files for absolute filesystem paths and confirm no matches.
- All markdown files render without broken internal links.

## Rollback Trigger

If the backlog update accidentally deletes or corrupts tracker rows, revert the docs commit.
