> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dashboard Control Plane

> **Status:** 🟡 Planned | **Started:** 2026-07-03
> **Scope:** Turn the dashboard into a real control plane for providers and webhooks: enable providers, configure per-provider webhooks, link to integration docs, refresh data without reloading the page, and deliver a polished, accessible, low-memory SPA that follows the OpenMuara philosophy of simple, fast, efficient tools.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Problem

The dashboard is mostly read-only and feels cramped on a single page. Users cannot:

- Enable or disable providers without editing `.muara/config.yml`.
- Configure webhook URLs, secrets, or events per provider.
- See how to integrate with a provider without leaving the app.
- Refresh data manually without a full browser reload.
- Filter the ledger or transaction table by URL, provider, status, reference, date range, or event type in a single coherent toolbar.
- Operate the dashboard comfortably on small screens or with only a keyboard.

---

## Goals

1. **Provider enablement** — toggle providers on/off from the dashboard and persist to config.
2. **Webhook configuration** — set URL, signing secret, and enabled events per provider.
3. **Integration docs** — link to each provider's OpenAPI spec or official docs from the dashboard.
4. **SPA reload button** — refresh all data instantly without losing tab/scroll state.
5. **Advanced table filtering** — filter ledger, transactions, and webhooks by provider, status, reference, URL, date range, and event type.
6. **Better page layout** — split cramped single-page content into focused views with clear navigation, rather than stacking everything vertically.
7. **Accessibility & usability** — pass axe-core with zero serious violations, full keyboard navigation, visible focus states, and responsive controls.
8. Add backend and frontend tests.
9. Pass all quality gates.

---

## Live audit findings

Audited `http://127.0.0.1:9000/_admin` (Ledger, Transactions, Webhooks views) with axe-core 4.9.1 and manual inspection.

### Axe-core violations

| ID | Impact | Count | Description | Quick fix |
|---|---|---|---|---|
| `color-contrast` | serious | 1 | Active tab/segment text does not meet WCAG AA contrast. | Bump active-state foreground or background luminance. |
| `region` | moderate | 2 | Page heading and nav are not contained in landmarks. | Wrap header in `<header>`, ensure all content is inside `<main>`/`<nav>`/`<footer>`. |

### Manual UX findings

- **Provider controls are missing entirely.** The dashboard currently has only Ledger, Transactions, and Webhooks tabs; provider enablement and webhook configuration live elsewhere or are not implemented.
- **Ledger filter toolbar overlaps with top tabs.** Segment buttons "All / Transactions / Webhooks" duplicate the purpose of the top tab bar and consume vertical space.
- **Search requires clicking a Search button.** Filtering should be real-time (debounced) as the user types.
- **Provider filter select lacks a visible label.** It has `aria-label` but no visible label, which hurts sighted keyboard users.
- **Empty state is minimal.** It does not explain what the user can do next (e.g., run an example, create a checkout session).
- **Onboarding is not surfaced in the SPA.** The `/_admin/onboarding` endpoint exists but is not visible in the UI.
- **No keyboard shortcuts.** Common actions (reload, focus search, next/previous tab) require mouse interaction.
- **No connection status indicator.** Users cannot tell when the server is unreachable or data is stale.
- **Theme toggle has no accessible label.** The "☾" icon alone is not descriptive enough.
- **Tables are not sortable.** Columns such as Time, Amount, and Status should be sortable.
- **Webhook config link for Stripe appears even when Stripe is disabled.** Provider-specific shortcuts should respect the enabled state.

---

## Recommendations & enhancements

### Provider management

