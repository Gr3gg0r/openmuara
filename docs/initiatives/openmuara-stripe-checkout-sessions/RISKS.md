# OpenMuara Stripe FPX & Card Payments — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | Stripe Checkout Session or PaymentIntents API shape changes after implementation. | Low | Medium | Scope to documented subset; keep contract/golden tests so deviations are caught. |
| R02 | Session IDs (`cs_test_*`) and PaymentIntent IDs (`pi_test_*`) collide in the ledger. | Low | High | Use distinct prefixes; store PaymentIntents in a dedicated in-memory store while still writing a ledger transaction with the PI reference. |
| R03 | Existing Stripe Checkout tests break when refactoring provider routes. | Medium | High | Run full `go test ./internal/stripe/...` after every sub-step; keep create/retrieve code paths stable during P01. |
| R04 | Card confirmation bypasses real 3-D Secure, diverging from production Stripe. | Low | Low | Document as a limitation; test-mode card payments succeed immediately, matching Stripe test behavior. |
| R05 | Webhook event type changes confuse existing consumers. | Low | Medium | Emit documented event types only; add tests asserting payload shape. |
| R06 | Webhook configuration UI cannot persist changes because dispatcher URL is set at startup. | Low | Medium | Persist to `.muara/config.yml` and update the running dispatcher; document in `DECISIONS.md`. |
| R07 | PaymentIntent required fields differ from Stripe SDK expectations, causing deserialization issues. | Low | Medium | Return the documented subset; include `id`, `object`, `amount`, `currency`, `status`, `client_secret`, `payment_method_types`, and `next_action` where applicable. |
| R08 | HTML smoke tests become brittle if checkout/authentication page markup changes. | Medium | Medium | Use semantic form fields and `data-testid` attributes; keep smoke tests focused on form action URLs and hidden fields. |
| R09 | Runtime write-back to `.muara/config.yml` races with manual edits or config reload. | Low | Low | Use atomic file write; reload config after save via existing Viper watch or explicit admin reload. |
| R10 | Adding `GET /v1/checkout/sessions/{id}/pay` changes the security model of the session URL path. | Low | Low | The path is already exposed in the `url` field; adding a handler just completes the contract. Keep id tokens unguessable (UUID). |
