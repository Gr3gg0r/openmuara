> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P05 — Provider Settings

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** P01, P04
> **Target files:** `web/dashboard/src/views/Settings.tsx`, `web/dashboard/src/views/ProviderDetail.tsx`, `web/dashboard/src/components/WebhookConfig.tsx`; `internal/server/admin_api.go`, `internal/server/config_admin.go`, `internal/config/wizard.go`, `internal/provider/provider.go`, `web/dashboard/src/api.ts`, `web/dashboard/src/types.ts`
> **Status:** 🟡

## Goal

Build the Settings view as a provider control plane: a card grid and a rich provider detail page with enable/disable toggle, version tabs, base URL, per-provider webhook URL, and environment variable reference. Per-provider webhook targets are edited here, replacing the global webhook config in the Webhooks delivery-log view.

## Tasks

- [ ] Create `Settings.tsx` with a provider card grid showing display name, status, description, and emulated providers.
- [ ] Create `ProviderDetail.tsx` that reads `provider` from URL state and fetches `/_admin/providers/{name}`.
- [ ] Add enable/disable toggle on the detail page; persist via `PATCH /_admin/config/providers`.
- [ ] Show base URL and sample endpoint for the selected version.
- [ ] Render v1/v2 tabs when the provider reports multiple versions; switching tabs updates the visible base URL and sample endpoint.
- [ ] Add per-provider webhook URL input and save via `PATCH /_admin/config/webhooks` (reusing/moving `WebhookConfig` logic into Provider Detail).
- [ ] Derive and display related environment variable names from `MUARA_<PROVIDER>_<CONFIG_KEY>`.
- [ ] Enrich backend provider metadata with `version_details` (base URL, sample route per version) and `env_vars`.
- [ ] Show a "restart required" notice after provider enablement changes.
- [ ] Remove or hide the global `WebhookConfig` panel from the top-level Webhooks view (delivery log only).

## Acceptance Criteria

- [ ] Settings shows provider cards with summary and enable status.
- [ ] Clicking a provider card opens its detail page.
- [ ] Provider detail has an enable/disable toggle that persists to config.
- [ ] Provider detail shows base URL and sample endpoint for the selected version.
- [ ] Versioned providers show v1/v2 tabs.
- [ ] Provider detail has a per-provider webhook URL input.
- [ ] Provider detail lists related environment variables as read-only reference.
- [ ] The top-level Webhooks view no longer contains global webhook configuration UI.
- [ ] Axe-core reports zero serious violations on Settings and ProviderDetail.

## Quality Gates

Run before committing:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
```

## Notes

- This prompt changes provider metadata and config persistence; treat it as a P0 integration change per `AGENTS.md` and get user sign-off before implementation.
- Keep env vars as names only; never expose secret values in the UI.
