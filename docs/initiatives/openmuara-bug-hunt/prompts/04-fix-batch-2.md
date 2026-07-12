> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P04 — Fix Batch 2

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** P03
> **Target files:** Per-bug target files
> **Status:** ⬜

## Goal

Fix the remaining confirmed bugs from P02, or document a clear rationale for deferral.

## Tasks

- [ ] Confirm user sign-off in `DECISIONS.md` for any remaining P0/P1 integration fixes before coding.
- [ ] Fix each remaining bug in the P04 batch with minimal changes.
- [ ] Add a regression test for each fix.
- [ ] For any bug intentionally not fixed, record it in `RISKS.md` with reason, impact, and target release.
- [ ] Run quality gates after **each** commit.
- [ ] Verify no dashboard redesign invariant is broken.
- [ ] Update `TRACKING.md` bug register with status and commit hashes.
- [ ] Update `HANDOFF.md` with what was fixed or deferred.

## Acceptance Criteria

- [ ] All confirmed P02 bugs are either fixed or explicitly deferred with rationale.
- [ ] Each fix has a regression test.
- [ ] No dashboard redesign invariant regressed.
- [ ] All quality gates pass.
- [ ] `DECISIONS.md` contains sign-off records for any integration fixes.

## Quality Gates

Run before committing:

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
```

## Notes

- If a deferred bug is P0 or P1, get user sign-off before deferring.
- If a fix reveals a deeper issue, stop, document it in `RISKS.md`, and decide with the user whether to expand scope.
- Do not opportunistically refactor unrelated code.
