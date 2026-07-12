> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix E — Bug Register Format

> **Updated:** 2026-07-06

Use these columns in `TRACKING.md` for every confirmed bug.

| ID | Severity | Area | Summary | Reproduction | Finding File | Root Cause Category | Regression Test | Status | Commit | Introduced By | Fixed By |
|----|----------|------|---------|--------------|--------------|---------------------|-------------------|--------|--------|---------------|----------|
| B001 | P1 | webhook | Dispatcher panics on nil store | `TestDispatcherNilStore` | `findings/B001-dispatcher-nil-store.md` | nil guard | `TestDispatcherNilStore` | open | — | unknown | — |

## Severity rubric

- **P0** — Crash, security vulnerability, data loss, or a completely broken primary flow.
- **P1** — Broken feature, UX regression, or incorrect provider behavior that blocks a common use case.
- **P2** — Polish, edge case, or cosmetic issue.

## Root cause categories

nil guard, race condition, config drift, validation gap, routing mismatch, provider contract drift, UI state bug, a11y markup, test flake, documentation gap, dependency vulnerability.

## Finding file naming convention

```
findings/BXXX-<short-kebab-description>.md
```

Example: `findings/B001-dispatcher-nil-store.md`
