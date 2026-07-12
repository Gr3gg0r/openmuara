> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Solid — Known Issues & Out-of-Scope List

> **Purpose:** Prevent the AI from wasting time on pre-existing bugs or out-of-scope problems.

---

## Pre-Existing Bugs (Do NOT Fix Unless Caused by This Initiative)

| ID | Issue | Location | Impact | Why Out of Scope |
|----|-------|----------|--------|------------------|
| K01 | Admin dashboard expects array responses | `internal/ui/index.html` | High | Fixed by **S01**. |
| K02 | OpenAPI spec drift (`/readyz`, pagination, 409) | `docs/openapi.yaml` | High | Fixed by **S02**. |
| K03 | Stripe simulation mutates status directly | `internal/stripe/simulation.go` | High | Fixed by **S03**. |
| K04 | Fawry escape does not update ledger | `internal/fawry/escape.go` | High | Fixed by **S04**. |
| K05 | Fawry incoming webhook ignores signature | `internal/fawry/webhook.go` | High | Fixed by **S04**. |
| K06 | Only active provider dispatcher is wired | `internal/cli/start.go` | Medium | Fixed by **S05**. |

---

## Out-of-Scope Areas (Hard Boundaries)

| Area | Reason | Boundary |
|------|--------|----------|
| App Store / Play Store / RevenueCat | Hard frozen for v2 per `DECISIONS.md` | Do not implement in v1. |
| Multi-port runtime | Deferred to v1.2+ | Do not implement unless added to scope. |
| MCP server | Deferred to v1.2+ | Do not implement unless added to scope. |
| SaaS / hosted service | Out of project vision | Do not implement. |

---

## How to Use This File

1. **Before starting a step:** Scan this file.
2. **During execution:** If you hit a pre-existing bug unrelated to your task, log it here and move on.
3. **In the prompt:** Reference this file if the step touches code near a known issue.
