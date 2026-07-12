> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1 — Decision Log

> **Purpose:** Running record of decisions made during this project.

---

## Formal Decisions

### D001 — Product name: OpenMuara

- **Date:** 2026-06-26 (recovered from previous session)
- **Context:** The project started as `muara`, a Fawry emulator. The vision expanded to a universal multi-provider payment emulator. The old name no longer fit.
- **Options Considered:**
  - `muara` (keep) — pros: no rename cost; cons: too narrow, sounds like a toy.
  - `Muara` alone — pros: elegant, short; cons: too abstract, conflicts with place names.
  - `OpenMuara` — pros: signals open-source, distinct from geography, matches "OpenRouter" positioning; cons: longer, does not say "payments" without tagline.
- **Decision:** Rename to **OpenMuara** with tagline *One local API for every payment provider.*
- **Consequences:** All module paths, binary names, CLI commands, config paths, docs, and scripts must be updated.
- **Reversible?** No — once published, reverting the name is costly.
- **Reversal Trigger:** N/A
- **Logged By:** AI Agent (recovered from session)

### D005 — CLI binary name: `muara`

- **Date:** 2026-06-27
- **Context:** The full project name is `OpenMuara`, but `openmuara` is long to type as a CLI command.
- **Options Considered:**
  - `openmuara` — matches project name; 9 characters.
  - `muara` — shorter, punchy, easier to type; 5 characters.
- **Decision:** Use **`muara`** as the CLI binary name while keeping the project and module name as `OpenMuara` / `github.com/openmuara/openmuara`.
- **Consequences:** Users type `muara start`, `muara webhook list`, etc. Docs and scripts reference `muara`. The module path remains `github.com/openmuara/openmuara`.
- **Reversible?** Yes — can add `openmuara` as an alias later.
- **Reversal Trigger:** If branding or distribution requires exact name match.
- **Logged By:** AI Agent

### D002 — Implementation language: Go

- **Date:** 2026-06-26 (recovered from previous session)
- **Context:** Need a lightweight, single-binary local developer tool.
- **Options Considered:**
  - Go — pros: single static binary, fast startup, low memory, built-in HTTP; cons: less ecosystem than Node for some tasks.
  - Node.js — pros: huge ecosystem; cons: runtime dependency, larger memory footprint.
  - Python — pros: rapid prototyping; cons: interpreter dependency, packaging friction.
- **Decision:** Stay with **Go**.
- **Consequences:** All v1 features implemented in Go 1.22+.
- **Reversible?** No — rewriting the runtime is not viable.
- **Reversal Trigger:** N/A
- **Logged By:** AI Agent (recovered from session)

### D003 — Operating model: local emulator, not live router

- **Date:** 2026-06-26 (recovered from previous session)
- **Context:** Ambiguity about whether OpenMuara would proxy traffic to real providers.
- **Options Considered:**
  - Local emulator only — pros: safe, offline, deterministic; cons: does not test real provider quirks.
  - Live relay/router — pros: tests real providers; cons: secrets, money, compliance, network flakiness.
- **Decision:** OpenMuara is a **local emulator**. It does not proxy to real providers by default.
- **Consequences:** Webhook relay forwards only locally-emulated webhooks to dev destinations.
- **Reversible?** Yes — a live-router mode could be added later behind explicit config.
- **Reversal Trigger:** If future users demand live sandbox routing.
- **Logged By:** AI Agent (recovered from session)

### D004 — Provider priority: small/local gateways first

- **Date:** 2026-06-26 (recovered from previous session)
- **Context:** Many providers could be emulated; need focus.
- **Options Considered:**
  - Stripe-first — pros: great docs, big market; cons: Stripe CLI already exists, less differentiation.
  - Small/local providers first — pros: high pain, no existing tooling, strong differentiation; cons: smaller market, fragmented docs.
- **Decision:** Focus v1 on **Fawry + SenangPay + Stripe** to prove both local-gateway and mainstream flows.
- **Consequences:** RevenueCat, App Store, Play Store are **hard frozen for v1 and targeted for v2**.
- **Reversible?** No for v1 — explicitly excluded from v1 scope.
- **Reversal Trigger:** N/A for v1; may be reconsidered when planning v2.
- **Logged By:** AI Agent (recovered from session)

### D006 — Persistence stays SQLite-first for v1

- **Date:** 2026-06-28
- **Context:** OpenMuara needs durable local state without adding external services.
- **Decision:** v1 uses SQLite as the default persistence layer, with an in-memory option for tests.
- **Consequences:**
  - Transaction and audit data live in `.muara/data/ledger.db` by default.
  - Single-writer concurrency limits are acceptable for a local-first tool.
- **Reversible?** Yes — future versions may add pluggable backends.
- **Reversal Trigger:** Demand for multi-writer or remote persistence.
- **Logged By:** AI Agent

### D007 — Provider registry is declarative

- **Date:** 2026-06-28
- **Context:** Providers need to be enabled and configured without code changes.
- **Decision:** Providers are activated through `providers.<name>.enabled` in `.muara/config.yml`.
- **Consequences:**
  - New providers can be registered at compile time and enabled via config.
  - Provider configs are cloned before `Init` to prevent mutation.
- **Reversible?** Yes — schema can be extended.
- **Reversal Trigger:** Need for dynamic provider loading.
- **Logged By:** AI Agent

### D008 — CSRF for admin UI uses double-submit cookie

- **Date:** 2026-06-28
- **Context:** The embedded admin dashboard issues state-changing requests.
- **Decision:** Use a non-`HttpOnly` cookie `openmuara_csrf` plus `X-CSRF-Token` / `csrf_token`. `/_admin/webhook-receiver` is exempt.
- **Consequences:**
  - Dashboard reads the token from a meta tag.
  - CSRF can be disabled via config in isolated environments.
- **Reversible?** Yes — can switch to session tokens later.
- **Reversal Trigger:** Auth/session layer added.
- **Logged By:** AI Agent

### D009 — Prometheus metrics use the `openmuara_` prefix

- **Date:** 2026-06-28
- **Context:** Metrics need a namespace to avoid collision with application metrics.
- **Decision:** All Prometheus metrics are prefixed with `openmuara_`.
- **Consequences:**
  - Scraping configs and alerts can filter by prefix.
  - Label cardinality is bounded by route count and provider list.
- **Reversible?** Yes — prefix can be changed in a major release.
- **Reversal Trigger:** Organizational naming standard changes.
- **Logged By:** AI Agent

### D010 — Transaction and audit stores share one SQLite connection

- **Date:** 2026-06-28
- **Context:** Two separate SQLite connections caused `database is locked` errors at startup.
- **Decision:** `internal/cli/start.go` opens a single `*sql.DB` and shares it between transaction and audit stores.
- **Consequences:**
  - Startup contention is eliminated.
  - Long-running writes still block; acceptable for a local emulator.
- **Reversible?** Yes — can switch to connection pooling later.
- **Reversal Trigger:** Need for higher write concurrency.
- **Logged By:** AI Agent

---

## Auto-Decisions (AFK Mode)

| ID | Prompt | Decision | Options Considered | Why Chosen | Reversible? | Date |
|----|--------|----------|-------------------|------------|-------------|------|
| A01 | Plan creation | Work only on `dev`, no feature branches | Feature branches vs. direct `dev` | Human explicitly requested consolidation on `dev` | Yes — can branch later if needed | 2026-06-27 |
| A02 | Plan creation | Use `docs/projects/openmuara-v1/` for planning | `.agents/project/` vs. `docs/projects/` | `docs/projects/` matches recovered mkp template and is committed to repo | Yes | 2026-06-27 |
