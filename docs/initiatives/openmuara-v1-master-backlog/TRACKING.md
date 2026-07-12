> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Master Backlog

> **Updated:** 2026-07-06 | **Status:** 🟡 Active
>
> **Scope:** Consolidated, priority-ranked view of all v1 work.
> **AI Agent:** Update this file after every product-code change.
> **Product Branch:** `dev`
> **Last Agent Action:** Verified VAL01 (client transaction validation endpoints) and closed it in the master backlog.
> **Next Agent Action:** Payment-gateway page UI/UX refresh completed. Next live options are OpenMuara Documentation Website, deferred `muara provider validate`, or a user-picked item.

---

## Legend

| Icon | Meaning |
|------|---------|
| ⬜ | To Do |
| 🟡 | In Progress |
| ✅ | Completed |
| ❌ | Blocked |
| ⏸️ | Deferred |
| ❄️ | Frozen for v2 |
| 🔀 | Parallel Safe (not used in this tracker) |

---

## Priority Rules

| Priority | Rule |
|----------|------|
| **High** | Regression, API contract break, core runtime gap, or daily-use blocker. |
| **Medium** | Provider hardening, observability, packaging, or docs that improve solid v1. |
| **Low** | Historical/completed work, deferred items, or explicitly frozen for v2. |

---

## How to Read This Table

1. Start at the top of the **High** lane.
2. Click or open the **Entry Point** — it is the executable prompt or task spec for that item.
3. After finishing, update **Status** and **Notes**, then update the source tracker.

---

## Master Backlog