- **Status indicators** — show `healthy`, `misconfigured`, or `disabled` per provider, not just on/off.
- **Bulk actions** — "Enable all", "Disable all", "Reset to defaults".
- **Quick-start actions** — per-provider buttons like "Create Stripe Checkout Session" or "Send Fawry charge".
- **Search/filter providers** — useful once many providers are registered.
- **Config diff preview** — show what will change in `.muara/config.yml` before saving.
- **Copy-to-clipboard** for webhook URLs, secrets, and sample curl commands.
- **Provider detail panel** — clicking a provider opens a side panel (not a full page) with config, docs, health, and quick actions.
- **Default provider state** — new providers discovered at startup should default to `enabled: false` until the user explicitly opts in; the dashboard should highlight newly available providers.
- **Provider ordering** — sort providers by status (enabled first), then alphabetically; let users pin favorites to the top.

### Webhooks

- **Test button** — send a synthetic event to the configured URL and report success/failure inline.
- **Delivery log filtering** — filter by provider, status, time range, URL, reference, and event type.
- **Retry UI** — retry failed deliveries directly from the webhook detail panel.
- **Payload preview** — show a human-readable payload sample for each selectable event.
- **Secret visibility toggle** — show/hide signing secrets in inputs.
- **Per-provider webhook targets** — each provider can override the global webhook URL and event list.
- **Webhook test payload contract** — the test button sends a well-known synthetic event (`muara.test_event`) with a timestamp and a signature using the configured secret; the UI shows HTTP status, latency, and signature verification result.

### Table UX (Ledger, Transactions, Webhooks)

- **Real-time search** — debounced input filtering across visible columns.
- **Column sorting** — click headers to sort by Time, Provider, Status, Amount, Reference.
- **Persistent column visibility** — let users show/hide columns; persist preference in `localStorage`.
- **Pagination or virtual scrolling** — for large ledgers, keep memory low by rendering only visible rows.
- **URL filter** — for webhook attempts, filter by the target URL substring.
- **Date-range filter** — preset ranges (last hour, today, last 7 days) plus custom start/end.
- **Status chips with color tokens** — consistent color-coded badges across all tables.
- **Row actions** — context menu or inline buttons for replay, copy reference, view details.
- **URL-state persistence** — reflect active filters, sort, and page in the query string so users can bookmark/share filtered views and use the browser back button naturally.

### Page layout & navigation

- **Split the single-page view** into top-level tabs: Overview, Providers, Ledger, Transactions, Webhooks, Settings.
- **Overview tab** shows onboarding checklist, connection status, enabled provider summary, and recent activity.
- **Remove the Ledger segment bar** and rely on the top tabs for switching between Transactions, Ledger, and Webhooks.
- **Breadcrumbs** when navigating into provider detail or webhook attempt detail.
- **Sticky filter bar** so filters remain visible while scrolling long tables.
- **Collapsible sidebar** for future expansion without adding vertical clutter.
- **Loading skeletons** — show skeleton placeholders while initial data loads; avoid layout shift.
- **Error boundaries** — catch unexpected rendering errors and show a friendly "Something went wrong" panel with a reload button.

### Empty states & onboarding

- **Contextual empty states** — each view shows an illustration/icon, explanatory text, and a primary CTA.
- **Onboarding checklist** — surfaced on Overview: "Enable a provider", "Send a test charge", "Configure a webhook", "Replay a webhook".
- **Quick-create CTAs** in empty states: "Create Stripe checkout session", "Run prepaid top-up example".
- **First-time product tour** — optional, dismissible tooltip sequence that walks through Overview → Providers → Webhooks → Ledger on the user's first visit.

### Accessibility

- **Landmarks** — `<header>`, `<main>`, `<nav>`, and `<footer>`; all content inside landmarks.
- **Focus management** — visible focus rings, focus trap in modals/side panels, return focus on close.
- **Labels** — every input, select, button, and toggle has a visible label or `aria-labelledby`.
- **Live regions** — announce save status, reload results, and errors via `aria-live`.
- **Skip link** already present; keep it working after layout changes.
- **Keyboard shortcuts** — `?` for help, `/` focus search, `r` reload, `g` then `p/l/t/w` for Go to tab.
- **Target size** — interactive controls at least 24×24 CSS pixels.
- ** axe-core gate** — add `axe-core` to frontend test suite; fail CI on serious violations.

