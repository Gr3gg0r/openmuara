---
id: bug-hunt-process
title: OpenMuara Bug-Hunt Process
---

> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Bug-Hunt Process

> **Updated:** 2026-07-06

This document describes the recurring bug-hunt practice for OpenMuara and how to run optional mutation testing.

## Recurring Bug-Hunt Sprints (E6)

Before every release, run a time-boxed bug hunt using the `docs/initiatives/openmuara-bug-hunt/` prompts:

1. **Schedule:** one focused sprint before the release branch is cut.
2. **Time-box:** 1–2 sessions for reconnaissance, 1 session for triage, and 2–4 sessions for fixes.
3. **Entry point:** `docs/initiatives/openmuara-bug-hunt/prompts/01-reconnaissance.md`.
4. **Exit criteria:**
   - At least 5 bugs discovered, reproduced, and documented.
   - All P0/P1 bugs fixed or explicitly deferred with user sign-off.
   - All quality gates pass.
   - P06 visual sign-off is complete.
5. **Handoff:** use `docs/initiatives/openmuara-bug-hunt/HANDOFF.md` between sessions.
6. **Post-sprint:** update `CHANGELOG.md`, root `TRACKING.md`, and `KNOWN_ISSUES.md`.

### Automated reminders

- The scheduled workflow `.github/workflows/bug-hunt-prep.yml` opens a `[bug-hunt] Prep:` issue on the first of each month. It is idempotent: it skips creation if an open bug-hunt issue already exists.
- The workflow can also be triggered manually from the Actions tab.

### KNOWN_ISSUES sync

- `scripts/check-known-issues.sh` compares deferred bug IDs in `docs/initiatives/openmuara-bug-hunt/KNOWN_ISSUES.md` with root `KNOWN_ISSUES.md`.
- The sync check runs in CI as a warning (`continue-on-error: true`) until it is proven stable; then it becomes required.
- To intentionally keep a root-only entry without a matching bug-hunt entry, add `<!-- check-known-issues:ignore -->` to `KNOWN_ISSUES.md`.

## Mutation Testing (E5)

Mutation testing validates that regression tests actually detect bugs, not just cover lines.

### Option A — Gremlins (recommended)

```bash
# Install gremlins
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest

# Run on the packages most changed by recent fixes
gremlins unleash --tags=test ./internal/webhook ./internal/engine ./internal/fawry
```

### Option B — go-mutesting

```bash
# Install go-mutesting
go install github.com/zimmski/go-mutesting/...@latest

# Run on a single package
go-mutesting --verbose --exec "go test ./internal/engine" ./internal/engine/...
```

### When to run

- After a bug-hunt fix batch, target the changed packages.
- If mutation score is below 70%, add or strengthen regression tests before merging.
- Do not block CI on mutation testing until the project maintains a stable 70%+ score.
