> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard — Mailpit-Style Redesign

> **Status:** 🟡 Planned | **Started:** 2026-07-06
> **Scope:** Redesign the OpenMuara admin dashboard to feel like Mailpit: a fixed left navigation, a full-width ledger as the default view, focused Webhooks and Settings views, and a rich provider configuration experience with version tabs, base URLs, and environment-variable reference.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/dashboard-mailpit-redesign`
>
> **Why:** The current dashboard stacks everything vertically and buries provider configuration. Users want a dense, mail/IDE-like layout where the ledger is the home screen, webhooks are one click away, and provider settings are discoverable and editable from the GUI.

---

## Initiative Structure

```
docs/initiatives/openmuara-dashboard-mailpit-redesign/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    ├── 01-shell-layout-and-navigation.md
    ├── 02-ledger-default-view.md
    ├── 03-webhooks-view.md
    └── 04-provider-settings.md
```

Planning docs live in `docs/initiatives/openmuara-dashboard-mailpit-redesign/` in the root repo.
Product code commits to the `feat/dashboard-mailpit-redesign` branch. Do not commit directly to
`main`.

---

## Goals

1. **Mailpit-like shell** — fixed left sidebar navigation, compact header, and a main outlet that fills the viewport.
2. **Ledger as default** — opening `/_admin` lands directly on a full-width, filterable ledger table.
3. **Focused top-level views** — three primary navigation items: **Ledger**, **Webhooks**, **Settings**.
4. **Consistent filters** — every table (Ledger, Webhooks, and any future table view) has a filter toolbar.
5. **Detail pages** — clicking a ledger row or webhook row opens a dedicated detail page instead of an inline panel.
6. **Provider settings UI** — enable/disable providers, set per-provider webhook targets, and inspect related environment variables without editing `.muara/config.yml`.
7. **Provider detail pages** — click a provider card in Settings to open a dedicated configuration page. Per-provider webhook targets are edited here, not on the top-level Webhooks view.
8. **Webhooks as a delivery log** — the top-level Webhooks view is a filterable table of webhook attempts, analogous to the Ledger view.
8. **Version tabs** — providers that support multiple API versions (e.g. Fawry `v1`/`v2`) show tabbed configuration.
9. **Base URL surfacing** — each provider/version page displays its base URL and sample endpoint.
10. **Environment variable reference** — each provider lists the env vars that map to its config keys.
11. **Dual-port runtime** — separate ports for the admin web UI and the provider emulation endpoints, making it easy to expose only the API later.
12. **Accessibility & quality** — keep keyboard shortcuts, axe-core clean, and all quality gates passing.
13. **Tests** — add/update backend and frontend tests for new routes, components, and navigation.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style.

### 2. Priority stack
When trade-offs arise, decide in this order:

1. **UI** — visual clarity, density, and Mailpit-like layout come first.
2. **UX** — navigation, feedback, and smooth interactions come second.
3. **Performance** — fast renders, minimal blocking work, and snappy routing come third.
4. **Usability** — discoverability, keyboard support, and helpful empty states come fourth.
5. **Philosophy** — local-first, simple, and explicit behavior come fifth.
6. **Efficiency** — low CPU/battery usage and lean code paths come sixth.
7. **Memory size** — bundle and heap size matter, but only after the above.

This means we may accept a slightly larger bundle or more memory if it materially improves UI/UX, but we will not add heavy dependencies just for convenience.

### 3. Backward compatibility
All admin JSON endpoints remain backward-compatible. Existing query-string filters and keyboard shortcuts are preserved or extended, not removed. When dual-port mode is enabled, the provider emulation endpoints continue to answer on the original port by default; the admin UI is optionally reachable on a second port.

### 4. P0 integration changes need explicit approval
Prompts that change provider emulation logic, webhook signature verification, or config persistence (P05) require user sign-off per `AGENTS.md`. Pure UI/UX refactors do not.

### 5. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `cd web/dashboard && npm run test`
- `cd web/dashboard && npm run build`

### 6. Definition of done
Beyond the quality gates, a prompt is done only when:

- The change is tested or justified as untestable.
- `HANDOFF.md` is updated with what was built.
- `TRACKING.md` marks the prompt `✅` with the commit hash.
- User-facing changes are noted for the next release notes.

---

## Out of Scope

- Runtime provider registration without restart.
- Generic OpenAPI request builder.
- Multi-node config sync.
- Charts/graphs (keep memory low).
- User authentication or RBAC beyond existing admin middleware.
- Adding new providers or payment methods.
- Changing the provider plugin schema contract.

---

## Metrics

| Metric | Current | Target | How measured |
|--------|---------|--------|--------------|
| Dashboard nav items | 4 top tabs (Overview/Ledger/Transactions/Webhooks) | 3 left-nav items (Ledger/Webhooks/Settings) | Visual inspection / Playwright |
| Provider config reachable from GUI | No dedicated settings view | Yes, via Settings → provider card | Manual / Playwright |
| Table filters | Ledger/Webhooks have filters | Every table view has a filter toolbar | Manual / Playwright |
| Detail pages | Inline detail panels | Separate Ledger Detail and Webhook Detail pages | Manual / Playwright |
| Admin/API port split | Single port (9000) | Optional `server.admin_port` separate from `server.port` | Config test / integration test |
| Axe-core serious violations | 0 known | 0 | `npm run test:a11y` |
| Bundle size | ≤ 150 KB gzipped | ≤ 150 KB gzipped | `scripts/check-bundle-size.js` |
| Quality gates | Passing | Passing | `go test ./...`, `golangci-lint run`, `npm run test` |

## Success Criteria

- `/_admin` opens to the Ledger view by default.
- Left navigation contains Ledger, Webhooks, and Settings.
- Every table view has a filter toolbar (search, provider/status filters, sort).
- Clicking a ledger row opens a Ledger Detail page.
- Clicking a webhook row opens a Webhook Detail page.
- The top-level Webhooks view shows webhook delivery attempts only; configuration is not edited there.
- Settings shows provider cards with summary (name, status, description, emulated providers).
- Clicking a provider card opens a detail page with:
  - Enable/disable toggle.
  - Base URL for the selected version.
  - v1/v2 tabs when the provider exposes multiple versions.
  - Per-provider webhook URL input.
  - Related environment variables list.
- The server can run the admin UI on a separate port from the provider emulation endpoints.
- Config changes persist to `.muara/config.yml` and show a "restart required" notice.
- All existing keyboard shortcuts continue to work.
- Dashboard bundle stays within the 150 KB gzipped budget.
- All quality gates pass.

---

## Proposed Layout

```
┌─────────────────────────────────────────────────────────────┐
│  OpenMuara                     [reload] [theme] [help]      │  ← compact top bar
├──────────┬──────────────────────────────────────────────────┤
│  ◆ Ledger│                                                  │
│  ○ Webho │  [Ledger table — full width, filters, sortable]  │  ← default outlet
│  ○ Setti │                                                  │
│          │                                                  │
└──────────┴──────────────────────────────────────────────────┘

