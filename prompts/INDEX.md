# Prompt Index — OpenMuara v1

This directory contains the execution prompts for building OpenMuara v1. Each prompt is designed to
be self-contained and actionable by a single agent session.

## Phase 1 — Foundation

| # | Prompt | Focus | Status |
|---|--------|-------|--------|
| 01 | [Project Bootstrap](01-project-bootstrap.md) | Rebrand repo, directory conventions, config path migration | ✅ |
| 02 | [SQLite Persistence Layer](02-sqlite-persistence.md) | Schema, migrations, repository pattern | ✅ |
| 03 | [Configuration & Environment](03-configuration-environment.md) | `.muara/config.yml`, env vars, validation | ✅ |
| 04 | [Provider Registry Refactor](04-provider-registry-refactor.md) | Unify provider and plugin runtime models | ✅ |
| 05 | [Core HTTP Router](05-core-http-router.md) | Provider endpoint matching, middleware, request IDs, `/readyz` | ✅ |
| 06 | [Fawry Provider Hardening](06-fawry-provider-hardening.md) | SHA256 signature, charge, webhook, escape flow | ✅ |
| 07 | [Stripe Checkout Provider](07-stripe-checkout-provider.md) | Sessions, signatures, success simulation | ✅ |

## Phase 2 — Operations & Observability

| # | Prompt | Focus | Status |
|---|--------|-------|--------|
| 08a | [Health & Readiness](08a-health-readiness.md) | `/healthz`, `/readyz`, liveness/readiness probes | ✅ |
| 08b | [Prometheus Metrics](08b-prometheus-metrics.md) | Metric names, instrumentation, `/metrics` endpoint | ✅ |
| 09 | [Audit Logging](09-audit-logging.md) | Schema, middleware, CLI/API query, retention | ✅ |
| 10 | [Outgoing Webhooks](10-outgoing-webhooks.md) | Dispatcher, retries, delivery ledger, replay | ✅ |
| 11 | [Pagination, CORS & CSRF](11-pagination-cors-csrf.md) | API pagination, CORS config, double-submit cookie | ✅ |

## Phase 3 — Provider Ecosystem

| # | Prompt | Focus | Status |
|---|--------|-------|--------|
| 12 | [SenangPay Provider](12-senangpay-provider.md) | MD5 signature, charge/callback/webhook emulation | ✅ |
| 13 | [Receipt Validation Framework](13-receipt-validation-framework.md) | App Store / Play Store receipt emulation | ❄️ |
| 14 | [RevenueCat Emulation](14-revenuecat-emulation.md) | Subscriber status, entitlements, webhooks | ❄️ |
| 15 | [Docker & CI](15-docker-ci.md) | Dockerfile, compose, GitHub Actions workflow | ✅ |

## Phase 4 — Documentation & Release

| # | Prompt | Focus | Status |
|---|--------|-------|--------|
| 16a | [OpenAPI Specification](16a-openapi-spec.md) | Public API spec generation and validation | ✅ |
| 16b | [Release Workflow](16b-release-workflow.md) | Versioning, tags, changelog, artifacts | ✅ |
| 17 | [Finalization](17-finalization.md) | Docs checklist, quality gates, handoff | ✅ |
| 18 | [Migration Guide](18-migration-guide.md) | `muara` → `openmuara` migration | ✅ |
| 19 | [Post-Launch Monitoring](19-post-launch-monitoring.md) | Runbooks, alerting, known issues | ✅ |

## Legend

- ⬜ Not started
- 🟡 In progress / partial
- ✅ Done
- ❌ Blocked
- ⏸️ Deferred
- ❄️ Frozen for v2

## Cross-references

- Task specs: [`tasks/INDEX.md`](../tasks/INDEX.md)
- Decisions: [`DECISIONS.md`](../DECISIONS.md)
- Progress tracker: [`TRACKING.md`](../TRACKING.md)
