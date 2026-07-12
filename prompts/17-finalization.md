# Prompt 17 — Finalization

## Goal
Complete documentation and quality gates before v1 handoff.

## Acceptance Criteria
- [ ] All prompts complete
- [ ] `README.md` covers install, config, run, test
- [ ] `docs/` has architecture and provider guides
- [ ] `runbooks/` has common operations
- [ ] Quality gates pass:
  - `go build ./...`
  - `go test ./...`
  - `go vet ./...`
  - `golangci-lint run`
  - smoke test script
- [ ] `AGENTS.md` and `DECISIONS.md` are current
- [ ] `TRACKING.md` shows 19/19 done

## Files to Create/Change
- `README.md`
- `docs/architecture.md`
- `runbooks/*.md`
- `TRACKING.md`

## Response Shape
Return:
1. Documentation checklist
2. Quality gate command output
3. Known issues list

## Test Notes
- Run full quality gate suite
- Run smoke test end-to-end