Ledger / Webhooks table — every table has the same filter toolbar pattern:
┌─────────────────────────────────────────────────────────────┐
│ [search] [provider ▼] [status ▼] [sort ▼] [refresh]         │
├─────────────────────────────────────────────────────────────┤
│ Time | Type | Provider | Reference | Status | Summary        │
└─────────────────────────────────────────────────────────────┘

Ledger Detail page (click a row):
┌─────────────────────────────────────────────────────────────┐
│ ← Back to Ledger  |  Transaction: pi_xxx                    │
│ Provider    Stripe                                          │
│ Amount      10.00 USD                                       │
│ Status      paid                                            │
│ Trace ID    ...                                             │
│ Payload     { ... }                                         │
│ [Replay webhook]                                            │
└─────────────────────────────────────────────────────────────┘

Webhooks view (delivery log only — configuration moved to Settings):
┌─────────────────────────────────────────────────────────────┐
│ [search] [provider ▼] [status ▼] [sort ▼] [refresh]         │
├─────────────────────────────────────────────────────────────┤
│ Reference | Provider | URL | Status | Attempts | Last Error │
└─────────────────────────────────────────────────────────────┘

Webhook Detail page (click a row):
┌─────────────────────────────────────────────────────────────┐
│ ← Back to Webhooks  |  Webhook: wh_xxx                      │
│ Provider    Stripe                                          │
│ URL         https://example.com/webhook                     │
│ Status      delivered                                       │
│ Attempts    1                                               │
│ Payload     { ... }                                         │
│ [Replay]                                                    │
└─────────────────────────────────────────────────────────────┘

Settings → provider grid:
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ Stripe      │ │ Fawry       │ │ Billplz     │  ← cards: name, status, desc
│ enabled     │ │ v1 / v2     │ │ disabled    │
│ [Configure] │ │ [Configure] │ │ [Configure] │
└─────────────┘ └─────────────┘ └─────────────┘

