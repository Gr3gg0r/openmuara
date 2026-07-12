---
id: architecture
title: OpenMuara Architecture
---

# OpenMuara Architecture

OpenMuara is a local-first payment virtualization layer. It emulates payment providers so that applications can exercise billing flows offline, fast, and deterministically.

## High-level flow

```
┌─────────────────┐      HTTP       ┌──────────────────────┐
│   Test client   │ ───────────────▶ │  OpenMuara server    │
│                 │ ◀─────────────── │  (internal/server)   │
└─────────────────┘                  └──────────────────────┘
                                              │
           ┌──────────────────────────────────┼──────────────────────────────────┐
           ▼                                  ▼                                  ▼
  ┌─────────────────┐              ┌─────────────────┐                ┌─────────────────┐
  │ Provider plugin │              │ Engine ledger   │                │ Webhook         │
  │ (fawry/stripe/  │              │ (SQLite/memory) │                │ dispatcher      │
  │  senangpay/...) │              │                 │                │                 │
  └─────────────────┘              └─────────────────┘                └─────────────────┘
```

1. **Incoming request** — `internal/server/router.go` routes provider-specific or universal API requests.
2. **Provider emulation** — each provider validates signatures, creates a transaction, and returns a provider-shaped response.
3. **Ledger** — `internal/engine` stores transactions durably in SQLite (default) or memory.
4. **Outgoing webhooks** — the dispatcher emits provider-style webhook payloads to a configured local URL.
5. **Admin UI / API** — `/_admin` endpoints and the embedded dashboard expose transactions, webhooks, and replay controls.

## Key packages

| Package | Responsibility |
|---------|----------------|
| `internal/server` | HTTP router, middleware (request ID, logging, CORS, CSRF, metrics, body limit), admin handlers |
| `internal/api` | Universal payment API (`/v1/pay`, `/v1/refund`) |
| `internal/engine` | Transaction ledger, state machine, SQLite persistence |
| `internal/provider` | Provider registry and shared interfaces |
| `internal/fawry` | Fawry Express Checkout emulation |
| `internal/stripe` | Stripe Checkout emulation |
| `internal/senangpay` | SenangPay charge/callback/webhook emulation |
| `internal/webhook` | Outgoing webhook delivery, replay, attempt store |
| `internal/audit` | Structured audit log storage and logger |
| `internal/config` | YAML/env configuration loader |
| `internal/ui` | Embedded admin HTML dashboard |

## Provider versioning

Providers that emulate more than one API version use a versioned sub-package layout:

```
internal/<provider>/
├── provider.go        # version dispatcher; selects v1/v2 from config
├── charge.go          # version-agnostic charge API (if shared across versions)
├── escape.go          # admin escape page
├── v1/
│   └── webhook.go     # v1 webhook receiver and payload builder
│   └── provider.go    # v1 concrete provider
└── v2/
    └── webhook.go     # v2 webhook receiver and payload builder
    └── provider.go    # v2 concrete provider
```

- The top-level provider reads `cfg["version"]` (default `"v1"`) and delegates webhook/payload behavior to the selected version.
- Legacy unversioned routes (e.g., `/fawry/charge`, `/fawry/webhook`) continue to work and map to the configured default version.
- Explicit versioned routes (e.g., `/fawry/v1/webhook`, `/fawry/v2/webhook`) are always registered so tests and clients can target a specific version. This lets you pilot a new version while production traffic stays on the current default.
- Single-version providers stay flat; they do not need a `v1/` sub-package until a second version appears.
- The admin dashboard `/_admin` surfaces the active version and supported versions for each enabled provider, and `muara.yml.example` documents the `version` config key.

## State machine

Transactions move through states using `engine.Transition`. Invalid transitions return HTTP 409. Common states:

- `new`
- `paid`
- `unpaid`
- `refunded`

## Middleware chain

Request processing order (outermost → innermost):

1. Request ID
2. Logging
3. CORS
4. CSRF guard
5. Metrics
6. Max body size
7. Route handler

`/metrics` is excluded from request counters. `/_admin/webhook-receiver` is excluded from CSRF because it receives server-to-server webhooks.

## Observability

- **Health:** `GET /healthz` and `GET /readyz`
- **Metrics:** `GET /metrics` in Prometheus format
- **Audit:** `GET /_admin/audit` and `muara audit list`
- **Dashboard:** `GET /_admin`
