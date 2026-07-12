> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 01 — Fawry State Extensions

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** —

---

## Goal

Add `canceled` and `expired` transaction states and wire them through the Fawry
escape page and webhook payload so MKP can test all four `OrderStatus` values.

## Why now

MKP's Fawry handler expects `PAID`, `UNPAID`, `CANCELED`, and `EXPIRED`
`OrderStatus` values. OpenMuara currently only models `new`, `paid`, `unpaid`,
and `refunded`.

## Scope

### In scope

- Add `TransactionStatusCanceled` and `TransactionStatusExpired` to
  `internal/engine/transaction.go`.
- Update `validTransitions` so a `new` charge can transition to `canceled` or
  `expired`, and a `paid` charge cannot.
- Update `internal/fawry/escape.go` to accept `status=CANCELED` and
  `status=EXPIRED` from the escape form.
- Map these statuses to the Fawry V2 webhook payload `OrderStatus` field.
- Update tests for state transitions and escape handling.

### Out of scope

- Stripe or RevenueCat states.
- Changing the Fawry V1 payload contract.

## Acceptance criteria

- [ ] A charge can be created in `new` state and escaped to `PAID`, `UNPAID`,
      `CANCELED`, or `EXPIRED`.
- [ ] Invalid transitions (e.g., `paid` → `canceled`) return HTTP 409.
- [ ] The Fawry V2 webhook payload carries the matching `OrderStatus`.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Keep status strings lower-case internally (`canceled`, `expired`) and map to
  upper-case Fawry values only in the payload builder.
- Fuzz tests for state transitions may need new seed values.

## Deliverables

- Code changes on `feat/mkp-fawry`.
- Updated tests for transitions and escape.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
