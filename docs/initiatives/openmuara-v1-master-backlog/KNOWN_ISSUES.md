> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara v1 Master Backlog — Known Issues & Boundaries

> **Purpose:** Prevent wasted effort on pre-existing bugs or out-of-scope problems.

---

## Resolved Bugs

The following issues were tracked as pre-existing bugs and have been resolved in v1.

| ID | Issue | Resolution |
|----|-------|------------|
| K01 | Admin dashboard expected array responses | Fixed by S01 — pagination envelope `{ limit, offset, results }` |
| K02 | OpenAPI spec drift (`/readyz`, pagination, 409) | Fixed by S02 — spec synced and sync test added |
| K03 | Stripe simulation mutated status directly | Fixed by S03 — `engine.Transition` state machine |
| K04 | Fawry escape did not update ledger | Fixed by S04 — escape now transitions transaction status |
| K05 | Fawry incoming webhook ignored signature | Fixed by S04 — HMAC signature verification added |
| K06 | Only active provider dispatcher was wired | Fixed by S05 — per-provider dispatcher map |

---

## Operational Limitations

These are expected boundaries of OpenMuara v1, not bugs to fix.

| ID | Limitation | Impact | Mitigation |
|----|------------|--------|------------|
| L01 | No built-in authentication on `/_admin`, `/metrics`, or provider routes | Anyone with network access can view/modify state | Run behind a reverse proxy or on `127.0.0.1` only |
| L02 | SQLite is a single-writer store | High write concurrency can return `database is locked` | Use one instance per environment; avoid parallel load tests |
| L03 | Audit log grows unbounded | Disk usage increases over time | Periodically archive or prune old `audit_logs` rows |
| L04 | Webhook retries are immediate, no backoff or dead-letter queue | Bursty failures may retry quickly | Fix the consumer promptly; replay manually after fixing |
| L05 | Metrics endpoint is unauthenticated | Metric counts (not payloads) are public within network scope | Bind to localhost or protect with a reverse proxy |
| L06 | CORS and CSRF settings are global | Cannot configure per-provider or per-route rules | Set origins for the whole server; disable CSRF only in isolated environments |

---

## Out-of-Scope Areas (Hard Boundaries)

| Area | Reason | Boundary |
|------|--------|----------|
| App Store / Play Store receipt validation | Hard frozen for v2 | Do not implement in v1. |
| RevenueCat emulation | Moved to v2 initiative `docs/initiatives/openmuara-v2-revenuecat/` | Do not implement in v1. |
| Multi-port runtime | Deferred to v1.2+ | Do not implement unless added to scope. |
| MCP server | Deferred to v1.2+ | Do not implement unless added to scope. |
| SaaS / hosted service | Out of project vision | Do not implement. |

---

## How to Use This File

1. Before starting a backlog item, scan this file.
2. If you hit a pre-existing bug unrelated to your task, log it here and move on.
3. If the bug is in-scope for an active prompt, fix it there and update status.
