> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt P04 — Docs and Release Notes

> **Initiative:** OpenMuara Dark Mode
> **Target:** `<repo-root>/`
> **Branch:** `feat/dark-mode`
> **Depends on:** P01, P02, P03

---

## Goal

Document the new dark mode for users and contributors, and prepare a release-notes snippet.

## Why now

A feature that users cannot discover or contributors cannot extend is only half done. Clear docs close the loop.

## Scope

### In scope

- Update `README.md` to mention dark mode support in the dashboard, provider pages, and examples.
- Add a short section to `docs/initiatives/openmuara-dark-mode/README.md` summarizing the shipped implementation.
- Add a `CHANGELOG.md` entry under the next unreleased version.
- Document the design-token naming convention for future contributors (can live in `docs/initiatives/openmuara-dark-mode/README.md`).
- Document the keyboard shortcut (`d`) for toggling theme in the dashboard.
- Update the initiative's `HANDOFF.md`, `DECISIONS.md`, and `KNOWN_ISSUES.md` to reflect the completed state.

### Out of scope

- Rewriting the full user guide.
- Adding screenshots (optional; if added, keep them out of git).

## Acceptance criteria

- [ ] `README.md` mentions dashboard dark mode.
- [ ] `CHANGELOG.md` has an unreleased entry describing the feature.
- [ ] Contributor docs explain how to use theme tokens.
- [ ] All quality gates still pass.

## Deliverables

- Doc changes on `feat/dark-mode`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
