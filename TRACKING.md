# OpenMuara v1 — Tracking

## Current Status

- **Prompts:** 18/19 implemented (0 in progress, 0 not started, 0 deferred; 1 frozen for v2)
- **Tasks:** 2/2 implemented
- **Target:** Complete v1 foundation + provider ecosystem
- **v1 Philosophy:** Emulator focuses on **single charge items** — one-time payments only. Subscription/entitlement emulation (RevenueCat, mobile receipt validation) belongs to v2.

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Hard freeze — v2 only |

## Prompt Progress

| # | Prompt | Status | Notes |
|---|--------|--------|-------|
| 01 | Project Bootstrap & Rebrand | ✅ | Module path `github.com/openmuara/openmuara`, binary/CLI `muara`, workspace `.muara/`. Commits `32c53e0`–`766367d`. |
| 02 | SQLite Persistence Layer | ✅ | `engine.SQLiteStore` default at `.muara/data/ledger.db`. Commit `27b13d9`. |
| 03 | Configuration & Environment | ✅ | Viper-based loader, env prefix `MUARA_`, validation, bundled defaults. Commit `5a09e40`. |
| 04 | Provider Registry Refactor | ✅ | Declarative provider activation via `config.LoadEnabledProviders`. Commit `53db749`. |
| 05 | Core HTTP Router | ✅ | Router, middleware, idempotency, admin routes. Commit `1b2f9e6` area. |
| 06 | Fawry Provider Hardening | ✅ | Charge + escape + webhook receiver. Commit `6af3bd4` area. |
| 07 | Stripe Checkout Provider | ✅ | Checkout sessions + simulation endpoints. Commit `64bbbe2`. |
| 08a | Health & Readiness | ✅ | `/readyz` endpoint added. |
| 08b | Prometheus Metrics | ✅ | `/metrics` endpoint + request/webhook/transaction counters. |
| 09 | Audit Logging | ✅ | SQLite audit table, middleware, CLI/API list, key events logged. |
| 10 | Outgoing Webhooks | ✅ | Relay, replay, admin list. Commit `ce3cf63`–`d074bd6`. |
| 11 | Pagination, CORS & CSRF | ✅ | CORS config, CSRF double-submit cookie, admin forms protected. |
| 12 | SenangPay Provider | ✅ | Charge + callback + webhook with MD5 signature. Commit `6af3bd4`. |
| 13 | Receipt Validation Framework | ❄️ | Hard freeze — v2 only |
| 14 | RevenueCat Emulation | ❄️ | Moved to v2 initiative: `docs/initiatives/openmuara-v2-revenuecat/` |
| 15 | Docker & CI | ✅ | CI badge in README; Docker/Compose already present. |
| 16a | OpenAPI Specification | ✅ | Initial spec + live endpoint. Commit `6e99349`; see `v1-solid` S02 for drift sync. |
| 16b | Release Workflow | ✅ | `VERSION`, `CHANGELOG.md`, release workflow, cross-compile task. |
| 17 | Finalization | ✅ | Docs sweep, runbooks, and quality gates complete. |
| 18 | Migration Guide | ✅ | `docs/migration/openmuara-to-openmuara.md`, `scripts/migrate-openmuara.sh`, and `muara migrate` CLI implemented. |
| 19 | Post-Launch Monitoring | ✅ | Operations guide, on-call/debugging runbooks, risks, and known issues documented. |

## Task Progress

| # | Task | Status | Related Prompt |
|---|------|--------|----------------|
| T01 | SenangPay Signature | ✅ | 12 — MD5 signature implemented; spec aligned with code. |
| T02 | OpenMuara Migration Guide | ✅ | 18 — Migration guide, script, and CLI command implemented. |

## Decisions

- D028 ✅ Repo directory stays `muara` for v1; module path is `github.com/openmuara/openmuara`; binary and CLI commands are `muara`.
- D029 ✅ Runtime config lives at `.muara/config.yml` (user-local, ignored by git); bundled defaults embedded in `internal/config`.
- D030 ✅ Web UI shares same port under `/_admin`

## Active Initiatives

