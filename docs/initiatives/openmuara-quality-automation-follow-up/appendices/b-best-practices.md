> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix B — Best Practices for Quality Gates

> **Updated:** 2026-07-06

## Visual Regression

1. **Determinism first.** Hide all timestamps, animated elements, and non-deterministic data before capturing screenshots.
2. **No secrets in screenshots.** Never capture API keys, webhook secrets, or PII in baseline images; use deterministic fake data.
3. **Isolated data.** Use a fresh SQLite path or an in-memory store for visual baseline runs so prior state does not leak.
4. **Reviewable diffs.** Store baseline images in git so PRs show visual changes as binary diffs.
5. **Intentional updates.** Provide a documented `--update-snapshots` command and require the updated images to be committed separately.
6. **Start non-blocking.** Visual diffs are inherently sensitive; promote to required only after proven stable.

## Mutation Testing

1. **Target changed packages.** Running mutation testing on the whole repo is too slow for PR feedback.
2. **Pin the tool version.** `gremlins` is evolving; pin a version in CI to avoid surprise behavior changes.
3. **Set a realistic threshold.** A too-high threshold encourages low-value tests; a too-low threshold misses bugs.
4. **Treat mutants as signals.** If a mutant survives, ask whether the code is too complex or the test is too weak.

## Coverage Regression

1. **Module-level granularity.** Global coverage can hide drops in a changed package.
2. **Baseline from target branch.** Compare the PR coverage against `main`/`dev`, not the PR’s own partial run.
3. **Allow documented overrides.** A refactor that removes dead code may legitimately drop coverage; record the rationale.

## Error Codes

1. **Additive only.** Add codes to existing errors; do not remove or rewrite public error messages without sign-off.
2. **Stable codes.** Once a code is public, treat it as part of the API contract.
3. **Domain grouping.** Keep codes grouped by area (E1xx generic, E2xx config, E3xx provider, etc.).
4. **Log and return.** Include the code in logs and in API responses where appropriate.

## Recurring Process

1. **Idempotent scheduling.** A scheduled workflow should not create duplicate open issues.
2. **Actionable content.** The auto-created issue should link to the bug-hunt prompts and the current release milestone.
3. **Easy to skip.** Provide a label or manual override for releases that do not need a bug hunt.
