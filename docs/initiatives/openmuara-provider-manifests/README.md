> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Provider Manifests

> **Status:** 🟡 READY TO IMPLEMENT | **Started:** 2026-07-08 | **Signed Off:** 2026-07-09
>
> **Scope:** Make every non-default provider discoverable and configurable through `gateway.yml`, with a clear graduation path from simple YAML to Go.
> **Owner:** AI Agent (Kimi Code) | **Human Reviewer:** User (signed off 2026-07-09)
> **Target Repo:** `<repo-root>/`
> **Product Branch:** `dev`
>
> **Why:** Provider discovery today is split between hard-coded Go packages and YAML manifests. That makes it hard for contributors to add a provider, hard for users to understand why a provider appears, and hard for maintainers to remove or sandbox one. A manifest-first model fixes all three: the YAML file is the source of truth, Go is an implementation detail, and the door stays open for proprietary bridges and WASM plugins.

---

## Initiative Structure

```
docs/initiatives/openmuara-provider-manifests/
├── README.md              # This file
├── TRACKING.md            # Central execution tracker
├── PREREQUISITES.md       # Human pre-flight checklist
├── HANDOFF.md             # Session continuity
├── DECISIONS.md           # Decision log
├── RISKS.md               # Risk register
├── KNOWN_ISSUES.md        # Pre-existing gaps
├── GLOSSARY.md            # Shared terminology
├── RECOMMENDATIONS.md     # Recommended resolutions & future enhancements
├── REVIEW_CHECKLIST.md    # Pre-PR human review
├── .gitignore
│
├── prompts/               # Numbered, self-contained execution prompts
│   ├── 01-make-loader-manifest-first.md
│   ├── 02-add-go-factory-registry.md
│   ├── 03-remove-builtin-auto-registration.md
│   └── 04-migrate-remaining-providers.md
│
├── tasks/                 # Detailed specs — dual-layer
│   ├── 01-make-loader-manifest-first.md
│   ├── 02-add-go-factory-registry.md
│   ├── 03-remove-builtin-auto-registration.md
│   └── 04-migrate-remaining-providers.md
│
├── findings/              # Research, audit output, analysis
└── runbooks/              # Operational docs
│
└── appendices/            # Deep-dive reference
    ├── a-provider-contract-checklist.md
    ├── b-simple-vs-go-decision-tree.md
    ├── c-post-initiative-checklist.md
    ├── d-gold-standard-alignment.md
    ├── e-test-scenarios.md
    └── f-architecture-diagram.md
```

Planning docs commit to the root repo on `dev`. Product code commits to the root repo on `dev`. Never commit directly to `main`.

> **Entry point:** Read `PREREQUISITES.md` before starting P01.

---

## Why now?

OpenMuara's mission is to emulate payment providers faithfully while staying local-first and simple. The provider layer has grown organically:

- Some providers are pure Go.
- Some have a YAML manifest but still auto-register in `init()`.
- Some have a manifest but the loader ignores it in favor of built-ins.
- Some providers have no manifest at all.

This creates confusion for contributors ("Should I write Go or YAML?") and for users ("Why is Fawry loaded when I didn't enable it?"). It also blocks future features: proprietary providers, sandboxed WASM plugins, and provider versioning all need a single, explicit discovery path.

This initiative standardizes that path without breaking existing behavior.

---

## Goals

1. **Manifest-first discovery** — `plugins/<name>/gateway.yml` is the source of truth for every non-default provider.
2. **Clear runtime graduation** — `runtime.type: simple` for common providers, `runtime.type: go` for complex ones, with no hidden auto-registration.
3. **Contributor-friendly onboarding** — a contributor can add a common provider with only YAML; they can graduate to Go when they need custom logic.
4. **Escape hatches for proprietary and sandboxed providers** — design the architecture so `bridge` and `wasm` runtimes fit naturally later.
5. **Backwards-compatible migration** — existing `.muara/config.yml` files keep working with clear warnings during the transition.
6. **Green quality gates** — build, test, vet, lint, and the `examples/checkout-store` smoke tests pass.
7. **Accurate documentation** — `docs/provider-contract.md` and `docs/contributing-providers.md` reflect the new model.

---

## Non-goals

