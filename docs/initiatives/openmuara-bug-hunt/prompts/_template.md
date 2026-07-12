> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P0X — Prompt Title

> **Initiative:** OpenMuara Bug Hunt
> **Depends on:** —
> **Target files:** `path/to/file.go`
> **Status:** ⬜

## Goal

One-sentence goal.

## Tasks

- [ ] Task one.
- [ ] Task two.
- [ ] Update `TRACKING.md` with status and commit hash.
- [ ] Update `HANDOFF.md` with what changed and any caveats.
- [ ] Update `CHANGELOG.md` if the prompt fixes user-facing bugs.
- [ ] Update `KNOWN_ISSUES.md` if any bug is deferred.

## Acceptance Criteria

- [ ] Criterion one.
- [ ] Criterion two.
- [ ] All dashboard redesign invariants preserved (for UI-impacting prompts).
- [ ] Quality gates pass.

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

- Any context or caveats.
- If this prompt touches provider emulation, signature verification, config/auth/PII, or the provider plugin schema contract, obtain user sign-off and record it in `DECISIONS.md` before implementing.
