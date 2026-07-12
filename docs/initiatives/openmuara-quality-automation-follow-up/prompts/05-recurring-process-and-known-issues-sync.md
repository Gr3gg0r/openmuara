> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P05 — Recurring Process & KNOWN_ISSUES Sync

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** P02–P04
> **Target files:** `.github/workflows/`, `scripts/check-known-issues.sh`, `docs/bug-hunt-process.md`
> **Status:** ⬜

## Goal

Automate the recurring bug-hunt reminder and add a CI check that keeps root `KNOWN_ISSUES.md` in sync with the bug-hunt risk register.

## Tasks

- [ ] Add a scheduled GitHub workflow that opens a bug-hunt prep issue before each release.
- [ ] Ensure the workflow is idempotent: skip if an open bug-hunt prep issue already exists.
- [ ] Add a script that compares deferred items in `docs/initiatives/openmuara-bug-hunt/RISKS.md` with root `KNOWN_ISSUES.md`.
- [ ] Run the sync script in CI as a warning initially; promote to required once stable.
- [ ] Update `docs/bug-hunt-process.md` with the new automation links.

## Acceptance Criteria

- [ ] Scheduled issue is created with the correct labels and links to prompts.
- [ ] Duplicate issues are not created.
- [ ] Sync script reports any missing deferred items clearly.
- [ ] Documentation reflects the automated workflow.

## Completion Checklist

- [ ] Workflow and script added and tested.
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and trigger cadence.
- [ ] `HANDOFF.md` updated with what was done and P06 next steps.
- [ ] `DECISIONS.md` updated if the sync rules or skip conditions changed.

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

- Use a scheduled cron workflow; avoid :00 and :30 to reduce GitHub load.
- The sync script should allow an explicit "intentionally not listed" marker to avoid false positives.
