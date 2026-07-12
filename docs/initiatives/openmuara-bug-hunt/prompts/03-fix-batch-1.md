> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P03 — Fix Batch 1

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** P02
> **Target files:** Per-bug target files
> **Status:** ⬜

## Goal

Fix the 2–3 highest-impact, lowest-risk bugs identified in P02, each with a regression test and a focused commit.

## Per-Bug Fix Template

For each bug in the P03 batch, create or update exactly one focused commit with this shape:

1. **Reproduction first** — add a failing regression test or reproduction script.
2. **Minimal fix** — change only what is required to fix the bug.
3. **Verify** — run the full quality gate suite and confirm the regression test passes.
4. **Commit** — `fix(scope): short description (Bxxx)`.
5. **Track** — update `TRACKING.md` bug register with status `fixed`, commit hash, and `Fixed By`.
6. **Handoff** — update `HANDOFF.md` with what changed and any caveats.

## Tasks

- [ ] Confirm user sign-off in `DECISIONS.md` for any P0/P1 integration fix before coding.
- [ ] Fix each bug in the P03 batch with minimal changes.
- [ ] Add a regression test for each fix.
- [ ] Run quality gates after **each** commit.
- [ ] Verify no dashboard redesign invariant is broken.
- [ ] Update `TRACKING.md` bug register with status and commit hash.
- [ ] Update `HANDOFF.md` with what was fixed.

## Acceptance Criteria

- [ ] All P03 bugs marked fixed in `TRACKING.md`.
- [ ] Each fix has a regression test that fails before and passes after.
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

- One logical change per commit.
- If a fix turns out to be riskier than expected, stop and move it to P04 or deferred; update `TRACKING.md` and `RISKS.md`.
- Do not opportunistically refactor unrelated code.