| ID | Title | Source | Priority | Status | Depends On | Owner | Entry Point | Notes |
|----|-------|--------|----------|--------|------------|-------|-------------|-------|
| **S01** | Fix admin dashboard for paginated responses | v1-solid | High | ✅ | — | AI Agent | `docs/initiatives/openmuara-v1-solid/prompts/01-fix-admin-dashboard-pagination.md` | Regression from pagination change. Target: `internal/ui/index.html`, tests. Commit `c55fd81`. |
| **S02** | Sync OpenAPI spec with current API | v1-solid | High | ✅ | — | AI Agent | `docs/initiatives/openmuara-v1-solid/prompts/02-sync-openapi-spec.md` | Added `/readyz`, paginated envelopes, 409 responses. Commit `e197dd4`. |
| **S03** | Apply state machine to Stripe simulation | v1-solid | High | ✅ | — | AI Agent | `docs/initiatives/openmuara-v1-solid/prompts/03-apply-state-machine-to-stripe.md` | Replaced direct `tx.Status` with `engine.Transition`; 409 on invalid transition. Commit `807a300`. |
| **S04** | Fawry escape updates ledger + webhook signature verification | v1-solid | High | ✅ | S03 | AI Agent | `docs/initiatives/openmuara-v1-solid/prompts/04-fawry-escape-ledger-and-webhook-sigs.md` + `docs/initiatives/openmuara-v1-solid/tasks/04-fawry-escape-ledger-and-webhook-sigs.md` | Escape updates ledger; webhook verifies HMAC signature. Commit `8f099cc`. |
| **P01** | Project bootstrap & rebrand (`muara` → `openmuara`) | root | Low | ✅ | — | AI Agent | `<repo-root>/prompts/01-project-bootstrap.md` | Module path, binary/CLI `muara`, workspace `.muara/` all implemented. |
| **S05** | Improve dispatcher wiring + update runbooks | v1-solid | Medium | ✅ | S04 | AI Agent | `docs/initiatives/openmuara-v1-solid/prompts/05-dispatcher-wiring-and-runbooks.md` | Per-provider dispatcher map; replay uses correct provider. Commit `bf79dea`. |
| **P04** | Provider registry refactor | root | Medium | ✅ | P03 | AI Agent | `<repo-root>/prompts/04-provider-registry-refactor.md` | Declarative provider activation implemented. |
| **P08b** | Prometheus metrics | root | Medium | ✅ | P05 | AI Agent | `<repo-root>/prompts/08b-prometheus-metrics.md` | `/metrics` + request/webhook/transaction counters. |
| **P09** | Audit logging | root | Medium | ✅ | P05 | AI Agent | `<repo-root>/prompts/09-audit-logging.md` | SQLite audit table, middleware, CLI/API, event logging. |
| **P11** | Pagination, CORS & CSRF | root | Medium | ✅ | P05, P10 | AI Agent | `<repo-root>/prompts/11-pagination-cors-csrf.md` | CORS config + CSRF double-submit cookie. |
| **SEC01** | Security Hardening | security-hardening | Medium | ✅ | P11 | AI Agent | `docs/initiatives/archive/openmuara-security-hardening/README.md` | Admin auth, TLS, rate limiting, security headers, audit logging, CLI helpers, CI gates. Merged to `dev`. |
| **WEB01** | Web UI SPA | web-ui-spa | Medium | ✅ | P11 | AI Agent | `docs/initiatives/archive/openmuara-web-ui-spa/README.md` | Dashboard migrated to Vite + Preact SPA embedded in Go binary. Escape/pay pages remain server-rendered. Bundle-size budget, UI tests, CI jobs, a11y/error states, CSP compatibility. Merged to `dev`. Future: dark mode, offline caching, visual regression tests. |
| **P15** | Docker & CI | root | Medium | ✅ | P05 | AI Agent | `<repo-root>/prompts/15-docker-ci.md` | CI badge + existing Docker/Compose support. |
| **P16b** | Release workflow | root | Medium | ✅ | P15, P17 | AI Agent | `<repo-root>/prompts/16b-release-workflow.md` | `VERSION`, `CHANGELOG.md`, tag-triggered release workflow. |
| **P17** | Finalization | root | Medium | ✅ | All core | AI Agent | `<repo-root>/prompts/17-finalization.md` | Docs sweep, runbooks, quality gates, and trackers updated. |
| **T01** | SenangPay signature spec | root | Low | ✅ | P12 | AI Agent | `<repo-root>/tasks/senangpay-signature.md` | MD5 signature implemented; spec aligned with code. |
| **P02** | SQLite persistence layer | root | Low | ✅ | — | AI Agent | `<repo-root>/prompts/02-sqlite-persistence.md` | Completed (`27b13d9` in project tracker). |
| **P03** | Configuration & environment | root | Low | ✅ | P02 | AI Agent | `<repo-root>/prompts/03-configuration-environment.md` | Config loader + validation implemented. |
| **P05** | Core HTTP router | root | Low | ✅ | P03 | AI Agent | `<repo-root>/prompts/05-core-http-router.md` | Router, middleware, idempotency implemented. |
| **P06** | Fawry provider hardening | root | Low | ✅ | P05 | AI Agent | `<repo-root>/prompts/06-fawry-provider-hardening.md` | Charge + escape + webhook receiver exist. |
| **P07** | Stripe Checkout provider | root | Low | ✅ | P05 | AI Agent | `<repo-root>/prompts/07-stripe-checkout-provider.md` | Checkout sessions + simulation endpoints. |
| **P08a** | Health & readiness | root | Low | ✅ | P05 | AI Agent | `<repo-root>/prompts/08a-health-readiness.md` | `/readyz` added. |
| **P10** | Outgoing webhooks | root | Low | ✅ | P06 | AI Agent | `<repo-root>/prompts/10-outgoing-webhooks.md` | Relay, replay, admin list implemented. |
| **P12** | SenangPay provider | root | Low | ✅ | P05 | AI Agent | `<repo-root>/prompts/12-senangpay-provider.md` | Charge + callback + webhook with MD5 signature. |
| **P13** | Receipt validation framework | root | Low | ❄️ | — | — | `<repo-root>/prompts/13-receipt-validation-framework.md` | Hard frozen for v2 (App Store / Play Store receipts). |
| **P14** | RevenueCat emulation | root | Low | ❄️ | — | — | `<repo-root>/prompts/14-revenuecat-emulation.md` | Moved to v2 initiative: `docs/initiatives/openmuara-v2-revenuecat/`. |
| **P18** | Migration guide | root | Low | ✅ | P17 | AI Agent | `<repo-root>/prompts/18-migration-guide.md` | `docs/migration/openmuara-to-openmuara.md`, `scripts/migrate-openmuara.sh`, `muara migrate` implemented. |
| **P19** | Post-launch monitoring | root | Low | ✅ | P17 | Human | `<repo-root>/prompts/19-post-launch-monitoring.md` | `docs/operations.md`, `runbooks/on-call.md`, `runbooks/debugging.md`, risks, and known issues complete. |
| **T02** | OpenMuara-to-OpenMuara migration guide | root | Low | ✅ | P18 | AI Agent | `<repo-root>/tasks/openmuara-migration-guide.md` | Migration guide, script, and CLI command implemented. |
| **VAL01** | Client transaction validation endpoints | root | High | ✅ | — | AI Agent | `internal/fawry/status.go`, `internal/senangpay/status.go` | Signed payment-status query endpoints implemented and tested: Fawry `GET /fawry/payment-status` (`6e0abe3`), SenangPay `GET /senangpay/query` (`463563f`). OpenAPI spec updated. |
| **MKP01** | MKP Fawry integration gaps | mkp-fawry | High | ⬜ | — | AI Agent | `docs/initiatives/openmuara-mkp-fawry/README.md` | New initiative tracking Fawry gaps for MKP v2: extended states, response delay, billing type. Payment-status endpoint added under VAL01 addresses reference-to-webhook correlation. |
| **DASH01** | Dashboard Mailpit-style redesign | dashboard-mailpit-redesign | Medium | ✅ | — | AI Agent | `docs/initiatives/openmuara-dashboard-mailpit-redesign/README.md` | Merged to `dev` (commits a501b1d–5ad83c9). Left navigation, Ledger default, Webhooks log, Settings with provider cards, version tabs, base URLs, env var reference. |
| **SG01** | v1 solid gold polish | v1-solid-gold | Medium | ✅ | — | AI Agent | `docs/initiatives/openmuara-v1-solid-gold/README.md` | v1 hygiene, testing, debuggability, and usability polish. Merged to `dev` at `42ae9f1`. |
| **REL01** | CI & release audit | readiness-ci-release | Medium | ✅ | P15, P16b | AI Agent | `docs/initiatives/openmuara-readiness-ci-release-audit/README.md` | SLSA provenance, cosign signing, SBOMs, verified install script, `muara health`, hardened Docker image, full `task quality` pass. |
| **DOCS01** | Documentation completeness audit | readiness-docs-completeness | Medium | ✅ | — | AI Agent | `docs/initiatives/openmuara-readiness-docs-completeness-audit/README.md` | Gold-standard OSS docs delivered: accuracy sweep, provider doc hardening, CLI reference, governance, website discoverability, verification gates. Includes small `internal/provider/simple` secret-resolution fix so SenangPay example works. |

