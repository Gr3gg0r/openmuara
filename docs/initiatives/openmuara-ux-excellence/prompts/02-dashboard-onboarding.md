> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P02 — Dashboard Onboarding Checklist

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** P01

---

## Goal

Show a getting-started checklist on `/_admin` that guides new users from zero to a verified webhook.

## Why now

The dashboard currently opens to empty tables. After P08 the default view will be the ledger, but users still need guidance on whether the server is ready, which provider is active, and what request to send first.

## Scope

### In scope

- Add a `GET /_admin/onboarding` endpoint that returns checklist state derived from existing data:
  - `server_ready` — always true when the endpoint responds.
  - `providers_enabled` — at least one provider is enabled.
  - `first_transaction` — at least one transaction exists.
  - `first_webhook_received` — at least one delivered/attempted webhook exists.
- Render the checklist at the top of `internal/ui/index.html`, above the ledger view.
- Show a short provider-specific "next step" based on the active provider (e.g., "Send POST /fawry/charge" or "Create a Checkout Session with POST /v1/checkout/sessions").
- Keep the checklist stateless; derive it from transactions and webhooks.
- Collapse or hide the checklist once all items are complete, with a way to expand it again.
- Handle API errors gracefully (e.g., show a retry message if `/_admin/onboarding` fails).

### Out of scope

- Persisting onboarding progress in a separate file.
- Tooltips, tours, or animations.

## Acceptance criteria

- [ ] `GET /_admin/onboarding` returns a JSON checklist.
- [ ] Dashboard renders the checklist and updates it live.
- [ ] Provider-specific next step is shown.
- [ ] Tests cover the new endpoint.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Reuse existing stores (`internal/engine`, `internal/webhook`).
- The provider-specific hint can be a static map in the admin handler keyed by active provider.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated `internal/server/admin_api_test.go` and dashboard tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
