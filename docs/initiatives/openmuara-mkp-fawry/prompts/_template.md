> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt NN — <Title>

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** —

---

## Goal

<One-sentence outcome.>

## Why now

<Context and gap this closes for MKP.>

## Scope

### In scope

- <item>

### Out of scope

- <item>

## Acceptance criteria

- [ ] <criterion>
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- <hint>

## Deliverables

- Code changes on `feat/mkp-fawry`.
- Updated tests (happy path, error path, one edge case).
- Updated smoke test if routes, CLI flags, or defaults changed.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Updated provider docs if user-facing behavior changed.
- Git commit with a clear message.
