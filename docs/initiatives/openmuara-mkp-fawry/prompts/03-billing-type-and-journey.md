> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 03 — Billing Type and Journey

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** Prompt 01

---

## Goal

Let `/fawry/charge` accept a `billing_type` hint and shape the outgoing Fawry V2
webhook payload for subscription (`recurring`) vs prepaid (`one_time`) journeys.

## Why now

MKP has two product journeys. The Fawry handler needs to know whether a webhook
represents a subscription or a one-time purchase.

## Scope

### In scope

- Accept optional `billing_type` field in `/fawry/charge` (values:
  `recurring`, `one_time`).
- Store `billing_type` on the transaction (e.g., `Transaction.Type` or a new
  field).
- Pass billing type to the Fawry V2 payload builder.
- Adjust the webhook payload to include subscription-relevant fields when
  `recurring` (e.g., `orderItems` may include subscription metadata) and keep it
  minimal for `one_time`.
- Default to `one_time` when omitted, preserving current behavior.
- Add tests for both journey shapes.

### Out of scope

- Real subscription lifecycle (renewal, expiration events) beyond the initial
  webhook shape.
- Changes to the Fawry V1 contract.

## Acceptance criteria

- [ ] `POST /fawry/charge` with `billing_type=recurring` stores the hint.
- [ ] The outgoing Fawry V2 webhook includes journey-specific fields for
      recurring and a simpler shape for one_time.
- [ ] Omitted `billing_type` defaults to one-time behavior.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Use a typed string for `BillingType` to avoid magic strings.
- The payload builder can branch on `tx.Type` or a dedicated field.

## Deliverables

- Code changes on `feat/mkp-fawry`.
- Updated charge and webhook tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