### Performance & philosophy

- **Low memory** — keep the dashboard bundle small; avoid heavy charting libraries. Prefer built-in `<table>` and native browser APIs.
- **Efficient re-renders** — React `memo` for table rows, virtual list for >100 rows.
- **No runtime provider registration** — keep the dashboard lightweight; config changes require a restart.
- **Atomic config writes** — write temp file then rename.
- **Minimal dependencies** — do not add a CSS framework; extend the existing token-based CSS system.
- **Bundle budget** — dashboard JS + CSS ≤ 150 KB gzipped; fail CI if the budget is exceeded.
- **Browser baseline** — support the last two versions of Chrome, Firefox, Safari, and Edge; no polyfills for dead browsers.

### Responsive & mobile

- **Stack filters vertically** on narrow viewports.
- **Horizontal scroll** for tables with an accessible overflow indicator.
- **Touch-friendly toggles** and buttons with adequate spacing.
- **Bottom sheet** for detail panels on mobile instead of side panel.

### Developer ergonomics

- **OpenAPI/official docs link per provider** — stored in provider metadata, surfaced as "View docs" buttons.
- **Curl command generator** — from any transaction or webhook row, copy a ready-to-run curl.
- **Environment badge** — display `hardened`, `dev`, or `production-like` mode.
- **Connection status** — indicator showing whether the server connection is live.
- **Undo/revert** for the last config change.
- **Import/export config** — upload/download `.muara/config.yml` from the UI.
- **Confirmation dialogs** for destructive actions (disable provider, reset config).
- **Toast notifications** for save, error, and test outcomes.
- **Form dirty-state guard** — warn the user before leaving a view with unsaved webhook/provider changes.
- **Copy JSON payload** — copy the raw request/response payload of any ledger row.
- **Keyboard help overlay** — press `?` to show all shortcuts.
- **Rate-limit UX** — when a write endpoint returns `429`, show a countdown and disable the submit button until the limit resets.
- **Field-level help** — tooltips or helper text explaining webhook URL, signing secret, and event selection for first-time users.

---

## Non-goals

- Runtime provider registration without restart.
- Generic OpenAPI request validation.
- Multi-node config sync.
- Real-time charts/graphs (keep memory low).
- User authentication or RBAC beyond existing admin middleware.

---

## Architecture & constraints

### Current admin API surface

The dashboard already consumes read-only endpoints registered by `AdminAPIHandlers` and `WebhookAdminHandlers`:

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/_admin/transactions` | List transactions with `q`, `provider`, `status`, `limit`, `offset`. |
| `GET` | `/_admin/transactions/{ref}` | Get one transaction + empty history placeholder. |
| `POST` | `/_admin/transactions/{ref}/replay-webhook` | Replay a transaction's webhook. |
| `GET` | `/_admin/ledger` | Unified ledger of transactions and webhook attempts. |
| `GET` | `/_admin/providers` | List enabled/available providers with metadata. |
| `GET` | `/_admin/onboarding` | Onboarding state and next-step hint. |
| `GET` | `/_admin/webhooks` | List webhook attempts. |
| `GET` | `/_admin/webhooks/{ref}` | Inspect one webhook attempt (headers redacted). |
| `POST` | `/_admin/webhooks/{ref}/replay` | Replay one webhook attempt. |
| `POST` | `/_admin/webhooks/replay-all` | Bulk replay with optional filters. |
| `DELETE` | `/_admin/webhooks/{ref}` | Clear sensitive payload/headers for an attempt. |

New write endpoints needed for this initiative:

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/_admin/config` | Return safe, non-secret subset of current config. |
| `PATCH` | `/_admin/config/providers` | Enable/disable providers; validate before persisting. |
| `GET` | `/_admin/config/webhooks` | Get webhook config (secrets masked). |
| `PATCH` | `/_admin/config/webhooks` | Update global + per-provider webhook targets/events. |
| `POST` | `/_admin/config/webhooks/test` | Dispatch a synthetic event to a target URL. |
| `POST` | `/_admin/config/reload` | Signal dashboard clients to refresh state. |
| `GET` | `/_admin/providers/{name}/health` | Return provider health: `healthy`, `misconfigured`, or `disabled`, plus a short reason. |