---

## Execution Rules

1. Start at the top of the **High** priority lane and work down.
2. Open the **Entry Point** for the item; that prompt/task is the executable artifact.
3. After completing any item:
   - Update its status and commit hash in this file.
   - Update the source tracker (`root`, `project`, or `v1-solid`).
   - Update `HANDOFF.md` in this initiative.
4. Do **not** start Low-priority historical/frozen items unless explicitly asked.

---

## Quality Gates

Quality gates are run against **product code** in the relevant source initiative, not in this meta tracker. Each `Entry Point` lists the required gate commands. When this backlog is updated, verify:

- No absolute filesystem paths were introduced.
- All cross-referenced tracker files still exist.

---

## Cross-Reference Map

| Tracker | Path | What It Contains |
|---------|------|------------------|
| This backlog | `docs/initiatives/openmuara-v1-master-backlog/TRACKING.md` | Consolidated priority view |
| Original prompt tracker | `<repo-root>/TRACKING.md` | Prompts 01–19, tasks T01–T02 |
| v1 execution tracker | `<repo-root>/docs/projects/openmuara-v1/TRACKING.md` | Phases, commit hashes, gate results |
| v1-solid initiative | `<repo-root>/docs/initiatives/openmuara-v1-solid/TRACKING.md` | Active regression-fix prompts S01–S05 |
