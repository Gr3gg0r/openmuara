> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Bug Hunt — Prerequisites & Assumptions

> **Updated:** 2026-07-06

## Required Tools

- Go toolchain (version matching `go.mod`).
- Node.js + npm (for `web/dashboard`).
- `golangci-lint` installed and available on `$PATH`.
- `govulncheck` (optional but recommended for P01).
- `curl` or `httpie` for runtime provider checks.
- Playwright MCP access for P01 baseline and P06 visual sign-off.
- Git with the ability to create branches and commits.

## Branch Base

This initiative lives on `feat/bug-hunt`, which was branched from `feat/dashboard-mailpit-redesign`. That means the branch already contains the completed Mailpit-style dashboard redesign and its P01–P06 implementation.

Do **not** rebase onto `dev` unless the dashboard redesign has been merged and you are explicitly told to do so.

## Assumptions

- The dashboard redesign branch is green at the start of P01.
- No new features or providers will be added during this bug hunt.
- The user is available for P0/P1 integration-fix sign-off within a reasonable time frame.
- All provider emulation tests can run offline (no real provider calls).
- The local environment has enough resources to run `go test -race ./...`.

## Time-box Guidance

| Prompt | Suggested Effort | Why |
|--------|------------------|-----|
| P01 Reconnaissance | 1–2 sessions | Running gates, searching code, exercising endpoints, and capturing baselines takes time. |
| P02 Triage | 1 session | Reproduction validation and batch planning. |
| P03 Fix Batch 1 | 1–2 sessions | Highest-impact fixes with regression tests. |
| P04 Fix Batch 2 | 1–2 sessions | Remaining fixes or documented deferrals. |
| P05 Regression Tests & Quality Gates | 1 session | Integration tests, coverage, full gate suite. |
| P06 Visual Sign-off | 1 session | Playwright MCP capture and comparison. |

If a single bug consumes more than one session, escalate in `HANDOFF.md` and `RISKS.md`.

## Baseline Capture

Before P01 starts, record the actual environment here:

```bash
go version:
node version:
npm version:
golangci-lint --version:
OS:
commit at start of P01:
```

---

## Communication & Escalation

- **Daily checkpoint:** Update `HANDOFF.md` after every prompt, even if the only update is "no progress today."
- **Blockers:** If a quality gate fails and cannot be resolved within one session, mark the prompt `❌` in `TRACKING.md`, log the blocker in `RISKS.md`, and summarize it in `HANDOFF.md`.
- **User sign-off:** For P0/P1 integration fixes, request sign-off in `DECISIONS.md` before writing code. If sign-off is delayed, move the bug to deferred and update `KNOWN_ISSUES.md`.
- **Scope questions:** If a bug is a missing feature, a large refactor, or touches v2 scope, defer it and link to the relevant backlog initiative.
- **Visual regressions:** Any UI change discovered during P06 must be documented with a Playwright MCP screenshot and either fixed or escalated.
