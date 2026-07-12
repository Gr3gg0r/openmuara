# Prompt 18 — Migration Guide

> **Status:** ✅ Completed  
> **Scope:** Document how users migrate from the legacy `muara` layout to the current `openmuara` layout.

## Goal

Provide clear, versioned migration guidance for users upgrading OpenMuara:

- What changed between legacy `muara` and `openmuara` (module path, binary name, config path).
- How to back up and restore the `.muara/` workspace.
- How to update environment variables and provider configs.
- Where to get help if a migration fails.

## Deliverables

- `docs/migration/openmuara-to-openmuara.md` — primary migration guide.
- `tasks/openmuara-migration-guide.md` — detailed task spec referenced from this prompt.
- `CHANGELOG.md` entry under `[Unreleased]` summarizing migration-relevant changes.

## Acceptance criteria

- A user reading the guide can migrate an existing `.muara/` workspace without reading source code.
- All example commands use `127.0.0.1` and clearly fake secrets.
- The guide is linked from `README.md` and the website sidebar.

## See also

- `docs/projects/openmuara-v1/TRACKING.md`
- `tasks/INDEX.md`
