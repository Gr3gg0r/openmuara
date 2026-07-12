> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P05 — Regression Tests and Quality Gates

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** P03, P04
> **Target files:** Tests across affected modules
> **Status:** ⬜

## Goal

Add integration-level regression coverage for the fixed flows, verify coverage did not drop, and confirm the entire branch passes all quality gates.

## Tasks

- [ ] Add or strengthen integration tests for the most critical fixed flows.
- [ ] Run the full gate suite: build, test, race, vet, lint, frontend tests/build, bundle size, a11y contrast, smoke test.
- [ ] Review `coverage.out` (or `go test -cover`) for any dropped coverage on changed modules.
- [ ] If coverage dropped, add tests or document a justification in `DECISIONS.md`.
- [ ] Verify the dashboard redesign invariant checklist in `TRACKING.md` is still green (functional check; visual check happens in P06).
- [ ] Add a `CHANGELOG.md` release-notes snippet for fixed bugs.
- [ ] Review `RISKS.md`, `DECISIONS.md`, and `KNOWN_ISSUES.md` for accuracy.
- [ ] Update `HANDOFF.md` with final state and readiness for P06.
- [ ] Update `TRACKING.md` to mark P05 complete.

## Acceptance Criteria

- [ ] All quality gates pass.
- [ ] No test coverage drop on changed modules.
- [ ] `TRACKING.md` shows P05 complete.
- [ ] `CHANGELOG.md` has a release-notes snippet.
- [ ] Branch is ready for P06 visual sign-off.

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
./scripts/smoke-test.sh || true
```

## Notes

- Do not introduce new fixes in this prompt unless a regression is found.
- If a gate fails, fix it in this prompt or move the fix to P04 and update trackers.
- Record the final coverage delta in `HANDOFF.md`.
