# OpenMuara Stripe FPX — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | New `/v1/stripe/fpx/*` routes are confused with real Stripe API paths and break client expectations. | Low | Medium | Document clearly in README and runbooks that these are OpenMuara-emulated routes; keep response shapes Stripe-like but mark ids with `test_` prefix. |
| R02 | Webhook signature scheme diverges from Stripe’s real scheme. | Low | High | Reuse the existing `internal/stripe/signature.go` `SignPayload` implementation and verify with existing signature tests. |
| R03 | FPX escape page collides with existing `/_admin/*` or `/v1/*` routing. | Low | Medium | Register under `/v1/stripe/fpx/*` consistently; verify in `internal/server/router.go` that provider routes are mounted without collision. |
| R04 | State-machine transitions fail because FPX uses different status strings than Checkout Sessions. | Medium | Medium | Map `PAID` → `paid`, `CANCELED` → `canceled`; use `engine.Transition` and test invalid transitions. |
| R05 | Existing Stripe Checkout tests regress due to shared provider state. | Low | High | Run full `go test ./internal/stripe/...` after changes; use isolated test registries. |

## Post-supersession risks / lessons learned

| ID | Risk / Lesson | Mitigation |
|----|---------------|------------|
| L01 | Custom provider routes create lock-in and migration work. | Start with real provider API paths; only add OpenMuara-specific extensions when absolutely necessary. |
| L02 | Missing OpenAPI updates make it hard to discover emulated endpoints. | Update `docs/openapi.yaml` as part of every prompt that adds public routes. |
| L03 | Lack of example apps slows user adoption. | Create a minimal example directory for each new provider flow. |
| L04 | Minimal escape pages without admin visibility complicate debugging. | Host payment pages under `/_admin` or expose them in the dashboard ledger for replay/inspect. |
| L05 | Superseded initiatives left in "Active" status confuse future agents. | Archive promptly, link to successor, and reconcile status across README, TRACKING, and HANDOFF. |