| Initiative | Status | Tracker |
|------------|--------|---------|
| OpenMuara Security Hardening | ✅ Completed / Archived | `docs/initiatives/archive/openmuara-security-hardening/TRACKING.md` |
| OpenMuara Web UI SPA | ✅ Completed / Archived | `docs/initiatives/archive/openmuara-web-ui-spa/TRACKING.md` |
| OpenMuara Documentation Website | 🟡 In Progress | `docs/initiatives/openmuara-docs-website/TRACKING.md` — branch `feat/docs-website` |
| OpenMuara MKP Fawry | ⏸️ Suspended | `docs/initiatives/openmuara-mkp-fawry/TRACKING.md` — work stashed on `feat/mkp-fawry` |
| OpenMuara UX Excellence | ✅ Completed | `docs/initiatives/openmuara-ux-excellence/TRACKING.md` — merged to `dev` |
| OpenMuara Dark Mode | ✅ Completed | `docs/initiatives/openmuara-dark-mode/TRACKING.md` — merged to `dev` |
| OpenMuara Testing Gold Standard | ✅ Completed | `docs/initiatives/openmuara-testing-gold-standard/TRACKING.md` |
| OpenMuara Stripe FPX | ✅ Superseded | `docs/initiatives/openmuara-stripe-fpx/TRACKING.md` — replaced by Checkout Sessions initiative |
| OpenMuara Stripe FPX & Card Payments | ✅ Completed | `docs/initiatives/openmuara-stripe-checkout-sessions/TRACKING.md` |
| OpenMuara Provider API Versioning | ✅ Completed | `docs/initiatives/openmuara-provider-versioning/TRACKING.md` — branch `feat/provider-versioning` |
| OpenMuara v1 Solid Gold | ✅ Completed | `docs/initiatives/openmuara-v1-solid-gold/TRACKING.md` — branch `feat/v1-solid-gold` |
| OpenMuara CLI & TUI Polish | ✅ Completed | `docs/initiatives/archive/openmuara-cli-tui-polish/README.md` |
| OpenMuara Dashboard Control Plane | ✅ Completed | `docs/initiatives/archive/openmuara-dashboard-control-plane/README.md` |
| OpenMuara Dashboard Design Refresh | ✅ Completed | `docs/initiatives/archive/openmuara-dashboard-design-refresh/README.md` |
| OpenMuara Dashboard Mailpit-Style Redesign | ✅ Completed | `docs/initiatives/openmuara-dashboard-mailpit-redesign/README.md` — merged to `dev` |
| OpenMuara Bug Hunt | 🟢 Completed | `docs/initiatives/openmuara-bug-hunt/README.md` — branch `feat/bug-hunt` |
| OpenMuara Quality Automation Follow-Up | 🟢 Completed | `docs/initiatives/openmuara-quality-automation-follow-up/README.md` |
| OpenMuara Provider Manifests | ✅ Completed | `docs/initiatives/openmuara-provider-manifests/TRACKING.md` |

## v1.1 Initiatives

| Initiative | Status | Tracker |
|------------|--------|---------|
| OpenMuara v1.1 Subscriptions | ⏸️ Suspended | `docs/initiatives/openmuara-v1-1-subscriptions/TRACKING.md` — branch `feat/v1-1-subscriptions` |

## v2 Initiatives

| Initiative | Status | Tracker |
|------------|--------|---------|
| OpenMuara v2 RevenueCat | ⏸️ Suspended | `docs/initiatives/openmuara-v2-revenuecat/TRACKING.md` — branch `feat/v2-revenuecat` |

## Archived Initiatives

| Initiative | Status | Tracker |
|------------|--------|---------|
| OpenMuara Billplz | ✅ Completed | `docs/initiatives/archive/openmuara-billplz/TRACKING.md` |
| OpenMuara ToyyibPay | ✅ Completed | `docs/initiatives/archive/openmuara-toyyibpay/TRACKING.md` |
| OpenMuara iPay88 | ✅ Completed | `docs/initiatives/archive/openmuara-ipay88/TRACKING.md` |

## Known Blockers

- D001–D027 decisions lost; backfill if history recovered.

## Next Session

Start with the highest-priority open item in `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` or land `feat/v1-solid-gold` via PR to `dev`.
