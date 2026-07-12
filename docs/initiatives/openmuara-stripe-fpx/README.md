> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Stripe FPX

> **Status:** ❄️ Archived / Superseded | **Started:** 2026-06-30 | **Archived:** 2026-07-03
> **Superseded by:** [`docs/initiatives/openmuara-stripe-checkout-sessions/`](../openmuara-stripe-checkout-sessions/)
> **Scope:** Emulate Stripe FPX (Malaysian online bank transfer) payment flows using the same charge + escape pattern as the Fawry provider.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** ___________
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`

---

## Initiative Structure

```
docs/initiatives/openmuara-stripe-fpx/
├── README.md              # This file
├── HOWTO.md               # Decomposition guide for AI
├── PREREQUISITES.md       # Human pre-flight checklist
├── TRACKING.md            # Central execution tracker
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing bugs / out-of-scope
├── REFERENCES.md          # Links to specs, runbooks, vendor docs
├── .gitignore             # Ignore screenshots, logs, temp files
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── _template.md
│   └── 01-stripe-fpx-charge-and-escape.md
│
├── tasks/                 # (Optional) Detailed specs — dual-layer
├── findings/              # Research, audit output, analysis
│   └── 2026-07-03-stripe-fpx-audit.md
├── runbooks/              # Operational docs
├── screenshots/           # QA evidence (gitignored)
├── qa/                    # Validation artifacts (gitignored)
└── state/                 # Agent state snapshots (gitignored)
```

Planning docs live in `docs/initiatives/openmuara-stripe-fpx/` in the root repo. Product code commits to the `dev` branch. Do not commit directly to `main`.

---

## Why FPX?

FPX is Malaysia’s dominant online bank transfer rail and a first-class Stripe payment method. Supporting it lets Malaysian developers test their Stripe integration locally without real MYR transactions.

---

## Goals

Follow the Fawry charge + escape pattern for both FPX and card payments:

1. `POST /v1/stripe/fpx/charge` creates an FPX charge, records it in the ledger, and returns a reference.
2. `GET /v1/stripe/fpx/escape` renders a minimal bank selector page.
3. `POST /v1/stripe/fpx/escape` confirms or cancels the FPX payment, updates the ledger, dispatches a Stripe-signed webhook, and redirects to the caller.
4. `POST /v1/stripe/card/charge` creates a card charge, records it in the ledger, and returns a reference.
5. `GET /v1/stripe/card/escape` renders a minimal card confirmation page.
6. `POST /v1/stripe/card/escape` confirms or cancels the card payment, updates the ledger, dispatches a Stripe-signed webhook, and redirects to the caller.
7. Preserve all existing Stripe Checkout session behavior.

---

## Conventions (MANDATORY)

### 1. Read `AGENTS.md` first
Root `AGENTS.md` governs branch rules, quality gates, autonomy boundaries, and code style. This initiative does not repeat every rule.

### 2. Provider contract fidelity
FPX emulation must match Stripe’s documented behavior for the subset we implement, including:
- FPX charge request/response shape
- Stripe webhook signature header (`Stripe-Signature`)
- Webhook events: `checkout.session.completed` for success, `payment_intent.canceled` for cancel

### 3. Quality gates
Every prompt must pass:
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run`
- `./scripts/smoke-test.sh`

---

## Supersession note

This initiative was completed but later **superseded** by
[`docs/initiatives/openmuara-stripe-checkout-sessions/`](../openmuara-stripe-checkout-sessions/).
The custom `/v1/stripe/fpx/*` and `/v1/stripe/card/*` routes implemented in P01 were
removed because they are OpenMuara-specific and break Stripe SDK compatibility.
The replacement initiative implements the real Stripe APIs:

- `POST /v1/checkout/sessions` with `payment_method_types: ["fpx"]` / `["card"]`
- `POST /v1/payment_intents` with test payment method tokens
- Local OpenMuara-hosted checkout / authentication pages
- Stripe-compatible webhook configuration UI

This directory is retained as an archive and decision record.

## Out of Scope

- Real bank selection UI polish beyond a minimal HTML picker
- FPX-specific disputes/refunds (reuse Stripe simulation endpoints)
- Non-Stripe FPX providers (iPay88, Razer Merchant Services, etc.)
- FPX via Stripe Checkout Sessions

## Audit findings & recommendations

A post-completion audit identified several gaps that were not addressed before
supersession. These are documented below as lessons for future provider
emulation initiatives.

| # | Finding | Priority | Recommendation | Status |
|---|---------|----------|----------------|--------|
| R1 | Custom routes break Stripe SDK parity | High | Prefer real Stripe API paths (`/v1/checkout/sessions`, `/v1/payment_intents`) so client code works unchanged against real Stripe. | ✅ Addressed by successor |
| R2 | No OpenAPI spec updates | High | Every new public route must be added to `docs/openapi.yaml` before the prompt is closed. | ✅ Addressed by successor |
| R3 | No example app or usage docs | Medium | Provide a minimal example (`examples/stripe-fpx/`) showing request/response flow and webhook handling. | ❌ Not addressed |
| R4 | No operational runbook | Medium | Add a runbook covering common FPX test scenarios, bank selection, and webhook signature verification. | ❌ Not addressed |
| R5 | Limited edge-case coverage | Medium | Tests for idempotency, duplicate confirm/cancel, invalid bank codes, and CSRF on escape pages. | ❌ Not addressed |
| R6 | No dashboard integration | Medium | Escape pages should be inspectable/replayable from `/_admin` alongside other transactions. | ✅ Addressed by successor |
| R7 | Currency defaults undocumented | Low | Document that FPX defaults to `myr` and card defaults to `usd` in request/response examples. | ❌ Not addressed |
| R8 | Webhook event scope narrow | Low | Explicitly document which webhook events are emitted and which are deferred to generic simulation endpoints. | ✅ Partially addressed |
| R9 | Initiative not archived | High | Mark superseded initiatives as archived, link to the replacement, and record the decision. | ✅ Addressed by this audit |
| R10 | Status inconsistency | High | Keep `README.md`, `TRACKING.md`, and `HANDOFF.md` status in sync. | ✅ Addressed by this audit |

## Self-rating

Pre-audit: **6/10** — functionality was delivered and tested, but the initiative
lacked OpenAPI updates, examples, runbooks, edge-case coverage, and a clear
supersession path.

Post-audit (after applying R9, R10): **8.5/10** — archival and cross-linking are
solid; remaining recommendations (R3–R8) are valuable but no longer actionable
because the implementation was replaced. They are preserved as lessons learned.