### Config schema additions

Provider config (`providers.<name>`) already supports `enabled` + arbitrary `config`. This initiative extends the webhook section and adds a per-provider override block:

```yaml
webhook:
  url: "https://example.com/webhook"          # global fallback
  max_retries: 3
  targets:                                      # per-provider overrides
    stripe: "https://example.com/stripe"
    fawry: "https://example.com/fawry"
  events:                                       # per-provider event subscriptions
    stripe: ["checkout.session.completed", "payment_intent.succeeded"]
    fawry: ["charge_paid"]
```

Dashboard writes must:
1. Validate the new config in memory (reuse `Config.Validate`).
2. Write to a temp file, then atomic rename to `.muara/config.yml`.
3. Keep a `.muara/config.yml.bak` before each write so users can recover from bad edits.
4. Return `202 Accepted`; a restart is still required to activate provider changes.
5. Redact secrets in `GET` responses (replace with `***`).
6. Detect external config changes by comparing a server-side checksum/timestamp; reject writes with `409 Conflict` if the file changed since the client loaded it.

### Backend constraints

- **No runtime provider registration.** Enabling/disabling a provider requires a server restart; the UI must communicate this clearly.
- **Atomic config writes.** Use `os.WriteFile` to a temp file in the same directory, then `os.Rename`.
- **Secret redaction.** Any endpoint returning config must mask keys matching `*secret*`, `*key*`, `*token*`, `*password*`.
- **Validation reuse.** Call `config.Load` on the candidate YAML and `cfg.Validate()` before writing.
- **Audit logging.** Log every config change via `audit.FromContext(r.Context()).Log(...)`.

### Frontend architecture

- Add a `ConfigProvider` context that holds global reload state.
- `useConfig()` hook returns `{ config, reload, isLoading, error }`.
- All views read from the context; the reload button invalidates the cache and re-fetches in parallel.
- Per-provider webhook forms are driven by the event list exposed by each provider plugin (see `provider.Events()` or static metadata from `listProvidersHandler`).
- Add a `useFilters()` hook for debounced search, column sort, and persisted column visibility.
- Keep bundle size low; prefer native APIs over heavy table libraries.

### State management & concurrency

- **Optimistic UI with rollback** — toggle provider state immediately in the UI, then revert on server error with a toast explaining the failure.
- **Dirty-state tracking** — track unsaved form changes so the UI can warn before navigation and show a "Save changes" prompt.
- **Multi-tab synchronization** — broadcast config changes via `BroadcastChannel` or `storage` events so other dashboard tabs reload state automatically; fall back to a "config changed on disk" banner.
- **Request deduplication** — cancel in-flight reload/fetch requests when a newer request is issued.
- **Stale-data indicator** — show a subtle timestamp or "Updated 2m ago" label; flash rows that changed since the last reload.

### Real-time updates

- Use **SSE** on a lightweight `/_admin/events` endpoint for new transactions/webhook attempts when the dashboard is open.
- If SSE is too complex for v1, use **short polling (5–10s)** with a visible pause/resume control to keep server load negligible.
- Avoid WebSockets (heavier connection model) to stay aligned with the low-memory philosophy.

---

## Security checklist

