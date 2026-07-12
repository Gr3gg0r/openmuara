> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests — Risk Register

## Risk Scoring

| Score | Likelihood | Impact |
|---|---|---|
| 1 | Rare | Negligible |
| 2 | Unlikely | Minor |
| 3 | Possible | Moderate |
| 4 | Likely | Major |
| 5 | Almost certain | Critical |

**Risk Score = Likelihood × Impact.** Scores ≥ 12 are high priority and require active mitigation.

---

## Risk Register

| ID | Risk | Likelihood | Impact | Score | Owner | Mitigation | Contingency | Status |
|---|---|---|---:|---|---|---|---|---|
| R001 | Removing `init()` registrations breaks tests that rely on `provider.Get("fawry")` | 4 | 3 | 12 | Agent | Update tests to load manifests explicitly or use factory registry. | If tests are too entangled, add a temporary test helper that loads built-ins via manifest fixtures. | Open |
| R002 | `runtime.type: go` factory registry design is too complex | 3 | 3 | 9 | Agent | Start with a map keyed by provider name; defer interface over-engineering. | If registry becomes unwieldy, split into `internal/provider/factory/` with minimal interface. | Open |
| R003 | Existing `plugins/*/gateway.yml` files are inconsistent | 3 | 3 | 9 | Agent | Audit and normalize manifests as part of P04; add schema validation. | If inconsistency is deep, split P04 into per-provider sub-prompts. | Open |
| R004 | Stripe has no manifest today | 4 | 2 | 8 | Agent | Create `plugins/stripe/gateway.yml` with `runtime.type: go` in P04. | If Stripe is too complex for current factory, defer to a follow-up task with explicit note. | Open |
| R005 | Loader ordering changes break existing `.muara/config.yml` users | 3 | 4 | 12 | Agent | Document migration in `docs/migration/` and print clear startup warnings. | If breakage is severe, add a compatibility mode flag for one release cycle. | Open |
| R006 | Parallel agent sessions produce conflicting product-code changes | 3 | 3 | 9 | Agent | Re-run gates, diff carefully, and commit one logical change at a time. | If conflicts are severe, branch the parallel work and merge explicitly. | Active |
| R007 | Contributors misunderstand when to use `simple` vs `go` | 3 | 2 | 6 | Agent | Add decision tree to `docs/contributing-providers.md` and `appendices/b-simple-vs-go-decision-tree.md`. | Add `muara provider diagnose` command in future to suggest runtime type. | Open |
| R008 | WASM sandboxed runtime becomes a premature abstraction | 2 | 4 | 8 | Agent | Explicitly defer to a future initiative; keep architecture door open. | If pressure builds, create a spike initiative rather than expanding scope. | Open |
| R009 | Provider protocol emulation regresses during refactor | 3 | 4 | 12 | Domain expert / Agent | Add conformance tests for each provider before changing loader; run checkout-store smoke tests. | If regression found, revert and re-plan the refactor. | Open |
| R010 | Documentation becomes stale during implementation | 3 | 2 | 6 | Agent | Make docs updates part of P04 acceptance criteria; include docs in review checklist. | If docs drift, block PR until corrected. | Open |
| R011 | Factory registry introduces global mutable state | 2 | 3 | 6 | Agent | Registry is read-only after `init()`; use `sync.Once` or explicit initialization; no runtime registration. | If mutable state leaks, refactor registry to be owned by loader and passed explicitly. | Open |
| R012 | Coverage drops on changed modules | 3 | 2 | 6 | Agent | Require tests for new loader paths and factory registry; gate on coverage diff. | If coverage drops, add targeted tests before proceeding. | Open |

---

## Watch List

Files and behaviors to monitor closely:

- `internal/provider/conform/conform_test.go` — likely depends on global provider registration.
- `internal/server/providers_test.go` — may assume built-ins are pre-loaded.
- `internal/cli/start.go` — startup path must load manifests correctly.
- `internal/config/provider_loader.go` — core of the manifest-first change.
- `internal/server/router.go` — must not depend on global provider registry.
- `examples/checkout-store/` — end-to-end validation for Fawry and Stripe.
- `docs/contributing-providers.md` — must stay accurate.

---

## Trigger Conditions

| Risk | Trigger | Response |
|---|---|---|
| R001 | `go test ./...` fails in `internal/provider/conform` after removing `init()` | Add manifest-loading test helper; update affected tests. |
| R005 | Existing config fails to start after loader change | Add compatibility warning; document migration. |
| R009 | Checkout-store smoke test fails | Stop; identify regression; add targeted conformance test. |
| R006 | `git status` shows unexpected conflicts | Pause; run gates; merge or stash parallel work explicitly. |

---

## Residual Risks

After mitigation, the following risks remain:

- R005: Users with old configs may still be surprised by loader changes until migration docs are read.
- R007: Contributors may still choose the wrong runtime type without human review.
- R008: WASM remains a future promise; third-party extensibility is not yet implemented.

These are accepted and documented.
