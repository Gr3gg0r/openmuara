> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# OpenMuara Quality Automation Follow-Up — Prerequisites & Assumptions

> **Updated:** 2026-07-06

## Required Tools

- Go toolchain (version matching `go.mod`).
- Node.js + npm (for `web/dashboard`).
- `golangci-lint` installed and available on `$PATH`.
- `govulncheck` (used by existing CI; confirm it still passes after changes).
- `npx playwright install chromium` for local visual baseline runs.
- `gremlins` (to be installed during P02):
  ```bash
  go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
  ```
- Git with the ability to create branches and commits.
- GitHub CLI (`gh`) or web access if workflow changes need to be validated.

## Branch Base

This initiative lives on `feat/quality-automation-follow-up`, branched from `dev`.
`dev` already contains the completed bug-hunt enhancements (E1–E12).

Do **not** commit directly to `dev` or `main`.

## Assumptions

- The bug-hunt branch is merged and `dev` is green at the start of P01.
- No new providers or features will be added during this initiative.
- GitHub Actions runners can install Playwright Chromium and `gremlins` within the 10-minute CI budget.
- All provider emulation tests can run offline (no real provider calls).
- The local environment has enough resources to run `go test -race ./...` and Playwright.

## Time-box Guidance

| Prompt | Suggested Effort | Why |
|---|---|---|
| P01 CI visual baseline | 1 session | Wire existing script into CI; tune flake tolerance. |
| P02 Mutation testing gate | 1–2 sessions | Install, tune threshold, pick packages, handle flakes. |
| P03 Coverage regression gate | 1 session | Convert comment bot into required check. |
| P04 Provider errcode adoption | 1–2 sessions | Wrap errors without changing behavior; add tests. |
| P05 Recurring process & KNOWN_ISSUES sync | 1 session | Scheduled workflow and sync script. |
| P06 Final gates & documentation | 1 session | Full gate suite, runbook update, handoff. |

If any gate proves unstable for more than one session, mark the prompt `❌` in `TRACKING.md`, log the blocker in `RISKS.md`, and summarize in `HANDOFF.md`.

## Baseline Capture

Before P01 starts, record the actual environment here:

```bash
go version:
node version:
npm version:
golangci-lint --version:
gremlins version (after install):
OS:
commit at start of P01:
```

## Communication & Escalation

- **Daily checkpoint:** Update `HANDOFF.md` after every prompt, even if the only update is "no progress today."
- **Blockers:** If a quality gate fails and cannot be resolved within one session, mark the prompt `❌` in `TRACKING.md`, log the blocker in `RISKS.md`, and summarize it in `HANDOFF.md`.
- **User sign-off:** P04 provider errcode adoption is additive, but if any change alters an existing error message used by clients, request sign-off in `DECISIONS.md` before merging.
- **Flaky gates:** A flaky required check must be demoted to commentary until it is stable. Document the demotion in `DECISIONS.md`.
