> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# 04 — Fawry Escape Ledger Update and Webhook Signature Verification

## Objective

When a user completes a Fawry payment via the escape page, the shared transaction ledger must reflect the outcome. Incoming Fawry V2 webhooks must have their signatures verified when a webhook secret is configured.

## Background

`POST /_admin/fawry-escape` currently only dispatches an outgoing webhook. The transaction remains `new` in the ledger, which breaks dashboard consistency and refund/scenario workflows that depend on accurate status.

Incoming `POST /fawry/webhook` accepts any JSON body without verifying `messageSignature`, making it easy to inject fake webhook events during local testing.

## Constraints

- Use the existing `engine.TransactionStore` interface.
- Use `engine.Transition` to enforce the state machine.
- Webhook signature verification must be **optional**: skip if `webhook_secret` is empty.
- Do not break the existing smoke test.
- App Store / Play Store / RevenueCat remain frozen for v2.

## Acceptance Criteria

- [ ] `POST /_admin/fawry-escape` updates the matching transaction to `paid` or `unpaid`.
- [ ] If the transaction reference does not exist, the handler returns 404.
- [ ] `POST /fawry/webhook` returns 401 when signature verification is enabled and the signature is invalid.
- [ ] `POST /fawry/webhook` returns 200 and processes the payload when the signature is valid.
- [ ] When `webhook_secret` is empty, signature verification is skipped and the payload is processed.

## Test Expectations

- Unit test: escape action updates ledger status to `paid`.
- Unit test: escape action returns 404 for missing reference.
- Unit test: webhook handler rejects invalid signature with 401.
- Unit test: webhook handler accepts valid signature with 200.
- Smoke test passes end-to-end.

## Rollback Trigger

If the smoke test fails after this change and cannot be fixed within 30 minutes, revert the commit.

## References

- Fawry signature logic: `internal/fawry/signature.go`
- Escape handler: `internal/fawry/escape.go`
- Webhook handler: `internal/fawry/webhook.go`
