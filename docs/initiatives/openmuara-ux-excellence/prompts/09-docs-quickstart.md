> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P09 — Quick-Start Documentation

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** P01, P04, P07, P08

---

## Goal

Create a single quick-start page that gets a new user from zero to their first successful charge and webhook in under five minutes.

## Why now

Documentation is spread across `runbooks/local-development.md`, `docs/providers.md`, `docs/operations.md`, and per-provider prompts. A new user has to piece together the path.

## Scope

### In scope

- Create `docs/quickstart.md` with per-audience paths:
  - **Developer:** Install/build → `muara init` wizard → start server → send first charge with `curl` → open `/_admin/ledger` and confirm transaction + webhook.
  - **AI Agent:** Install/build → `muara init --defaults` → `muara start --quiet` → query `/_admin/ledger` or `muara --json` commands.
  - **Tester:** Start server → open `/_admin/ledger` → trigger a payment → inspect payload, signature status, and replay webhook.
  - **Contributor:** Read provider checklist → add provider metadata → verify it appears in wizard, dashboard, and provider guide.
- Add per-provider copy-paste examples for at least Stripe, Fawry, and Billplz.
- Update root `README.md` to link prominently to `docs/quickstart.md`.
- Keep the doc under 250 lines.

### Out of scope

- Video or screenshot guides.
- Provider-specific deep dives.

## Acceptance criteria

- [ ] `docs/quickstart.md` exists and is under 250 lines.
- [ ] README links to quickstart in the first section.
- [ ] Each provider example is copy-paste runnable against a default `muara init` config.
- [ ] Each persona path is represented in the doc.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use the existing smoke-test request shapes as the canonical examples.
- Link to `docs/providers/<name>.md` for advanced configuration.

## Deliverables

- Docs changes on `feat/ux-excellence`.
- Updated `README.md`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