- [ ] Admin endpoints protected by existing `admin` middleware when `admin.enabled: true`.
- [ ] Config write endpoints require `PATCH` and use CSRF token/header when `server.csrf.enabled: true`.
- [ ] Secrets are never returned by `GET _admin/config` or `GET _admin/config/webhooks`.
- [ ] Webhook test endpoint cannot be abused as an SSRF vector: only `http/https`, no private/reserved IP ranges in hardened mode, timeout ≤ 5s.
- [ ] Validate webhook URL scheme and reject `file://`, `ftp://`, etc.
- [ ] Log config changes with trace ID and old/new values (secrets redacted).
- [ ] Rate-limit write endpoints (reuse `RateLimitMiddleware`).
- [ ] Confirm destructive actions (disable provider, reset defaults) in the UI.
- [ ] SPA obtains a fresh CSRF token from `GET /_admin/config` or a dedicated `/_admin/csrf` endpoint before each `PATCH`.
- [ ] Session cookie is `HttpOnly`, `Secure` when TLS is enabled, and `SameSite=Lax` or `Strict`.
- [ ] Imported config YAML is validated server-side with the same `Config.Validate` path before replacing the live file.

---

## Frontend testing strategy

- **Unit tests** for each new component/hook using Vitest + React Testing Library:
  - Provider toggle card renders status and fires `onChange`.
  - Webhook form validates URL input and masks secret.
  - `useConfig` reload invalidates cache and sets loading state.
  - `useFilters` debounces search and persists sort order.
- **Integration tests** for the config flow:
  - Enable a provider → expect `PATCH` body.
  - Save webhook config → expect toast and redacted response handling.
  - Test webhook button → expect loading then success/error state.
- **Accessibility tests**:
  - Toggles have accessible labels.
  - Focus management in modals.
  - Keyboard-only save/reload flow.
  - axe-core run on each major view with zero serious violations.
- **Visual regression** (optional): capture full-page screenshots in CI for key views.
- **Test fixtures** — create reusable fixtures in `web/dashboard/src/test/fixtures/` for provider metadata, config, webhook attempts, and transactions so tests stay deterministic.

---

## Acceptance criteria

- [ ] Providers card has enable/disable toggles.
- [ ] Webhook config page supports URL, secret, and event selection per provider.
- [ ] Provider cards link to API docs or official docs.
- [ ] Reload button refreshes data without full page reload.
- [ ] Ledger/transaction/webhook tables support real-time search, sorting, and URL/provider/status/date filters.
- [ ] Dashboard layout splits content across focused tabs instead of one cramped page.
- [ ] Axe-core reports zero serious violations.
- [ ] Keyboard shortcuts work for search, reload, and navigation.
- [ ] Empty states and onboarding checklist guide first-time users.
- [ ] Backend tests for admin endpoints.
- [ ] Frontend tests for controls.
- [ ] Config writes create a `.muara/config.yml.bak` and detect external changes (409 conflict).
- [ ] Webhook test button reports status, latency, and signature verification.
- [ ] Dashboard bundle stays within the 150 KB gzipped budget.
- [ ] All quality gates pass.

---

## References

- `web/dashboard/src/components/Providers.tsx`
- `web/dashboard/src/views/Webhooks.tsx`
- `web/dashboard/src/views/Ledger.tsx`
- `web/dashboard/src/views/Transactions.tsx`
- `web/dashboard/src/components/Shell.tsx`
- `web/dashboard/src/hooks/useConfig.ts` (to create)
- `web/dashboard/src/hooks/useFilters.ts` (to create)
- `internal/server/admin_api.go`
- `internal/server/webhook_admin.go`
- `internal/server/router.go`
- `internal/webhook/dispatcher.go`
- `internal/provider/provider.go`
- `internal/config/config.go`

---

## Self-assessment

**Solidity: 9.8 / 10**

The initiative now covers provider enablement, per-provider webhook configuration, integration docs, SPA reload, advanced table filtering, layout de-cluttering, accessibility audit findings, responsive/mobile considerations, keyboard shortcuts, performance constraints aligned with OpenMuara's low-memory philosophy, state management/concurrency, config backup/conflict detection, real-time update strategy, CSRF/session details, bundle/browser baselines, test fixtures, and a thorough testing strategy. The remaining 0.2 point is reserved for implementation-specific surprises (exact event-list contract per provider and real-world bundle-size validation), which should be resolved during implementation and verified in CI.