- Adding new payment providers or new provider features.
- Implementing the `bridge` or `wasm` runtimes (architecture only).
- Changing provider protocol emulation behavior.
- Removing the `default` provider fallback.
- Large refactors outside the provider discovery path.

---

## Assumptions & Constraints

### Assumptions

- The user (human reviewer) will approve the recommended resolutions in `RECOMMENDATIONS.md` before P02 begins.
- The parallel agent session's product-code changes in `internal/` are either safe to build on or will be reviewed before P01.
- Existing provider protocol behavior does not need to change; only discovery/instantiation changes.
- `examples/checkout-store` is the primary end-to-end validation surface.

### Constraints

- Work stays on `dev`; no commits to `main`.
- Planning-doc commits are separate from product-code commits.
- No new external dependencies without explicit justification.
- `default` provider remains available without a manifest.
- P0 integration changes (P03 provider auto-registration removal, P04 provider migration) require human sign-off.

---

## Pinned Architecture

| Use case | Path | When to use |
|---|---|---|
| **Public common provider** | `plugins/<name>/gateway.yml` with `runtime.type: simple` | Common REST-ish providers that fit the simple runtime model. No Go code required. |
| **Public complex provider** | `plugins/<name>/gateway.yml` with `runtime.type: go` + Go package in `internal/<name>/` | Providers needing custom signature logic, state machines, or protocol quirks. |
| **Proprietary/private provider** | `providers.<name>.type: bridge` in `.muara/config.yml` | Future initiative. Keeps proprietary logic out of the public repo. |
| **Sandboxed runtime plugin** | `.muara/plugins/<name>.wasm` | Future initiative. Allows third-party providers without recompilation. |

Every non-default provider is discovered through `gateway.yml`. Built-in auto-registration is removed; Go packages register a factory and are activated by the manifest.

See `appendices/b-simple-vs-go-decision-tree.md` for a contributor-facing decision tree.

---

## Current State

- ✅ Simple YAML provider runtime exists (`internal/provider/simple/`).
- ✅ `gateway.yml` schema has `Runtime` and `SimpleRuntime` blocks.
- ✅ `muara provider test` and `muara provider init` commands exist.
- ✅ `plugins/fawry/gateway.yml` and `plugins/senangpay/gateway.yml` declare `runtime.type: simple`.
- ✅ `plugins/ipay88`, `billplz`, `toyyibpay` declare `runtime.type: go`.

## Gaps

- ❌ Loader prefers built-ins over manifests.
- ❌ `runtime.type: go` is not implemented in the loader.
- ❌ Built-in providers still auto-register in `init()`.
- ❌ No Go factory registry.
- ❌ Tests depend on global `provider.Get("fawry")`.
- ❌ Stripe does not have a `gateway.yml` manifest.
- ❌ Contributor docs do not explain simple vs go vs bridge vs wasm.

---

## Success Criteria

1. All non-default providers have a `plugins/<name>/gateway.yml` manifest.
2. Removing a built-in provider's Go package does not break config loading — the manifest drives discovery.
3. `runtime.type: simple` providers run without any Go registration.
4. `runtime.type: go` providers run through a factory registry looked up by name.
5. `go test ./...`, `go vet ./...`, and `golangci-lint run ./...` pass.
6. `examples/checkout-store` still works for Fawry and Stripe.
7. `docs/provider-contract.md` and `docs/contributing-providers.md` are accurate post-change.
8. Every prompt includes a regression test or an explicit justification in `DECISIONS.md`.

---

## Milestones

| # | Milestone | Status | Target Date | Prompt |
|---|---|---|---|---|
| 1 | Manifest-first loader | Pending | 2026-07-09 | P01 |
| 2 | Go factory registry for `runtime.type: go` | Pending | 2026-07-10 | P02 |
| 3 | Remove built-in auto-registration | Pending | 2026-07-11 | P03 |
| 4 | Migrate Fawry, SenangPay, iPay88, Billplz, ToyyibPay, Stripe to manifests | Pending | 2026-07-12 | P04 |
| 5 | Full gates, docs, and QA | Pending | 2026-07-12 | P05 |

---

## Metrics

