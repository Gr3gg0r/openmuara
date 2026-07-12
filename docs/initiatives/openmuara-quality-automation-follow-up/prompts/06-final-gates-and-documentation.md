> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P06 — Final Gates & Documentation

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** P01–P05
> **Target files:** `runbooks/quality-gates.md`, `TRACKING.md`, `HANDOFF.md`, `REVIEW_CHECKLIST.md`
> **Status:** ⬜

## Goal

Run the full quality matrix, update all runbooks and trackers, promote stable gates to required, and prepare the branch for PR to `dev`.

## Tasks

- [ ] Run every gate locally and in CI.
- [ ] Verify CI wall time remains under 10 minutes.
- [ ] Promote gates that have been stable for at least three PRs to required checks.
- [ ] Update `runbooks/quality-gates.md` with every new gate and its local command.
- [ ] Complete `TRACKING.md`, `HANDOFF.md`, `DECISIONS.md`, `RISKS.md`, and `REVIEW_CHECKLIST.md`.
- [ ] Update `CHANGELOG.md` with a quality-automation snippet.

## Acceptance Criteria

- [ ] All prompts marked `✅` in `TRACKING.md`.
- [ ] All required gates pass on the PR.
- [ ] `runbooks/quality-gates.md` documents every new gate.
- [ ] `REVIEW_CHECKLIST.md` is complete and signed off.
- [ ] `CHANGELOG.md` has a release-notes snippet.

## Completion Checklist

- [ ] Full gate suite run locally and in CI.
- [ ] `TRACKING.md` updated with final status and all commits.
- [ ] `HANDOFF.md` updated with final state and no open blockers.
- [ ] `RISKS.md` updated with all resolved risks closed.
- [ ] `DECISIONS.md` updated with final statuses.
- [ ] Root `TRACKING.md` updated to `🟢 Completed`.
- [ ] Feature branch deleted after merge.

## Quality Gates

Run the full suite before committing:

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test:ci
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
cd web/dashboard && node scripts/a11y-contrast-check.js
cd web/dashboard && npm run test:visual-baseline
```

## Notes

- Any gate that is still flaky must remain non-blocking with a documented plan in `RISKS.md`.
- Do not merge until the human reviewer has completed `REVIEW_CHECKLIST.md`.
