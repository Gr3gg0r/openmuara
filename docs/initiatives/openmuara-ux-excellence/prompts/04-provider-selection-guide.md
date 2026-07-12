> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P04 — Provider Selection Guide

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** P02

---

## Goal

Help users choose and configure the right OpenMuara provider for their real payment gateway.

## Why now

Users currently need to read `docs/providers.md` to map Stripe/Fawry/Billplz/etc. to OpenMuara provider names and routes. This should be discoverable inside the dashboard.

## Scope

### In scope

- Enrich `GET /_admin/providers` with per-provider metadata using a stable schema:
  - `name` — provider key.
  - `description` — one-line summary.
  - `real_providers` — list of real providers it emulates (e.g., `Stripe`, `Stripe Checkout`, `Stripe PaymentIntents`).
  - `sample_route` — first API route to try.
  - `sample_method` — HTTP method for the sample route.
  - `docs_path` — link to provider doc.
  - `category` — e.g., `card`, `redirect`, `wallet`, `regional`.
  - `is_recommended_for_first_time` — boolean.
- Update `internal/ui/index.html` providers grid to show this metadata and a "How to use" snippet.
- Add a short provider doc per provider under `docs/providers/<name>.md` if missing.
- Update the top-level `docs/providers.md` index.
- Add a **Contributor checklist** to `docs/providers.md` explaining how to register a new provider so it appears in the wizard, dashboard, and provider guide.

### Out of scope

- Auto-detecting the user's SDK.
- Provider logos or branding assets.

## Acceptance criteria

- [ ] `/_admin/providers` includes provider metadata following the stable schema.
- [ ] Dashboard shows description, emulated providers, category, and a sample route for each provider.
- [ ] Recommended first-time provider is visually highlighted.
- [ ] Provider-specific docs exist for all built-in providers.
- [ ] Contributor checklist exists in `docs/providers.md`.
- [ ] Tests cover the enriched endpoint.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Metadata can live in a small map in `internal/provider` keyed by provider name.
- The provider plugin manifests (`plugins/*/gateway.yml`) are a natural future home for this; for now, keep it simple and in-code.

## Deliverables

- Code changes on `feat/ux-excellence`.
- New/updated provider docs.
- Updated tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
