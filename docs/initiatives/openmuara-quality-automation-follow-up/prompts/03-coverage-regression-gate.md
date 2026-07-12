> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P03 — Coverage Regression Gate

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** P01
> **Target files:** `.github/workflows/coverage-comment.yml`, `scripts/check-coverage-regression.sh`
> **Status:** ⬜

## Goal

Evolve the existing PR coverage-comment bot into a required check that fails when any changed module drops coverage compared to the target branch.

## Tasks

- [ ] Review the existing `.github/workflows/coverage-comment.yml` output format.
- [ ] Add a step that computes per-module coverage for both the PR and the target branch (`main`/`dev`).
- [ ] Identify modules changed by the PR (e.g., via `git diff --name-only`).
- [ ] Fail the check if any changed module’s coverage is lower than the target branch baseline.
- [ ] Keep the friendly comment but update it to mention the gate status.
- [ ] Start as non-blocking; promote after three stable PRs.

## Acceptance Criteria

- [ ] Coverage gate reports per-module deltas in the PR comment.
- [ ] Gate fails when a changed module loses coverage.
- [ ] Gate passes when coverage is unchanged or improved.
- [ ] A documented override path exists for intentional coverage drops.

## Completion Checklist

- [ ] Coverage workflow updated and tested on a PR with a coverage drop.
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and override process.
- [ ] `HANDOFF.md` updated with what was done and P04 next steps.
- [ ] `DECISIONS.md` updated if the comparison baseline or override rules changed.

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

- Use `go test -coverprofile` and `go tool cover -func` for per-package numbers.
- The override path should require a `DECISIONS.md` entry, not a magic label.
