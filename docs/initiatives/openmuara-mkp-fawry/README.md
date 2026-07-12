> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara MKP Fawry Integration

> **Status:** ⬜ Not Started | **Started:** —
> **Scope:** Close the gaps between OpenMuara's Fawry emulation and MKP v2's billing requirements so MKP can eject its internal Fawry simulator.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `feat/mkp-fawry`
>
> **Consumer:** Muslim Kids Platform (MKP) v2
> **Consumer Handler:** `services/mkp-v2-api/internal/api/billing/fawry.go`
> **Source Requirements:** `docs/mkp-billing-requirements.md`

---

## Initiative Structure

```
docs/initiatives/openmuara-mkp-fawry/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
│
└── prompts/               # Numbered, self-contained execution prompts
    ├── _template.md
    ├── 01-fawry-state-extensions.md
    ├── 02-response-delay.md
    ├── 03-billing-type-and-journey.md
    ├── 04-escape-page-and-webhook-shape.md
    └── 05-docs-and-ci.md
```

Planning docs live in `docs/initiatives/openmuara-mkp-fawry/` in the root repo.
Product code commits to the `feat/mkp-fawry` branch. Do not commit directly to
`main`.

---

## Why this initiative?

MKP v2 currently maintains its own internal Fawry simulator. OpenMuara already
emulates the core Fawry flow, but a few mismatches prevent MKP from pointing
its handler at OpenMuara:

1. MKP tests expect `OrderStatus` values `CANCELED` and `EXPIRED`, not just
   `PAID` / `UNPAID`.
2. MKP wants a configurable `response_delay_ms` to simulate slow gateways.
3. MKP has two journey types — `recurring` subscription and `one_time` prepaid —
   and the webhook payload should signal which one is in play.
4. The admin escape page and docs need to expose these new capabilities.

This initiative closes those gaps while keeping the existing Fawry API
backward-compatible.

---

## Goals

1. **Extended Fawry states** — Support `PAID`, `UNPAID`, `CANCELED`, and
   `EXPIRED` order statuses end-to-end (ledger, escape page, webhook payload).
2. **Configurable response delay** — Add `fawry.response_delay_ms` and apply it
   to the outgoing webhook dispatch.
3. **Billing-type journey hint** — Accept `billing_type` on `/fawry/charge` and
   shape the Fawry V2 payload for subscription vs prepaid.
4. **Escape page update** — Add canceled / expired actions and surface the
   billing type in the UI.
5. **Docs and tests** — Update the Fawry provider guide, add MKP-focused
   examples, and add tests for all new states and delays.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and
code style.

### 2. Backward compatibility
Existing Fawry routes, statuses, and configs keep working. New fields are
optional and default to current behavior.

### 3. No external services
OpenMuara stays local-first. Delays are in-process sleeps, not external calls.

### 4. Three-interface parity
New capabilities must be reachable from:

- The HTTP API (`/fawry/charge`, `/_admin/fawry-escape`).
- The dashboard escape page.
- The configuration file (and env var overrides).

### 5. Quality gates
Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

### 6. Definition of done
Beyond the quality gates, a prompt is done only when:

- The feature works end-to-end for the MKP Fawry handler.
- Tests cover happy path, error path, and edge cases.
- The smoke test is updated if routes, CLI flags, or default behavior change.
- `HANDOFF.md` is updated with what was built and what changed.
- `TRACKING.md` marks the prompt `✅` with the commit hash.

---

## Out of Scope

- Stripe or RevenueCat emulation for MKP.
- Changes to the provider plugin schema contract.
- Real Fawry sandbox connectivity.
- Dashboard authentication.

---

## Success Criteria

- MKP can replace its internal Fawry simulator with OpenMuara for local dev
  and CI.
- All four `OrderStatus` values round-trip through charge → escape → webhook.
- `response_delay_ms` delays the outgoing webhook without blocking the escape
  redirect.
- `billing_type=recurring` produces a subscription-shaped webhook payload;
  `billing_type=one_time` produces a prepaid-shaped payload.
- Existing Fawry tests and the smoke test still pass.
- `docs/providers/fawry.md` explains the new options with copy-paste examples.