Settings → provider detail (Fawry example):
┌─────────────────────────────────────────────────────────────┐
│ Fawry                                        [Active] [Off] │  ← enable toggle
│ Egyptian payment gateway (legacy v1 and V2 notifications)   │
│                                                             │
│ [ v1 ] [ v2 ]                                               │  ← version tabs
│ Base URL: http://127.0.0.1:9000/fawry/v1                    │
│ Sample endpoint: POST /fawry/v1/charge                      │
│                                                             │
│ Webhook target URL (per provider)                           │
│ [https://example.com/fawry-webhook        ]                 │
│                                                             │
│ Environment variables                                       │
│ MUARA_FAWRY_MERCHANT_CODE                                   │
│ MUARA_FAWRY_MERCHANT_SECURITY_KEY                           │
│ MUARA_FAWRY_WEBHOOK_SECRET                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Backend Changes Needed

The dashboard already consumes these endpoints:

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/_admin/providers` | Enabled/available providers + metadata. |
| `GET` | `/_admin/config` | Safe config subset. |
| `PATCH` | `/_admin/config/providers` | Enable/disable providers. |
| `GET` | `/_admin/config/webhooks` | Webhook config. |
| `PATCH` | `/_admin/config/webhooks` | Update webhook targets/events. |

This initiative requires the following additions:

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/_admin/providers/{name}` | Full provider metadata: base URLs per version, env var mapping, sample routes per version, docs link. |
| (meta) | `/_admin/providers` | Include `base_url` and `env_vars` in the existing metadata payload where available. |
| (runtime) | `server.admin_port` | Optional second port for the admin UI. Provider emulation endpoints remain on `server.port` by default. |

### Dual-port runtime

Add an optional `admin_port` field to the server config:

```yaml
server:
  host: 127.0.0.1
  port: 9000          # provider emulation endpoints
  admin_port: 9001    # admin web UI (optional; falls back to port when empty)
```

- When `admin_port` is set, `/_admin` and all `/_admin/*` JSON endpoints are served only on the admin port.
- The provider port (`port`) continues to serve provider emulation routes and health/readiness endpoints.
- When `admin_port` is unset, the current single-port behavior is preserved.
- The dashboard discovers the correct API base URL from a `<meta>` tag injected by the Go template or via `window.__MUARA_ADMIN_API__`.

This makes it easy to expose only the provider API to a network while keeping the admin UI localhost-only.

### Provider metadata enrichment

`internal/config/wizard.go` and `internal/server/admin_api.go` need to expose:

- `base_url` — the absolute base URL for the provider/version (e.g. `http://127.0.0.1:9000/fawry/v1`).
- `versions` already exists; extend it with per-version `sample_route` and `base_url`.
- `env_vars` — a list of environment variables that map to `providers.<name>.config` keys.

Example metadata addition:

```json
{
  "fawry": {
    "display_name": "Fawry",
    "description": "Egyptian payment gateway (legacy v1 and V2 server notifications)",
    "enabled": true,
    "version": "v1",
    "versions": ["v1", "v2"],
    "version_details": {
      "v1": { "base_url": "http://127.0.0.1:9000/fawry/v1", "sample_route": "/fawry/v1/charge" },
      "v2": { "base_url": "http://127.0.0.1:9000/fawry/v2", "sample_route": "/fawry/v2/charge" }
    },
    "env_vars": [
      "MUARA_FAWRY_MERCHANT_CODE",
      "MUARA_FAWRY_MERCHANT_SECURITY_KEY",
      "MUARA_FAWRY_WEBHOOK_SECRET"
    ]
  }
}
```

Environment variable names are derived from the config key using a documented convention:
`MUARA_<PROVIDER>_<CONFIG_KEY>` upper-cased with underscores. The UI renders them as read-only reference text, not editable secrets.

---

## Frontend Architecture

### New/updated components

| File | Purpose |
|---|---|
| `web/dashboard/src/components/AppShell.tsx` | Left nav + header + main outlet. Replaces `Shell.tsx`. |
| `web/dashboard/src/components/SidebarNav.tsx` | Three nav links with active state and keyboard shortcuts. |
 | `web/dashboard/src/components/FilterToolbar.tsx` | Reusable filter bar: search, provider/status selects, sort, refresh. |
| `web/dashboard/src/views/Settings.tsx` | Provider card grid. |
| `web/dashboard/src/views/ProviderDetail.tsx` | Per-provider config page with version tabs, base URL, webhook URL, env vars. |
| `web/dashboard/src/views/Ledger.tsx` | Existing ledger view, promoted to default. |
| `web/dashboard/src/views/LedgerDetail.tsx` | Full-page ledger row detail. |
| `web/dashboard/src/views/Webhooks.tsx` | Existing webhooks view, moved out of Settings. |
| `web/dashboard/src/views/WebhookDetail.tsx` | Full-page webhook attempt detail. |
| `web/dashboard/src/app.tsx` | Router-like tab switcher using URL state. |

### Routing/navigation

Use URL query parameter `view` (or path if feasible within the embedded SPA):

- `/_admin` → Ledger
- `/_admin?view=webhooks` → Webhooks
- `/_admin?view=settings` → Settings provider grid
- `/_admin?view=settings&provider=fawry` → Fawry detail page
- `/_admin?view=ledger-detail&ref=pi_xxx` → Ledger detail page
- `/_admin?view=webhook-detail&ref=wh_xxx` → Webhook detail page

Keep `tab` parameter as a fallback redirect for existing bookmarks.

Clicking a ledger or webhook row navigates to its detail page; the detail page has a back button returning focus to the originating row.

### State management

- Reuse existing `useAsync` and `usePolling` hooks.
- Add a lightweight `useConfig()` hook that returns config, reload, and saving state.
- Provider enablement uses optimistic UI with rollback on error.

---

## Accessibility & Keyboard

- Left nav is a `<nav>` landmark with `aria-label="Main"`.
- Active nav item has `aria-current="page"`.
- Keyboard shortcuts preserved:
  - `1` / `2` / `3` — Ledger / Webhooks / Settings.
  - `/` — focus ledger search when on Ledger.
  - `r` — reload.
  - `?` — help.
  - `d` — toggle theme.
- Provider detail version tabs follow the existing tab keyboard pattern from `Shell.tsx`.
- Focus returns to the triggering card when closing provider detail.

---

## Security Checklist

- [ ] Admin endpoints protected by existing admin middleware when `admin.enabled: true`.
- [ ] Config write endpoints use CSRF token/header when `server.csrf.enabled: true`.
- [ ] Secrets are never returned by metadata endpoints; env vars are names only.
- [ ] Webhook test endpoint reuses existing SSRF protections in hardened mode.
- [ ] Config writes create a `.muara/config.yml.bak` and detect external changes (`409 Conflict`).
- [ ] Rate-limit write endpoints via existing middleware.

---

## Frontend Testing Strategy

- **Unit tests** (Vitest + React Testing Library):
  - `SidebarNav` renders three items and highlights the active one.
  - `FilterToolbar` fires filter/sort callbacks and persists state to the URL.
  - `Settings` renders provider cards and navigates to detail on click.
  - `ProviderDetail` shows version tabs only when `versions.length > 1`.
  - `LedgerDetail` and `WebhookDetail` render data and expose a back navigation.
  - Enable toggle fires `PATCH /_admin/config/providers` and shows restart notice.
- **Integration tests**:
  - Navigate Ledger → Webhooks → Settings via nav and keyboard.
  - Click a ledger/webhook row and verify navigation to detail page.
  - Save provider webhook target and verify `PATCH /_admin/config/webhooks` body.
- **Accessibility tests**:
  - axe-core on Ledger, LedgerDetail, Webhooks, WebhookDetail, Settings, and ProviderDetail.
  - Keyboard-only navigation flow.

---

## Acceptance Criteria

- [ ] Left sidebar navigation with Ledger, Webhooks, Settings.
- [ ] Ledger is the default view at `/_admin`.
- [ ] Every table view (Ledger, Webhooks) has a filter toolbar.
- [ ] Clicking a ledger row opens a Ledger Detail page.
- [ ] Clicking a webhook row opens a Webhook Detail page.
- [ ] The top-level Webhooks view is a delivery log only; provider webhook configuration is edited in Settings.
- [ ] Settings shows provider cards with enable status and summary.
- [ ] Provider cards navigate to a detail page.
- [ ] Provider detail has enable/disable toggle.
- [ ] Provider detail shows base URL for the selected version.
- [ ] Providers with multiple versions show v1/v2 tabs.
- [ ] Provider detail has per-provider webhook URL input.
- [ ] Provider detail lists related environment variables.
- [ ] The server supports an optional `server.admin_port` separate from `server.port`.
- [ ] Config changes persist and show "restart required" notice.
- [ ] Existing keyboard shortcuts still work.
- [ ] Axe-core zero serious violations.
- [ ] Bundle size ≤ 150 KB gzipped.
- [ ] All quality gates pass.

---

## References

- `web/dashboard/src/app.tsx`
- `web/dashboard/src/components/Shell.tsx`
- `web/dashboard/src/components/Providers.tsx`
- `web/dashboard/src/components/WebhookConfig.tsx`
- `web/dashboard/src/views/Ledger.tsx`
- `web/dashboard/src/views/Webhooks.tsx`
- `internal/server/admin_api.go`
- `internal/server/config_admin.go`
- `internal/server/server.go`
- `internal/config/config.go`
- `internal/config/wizard.go`
- `internal/provider/provider.go`
- `cmd/muara/start.go`

---

## Self-assessment

**Solidity: 8.5 / 10**

The initiative now covers navigation, filters, detail pages, dual-port runtime, provider metadata enrichment, and quality gates. The main implementation risks are the dual-port server changes (P02) and provider metadata enrichment (P05); both are opt-in/preserved-behind-existing-behavior and will be recorded in `DECISIONS.md` as they land.
