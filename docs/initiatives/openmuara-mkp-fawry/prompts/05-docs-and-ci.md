> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 05 — Docs and CI

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** Prompts 01–04

---

## Goal

Update Fawry provider documentation, the MKP requirements doc, and CI coverage
so MKP engineers can adopt the emulator.

## Why now

Features are only adoptable once they are documented. This prompt closes the
loop for MKP's migration.

## Scope

### In scope

- Update `docs/providers/fawry.md` with:
  - New `OrderStatus` values.
  - `response_delay_ms` config option.
  - `billing_type` charge field and payload differences.
  - Copy-paste examples for MKP.
- Update `docs/mkp-billing-requirements.md` to mark Fawry gaps as resolved.
- Update `runbooks/local-development.md` with MKP Fawry testing notes if
  needed.
- Add or extend smoke-test coverage for canceled and expired flows.
- Update `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` if
  appropriate.

### Out of scope

- Stripe or RevenueCat documentation.
- Changes to the contribution guide.

## Acceptance criteria

- [ ] `docs/providers/fawry.md` explains all new options with examples.
- [ ] `docs/mkp-billing-requirements.md` reflects that Fawry gaps are closed.
- [ ] Smoke test exercises at least one new status or billing type.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Keep examples runnable against `http://127.0.0.1:9000`.
- Use `curl` snippets so MKP engineers can copy them into CI.

## Deliverables

- Documentation changes on `feat/mkp-fawry`.
- Updated smoke test if needed.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
