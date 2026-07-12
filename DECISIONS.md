# Decisions

This file records architectural and product decisions for OpenMuara. Decisions are numbered sequentially (`D001`, `D002`, ...).

> **Note:** Earlier decisions (D001–D027) were lost in a container restart. They should be backfilled from project history if available. New decisions continue from D028.

## Decisions

### D028 — Local repo directory stays `toyol` until organization transfer

- **Status:** Decided
- **Context:** The project identity is now `OpenMuara` (`openmuara`). The module path, binary name, and user-facing config directory have already changed. The local repository directory name remains `toyol` temporarily to avoid breaking existing clones, worktrees, and local tooling. The repository will be transferred to the `openmuara` GitHub organization, at which point the directory name can be updated.
- **Decision:** Keep the local directory name as `toyol` until the GitHub organization transfer. All tracked content must use `OpenMuara` / `openmuara` / `muara` branding.
- **Consequences:**
  - Module path is `github.com/openmuara/openmuara`.
  - Binary and CLI commands are `muara`.
  - User workspace directory is `.muara/`.
  - The local directory name is the only remaining `toyol` reference and is not committed content.

### D029 — Runtime config at `.muara/config.yml`; example in `muara.yml.example`

- **Status:** Decided
- **Context:** Need a user config location without committing secrets. An example file documents the schema without risking leaked credentials.
- **Decision:**
  - Runtime config lives at `.muara/config.yml` (user-local, ignored by git).
  - Committed example/template at `muara.yml.example`.
  - Bundled default config is embedded in `internal/config.DefaultYAML()` and written by `muara init`.
- **Consequences:**
  - Users can copy `muara.yml.example` to `.muara/config.yml` or run `muara init`.
  - No secrets committed.
  - Defaults stay in sync with code.

### D030 — Web UI shares the same port under `/_admin`

- **Status:** Decided
- **Context:** OpenMuara can expose a web dashboard. Option A: separate port for UI. Option B: same port under `/_admin`.
- **Decision:** Web UI is served on the same HTTP port as the API, under the `/_admin/` path prefix.
- **Consequences:**
  - Simpler deployment (single port).
  - Admin routes must be protected (CSRF, future auth).
  - Static assets served from `web/`.

### D031 — Persistence stays SQLite-first for v1

- **Status:** Decided
- **Context:** OpenMuara needs durable local state without adding external services.
- **Decision:** v1 uses SQLite as the default persistence layer. An in-memory store remains available for tests.
- **Consequences:**
  - Transaction and audit data live in `.muara/data/ledger.db` by default.
  - Single-writer concurrency limits are acceptable for a local-first tool.
  - Future versions may add pluggable backends.

### D032 — Provider registry is declarative

- **Status:** Decided
- **Context:** Providers need to be enabled and configured without code changes.
- **Decision:** Providers are activated through `providers.<name>.enabled` in `.muara/config.yml` and loaded by `config.LoadEnabledProviders`.
- **Consequences:**
  - New providers can be registered at compile time and enabled via config.
  - Provider configs are cloned before `Init` to prevent accidental mutation.

### D033 — CSRF for admin UI uses double-submit cookie

- **Status:** Decided
- **Context:** The embedded admin dashboard issues state-changing requests and needs CSRF protection.
- **Decision:** Use a non-`HttpOnly` cookie `openmuara_csrf` plus `X-CSRF-Token` header / `csrf_token` form field. Server-to-server `/_admin/webhook-receiver` is exempt.
- **Consequences:**
  - Dashboard reads the token from a meta tag.
  - CSRF can be disabled via config if running in an isolated environment.

### D034 — Prometheus metrics use the `openmuara_` prefix

- **Status:** Decided
- **Context:** Metrics must be namespaced to avoid collision with application metrics.
- **Decision:** All Prometheus metrics are prefixed with `openmuara_` and expose provider/method/path/status labels.
- **Consequences:**
  - Scraping configs and alerts can filter by prefix.
  - Label cardinality is bounded by route count and provider list.

### D035 — Transaction and audit stores share one SQLite connection

- **Status:** Decided
- **Context:** Opening two separate SQLite connections to the same database caused `database is locked` errors at startup.
- **Decision:** `internal/cli/start.go` opens a single `*sql.DB` and passes it to both `engine.NewSQLiteStoreFromDB` and `audit.NewSQLiteStoreFromDB`.
- **Consequences:**
  - Concurrent writes are serialized through one connection, eliminating startup contention.
  - Long-running writes still block; this is acceptable for a local emulator.

### D036 — Trace-ID propagation and optional pprof for v1 debuggability

- **Status:** Decided
- **Context:** Debugging provider-to-webhook flows requires correlating incoming requests with outgoing webhooks. Go pprof is useful for local profiling but should not be exposed by default.
- **Decision:**
  - Every incoming request gets a `trace_id` via `RequestIDMiddleware` and the `X-Trace-Id` header.
  - The trace ID is stored on `engine.Transaction` and `webhook.Attempt`.
  - Outgoing webhooks include the same trace ID in the `X-Trace-Id` header.
  - `muara transaction inspect <ref>` and dashboard detail panels expose the trace ID.
  - Go pprof endpoints are mounted under `/_admin/debug/pprof/*` only when `server.pprof: true`.
- **Consequences:**
  - Users can correlate a provider request with its webhook using a single ID.
  - The SQLite schema gains an additive `trace_id` column; existing databases migrate automatically.
  - pprof remains off by default, avoiding accidental exposure.

### D037 — v1 focuses on single charge items; RevenueCat moved to v2

- **Status:** Decided
- **Context:** OpenMuara v1 has been growing provider coverage. RevenueCat emulation (subscriptions, offerings, entitlements, mobile receipts) was originally planned as Prompt 14 for v1, but it represents a different product surface from one-time payment emulation.
- **Decision:**
  - v1's scope is **single charge item emulation** only: one-time payments through providers such as Stripe Checkout, Fawry, SenangPay, iPay88, Billplz, and ToyyibPay.
  - Subscription, entitlement, and mobile receipt validation are explicitly out of scope for v1.
  - RevenueCat emulation is moved to v2 under `docs/initiatives/openmuara-v2-revenuecat/`.
- **Consequences:**
  - v1 stays focused and shippable.
  - No RevenueCat endpoints, state, or migrations are added to v1.
  - v2 will introduce subscription and purchase emulation when development begins.

## Backlog / Proposed

| ID | Topic | Status |
|----|-------|--------|

## How to Add a Decision

1. Use the next sequential ID.
2. Include status, context, decision, and consequences.
3. Update this file in the same commit as the code change, or in a dedicated docs commit.
