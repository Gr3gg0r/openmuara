> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P02 — Mutation Testing Gate

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** P01
> **Target files:** `.github/workflows/ci.yml`, `scripts/`, `internal/webhook/`, `internal/engine/`, `internal/fawry/`
> **Status:** ⬜

## Goal

Add a mutation-testing job using `gremlins` to verify that regression tests actually detect bugs, not just cover lines.

## Tasks

- [ ] Install `gremlins` locally and measure the baseline mutation score for `internal/webhook`, `internal/engine`, and `internal/fawry`.
- [ ] Choose a threshold (start at 70%) and document it in `DECISIONS.md`.
- [ ] Add a CI job that runs `gremlins` on packages changed by the PR (or the curated list if no reliable changed-package detection exists).
- [ ] Cache the `gremlins` binary to keep CI fast.
- [ ] Start the job as non-blocking commentary; promote after baseline is stable.

## Acceptance Criteria

- [ ] Mutation testing runs in CI and reports a score per package.
- [ ] Score threshold is documented and achievable.
- [ ] CI time remains under the 10-minute budget.
- [ ] Local reproduction command is documented.

## Completion Checklist

- [ ] CI job added and baseline scores captured.
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and threshold.
- [ ] `HANDOFF.md` updated with what was done and P03 next steps.
- [ ] `DECISIONS.md` updated if the threshold or target packages changed.

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

- If a package cannot reach 70% without major test work, lower the threshold for that package and record the rationale in `DECISIONS.md`.
- Do not block PRs on mutation testing until the score is stable for at least three PRs.
