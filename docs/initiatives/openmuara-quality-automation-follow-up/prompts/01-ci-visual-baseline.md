> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P01 — CI Visual Baseline

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** —
> **Target files:** `.github/workflows/ci.yml`, `web/dashboard/e2e/visual-baseline.spec.ts`, `web/dashboard/playwright.config.ts`
> **Status:** ⬜

## Goal

Run the existing `npm run test:visual-baseline` in CI on PRs that touch the dashboard, and make it a reliable, reviewable gate.

## Tasks

- [ ] Confirm `npm run test:visual-baseline` passes locally with a clean worktree.
- [ ] Add a path-filtered CI job that runs the visual baseline only when `web/dashboard/**` or `internal/ui/**` changes.
- [ ] Ensure the CI job installs Playwright Chromium and caches it.
- [ ] Document how to update baselines intentionally (`--update-snapshots`).
- [ ] Start the job as non-blocking; collect stability data.

## Acceptance Criteria

- [ ] Visual baseline job runs on a dashboard PR and produces a clear diff on intentional UI change.
- [ ] Baseline images remain in `docs/initiatives/openmuara-bug-hunt/findings/visual-baseline/`.
- [ ] CI failure message explains how to update snapshots if the change is intentional.
- [ ] Local and CI commands are identical.

## Completion Checklist

- [ ] CI job added and tested (can use a draft PR).
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and stability notes.
- [ ] `HANDOFF.md` updated with what was done and P02 next steps.
- [ ] `DECISIONS.md` updated if the path filter or promotion timing changed.

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

- Do not promote the job to required until it has passed reliably for at least three unrelated PRs.
- If flakiness appears, log it in `RISKS.md` and consider hiding additional dynamic elements.
