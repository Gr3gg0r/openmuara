> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# PXX — Title

> **Initiative:** OpenMuara Quality Automation Follow-Up
> **Depends on:** (prompt numbers or "—")
> **Target files:** (file paths)
> **Status:** ⬜

## Goal

(One-paragraph description of what this prompt achieves.)

## Tasks

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2

## Completion Checklist

- [ ] Code/tests written and passing locally.
- [ ] `TRACKING.md` updated with status `✅`, commit hash, and notes.
- [ ] `HANDOFF.md` updated with what was done and what is next.
- [ ] `RISKS.md` updated if any risk was mitigated or discovered.
- [ ] `DECISIONS.md` updated if any threshold or scope decision changed.

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

- (Any special considerations, sign-off needs, or rollback plans.)
