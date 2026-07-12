> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix C — Post-Bug-Hunt Checklist

> **Updated:** 2026-07-06

After P06 is complete and the branch is green, run through these steps before opening a PR to `dev`.

## 1. Final Documentation Sweep

- [ ] `TRACKING.md` shows all prompts `✅` with commit hashes.
- [ ] `HANDOFF.md` has the final state, visual sign-off summary, and any open questions.
- [ ] `DECISIONS.md` has final statuses for all decisions (`Decided`, `Approved`, `Deferred`).
- [ ] `RISKS.md` has closed or accepted statuses for resolved risks.
- [ ] `KNOWN_ISSUES.md` lists any deferred bugs with rationale and target release.

## 2. Release Notes

- [ ] Add a `CHANGELOG.md` snippet summarizing fixed bugs by ID and area.
- [ ] Note any deferred P0/P1 bugs and required user actions (e.g., config changes).

## 3. Root Tracker Update

- [ ] Update `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` if the bug hunt resolves any backlog items.
- [ ] Update root `TRACKING.md` active initiatives table for `OpenMuara Bug Hunt` status (e.g., from `🟡 In Progress` to `✅ Completed`).

## 4. Branch Hygiene

- [ ] Ensure all commits on `feat/bug-hunt` are logical and squashed if necessary (do not squash dashboard redesign commits that belong to that initiative).
- [ ] Rebase onto the latest `feat/dashboard-mailpit-redesign` if new dashboard commits landed while the bug hunt was running.
- [ ] Confirm the diff contains only bug fixes, tests, and docs related to this initiative.

## 5. PR Preparation

- [ ] Open PR from `feat/bug-hunt` to `dev` (or to `feat/dashboard-mailpit-redesign` if the dashboard redesign has not yet merged).
- [ ] PR description references the bug IDs fixed and any deferred bugs.
- [ ] Attach P06 Playwright MCP screenshots or link to `HANDOFF.md`.
- [ ] Request human review using `REVIEW_CHECKLIST.md`.

## 6. Handoff to Human Reviewer

- [ ] Summarize the bug hunt: number found, number fixed, number deferred, coverage delta, and visual sign-off result.
- [ ] Highlight any P0/P1 integration fixes that required sign-off.
- [ ] Provide the list of files changed outside of `docs/initiatives/openmuara-bug-hunt/`.
