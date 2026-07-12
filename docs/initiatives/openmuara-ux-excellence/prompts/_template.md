> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt NN — <Title>

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

<One-sentence outcome.>

## Why now

<Context and user pain point.>

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

- Code changes on `feat/ux-excellence`.
- Updated tests (happy path, error path, one edge case).
- Updated smoke test if routes, CLI flags, or defaults changed.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Release-notes snippet describing user-facing changes.
- Git commit with a clear message.