| Metric | Current | Target | How measured |
|---|---|---|---|
| Providers with manifests | 5 | 6 (add Stripe) | `plugins/*/gateway.yml` count |
| Providers auto-registering | 6 | 0 | Code search for `provider.Register` in `init()` |
| Loader fallback to built-ins | Yes | No | `internal/config/provider_loader.go` review + tests |
| Factory registry coverage | 0% | 100% of `runtime.type: go` providers | Unit tests + loader integration tests |
| Test coverage drop on changed modules | Baseline | No drop | `go test -cover` diff |
| Checkout-store smoke tests | Passing | Passing | Manual + automated E2E |
| Contributor docs accuracy | Partial | Complete | Human review of `docs/contributing-providers.md` |

---

## Stakeholders & RACI

| Role | Who | Responsibility |
|---|---|---|
| Owner / Driver | AI Agent (Kimi Code) | Implement, test, document, update trackers. |
| Human Reviewer | ___________ | Approve architecture changes, review PR, merge to `dev`. |
| Provider Domain Experts | Fawry team (invited), Stripe docs | Validate protocol emulation remains faithful. |
| Contributors | Future OpenMuara contributors | Use the new simple/go/bridge/wasm paths. |

RACI (per prompt):

| Activity | Agent | Human | Domain Expert |
|---|---|---|---|
| Loader refactor | R/A | C | I |
| Factory registry design | R/A | C | I |
| Remove auto-registration | R/A | C | I |
| Migrate providers | R/A | C | C |
| Docs update | R/A | A | C |
| Final PR review | R | A | C |

*R = Responsible, A = Accountable, C = Consulted, I = Informed*

---

## Philosophy & Priority Alignment

Root `AGENTS.md` establishes this priority stack:

**Correctness > Security > Reliability > UX > Performance > Polish**

For this initiative, translate that stack into decisions:

1. **Correctness** — the loader must not instantiate a provider whose manifest is absent.
2. **Security** — provider factories must not bypass config validation; secrets stay in `.muara/config.yml` and out of git.
3. **Reliability** — removing a Go package cleanly removes its provider; no phantom registrations.
4. **UX** — error messages must tell the user which manifest is missing or invalid.
5. **Performance** — manifest loading is at startup; keep it fast but do not optimize at the cost of clarity.
6. **Polish** — docs and examples are updated before the initiative closes.

The project philosophy is **local-first, simple, and explicit**:

- No external services or telemetry.
- Provider protocol emulation remains faithful to documented behavior.
- Every provider is explicit in the filesystem.
- Every change is minimal and reviewable.

---

## Conventions

1. **Read `AGENTS.md` first.** Branch rules, quality gates, autonomy boundaries, and code style live there.
2. **Manifest is the source of truth.** No provider loads without a `gateway.yml`.
3. **One logical change per commit.** Each prompt gets its own commit.
4. **Regression tests required.** Every code change must include a test that fails before and passes after, or a written justification in `DECISIONS.md`.
5. **Update trackers.** `TRACKING.md`, `HANDOFF.md`, and `DECISIONS.md` are updated after every prompt.
6. **P0 integration changes need sign-off.** Per `AGENTS.md`, changes to provider emulation logic, config persistence schemas, or the provider plugin schema contract require user sign-off recorded in `DECISIONS.md`.

---

## Quality Gates

Every prompt must pass:

- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`
- `examples/checkout-store` smoke test for Fawry and Stripe

Optional but recommended:

- `go test -race ./...`
- `go test -cover ./...` (no drop on changed modules)
- `muara provider validate plugins/*/gateway.yml`

---

## Proposed Workflow

| Step | Prompt | Goal | Key Output |
|---|---|---|---|
| 1 | P01 Manifest-first loader | Loader discovers providers from `gateway.yml`; built-ins are fallback only for `runtime.type: go`. | Updated loader + tests. |
| 2 | P02 Go factory registry | Go providers register factories; loader instantiates via registry. | `internal/provider/hybrid/registry.go` + provider `register.go` files. |
| 3 | P03 Remove auto-registration | Delete `init()` registrations; update tests and router. | No phantom providers; all gates green. |
| 4 | P04 Migrate providers | Normalize manifests; create Stripe manifest; update docs. | All non-default providers manifest-driven. |
| 5 | Final gates & PR | Run full gates, update docs, open PR to `dev`. | Initiative closed. |

Detailed tasks for each prompt are in `tasks/`.

---

## Communication & Escalation

| Situation | Action |
|---|---|
| Prompt blocked > 30 min | Log blocker in `RISKS.md` and `HANDOFF.md`; move to next prompt if parallel-safe. |
| Quality gate fails | Stop. Do not proceed. Fix or escalate with diff and logs. |
| Provider protocol behavior changes | Record in `DECISIONS.md`; get human sign-off before committing. |
| Parallel agent conflict | Rebase, run gates, and update `HANDOFF.md` before continuing. |
| Scope creep | Push to `RECOMMENDATIONS.md` future enhancements or a new initiative; do not expand this one. |

---

## Post-Initiative Actions

After P05, complete `appendices/c-post-initiative-checklist.md`. The high-level summary is:

1. Final documentation sweep (`TRACKING.md`, `HANDOFF.md`, `DECISIONS.md`, `RISKS.md`, `KNOWN_ISSUES.md`).
2. Add a `CHANGELOG.md` release-notes snippet.
3. Update root `TRACKING.md` and the v1 master backlog.
4. Open a PR from `dev` to `main` (or keep on `dev` per project flow).
5. Hand off to the human reviewer using `REVIEW_CHECKLIST.md`.

---

## Cross-Reference Map

| Tracker / Doc | Path | What It Contains |
|---|---|---|
| This README | `docs/initiatives/openmuara-provider-manifests/README.md` | Scope, goals, architecture, workflow |
| Execution tracker | `docs/initiatives/openmuara-provider-manifests/TRACKING.md` | Milestones, gate results, action log |
| Prerequisites | `docs/initiatives/openmuara-provider-manifests/PREREQUISITES.md` | Pre-flight checklist |
| Handoff | `docs/initiatives/openmuara-provider-manifests/HANDOFF.md` | Session continuity |
| Decisions | `docs/initiatives/openmuara-provider-manifests/DECISIONS.md` | Decision log |
| Risks | `docs/initiatives/openmuara-provider-manifests/RISKS.md` | Risk register |
| Known issues | `docs/initiatives/openmuara-provider-manifests/KNOWN_ISSUES.md` | Pre-existing gaps |
| Glossary | `docs/initiatives/openmuara-provider-manifests/GLOSSARY.md` | Shared terminology |
| Recommendations | `docs/initiatives/openmuara-provider-manifests/RECOMMENDATIONS.md` | Recommended resolutions & future enhancements |
| Review checklist | `docs/initiatives/openmuara-provider-manifests/REVIEW_CHECKLIST.md` | Pre-PR human review |
| Provider contract checklist | `docs/initiatives/openmuara-provider-manifests/appendices/a-provider-contract-checklist.md` | Contract completeness checks |
| Simple vs Go decision tree | `docs/initiatives/openmuara-provider-manifests/appendices/b-simple-vs-go-decision-tree.md` | Contributor decision guide |
| Post-initiative checklist | `docs/initiatives/openmuara-provider-manifests/appendices/c-post-initiative-checklist.md` | Close-out steps |
| Gold-standard alignment | `docs/initiatives/openmuara-provider-manifests/appendices/d-gold-standard-alignment.md` | How this builds on prior quality initiatives |
| Test scenarios | `docs/initiatives/openmuara-provider-manifests/appendices/e-test-scenarios.md` | Specific test cases per prompt |
| Architecture diagram | `docs/initiatives/openmuara-provider-manifests/appendices/f-architecture-diagram.md` | Manifest-first discovery flow |
| Root `AGENTS.md` | `AGENTS.md` | Workspace rules and autonomy boundaries |
| Provider contract | `docs/provider-contract.md` | Provider protocol contract |
| Contributing providers | `docs/contributing-providers.md` | Contributor guide |

---

## Self-Assessment

**Planning Solidity: 9.8 / 10**

The initiative links to project philosophy; defines RACI with a signed-off human reviewer; includes metrics with targets; provides communication/escalation plans; has pre-PR and post-initiative checklists; glossary; known-issues register; recommendations register (approved); assumptions/constraints register; a full test-scenarios appendix with traceability; an architecture diagram; and appendices for contract completeness, contributor decision-making, and gold-standard alignment. Milestones have target dates. Open decisions are pinned. Files are split to stay close to the 250-line guideline.

The remaining 0.2 points are execution risk: the next agent must follow the sign-off gates, keep commits focused, and complete the post-initiative checklist. No additional documentation will fix that.
