> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix A — Gaps & Enhancements

> **Updated:** 2026-07-06

This appendix documents the gaps left after the bug-hunt E1–E12 implementation and the enhancements this follow-up initiative adds to close them.

## Gap 1: Visual baseline is not enforced in CI

**Current state:** `npm run test:visual-baseline` exists and produces screenshots, but it is not run in CI.
**Risk:** Unintended UI changes can merge undetected.
**Enhancement:** Add a CI job that runs the visual baseline diff on PRs touching `web/dashboard/`. Start as non-blocking, promote to required after stability is proven.
**Gold-standard touch:** aligns with `openmuara-a11y-usability-polish` (layout regressions are caught early).

## Gap 2: Mutation testing is documented but not run

**Current state:** `docs/bug-hunt-process.md` describes how to run `gremlins`, but CI does not run it.
**Risk:** Tests may cover lines without catching real bugs.
**Enhancement:** Add a CI mutation-testing job for packages changed by the PR, with a 70% threshold.
**Gold-standard touch:** aligns with `openmuara-testing-gold-standard` (tests must actually kill mutants).

## Gap 3: Coverage bot only comments; it does not block

**Current state:** `.github/workflows/coverage-comment.yml` posts a coverage summary comment.
**Risk:** Module-level coverage regressions can still merge.
**Enhancement:** Extend the workflow to compute per-module coverage delta and fail the check when a changed module drops.
**Gold-standard touch:** aligns with `openmuara-v1-solid-gold` (≥80% module coverage target).

## Gap 4: Error-code taxonomy is only used in the webhook dispatcher

**Current state:** `internal/errcode` exists and is used in `internal/webhook/dispatcher.go`.
**Risk:** Provider and API errors still lack stable codes, making debugging and bug classification harder.
**Enhancement:** Adopt `errcode` in all provider packages and in API error responses, without changing existing messages.
**Gold-standard touch:** aligns with `openmuara-ux-excellence` (plain-language, debuggable errors).

## Gap 5: Recurring bug-hunt process is manual

**Current state:** `docs/bug-hunt-process.md` describes when and how to run bug hunts.
**Risk:** The process is forgotten between releases.
**Enhancement:** Add a scheduled GitHub workflow that opens a bug-hunt prep issue before each release.
**Gold-standard touch:** institutionalizes quality discipline without adding ceremony.

## Gap 6: `KNOWN_ISSUES.md` can drift from the bug-hunt register

**Current state:** Root `KNOWN_ISSUES.md` is kept in sync manually.
**Risk:** Deferred bugs are documented in one place but missing from the user-facing file.
**Enhancement:** Add a CI script that warns when deferred items in `docs/initiatives/openmuara-bug-hunt/RISKS.md` are not present in root `KNOWN_ISSUES.md`.
**Gold-standard touch:** honest, user-facing communication.

## Good-to-Have Enhancements (Not Required)

- Visual baseline per theme (light/dark) and per viewport size.
- Mutation testing expansion to all packages once initial targets are stable.
- Automated changelog generation from PR labels.
- Provider golden-file diff guard (fail if golden files change without explicit update).
- CI job summaries with direct links to failing gate logs.
- Nightly full quality matrix to catch tooling drift.

See `appendices/d-recommendations-roadmap.md` for the complete register of future and low-priority items.
