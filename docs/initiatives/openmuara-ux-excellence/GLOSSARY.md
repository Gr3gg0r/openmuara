> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara UX Excellence — Glossary

> **Created:** 2026-07-01

---

| Term | Meaning in this initiative |
|------|---------------------------|
| **Ledger** | The unified dashboard view of transactions and webhook events, named after the financial transaction ledger in `internal/engine`. |
| **Wizard** | The interactive `muara init` flow that asks questions and generates a tailored config. |
| **Provider metadata** | Stable schema exposed by `GET /_admin/providers` describing each provider (name, description, emulated real providers, sample route, category, etc.). |
| **Signature status** | Whether OpenMuara considers a webhook payload's signature valid, stored on `webhook.Attempt`. |
| **Three-interface parity** | The rule that web UI, CLI, and admin API should expose the same capabilities so no persona is forced into the wrong tool. |
| **Zero-data state** | The dashboard UI shown before any transactions or webhooks exist, with guidance on how to generate the first event. |
| **Additive UX** | New features layer on top of existing behavior without breaking configs, routes, or commands. |
